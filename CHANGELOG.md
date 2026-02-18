# Changelog

All notible changes to SOLACE.sys will be documented in this file.
this format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html)

## [Unreleased]

### Added

- Initial repository setup and Go module initialisation.
- Added Smoke Test for immediate CI/CD pipeline testing
- Github Actions CI/CD pipeline with `golangci-lint` and strict race detection.

## [0.1.0] - 2026-02-16

### Added

- Installed go-fiber for the backend API
- Setup a simple backend API
- Embedded index html file and the static files
- Added htmx import tag to index.html
- initalised the frontend serving

## [0.2.1] - 2026-02-17

### Added

- Installed Pion for native WebRTC functionaltiy
- Added SDP handshake from frontend and backend
- Added creating an answer for the frontend
- Added DataChannel listener for peerConnection
- Added timeout for ICE gathering
- Added Detailed state feedback for peerConnection State

### Fixed

- Kept peerConnection alive after the scope of the webrtc/offer endpoint lifetime
- Waited for ICE candidates to be acquired
- A Few spelling mistakes

### Replaced

- Replaced the global hashmap for peerConnections to sync.Map to avoid concurrent map writes panic
-

## [0.3.0] - 2026-02-17

### Added

- Improved the code quality on the frontend
- Encapsulated the free code for handshake and event listeners
- data incoming and outgoing event listeners established
