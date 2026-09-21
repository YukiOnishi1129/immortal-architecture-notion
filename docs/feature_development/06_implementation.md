# ⑥ 実装（AI活用）

> **この章のゴール**
> ここまでの成果物を、そのままAIへの指示に変える。
> **設計が終わっていれば、実装はAIの仕事**だと体感する。

---

## なぜここまで設計してきたのか

```
❌ 設計なしでAIに投げる:
   「Notion連携機能を作って」
      ↓
   AI「どこに作りますか？」
   AI「既存のstatusと連動させますか？」
   AI「失敗したらどうしますか？」
      ↓
   AIが勝手に決める → 意図と違うものができる
      ↓
   手直しで結局2時間

✅ 設計してからAIに投げる:
   「02_design_policy.md の論点Cに従って、
     note_notion_syncs テーブルのリポジトリを作って」
      ↓
   AI「はい、ポートはこれですね」
      ↓
   一発で意図どおり → レビュー20分
```

**AIは「何を作るか」は決められません。「どう作るか」は得意です。**

---

## ステップ1: AIに文脈を読ませる

最初に、このプロジェクトのルールを読ませます。

```
以下のドキュメントを読んで、このプロジェクトの設計方針を理解してください。

【アーキテクチャ】
- backend-clean/docs/02_clean_architecture_guide.md
- backend-clean/docs/08_cqrs_architecture_guide.md

【今回の機能の設計】
- docs/feature_development/02_design_policy.md
- docs/feature_development/03_impact_analysis.md
- docs/feature_development/04_test_strategy.md
- docs/feature_development/05_work_breakdown.md

【全体設計】
- docs/global_design/05_domain_design.md
- docs/global_design/06_database_design.md
```

> **なぜ設計ドキュメントを読ませるのか**
> AIはコードだけ見ても「なぜこうなっているか」がわかりません。
> 「論点Aで公開と分離すると決めた」を知っていれば、
> **勝手に ChangeStatus を触ろうとしません。**

---

## ステップ2: PR単位で投げる

**一度に全部投げないこと。** PR④の作業中にPR⑥の話をすると、AIが混乱します。

### PR⓪ 既存テストの追加

```
docs/feature_development/05_work_breakdown.md の PR⓪ を実装してください。

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
- ChangeStatus: Draft → Publish、Publish → Draft
- ChangeStatus: 他人のノートで ErrUnauthorized
- ChangeStatus: 不正な遷移で ErrInvalidStatusChange
- ChangeStatus: ReadModel が Upsert される
- Update: タイトル・セクションの更新
- Update: 他人のノートで ErrUnauthorized
- Update: 空タイトルで ErrTitleRequired
- Create / Delete も同様

【参考にすべき既存コード】
internal/usecase/template_interactor_test.go と同じ形式で、
gomock を使って書いてください。

【目標】
usecase パッケージのカバレッジを 37.2% から 70%以上へ
```

### PR① マイグレーション

```
docs/feature_development/05_work_breakdown.md の PR① を実装してください。

【作るもの】
backend-clean/migrations/ に up と down のSQLファイル

【要件】
- notion_destinations テーブル（送り先マスタ）
- note_notion_syncs テーブル（連携状態）
- templates に notion_destination_id 列を追加

【制約】
- 既存の templates 行を壊さないこと（NULL許容にする）
- note_notion_syncs は notes 削除時に CASCADE で消えること
- status は pending/synced/failed のみ許可（CHECK制約）
- 既存のマイグレーションファイルの命名規則に合わせること

【完了後にやること】
make sqlc-generate を実行し、生成物も一緒にコミット対象にする
```

### PR② ドメイン層

