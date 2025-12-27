package parts

import (
	"fmt"
	"os"
	"path/filepath"

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
		Channels          int     `yaml:"channels"`
		MotorVoltageMin   float64 `yaml:"motor_voltage_min"`
		MotorVoltageMax   float64 `yaml:"motor_voltage_max"`
		LogicVoltageMin   float64 `yaml:"logic_voltage_min"`
		LogicVoltageMax   float64 `yaml:"logic_voltage_max"`
		CurrentContinuous float64 `yaml:"current_continuous"`
		CurrentPeak       float64 `yaml:"current_peak"`
	} `yaml:"motor_driver"`
}

// MotorPartFile represents the YAML structure for a motor.
// Example: parts/motors/generic_dc_12v_gearmotor.yaml
type MotorPartFile struct {
	PartID string `yaml:"part_id"`
	Type   string `yaml:"type"`
	Name   string `yaml:"name"`

	Motor struct {
		NominalCurrent float64 `yaml:"nominal_current"`
		StallCurrent   float64 `yaml:"stall_current"`
	} `yaml:"motor"`
}

// MCUPartFile represents the YAML structure for an MCU.
// Example: parts/mcus/esp32s3.yaml
type MCUPartFile struct {
	PartID string `yaml:"part_id"`
	Type   string `yaml:"type"`
	Name   string `yaml:"name"`

	MCU struct {
		LogicVoltage float64 `yaml:"logic_voltage"`
	} `yaml:"mcu"`
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
