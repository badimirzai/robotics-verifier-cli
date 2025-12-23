# Contributing to RoboStack

Thanks for your interest in contributing to RoboStack.

RoboStack is an open-core project focused on deterministic, explainable
verification of robotics hardware stacks. Contributions that improve
correctness, clarity, and robustness are especially welcome.

---

## How to contribute

You can contribute by:
- Reporting bugs
- Improving documentation
- Adding validation rules
- Extending existing checks
- Adding examples and tests

Please open an issue before starting large changes.

---

## Code quality

Contributions should:
- Be deterministic and reproducible
- Avoid introducing hidden assumptions
- Prefer explicit checks over heuristics
- Include tests where appropriate
- Keep outputs explainable

If a rule cannot be explained, it should not exist.

---

## License

By contributing to this repository, you agree that your contributions
will be licensed under the Mozilla Public License 2.0 (MPL 2.0).

You retain copyright to your contributions.

---

## Scope and direction

RoboStack maintains a strict separation between:
- the open verification core
- potential future extensions and services

Contributions should focus on improving the open core.

If you are unsure whether a change fits the project’s scope, open an issue
to discuss it first.