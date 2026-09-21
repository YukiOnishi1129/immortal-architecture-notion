# 🛠️ ローカル環境セットアップガイド

> 💡 **このドキュメントのゴール**
> backend-cleanをローカル環境で動かし、開発を始められる状態にする

---

## 🎯 前提条件

以下がインストールされていることを確認してください:

```
✅ Docker & Docker Compose
✅ Go 1.25以上
✅ Node.js 22以上 & pnpm 10以上
✅ make コマンド
✅ migrate コマンド（golang-migrate）
```

---

## 📦 ステップ1: 環境変数の設定

```bash
# プロジェクトroot直下に.envファイルがあることを確認
# .env.exampleをコピーして.envを作成

cp .env.example .env
```

`.env`の内容を確認・編集:

```bash
# データベース設定
DB_HOST=localhost
DB_PORT=5432
DB_USER=mini_notion_user
DB_PASSWORD=mini_notion_password
DB_NAME=mini_notion_db

# API設定
API_PORT=8080
```

---

## 🐳 ステップ2: Docker Composeでコンテナを起動

```bash
# プロジェクトroot直下でDocker Composeを起動
# ※ backend-cleanディレクトリではなく、root直下で実行してください

docker compose up -d
```

これで以下のコンテナが起動します:
- `db`: PostgreSQL 15（ポート: 5432）
- `api`: backend-clean API（ポート: 8080）

**確認:**

```bash
# コンテナが起動しているか確認
docker compose ps

# ログを確認
docker compose logs -f api
```

---

## 📝 ステップ3: OpenAPI定義からGoコードを生成

### 3-1. TypeSpecからOpenAPI YAMLを生成

```bash
# api-schemaディレクトリでTypeSpecからOpenAPI YAMLを生成
cd api-schema
pnpm install
pnpm run generate:openapi
```

**生成されるファイル:**
- `api-schema/generated/openapi.yaml` - OpenAPI仕様書

### 3-2. OpenAPIからGoコードを生成

```bash
# backend-cleanディレクトリでOpenAPIからGoコードを生成
cd ../backend-clean
make oapi
```

**生成されるファイル:**
- `backend-clean/internal/driver/oas/server.gen.go` - OpenAPIサーバーコード
- `backend-clean/internal/driver/oas/types.gen.go` - リクエスト/レスポンス型

---

## 🗄️ ステップ4: データベースマイグレーション

```bash
# backend-cleanディレクトリでマイグレーション実行
cd backend-clean
make migrate-up
```

これで以下のテーブルが作成されます:
- `accounts` - アカウント
- `templates` - テンプレート
- `fields` - フィールド
- `notes` - ノート
- `sections` - セクション

**確認（方法1: psqlコマンド）:**

```bash
# PostgreSQLに接続してテーブルを確認
docker compose exec db psql -U mini_notion_user -d mini_notion_db

# テーブル一覧を表示
\dt

# 終了
\q
```

**確認（方法2: TablePlus）:**

[TablePlus](https://tableplus.com/) を使うと、GUIで簡単にテーブルを確認できます。

```
接続情報:
- Host: localhost
- Port: 5432
- User: mini_notion_user
- Password: mini_notion_password
- Database: mini_notion_db
```

TablePlusで接続後、以下を確認:
- テーブル一覧（accounts, templates, fields, notes, sections）
- テーブルのカラム定義
- 実際のデータ（INSERT後）

---

## ✅ ステップ5: 動作確認

コンテナのログでAPIが起動しているか確認:

```bash
# APIコンテナのログを確認
docker compose logs api

# 起動成功のログが表示されることを確認
# 例: "Server started on :8080"
```

**APIの疎通確認:**

APIエンドポイントは認証が必要なため、ログで起動確認ができればOKです。
実際のAPI呼び出しは認証実装後に確認できます。

---

## 🧪 ステップ6: テストの実行

```bash
cd backend-clean

# 全テストを実行
make test-all

# カバレッジを確認
go test -cover ./...

# 特定パッケージのテストを実行
make test-pkg PKG=./internal/domain/note

# 特定のテストを実行
make test-run PKG=./internal/domain/note RUN=TestValidateNoteForCreate
```

**テストが通ればセットアップ完了！**

---

## 🔧 よく使うコマンド

### Docker関連

```bash
# コンテナを起動
docker compose up -d

# コンテナを停止
docker compose down

# ログを確認
docker compose logs -f api

# DBに接続
docker compose exec db psql -U mini_notion_user -d mini_notion_db
```

### マイグレーション関連

```bash
# マイグレーションを実行
make migrate-up

# マイグレーションを1つ戻す
make migrate-down

# 新しいマイグレーションを作成
make migrate-new NAME=add_tags_table
```

### OpenAPI生成関連

```bash
# TypeSpecからOpenAPI YAMLを生成
cd api-schema
pnpm run generate:openapi

# OpenAPIからGoコードを生成
cd ../backend-clean
make oapi
```

### テスト関連

```bash
# 全テストを実行
make test-all

# カバレッジを確認
go test -cover ./...

# 特定パッケージのテストを実行
make test-pkg PKG=./internal/domain/note

# 特定のテストを実行
make test-run PKG=./internal/domain/note RUN=TestValidateNoteForCreate
```

### ビルド関連

```bash
# ビルド
make build

# Lint
make lint
```

---

## 🚨 トラブルシューティング

### 問題1: ポートが既に使われている

```bash
# エラー: Bind for 0.0.0.0:8080 failed: port is already allocated

# 解決策: .envのAPI_PORTを変更
API_PORT=8081
```

### 問題2: マイグレーションが失敗する

```bash
# エラー: error: dial tcp: lookup db: no such host

# 解決策: DBコンテナが起動しているか確認
docker compose ps

# DBコンテナを再起動
docker compose restart db
```

### 問題3: OpenAPI生成が失敗する

```bash
# エラー: OpenAPI spec not found at ../api-schema/generated/openapi.yaml

# 解決策: TypeSpecからOpenAPI YAMLを生成
cd api-schema
pnpm run generate:openapi
```

### 問題4: テストが失敗する

```bash
# エラー: mock_note_repository.go: no such file or directory

# 解決策: モックを生成
go generate ./...
```

---

## 🎯 次のステップ

1. **既存コードを読む**
   - `internal/domain/note/` から読み始める
   - `internal/usecase/note_interactor.go` を読む

2. **テストを書いてみる**
   - Domainのバリデーションから始める
   - テーブル駆動テストで書く

3. **新しい機能を追加してみる**
   - 「コメント機能」「タグ機能」など
   - AIと一緒に実装（[04_ai_driven_development.md](./04_ai_driven_development.md)を参照）

**Happy Coding!** 🎉
