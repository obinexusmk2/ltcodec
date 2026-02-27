# Publishing ltcodec to the Go Registry

## Pre-flight checklist

Before running these commands, confirm:

- [ ] `github.com/obinexusmk2/ltcodec` repository is **public** on GitHub
- [ ] Your local remote points to it: `git remote -v`
- [ ] All files are committed: `git status`

---

## Step 1 — Tidy and verify

```powershell
cd C:\Users\OBINexus\Projects\ltcodec\ltcodec

# Remove stale uuid dep (not imported), regenerate clean go.sum
go mod tidy

# Confirm the module builds and the CLI works
go build ./...
go run cmd/ltcodec/main.go state
```

---

## Step 2 — Commit everything

```powershell
git add go.mod go.sum doc.go LICENSE CHANGELOG.md
git add pkg/gps/gps.go pkg/identity/identity.go
git add pkg/spacetime/spacetime.go pkg/codec/codec.go
git add cmd/ltcodec/main.go
git commit -m "chore: prepare v1.0.0 for Go registry publication"
```

---

## Step 3 — Tag v1.0.0

The Go module proxy indexes **annotated tags** that follow `vMAJOR.MINOR.PATCH`.

```powershell
git tag -a v1.0.0 -m "ltcodec v1.0.0 — LTF spacetime-anchored codec, OBINexus Computing"
```

---

## Step 4 — Push branch + tag

```powershell
git push origin main
git push origin v1.0.0
```

---

## Step 5 — Request indexing from the Go proxy

The proxy picks up the tag automatically within minutes, but you can
force immediate indexing:

```powershell
# Open this URL in your browser (or curl it):
# https://sum.golang.org/lookup/github.com/obinexusmk2/ltcodec@v1.0.0
# https://pkg.go.dev/github.com/obinexusmk2/ltcodec@v1.0.0
```

Or trigger it from any machine with Go installed:

```bash
GOPROXY=https://proxy.golang.org go list -m github.com/obinexusmk2/ltcodec@v1.0.0
```

---

## Step 6 — Verify on pkg.go.dev

Within ~5 minutes the package page will appear at:

**https://pkg.go.dev/github.com/obinexusmk2/ltcodec**

It will show:
- The `doc.go` overview as the module description
- All four sub-packages (`gps`, `identity`, `spacetime`, `codec`)
- The MIT licence badge
- `go install github.com/obinexusmk2/ltcodec/cmd/ltcodec@latest` install command

---

## What others can do after publication

```bash
# Install the CLI
go install github.com/obinexusmk2/ltcodec/cmd/ltcodec@latest

# Import the library
import "github.com/obinexusmk2/ltcodec/pkg/codec"
```

---

## Next version plan

| Version | Feature |
|---------|---------|
| v1.1.0  | NSIGII Trident channel writer |
| v1.2.0  | SMTP relay for LTF legal packets |
| v2.0.0  | CGo bridge — riftlang `.so.a` linkage |
