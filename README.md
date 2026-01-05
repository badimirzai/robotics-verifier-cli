# Robotics Verifier (rv-cli) [![CI](https://github.com/badimirzai/robotics-verifier-cli/actions/workflows/ci.yaml/badge.svg?branch=main)](https://github.com/badimirzai/robotics-verifier-cli/actions/workflows/ci.yaml) [![Release](https://img.shields.io/github/v/release/badimirzai/robotics-verifier-cli?label=release)](https://github.com/badimirzai/robotics-verifier-cli/releases)

A **hardware compatibility linter** for robotics.

Run `rv check` on a YAML spec to catch electrical integration failures **before you build**. Stop frying components and wasting weeks on shipping cycles.

```bash
# Verify your build in seconds
rv check robot_spec.yaml 
```

### The Problem
Most robotics failures happen before the first line of code ever runs. Incorrect voltage ranges, drivers that can't handle stall currents, and logic level mismatches waste time and damage expensive parts.

This tool enforces a hardware contract so these "silent killers" surface immediately-both on your local machine and in CI.

* **Catch Electrical Mismatches**: Detects logic level gaps (e.g., 3.3V vs 5V), I2C address conflicts, and voltage range violations.

* **Mechanical Safety**: Validates motor torque/current requirements against driver peak and continuous limits.

* **Power Budgeting**: Checks battery C-rate and discharge limits against total peak stall currents.

* **Deterministic & Private**: No AI hallucinations and no network calls. It only validates the specs you provide.

---

## How it works (at a glance)

![Architecture Overview](./assets/rvcli_architecture.png)

---

## Quick start (90 seconds)

### Install

Requires Go **1.25.5** or newer (https://go.dev/dl/).

```bash
go install github.com/badimirzai/robotics-verifier-cli/cmd/rv@latest
rv --help
```

### Try it in 30 seconds

```bash
rv init --list
# templates:
# - 4wd-problem
# - 4wd-clean
rv init --template 4wd-problem
# Wrote robot.yaml (template: 4wd-problem)

rv check robot.yaml
# shows multiple ERROR/WARN findings, exit code 2

rv init --template 4wd-clean --out robot.yaml --force
# Wrote robot.yaml (template: 4wd-clean)

rv check robot.yaml
# clean (or WARN-only if you intentionally keep some warnings), exit code 0
```



### Minimal example

Create a file named `spec.yaml`:

```yaml
name: "minimal-voltage-mismatch"

power:
  battery:
    voltage_v: 12
    max_current_a: 10
  logic_rail:
    voltage_v: 3.3
    max_current_a: 1

mcu:
  name: "Generic MCU"
  logic_voltage_v: 3.3
  max_gpio_current_ma: 12

motor_driver:
  name: "TB6612FNG-like"
  motor_supply_min_v: 18
  motor_supply_max_v: 24
  continuous_per_channel_a: 0.6
  peak_per_channel_a: 6
  channels: 1
  logic_voltage_min_v: 3.0
  logic_voltage_max_v: 5.5

motors:
  - name: "DC motor"
    count: 1
    voltage_min_v: 6
    voltage_max_v: 12
    stall_current_a: 5
    nominal_current_a: 1
```

Run:

```bash
rv check spec.yaml
```

Example output:

<pre>
rv check
--------------
<span style="color:#00a6d6">INFO</span> DRV_CHANNELS_OK: channels OK: 1 motors &lt;= 1 motor_driver.channels
<span style="color:#d10f1a">ERROR</span> DRV_SUPPLY_RANGE: battery 12.00V outside motor_driver motor supply range [18.00, 24.00]V
<span style="color:#c99200">WARN</span> DRV_CONT_LOW_MARGIN: motor_driver.continuous_per_channel_a 0.60A may be low for motor DC motor nominal 1.00A (want &gt;= 1.25A)
<span style="color:#00a6d6">INFO</span> RAIL_BUDGET_NOTE: logic rail budget set to 1.00A (v1 does not estimate MCU+driver logic draw yet)
exit code: 2
</pre>

Human-readable output uses color in terminals (header/OK green, INFO cyan, WARN yellow, ERROR red). Disable with `--no-color` or the standard `NO_COLOR` environment variable (set to any non-empty value, e.g. `NO_COLOR=1`).

Interpretation:
- The supply voltage cannot power the driver. This is a hard stop.
- The driver continuous current is lower than motor nominal current margin. Proceeding is risky.


## Example (video)
https://github.com/user-attachments/assets/3c73410f-bda8-49a3-9171-b888dff7446e


Example run catching voltage, current, battery C-rate, logic-level, and I2C conflicts in a 4-wheel mobile robot before anything gets built. 

---





## What it checks today

- Voltage compatibility between supply and drivers
- Current sufficiency for stall and nominal loads
- Driver to motor channel allocation
- Basic logic level consistency
- Logic rail compatibility between MCU and motor driver
- Battery C rate vs total peak stall current (motors)
- Total motor stall current vs driver peak current across all channels
- Simple I2C address conflicts on a single bus (duplicate device addresses)


### Core commands

```text
rv check <file.yaml>       Run analysis
rv version                 Show installed version
rv check --output json     Emit JSON findings
rv --help                  Show all commands and flags
rv check --help            Show check command options
```

**Note**: Checks are skipped when required inputs are missing (zero). This keeps partial specs usable.

Findings:
- INFO for context
- WARN for risk
- ERROR for violations (non zero exit code)

CI example:

```yaml
steps:
  - name: Verify hardware spec
    run: rv check specs/robot.yaml
```

JSON example:

```bash
rv check specs/robot.yaml --output json
```

JSON file output example:

```bash
rv check specs/robot.yaml --output json --out-file report.json
```

JSON pretty output example:

```bash
rv check specs/robot.yaml --output json --pretty
```

JSON pretty + file output example:

```bash
rv check specs/robot.yaml --output json --pretty --out-file report.json
```

When using `--output json` or `--output json --pretty` in a terminal, severity values are colorized for readability. Colors are never used for JSON files or non-TTY output.

---

## Supported configurations (v0.1)

Focused on early stage mobile robots.

Supported:
- DC motors (one motor per driver channel)
- TB6612FNG and L298 class H bridge drivers
- Single logic rail
- Basic YAML part inheritance

Not supported yet:
- Stepper, BLDC, ESC
- Multi rail power trees
- Thermal derating
- Serial or IO protocol arbitration

This is a linter. Not a simulator or optimizer.

---

## YAML specification

The core fields used in validation are:

- power.battery.voltage_v
- power.battery.capacity_ah
- power.battery.c_rating
- power.battery.max_discharge_a
- power.battery.max_current_a
- motor_driver.motor_supply_min_v
- motor_driver.motor_supply_max_v
- motors[].stall_current_a
- motor_driver.peak_per_channel_a

Battery max discharge uses the following precedence:
1) power.battery.max_discharge_a
2) power.battery.capacity_ah * power.battery.c_rating
3) power.battery.max_current_a

