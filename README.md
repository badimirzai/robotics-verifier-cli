# Robotics Verifier CLI [![CI](https://github.com/badimirzai/robotics-verifier-cli/actions/workflows/ci.yaml/badge.svg?branch=main)](https://github.com/badimirzai/robotics-verifier-cli/actions/workflows/ci.yaml) [![Release](https://img.shields.io/github/v/release/badimirzai/robotics-verifier-cli?label=release)](https://github.com/badimirzai/robotics-verifier-cli/releases)

**A specification and verification layer for robotics hardware.**
It turns electrical architecture into a machine-checkable spec, surfacing voltage, current, and logic-level mistakes before firmware or fabrication.

Deterministic today by design. Assistive layers may come later only if the core earns trust.

*(Codename: Architon — tentative, not part of interface guarantees.)*


---

## Why this exists

Most robotics failures start **before firmware runs**: mismatched voltages, undersized current paths, logic‑level mistakes.  
This CLI provides a **contract** for electrical architecture so these failures surface early, locally and in CI.


![Architecture Overview](./assets/rvcli_architecture.png)


---

## What it does today

- Parse hardware specs defined in YAML
- Check voltage, current, and logic‑level compatibility
- Emit structured findings: `INFO`, `WARN`, `ERROR`
- Exit non‑zero on `ERROR` so CI can block merges
- CLI‑friendly outputs for terminals and pipelines

Example:

```text
rv check
--------------
INFO DRV_CHANNELS_OK: channels OK: 2 motors <= 2 driver channels
ERROR DRV_SUPPLY_RANGE: battery 14.80V outside driver motor supply range [2.50, 13.50]V
ERROR DRV_PEAK_LT_STALL: driver peak 3.20A < motor stall 5.00A
WARN DRV_CONT_LOW_MARGIN: driver continuous 1.20A may be low for motor nominal 1.50A
```

---
## Start here (5 steps)

0. Install: `make install`
1. Copy template: `cp examples/amr_basic.yaml specs/robot.yaml`
2. Edit only these fields first:
   - battery_voltage
   - driver_peak_current
   - motor_stall_current
3. Run: `rv check specs/robot.yaml`
4. Fix any ERROR. Commit only when clean.

If you hit a shortage of fields or missing rules, open an issue. Scope must be tightened intentionally.

---

## Architecture in one paragraph

1. **Spec** – YAML captures robot hardware intent.  
2. **Engine** – deterministic checks with explainable rules.  
3. **Assistants (future)** – optional helpers that operate **on top**, never replacing rule logic.

**Principle:** trust first, automation second.

---

## Installation

### Prerequisites

- Go **1.25.5** or newer: https://go.dev/dl/
- GOPATH/bin in your PATH

---

## Quick start

```bash
git clone https://github.com/badimirzai/robotics-verifier-cli.git
cd robotics-verifier-cli
make build
make install     # installs robotics-verifier-cli and rv symlink
rv --help
```

---

### Run checks on an example

```bash
rv check ./examples/amr_parts.yaml
```

or

```bash
make verify
```

---

### Local build (no install)

```bash
make build
./bin/robotics-verifier-cli check ./examples/amr_basic.yaml
```

---

## CLI behavior

- `INFO`: contextual notes
- `WARN`: non‑ideal but non‑blocking
- `ERROR`: hard violations, non‑zero exit code
- Designed for terminals & CI
- **No hidden network calls**

Typical CI:

```yaml
steps:
  - name: Run hardware checks
    run: rv check specs/amr.yaml
```

---

## Current Scope & Limitations (v1)

v1 is intentionally narrow: lint early-stage mobile robots built around **DC gearmotors**, **H‑bridge drivers**, and a **single logic rail**.  
Goal: prevent category‑error mistakes **before money is spent**.

### ✔️ Supported (v1)
- DC motors (1 motor per driver channel)
- TB6612FNG / L298‑class drivers
- Single logic rail verification
- Basic electrical compatibility checks

### ❌ Not in core yet
- Stepper motors / BLDC / ESC
- Multi‑rail trees
- Thermal/derating models
- API‑assisted part import
- IO‑level protocol arbitration

### ⚠️ Assumptions
- One motor per driver channel
- Zero in YAML means "unset" → pulled from `part:`
- 25 percent current margin for continuous current checks
- Uses nominal battery voltage

This is a **linter** — not a simulator or optimizer.

---

## Roadmap direction (not promise)

- richer rule sets (AWG, derating, interface compatibility)
- supplier adapters → canonical part model
- KiCad boilerplates from specs
- ROS2 scaffolding
- assistive tooling (post‑guarantee)

---

## Contributing

By contributing you agree to the CLA in `CLA.md` (relicensing & commercial rights).  
Start by opening an issue to align scope **before** coding.

See `CONTRIBUTING.md`.

---

## License

Apache‑2.0. See `LICENSE`.

---

## CI Overview

GitHub Actions runs build + tests + smoke checks.

Exit codes:
- `0`: clean
- `2`: spec violations (expected if invalid input)
- `>=3`: crash (CI fail)

Local:

```bash
make ci
```

---

## Disclaimer

This does **not** replace datasheets or engineering judgment. 
Intended for early-stage verification and decision support.  
Early alpha. Interfaces will break. Rules will evolve.  
Not for safety‑critical systems.
