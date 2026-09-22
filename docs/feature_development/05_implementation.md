# ⑤ 実装（AI活用）

> **この章でやること**
> ここまでの成果物を、そのままAIへの指示に変える。

---

## 設計が終わっていれば実装は速い

```
❌ 設計なしで投げる
   「Notion連携機能を作って」
   → AIが勝手に決める → 意図と違うものができる → 手直しで2時間

✅ 設計してから投げる
   「04_work_breakdown.md の PR4 を実装して」
   → 一発で意図どおり → レビュー20分
```

**AIは「何を作るか」は決められません。「どう作るか」は得意です。**

---

## Step 1. 文脈を読ませる

```
以下のドキュメントを読んで、このプロジェクトの設計方針を理解してください。

【アーキテクチャ】
- backend-clean/docs/02_clean_architecture_guide.md
- backend-clean/docs/08_cqrs_architecture_guide.md

【今回の設計】
- docs/feature_development/02_design_policy.md
- docs/feature_development/02_design_policy.md
- docs/feature_development/03_test_strategy.md
- docs/feature_development/04_work_breakdown.md
```

> **なぜ設計ドキュメントを読ませるのか**
> AIはコードだけ見ても「なぜこうなっているか」がわかりません。
> 「②で失敗したらロールバックすると決めた」を知っていれば、
> 勝手に「失敗しても公開は成功させる」実装をしません。

---

## Step 2. PR単位で投げる

**一度に全部投げないこと。** PR3の作業中にPR5の話をすると混乱します。

### PR0 既存テストの追加

```
04_work_breakdown.md の PR0 を実装してください。

【作るもの】
backend-clean/internal/usecase/note_command_interactor_test.go

【目的】
これから ChangeStatus と Update を変更するので、
先に現状の振る舞いを固定するテストを書きます。

【🚨 最重要の制約】
プロダクションコードを1行も変更しないでください。
Notion連携のコードも書かないでください。
このPRは「いまの動きをテストで記録する」だけです。

【テストすべき振る舞い】
- ChangeStatus: Draft⇄Publish、他人のノートで ErrUnauthorized、
                不正な遷移で ErrInvalidStatusChange、ReadModel が Upsert される
- Update: タイトル・セクションの更新、権限エラー、空タイトル
- Create / Delete も同様

【参考にすべき既存コード】
internal/usecase/template_interactor_test.go と同じ形式で、gomock を使う

【目標】
usecase のカバレッジを 37.2% から 70%以上へ
```

### PR1-a マイグレーション

```
04_work_breakdown.md の PR1-a を実装してください。

【作るもの】
backend-clean/migrations/ に up と down のSQLファイル

【内容】
- templates に notion_parent_page_id 列を追加
- notes に notion_page_id / notion_page_url / notion_synced_at 列を追加

【制約】
- 既存行を壊さないこと（すべて NULL 許容）
- 既存のマイグレーションファイルの命名規則に合わせること

【完了後】
make sqlc-generate を実行し、生成物も一緒にコミット対象にする
```

### PR1-b API定義

```
04_work_breakdown.md の PR1-b を実装してください。

【変更するもの】
- api-schema/typespec/models/template.tsp
- api-schema/typespec/models/note.tsp

【追加するフィールド】
TemplateResponse / CreateTemplateRequest / UpdateTemplateRequest
  + notionParentPageUrl

NoteResponse
  + notionPageUrl

【🚨 制約】
すべて optional にしてください。
必須にすると、実装が追いつくまで既存のクライアントが壊れます。

【生成の順序】
1. typespec を編集
2. api-schema で pnpm openapi   → openapi.yaml
3. backend-clean で make oapi   → Goの型
4. frontend で pnpm openapi     → TSの型

生成物はすべてコミットしてください。
再生成しても差分が出ないことを確認してください。

【このPRの目的】
先に契約を確定させることで、バックエンドとフロントを並行して書けるようにします。
実装はまだ入れません。
```

### PR2-a ドメイン層

```
04_work_breakdown.md の PR2-a を実装してください。

【作るもの】
backend-clean/internal/domain/notion/ に entity.go, logic.go, logic_test.go

【要件】
- Sync エンティティ（NotionPageID は *string で、nil なら未連携）
- IsFirstSync() … NotionPageID が nil なら true
- CanSync() … オーナー以外は ErrUnauthorized

【🚨 制約】
ドメイン層なので、net/http も database/sql も外部ライブラリも import しないこと

【参考にすべき既存コード】
internal/domain/note/logic.go の書き方に揃える
エラーは internal/domain/errors に追加する

【テスト】
internal/domain/note/logic_test.go と同じテーブル駆動テスト形式。カバレッジ90%以上
```

### PR2-b Notion APIクライアント

