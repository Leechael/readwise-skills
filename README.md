# readwise-cli

`readwise-cli` is a command-line tool for working with Readwise highlights/books and Reader documents.

It supports:
- highlight CRUD and tag management
- book listing and tag management
- full highlight export with cursor pagination
- daily review highlights
- Reader document save/list/update/delete
- Reader tag listing

---

## Install

### Option A: Download from GitHub Releases

```bash
gh release list -R Leechael/readwise-skills
TAG="vX.Y.Z"
gh release download "$TAG" -R Leechael/readwise-skills --pattern "readwise-cli-*.tar.gz"
```

Extract the archive for your platform and place `readwise-cli` in your `PATH`.

### Option B: Build from source

```bash
git clone git@github.com:Leechael/readwise-skills.git
cd readwise-skills
go build -o bin/readwise-cli ./cmd/readwise
```

---

## Required configuration

Set credentials via environment variable:

```bash
export READWISE_API_TOKEN="<token>"
```

Alternatively, pass `--token` on every invocation.

Validate setup before use:

```bash
readwise-cli status
readwise-cli status --json
```

---

## Commands

### v2 — Highlights & Books

- `status` — check credentials and API connectivity
- `highlight list` — list highlights with filters
- `highlight get <id>` — get a single highlight
- `highlight create` — create highlights (flags or stdin)
- `highlight update <id>` — update highlight fields
- `highlight delete <id>` — delete a highlight
- `highlight tag list|add|update|delete` — manage highlight tags
- `book list` — list books/sources with filters
- `book get <id>` — get a single book
- `book tag list|add|update|delete` — manage book tags
- `export` — full highlight export with cursor pagination
- `review` — daily review highlights

### v3 — Reader

- `reader list` — list Reader documents with filters
- `reader save` — save a URL/document to Reader
- `reader update <id>` — update document metadata
- `reader delete <id>` — delete a document
- `reader tag list` — list all Reader tags

### Output modes

- Parseable output is available via `--json`.
- Human-readable output is available via `--plain`.
- `--json` and `--plain` are mutually exclusive.
- `--jq` requires `--json` for filtered JSON output.

---

## Usage examples

```bash
# status
readwise-cli status
readwise-cli status --json

# highlights
readwise-cli highlight list --json
readwise-cli highlight list --book-id 123 --page-size 10 --json
readwise-cli highlight get 456 --plain
readwise-cli highlight create --text "Quote" --title "Book" --author "Author" --json
readwise-cli highlight update 456 --note "updated note" --json
readwise-cli highlight delete 456

# highlight tags
readwise-cli highlight tag list 456 --json
readwise-cli highlight tag add 456 --name "favorite" --json

# books
readwise-cli book list --category books --json
readwise-cli book get 123 --json

# export & review
readwise-cli export --json
readwise-cli export --updated-after "2025-01-01T00:00:00Z" --json
readwise-cli review --json

# reader documents
readwise-cli reader list --json
readwise-cli reader list --category article --location new --json
readwise-cli reader save --url "https://example.com/article" --json
readwise-cli reader update abc123 --title "New Title" --location archive
readwise-cli reader delete abc123

# reader tags
readwise-cli reader tag list --json
```

---

## Install the Agent Skill

This repository also ships an Agent Skill package under `skills/readwise`.

Install with:

```bash
npx skills add Leechael/readwise-skills
```

After installation, your agent can load and use the `readwise` skill instructions.

---

## Recommended secret handling

Use 1Password CLI to inject credentials at runtime:

- https://developer.1password.com/docs/service-accounts/use-with-1password-cli

Example:

```bash
op run --env-file=.env -- readwise-cli status
op run --env-file=.env -- readwise-cli highlight list --json
```
