---
id: todo-rvwr
title: tag-driven release workflow in GitHub Actions
status: done
priority: 5
parents:
    - todo-96ke
depends:
    - todo-gag2
createdAt: "2026-03-28T18:50:47Z"
lastModified: "2026-03-30T21:02:28Z"
---

# User Story: tag-driven release workflow in GitHub Actions

Add a GitHub Actions workflow that responds to a newly pushed version tag by building release binaries for macOS, Linux, and Windows on amd64 and arm64 and creating the corresponding GitHub Release.

## User story

As a maintainer,
I want pushing a new version tag to trigger the full release workflow in GitHub Actions,
so that builds happen automatically from the tagged revision and the release record is created without manual steps.

## Goal

This story covers the GitHub Actions release workflow end to end. When a new release tag is pushed to the repository, GitHub Actions should run the build inside the workflow, compile `todolist` for the supported platform matrix, package the outputs, generate checksums, and create the GitHub Release for that tag.

## Technical direction

- trigger on newly pushed semantic version tags
- perform all release builds inside GitHub Actions in response to the tag push
- use a GitHub Actions matrix for:
  - `darwin/amd64`
  - `darwin/arm64`
  - `linux/amd64`
  - `linux/arm64`
  - `windows/amd64`
  - `windows/arm64`
- build the CLI with release version information injected so `todolist version` reports the pushed tag version in produced binaries
- package each binary into a conventional archive format appropriate for the target platform
- use stable, predictable artifact names that include version, OS, and architecture
- produce a checksum file covering every generated archive
- create a GitHub Release whose tag/version matches the pushed tag
- keep the workflow fully tag-driven so the maintainer action is pushing the tag

## Acceptance criteria

### Tag-triggered release workflow

Given:

- a maintainer pushes a new release tag to the repository

When:

- GitHub Actions receives the tag push event

Then:

- the release workflow starts automatically
- the build is performed inside GitHub Actions from the tagged revision
- ordinary branch pushes do not run this release path

### Cross-platform build matrix

Given:

- the release workflow is triggered by a new tag

Then:

- binaries are built for `darwin/amd64`, `darwin/arm64`, `linux/amd64`, `linux/arm64`, `windows/amd64`, and `windows/arm64`
- each matrix entry produces a runnable `todolist` binary for its target platform

### Packaged release archives

Given:

- a target binary has been built inside the workflow

Then:

- the workflow packages it into a release archive
- archive names are deterministic and include enough metadata for users to identify the correct platform build
- Windows artifacts use the correct executable naming conventions for that platform

### Checksums

Given:

- all release archives have been generated

Then:

- the workflow emits SHA-256 checksums for all archives
- the checksum output is retained as part of the release workflow output

### GitHub Release creation

Given:

- the workflow is running for a newly pushed release tag

Then:

- a GitHub Release is created automatically for that tag
- the created release corresponds to the same version that triggered the build

### Maintainer workflow

Given:

- the repository is configured correctly for releases

Then:

- pushing a new git tag is sufficient to trigger the build and create the GitHub Release
- no manual local build step is required

## Decisions

1. The release workflow is initiated by pushing a new version tag to the repository.
2. All release builds are performed inside GitHub Actions.
3. The supported release matrix is macOS, Linux, and Windows on amd64 and arm64.
4. Checksums use SHA-256.
5. This story combines artifact building and GitHub Release creation into one tag-driven workflow.

## Out of scope

- manual local release builds
- Homebrew/Linuxbrew formula publishing
