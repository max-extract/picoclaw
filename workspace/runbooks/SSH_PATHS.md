# SSH Paths For Telegram And Local Execution

There is no dedicated `ssh/` folder in this repo.
SSH execution is handled by script env/flags in `scripts/me_ops_exec.sh`.

## Why path errors happen

Telegram-triggered sessions can run from a cwd that is not your MaxExtract root.
Hardcoded absolute paths then fail.

## Path-safe command pattern

From `picoclaw-deploy`, use the wrapper:

```bash
MAXEXTRACT_USE_SSH=1 ./workspace/bin/me.sh me_bots_inventory.sh --context mycoolify --mode all --json
MAXEXTRACT_USE_SSH=1 ./workspace/bin/me.sh me_bots_api_state.sh --context mycoolify --mode all --json
MAXEXTRACT_USE_SSH=1 ./workspace/bin/me.sh me_bots_digest.sh --context mycoolify --mode all --days auto
```

## If root cannot be auto-detected

```bash
export MAXEXTRACT_ROOT="/path/to/maxextract"
MAXEXTRACT_USE_SSH=1 ./workspace/bin/me.sh me_bots_inventory.sh --context mycoolify --mode all --json
```

## Force local mode

```bash
MAXEXTRACT_USE_SSH=0 ./workspace/bin/me.sh me_bots_inventory.sh --context mycoolify --mode all --json
```

## SSH auth notes

If you see `Permission denied (publickey,password)`, path resolution is working but SSH credentials are not.

Use one of these:

```bash
export MAXEXTRACT_SSH_USER=root
export MAXEXTRACT_SSH_HOST=84.32.104.5
export MAXEXTRACT_SSH_PORT=22
export MAXEXTRACT_SSH_TARGET="root@84.32.104.5"
```

To fail fast instead of local fallback:

```bash
export MAXEXTRACT_SSH_STRICT=1
```
