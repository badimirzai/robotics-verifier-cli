package validate

import (
	"fmt"
	"math"
	"strings"

	"github.com/badimirzai/robotics-verifier-cli/internal/model"
)

type Severity string

const (
	SevError Severity = "ERROR"
	SevWarn  Severity = "WARN"
	SevInfo  Severity = "INFO"
)

type Finding struct {
	Severity Severity
	Code     string
	Message  string
	Location *Location
}

type Report struct {
	Findings []Finding
}

type Location struct {
	File string
	Line int
}

func (r Report) HasErrors() bool {
	for _, f := range r.Findings {
		if f.Severity == SevError {
			return true
		}
	}
	return false
}

func RunAll(spec model.RobotSpec, locs map[string]Location) Report {
	var r Report
	r.Findings = append(r.Findings, ruleDriverChannels(spec, locs)...)
	r.Findings = append(r.Findings, ruleMotorSupplyVoltage(spec, locs)...)
	r.Findings = append(r.Findings, ruleDriverCurrentHeadroom(spec, locs)...)
	r.Findings = append(r.Findings, ruleLogicVoltageCompat(spec, locs)...)
	r.Findings = append(r.Findings, ruleRailCurrentBudget(spec, locs)...)
	return r
}

func withLocation(locs map[string]Location, path string, f Finding) Finding {
	if locs == nil {
		return f
	}
	if loc, ok := findLocation(locs, path); ok {
		f.Location = &loc
	}
	return f
}

func findLocation(locs map[string]Location, path string) (Location, bool) {
	if loc, ok := locs[path]; ok {
		return loc, true
	}
	for path != "" {
		if idx := strings.LastIndex(path, "."); idx >= 0 {
			path = path[:idx]
		} else if idx := strings.LastIndex(path, "["); idx >= 0 {
			path = path[:idx]
		} else {
			path = ""
		}
		if loc, ok := locs[path]; ok {
			return loc, true
		}
	}
	return Location{}, false
}

func ruleDriverChannels(spec model.RobotSpec, locs map[string]Location) []Finding {
	totalMotors := 0
	for _, m := range spec.Motors {
		totalMotors += m.Count
	}
	if spec.Driver.Channels <= 0 {
		return []Finding{withLocation(locs, "motor_driver.channels", Finding{
			SevError,
			"DRV_CHANNELS_INVALID",
			"motor_driver.channels must be > 0",
			nil,
		})}
	}
	if totalMotors > spec.Driver.Channels {
		return []Finding{withLocation(locs, "motor_driver.channels", Finding{
			SevError,
			"DRV_CHANNELS_INSUFFICIENT",
			fmt.Sprintf("motors require %d channels but motor_driver.channels is %d", totalMotors, spec.Driver.Channels),
			nil,
		})}
	}
	return []Finding{withLocation(locs, "motor_driver.channels", Finding{
		SevInfo,
		"DRV_CHANNELS_OK",
		fmt.Sprintf("channels OK: %d motors <= %d motor_driver.channels", totalMotors, spec.Driver.Channels),
		nil,
	})}
}

func ruleMotorSupplyVoltage(spec model.RobotSpec, locs map[string]Location) []Finding {
	batV := spec.Power.Battery.VoltageV
	if batV <= 0 {
		return []Finding{withLocation(locs, "power.battery.voltage_v", Finding{
			SevError,
			"BAT_V_INVALID",
			"power.battery.voltage_v must be > 0",
			nil,
		})}
	}
	if batV < spec.Driver.MotorSupplyMinV || batV > spec.Driver.MotorSupplyMaxV {
		return []Finding{withLocation(locs, "power.battery.voltage_v", Finding{
			SevError,
			"DRV_SUPPLY_RANGE",
			fmt.Sprintf("battery %.2fV outside motor_driver motor supply range [%.2f, %.2f]V",
				batV, spec.Driver.MotorSupplyMinV, spec.Driver.MotorSupplyMaxV),
			nil,
		})}
	}
	return nil
}

