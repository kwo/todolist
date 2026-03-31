---
id: todo-qwd4
title: support Homebrew and Linuxbrew tap installation from release artifacts
status: done
priority: 2
parents:
    - todo-96ke
depends:
    - todo-rvwr
createdAt: "2026-03-28T18:50:55Z"
lastModified: "2026-03-31T19:03:26Z"
---

# User Story: support Homebrew and Linuxbrew tap installation from release artifacts

Publish release assets in a layout that works for both macOS Homebrew and Linuxbrew formulas in a personal tap.

## Acceptance criteria

- release asset naming is stable and formula-friendly
- tap formula can select macOS and Linux archives by CPU architecture
- install path works for Homebrew and Linuxbrew users from the same tap
- usage or release docs explain how to update the tap formula

## Exact plan

1. Update the GitHub release workflow to keep the current asset naming convention:
   - `todolist_v<version>_<goos>_<goarch>.tar.gz`
   - example: `todolist_v1.2.3_darwin_amd64.tar.gz`
2. Change release packaging for all platforms to be binary-only:
   - macOS/Linux archives contain only `todolist`
   - Windows archives contain only `todolist.exe`
   - archives should extract the binary at the archive root, not inside a versioned directory
3. Keep Homebrew/Linuxbrew support scoped to the exact released targets:
   - `darwin/amd64`
   - `darwin/arm64`
   - `linux/amd64`
   - `linux/arm64`
4. Add a maintainer helper script at `scripts/homebrew-formula` that prints a complete ready-to-commit `todolist.rb` formula to stdout.
5. The helper script should:
   - use `gh` from within this repo
   - default to the latest published non-draft, non-prerelease GitHub release
   - support an explicit `--tag vX.Y.Z` override
   - inspect the selected release with `gh release list` and `gh release view`
   - download and parse the `SHA256SUMS` release asset
   - require all four macOS/Linux archives and `SHA256SUMS`; otherwise fail with a clear error
   - fail clearly if `gh` is missing or unusable
6. Generate a formula named `todolist.rb` with class `Todolist` using the existing `homebrew-tools` style:
   - `desc "Local-first CLI for managing todos stored as Markdown files"`
   - `homepage "https://github.com/kwo/todolist"`
   - `on_macos` and `on_linux` blocks
   - exact-match support for macOS/Linux amd64/arm64 only
   - clear failure message for unsupported platform/architecture combinations
   - no `test do` block
7. Add a new maintainer document at `BUILD.md` that covers only the release/build/tap update workflow:
   - tagged release expectations
   - release artifact naming/layout
   - how to run `scripts/homebrew-formula`
   - how to update `kwo/homebrew-tools`
   - how users install from `brew tap kwo/tools`

## Decisions

- keep current release asset naming instead of renaming to match existing tap formulas
- use a flat archive layout with only the binary in each release archive
- apply binary-only packaging to Windows too for consistency
- use GitHub Releases from `kwo/todolist` as the formula download source
- document tap maintenance in `BUILD.md`, not `README.md`
- do not modify `../homebrew-tools` as part of this todo; this story prepares the repo and generator only

## Out of scope

- creating or committing `todolist.rb` in `kwo/homebrew-tools`
- automating commits or pull requests against the tap repo
