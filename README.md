# ticky — TickTick in your terminal.

[\[日本語\]](README_ja.md)

A CLI tool for managing TickTick tasks via the Open API. OAuth 2.0 authentication, project/task CRUD, tag management, and script-friendly design.

## Features

- **Tasks** — create, get, update, complete, delete tasks with priority, due dates, and tags
- **Projects** — list and view project details
- **Tags** — aggregate tags across all projects
- **Flexible due dates** — `today`, `tomorrow`, `+3d`, `YYYY-MM-DD`
- **Priority levels** — `none`, `low`, `medium`, `high`
- **Multiple output formats** — human-readable text, JSON, and TSV
- **OAuth 2.0** — browser-based login with token auto-refresh

## Installation

### Homebrew

```bash
brew install tackeyy/tap/ticky
```

### Go

```bash
go install github.com/tackeyy/ticky@latest
```

### Build from source

```bash
git clone https://github.com/tackeyy/ticky.git
cd ticky
go build -o ticky .
```

## Quick Start

### 1. Create a TickTick App

1. Go to [TickTick Developer Portal](https://developer.ticktick.com/manage) and click **+Create App**
2. Name your app (e.g., `ticky`)

### 2. Configure OAuth Settings

In your app settings:

| Setting | Value |
|---|---|
| Redirect URL | `http://localhost:18080/callback` |
| Scopes | `tasks:read`, `tasks:write` |

### 3. Set Environment Variables

```bash
export TICKTICK_CLIENT_ID=your_client_id
export TICKTICK_CLIENT_SECRET=your_client_secret
```

### 4. Login and Run

```bash
ticky auth login
ticky tasks list --json
```

## Commands

### `auth login` — Login via OAuth

```bash
ticky auth login
```

Opens a browser for TickTick authorization. Token is saved to `~/.config/ticky/token.json`.

### `auth status` — Check authentication status

```bash
ticky auth status [--json] [--plain]
```

### `auth logout` — Remove saved token

```bash
ticky auth logout
```

### `tasks list` — List tasks

```bash
ticky tasks list [--project <id>] [--json] [--plain]
```

| Flag | Required | Description |
|---|---|---|
| `--project <id>` | No | Project ID (default: Inbox) |

### `tasks get` — Get task details

```bash
ticky tasks get <task_id> --project <id> [--json] [--plain]
```

| Flag | Required | Description |
|---|---|---|
| `<task_id>` | Yes | Task ID |
| `--project <id>` | Yes | Project ID |

### `tasks create` — Create a task

```bash
ticky tasks create --title <title> [--project <id>] [--content <text>] [--priority <level>] [--due <date>] [--tags <tags>] [--json] [--plain]
```

| Flag | Required | Description |
|---|---|---|
| `--title <title>` | Yes | Task title |
| `--project <id>` | No | Project ID (default: Inbox) |
| `--content <text>` | No | Task content/description |
| `--priority <level>` | No | `none`, `low`, `medium`, `high` |
| `--due <date>` | No | `today`, `tomorrow`, `+3d`, `YYYY-MM-DD` |
| `--tags <tags>` | No | Comma-separated tags |

Examples:

```bash
ticky tasks create --title "Review PR" --priority high --due tomorrow
ticky tasks create --title "Buy milk" --tags "shopping,personal" --json
```

### `tasks update` — Update a task

```bash
ticky tasks update <task_id> --project <id> [--title <title>] [--content <text>] [--priority <level>] [--due <date>] [--clear-due] [--tags <tags>] [--add-tags <tags>] [--remove-tags <tags>] [--json] [--plain]
```

| Flag | Required | Description |
|---|---|---|
| `<task_id>` | Yes | Task ID |
| `--project <id>` | Yes | Project ID |
| `--title <title>` | No | New title |
| `--content <text>` | No | New content |
| `--priority <level>` | No | `none`, `low`, `medium`, `high` |
| `--due <date>` | No | New due date |
| `--clear-due` | No | Clear the due date |
| `--tags <tags>` | No | Replace all tags |
| `--add-tags <tags>` | No | Add tags |
| `--remove-tags <tags>` | No | Remove tags |

Examples:

```bash
ticky tasks update abc123 --project def456 --priority high --due +3d
ticky tasks update abc123 --project def456 --add-tags "urgent" --json
```

### `tasks complete` — Complete a task

```bash
ticky tasks complete <task_id> --project <id> [--json] [--plain]
```

| Flag | Required | Description |
|---|---|---|
| `<task_id>` | Yes | Task ID |
| `--project <id>` | Yes | Project ID |

### `tasks delete` — Delete a task

```bash
ticky tasks delete <task_id> --project <id> [--json] [--plain]
```

| Flag | Required | Description |
|---|---|---|
| `<task_id>` | Yes | Task ID |
| `--project <id>` | Yes | Project ID |

### `projects list` — List projects

```bash
ticky projects list [--json] [--plain]
```

### `projects get` — Get project details

```bash
ticky projects get <project_id> [--json] [--plain]
```

| Flag | Required | Description |
|---|---|---|
| `<project_id>` | Yes | Project ID |

### `tags list` — List all tags

```bash
ticky tags list [--json] [--plain]
```

Aggregates tags from tasks across all projects, sorted by usage count.

## Configuration

### Environment Variables

| Variable | Required | Description |
|---|---|---|
| `TICKTICK_CLIENT_ID` | Yes | OAuth client ID |
| `TICKTICK_CLIENT_SECRET` | Yes | OAuth client secret |
| `TICKTICK_ACCESS_TOKEN` | No | Direct access token (skips token file, useful for CI/agents) |

### Token Storage

After `ticky auth login`, the OAuth token is saved to `~/.config/ticky/token.json` with `0600` permissions. Token refresh is handled automatically.

If `TICKTICK_ACCESS_TOKEN` is set, the token file is ignored.

## Output Formats

### Text (default)

```
abc123def456789012345678 Review PR [high] (due: 2026-02-12) #work
```

### JSON (`--json`)

```json
[
  {
    "id": "abc123def456789012345678",
    "projectId": "inbox123",
    "title": "Review PR",
    "priority": 5,
    "dueDate": "2026-02-12T14:59:59.000+0000",
    "tags": ["work"]
  }
]
```

### TSV (`--plain`)

```
abc123def456789012345678	inbox123	Review PR	high	2026-02-12T14:59:59.000+0000	work
```

## Development

```bash
# Build
go build -o ticky .

# Run tests
go test -v ./...

# Run tests with coverage
go test -cover ./...
```

See [docs/TESTING.md](docs/TESTING.md) for detailed testing guide.

## License

MIT

## Contributing

Contributions are welcome! Please read our [Contributing Guide](CONTRIBUTING.md) before submitting a Pull Request.

See also:
- [Testing Guide](docs/TESTING.md)
- [Code of Conduct](CODE_OF_CONDUCT.md)

## Links

- [GitHub Repository](https://github.com/tackeyy/ticky)
- [TickTick Open API Documentation](https://developer.ticktick.com/api)
