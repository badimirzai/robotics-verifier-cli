# Contributing to Architon CLI

Thank you for considering contributing. This project is early stage and focused on correctness, clarity, and predictable behavior.

## Ground Rules

- By contributing, you agree to the MIT License and the CLA in `CLA.md`.
- All contributions must be your own work.
- Do not submit code you cannot certify for commercial use.
- No feature creep. Open an issue before building new features.
- Keep scope tight and avoid premature abstractions.

## Code Expectations

- Keep functions small and single purpose.
- Write clear validation logic with explicit error messages.
- Favor directness over cleverness.
- No unexplained magic constants. Inline comments if necessary.
- Consistent naming: verbs for commands, nouns for data structures.

## Tests

- Every contribution that changes behavior must have tests.
- Tests must be deterministic and not rely on remote resources.
- If you fix a bug, write a test that fails before the fix and passes after.

## Pull Requests

1. Create a branch from `main`.
2. Ensure all tests pass locally.
3. Add or update relevant documentation.
4. Keep diffs minimal and focused on one concern.
5. Request review via PR. Do not merge without approval.

## Commit Style

- Use clear prefixes:
  - `feat:` new feature
  - `fix:` bug fix
  - `docs:` documentation changes
  - `test:` tests only
  - `refactor:` structural changes without behavior change
- Keep messages brief and to the point.

## Communication

- Open an issue for discussions larger than a few lines of change.
- Stay civil and technical. Critique code, not people.
- Decisions prioritize stability, clarity, and maintainability.

## Licensing

- All contributions are under MIT License.
- CLA must be signed for PRs to be accepted.
- You retain copyright to your contributions.

## Rejection Criteria

Your PR will be rejected if:

- It increases complexity without clear benefit.
- It introduces magic behavior or silent failures.
- It duplicates existing functionality without justification.
- It adds heavy external dependencies without approval.
- It solves a problem that is not validated or discussed.

## Final Note

This project is evolving. Expect rules to refine over time. Drive the project forward with clarity and discipline.
