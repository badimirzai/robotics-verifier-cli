package parts

import (
	"path/filepath"
	"runtime"
	"testing"
)

func testPartsDir(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to locate test file path")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", "..", "parts"))
}

func TestStore_LoadDriver_TB6612FNG(t *testing.T) {
	store := NewStore(testPartsDir(t))

	drv, err := store.LoadDriver("drivers/tb6612fng")
	if err != nil {
		t.Fatalf("LoadDriver(tb6612fng) returned error: %v", err)
	}

	if drv.PartID != "drivers/tb6612fng" {
		t.Errorf("expected PartID=drivers/tb6612fng, got %q", drv.PartID)
	}
	if drv.Type != "motor_driver" {
		t.Errorf("expected Type=motor_driver, got %q", drv.Type)
	}
	if drv.Name == "" {
		t.Errorf("expected Name to be set")
	}

	if drv.MotorDriver.Channels != 2 {
		t.Errorf("expected Channels=2, got %d", drv.MotorDriver.Channels)
	}
	if drv.MotorDriver.MotorVoltageMin <= 0 || drv.MotorDriver.MotorVoltageMax <= 0 {
		t.Errorf("expected motor voltage range to be set, got [%.2f, %.2f]",
			drv.MotorDriver.MotorVoltageMin, drv.MotorDriver.MotorVoltageMax)
	}
	if drv.MotorDriver.CurrentContinuous <= 0 || drv.MotorDriver.CurrentPeak <= 0 {
		t.Errorf("expected currents to be >0, got continuous=%.2f, peak=%.2f",
			drv.MotorDriver.CurrentContinuous, drv.MotorDriver.CurrentPeak)
	}
}

func TestStore_LoadMotor_Generic12V(t *testing.T) {
	store := NewStore(testPartsDir(t))

	m, err := store.LoadMotor("motors/generic_dc_12v_gearmotor")
	if err != nil {
		t.Fatalf("LoadMotor(generic_dc_12v_gearmotor) returned error: %v", err)
	}

	if m.Type != "motor" {
		t.Errorf("expected Type=motor, got %q", m.Type)
	}
	if m.Name == "" {
		t.Errorf("expected Name to be set")
	}

	if m.Motor.NominalCurrent <= 0 || m.Motor.StallCurrent <= 0 {
		t.Errorf("expected non-zero currents, got nominal=%.2f, stall=%.2f",
			m.Motor.NominalCurrent, m.Motor.StallCurrent)
	}
}

func TestStore_LoadMCU_ESP32S3(t *testing.T) {
	store := NewStore(testPartsDir(t))

	mcu, err := store.LoadMCU("mcus/esp32s3")
	if err != nil {
		t.Fatalf("LoadMCU(esp32s3) returned error: %v", err)
	}

	if mcu.Type != "mcu" {
		t.Errorf("expected Type=mcu, got %q", mcu.Type)
	}
	if mcu.Name == "" {
		t.Errorf("expected Name to be set")
	}
	if mcu.MCU.LogicVoltage <= 0 {
		t.Errorf("expected non-zero logic voltage, got %.2f", mcu.MCU.LogicVoltage)
	}
}

func TestStore_LoadMissingPart_ReturnsError(t *testing.T) {
	store := NewStore(testPartsDir(t))

	if _, err := store.LoadDriver("drivers/does_not_exist"); err == nil {
		t.Fatalf("expected error when loading missing driver, got nil")
	}
}
