# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Roadmap

## [0.0.4] - 2021-10-25
### Changed
- Using `math.MaxInt32` instead of `math.MaxInt`. `math.MaxInt` is only available in Go 1.17+.

## [0.0.3] - 2021-10-25
### Added
- Added the ability to generate X amount of random numbers.

## [0.0.2] - 2021-10-01
### Added

- More tests.

### Changed
- Improved validation for `min` and `max`.
- Fixed `Generate` panic in case of `rand.Int` failure.

## [0.0.1] - 2021-09-24
### Added
- [x] Ability to create custom errors
- [x] Ability to create custom errors with code
- [x] Ability to create custom errors with status code
- [x] Ability to create custom errors with message
- [x] Ability to create custom errors wrapping an error
- [x] Ability to create static (pre-created) custom errors
- [x] Ability to create dynamic (in-line) custom errors
- [x] Ability to print a custom error with a dynamic, and custom message

### Checklist
- [x] CI Pipeline:
  - [x] Lint
  - [x] Tests
- [x] Documentation:
  - [x] Package's documentation (`doc.go`)
  - [x] Meaningful code comments, and symbol names (`const`, `var`, `func`)
  - [x] `GoDoc` server tested
  - [x] `README.md`
  - [x] `LICENSE`
    - [x] Files has LICENSE in the header
  - [x] Useful `CHANGELOG.md`
  - [x] Clear `CONTRIBUTION.md`
- Automation:
  - [x] `Makefile`
- Testing:
  - [x] Coverage 80%+
  - [x] Unit test
  - [x] Real testing
- Examples:
  - [x] Example's test file