func ruleDriverCurrentHeadroom(spec model.RobotSpec, locs map[string]Location) []Finding {
	var out []Finding
	for i, m := range spec.Motors {
		if m.Count <= 0 {
			path := fmt.Sprintf("motors[%d].count", i)
			out = append(out, withLocation(locs, path, Finding{
				SevError,
				"MOTOR_COUNT_INVALID",
				fmt.Sprintf("motors[%d].count must be > 0", i),
				nil,
			}))
			continue
		}
		// Worst-case per channel: stall current. If you want to be conservative, require peak >= stall.
		if spec.Driver.PeakPerChA < m.StallCurrentA {
			out = append(out, withLocation(locs, "motor_driver.peak_per_channel_a", Finding{
				SevError,
				"DRV_PEAK_LT_STALL",
				fmt.Sprintf("motor_driver.peak_per_channel_a %.2fA < motor %s stall %.2fA (per channel)",
					spec.Driver.PeakPerChA, m.Name, m.StallCurrentA),
				nil,
			}))
		}
		// Continuous should exceed nominal with margin
		margin := 1.25
		if spec.Driver.ContinuousPerChA < margin*m.NominalCurrentA {
			out = append(out, withLocation(locs, "motor_driver.continuous_per_channel_a", Finding{
				SevWarn,
				"DRV_CONT_LOW_MARGIN",
				fmt.Sprintf("motor_driver.continuous_per_channel_a %.2fA may be low for motor %s nominal %.2fA (want >= %.2fA)",
					spec.Driver.ContinuousPerChA, m.Name, m.NominalCurrentA, margin*m.NominalCurrentA),
				nil,
			}))
		}
	}
	return out
}

func ruleLogicVoltageCompat(spec model.RobotSpec, locs map[string]Location) []Finding {
	lv := spec.Power.Rail.VoltageV
	if lv <= 0 {
		return []Finding{withLocation(locs, "power.logic_rail.voltage_v", Finding{
			SevError,
			"RAIL_V_INVALID",
			"power.logic_rail.voltage_v must be > 0",
			nil,
		})}
	}
	if lv < spec.Driver.LogicVoltageMinV || lv > spec.Driver.LogicVoltageMaxV {
		return []Finding{withLocation(locs, "power.logic_rail.voltage_v", Finding{
			SevError,
			"LOGIC_V_DRIVER_MISMATCH",
			fmt.Sprintf("logic rail %.2fV outside motor_driver logic range [%.2f, %.2f]V",
				lv, spec.Driver.LogicVoltageMinV, spec.Driver.LogicVoltageMaxV),
			nil,
		})}
	}
	if math.Abs(spec.MCU.LogicVoltageV-lv) > 0.25 {
		return []Finding{withLocation(locs, "mcu.logic_voltage_v", Finding{
			SevWarn,
			"LOGIC_V_MCU_MISMATCH",
			fmt.Sprintf("MCU logic %.2fV differs from rail %.2fV, check level shifting",
				spec.MCU.LogicVoltageV, lv),
			nil,
		})}
	}
	return nil
}

func ruleRailCurrentBudget(spec model.RobotSpec, locs map[string]Location) []Finding {
	railMax := spec.Power.Rail.MaxCurrentA
	if railMax <= 0 {
		return []Finding{withLocation(locs, "power.logic_rail.max_current_a", Finding{
			SevWarn,
			"RAIL_I_UNKNOWN",
			"power.logic_rail.max_current_a not set, cannot budget logic rail current",
			nil,
		})}
	}
	// For v1 we do not model currents precisely. We just nudge the user.
	return []Finding{withLocation(locs, "power.logic_rail.max_current_a", Finding{
		SevInfo,
		"RAIL_BUDGET_NOTE",
		fmt.Sprintf("logic rail budget set to %.2fA (v1 does not estimate MCU+driver logic draw yet)", railMax),
		nil,
	})}
}
