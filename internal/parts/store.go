package parts

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/badimirzai/robotics-verifier-cli/internal/model"
	"gopkg.in/yaml.v3"
)

// DriverPartFile represents the YAML structure for a motor driver part.
// Example: parts/drivers/tb6612fng.yaml
type DriverPartFile struct {
	PartID string `yaml:"part_id"`
	Type   string `yaml:"type"`
	Name   string `yaml:"name"`
	MPN    string `yaml:"mpn"`

	MotorDriver struct {
		Channels         int     `yaml:"channels"`
		MotorSupplyMinV  float64 `yaml:"motor_supply_min_v"`
		MotorSupplyMaxV  float64 `yaml:"motor_supply_max_v"`
		LogicVoltageMinV float64 `yaml:"logic_voltage_min_v"`
		LogicVoltageMaxV float64 `yaml:"logic_voltage_max_v"`
		ContinuousPerChA float64 `yaml:"continuous_per_channel_a"`
		PeakPerChA       float64 `yaml:"peak_per_channel_a"`
	} `yaml:"motor_driver"`
}

// MotorPartFile represents the YAML structure for a motor.
// Example: parts/motors/generic_dc_12v_gearmotor.yaml
type MotorPartFile struct {
	PartID string `yaml:"part_id"`
	Type   string `yaml:"type"`
	Name   string `yaml:"name"`

	Motor struct {
		VoltageMinV     float64 `yaml:"voltage_min_v"`
		VoltageMaxV     float64 `yaml:"voltage_max_v"`
		NominalCurrentA float64 `yaml:"nominal_current_a"`
		StallCurrentA   float64 `yaml:"stall_current_a"`
	} `yaml:"motor"`
}

// MCUPartFile represents the YAML structure for an MCU.
// Example: parts/mcus/esp32s3.yaml
type MCUPartFile struct {
	PartID string `yaml:"part_id"`
	Type   string `yaml:"type"`
	Name   string `yaml:"name"`

	MCU struct {
		LogicVoltageV float64 `yaml:"logic_voltage_v"`
	} `yaml:"mcu"`
}

// I2CSensorPartFile represents the YAML structure for an I2C sensor device.
// Example: parts/sensors/mpu6050.yaml
type I2CSensorPartFile struct {
	PartID string `yaml:"part_id"`
	Type   string `yaml:"type"`
	Name   string `yaml:"name"`
	MPN    string `yaml:"mpn"`

	I2CDevice struct {
		AddressHex model.I2CAddress `yaml:"address_hex"`
	} `yaml:"i2c_device"`
}

// Store knows how to load part files from a base directory (usually "./parts").
type Store struct{ BaseDir string }

// NewStore creates a new part store rooted at baseDir (e.g. "parts").
func NewStore(baseDir string) *Store {
	return &Store{BaseDir: baseDir}
}

// LoadDriver loads a motor driver part by ID, e.g. "drivers/tb6612fng".
func (s *Store) LoadDriver(partID string) (DriverPartFile, error) {
	var part DriverPartFile
	if err := s.loadPart(partID, &part); err != nil {
		return DriverPartFile{}, err
	}
	if part.Type != "motor_driver" {
		return DriverPartFile{}, fmt.Errorf("expected type motor_driver, got %q", part.Type)
	}
	return part, nil
}

// LoadMotor loads a motor part by ID, e.g. "motors/generic_dc_12v_gearmotor".
func (s *Store) LoadMotor(partID string) (MotorPartFile, error) {
	var part MotorPartFile
	if err := s.loadPart(partID, &part); err != nil {
		return MotorPartFile{}, err
	}
	if part.Type != "motor" {
		return MotorPartFile{}, fmt.Errorf("expected type motor, got %q", part.Type)
	}
	return part, nil
}

// LoadMCU loads an MCU part by ID, e.g. "mcus/esp32s3".
func (s *Store) LoadMCU(partID string) (MCUPartFile, error) {
	var part MCUPartFile
	if err := s.loadPart(partID, &part); err != nil {
		return MCUPartFile{}, err
	}
	if part.Type != "mcu" {
		return MCUPartFile{}, fmt.Errorf("expected type mcu, got %q", part.Type)
	}
	return part, nil
}

// LoadI2CSensor loads an I2C sensor part by ID, e.g. "sensors/mpu6050".
func (s *Store) LoadI2CSensor(partID string) (I2CSensorPartFile, error) {
	var part I2CSensorPartFile
	if err := s.loadPart(partID, &part); err != nil {
		return I2CSensorPartFile{}, err
	}
	if part.Type != "i2c_sensor" {
		return I2CSensorPartFile{}, fmt.Errorf("expected type i2c_sensor, got %q", part.Type)
	}
	return part, nil
}

// loadPart is a small helper to read and unmarshal a YAML file.
func (s *Store) loadPart(partID string, out any) error {
	path := filepath.Join(s.BaseDir, filepath.FromSlash(partID)+".yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(data, out); err != nil {
		return err
	}
	return nil
}