I2C bus structure (addresses accept decimal or 0x hex):

```yaml
i2c_buses:
  - name: "bus0"
    devices:
      - name: "imu_left"
        address_hex: 0x68
      - name: "imu_right"
        address_hex: 104
```

Unset or missing fields are treated as unknown. Some required values will surface as errors during resolution.

More examples are available in the `examples/` directory.

---

## Versioning and stability

The interface is still evolving. Breaking changes may happen before 1.0.

Exit codes and rule identifiers are stable within a minor version:
- 0 clean
- 2 rule violations
- 3+ parser or internal errors

---

## CLI output options

Output control flags (check command):
- --output json: machine readable JSON to stdout
- --pretty: pretty print JSON to stdout (requires --output json)
- --out-file <path>: write compact JSON to file (requires --output json)
- --debug: enable debug mode (or use RV_DEBUG=1)

Output behavior matrix:
- rv check spec.yaml: human-readable output
- rv check spec.yaml --output json: compact JSON to stdout
- rv check spec.yaml --output json --pretty: pretty JSON to stdout
- rv check spec.yaml --output json --out-file result.json: writes compact JSON and prints "Written to result.json"
- rv check spec.yaml --output json --pretty --out-file result.json: pretty JSON to stdout and compact JSON to file

See `CHEATSHEET.md` for a quick command reference.

---

## Determinism first

Trust is the primary feature.

This tool does not guess, fetch, or infer part data. It validates what you specify. Assistive or automated layers may come later, but only on top of a proven deterministic core.

---

## Contributing

Open an issue before starting work so scope can be aligned.

By contributing you agree to the CLA in `CLA.md`.

---

## License

Apache 2.0. See `LICENSE`.

---

## Disclaimer

This tool does not replace datasheets or engineering judgement.
Not suitable for safety critical systems.
Use at your own risk. Early alpha.
