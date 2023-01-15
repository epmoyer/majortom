# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## 1.4.0 - 2023-01-15
### Changed
- `-init` will create parent directories if they don't exist.
    - NOTE: for a brief time there was a `-force` option (to force `-init` without prompting) but it turned out that prompts were not visible to thn user because thy occur within the (output capturing) shell function to(). So now `-init` creates the file (and paths) unconditionally, so long as the config file does not already exist.t

## 1.3.0 - 2023-01-14
### Added
- Show `(DOES NOT EXIST)` for paths that do not exist if in `-no-color` mode.

## 1.2.0 - 2023-01-14
### Added
- Color options `-color=` (`16`, `256`, `16m`) and `-no-color`
### Changed
- Use `gtar` for tar archive creation, if available

## 1.1.0 - 2023-01-06
### Changed
- Adopt app name `majortom`

## 1.0.0 - 2023-01-06
Initial Releare