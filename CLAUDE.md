# CLAUDE.md

このファイルは、このリポジトリでコードを扱う際のClaude Code（claude.ai/code）へのガイダンスを提供します。

## プロジェクト概要

これは「rrk」（履歴）という名前のGoプロジェクトです。プロジェクトは現在初期セットアップ段階にあります。

## よく使う開発コマンド

[`docs/DEVELOPERS.md`](./docs/DEVELOPERS.md) を参照してください。

## プロジェクト構造

```
rrk/
├── main.go              # エントリーポイント
├── go.mod              # Goモジュール定義
├── Makefile            # ビルド・リリース用コマンド
├── scripts/            # ユーティリティスクリプト
│   └── bump-version.sh # バージョン管理スクリプト
├── docs/               # ドキュメント
│   ├── DEVELOPERS.md   # 開発者向けドキュメント
│   └── REQ.md         # 要件定義
├── .github/            # GitHub Actions設定
│   └── workflows/
│       ├── build.yml   # CI/CDビルド
│       └── release.yml # リリース自動化
└── CLAUDE.md          # このファイル
```
