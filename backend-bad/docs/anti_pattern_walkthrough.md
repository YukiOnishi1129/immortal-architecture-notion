# Lesson 2: アンチパターン実装ウォークスルー

ここからは実際のコードを開きながら「どこで何が壊れているのか」を確かめます。各セクションに **Try it** を挟んでいるので、IDE を開きつつ同じファイルを覗いてください。ファイル名は `backend/` を基準に書いてあります。

---

## 0. ざっくり診断表

| どこで壊れている? | 何が起きている? | まず見るファイル |
| --- | --- | --- |
| Controller が肥大化 | 1 メソッドで HTTP 入出力〜DB 詰め替えまで全部担当 | `internal/controller/*_controller.go` |
| Service が外部依存まみれ | OpenAPI 型・sqlc 型・`echo.Context` をそのまま握る | `internal/service/*.go` |
| トランザクションを素手で書く | 全ユースケースで `Begin/Commit/Rollback` をコピペ | `internal/service/note_service.go` ほか |
| if/switch で泥沼 | ステータスや検索条件の分岐が延々と続く | `internal/service/note_service.go` |
| エラーの顔が全部一緒 | `failed to ...` しか返らず原因不明 | `internal/controller/helpers.go` |
| テストがほぼ無理 | Service テストがコメントで「書きたいけど無理」と叫ぶ | `internal/service/template_service_test.go` |

---

## 1. Controller が何でも屋

- 例: `internal/controller/template_controller.go` の `TemplatesCreateTemplate`
- JSON バインド → 追加バリデーション → DB 用 struct に詰め替え → Service 呼び出し → レスポンス整形を 1 関数で実施。
- Echo の `Context`、OpenAPI の DTO、sqlc のモデルが同居しており、責務が 1 ミリも分離されていない。
- **痛み**: 仕様追加ごとに Controller を修正し、同じバリデーションを複数メソッドへコピペ。テストを書くには Echo Context と DB 型を両方構築する必要がある。

> ✅ **Try it**: `rg "TemplatesCreateTemplate" -n internal/controller` を実行し、1 関数の行数と `sqldb.Field` を組み立てる処理を確認してみてください。

---

## 2. Service が Echo / DB にベッタリ

- 例: `internal/service/template_service.go`
- メソッドシグネチャが `CreateTemplate(ctx echo.Context, ...)` という時点で Web フレームワーク依存。
- OpenAPI で生成した `openapi.Models*` と sqlc が吐く `sqldb.*` をそのまま受け渡し、ドメイン独自の型が存在しない。
- **痛み**: API/DB の仕様変更が全層へ波及。ユニットテストでは Echo Context も sqlc モックも必要で、テストコードだけで疲弊します。

> ✅ **Try it**: `rg "CreateTemplate(ctx echo.Context" -n` で該当メソッドを開き、引数と戻り値に何種類の型が混ざっているか数えてみてください。

---

## 3. トランザクションを素手で管理

- 例: `internal/service/note_service.go` の `CreateNote`, `UpdateNote`, `ChangeStatus`
- どのユースケースも `tx, err := s.pool.Begin(ctx)` から始まり、`Rollback`/`Commit` を手書き。わずかな変更でロールバック忘れが発生します。
- `make lint` を実行すると `Rollback` のエラーチェック漏れ警告が大量に出るのは、コピペ構造のせい。
- **痛み**: Repository 層が存在しないため、ユニットテストでトランザクションを挙動させる術がなく、DB を立ち上げた統合テスト頼みになります。

> ✅ **Try it**: `rg "Begin(ctx.Request" -n internal/service` を実行し、何ファイルで同じパターンが繰り返されているか数えてみる。

---

## 4. if 文だらけのユースケース

- 例: `internal/service/note_service.go` の `ListNotes` や `ChangeStatus`
- ステータスや検索条件を素の文字列で比較し、`if filters.OwnerID != nil { ... }` が延々続く。
- 新しい状態や絞り込み条件を追加すると、関係しそうな if 文を全部探して修正する羽目に。
- **痛み**: 条件の抜け漏れが発生しやすく、レビュー側も全体像を把握できません。

> ✅ **Try it**: `note_service.go` で `status != "Draft" && status != "Publish"` を検索し、別の関数でも同じようなチェックが散らばっていることを確認しましょう。

---

## 5. 役に立たないエラーレスポンス

- 例: `internal/controller/helpers.go` の `respondError`
- どのエラーでも `Code = HTTP ステータス文字列`, `Message = 固定文言` で返却。`template not found` と `internal error` の区別がつかないケースがほとんど。
- **痛み**: クライアントはログを覗かないと原因が分からず、API の利用者からは「常に failed to ...」としか見えない。

> ✅ **Try it**: `respondError` の呼び出し元を `rg "respondError" -n internal/controller` で一覧し、すべて固定メッセージになっていることを確認してください。

---

## 6. テストで詰む様子をそのまま展示

- 例: `internal/service/template_service_test.go`
- ユニットテストには `ownerID` の形式チェックなど一部だけを書き、書きたいテストはコメントで列挙。「トランザクションがロールバックされるか」など DB 依存のケースは再現不能だからです。
- **痛み**: 「本当はここまで保証したい」という想いがコメントに滲むだけで、品質は上がらない。DB を使った統合テスト無しでは何も担保できない構造だと再認識できます。

> ✅ **Try it**: `template_service_test.go` のコメントアウトされたケースを読み、「なぜモックで再現できないのか？」を自分の言葉で説明できるか挑戦してみましょう。

---

## 7. 今日の気づきをメモしよう

| アンチパターン | 見つけたファイル/行 | どんな痛みを感じたか | どう直したいかのメモ |
| --- | --- | --- | --- |
| Controller 何でも屋 | | | |
| Service が Echo/DB 依存 | | | |
| 手書きトランザクション | | | |
| if/switch 地獄 | | | |
| 使えないエラーレスポンス | | | |
| テスト不能 | | | |

※ 上の表は学習ノート用。コピーして使ってください。`anti_pattern_fix.md` を読むときに、このメモが役立ちます。

---

ここまでで「どこが辛いか」が見えてきたはずです。次は `anti_pattern_fix.md` を開き、同じ問題をどうやって改善するのかをイメージしていきましょう。EOF
