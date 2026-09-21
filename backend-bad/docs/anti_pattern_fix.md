# Lesson 3: アンチパターン改善ガイド

Lesson 1・2 で「どこが辛いか」を観察したら、次は「じゃあどう直すか？」です。ここでは各アンチパターンに対して **いま何が起きているか → どう直すか → どんな効果があるか** をセットでまとめ、さらに “Before/After” のコード片も添えています。用語はできるだけ噛み砕いているので、迷ったらここに戻ってください。

---

## 0. この設計で運用を続けたときの主なリスク

| リスク | どう困るか |
| --- | --- |
| 変更の波及が常に最大 | OpenAPI 型と sqlc 型を直に握っているため、API/DB の小変更でも Controller/Service すべてを書き換え。成長期に改修コストが跳ね上がる |
| トランザクション事故 | 各メソッドが手書きで Begin/Commit/Rollback。コピペ増殖で Rollback 漏れや二重呼び出しが埋もれる |
| ビジネスルールの所在不明 | if/switch が散乱し、条件変更や追加のたびに全文検索が必要。既存ルールを壊しやすい |
| エラーが原因を隠す | どこでも `failed to ...` で返すため、障害調査やサポートに時間がかかる。クライアントも原因を掴めない |
| チーム拡張に弱い | 境界が曖昧で認知負荷が高く、新メンバーのオンボーディングやレビューが重くなる |
| テスト戦略の崩壊 | モックが差し込めず、常に DB 付き統合テスト頼み。CI 時間が伸び、フィードバックが遅い |

これらを頭に置きながら、以下の改善策を読むと「なぜ直す必要があるのか」が明確になります。

---

## 1. Controller をシンプルにする

| いま起きていること | やること | どう良くなるか |
| --- | --- | --- |
| 1 つの関数に JSON パース / 入力チェック / DB 用データ作成 / レスポンス作成まで全部詰め込んでいる | Controller では「HTTP からアプリ内 DTO（ユースケース用の構造体）に変換する」ところまでにして、DB モデル生成やビジネスルールは UseCase に渡す | Controller の行数が減り、読むのも直すのも楽になる。仕様変更があっても UseCase だけ触ればよくなる |
| OpenAPI の型と DB の型が同じファイルに混在 | `adapter/http` などに「OpenAPI → UseCase DTO」の変換関数を置く | API 定義や DB 変更の影響範囲を小さくできる |
| テストで Echo の Context や `sqldb.Field` まで準備しないとダメ | Controller テストでは「HTTP 入力がアプリ内 DTO に正しく写っているか」だけを見る（例: `/templates?q=foo` が `ListTemplatesInput.Query = \"foo\"` になるか）。UseCase はモックで代用 | テストケースがすっきりし、ユニットテストを書く気力がわく |

#### 今のアンチパターン実装（`internal/controller/template_controller.go`）
```go
func (c *Controller) TemplatesCreateTemplate(ctx echo.Context) error {
    var body openapi.TemplatesCreateTemplateJSONRequestBody
    if err := ctx.Bind(&body); err != nil {
        return respondError(ctx, http.StatusBadRequest, "invalid payload")
    }

    fields := make([]sqldb.Field, len(body.Fields))
    for i, field := range body.Fields {
        fields[i] = sqldb.Field{
            Label:      field.Label,
            Order:      field.Order,
            IsRequired: field.IsRequired,
        }
    }

    template, err := c.templateService.CreateTemplate(ctx, body.OwnerId.String(), body.Name, fields)
    ...
}
```

```go
// Controller 側: HTTP → DTO 変換に専念
func (c *Controller) TemplatesListTemplates(ctx echo.Context) error {
    input := ListTemplatesInput{}
    if q := ctx.QueryParam("q"); q != "" {
        input.Query = &q
    }
    templates, err := c.templateUseCase.ListTemplates(ctx.Request().Context(), input)
    if err != nil {
        return handleError(ctx, err)
    }
    return ctx.JSON(http.StatusOK, templates)
}

// テストでは UseCase をモックし、DTO の値を検証
mockUC.On("ListTemplates", mock.Anything, ListTemplatesInput{Query: ptr("foo")}).Return(nil, nil)
```

---

## 2. Service をフレームワークと DB から引きはがす

| いま起きていること | やること | どう良くなるか |
| --- | --- | --- |
| `CreateTemplate(ctx echo.Context, ...)` のように Echo の型をそのまま受け取っている | Service/UseCase は標準の `context.Context` だけを受け取り（例: `CreateTemplate(ctx context.Context, input TemplateCreateInput)`）、Echo などフレームワークの話は Controller で完結させる | 将来 Fiber や chi に変えても Service のコードは据え置きにできる |
| OpenAPI の型や `sqldb.*` がそのまま登場 | ドメイン独自の構造体（`Template`, `TemplateField` など）を定義し、外部の型は Adapter 層で変換（例: `openapi.TemplatesCreateTemplateJSONRequestBody` → `TemplateCreateInput`） | API/DB 変更が Service まで波及しなくなる |
| Service が直接 `sqlc` のクエリを呼んでいる | Repository インターフェースを用意し、Service からは `TemplateRepository.Create(ctx, Template)` のようなメソッドだけを呼ぶ | DB 切り替えやモック化が簡単になり、テストでも Repository を差し替えるだけで済む |

