# Architon

**Specification layer for robot hardware**  
Architon turns your robot's electrical architecture into a machine‑checkable spec.  
It replaces scattered spreadsheets and implicit assumptions with structured YAML verified by a rule engine before you ever order parts.

Deterministic today by design. Assistive AI can enter later when guarantees exist.

---

## Why this exists

Most robotics failures start before firmware runs: mismatched voltages, undersized current paths, logic‑level mistakes.  
Architon provides a **contract** for electrical architecture so these failures surface early, locally and in CI.

---

## What it does today

- Parse hardware specs defined in YAML
- Check voltage, current, and logic‑level compatibility
- Emit structured findings: `INFO`, `WARN`, `ERROR`
- Exit non‑zero on `ERROR` so CI can block merges
- CLI‑friendly outputs for terminals and pipelines

Example output:

```text

make validate
go run . validate -f examples/amr_basic.yaml
robostack validate
--------------
INFO DRV_CHANNELS_OK: channels OK: 2 motors <= 2 driver channels
ERROR DRV_SUPPLY_RANGE: battery 14.80V outside driver motor supply range [2.50, 13.50]V
ERROR DRV_PEAK_LT_STALL: driver peak 3.20A < motor DC gearmotor stall 5.00A (per channel)
WARN DRV_CONT_LOW_MARGIN: driver continuous 1.20A may be low for motor DC gearmotor nominal 1.50A (want >= 1.88A)
INFO RAIL_BUDGET_NOTE: logic rail budget set to 2.00A (v1 does not estimate MCU+driver logic draw yet)

exit status 2
make: *** [validate] Error 1
```

---

## Architecture in one paragraph

1. **Spec** – YAML definitions that capture robot hardware intent.  
2. **Engine** – deterministic checks with explainable rules.  
3. **Assistants (future)** – optional helpers that operate *on top* once guarantees exist.  
   They **do not replace** the rule engine. AI will never silently assert correctness.

This is a sequencing principle, not a permanent ban on ML. Trust first, automation second.

---

## Installation

### Prerequisites

- Go **1.25.5** or newer: https://go.dev/dl/
- GOPATH/bin added to your PATH

## Quick start

### Install

```bash
git clone https://github.com/badimirzai/robostack-cli.git
cd robostack-cli

# build and install the CLI
go install ./...
```

Make sure your `GOPATH/bin` is on your `PATH`, then you should have an `architon` binary available.

### Run validate on an example

```bash
robostack-cli validate -f ./examples/amr_basic.yaml
```

Example output:

```text
ERROR DRV_SUPPLY_RANGE: battery 14.80V outside driver motor supply range [2.50, 13.50]V
ERROR DRV_PEAK_LT_STALL: driver peak 3.20A < motor DC gearmotor stall 5.00A (per channel)
WARN DRV_CONT_LOW_MARGIN: driver continuous 1.20A may be low for motor DC gearmotor nominal 1.50A (want >= 1.88A)
INFO RAIL_BUDGET_NOTE: logic rail budget set to 2.00A (v1 does not estimate MCU+driver logic draw yet)
```

Non-zero exit codes are used when there is at least one `ERROR`, so you can plug this straight into CI pipelines.


---

## CLI behavior

- `INFO`: contextual notes
- `WARN`: non‑ideal but non‑blocking
- `ERROR`: hard violations, non‑zero exit code
- Designed for terminals and CI; no spinner fluff or hidden state

Typical CI usage:

```yaml
steps:
  - name: Run hardware checks
    run: Architon check specs/amr.yaml
```

---

## Current Scope & Limitations (v1)

Architon v1 is intentionally narrow. It is designed to lint early-stage mobile robots built around **DC gearmotors**, **H-bridge motor drivers**, and a **single logic rail**. The goal is to prevent category-error mistakes before money is spent.

### ✔️ Supported in v1
- DC motors (1 motor per driver channel)
- H-bridge motor drivers (TB6612FNG, L298 class)
- Single logic rail verification (`power.logic_rail`)
- Basic electrical compatibility checks:
  - Battery voltage vs driver logic/motor voltage ranges
  - Motor stall current vs driver peak current
  - Motor nominal current vs driver continuous current (with margin)
  - MCU logic voltage vs logic rail (level shifting risk)
- `part:` references and default value merging from the parts library

### ❌ Not supported (yet)
- Stepper motors (2-phase, 4-phase)
- BLDC / ESC drivers (3-phase)
- Multi-rail power trees (24 V bus, 5 V logic, 3.3 V sensor rails)
- Thermal/derating or PCB trace current modeling
- Battery chemistry discharge curves, IR drop modeling
- Integration with Mouser/DigiKey/Octopart APIs (planned v2)
- IO-level protocol compatibility (UART/SPI/I2C voltage domain arbitration)

### ⚠️ Assumptions (v1)
- One motor per driver channel (DC only)
- Zero values in YAML mean "unset" and will be filled from `part:`
- 25 percent current margin heuristic for continuous current checks
- Validation uses **nominal** battery voltage (not max/min chemistry curves)
- Errors and warnings reflect deterministic rules, not probabilistic models

---

Architon v1 is a **linter** — not a simulator and not an optimizer.  
It focuses on correctness over completeness and prioritizes **explainable rule-based checks**.  
Future versions (v2+) will expand into multi-rail systems, API-assisted part selection, and HIL-oriented verification workflows.

---

## Roadmap

This is **not a promise**, it is direction. Order and scope may change.

- richer rule sets (AWG, simple derating, interface compatibility)
- supplier adapters → canonical part model (surface missing fields)
- KiCad boilerplates from specs
- ROS2 scaffolding aligned to the same spec
- assistive tooling for part selection / impact analysis **on top** of rules

> Everything here remains deterministic first.  
> AI/ML may appear once it can provide value without eroding guarantees.

---

## Contributing

Focus areas right now:

- real-world example specs
- missing rule proposals
- feedback on naming and structure

Small, surgical PRs are preferred.

---

## License
This project is open source under the MIT license. See `LICENSE` for details.

## Contributions
By contributing, you agree to the Contributor License Agreement in `CLA.md`.
This allows the project to relicense or commercialize the work in the future.

---

## Disclaimer
Archeon does not replace datasheets, safety analysis, or engineering judgment.  
It is intended for **early-stage verification and decision support**.
This is **early alpha**.  
Interfaces will break, specs will evolve, and rules will change.  
Do not depend on this for safety‑critical systems yet.
