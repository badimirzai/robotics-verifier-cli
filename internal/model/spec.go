package model

import (
	"fmt"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

type RobotSpec struct {
	Name     string      `yaml:"name"`
	Power    PowerSpec   `yaml:"power"`
	Motors   []Motor     `yaml:"motors"`
	Driver   MotorDriver `yaml:"motor_driver"`
	MCU      MCU         `yaml:"mcu"`
	I2CBuses []I2CBus    `yaml:"i2c_buses"`
}

type PowerSpec struct {
	Battery Battery `yaml:"battery"`
	Rail    Rail    `yaml:"logic_rail"` // main logic rail after regulation
}

type Battery struct {
	Chemistry     string  `yaml:"chemistry"` // e.g. "Li-ion"
	VoltageV      float64 `yaml:"voltage_v"` // nominal
	MaxCurrentA   float64 `yaml:"max_current_a"`
	CapacityAh    float64 `yaml:"capacity_ah"`
	CRating       float64 `yaml:"c_rating"`
	MaxDischargeA float64 `yaml:"max_discharge_a"`
}

type Rail struct {
	VoltageV    float64 `yaml:"voltage_v"`     // e.g. 5.0
	MaxCurrentA float64 `yaml:"max_current_a"` // regulator output capability
}

type Motor struct {
	Part            string  `yaml:"part,omitempty"`
	Name            string  `yaml:"name"`
	Count           int     `yaml:"count"`
	VoltageMinV     float64 `yaml:"voltage_min_v"`
	VoltageMaxV     float64 `yaml:"voltage_max_v"`
	StallCurrentA   float64 `yaml:"stall_current_a"`
	NominalCurrentA float64 `yaml:"nominal_current_a"`
}

type MotorDriver struct {
	Part             string  `yaml:"part,omitempty"`
	Name             string  `yaml:"name"`
	MotorSupplyMinV  float64 `yaml:"motor_supply_min_v"`
	MotorSupplyMaxV  float64 `yaml:"motor_supply_max_v"`
	ContinuousPerChA float64 `yaml:"continuous_per_channel_a"`
	PeakPerChA       float64 `yaml:"peak_per_channel_a"`
	Channels         int     `yaml:"channels"`
	LogicVoltageMinV float64 `yaml:"logic_voltage_min_v"`
	LogicVoltageMaxV float64 `yaml:"logic_voltage_max_v"`
}

type MCU struct {
	Part             string  `yaml:"part,omitempty"`
	Name             string  `yaml:"name"`
	LogicVoltageV    float64 `yaml:"logic_voltage_v"` // usually 3.3 for ESP32
	MaxGPIOCurrentmA float64 `yaml:"max_gpio_current_ma"`
}

type I2CBus struct {
	Name    string      `yaml:"name"`
	Devices []I2CDevice `yaml:"devices"`
}

type I2CAddress uint16

func (a *I2CAddress) UnmarshalYAML(value *yaml.Node) error {
	// Allow decimal or 0x-prefixed hex addresses during YAML decoding.
	if value.Kind != yaml.ScalarNode {
		return fmt.Errorf("i2c address must be a scalar")
	}
	raw := strings.TrimSpace(value.Value)
	if raw == "" {
		return fmt.Errorf("i2c address is empty")
	}
	parsed, err := strconv.ParseUint(raw, 0, 16)
	if err != nil {
		return fmt.Errorf("invalid i2c address %q: %w", raw, err)
	}
	*a = I2CAddress(parsed)
	return nil
}

type I2CDevice struct {
	Part       string     `yaml:"part,omitempty"`
	Name       string     `yaml:"name"`
	AddressHex I2CAddress `yaml:"address_hex"`
}
