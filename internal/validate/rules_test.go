package validate

import (
	"testing"

	"github.com/badimirzai/robotics-verifier-cli/internal/model"
)

func baseSpec() model.RobotSpec {
	return model.RobotSpec{
		Power: model.PowerSpec{
			Battery: model.Battery{
				VoltageV: 12,
			},
			Rail: model.Rail{
				VoltageV:    5,
				MaxCurrentA: 1,
			},
		},
		MCU: model.MCU{
			LogicVoltageV: 5,
		},
		Driver: model.MotorDriver{
			Channels:         2,
			MotorSupplyMinV:  6,
			MotorSupplyMaxV:  16,
			LogicVoltageMinV: 4.5,
			LogicVoltageMaxV: 5.5,
			ContinuousPerChA: 2,
			PeakPerChA:       6,
		},
		Motors: []model.Motor{
			{
				Name:            "M",
				Count:           2,
				NominalCurrentA: 1,
				StallCurrentA:   5,
			},
		},
	}
}

func reportCodes(r Report) map[string]bool {
	codes := make(map[string]bool)
	for _, f := range r.Findings {
		codes[f.Code] = true
	}
	return codes
}

func requireHasCode(t *testing.T, codes map[string]bool, code string) {
	t.Helper()
	if !codes[code] {
		t.Fatalf("expected code %q, got %#v", code, codes)
	}
}

func requireNoCode(t *testing.T, codes map[string]bool, code string) {
	t.Helper()
	if codes[code] {
		t.Fatalf("did not expect code %q, got %#v", code, codes)
	}
}

func TestRuleDriverChannels(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*model.RobotSpec)
		want   []string
		not    []string
	}{
		{
			name: "invalid_channels",
			mutate: func(s *model.RobotSpec) {
				s.Driver.Channels = 0
			},
			want: []string{"DRV_CHANNELS_INVALID"},
			not:  []string{"DRV_CHANNELS_OK", "DRV_CHANNELS_INSUFFICIENT"},
		},
		{
			name: "insufficient_channels",
			mutate: func(s *model.RobotSpec) {
				s.Driver.Channels = 2
				s.Motors[0].Count = 3
			},
			want: []string{"DRV_CHANNELS_INSUFFICIENT"},
			not:  []string{"DRV_CHANNELS_OK", "DRV_CHANNELS_INVALID"},
		},
		{
			name:   "ok_channels",
			mutate: func(s *model.RobotSpec) {},
			want:   []string{"DRV_CHANNELS_OK"},
			not:    []string{"DRV_CHANNELS_INVALID", "DRV_CHANNELS_INSUFFICIENT"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec := baseSpec()
			tt.mutate(&spec)
			codes := reportCodes(RunAll(spec, nil))
			for _, c := range tt.want {
				requireHasCode(t, codes, c)
			}
			for _, c := range tt.not {
				requireNoCode(t, codes, c)
			}
		})
	}
}

func TestRuleMotorSupplyVoltage(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*model.RobotSpec)
		want   []string
		not    []string
	}{
		{
			name: "battery_invalid",
			mutate: func(s *model.RobotSpec) {
				s.Power.Battery.VoltageV = 0
			},
			want: []string{"BAT_V_INVALID"},
			not:  []string{"DRV_SUPPLY_RANGE"},
		},
		{
			name: "battery_out_of_range",
			mutate: func(s *model.RobotSpec) {
				s.Power.Battery.VoltageV = 20
			},
			want: []string{"DRV_SUPPLY_RANGE"},
			not:  []string{"BAT_V_INVALID"},
		},
		{
			name:   "battery_in_range",
			mutate: func(s *model.RobotSpec) {},
			not:    []string{"BAT_V_INVALID", "DRV_SUPPLY_RANGE"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec := baseSpec()
			tt.mutate(&spec)
			codes := reportCodes(RunAll(spec, nil))
			for _, c := range tt.want {
				requireHasCode(t, codes, c)
			}
			for _, c := range tt.not {
				requireNoCode(t, codes, c)
			}
		})
	}
}

