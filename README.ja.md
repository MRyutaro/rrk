# rrk

[![GitHub release](https://img.shields.io/github/release/MRyutaro/rrk.svg)](https://github.com/MRyutaro/rrk/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**rrk**（rireki）は、Go製の単一実行ファイルCLIツールで、bash/zshのシェル履歴を**セッション**・**ディレクトリ**単位で論理的にグループ化し、過去のコマンドを「**そのまま再実行可能**」な形で抽出・管理できる履歴管理ツールです。

> 📖 **English Documentation** - [README.md](./README.md)

## 特徴

- 📁 **ディレクトリ単位**での履歴管理
- 🪟 **セッション単位**での履歴管理  
- 🔄 **ワンコマンド再実行** - `rrk rerun <ID>`
- 🚀 **単一実行ファイル** - 依存関係なし
- 💾 **軽量** - データベース不要のファイルベース保存
- 🐚 **シェル統合** - bash/zsh対応（自動セットアップ）
- 🔄 **自動アップデート** - GitHub Releasesとの統合アップデート機能
- 🗑️ **簡単削除** - データ保存オプション付きクリーンアンインストール

## インストール

### リリースからダウンロード（推奨）

1. [GitHub Releases](https://github.com/MRyutaro/rrk/releases)からシステムに適したバイナリをダウンロード
2. 実行可能にしてPATHに配置：

```bash
# Linux/macOSの例
chmod +x rrk-<OS>-<ARCH>
sudo mv rrk-<OS>-<ARCH> /usr/local/bin/rrk

# シェル統合をセットアップ
rrk setup
```

### ソースからビルド

```bash
git clone https://github.com/MRyutaro/rrk.git
cd rrk
make build
sudo mv rrk /usr/local/bin/

# シェル統合をセットアップ
rrk setup
```

### シェル統合セットアップ

インストール後、自動履歴記録を有効にするためにセットアップコマンドを実行：

```bash
rrk setup
```

これにより以下が実行されます：
- シェル（bash/zsh）の自動検出
- シェル設定（`~/.bashrc`または`~/.zshrc`）への統合フック追加
- `~/.rrk/`内の必要な設定ファイル作成

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
rrk s list

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
rrk d list

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

updateコマンドは以下を実行します：
- GitHub Releasesから最新バージョンをダウンロード
- 現在のバイナリを置き換え
- インストールを検証
- アップデート通知キャッシュをクリア

### バージョン確認

```bash
# バージョン情報を表示（アップデート通知も含む）
rrk version
rrk -v
rrk --version
```

### アンインストール

```bash
# シェル統合と全データを削除
rrk uninstall

# 確認なしでシェル統合と全データを削除
rrk uninstall -y
```

uninstallコマンドは以下を実行します：
- `~/.bashrc`/`~/.zshrc`からシェル統合を削除
- `~/.rrk/`から全rrkデータを削除
- バイナリ削除の手順を表示

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

- 履歴データは `~/.rrk/history.jsonl`（JSONL形式）に保存
- セッション情報は `~/.rrk/current_session` に保存
- シェル統合スクリプトは `~/.rrk/hook.sh` に保存
- バージョンキャッシュは `~/.rrk/.rrk_version_cache` に保存
- 外部データベース不要

## 高度な使用方法

### 手動シェル統合

手動セットアップやカスタム設定が必要な場合：

```bash
# シェル統合スクリプトを生成
rrk hook init bash > ~/.rrk_integration.sh
rrk hook init zsh > ~/.rrk_integration.sh

# シェル設定でソース
echo "source ~/.rrk_integration.sh" >> ~/.bashrc  # または ~/.zshrc
```

### 手動履歴記録

```bash
# コマンドを手動で記録
rrk hook record "your command here"

# 新しいセッションを初期化
rrk hook session-init
```

## CI/CD統合

rrkには自動リリース管理が含まれています：

- **プルリクエストマージ**: 自動的にパッチリリースを作成
- **手動タグ付け**: 全プラットフォームでリリースビルドをトリガー
- **マルチプラットフォームビルド**: Linux、macOS、Windows（AMD64/ARM64）
- **自動アップデート**: 内蔵のアップデート通知とインストール

## 開発者向け

詳細は [`docs/DEVELOPERS.md`](./docs/DEVELOPERS.md) と [`docs/REQ.md`](./docs/REQ.md) を参照してください。

### コントリビューション

- `main`ブランチへのプルリクエストのマージは自動的にパッチバージョンリリースをトリガーします
- CI/CDパイプラインがバージョン管理とGitHubリリースを自動で処理します
- ローカルでのバージョン管理には `make patch`、`make minor`、`make major` を使用

## ライセンス

MIT License