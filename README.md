# RoboStack CLI

A deterministic verifier and part-evaluation tool for early-stage robotics BOMs.

RoboStack helps you catch **electrical and control-level incompatibilities before you buy parts**.  
It is intentionally strict, explainable, and automation-friendly.

This is **not** a parts marketplace, optimizer, or AI recommender.  
It is a verifier first, a decision-support tool second.

---

## Why this exists

Most robotics failures happen *before* firmware is written:
- Wrong motor driver voltage range
- Undersized current limits vs stall current
- Logic-level mismatches between MCU, rails, and drivers
- Silent assumptions hidden in datasheets

RoboStack turns these assumptions into **explicit, machine-checkable rules**.

---

## What it does today

### 1. Deterministic stack verification
RoboStack consumes a small YAML spec describing your robot stack:
- Battery
- Logic rail
- MCU
- Motor driver
- Motors

It then runs cross-domain checks, including:
- Driver channel count vs motor count
- Battery voltage vs driver motor supply range
- Driver peak/continuous current vs motor stall/nominal current
- Logic voltage compatibility (MCU ↔ rail ↔ driver)
- Basic power-rail sanity

### 2. Structured, CI-friendly output
- Findings are emitted as `INFO`, `WARN`, or `ERROR`
- A non-zero exit code is returned if any `ERROR` is present
- Output is deterministic and explainable

This makes RoboStack safe to use in:
- CI pipelines
- Design checklists
- Automated BOM validation

---

## What it deliberately does *not* do (yet)

- No AI-based decision making
- No opaque scoring
- No automatic purchasing
- No optimization without explanations

Those come **later**, if and only if the deterministic core proves useful.

---

## Installation

```sh
git clone https://github.com/badi96/robostack-cli
cd robostack-cli
go build
```

Or run directly:

```sh
go run . verify -f examples/amr_basic.yaml
```

---

## Spec format

The input spec is a small YAML file describing your robot stack components and their key electrical properties.

Example:

```yaml
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

See `examples/amr_basic.yaml` for a complete, working example.

---

## Output and exit codes

- Results are grouped by severity: `INFO`, `WARN`, `ERROR`
- Exit code:
  - `0`: No errors
  - `2`: One or more `ERROR` entries detected

This makes RoboStack suitable for automated enforcement in CI.

---

## Roadmap (high level)

Near-term:
- Canonical `Part` data model with confidence tracking
- Supplier adapters (starting with Mouser)
- Data completeness reporting

Mid-term:
- Candidate filtering and transparent ranking
- Reason traces for every score contribution

Long-term:
- Optional backend for advanced ranking
- AI-generated explanations only (never decisions)

---

## Design principles

- **Deterministic over clever**
- **Explainable over optimized**
- **Fail fast and loudly**
- **No hidden assumptions**
- **No AI in the decision loop**

If the tool cannot explain *why* a part is rejected or ranked lower, it is considered broken.

---

## License

MIT (subject to change)

---

## Disclaimer

RoboStack does not replace datasheets or engineering judgment.  
It exists to surface mistakes early, not to remove responsibility.