```
04_work_breakdown.md の PR2-b を実装してください。

【作るもの】
- internal/port/notion_port.go
- internal/adapter/gateway/externalapi/notion/client.go, converter.go, client_test.go

【置き場所の理由】
gateway/externalapi/ は設計ガイドで「外部API（将来用）」として
用意されている空ディレクトリです。今回そこに実装します。
gateway/db/ が DB を担うのと同じ位置づけです。

【Notion API の仕様】
01_requirement_hearing.md 末尾の「参考: Notion API の事実」を参照。
- POST https://api.notion.com/v1/pages
- ヘッダー: Authorization: Bearer / Notion-Version: 2022-06-28
- 親の指定: { "parent": { "page_id": "xxx" } }
- ページ配下では title のみ設定可。本文は children にブロック配列で渡す
- ゴミ箱: PATCH で in_trash: true / 復元は false

【🚨 必須要件】
1. http.Client に必ず Timeout を設定（10秒）
   デフォルトの http.Client はタイムアウト無限で危険です
2. リトライは 429 と 5xx のみ。4xx はリトライしないこと
3. リトライは指数バックオフ（1秒→2秒→4秒、最大3回）
4. APIキーをログやエラーメッセージに出さないこと
5. ベースURLを差し替え可能にすること（WithBaseURL）… テスト用

【テスト】
🚨 実際のNotion APIを呼ばないこと。httptest.NewServer でモックする。
網羅するケース: 正常系 / 401（リトライしない）/ 429（リトライする）/
                500 / タイムアウト / 壊れたJSON
```

### PR3 リポジトリ + 親ページURLの設定

```
04_work_breakdown.md の PR3 を実装してください。

【作るもの】
- note_repository.go … page_id の読み書きを追加
- template_repository.go … 親ページIDの読み書きを追加
- template_interactor.go … URLからIDを抽出する処理
- template_controller.go / presenter … API定義（PR1-b）に実装を合わせる

【URLからIDを抽出する】
https://notion.so/workspace/1429989fe8ac4effbc8f57f56486db54?v=...
                            └──── この32文字がID ────┘
ハイフンあり・なし両方を受け付けること。

【制約】
既存の note_command_interactor.go は触らないでください。
このPRでは、まだ公開処理に連携を組み込みません。

【確認】
PR0で追加したテストが全てパスすること
```

### PR4 ChangeStatus / Update への組み込み

```
04_work_breakdown.md の PR4 を実装してください。

【変更するファイル】
- internal/usecase/note_command_interactor.go  ← 既存を変更
- internal/domain/note/read_model.go … NotionPageURL を追加
- internal/usecase/note_command_interactor.go の toReadModel() … コピー処理
- internal/adapter/gateway/db/sqlc/note_read_model_repository.go
- internal/adapter/http/presenter/note_helpers.go … レスポンス変換
- internal/driver/config/config.go
- internal/driver/initializer/api/initializer.go

【🚨 CQRS の注意】
このアプリは画面が note_read_models しか読みません。
notes に値を入れても、toReadModel() でコピーしないと画面に出ません。
コピー漏れはコンパイルが通るため、テストで検出してください。

【🚨 制約1: 呼び出し順序】
Notionを先に呼び、成功したらDBを更新してください。

❌ DBを更新 → Notion呼び出し → 失敗したらDBを戻す
✅ Notion呼び出し → 成功したらDBをトランザクションで更新
   （失敗したらDBを触っていないので、戻す処理が不要）

【🚨 制約2: トランザクション】
トランザクションの中で Notion API を呼ばないでください。
ロックが長時間保持され、コネクションプールが枯渇します。

【🚨 制約3: 孤児ページの掃除】
Notion作成に成功したあとDB更新が失敗したら、
作ったページをゴミ箱に入れてください（page_id が保存されず、
二度と特定できなくなるため）。
掃除自体が失敗したら、ログに page_id を出力します。

【🚨 制約4: フォールバック】
NOTION_API_KEY が未設定のときは、連携をスキップして従来どおりの公開処理を
行ってください（移行期間中の切り替えスイッチを兼ねます）。
config.go で必須チェックにしないこと。

【ChangeStatus の振る舞い】
公開する場合:
  1. テンプレートから親ページIDを取得。未設定ならエラー（公開しない）
     ※画面側でもボタンを非活性にしますが、APIでも必ず検証します
  2. Notionにページ作成（既存 page_id があればゴミ箱から復元）
  3. 成功したらトランザクションで status 更新 + page_id 保存
  4. Notion失敗なら status を変更しない

非公開に戻す場合:
  1. Notionページをゴミ箱へ（in_trash: true）
  2. 成功したらトランザクションで status 更新
  3. Notion失敗なら status を変更しない

【Update の振る舞い】
  - 下書きのノート: 従来どおり（Notionは触らない）
  - 公開中のノート: 先にNotionを更新し、成功したらDBを更新

【🚨 確認必須】
PR0で追加した既存テストが全てパスすること。
1つでも落ちたら、既存の振る舞いを壊しています。
```

---

## テストも同じPRで書く

第3章で決めた方針どおり、**実装と同じPRでテストを書きます**。
「あとでまとめて」は、まず実現しません。

### UseCase層（gomock）

