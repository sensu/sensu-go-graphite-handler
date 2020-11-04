# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed
- Improved test coverage
- Updated SDK to v0.11.0
- Fix entity name dot replacement to all dots

## [ 0.4.0] - 2020-02-07

### Changed
- Fixed goreleaser deprecated archive to use archives
- Converted from Travis CI to GitHub Actions

### Added
- Added classic parity support for `--prefix-source`.
- Added classic parity support for `--annotation-prefix`.
- Added classic parity support for annotation metrics of check attributes.

### Fixed
- Fixed a panic that could arise from statsd metrics.
- Implemented some go styling best practices.

## [ 0.3.1] - 2019-12-17

### Changed
- Reformatted README for [Plugin Style Guide](https://github.com/sensu-plugins/community/blob/master/PLUGIN_STYLEGUIDE.md)

## [ 0.3.0] - 2019-08-22

### Changed
- Fixed documentation to be more consistent with naming
- Fixed default port to be 2003
- Compile with 1.12.9
- Switch to Go Modules
- Fixed goreleaser to remove v in version

## [ 0.2.1] - 2019-07-15

### Changed
- Fixed zeroing out prefix if -n --no-prefix specified

## [ 0.2.0] - 2019-07-15

### Added
- Added no-prefix option to allow bare metric names

## [ 0.1.1] - 2019-05-10

### Changed
- Fixed default port to be 2013

## [ 0.1.0] - 2019-04-16
- Initial release