```
docs/feature_development/05_work_breakdown.md の PR② を実装してください。

【作るもの】
backend-clean/internal/domain/notion/ 配下に
entity.go, logic.go, logic_test.go

【要件】
- SyncStatus 型（pending/synced/failed）
- Sync エンティティ
- IsFirstSync() … NotionPageID が nil なら true
- CanSync() … オーナー以外は ErrUnauthorized

【制約】
🚨 ドメイン層なので、以下を import しないこと:
   - net/http
   - database/sql
   - 外部ライブラリ全般

【参考にすべき既存コード】
internal/domain/note/logic.go の書き方に揃えてください。
エラーは internal/domain/errors に追加してください。

【テスト】
internal/domain/note/logic_test.go と同じ
テーブル駆動テストの形式で書いてください。
カバレッジ90%以上を目指してください。
```

> **「参考にすべき既存コード」を必ず指定する**
> これがないと、AIは一般的なGoの書き方をします。
> プロジェクトの流儀に揃えさせるには、**見本を示す**のが確実です。

### PR③ Notion APIクライアント

```
docs/feature_development/05_work_breakdown.md の PR③ を実装してください。

【作るもの】
- internal/port/notion_port.go
- internal/adapter/gateway/notion/client.go
- internal/adapter/gateway/notion/converter.go
- internal/adapter/gateway/notion/client_test.go

【Notion API の仕様】
01_requirement_hearing.md 末尾の「参考: Notion API の事実」を参照。
呼び出しに必要な情報:
- POST https://api.notion.com/v1/pages
- ヘッダー: Authorization: Bearer / Notion-Version: 2022-06-28
- 親の指定: { "parent": { "page_id": "xxx" } }
- ページ配下では title のみ設定可。本文は children にブロック配列で渡す

【🚨 必須要件】
1. http.Client に必ず Timeout を設定すること（10秒）
   デフォルトの http.Client はタイムアウト無限で危険です
2. リトライは 429 と 5xx のみ。4xx はリトライしないこと
3. リトライは指数バックオフ（1秒→2秒→4秒、最大3回）
4. APIキーをログやエラーメッセージに出さないこと

【テスト】
🚨 実際のNotion APIを呼ばないこと。
httptest.NewServer でモックサーバーを立ててテストしてください。

以下のケースを網羅してください:
- 正常系（201が返る）
- 401（トークン不正）→ リトライしない
- 429（レート制限）→ リトライする
- 500 → リトライする
- タイムアウト
- レスポンスのJSONが壊れている
```

### PR④ リポジトリ + 送り先の管理

```
docs/feature_development/05_work_breakdown.md の PR④ を実装してください。

【作るもの】
- internal/adapter/gateway/db/sqlc/notion_sync_repository.go
- internal/adapter/gateway/db/sqlc/notion_destination_repository.go
- internal/usecase/notion_destination_interactor.go
- internal/adapter/http/controller/notion_destination_controller.go

【機能】
- 送り先（notion_destinations）の登録・一覧・編集
- テンプレートへの送り先の割り当て
- テンプレート作成時に送り先を必須にする

【制約】
既存の note_command_interactor.go は触らないでください。
このPRでは、まだ公開処理に連携を組み込みません。

【参考にすべき既存コード】
internal/usecase/template_interactor.go
internal/adapter/gateway/db/sqlc/template_repository.go

【確認】
PR⓪で追加したテストが全てパスすること
```

### PR⑤ ChangeStatus / Update への組み込み

