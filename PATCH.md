# Pass-through capture patch (fork branch, unmerged)

This branch is CLIProxyAPI with one small fork patch on top. It is **not merged** into this fork's `main`; merging is Kelly's decision.

## What the patch does

- **Forwards the caller's own credential.** A `claude-api-key` entry whose `api-key` is the sentinel `passthrough-client-credential` sends the incoming request's own `Authorization: Bearer` token (or `X-Api-Key` when there is no bearer) upstream, instead of a stored key. The proxy holds no sign-in and no key. (`internal/runtime/executor/claude_passthrough.go`, hooked in `claude_executor_request.go`.)
- **Never cools down the pass-through entry.** Each request carries a different caller's credential, so one caller's 401 or 429 must not bench the entry for everyone. Results for that entry skip the cooldown bookkeeping and only fire the result hook. (`sdk/cliproxy/auth/conductor_cooldown.go`.)

With `request-log: true`, every proxied call is written to `<auth-dir>/logs/` as a capture file. The request logger masks auth headers already (for example `Bearer dumm...lpha`).

## Config

```yaml
claude-api-key:
  - api-key: "passthrough-client-credential"
```

Use an empty `auth-dir` of its own, so no stored sign-in loads as a second Claude entry and gets round-robined in.

## Tested

Built and tested with `go test ./...` on this branch, then run on scratch ports against a local dummy upstream (dummy tokens, no Anthropic call). Patched: each caller's token reached the upstream, a 401 for one caller didn't affect the next call, and one capture file was written per call. Stock: the literal sentinel went upstream.
