# push-treasurenet-mainnet-artifacts-to-s3 Workflow Guide

The workflow defined in `.github/workflows/push-treasurenet-mainnet-artifacts-to-s3.yml` compiles multi-arch treasurenet binaries on a self-hosted runner, initializes node data, packages the node directories, and uploads every artifact to `s3://treasurenet-node-binary-us-west-1`. Use the checklist below to make sure it can run inside `treasurenet2` without manual intervention.

## Self-hosted runner requirements
- **OS and permissions**: Ubuntu (or compatible) runner with `sudo` access. The workflow writes into `/usr/bin` and `/data`, so both paths must permit the runner user to create and overwrite files.
- **Pre-installed tooling**: Go 1.18, Node.js 8.x, GNU Make, tar, AWS CLI v2, and access to `aws sts`.
- **Filesystem layout**: `/data` must exist before the workflow runs; the job rewrites node sub-directories underneath it.

## Repository dependencies
- `scripts/shell/` must contain:
  - `init_nodes.sh`
  - `init_node_template.sh`
  - `init_seednode.sh`
  - `new_genesis.sh`
  - `change_app.toml.sh`
  - `change_config.toml.sh`
- These scripts already exist at `treasurenet2/scripts/shell` (copied from the source repo). Keep their interfaces stable if you edit them because the workflow expects the current CLI flags and paths.

## GitHub secrets
Create or update the secrets below in the `treasurenet2` repository:
- `keyring_secret`: passphrase required for node initialization and genesis generation.
- AWS parameters such as `AWS_ROLE_ARN`, bucket names, and regions are hard-coded in the workflow `env` block. Update both the env values and the corresponding AWS resources if you need a different account or bucket.

## Triggers
- `push` events: any push to the `main` branch kicks off a run.
- `workflow_dispatch`: the Actions tab now exposes a **Run workflow** button named after this workflow. An optional `build_target` input lets you annotate manual runs (default `mainnet`).

Push to `main` or manually dispatch the workflow after the above prerequisites are satisfied. If a run fails, inspect the self-hosted runner logs first to confirm the runner meets the dependency and permission requirements.