#### 今のアンチパターン実装（`internal/service/template_service.go`）
```go
func (s *TemplateService) CreateTemplate(ctx echo.Context, ownerID string, name string, fields []sqldb.Field) (*openapi.ModelsTemplateResponse, error) {
    pgOwner, err := parseUUID(ownerID)
    ...
    tx, err := s.pool.Begin(ctx.Request().Context())
    ...
    template, err := txQueries.CreateTemplate(ctx.Request().Context(), &sqldb.CreateTemplateParams{ ... })
    ...
    _, err = txQueries.CreateField(ctx.Request().Context(), &sqldb.CreateFieldParams{ ... })
    ...
    if err := tx.Commit(ctx.Request().Context()); err != nil {
        if rbErr := tx.Rollback(ctx.Request().Context()); rbErr != nil {
            return nil, rbErr
        }
        return nil, err
    }
    return s.composeTemplateResponse(ctx.Request().Context(), template.ID, template.Name, template.OwnerID, template.UpdatedAt, false)
}
```

```go
type TemplateUseCase struct {
    repo TemplateRepository
}

func (u *TemplateUseCase) CreateTemplate(ctx context.Context, input TemplateCreateInput) error {
    tpl := Template{
        OwnerID: input.OwnerID,
        Name:    input.Name,
        Fields:  input.Fields,
    }
    return u.repo.Create(ctx, tpl)
}

// Adapter 層で OpenAPI → TemplateCreateInput へ変換
func toTemplateCreateInput(body openapi.TemplatesCreateTemplateJSONRequestBody) TemplateCreateInput {
    // ここで詰め替え
}
```

---

## 3. トランザクションを共通化する

| いま起きていること | やること | どう良くなるか |
| --- | --- | --- |
| どのメソッドも `Begin` → `Rollback` → `Commit` を手書きしている | 「トランザクション付きで処理を実行する関数」を 1 つ作り、UseCase からは `RunInTx(ctx, func(ctx context.Context) error { ... })` のように呼ぶだけにする | Rollback の書き忘れやコピペが減る。テストでもこの共通関数だけ確認すれば安心できる |
| `sqlc` の `WithTx` を直接触っていてテストが難しい | Repository も「トランザクション中です」と分かるコンテキストを受け取る設計にする。モックでは任意のエラーや成功を再現できるようにする | 「途中で失敗したらロールバックされるべき」といったケースをユニットテストで書ける |

#### 今のアンチパターン実装（`internal/service/note_service.go`）
```go
func (s *NoteService) ChangeStatus(ctx echo.Context, id, status string) (*openapi.ModelsNoteResponse, error) {
    ...
    tx, err := s.pool.Begin(ctx.Request().Context())
    ...
    note, err := txQueries.UpdateNoteStatus(ctx.Request().Context(), &sqldb.UpdateNoteStatusParams{
        ID:     noteID,
        Status: status,
    })
    if err != nil {
        if rbErr := tx.Rollback(ctx.Request().Context()); rbErr != nil {
            return nil, rbErr
        }
        ...
    }
    if err := tx.Commit(ctx.Request().Context()); err != nil {
        if rbErr := tx.Rollback(ctx.Request().Context()); rbErr != nil {
            return nil, rbErr
        }
        return nil, err
    }
    return s.composeNoteResponse(ctx.Request().Context(), noteResponseInput{ ... })
}
```

```go
func RunInTx(ctx context.Context, db DB, fn func(ctx context.Context) error) error {
    tx, err := db.Begin(ctx)
    if err != nil {
        return err
    }
    txCtx := context.WithValue(ctx, txKey{}, tx)
    if err := fn(txCtx); err != nil {
        _ = tx.Rollback(txCtx)
        return err
    }
    return tx.Commit(txCtx)
}

func (u *TemplateUseCase) UpdateTemplate(ctx context.Context, input TemplateUpdateInput) error {
    return RunInTx(ctx, u.db, func(txCtx context.Context) error {
        return u.repo.Update(txCtx, input)
    })
}
```

---

## 4. if/switch の沼から抜ける

| いま起きていること | やること | どう良くなるか |
| --- | --- | --- |
| ステータスや検索条件を文字列と if でベタ書き | `NoteStatus` などの列挙型や値オブジェクトを定義し、状態遷移をメソッドで表現する | 新ステータスを追加しても修正箇所がすぐ分かる。テストも「状態遷移が正しいか」だけに集中できる |
| Controller と Service で同じ条件分岐が二重に登場 | UseCase 内で検索条件オブジェクトを作り、Repository にまとめて渡す。Controller はリクエスト値を詰めるだけにする | 分岐の場所が整理され、レビューしやすい |

