package parts

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
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
	if drv.MotorDriver.MotorSupplyMinV <= 0 || drv.MotorDriver.MotorSupplyMaxV <= 0 {
		t.Errorf("expected motor voltage range to be set, got [%.2f, %.2f]",
			drv.MotorDriver.MotorSupplyMinV, drv.MotorDriver.MotorSupplyMaxV)
	}
	if drv.MotorDriver.ContinuousPerChA <= 0 || drv.MotorDriver.PeakPerChA <= 0 {
		t.Errorf("expected currents to be >0, got continuous=%.2f, peak=%.2f",
			drv.MotorDriver.ContinuousPerChA, drv.MotorDriver.PeakPerChA)
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

	if m.Motor.NominalCurrentA <= 0 || m.Motor.StallCurrentA <= 0 {
		t.Errorf("expected non-zero currents, got nominal=%.2f, stall=%.2f",
			m.Motor.NominalCurrentA, m.Motor.StallCurrentA)
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
	if mcu.MCU.LogicVoltageV <= 0 {
		t.Errorf("expected non-zero logic voltage, got %.2f", mcu.MCU.LogicVoltageV)
	}
}

func TestStore_LoadMissingPart_ReturnsError(t *testing.T) {
	store := NewStore(testPartsDir(t))

	if _, err := store.LoadDriver("drivers/does_not_exist"); err == nil {
		t.Fatalf("expected error when loading missing driver, got nil")
	}
}

func TestStore_LoadMotor_PrefersEarlierDir(t *testing.T) {
	tmp := t.TempDir()
	localDir := filepath.Join(tmp, "rv_parts")
	builtInDir := filepath.Join(tmp, "parts")
	if err := os.MkdirAll(filepath.Join(localDir, "motors"), 0o755); err != nil {
		t.Fatalf("mkdir local parts: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(builtInDir, "motors"), 0o755); err != nil {
		t.Fatalf("mkdir built-in parts: %v", err)
	}

	partID := "motors/test_motor"
	localPart := `part_id: motors/test_motor
type: motor
name: Local Motor
motor:
  voltage_min_v: 6
  voltage_max_v: 12
  nominal_current_a: 0.5
  stall_current_a: 1.2
`
	builtInPart := `part_id: motors/test_motor
type: motor
name: Built-in Motor
motor:
  voltage_min_v: 4
  voltage_max_v: 9
  nominal_current_a: 0.4
  stall_current_a: 1.0
`
	if err := os.WriteFile(filepath.Join(localDir, "motors", "test_motor.yaml"), []byte(localPart), 0o644); err != nil {
		t.Fatalf("write local part: %v", err)
	}
	if err := os.WriteFile(filepath.Join(builtInDir, "motors", "test_motor.yaml"), []byte(builtInPart), 0o644); err != nil {
		t.Fatalf("write built-in part: %v", err)
	}

	store := NewStoreWithDirs([]string{localDir, builtInDir})
	part, err := store.LoadMotor(partID)
	if err != nil {
		t.Fatalf("LoadMotor returned error: %v", err)
	}
	if part.Name != "Local Motor" {
		t.Fatalf("expected local part to win, got %q", part.Name)
	}
}

func TestStore_LoadMissingPart_ReportsSearchDirs(t *testing.T) {
	tmp := t.TempDir()
	localDir := filepath.Join(tmp, "rv_parts")
	builtInDir := filepath.Join(tmp, "parts")
	store := NewStoreWithDirs([]string{localDir, builtInDir})

	_, err := store.LoadMotor("motors/missing")
	if err == nil {
		t.Fatalf("expected error when loading missing motor, got nil")
	}
	msg := err.Error()
	if !strings.Contains(msg, "motors/missing") {
		t.Fatalf("expected error to mention part id, got %q", msg)
	}
	if !strings.Contains(msg, localDir) || !strings.Contains(msg, builtInDir) {
		t.Fatalf("expected error to list search dirs, got %q", msg)
	}
}

func TestStore_LoadI2CSensor_MPU6050(t *testing.T) {
	store := NewStore(testPartsDir(t))

	sensor, err := store.LoadI2CSensor("sensors/mpu6050")
	if err != nil {
		t.Fatalf("LoadI2CSensor(mpu6050) returned error: %v", err)
	}

	if sensor.Type != "i2c_sensor" {
		t.Errorf("expected Type=i2c_sensor, got %q", sensor.Type)
	}
	if sensor.Name == "" {
		t.Errorf("expected Name to be set")
	}
	if sensor.I2CDevice.AddressHex == 0 {
		t.Errorf("expected non-zero address, got 0")
	}
}
