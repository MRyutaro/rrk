# rrk for developers

## セットアップ

```bash
git clone https://github.com/MRyutaro/rrk.git
cd rrk
go mod download
```

## ビルド

```bash
make build    # バイナリをビルド
./rrk         # ローカルで実行
```

## テスト

```bash
make test     # 全テスト実行
```

## リリース

セマンティックバージョニングに従ってリリース:

```bash
make patch    # バグ修正 (0.0.1 → 0.0.2)
make minor    # 新機能追加 (0.0.2 → 0.1.0)
make major    # 破壊的変更 (0.1.0 → 1.0.0)
```

タグをプッシュすると、GitHub Actionsが自動でマルチプラットフォームバイナリをビルド・配布します。