```
docs/feature_development/05_work_breakdown.md の PR⑤ を実装してください。

【変更するファイル】
- internal/usecase/note_command_interactor.go  ← 既存を変更
- internal/driver/config/config.go
- internal/driver/initializer/api/initializer.go

【🚨 最重要の制約1: 呼び出し順序】
Notionを先に呼び、成功したらDBを更新してください。

❌ ダメな順序:
   DBを更新 → Notion呼び出し → 失敗したらDBを戻す

✅ 正しい順序:
   Notion呼び出し → 成功したらDBをトランザクションで更新
   （失敗したらDBを触っていないので、戻す処理が不要）

【🚨 最重要の制約2: トランザクション】
トランザクションの中で Notion API を呼ばないでください。
ロックが長時間保持され、コネクションプールが枯渇します。

【🚨 最重要の制約3: フォールバック】
NOTION_API_KEY が未設定のときは、
連携をスキップして従来どおりの公開処理を行ってください。
（移行期間中のフィーチャーフラグを兼ねます）

【ChangeStatus の振る舞い】
公開する場合:
  1. テンプレートから送り先を取得。未設定ならエラー（公開しない）
  2. Notionにページ作成（既存ページがあればゴミ箱から復元）
  3. 成功したらトランザクションで status 更新 + 連携情報保存
  4. Notion失敗なら status を変更しない

非公開に戻す場合:
  1. Notionページをゴミ箱へ（in_trash: true）
  2. 成功したらトランザクションで status 更新
  3. Notion失敗なら status を変更しない

【Update の振る舞い】
  - 下書きのノート: 従来どおり（Notionは触らない）
  - 公開中のノート: 先にNotionを更新し、成功したらDBを更新

【🚨 確認必須】
PR⓪で追加した既存テストが全てパスすること。
1つでも落ちたら、既存の振る舞いを壊しています。

【参考にすべき既存コード】
internal/usecase/note_command_interactor.go の現在の実装
```

---

## テストも同時に書く

第4章で決めた方針に沿って、**実装と同じPRでテストを書きます**。
「あとでまとめて書く」は、まず実現しません。

### UseCase層（gomock）— 処理の流れを検証

```go
func TestNotionSyncInteractor_Sync(t *testing.T) {
    tests := []struct {
        name      string
        setupMock func(*mock.MockNotionClient, *mock.MockNotionSyncRepository)
        wantErr   error
    }{
        {
            name: "初回連携では CreatePage が呼ばれる",
            setupMock: func(nc *mock.MockNotionClient, sr *mock.MockNotionSyncRepository) {
                sr.EXPECT().Get(gomock.Any(), "note-1").Return(nil, nil)  // 履歴なし
                sr.EXPECT().MarkPending(gomock.Any(), "note-1").Return(nil)
                nc.EXPECT().CreatePage(gomock.Any(), "parent-1", "タイトル", gomock.Any()).
                    Return(&port.PageResult{ID: "page-1", URL: "https://..."}, nil)
                sr.EXPECT().MarkSynced(gomock.Any(), "note-1", "page-1", gomock.Any()).Return(nil)
            },
        },
        {
            name: "2回目は UpdatePage が呼ばれる",
            setupMock: func(nc *mock.MockNotionClient, sr *mock.MockNotionSyncRepository) {
                pageID := "page-1"
                sr.EXPECT().Get(gomock.Any(), "note-1").
                    Return(&notion.Sync{NotionPageID: &pageID}, nil)
                sr.EXPECT().MarkPending(gomock.Any(), "note-1").Return(nil)
                nc.EXPECT().UpdatePage(gomock.Any(), "page-1", gomock.Any(), gomock.Any()).
                    Return(&port.PageResult{ID: "page-1"}, nil)
                sr.EXPECT().MarkSynced(gomock.Any(), "note-1", "page-1", gomock.Any()).Return(nil)
            },
        },
        {
            name: "Notion失敗時も failed が記録される",
            setupMock: func(nc *mock.MockNotionClient, sr *mock.MockNotionSyncRepository) {
                sr.EXPECT().Get(gomock.Any(), "note-1").Return(nil, nil)
                sr.EXPECT().MarkPending(gomock.Any(), "note-1").Return(nil)
                nc.EXPECT().CreatePage(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
                    Return(nil, errors.New("notion is down"))
                sr.EXPECT().MarkFailed(gomock.Any(), "note-1", gomock.Any()).Return(nil)
            },
            wantErr: ErrNotionSyncFailed,
        },
        {
            name:    "他人のノートは連携できない",
            // Notionクライアントを EXPECT しない = 呼ばれたら失敗する
            wantErr: domainerr.ErrUnauthorized,
        },
    }
    // ...
}
```

