# rrk

## インストール

### クイックインストール（推奨）

```bash
curl -LsSf https://raw.githubusercontent.com/MRyutaro/rrk/main/install.sh | sh
```

### 手動ダウンロード

[Releases](https://github.com/MRyutaro/rrk/releases)ページから、お使いのOSに対応したバイナリをダウンロードしてください。

```bash
# macOS (Apple Silicon)の例
curl -L https://github.com/MRyutaro/rrk/releases/latest/download/rrk-darwin-arm64 -o rrk
chmod +x rrk
sudo mv rrk /usr/local/bin/
```

### ソースからビルド

```bash
git clone https://github.com/MRyutaro/rrk.git
cd rrk
make build
sudo mv rrk /usr/local/bin/
```