func TestRuleDriverCurrentHeadroom(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*model.RobotSpec)
		want   []string
		not    []string
	}{
		{
			name: "motor_count_invalid",
			mutate: func(s *model.RobotSpec) {
				s.Motors[0].Count = 0
			},
			want: []string{"MOTOR_COUNT_INVALID"},
			not:  []string{"DRV_PEAK_LT_STALL", "DRV_CONT_LOW_MARGIN"},
		},
		{
			name: "peak_less_than_stall",
			mutate: func(s *model.RobotSpec) {
				s.Driver.PeakPerChA = 3
			},
			want: []string{"DRV_PEAK_LT_STALL"},
			not:  []string{"DRV_CONT_LOW_MARGIN"},
		},
		{
			name: "continuous_low_margin",
			mutate: func(s *model.RobotSpec) {
				s.Driver.ContinuousPerChA = 1.0
			},
			want: []string{"DRV_CONT_LOW_MARGIN"},
			not:  []string{"DRV_PEAK_LT_STALL"},
		},
		{
			name: "peak_and_continuous_issues",
			mutate: func(s *model.RobotSpec) {
				s.Driver.PeakPerChA = 3
				s.Driver.ContinuousPerChA = 1.0
			},
			want: []string{"DRV_PEAK_LT_STALL", "DRV_CONT_LOW_MARGIN"},
			not:  []string{},
		},
		{
			name:   "headroom_ok",
			mutate: func(s *model.RobotSpec) {},
			not:    []string{"MOTOR_COUNT_INVALID", "DRV_PEAK_LT_STALL", "DRV_CONT_LOW_MARGIN"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec := baseSpec()
			tt.mutate(&spec)
			codes := reportCodes(RunAll(spec, nil))
			for _, c := range tt.want {
				requireHasCode(t, codes, c)
			}
			for _, c := range tt.not {
				requireNoCode(t, codes, c)
			}
		})
	}
}

func TestRuleLogicVoltageCompat(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*model.RobotSpec)
		want   []string
		not    []string
	}{
		{
			name: "rail_invalid",
			mutate: func(s *model.RobotSpec) {
				s.Power.Rail.VoltageV = 0
			},
			want: []string{"RAIL_V_INVALID"},
			not:  []string{"LOGIC_V_DRIVER_MISMATCH", "LOGIC_V_MCU_MISMATCH"},
		},
		{
			name: "rail_outside_driver_range",
			mutate: func(s *model.RobotSpec) {
				s.Power.Rail.VoltageV = 7
			},
			want: []string{"LOGIC_V_DRIVER_MISMATCH"},
			not:  []string{"RAIL_V_INVALID", "LOGIC_V_MCU_MISMATCH"},
		},
		{
			name: "mcu_mismatch",
			mutate: func(s *model.RobotSpec) {
				s.MCU.LogicVoltageV = 3.3
			},
			want: []string{"LOGIC_V_MCU_MISMATCH"},
			not:  []string{"RAIL_V_INVALID", "LOGIC_V_DRIVER_MISMATCH"},
		},
		{
			name:   "logic_ok",
			mutate: func(s *model.RobotSpec) {},
			not:    []string{"RAIL_V_INVALID", "LOGIC_V_DRIVER_MISMATCH", "LOGIC_V_MCU_MISMATCH"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec := baseSpec()
			tt.mutate(&spec)
			codes := reportCodes(RunAll(spec, nil))
			for _, c := range tt.want {
				requireHasCode(t, codes, c)
			}
			for _, c := range tt.not {
				requireNoCode(t, codes, c)
			}
		})
	}
}

func TestRuleRailCurrentBudget(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*model.RobotSpec)
		want   []string
		not    []string
	}{
		{
			name: "rail_current_unknown",
			mutate: func(s *model.RobotSpec) {
				s.Power.Rail.MaxCurrentA = 0
			},
			want: []string{"RAIL_I_UNKNOWN"},
			not:  []string{"RAIL_BUDGET_NOTE"},
		},
		{
			name:   "rail_current_set",
			mutate: func(s *model.RobotSpec) {},
			want:   []string{"RAIL_BUDGET_NOTE"},
			not:    []string{"RAIL_I_UNKNOWN"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec := baseSpec()
			tt.mutate(&spec)
			codes := reportCodes(RunAll(spec, nil))
			for _, c := range tt.want {
				requireHasCode(t, codes, c)
			}
			for _, c := range tt.not {
				requireNoCode(t, codes, c)
			}
		})
	}
}
