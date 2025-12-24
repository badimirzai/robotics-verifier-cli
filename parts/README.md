# RoboStack Parts Library (v0)

These YAML files contain canonical, conservative electrical properties used by RoboStack validation rules.

Rules:
- Only include fields RoboStack validates today.
- Prefer conservative maxima/minima.
- Add `sources` with datasheet links.
- Add `confidence` per field where possible.

File naming:
- drivers/<mpn>.yaml (lowercase)
- motors/<slug>.yaml
- mcus/<slug>.yaml