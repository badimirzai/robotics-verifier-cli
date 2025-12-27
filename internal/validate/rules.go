package validate

import (
	"fmt"
	"math"

	"github.com/badimirzai/architon-cli/internal/model"
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
}

type Report struct {
	Findings []Finding
}

func (r Report) HasErrors() bool {
	for _, f := range r.Findings {
		if f.Severity == SevError {
			return true
		}
	}
	return false
}

func RunAll(spec model.RobotSpec) Report {
	var r Report
	r.Findings = append(r.Findings, ruleDriverChannels(spec)...)
	r.Findings = append(r.Findings, ruleMotorSupplyVoltage(spec)...)
	r.Findings = append(r.Findings, ruleDriverCurrentHeadroom(spec)...)
	r.Findings = append(r.Findings, ruleLogicVoltageCompat(spec)...)
	r.Findings = append(r.Findings, ruleRailCurrentBudget(spec)...)
	return r
}

func ruleDriverChannels(spec model.RobotSpec) []Finding {
	totalMotors := 0
	for _, m := range spec.Motors {
		totalMotors += m.Count
	}
	if spec.Driver.Channels <= 0 {
		return []Finding{{SevError, "DRV_CHANNELS_INVALID", "driver.channels must be > 0"}}
	}
	if totalMotors > spec.Driver.Channels {
		return []Finding{{
			SevError,
			"DRV_CHANNELS_INSUFFICIENT",
			fmt.Sprintf("motors require %d channels but driver has %d", totalMotors, spec.Driver.Channels),
		}}
	}
	return []Finding{{
		SevInfo,
		"DRV_CHANNELS_OK",
		fmt.Sprintf("channels OK: %d motors <= %d driver channels", totalMotors, spec.Driver.Channels),
	}}
}

func ruleMotorSupplyVoltage(spec model.RobotSpec) []Finding {
	batV := spec.Power.Battery.VoltageV
	if batV <= 0 {
		return []Finding{{SevError, "BAT_V_INVALID", "battery.voltage_v must be > 0"}}
	}
	if batV < spec.Driver.MotorSupplyMinV || batV > spec.Driver.MotorSupplyMaxV {
		return []Finding{{
			SevError,
			"DRV_SUPPLY_RANGE",
			fmt.Sprintf("battery %.2fV outside driver motor supply range [%.2f, %.2f]V",
				batV, spec.Driver.MotorSupplyMinV, spec.Driver.MotorSupplyMaxV),
		}}
	}
	return nil
}

func ruleDriverCurrentHeadroom(spec model.RobotSpec) []Finding {
	var out []Finding
	for _, m := range spec.Motors {
		if m.Count <= 0 {
			out = append(out, Finding{SevError, "MOTOR_COUNT_INVALID", fmt.Sprintf("%s count must be > 0", m.Name)})
			continue
		}
		// Worst-case per channel: stall current. If you want to be conservative, require peak >= stall.
		if spec.Driver.PeakPerChA < m.StallCurrentA {
			out = append(out, Finding{
				SevError,
				"DRV_PEAK_LT_STALL",
				fmt.Sprintf("driver peak %.2fA < motor %s stall %.2fA (per channel)",
					spec.Driver.PeakPerChA, m.Name, m.StallCurrentA),
			})
		}
		// Continuous should exceed nominal with margin
		margin := 1.25
		if spec.Driver.ContinuousPerChA < margin*m.NominalCurrentA {
			out = append(out, Finding{
				SevWarn,
				"DRV_CONT_LOW_MARGIN",
				fmt.Sprintf("driver continuous %.2fA may be low for motor %s nominal %.2fA (want >= %.2fA)",
					spec.Driver.ContinuousPerChA, m.Name, m.NominalCurrentA, margin*m.NominalCurrentA),
			})
		}
	}
	return out
}

func ruleLogicVoltageCompat(spec model.RobotSpec) []Finding {
	lv := spec.Power.Rail.VoltageV
	if lv <= 0 {
		return []Finding{{SevError, "RAIL_V_INVALID", "power.rail.voltage_v must be > 0"}}
	}
	if lv < spec.Driver.LogicVoltageMinV || lv > spec.Driver.LogicVoltageMaxV {
		return []Finding{{
			SevError,
			"LOGIC_V_DRIVER_MISMATCH",
			fmt.Sprintf("logic rail %.2fV outside driver logic range [%.2f, %.2f]V",
				lv, spec.Driver.LogicVoltageMinV, spec.Driver.LogicVoltageMaxV),
		}}
	}
	if math.Abs(spec.MCU.LogicVoltageV-lv) > 0.25 {
		return []Finding{{
			SevWarn,
			"LOGIC_V_MCU_MISMATCH",
			fmt.Sprintf("MCU logic %.2fV differs from rail %.2fV, check level shifting",
				spec.MCU.LogicVoltageV, lv),
		}}
	}
	return nil
}

func ruleRailCurrentBudget(spec model.RobotSpec) []Finding {
	railMax := spec.Power.Rail.MaxCurrentA
	if railMax <= 0 {
		return []Finding{{SevWarn, "RAIL_I_UNKNOWN", "power.rail.max_current_a not set, cannot budget logic rail current"}}
	}
	// For v1 we do not model currents precisely. We just nudge the user.
	return []Finding{{
		SevInfo,
		"RAIL_BUDGET_NOTE",
		fmt.Sprintf("logic rail budget set to %.2fA (v1 does not estimate MCU+driver logic draw yet)", railMax),
	}}
}
