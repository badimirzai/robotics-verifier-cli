# Robotics Verifier CLI
[![CI](https://github.com/badimirzai/robotics-verifier-cli/actions/workflows/ci.yaml/badge.svg?branch=main)](https://github.com/badimirzai/robotics-verifier-cli/actions/workflows/ci.yaml)

[![Release](https://img.shields.io/github/v/release/badimirzai/robotics-verifier-cli?label=release)](https://github.com/badimirzai/robotics-verifier-cli/releases)

**The Specification layer for robotics hardware.**  
Robotics Verifier CLI turns your robot's electrical architecture into a machine‑checkable spec.  
It replaces scattered spreadsheets and implicit assumptions with structured YAML verified by a rule engine before you ever order parts.

Deterministic today by design. Assistive AI can enter later when guarantees exist.

Project codename: Architon (tentative, may change).

---

## Why this exists

Most robotics failures start before firmware runs: mismatched voltages, undersized current paths, logic‑level mistakes.  
Robotics Verifier CLI provides a **contract** for electrical architecture so these failures surface early, locally and in CI.

---

## What it does today

- Parse hardware specs defined in YAML
- Check voltage, current, and logic‑level compatibility
- Emit structured findings: `INFO`, `WARN`, `ERROR`
- Exit non‑zero on `ERROR` so CI can block merges
- CLI‑friendly outputs for terminals and pipelines

Example output:

```text
rv check
--------------
INFO DRV_CHANNELS_OK: channels OK: 2 motors <= 2 driver channels
ERROR DRV_SUPPLY_RANGE: battery 14.80V outside driver motor supply range [2.50, 13.50]V
ERROR DRV_PEAK_LT_STALL: driver peak 3.20A < motor DC gearmotor stall 5.00A (per channel)
WARN DRV_CONT_LOW_MARGIN: driver continuous 1.20A may be low for motor DC gearmotor nominal 1.50A (want >= 1.88A)
INFO RAIL_BUDGET_NOTE: logic rail budget set to 2.00A (v1 does not estimate MCU+driver logic draw yet)
```

---

## Architecture in one paragraph

1. **Spec** – YAML definitions that capture robot hardware intent.  
2. **Engine** – deterministic checks with explainable rules.  
3. **Assistants (future)** – optional helpers that operate *on top* once guarantees exist.  
   They **do not replace** the rule engine. AI will never silently assert correctness.

This is a sequencing principle: **trust first, automation second.**

---

## Installation

### Prerequisites

- Go **1.25.5** or newer: https://go.dev/dl/
- GOPATH/bin added to your PATH

---

## Quick start

```bash
git clone https://github.com/badimirzai/robotics-verifier-cli.git
cd robotics-verifier-cli

# build and install the CLI
make build
make install     # installs robotics-verifier-cli and rv symlink
```

Verify install:

```bash
rv --help
```

---

### Run checks on an example

```bash
rv check ./examples/amr_parts.yaml
```

or:

```bash
make verify
```

---

### Build and run locally (no install)

```bash
make build
./bin/robotics-verifier-cli check ./examples/amr_basic.yaml
```

or:

```bash
go build -o ./bin/robotics-verifier-cli .
./bin/robotics-verifier-cli check ./examples/amr_basic.yaml
```

---

## CLI behavior

- `INFO`: contextual notes
- `WARN`: non‑ideal but non‑blocking
- `ERROR`: hard violations, non‑zero exit code
- Designed for terminals & CI; **no fluff**

Typical CI usage:

```yaml
steps:
  - name: Run hardware checks
    run: rv check specs/amr.yaml
```

---

## Current Scope & Limitations (v1)

v1 is intentionally narrow. It is designed to lint early-stage mobile robots built around **DC gearmotors**, **H-bridge motor drivers**, and a **single logic rail**. The goal is to prevent category-error mistakes **before money is spent**.

### ✔️ Supported in v1
- DC motors (1 motor per driver channel)
- H-bridge motor drivers (TB6612FNG, L298 class)
- Single logic rail verification (`power.logic_rail`)
- Basic electrical compatibility checks:
  - Battery voltage vs driver logic/motor voltage ranges
  - Motor stall current vs driver peak current
  - Motor nominal current vs driver continuous current (with margin)
  - MCU logic voltage vs logic rail (level shifting risk)
- `part:` references & default value merging from parts library

### ❌ Not supported yet
- Stepper motors
- BLDC / ESC
- Multi-rail trees
- Thermal/derating models
- API-assisted part import
- IO-level protocol arbitration

### ⚠️ Assumptions (v1)
- One motor per driver channel (DC only)
- Zero values in YAML mean "unset" and will be filled from `part:`
- 25 percent current margin heuristic for continuous current checks
- Validation uses **nominal** battery voltage (not max/min chemistry curves)
- Errors and warnings reflect deterministic rules, not probabilistic models

---


This tool is a **linter** — not a simulator and not an optimizer.  
It focuses on correctness over completeness and prioritizes **explainable rule-based checks**.

---

## Roadmap (direction, not promise)

- richer rule sets (AWG, derating, interface compatibility)
- supplier adapters → canonical part model
- KiCad boilerplates from specs
- ROS2 scaffolding aligned to the same spec
- assistive tooling (post-guarantee)

---

## Contributing

- real-world example specs
- missing rule proposals
- feedback on naming and structure

Small, surgical PRs preferred.

---

## License
MIT. See `LICENSE`.

---


## Continuous Integration (CI)

This repository uses GitHub Actions to ensure the CLI is always in a working state.

The pipeline performs:
1. Go setup and build
2. Unit tests for parts loader, resolver and rule logic
3. Smoke test of the CLI against example specs

Exit codes:
* 0 means no issues found
* 2 means the spec contains physical or configuration errors and the tool reported them
* 3 or higher indicates a crash or runtime fault and CI fails

CI does not fail on exit code 2 because it is the expected behavior for invalid specs. CI fails only if unit tests fail or if the CLI crashes.

### Local run

```
make ci
# or
go build ./...
go test ./...
go run . validate -f examples/amr_parts.yaml || true
```


## Contributor License Agreement
By contributing, you agree to the CLA in `CLA.md`.  
This allows future relicensing or commercialization.

---

## Disclaimer
Robotics Verifier CLI does **not** replace datasheets, safety analysis, or engineering judgment.  
It is intended for **early-stage verification and decision support**.  
This is **early alpha**. Interfaces will break. Specs will evolve. Rules will change.  
Do not use this for safety‑critical systems.
