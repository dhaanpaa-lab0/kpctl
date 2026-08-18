# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

`kpctl` (module `nexus-sds.com/kpctl`) is a Go CLI for managing K8S On Demand processing: selecting a Kubernetes context/namespace, syncing a GitOps repo, and applying `ExternalSecret` CRDs. Built with cobra (commands) and viper (config).

## Commands

```sh
go build -o dist/kpctl .          # build
task build                        # build via Taskfile (cleans dist/ first, output named dist/kpctl.<arch>.<os>)
go test ./...                     # run all tests
go test ./cmd -run TestName       # run a single test
go vet ./...
```

Requires Go 1.26+. No lint config beyond `go vet`.

## Architecture

**Layering**: `cmd/*.go` (cobra commands, thin) → `pkg/sysConfig` (viper-backed config/profile state) → `pkg/platform` (external systems: k8s API, git, xdg paths) → `pkg/core` (pure helpers, no external deps) → `pkg/helpers` (cobra-flag and bubbletea/pterm UI helpers).

Keep this layering when adding commands: a `cmd/*.go` file should read flags, call into `sysConfig`/`platform`, and print a result — business logic belongs in `pkg/`, not in the `RunE`/`Run` closure.

**Config and profiles** (`pkg/sysConfig`):
- Config is loaded by viper in `cmd/root.go`'s `initConfig()` (cobra `OnInitialize` hook), from `--config` or from `platform.GetConfigFolder()` (an XDG config dir) as `config.yaml`.
- Settings are grouped into named **profiles** under the viper key `profiles.<name>.<field>`, with `current_profile` marking the active one. Fields are always `k8s_context`, `k8s_cluster_ns`, `gitops_url`, `gitops_folder` (see `legacyProfileFields` in `pkg/sysConfig/profile.go`).
- `MigrateLegacyProfileConfig()` moves a pre-profile config (top-level fields, no `profiles` map) into a `default` profile in memory; `UpdateGlobalConfig()` (which calls `viper.WriteConfig`) then strips the now-redundant legacy top-level fields from disk via `stripLegacyProfileFields`. Any command that persists config should call `sysConfig.UpdateGlobalConfig()`, not `viper.WriteConfig()` directly, so this cleanup and the "config file not found → create it" fallback both run.
- **Known quirk**: viper lowercases map keys read back from a YAML file, so profile names round-tripped through disk become lowercase. Profile lookups (`ProfileExists`, `GetProfile`, `useProfile`, `deleteProfile`, `importProfile`) are case-sensitive against that map, so a profile created as `TeamA` may only be found as `teama` after a reload. This is a pre-existing, app-wide limitation, not specific to any one command — don't try to silently work around it in just one place.
- `ProfileExport`/`ExportProfile`/`ImportProfile` (also in `pkg/sysConfig/profile.go`) define the YAML file format used by the `exportProfile`/`importProfile` commands; keep the YAML tags in sync with `legacyProfileFields`/`Profile` if new profile fields are added.

**Global (non-profile) config**: flat key/value pairs under viper's `globals` map, managed by `pkg/sysConfig/global_config.go` and the `listGlobals`/`setGlobal`/`deleteGlobal` commands — unrelated to the profile system.

**Platform layer** (`pkg/platform`):
- `k8s_platform.go`: reads the local kubeconfig (via `client-go`/`clientcmd`) to list available contexts/namespaces for the interactive prompts in `initConfig`.
- `k8s_platform_ext_scrt.go`: builds and applies `ExternalSecret` CRD objects using the dynamic client (unstructured, since it's a CRD with no compiled Go type) — kept independent of the `external-secrets.io` Go module.
- `gitops_platform.go`: `SyncGitopsRepository` — clones or pulls the profile's `gitops_url` into `gitops_folder` (expanding `~`), handling three cases: existing `.git` (pull), empty target dir (clone), non-empty non-git dir (init + add remote + pull).
- `env_platform.go`: resolves the XDG config directory for the app.

**cmd conventions**: every file in `cmd/` carries the same MIT license header comment, registers itself via `init() { rootCmd.AddCommand(...) }`, uses `helpers.GetFlagAsString(cmd, "flagname")` to read flags, and many commands define kebab-case `Aliases` alongside the camelCase `Use` name (e.g. `syncGitops` aliased to `sync`/`sync-gitops`). Match this pattern for new commands.
