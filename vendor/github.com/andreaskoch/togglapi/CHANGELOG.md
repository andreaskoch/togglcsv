# Change Log
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

## [v0.4.1] - 2016-10-01

Fix return values of create functions

### Fixed
- Deserialize the created time entry, client and project entities correctly.

## [v0.4.0] - 2016-09-25

Rate limiting

### Added
- Add a request rate limiter to the Toggl REST client which prevents API quota errors. Only one request per second will be issued in the future (see: https://github.com/toggl/toggl_api_docs/commit/5d34c09fc17f2ba5ab97d2b04e97059dbba34ba0).

## [v0.3.0] - 2016-09-24

### Added
- Include the client ID in the project model

## [v0.2.0] - 2016-09-19

Rename repositories to "API"

### Changed
- Rename sub APIs for workspaces, clients, projects and time entries from "...Repository" to "...API"

## [v0.1.3] - 2016-09-18

Fix errors in example code

### Fixed
- Fix import cycle in API example
- Fix ProjectRepository example code

## [v0.1.2] - 2016-09-18

Improve godoc

### Added
- Add some godoc code examples

### Fixed
- Remove duplicate godoc package descriptions

## [v0.1.1] - 2016-09-18

Fix travis-ci build

### Fixed
- Add missing github.com/pkg/errors dependency

## [v0.1.0] - 2016-09-18

First working version

### Added
- Base implementation of the Toggl API methods for creating and retrieving clients, workspaces, projects and time entries
