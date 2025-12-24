# RoboStack CLI

A **deterministic verifier** for early-stage robotics Bills of Materials (BOMs).

RoboStack helps you catch electrical (and a few basic integration) incompatibilities **before you buy parts**.
It is intentionally strict, explainable, and automation-friendly.

This is **not** a parts marketplace or optimizer.
RoboStack enforces deterministic, explainable checks before any ranking or recommendation is considered.

---

## Status

**Alpha.** Rules and the YAML schema may change.  
If you rely on this tool, pin a commit.

---

## Why this exists

Most robotics failures happen before firmware is written:

- Wrong motor driver voltage range
- Undersized current limits vs motor stall current
- Logic-level mismatches between MCU, rails, and drivers
- Silent assumptions hidden in datasheets

RoboStack turns these assumptions into **explicit, machine-checkable rules**.

---

## What it does today

### Deterministic stack verification

RoboStack consumes a small YAML spec describing your robot stack:

- Battery
- Logic rail
- MCU
- Motor driver
- Motors

It then runs cross-domain checks, including:

- Driver channel count vs motor count
- Battery voltage vs driver motor supply range
- Driver peak and continuous current vs motor stall current (per channel)
- Driver continuous current vs motor nominal current (with margin warning)
- Basic power-rail sanity notes

### CI-friendly output

- Findings are emitted as `INFO`, `WARN`, or `ERROR`
- A non-zero exit code is returned if any `ERROR` is present
- Output is deterministic and explainable

---

## What it deliberately does not do (yet)

- No AI-based decision making
- No opaque scoring
- No automatic purchasing
- No optimization without explanations

Those come later **only if** the deterministic core proves useful.

---

## Prerequisites

- Go 1.25.5 (built with go1.25.5). Official releases: https://go.dev/dl/

## Install and run

### Option A: quickest (recommended)
Runs the included example spec.

```bash
make validate
```

Note: `make validate` will exit non-zero if the example triggers any `ERROR` rules.

### Option B: run any command via Go (no install)
Run validate on an example file:

```bash
go run . validate -f examples/amr_basic.yaml
```

Or via Make (this is the correct way to use the `run` target):

```bash
make run ARGS="validate -f examples/amr_basic.yaml"
```

### Option C: build a local binary (no PATH required)
Build:

```bash
make build
```

Run:

```bash
./bin/robostack validate -f examples/amr_basic.yaml
```

### Option D: install globally (puts `robostack` on your PATH)
Install:

```bash
make install
```

If your shell says `robostack: command not found`, add Go’s bin directory to PATH:

```bash
export PATH="$(go env GOPATH)/bin:$PATH"
```

---

## Example output

Example run:

```bash
go run . validate -f examples/amr_basic.yaml
```

Example output:

```text
robostack validate
--------------
INFO DRV_CHANNELS_OK: channels OK: 2 motors <= 2 driver channels
ERROR DRV_SUPPLY_RANGE: battery 14.80V outside driver motor supply range [2.50, 13.50]V
ERROR DRV_PEAK_LT_STALL: driver peak 3.20A < motor DC gearmotor stall 5.00A (per channel)
WARN DRV_CONT_LOW_MARGIN: driver continuous 1.20A may be low for motor DC gearmotor nominal 1.50A (want >= 1.88A)
INFO RAIL_BUDGET_NOTE: logic rail budget set to 2.00A (v1 does not estimate MCU+driver logic draw yet)
```

Exit codes:
- `0` = no `ERROR`
- `2` = one or more `ERROR`

---

## Spec format

The input spec is a small YAML file describing your robot stack components and key electrical properties.

```yaml
spec_version: 0.1

battery:
  voltage_nominal: 24
  voltage_max: 25.2

logic_rail:
  voltage: 5.0
  max_current: 3.0

mcu:
  logic_voltage: 3.3

motor_driver:
  channels: 2
  motor_voltage_min: 8
  motor_voltage_max: 30
  logic_voltage_min: 3.0
  logic_voltage_max: 5.5
  current_continuous: 1.5
  current_peak: 3.0

motors:
  count: 2
  nominal_current: 0.8
  stall_current: 2.2
```

See `examples/amr_basic.yaml` for a complete working example.

---

## Roadmap (high level)

**Near-term**
- Canonical part data model with confidence tracking
- Supplier adapters (starting with Mouser)
- Data completeness reporting

**Mid-term**
- Candidate filtering and transparent ranking
- Reason traces for every score contribution

**Long-term**
- Optional backend for advanced ranking
- AI-generated explanations only (never decisions)

---

## Design principles

- Deterministic over clever
- Explainable over optimized
- Fail fast and loudly
- No hidden assumptions
- No AI in the decision loop

If the tool cannot explain why a part is rejected or ranked lower, it is considered broken.

---

## License

This project is licensed under the **Mozilla Public License 2.0 (MPL 2.0)**.  
See the `LICENSE` file for details.

---

## Disclaimer

RoboStack does not replace datasheets, safety analysis, or engineering judgment.  
It is intended for **early-stage verification and decision support**.
