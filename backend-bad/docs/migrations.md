# Golang Migrations

`golang-migrate` CLI を使って `backend/migrations` 配下の `.sql` ファイルを実行します。  
Docker compose で起動した Postgres（`backend/.env` の `DATABASE_URL`）に対してコマンドを流す想定です。

## 事前準備

1. `golang-migrate` CLI をローカルにインストール  
   <https://github.com/golang-migrate/migrate/tree/master/cmd/migrate>
2. `backend/.env` をコピーして `DATABASE_URL` をセット
   ```env
   DATABASE_URL=postgresql://user:password@localhost:5432/IMMORTAL_ARCHITECTURE_BAD_API?sslmode=disable
   ```

## Makefile タスク

`backend/Makefile` で以下のターゲットを用意しています。`DATABASE_URL` は環境変数で上書き可能です。

```bash
# 1. 新しいマイグレーションファイルを作成（空の up/down）
make migrate-create name=add_notes_table

# 2. 最新まで適用
make migrate-up

# 3. 直前のマイグレーションを1つ戻す
make migrate-down

# 4. 失敗状態をリセット（バージョン番号を指定）
make migrate-force version=20240101010101
```

`migrate-create` は `backend/migrations` に `YYYYMMDDHHMMSS_add_notes_table.up.sql` / `.down.sql` を生成します。  
`make migrate-up`/`down` 実行前に DB コンテナが起動していることを確認してください。
