# ticky — TickTick をターミナルから。

[\[English\]](README.md)

TickTick Open API を使ってタスクを管理する CLI ツール。OAuth 2.0 認証、プロジェクト・タスクの CRUD、タグ管理、スクリプト連携に対応。

## 特徴

- **タスク** — 優先度・期日・タグ付きでタスクの作成・取得・更新・完了・削除
- **プロジェクト** — プロジェクト一覧と詳細の取得
- **タグ** — 全プロジェクトからタグを集約して一覧表示
- **柔軟な期日指定** — `today`、`tomorrow`、`+3d`、`YYYY-MM-DD`
- **優先度** — `none`、`low`、`medium`、`high`
- **複数の出力形式** — テキスト、JSON、TSV
- **OAuth 2.0** — ブラウザベースのログイン、トークン自動更新

## インストール

### Homebrew

```bash
brew install tackeyy/tap/ticky
```

### Go

```bash
go install github.com/tackeyy/ticky@latest
```

### ソースからビルド

```bash
git clone https://github.com/tackeyy/ticky.git
cd ticky
go build -o ticky .
```

## クイックスタート

### 1. TickTick アプリの作成

1. [TickTick Developer Portal](https://developer.ticktick.com/manage) にアクセスし **+Create App** をクリック
2. アプリ名を入力（例: `ticky`）

### 2. OAuth 設定

アプリ設定画面で以下を設定:

| 設定項目 | 値 |
|---|---|
| Redirect URL | `http://localhost:18080/callback` |
| Scopes | `tasks:read`、`tasks:write` |

### 3. 環境変数の設定

```bash
export TICKTICK_CLIENT_ID=your_client_id
export TICKTICK_CLIENT_SECRET=your_client_secret
```

### 4. ログインして実行

```bash
ticky auth login
ticky tasks list --json
```

## コマンド

### `auth login` — OAuth でログイン

```bash
ticky auth login
```

ブラウザで TickTick の認証画面を開きます。トークンは `~/.config/ticky/token.json` に保存されます。

### `auth status` — 認証状態を確認

```bash
ticky auth status [--json] [--plain]
```

### `auth logout` — トークンを削除

```bash
ticky auth logout
```

### `tasks list` — タスク一覧を取得

```bash
ticky tasks list [--project <id>] [--json] [--plain]
```

| フラグ | 必須 | 説明 |
|---|---|---|
| `--project <id>` | No | プロジェクト ID（デフォルト: Inbox） |

### `tasks get` — タスク詳細を取得

```bash
ticky tasks get <task_id> --project <id> [--json] [--plain]
```

| フラグ | 必須 | 説明 |
|---|---|---|
| `<task_id>` | Yes | タスク ID |
| `--project <id>` | Yes | プロジェクト ID |

### `tasks create` — タスクを作成

```bash
ticky tasks create --title <title> [--project <id>] [--content <text>] [--priority <level>] [--due <date>] [--tags <tags>] [--json] [--plain]
```

| フラグ | 必須 | 説明 |
|---|---|---|
| `--title <title>` | Yes | タスクのタイトル |
| `--project <id>` | No | プロジェクト ID（デフォルト: Inbox） |
| `--content <text>` | No | タスクの内容・説明 |
| `--priority <level>` | No | `none`、`low`、`medium`、`high` |
| `--due <date>` | No | `today`、`tomorrow`、`+3d`、`YYYY-MM-DD` |
| `--tags <tags>` | No | カンマ区切りのタグ |

例:

```bash
# 優先度と期日を指定して作成
ticky tasks create --title "PR レビュー" --priority high --due tomorrow

# タグ付きで作成（JSON 出力）
ticky tasks create --title "牛乳を買う" --tags "買い物,個人" --json
```

### `tasks update` — タスクを更新

```bash
ticky tasks update <task_id> --project <id> [--title <title>] [--content <text>] [--priority <level>] [--due <date>] [--clear-due] [--tags <tags>] [--add-tags <tags>] [--remove-tags <tags>] [--json] [--plain]
```

| フラグ | 必須 | 説明 |
|---|---|---|
| `<task_id>` | Yes | タスク ID |
| `--project <id>` | Yes | プロジェクト ID |
| `--title <title>` | No | 新しいタイトル |
| `--content <text>` | No | 新しい内容 |
| `--priority <level>` | No | `none`、`low`、`medium`、`high` |
| `--due <date>` | No | 新しい期日 |
| `--clear-due` | No | 期日をクリア |
| `--tags <tags>` | No | タグを全置換 |
| `--add-tags <tags>` | No | タグを追加 |
| `--remove-tags <tags>` | No | タグを削除 |

例:

```bash
ticky tasks update abc123 --project def456 --priority high --due +3d
ticky tasks update abc123 --project def456 --add-tags "緊急" --json
```

### `tasks complete` — タスクを完了

```bash
ticky tasks complete <task_id> --project <id> [--json] [--plain]
```

| フラグ | 必須 | 説明 |
|---|---|---|
| `<task_id>` | Yes | タスク ID |
| `--project <id>` | Yes | プロジェクト ID |

### `tasks delete` — タスクを削除

```bash
ticky tasks delete <task_id> --project <id> [--json] [--plain]
```

| フラグ | 必須 | 説明 |
|---|---|---|
| `<task_id>` | Yes | タスク ID |
| `--project <id>` | Yes | プロジェクト ID |

### `projects list` — プロジェクト一覧を取得

```bash
ticky projects list [--json] [--plain]
```

### `projects get` — プロジェクト詳細を取得

```bash
ticky projects get <project_id> [--json] [--plain]
```

| フラグ | 必須 | 説明 |
|---|---|---|
| `<project_id>` | Yes | プロジェクト ID |

### `tags list` — タグ一覧を取得

```bash
ticky tags list [--json] [--plain]
```

全プロジェクトのタスクからタグを集約し、使用数の多い順に表示します。

## 設定

### 環境変数

| 変数 | 必須 | 説明 |
|---|---|---|
| `TICKTICK_CLIENT_ID` | Yes | OAuth クライアント ID |
| `TICKTICK_CLIENT_SECRET` | Yes | OAuth クライアントシークレット |
| `TICKTICK_ACCESS_TOKEN` | No | アクセストークン直接指定（トークンファイルを無視。CI やエージェント向け） |

### トークンの保存

`ticky auth login` 実行後、OAuth トークンは `~/.config/ticky/token.json` にパーミッション `0600` で保存されます。トークンの更新は自動で行われます。

`TICKTICK_ACCESS_TOKEN` が設定されている場合、トークンファイルは無視されます。

## 出力形式

### テキスト（デフォルト）

```
abc123def456789012345678 PR レビュー [high] (due: 2026-02-12) #仕事
```

### JSON（`--json`）

```json
[
  {
    "id": "abc123def456789012345678",
    "projectId": "inbox123",
    "title": "PR レビュー",
    "priority": 5,
    "dueDate": "2026-02-12T14:59:59.000+0000",
    "tags": ["仕事"]
  }
]
```

### TSV（`--plain`）

```
abc123def456789012345678	inbox123	PR レビュー	high	2026-02-12T14:59:59.000+0000	仕事
```

## 開発

```bash
go build -o ticky .
```

## ライセンス

MIT

## リンク

- [GitHub リポジトリ](https://github.com/tackeyy/ticky)
- [TickTick Open API ドキュメント](https://developer.ticktick.com/api)