```go
{
    name: "初回連携では CreatePage が呼ばれる",
    setupMock: func(nc *mock.MockNotionClient, nr *mock.MockNoteRepository) {
        nr.EXPECT().Get(gomock.Any(), "note-1").
            Return(&note.WithMeta{/* notion_page_id は nil */}, nil)
        nc.EXPECT().CreatePage(gomock.Any(), "parent-1", "タイトル", gomock.Any()).
            Return(&port.PageResult{ID: "page-1"}, nil)
        nr.EXPECT().SaveNotionPage(gomock.Any(), "note-1", "page-1", gomock.Any())
    },
},
{
    name: "Notion失敗時は status を変更しない",
    setupMock: func(nc *mock.MockNotionClient, nr *mock.MockNoteRepository) {
        nc.EXPECT().CreatePage(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
            Return(nil, errors.New("notion is down"))
        // 🚨 UpdateStatus が呼ばれないことを検証（EXPECT を書かない）
    },
    wantErr: ErrNotionSyncFailed,
},
```

> **モックの利点: 「呼ばれなかったこと」を検証できる**
> `EXPECT()` を書かなければ、呼ばれた時点でテストが失敗します。
> 本物のAPIでは検証しづらい観点です。

**対応するQAケース:** QA-01, QA-02, QA-20

### Gateway層（httptest）

```go
tests := []struct {
    name       string
    statusCode int
    wantCalls  int   // リトライ回数
}{
    {name: "正常系",            statusCode: 200, wantCalls: 1},
    {name: "401はリトライしない", statusCode: 401, wantCalls: 1},
    {name: "429はリトライする",   statusCode: 429, wantCalls: 4},  // 初回+3回
    {name: "500はリトライする",   statusCode: 500, wantCalls: 4},
}
```

```go
srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    calls++
    w.WriteHeader(tt.statusCode)
}))
client := notion.NewClient("test-token", notion.WithBaseURL(srv.URL))
```

第3章で「ベースURLを差し替え可能にする」と決めたので、モックを差し込めています。

**対応するQAケース:** QA-31, QA-33

---

## AIが間違えやすいポイント

先回りして指示に含めてください。

| よくある間違い | なぜ起きるか | 先回りの指示 |
|---|---|---|
| トランザクション内で外部API | 一般的なコード例がそう | 「🚨 トランザクション外で」 |
| DB更新→Notion の順で書く | 直感的な順序がそう | 「🚨 Notionを先に呼ぶ」 |
| `http.Client{}` をそのまま使う | 入門コードがそう | 「Timeout を必ず設定」 |
| 4xx でもリトライする | 雑に実装する | 「429と5xxのみ」 |
| テストで本物のAPIを叩く | 指示がないと本物を使う | 「httptest を使う」 |
| 新しい設定を必須にする | バリデーションは善だと思う | 「必須にしない」 |
| 画面で止めるならAPIは省略 | 二重チェックを無駄だと思う | 「APIでも必ず検証」 |
| DB失敗時にNotionを放置 | 補償処理を思いつかない | 「失敗したらページをゴミ箱へ」 |

**この表自体がレビューのチェックリストになります。**

---

## 生成されたコードのレビュー

### 機械的に確認

```bash
go build ./...
go test ./internal/... -cover     # 既存のカバレッジを下回っていないか
make lint
make sqlc-generate && git diff --exit-code
```

### 人間が見る

```
□ 🚨 トランザクションの中で外部APIを呼んでいないか
   → grep -n "WithinTransaction" して中身を目で見る

□ 🚨 Notionを先に呼び、成功後にDBを更新しているか

□ タイムアウトが設定されているか
   → grep -n "http.Client"

□ 設定未設定でアプリが起動するか
   → NOTION_API_KEY を外して実際に起動してみる

□ 既存ファイルを勝手に触っていないか
   → git diff --stat で 02_design_policy.md の「変更しないもの」と照合

□ エラーハンドリングが握りつぶしになっていないか
   → grep -n "_ = " や "// ignore"
```

---

## AIに任せてよいこと・ダメなこと

```
✅ 任せてよい
   決まった設計に沿ったコードを書く
   既存パターンに揃えて書く
   テストケースを網羅的に書く
   ボイラープレート、SQL

❌ 任せてはいけない
   何を作るかを決める
   既存概念を拡張すべきかの判断
   失敗時にユーザーがどう振る舞うか
   影響範囲が許容できるかの判断
```

> **この境界線が、この教材の主題です**
> ①〜④が「任せてはいけないこと」、⑤が「任せてよいこと」。
> **①〜④をやらずに⑤だけAIに投げると事故ります。**

---

## チェックリスト

```
□ AIに設計ドキュメントを読ませたか
□ PR単位で投げているか
□ 「参考にすべき既存コード」を指定したか
□ 🚨 の制約（順序、トランザクション、タイムアウト、必須チェック）を明記したか
□ 第3章のQAケースに対応するテストを書いたか
□ テストを実装と同じPRに含めたか
□ git diff --stat で変更範囲を確認したか
```

---

## 次の章へ

コードとテストができました。
次は本物のNotionと繋いで確認します。

👉 [06_completion_check.md](./06_completion_check.md)
