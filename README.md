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

Run `initConfig` to set up or update the configuration interactively:

```sh
kpctl initConfig
```

This prompts for:
- `k8s_context` — the K8S context to use
- `k8s_cluster_ns` — the K8S cluster namespace
- `gitops_url` — the GitOps repository URL
- `gitops_folder` — the local folder to sync GitOps into

## Commands

| Command                           | Aliases                      | Description                                      |
|-----------------------------------|------------------------------|--------------------------------------------------|
| `initConfig`                      |                              | Setup or update the CLI configuration            |
| `listGlobals`                     | `listGlobal`, `list-globals` | List global configuration variables in a table   |
| `setGlobal --key <k> --value <v>` |                              | Set a global configuration variable              |
| `deleteGlobal --key <k>`          |                              | Delete a global configuration variable           |
| `syncGitops`                      | `sync`, `sync-gitops`        | Sync the GitOps repository into the local folder |

Run `kpctl --help` or `kpctl <command> --help` for full usage details.

## License

MIT — see [LICENSE](LICENSE).
