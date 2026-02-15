---
name: readwise
description: Readwise and Reader operations via readwise-cli. Use this skill for highlight/book retrieval, tag management, daily review, full export, and Reader document save/list/update/delete with environment-injected credentials.
compatibility: Requires readwise-cli in PATH and READWISE_TOKEN environment variable.
---

# Readwise Skill

Use this skill to perform Readwise (highlights/books) and Reader (documents) workflows through `readwise-cli`.

## Prerequisites

### 1) Ensure `readwise-cli` is installed

Preferred: install from GitHub Releases of this repository.

```bash
# 1) Inspect releases
gh release list -R Leechael/readwise-skills

# 2) Download latest artifacts
gh release download -R Leechael/readwise-skills --pattern 'readwise-cli-*.tar.gz'

# 3) Extract your platform artifact and install binary
tar -xzf readwise-cli-$(uname -s | tr '[:upper:]' '[:lower:]')-$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/').tar.gz
install -m 0755 readwise-cli /usr/local/bin/readwise-cli
```

Quick check:

```bash
readwise-cli --help
```

### 2) Configure credentials

Required environment variable:

```bash
export READWISE_TOKEN="<token>"
```

Alternatively, pass `--token` on every invocation.

### 3) Verify login/config status before operations

Use the built-in auth check command:

```bash
readwise-cli auth check
```

If credentials are missing or invalid, this command will fail with guidance. Do not continue with operations until `auth check` is successful.

## Recommended Credential Management (1Password CLI)

Prefer running `readwise-cli` with the 1Password CLI (`op`) so credentials are injected at runtime instead of stored in shell profiles.

Reference:
- https://developer.1password.com/docs/service-accounts/use-with-1password-cli

Example pattern:

```bash
op run --env-file=.env -- readwise-cli auth check
op run --env-file=.env -- readwise-cli highlight list --json
```

(Your `.env` should define `READWISE_TOKEN` with a 1Password secret reference.)

## Command Mapping

### v2 — Highlights & Books

- Check auth: `readwise-cli auth check`
- List highlights: `readwise-cli highlight list`
- Get highlight: `readwise-cli highlight get <id>`
- Create highlights: `readwise-cli highlight create`
- Update highlight: `readwise-cli highlight update <id>`
- Delete highlight: `readwise-cli highlight delete <id>`
- Highlight tags: `readwise-cli highlight tag list|add|update|delete`
- List books: `readwise-cli book list`
- Get book: `readwise-cli book get <id>`
- Book tags: `readwise-cli book tag list|add|update|delete`
- Export highlights: `readwise-cli export`
- Daily review: `readwise-cli review`

### v3 — Reader

- List documents: `readwise-cli reader list`
- Save document: `readwise-cli reader save`
- Update document: `readwise-cli reader update <id>`
- Delete document: `readwise-cli reader delete <id>`
- List Reader tags: `readwise-cli reader tag list`

## Recommended Workflow

1. Run `readwise-cli auth check` first.
2. Prefer read commands first: `highlight list`, `book list`, `reader list`.
3. Use filters to narrow results: `--category`, `--book-id`, `--updated-after`.
4. Pipe JSON output through `--jq` for scripting: `--json --jq '.results[]'`.
5. Use write commands when needed: `highlight create`, `reader save`.

## Usage Examples

### 1) Check authentication

```bash
readwise-cli auth check
readwise-cli auth check --json
```

### 2) List and filter highlights

```bash
readwise-cli highlight list --json
readwise-cli highlight list --book-id 12345 --page-size 10 --json
readwise-cli highlight list --json --jq '.results[] | {id, text}'
```

### 3) Get a single highlight

```bash
readwise-cli highlight get 67890 --json
readwise-cli highlight get 67890 --plain
```

### 4) Create a highlight

```bash
readwise-cli highlight create --text "Important quote" --title "Book Title" --author "Author Name" --json
echo '{"highlights":[{"text":"From stdin","title":"Source"}]}' | readwise-cli highlight create --stdin --json
```

### 5) Manage highlight tags

```bash
readwise-cli highlight tag list 67890 --json
readwise-cli highlight tag add 67890 --name "favorite" --json
readwise-cli highlight tag update 67890 111 --name "must-read"
readwise-cli highlight tag delete 67890 111
```

### 6) List and get books

```bash
readwise-cli book list --category books --json
readwise-cli book list --json --jq '.results[] | {id, title, author}'
readwise-cli book get 12345 --json
```

### 7) Export and daily review

```bash
readwise-cli export --json
readwise-cli export --updated-after "2025-01-01T00:00:00Z" --json
readwise-cli review --json
readwise-cli review --json --jq '.highlights[] | {text, title}'
```

### 8) Reader: list and save documents

```bash
readwise-cli reader list --json
readwise-cli reader list --category article --location new --json
readwise-cli reader save --url "https://example.com/article" --json
readwise-cli reader save --url "https://example.com" --title "Custom Title" --tag reading --tag tech --json
```

### 9) Reader: update and delete documents

```bash
readwise-cli reader update abc123 --title "New Title" --location archive
readwise-cli reader delete abc123
```

### 10) Reader tags

```bash
readwise-cli reader tag list --json
readwise-cli reader tag list --plain
```

## Error Handling Rules

- Missing credentials: explicitly report missing `READWISE_TOKEN`.
- API failures: include HTTP status code and response body.
- Not found: clearly include the identifier that was requested.
- Auth failures: exit code 2. Not found: exit code 3. General errors: exit code 1.

## Output Rules

- Use `--json` for machine-readable JSON output.
- Use `--plain` for tab-separated stable output.
- Use `--jq` with `--json` for filtered JSON output.
- Human hints are printed to stderr, data to stdout.
- Never invent Readwise data; only report real command results.
