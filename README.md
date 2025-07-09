# rrk

**rrk**（rireki）は、Go製の単一バイナリCLIツールで、bash/zshのシェル履歴を**セッション**・**ディレクトリ**単位で論理的にグループ化し、過去のコマンドを「**そのまま再実行可能**」な形で抽出・管理できる履歴管理ツールです。

## 特徴

- 📁 **ディレクトリ単位**での履歴管理
- 🪟 **セッション単位**での履歴管理  
- 🔄 **ワンコマンド再実行** - `rrk rerun <ID>`
- 🚀 **単一バイナリ** - 依存関係なし
- 💾 **軽量** - データベース不要のファイルベース保存
- 🐚 **シェル統合** - bash/zsh対応

## インストール

### クイックインストール（推奨）

```bash
curl -LsSf https://raw.githubusercontent.com/MRyutaro/rrk/main/install.sh | sh
```

### ソースからビルド

```bash
git clone https://github.com/MRyutaro/rrk.git
cd rrk
make build
sudo mv rrk /usr/local/bin/
```

## セットアップ

インストール後、シェル統合を有効にします：

```bash
# 自動セットアップ
rrk setup

# 手動セットアップ（bash）
echo 'eval "$(rrk hook init bash)"' >> ~/.bashrc
source ~/.bashrc

# 手動セットアップ（zsh）
echo 'eval "$(rrk hook init zsh)"' >> ~/.zshrc
source ~/.zshrc
```

## 使い方

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
rrk cwd show
rrk d show

# 特定のディレクトリの履歴を表示
rrk cwd show /path/to/directory

# 履歴があるディレクトリ一覧
rrk cwd list
```

### コマンド再実行

```bash
# 履歴IDを指定して再実行
rrk rerun <HISTORY_ID>

# 例：ID=1のコマンドを再実行
rrk rerun 1
```

## 使用例

```bash
# 現在のディレクトリの履歴を確認
$ rrk cwd show
ID  TIME      SESSION        COMMAND
1   14:30:12  abc123...      git status
2   14:30:45  abc123...      git add .
3   14:31:02  abc123...      git commit -m "fix bug"

# 特定のコマンドを再実行
$ rrk rerun 2
Re-running: git add .
Original directory: /Users/user/project
Current directory: /Users/user/project

# セッション履歴を確認
$ rrk session show
ID  TIME      DIRECTORY       COMMAND
1   14:30:12  ~/project       git status
2   14:30:45  ~/project       git add .
3   14:31:02  ~/project       git commit -m "fix bug"
4   14:32:15  ~/documents     vim README.md
```

## データ保存

- 履歴データは `~/rrk/history.jsonl` に保存
- セッション情報は `~/rrk/current_session` に保存
- 外部データベース不要

## 開発者向け

詳細は [`docs/DEVELOPERS.md`](./docs/DEVELOPERS.md) を参照してください。