#### 今のアンチパターン実装（`internal/service/note_service.go`）
```go
if status != "Draft" && status != "Publish" {
    return nil, errors.New("invalid status")
}
...
if filters.OwnerID != nil && *filters.OwnerID != "" {
    if id, err := parseUUID(*filters.OwnerID); err == nil {
        params.Column1 = id
    }
}
if filters.Query != nil && *filters.Query != "" {
    params.Column2 = *filters.Query
}
```

```go
type NoteStatus string

const (
    NoteStatusDraft  NoteStatus = "draft"
    NoteStatusPublic NoteStatus = "public"
)

func (s NoteStatus) CanPublish() bool {
    return s == NoteStatusDraft
}

func (u *NoteUseCase) ChangeStatus(ctx context.Context, id string, to NoteStatus) error {
    note, _ := u.repo.Get(ctx, id)
    if !note.Status.CanPublish() {
        return ErrInvalidTransition
    }
    note.Status = to
    return u.repo.Save(ctx, note)
}
```

---

## 5. エラーの出し方を決める

| いま起きていること | やること | どう良くなるか |
| --- | --- | --- |
| どんなエラーでも `failed to ...` など曖昧な文字列 | `AppError` のような自前のエラー型（`Code`, `Message`, `Status` を持つ struct）を作り、UseCase から返す（例: `Code=\"TemplateNotFound\"`） | Controller は `AppError` を見て HTTP レスポンスを決めるだけで済む。クライアントも原因を把握しやすい |
| エラーハンドリングが各所に分散 | Echo/Fiber のエラーハンドラーを 1 カ所に置き、`AppError` 以外は 500 にフォールバックする | エラー仕様を変えたいときに 1 ファイルだけ直せばよい |

#### 今のアンチパターン実装（`internal/controller/helpers.go`）
```go
func respondError(ctx echo.Context, status int, message string) error {
    return ctx.JSON(status, openapi.ModelsErrorResponse{
        Code:    http.StatusText(status),
        Message: message,
    })
}
```

```go
type AppError struct {
    Code    string
    Message string
    Status  int
}

var ErrTemplateNotFound = &AppError{
    Code:    "TemplateNotFound",
    Message: "テンプレートが見つかりません",
    Status:  http.StatusNotFound,
}

func respondError(ctx echo.Context, err error) error {
    if appErr, ok := err.(*AppError); ok {
        return ctx.JSON(appErr.Status, appErr)
    }
    return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "internal error"})
}
```

---

## 6. テストの役割分担をはっきりさせる

| いま起きていること | やること | どう良くなるか |
| --- | --- | --- |
| Service が Echo/DB と直結していてユニットテスト不能 | Repository やトランザクション管理のインターフェースをモックに差し替え、UseCase が持つロジックだけをテストする（例: `CreateTemplate` が不正な ownerID なら `ErrInvalidAccountID` を返すか） | 正常系・異常系をテーブルドリブンで網羅できる |
| 正常系の確認には毎回 DB を立てる必要がある | Integration/E2E テストは testcontainers などで別途実行する方針を決める。ユニットテストとは目的を切り分ける | CI の実行時間と開発者の手間が減る |
| 「本当はここまでテストしたかった」がコメントで終わっている | 設計を分割し直し、コメントで書いたケースを実際のテストコードとして追加する | 教材として「改善するとこういうテストが書ける」を示せる |

```go
func TestTemplateUseCase_Create(t *testing.T) {
    repo := new(MockTemplateRepository)
    repo.On("Create", mock.Anything, mock.Anything).Return(nil)

    uc := TemplateUseCase{repo: repo}
    err := uc.CreateTemplate(context.Background(), TemplateCreateInput{
        OwnerID: mustUUID("11111111-1111-1111-1111-111111111111"),
        Name:    "Template",
    })

    require.NoError(t, err)
    repo.AssertCalled(t, "Create", mock.Anything, mock.MatchedBy(func(tpl Template) bool {
        return tpl.Name == "Template"
    }))
}
```

---

## 7. 進め方のサンプル

1. **境界を作る**: Controller → UseCase → Repository の順に責務を分離する。
2. **依存の矢印をそろえる**: UseCase から見るとインターフェースしか見えないようにし、実装（Echo/sqlc）は Adapter に閉じ込める。
3. **トランザクションを共通化**: `RunInTx` のような関数を経由して Begin/Commit を一元管理する。
4. **テストから逆算**: 「ユニットテストで保証したい条件」を箇条書きし、それをモックで再現できるように設計を調整する。
5. **エラーの型を決める**: 共通のエラー struct を定義し、Controller はそれを HTTP 応答へ変換するだけにする。

この順序で進めれば、アンチパターンで紹介した問題点を少しずつ解消しながら、テストしやすく変更に強い構成に近づけられます。
