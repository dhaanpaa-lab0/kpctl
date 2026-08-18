# kpctl

K8S On Demand processing CLI — a utility for managing K8S On Demand processing.

## Installation

Requires Go 1.26+.

```sh
go build -o dist/kpctl .
```

Or with [Task](https://taskfile.dev):

```sh
task build
```

## Configuration

kpctl reads configuration from `$HOME/.kpctl.yaml` by default, or from a file passed via `--config`.

Settings are grouped into named **profiles**, stored in the config file under a `profiles` map, with `current_profile` marking which one is active. Config files written before profiles existed (with top-level `k8s_context`, `k8s_cluster_ns`, `gitops_url`, `gitops_folder` keys) are migrated automatically into a `default` profile the first time kpctl runs; the migrated layout is written back the next time you run a command that saves config (e.g. `initConfig`).

Run `initConfig` to set up or update a profile interactively:

```sh
kpctl initConfig
```

This prompts for a profile name (defaulting to the active profile), then:
- `k8s_context` — the K8S context to use
- `k8s_cluster_ns` — the K8S cluster namespace
- `gitops_url` — the GitOps repository URL
- `gitops_folder` — the local folder to sync GitOps into

and makes the configured profile active.

## Commands

| Command                           | Aliases                        | Description                                                                |
|-----------------------------------|--------------------------------|----------------------------------------------------------------------------|
| `initConfig`                      |                                | Setup or update a configuration profile                                    |
| `listProfiles`                    | `listProfile`, `list-profiles` | List configuration profiles in a table, marking the active one             |
| `useProfile --name <n>`           |                                | Set the active configuration profile                                       |
| `deleteProfile --name <n>`        |                                | Delete a configuration profile                                             |
| `listGlobals`                     | `listGlobal`, `list-globals`   | List global configuration variables in a table                             |
| `setGlobal --key <k> --value <v>` |                                | Set a global configuration variable                                        |
| `deleteGlobal --key <k>`          |                                | Delete a global configuration variable                                     |
| `syncGitops`                      | `sync`, `sync-gitops`          | Sync the GitOps repository (from the active profile) into the local folder |
| `exportProfile --file <f>`        | `export-profile`               | Export a configuration profile to a YAML file                              |
| `importProfile --file <f>`        | `import-profile`               | Import a configuration profile from a YAML file                            |
| `externalSecret`                  | `external-secret`               | Build an ExternalSecret manifest, prompting for missing flags; write it with `--file` and/or apply it to the cluster with `--apply` |

Run `kpctl --help` or `kpctl <command> --help` for full usage details.

## License

MIT — see [LICENSE](LICENSE).