> **モックの利点: 「呼ばれなかったこと」を検証できる**
> 最後のケースは、`EXPECT()` を書かないことで
> 「Notion APIが呼ばれたらテスト失敗」になります。
> 本物のAPIでは検証しづらい観点です。

**QAケースとの対応:** QA-01, QA-02, QA-10, QA-32

### Gateway層（httptest）— HTTPの扱いを検証

```go
func TestClient_CreatePage(t *testing.T) {
    tests := []struct {
        name       string
        statusCode int
        body       string
        wantErr    bool
        wantCalls  int  // リトライ回数の検証
    }{
        {name: "正常系", statusCode: 200,
         body: `{"id":"page-1","url":"https://notion.so/page-1"}`,
         wantErr: false, wantCalls: 1},

        {name: "401はリトライしない", statusCode: 401,
         wantErr: true, wantCalls: 1},          // 1回で諦める

        {name: "429はリトライする", statusCode: 429,
         wantErr: true, wantCalls: 4},          // 初回 + 3回

        {name: "500はリトライする", statusCode: 500,
         wantErr: true, wantCalls: 4},

        {name: "壊れたJSONはエラー", statusCode: 200,
         body: `{broken`, wantErr: true, wantCalls: 1},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            calls := 0
            srv := httptest.NewServer(http.HandlerFunc(
                func(w http.ResponseWriter, r *http.Request) {
                    calls++
                    if got := r.Header.Get("Notion-Version"); got != "2022-06-28" {
                        t.Errorf("Notion-Version = %q", got)
                    }
                    w.WriteHeader(tt.statusCode)
                    _, _ = w.Write([]byte(tt.body))
                }))
            defer srv.Close()

            client := notion.NewClient("test-token", notion.WithBaseURL(srv.URL))
            _, err := client.CreatePage(context.Background(), "parent", "title", nil)

            if (err != nil) != tt.wantErr {
                t.Errorf("err = %v, wantErr %v", err, tt.wantErr)
            }
            if calls != tt.wantCalls {
                t.Errorf("calls = %d, want %d", calls, tt.wantCalls)
            }
        })
    }
}
```

> 第4章で「ベースURLを差し替え可能にする」と決めておいたので、
> `WithBaseURL(srv.URL)` でモックサーバーを差し込めています。
> **テストしやすい構造を先に決めた成果**です。

#### タイムアウトの検証

```go
func TestClient_Timeout(t *testing.T) {
    srv := httptest.NewServer(http.HandlerFunc(
        func(w http.ResponseWriter, r *http.Request) {
            time.Sleep(2 * time.Second)  // わざと遅くする
        }))
    defer srv.Close()

    client := notion.NewClient("token",
        notion.WithBaseURL(srv.URL),
        notion.WithTimeout(100*time.Millisecond))

    if _, err := client.CreatePage(context.Background(), "p", "t", nil); err == nil {
        t.Fatal("タイムアウトすべきなのにエラーが返らない")
    }
}
```

**QAケースとの対応:** QA-31, QA-33

---

## ⚠️ AIが間違えやすいポイント

実際にAIに任せると、以下をよく間違えます。**先回りして指示に含めてください。**

| よくある間違い | なぜ起きるか | 先回りの指示 |
|---|---|---|
| トランザクション内で外部API | 一般的なコード例がそうなっている | 「🚨 トランザクション外で」と明記 |
| `http.Client{}` をそのまま使う | Goの入門コードがそうなっている | 「Timeout を必ず設定」と明記 |
| 4xx でもリトライする | リトライ実装を雑に書く | 「429と5xxのみ」と明記 |
| テストで本物のAPIを叩く | モックの指示がないと本物を使う | 「httptest を使う」と明記 |
| 新しい設定を必須にする | バリデーションは善だと思っている | 「必須にしない」と明記 |
| 既存の `ChangeStatus` を触る | 「公開時に連携」が自然に見える | 「既存ファイルは触らない」と明記 |
| エラーメッセージにAPIキーを含める | デバッグしやすさを優先する | 「キーをログに出さない」と明記 |

> **この表自体が、レビューのチェックリストになります**

---

## 生成されたコードのレビュー観点

AIが書いたコードは、**必ず人間が確認**します。

### 機械的に確認できること

```bash
# ビルドが通るか
go build ./...

# テストが通るか
go test ./internal/... -cover

# 既存のカバレッジを下回っていないか
# （domain/service は100%、controller は84.7%が基準）

# Lintが通るか
make lint

# 生成物に差分がないか
make sqlc-generate && git diff --exit-code
```

### 人間が見るべきこと

```
□ 🚨 トランザクションの中で外部APIを呼んでいないか
   → grep -n "WithinTransaction" して中身を目で見る

□ タイムアウトが設定されているか
   → grep -n "http.Client" して確認

□ 設定未設定でアプリが起動するか
   → NOTION_API_KEY を外して実際に起動してみる

□ 既存ファイルを勝手に触っていないか
   → git diff --stat で変更ファイル一覧を見る
   → 03_impact_analysis.md の「変更しないもの」と照合

□ エラーハンドリングが握りつぶしになっていないか
   → grep -n "_ = " や "// ignore" を探す

□ ドメイン層が外部に依存していないか
   → internal/domain/ の import 文を目で見る
```

### 特に厳しく見るべき箇所

```
PR④ の notion_sync_interactor.go
  → トランザクション境界（目に見えないバグ）

PR⑤ の config.go
  → 必須チェックを入れていないか

PR⑤ の initializer.go
  → 既存の配線を壊していないか
```

---

## AIに任せてよいこと・ダメなこと

```
┌──────────────────────────────────────────┐
│ ✅ AIに任せてよい                          │
├──────────────────────────────────────────┤
│ ・決まった設計に沿ったコードを書く           │
│ ・既存パターンに揃えて書く                  │
│ ・テストケースを網羅的に書く                │
│ ・ボイラープレート（ファクトリ、モック）      │
│ ・SQLの構文、マイグレーションの書き方        │
└──────────────────────────────────────────┘

┌──────────────────────────────────────────┐
│ ❌ AIに任せてはいけない                     │
├──────────────────────────────────────────┤
│ ・何を作るかを決める                        │
│ ・既存概念を拡張すべきかの判断               │
│ ・データ構造の正規化をどうするか             │
│ ・失敗時にユーザーがどう振る舞うか           │
│ ・影響範囲が許容できるかの判断               │
│ ・リリースしてよいかの判断                  │
└──────────────────────────────────────────┘
```

> **この境界線が、この教材の主題です**
> ①〜④が「AIに任せてはいけないこと」、
> ⑤が「AIに任せてよいこと」です。
>
> **①〜④をやらずに⑤だけAIに投げると、事故ります。**

---

## この章のチェックリスト

```
□ AIに設計ドキュメントを読ませたか
□ PR単位で投げているか（全部まとめて投げていないか）
□ 「参考にすべき既存コード」を指定したか
□ 🚨 の制約（トランザクション、タイムアウト、必須チェック）を明記したか
□ テストの形式と網羅すべきケースを指定したか
□ 第4章のQAケースに対応するテストを書いたか
□ テストを実装と同じPRに含めたか
□ 生成後、git diff --stat で変更範囲を確認したか
□ 「変更しないもの」が本当に変更されていないか照合したか
```

---

## 次の章へ

コードとテストができました。

次はリリース可否の判定と、リリース後に壊れたと気づける仕組みです。

👉 [07_completion_check.md](./07_completion_check.md)
