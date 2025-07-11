# rrk

[![GitHub release](https://img.shields.io/github/release/MRyutaro/rrk.svg)](https://github.com/MRyutaro/rrk/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**rrk**（rireki）は、Go製の単一実行ファイルCLIツールで、シェル履歴をディレクトリごとにツリー形式で表示する履歴可視化ツールです。

> 📖 **English Documentation** - [README.md](./README.md)

## 特徴

- 🌳 **ツリー表示** - コマンド履歴をディレクトリツリー形式で表示
- 📁 **ディレクトリ別整理** - 各ディレクトリでどんなコマンドを実行したかを一目で確認
- 🎯 **集中表示** - 特定ディレクトリの履歴のみを表示
- 🚀 **単一実行ファイル** - 依存関係なし
- 💾 **軽量** - データベース不要のファイルベース保存
- 🐚 **シェル統合** - bash/zsh対応（自動セットアップ）
- 🔄 **自動アップデート** - GitHub Releasesとの統合アップデート機能
- 🗑️ **簡単削除** - データ保存オプション付きクリーンアンインストール

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

### ツリー表示

```bash
# 全ての履歴をツリー形式で表示
rrk

# 特定ディレクトリの履歴を表示
rrk /path/to/directory

# 各ディレクトリで表示するコマンド数を制限
rrk -n 5
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
# コマンド履歴をツリー形式で表示
$ rrk
/home/user
├── project/
│   ├── git status
│   ├── git add .
│   └── git commit -m "fix bug"
├── scripts/
│   ├── ./deploy.sh
│   └── python backup.py
└── .config/
    └── vim init.vim

/var
└── log/
    ├── tail -f syslog
    ├── grep ERROR *.log
    └── journalctl -u nginx

# 特定ディレクトリの履歴を表示
$ rrk /home/user/project
├── git status
├── git add .
└── git commit -m "fix bug"

# ディレクトリごとのコマンド数を制限
$ rrk -n 2
/home/user
├── project/
│   ├── git add .
│   └── git commit -m "fix bug"
└── scripts/
    └── python backup.py
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