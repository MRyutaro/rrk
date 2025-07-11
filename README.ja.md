# rrk

**rrk**（rireki）は、Go製の単一実行ファイルCLIツールで、bash/zshのシェル履歴を**セッション**・**ディレクトリ**単位で論理的にグループ化し、過去のコマンドを「**そのまま再実行可能**」な形で抽出・管理できる履歴管理ツールです。

## 特徴

- 📁 **ディレクトリ単位**での履歴管理
- 🪟 **セッション単位**での履歴管理  
- 🔄 **ワンコマンド再実行** - `rrk rerun <ID>`
- 🚀 **単一実行ファイル** - 依存関係なし
- 💾 **軽量** - データベース不要のファイルベース保存
- 🐚 **シェル統合** - bash/zsh対応

## インストール

### クイックインストール（推奨）

```bash
curl -LsSf https://raw.githubusercontent.com/MRyutaro/rrk/main/install.sh | sh
```

このスクリプトは以下を自動的に実行します：
- システムに適したバイナリをダウンロード
- `~/.local/bin`（または`$INSTALL_DIR`）にインストール
- シェル統合（bash/zsh）を自動設定
- 必要に応じてインストールディレクトリをPATHに追加

### ソースからビルド

```bash
git clone https://github.com/MRyutaro/rrk.git
cd rrk
make build
sudo mv rrk /usr/local/bin/

# ソースからビルドした後は、シェル統合を設定：
rrk setup
```

## 使い方

### 全履歴表示

```bash
# 全ての履歴を表示
rrk list

# 最新20件のみ表示
rrk list -n 20
```

### セッション管理

```bash
# セッション一覧
rrk session list

# 現在のセッション履歴を表示
rrk session show
rrk s show

# 特定のセッション履歴を表示
rrk session show <SESSION_ID>
```

### ディレクトリ管理

```bash
# 現在のディレクトリの履歴を表示
rrk dir show
rrk d show

# 特定のディレクトリの履歴を表示
rrk dir show /path/to/directory

# 履歴があるディレクトリ一覧（番号付きID表示）
rrk dir list

# IDでディレクトリ履歴を表示
rrk dir show <ID>
rrk d show <ID>
```

### コマンド再実行

```bash
# 履歴IDを指定して再実行
rrk rerun <HISTORY_ID>

# 例：ID=1のコマンドを再実行
rrk rerun 1
```

### アップデート

```bash
# 最新版にアップデート
rrk update
```

### バージョン確認

```bash
# バージョン情報を表示（GitHubの最新リリース情報も表示）
rrk -v
rrk --version
```

> **注意**: バージョン比較は開発ビルドを正しく処理し、`rrk update`実行後は更新通知が自動的にクリアされます。

### アンインストール

```bash
# シェル統合のみ削除
rrk uninstall
```

## 使用例

```bash
# 現在のディレクトリの履歴を確認
$ rrk dir show
ID  TIME      SESSION        COMMAND
1   14:30:12  abc123...      git status
2   14:30:45  abc123...      git add .
3   14:31:02  abc123...      git commit -m "fix bug"

# 特定のコマンドを再実行
$ rrk rerun 2
git add .

# セッション履歴を確認
$ rrk session show
ID  TIME      DIRECTORY       COMMAND
1   14:30:12  ~/project       git status
2   14:30:45  ~/project       git add .
3   14:31:02  ~/project       git commit -m "fix bug"
4   14:32:15  ~/documents     vim README.md

# 全履歴を確認
$ rrk list
ID  TIME      DIRECTORY       SESSION        COMMAND
1   14:30:12  ~/project       abc123...      git status
2   14:30:45  ~/project       abc123...      git add .
3   14:31:02  ~/project       abc123...      git commit -m "fix bug"
4   14:32:15  ~/documents     def456...      vim README.md

# ディレクトリ一覧をIDで表示
$ rrk dir list
ID  DIRECTORY        STATUS
0   ~/project        (current)
1   ~/documents
2   /tmp

# ディレクトリIDで履歴を表示
$ rrk dir show 1
ID  TIME      SESSION        COMMAND
4   14:32:15  def456...      vim README.md
```

## データ保存

- 履歴データは `~/.rrk/history.jsonl` に保存
- セッション情報は `~/.rrk/current_session` に保存
- 外部データベース不要

## 開発者向け

詳細は [`docs/DEVELOPERS.md`](./docs/DEVELOPERS.md) を参照してください。

### コントリビューション

- `main`ブランチへのプルリクエストのマージは自動的にパッチバージョンリリースをトリガーします
- CI/CDパイプラインがバージョン管理とGitHubリリースを自動で処理します

## ライセンス

MIT License