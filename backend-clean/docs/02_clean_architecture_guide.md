# 🎯 クリーンアーキテクチャの処理の流れ - 小学生でもわかる完全ガイド

この図、見たことありますか？

![Clean Architecture](https://blog.cleancoder.com/uncle-bob/images/2012-08-13-the-clean-architecture/CleanArchitecture.jpg)

**クリーンアーキテクチャで最も有名な「同心円の図」**です。

でも、**「これ、実際どう動くの？」**って思いませんか？

このドキュメントでは、この図を使って、**小学生でもわかるレベル**で処理の流れを解説します。

---

## 📚 目次

1. [同心円の図って何？](#同心円の図って何)
2. [例え話: 学校で理解する](#例え話-学校で理解する)
3. [インターフェースって何？](#インターフェースって何)
4. [実際のリクエストの流れ](#実際のリクエストの流れ)
5. [Factoryって何？](#factoryって何)
6. [どのメソッドがどこを呼ぶのか](#どのメソッドがどこを呼ぶのか)

---

## 同心円の図って何？

```
        🎯 クリーンアーキテクチャの同心円

    ┌─────────────────────────────────────┐
    │   外側: Frameworks & Drivers        │  ← DB、Web、外部ツール
    │  ┌───────────────────────────────┐  │
    │  │  中: Interface Adapters       │  │  ← Controller、Gateway、Presenter
    │  │ ┌─────────────────────────┐   │  │
    │  │ │ 内: Application Logic   │   │  │  ← UseCase
    │  │ │ ┌───────────────────┐   │   │  │
    │  │ │ │ 中心: Entities    │   │   │  │  ← Domain
    │  │ │ └───────────────────┘   │   │  │
    │  │ └─────────────────────────┘   │  │
    │  └───────────────────────────────┘  │
    └─────────────────────────────────────┘

重要なルール:
矢印は「内側に向かってのみ」許可される
外側は内側を知ってOK、内側は外側を知っちゃダメ
```

### 各層の役割

| 層 | 何をする？ | 例 |
|----|-----------|-----|
| **Entities（中心）** | ビジネスルール | 「ノートは公開前に必須フィールドを埋める必要がある」 |
| **Use Cases（内側）** | アプリケーションの手順 | 「1. ノート取得 2. 検証 3. 保存 4. 通知」 |
| **Interface Adapters（中間）** | 変換 | 「HTTP → Domain型」「Domain型 → DB型」 |
| **Frameworks & Drivers（外側）** | 外部ツール | PostgreSQL、Echo、OpenAPI |

---

## 例え話: 学校で理解する

クリーンアーキテクチャを**学校**で例えます。

### 🏫 学校の組織図

```
┌─────────────────────────────────────────────┐
│        🏫 学校                               │
├─────────────────────────────────────────────┤
│                                             │
│  📞 受付（Controller）                       │
│    ↓ 「保護者から電話が来た」                │
│                                             │
│  📋 事務室（UseCase）                        │
│    ↓ 「どう対応するか決める」                │
│                                             │
│  📖 校則（Domain/ビジネスルール）            │
│    ↓ 「何が許されるか」                      │
│                                             │
│  🗂️ 倉庫（Repository/DB）                   │
│    「生徒の情報を保管」                      │
│                                             │
│  📣 掲示板（Presenter）                      │
│    「結果を保護者に伝える形式に整える」      │
└─────────────────────────────────────────────┘
```

### シナリオ: 「生徒の成績を見たい」

1. **📞 受付（Controller）が電話を受ける**
   ```
   保護者: 「うちの子の成績を教えてください」
   受付: 「はい、少々お待ちください」
   ```

2. **📋 事務室（UseCase）が対応を決める**
   ```
   事務室: 「この保護者は誰の親？」
   事務室: 「校則で確認できる？」 → 📖 校則をチェック
   事務室: 「成績を取ってきて」 → 🗂️ 倉庫に依頼
   ```

3. **📖 校則（Domain）をチェック**
   ```
   校則: 「保護者は自分の子の成績だけ見れる」
   校則: 「他人の子の成績は見れない」
   ```

4. **🗂️ 倉庫（Repository）が情報を取得**
   ```
   倉庫: 「生徒ID=123の成績を探す」
   倉庫: 「見つけた！国語90点、算数85点」
   ```

5. **📋 事務室（UseCase）が掲示板に渡す**
   ```
   事務室: 「この成績を保護者に伝えられる形にして」
   ```

6. **📣 掲示板（Presenter）が整形**
   ```
   掲示板: 「国語: 90点、算数: 85点」
   掲示板: 「保護者にわかりやすい形式で準備完了」
   ```

7. **📞 受付（Controller）が伝える**
   ```
   受付: 「お待たせしました。お子様の成績は...」
   ```

### ポイント

```
✅ 受付（Controller）は校則（Domain）を知らない
   → 事務室（UseCase）に丸投げ

✅ 校則（Domain）は電話のことを知らない
   → 純粋なルールだけ

✅ 倉庫（Repository）は校則に従って動く
   → 「誰でも見れる」じゃなくて「保護者だけ」

✅ それぞれの役割が明確
   → 受付が倉庫を直接見に行かない
```

---

## インターフェースって何？

インターフェースは**「約束事」**です。

### 🤝 例え: レストランのメニュー

```
┌─────────────────────────────┐
│      📋 メニュー（Interface） │
│                             │
│  - ピザを焼く                │
│  - パスタを作る              │
│  - サラダを盛る              │
└─────────────────────────────┘
         ↑ 「こういう料理を提供する」という約束
         │
    ┌────┴─────┐
    │          │
イタリア人    日本人
シェフ        シェフ
```

**メニュー（Interface）があれば:**
- お客さんは「誰が作るか」を気にしなくていい
- シェフを交代しても、メニューが同じならOK
- 新人シェフでも、メニュー通り作ればOK

### 💻 コードで見る

```go
// 📋 これがインターフェース（約束事）
type NoteRepository interface {
    Get(ctx context.Context, id string) (*note.Note, error)
    Create(ctx context.Context, n note.Note) error
    Delete(ctx context.Context, id string) error
}

// 👨‍🍳 これが実装その1（PostgreSQLシェフ）
type PostgresNoteRepository struct {
    pool *pgxpool.Pool
}

func (r *PostgresNoteRepository) Get(ctx context.Context, id string) (*note.Note, error) {
    // PostgreSQLから取得
    row, err := r.pool.QueryRow(ctx, "SELECT * FROM notes WHERE id = $1", id)
    // ...
}

// 👨‍🍳 これが実装その2（MySQLシェフ）
type MySQLNoteRepository struct {
    db *sql.DB
}

func (r *MySQLNoteRepository) Get(ctx context.Context, id string) (*note.Note, error) {
    // MySQLから取得
    row, err := r.db.QueryRowContext(ctx, "SELECT * FROM notes WHERE id = ?", id)
    // ...
}

// 👨‍🍳 これが実装その3（テスト用の偽物シェフ）
type MockNoteRepository struct {
    notes map[string]*note.Note
}

func (r *MockNoteRepository) Get(ctx context.Context, id string) (*note.Note, error) {
    // メモリから取得（DBアクセスなし）
    return r.notes[id], nil
}
```

**UseCaseはどれを使ってる？**

```go
type NoteInteractor struct {
    repo NoteRepository  // ← インターフェース（約束事）だけ知ってる
}

func (u *NoteInteractor) Get(ctx context.Context, id string) error {
    // どのシェフが作ってるか知らない！
    // でも約束通りに料理（データ）が出てくる
    note, err := u.repo.Get(ctx, id)
    // ...
}
```

**これがインターフェースの力！**

```
✅ PostgreSQL → MySQLに変更しても、UseCaseは変更不要
✅ テスト時はMockを使える（DB起動不要）
✅ 新しいDB追加も簡単（インターフェース実装するだけ）
```

---

## 実際のリクエストの流れ

では、実際に **GET /api/notes/123** というリクエストが来たら何が起こるか、詳しく見ていきます。

### 🚀 全体の流れ

```
🌍 ユーザー
    │
    │ GET /api/notes/123
    ↓
┌─────────────────────────────────────────┐
│  ① Echo (Webフレームワーク)              │  ← 外側: Framework
│     「リクエストが来た！Controllerへ」   │
└─────────────────────────────────────────┘
    │
    ↓
┌─────────────────────────────────────────┐
│  ② Controller (受付窓口)                 │  ← 中間: Interface Adapter
│     「ID=123のノート欲しいって」         │
│     「UseCaseに頼もう」                  │
└─────────────────────────────────────────┘
    │
    │ input.Get(ctx, "123")
    ↓
┌─────────────────────────────────────────┐
│  ③ UseCase (事務室)                      │  ← 内側: Application Logic
│     「1. Repositoryに取得依頼」          │
│     「2. Presenterに渡す」               │
└─────────────────────────────────────────┘
    │                           ↑
    │ Get("123")                │ 取得したNote
    ↓                           │
┌─────────────────────────────────────────┐
│  ④ Repository (倉庫)                     │  ← 中間: Interface Adapter
│     「DBから取得します」                 │
│     「DB型 → Domain型に変換」            │
└─────────────────────────────────────────┘
    │
    │ SELECT * FROM notes WHERE id = '123'
    ↓
┌─────────────────────────────────────────┐
│  ⑤ PostgreSQL (データベース)            │  ← 外側: Framework
│     「データ返すよ」                     │
└─────────────────────────────────────────┘
    │
    │ Row{id, title, ...}
    ↓
┌─────────────────────────────────────────┐
│  ⑥ Repository (倉庫)                     │
│     「pgx.Row → note.Note に変換」       │
└─────────────────────────────────────────┘
    │
    │ note.Note
    ↓
┌─────────────────────────────────────────┐
│  ⑦ UseCase (事務室)                      │
│     「Presenterに渡す」                  │
└─────────────────────────────────────────┘
    │
    │ PresentNote(note)
    ↓
┌─────────────────────────────────────────┐
│  ⑧ Presenter (掲示板)                    │  ← 中間: Interface Adapter
│     「note.Note → OpenAPI型に変換」      │
│     「結果を保存」                        │
└─────────────────────────────────────────┘
    │
    ↓
┌─────────────────────────────────────────┐
│  ⑨ Controller (受付窓口)                 │
│     「Presenterから結果を取り出す」      │
│     「HTTPレスポンスとして返す」         │
└─────────────────────────────────────────┘
    │
    │ HTTP 200 OK + JSON
    ↓
🌍 ユーザー
```

### 📝 コードで見る（抜粋）

#### ① Echo → Controller

```go
// Echo が自動でルーティング
// GET /api/notes/:id → NoteController.GetByID を呼ぶ
```

#### ② Controller

```go
// internal/adapter/http/controller/note_controller.go
func (c *NoteController) GetByID(ctx echo.Context, noteID string) error {
    // 1. Factory経由でUseCaseとPresenterを作る
    input, p := c.newIO()

    // 2. UseCaseを呼ぶ（詳細は知らない、丸投げ）
    if err := input.Get(ctx.Request().Context(), noteID); err != nil {
        return handleError(ctx, err)
    }

    // 3. Presenterから結果を取得
    return ctx.JSON(http.StatusOK, p.Note())
}

func (c *NoteController) newIO() (port.NoteInputPort, *presenter.NotePresenter) {
    // Factoryを呼んで新しいインスタンスを作る（後で詳しく説明）
    output := c.outputFactory()
    input := c.inputFactory(c.noteRepoFactory(), c.tplRepoFactory(), c.txFactory(), output)
    return input, output
}
```

#### ③ UseCase

```go
// internal/usecase/note_interactor.go
func (u *NoteInteractor) Get(ctx context.Context, id string) error {
    // 1. Repositoryに取得依頼（インターフェース経由）
    n, err := u.notes.Get(ctx, id)
    if err != nil {
        return err
    }

    // 2. Presenterに渡す（インターフェース経由）
    return u.output.PresentNote(ctx, n)
}
```

#### ④ Repository

```go
// internal/adapter/gateway/db/sqlc/note_repository.go
func (r *NoteRepository) Get(ctx context.Context, id string) (*note.WithMeta, error) {
    // 1. string → pgtype.UUID に変換
    pgID, err := toUUID(id)
    if err != nil {
        return nil, err
    }

    // 2. sqlc でDB問い合わせ
    row, err := queriesForContext(ctx, r.queries).GetNoteByID(ctx, pgID)
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, domainerr.ErrNotFound  // DBエラー → ドメインエラーに変換
        }
        return nil, err
    }

    // 3. Sections取得
    sections, err := r.listSections(ctx, row.ID)
    if err != nil {
        return nil, err
    }

    // 4. DB行 → ドメインモデルに変換
    return &note.WithMeta{
        Note: note.Note{
            ID:         uuidToString(row.ID),       // pgtype.UUID → string
            Title:      row.Title,
            TemplateID: uuidToString(row.TemplateID),
            OwnerID:    uuidToString(row.OwnerID),
            Status:     note.NoteStatus(row.Status),
            CreatedAt:  timestamptzToTime(row.CreatedAt),
            UpdatedAt:  timestamptzToTime(row.UpdatedAt),
        },
        TemplateName:   row.TemplateName,
        OwnerFirstName: row.FirstName,
        OwnerLastName:  row.LastName,
        OwnerThumbnail: thumbnail,
        Sections:       sections,
    }, nil
}
```

#### ⑧ Presenter

```go
// internal/adapter/http/presenter/note_presenter.go
func (p *NotePresenter) PresentNote(_ context.Context, n *note.WithMeta) error {
    // ドメインモデル → OpenAPI型に変換
    resp := toNoteResponse(*n)
    p.note = &resp  // 結果を保存
    return nil
}

func toNoteResponse(n note.WithMeta) openapi.ModelsNoteResponse {
    sections := make([]openapi.ModelsSection, 0, len(n.Sections))
    for _, s := range n.Sections {
        sections = append(sections, openapi.ModelsSection{
            Id:         s.Section.ID,
            FieldId:    s.Section.FieldID,
            FieldLabel: s.FieldLabel,
            Content:    s.Section.Content,
            IsRequired: s.IsRequired,
        })
    }

    return openapi.ModelsNoteResponse{
        Id:           n.Note.ID,
        Title:        n.Note.Title,
        TemplateId:   n.Note.TemplateID,
        TemplateName: n.TemplateName,
        OwnerId:      n.Note.OwnerID,
        Owner: openapi.ModelsAccountSummary{
            Id:        n.Note.OwnerID,
            FirstName: n.OwnerFirstName,
            LastName:  n.OwnerLastName,
            Thumbnail: n.OwnerThumbnail,
        },
        Status:    openapi.ModelsNoteStatus(n.Note.Status),
        Sections:  sections,
        CreatedAt: n.Note.CreatedAt,
        UpdatedAt: n.Note.UpdatedAt,
    }
}
```

### 🔄 データの変換の流れ

```
HTTP Request (JSON)
    │
    │ {"id": "123"}
    ↓
Controller: string "123"
    │
    │ noteID: string
    ↓
UseCase: string "123"
    │
    │ id: string
    ↓
Repository: pgtype.UUID
    │
    │ SELECT * FROM notes WHERE id = $1
    ↓
PostgreSQL: Row
    │
    │ {id: uuid, title: text, ...}
    ↓
Repository: note.WithMeta (Domain)
    │
    │ Note{ID: string, Title: string, ...}
    ↓
UseCase: note.WithMeta
    │
    │ 渡すだけ
    ↓
Presenter: openapi.ModelsNoteResponse
    │
    │ {id: string, title: string, ...}
    ↓
Controller: HTTP Response (JSON)
    │
    │ {"id": "123", "title": "..."}
    ↓
HTTP Response
```

---

## Factoryって何？

Factoryは**「工場」**です。

### 🏭 例え: おもちゃ工場

```
┌─────────────────────────────────┐
│      🏭 おもちゃ工場              │
│                                 │
│  設計図を見て、毎回新しい        │
│  おもちゃを作る                  │
│                                 │
│  設計図:                         │
│  - 車のおもちゃの作り方          │
│  - 人形の作り方                  │
│  - ブロックの作り方              │
└─────────────────────────────────┘
```

**工場の仕事:**
1. 注文が来る「車のおもちゃ1つください」
2. 設計図を見る
3. **新しい**おもちゃを作る
4. 渡す

### 💻 なぜFactoryが必要？

#### ❌ Factoryなしの場合

```go
// ❌ 悪い例: 起動時に1つだけ作る
func main() {
    // 起動時に1つだけPresenterを作る
    presenter := &NotePresenter{}

    // 全てのリクエストで同じPresenterを使い回す
    controller := &NoteController{
        presenter: presenter,  // ← 使い回し
    }

    // リクエスト1
    // presenter.note = ノートA
    // リクエスト2
    // presenter.note = ノートB  ← 上書き！
    // リクエスト1のレスポンス
    // return presenter.note  ← ノートBが返る！バグ！
}
```

#### ✅ Factoryありの場合

```go
// ✅ 良い例: リクエストごとに新しいPresenterを作る
type NoteController struct {
    // 工場（設計図）を保存
    presenterFactory func() *NotePresenter
}

func (c *NoteController) GetByID(ctx echo.Context, noteID string) error {
    // リクエストごとに新しいPresenterを作る
    presenter := c.presenterFactory()  // ← 工場を呼ぶ

    // このリクエスト専用のPresenterだから安全
    // presenter.note = ノートA
    // return presenter.note  // ノートAが返る
}
```

### 🔧 Factoryの実装

```go
// ① 工場の定義
// internal/driver/factory/presenter_factory.go
func NewNotePresenterFactory() func() *presenter.NotePresenter {
    return func() *presenter.NotePresenter {
        // 新しいPresenterを作る
        return presenter.NewNotePresenter()
    }
}

// ② Controllerに工場を渡す
// internal/driver/initializer/api/initializer.go
func BuildServer(ctx context.Context) (*echo.Echo, func(), error) {
    // ...

    // 工場を作る
    presenterFactory := factory.NewNotePresenterFactory()

    // Controllerに工場を渡す（おもちゃ本体じゃなくて、設計図を渡す）
    noteCtrl := controller.NewNoteController(
        noteInputFactory,
        presenterFactory,  // ← 設計図
        noteRepoFactory,
        tplRepoFactory,
        txFactory,
    )

    // ...
}

// ③ Controllerは工場を保存
type NoteController struct {
    presenterFactory func() *presenter.NotePresenter  // ← 設計図を保存
}

// ④ リクエストごとに工場を呼ぶ
func (c *NoteController) GetByID(ctx echo.Context, noteID string) error {
    // 工場を呼んで新しいPresenterを作る
    presenter := c.presenterFactory()  // ← 新しいおもちゃを作る

    // ...
}
```

---

## どのメソッドがどこを呼ぶのか

### 📊 呼び出し関係の全体図

```
main.go
  │
  │ BuildServer()
  ↓
Initializer
  │
  │ NewXxxFactory() で工場を作る
  │ NewXxxController() でControllerを作る
  ↓
Controller
  │
  │ リクエストごとに:
  │ 1. Factory を呼ぶ → UseCase と Presenter を作る
  │ 2. UseCase.Get() を呼ぶ
  │ 3. Presenter.Note() で結果を取得
  ↓
UseCase (Interactor)
  │
  │ 1. Repository.Get() を呼ぶ（インターフェース経由）
  │ 2. Presenter.PresentNote() を呼ぶ（インターフェース経由）
  ↓
Repository          Presenter
  │                   │
  │ DB問い合わせ      │ 結果を保存
  ↓                   ↓
PostgreSQL          内部フィールド
```

### 🔍 詳細な呼び出しチェーン

#### 1. アプリ起動時

```
main()
  → BuildServer()
    → NewPool() (DB接続)
    → NewNoteInputFactory() (工場を作る)
    → NewNotePresenterFactory() (工場を作る)
    → NewNoteRepositoryFactory() (工場を作る)
    → NewTxFactory() (工場を作る)
    → NewNoteController(工場たちを渡す)
    → NewServer(Controllerたちを渡す)
    → Echo起動
```

**ポイント:**
- 起動時は**工場を作るだけ**
- 実際のUseCase/Presenter/Repositoryは**まだ作らない**

#### 2. リクエスト受信時

```
Echo がリクエスト受信
  → ルーティング確認
  → NoteController.GetByID() を呼ぶ
```

#### 3. Controller内部

```
NoteController.GetByID(ctx, noteID)
  │
  │ このリクエスト専用のインスタンスを作る
  ↓
newIO() を呼ぶ
  │
  ├─ outputFactory() → 新しいPresenterを作る
  │
  ├─ noteRepoFactory() → 新しいRepositoryを作る
  │
  ├─ tplRepoFactory() → 新しいRepositoryを作る
  │
  ├─ txFactory() → 新しいTxManagerを作る
  │
  └─ inputFactory(repos..., output) → 新しいUseCaseを作る
       │
       └─ UseCase に Repository と Presenter を注入
  │
  ↓
input.Get(ctx, noteID) を呼ぶ
  │
  ↓
p.Note() で結果を取得
  │
  ↓
ctx.JSON(200, result) でHTTPレスポンス
```

#### 4. UseCase内部

```
NoteInteractor.Get(ctx, id)
  │
  ├─ u.notes.Get(ctx, id) → Repository を呼ぶ
  │    │
  │    └─ PostgreSQL問い合わせ
  │         DB行 → ドメインモデルに変換
  │         return note.WithMeta
  │
  ├─ u.output.PresentNote(ctx, note) → Presenter を呼ぶ
  │    │
  │    └─ ドメインモデル → OpenAPI型に変換
  │         p.note = result (保存)
  │
  └─ return nil
```

#### 5. Controller に戻る

```
Controller
  │
  ├─ input.Get() が完了
  │
  ├─ p.Note() で結果を取得
  │    │
  │    └─ Presenterが保存していた結果を返す
  │
  └─ ctx.JSON(200, result)
       │
       └─ Echo が HTTPレスポンスとして返す
```

### 🎯 重要なポイント

```
✅ 起動時:
   工場を作る（設計図を用意）
   ↓
   Controllerに工場を渡す
   ↓
   待機

✅ リクエストごと:
   工場を呼ぶ
   ↓
   新しいインスタンスを作る（UseCase、Presenter、Repository）
   ↓
   UseCaseに依存を注入
   ↓
   処理
   ↓
   結果を返す
   ↓
   インスタンスは破棄（次のリクエストでまた新しく作る）
```

---

## 🎓 まとめ

### 同心円の図の本質

```
┌──────────────────────────────────────┐
│  外側: 変わりやすいもの               │
│  - DB (PostgreSQL → MySQL)           │
│  - Web (Echo → Gin)                  │
│  - API (REST → GraphQL)              │
└──────────────────────────────────────┘
         ↓ 影響を受けない
┌──────────────────────────────────────┐
│  内側: 変わりにくいもの               │
│  - ビジネスルール                     │
│  - ドメインモデル                     │
└──────────────────────────────────────┘
```

### 処理の流れの本質

```
1. Controller: 外の世界 → 内の世界に変換
2. UseCase: 手順書通りに処理
3. Repository: 内の世界 → DBに変換
4. Presenter: 内の世界 → 外の世界に変換
```

### インターフェースの本質

```
約束事だけ決めておけば、実装は簡単に差し替えられる
→ テスト簡単、DB変更簡単、保守簡単
```

### Factoryの本質

```
リクエストごとに新しいインスタンスを作ることで、
状態が混ざらない、安全な並行処理ができる
```

---

## 💡 よくある質問と答え

### Q1: なんでこんなに複雑にするの？シンプルじゃダメなの？

**A:** 小さいアプリなら確かにオーバーエンジニアリングです。でも：

```
小規模 (1人、1ヶ月):
┌──────────────┐
│   main.go    │
│  全部ここ！  │ ← これでOK
└──────────────┘

中〜大規模 (チーム、長期保守):
┌──────────────────────────────────┐
│  クリーンアーキテクチャ            │
│  ・テストしやすい                  │
│  ・変更に強い                      │
│  ・チーム開発しやすい               │
│  ・フレームワーク変更できる         │
└──────────────────────────────────┘
```

**このプロジェクトは教材なので、わざと丁寧に分けています。**

---

### Q2: どのレイヤーにテストを書けばいい？

**A:** 全部！でも優先順位は：

```
最優先: ⭐⭐⭐
├─ Domain (ビジネスルール)
│  → ここが壊れたら致命的！
│
├─ UseCase (手順)
│  → ビジネスロジックのテスト
│
次点: ⭐⭐
├─ Presenter (変換)
│  → レスポンス形式が正しいか
│
├─ Gateway (DB変換)
│  → ドメインモデル変換が正しいか
│
最後: ⭐
└─ Controller (HTTPハンドラ)
   → 結合テストやE2Eで確認
```

---

### Q3: Repositoryのインターフェースはどこに置く？

**A:** このプロジェクトでは **Port（`internal/port`）** に置いています。

```
Option A: Portに置く (このプロジェクト)
┌────────────────┐
│   port/        │
│  - NoteRepository interface
│  - TemplateRepository interface
└────────────────┘
         ↑
         │ 実装（ORM別）
┌──────────────────────┐
│ gateway/db/sqlc/     │  ← sqlc実装
│  - NoteRepository struct
│ gateway/db/gorm/    │  ← GORM実装（準備済み）
│  - AccountRepository struct
└──────────────────────┘

Option B: Domainに置く
┌────────────────┐
│  domain/note/  │
│  - Note entity
│  - NoteRepository interface  ← ここ
└────────────────┘

どっちでもOK！プロジェクトで決めよう。
```

---

### Q4: トランザクションはどこで始める？

**A:** **UseCase で開始**します。

```go
func (u *NoteInteractor) Create(ctx context.Context, input port.NoteCreateInput) error {
    // ✅ UseCaseでトランザクション開始
    err := u.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
        // この中で複数のRepositoryを呼ぶ
        note, err := u.notes.Create(txCtx, ...)
        if err != nil {
            return err  // 自動ロールバック
        }

        err = u.notes.ReplaceSections(txCtx, ...)
        if err != nil {
            return err  // 自動ロールバック
        }

        return nil  // 正常終了 → コミット
    })
    return err
}
```

**なぜ？**
- **ビジネスロジックの境界 = トランザクションの境界**
- Repositoryは「データを取る・保存する」だけに集中

---

### Q5: エラーハンドリングはどうする？

```
ドメインエラー (internal/domain/errors/)
    ↓
ErrNotFound, ErrUnauthorized, ErrTitleRequired...
    ↓
Controller (handleError関数)
    ↓
HTTPステータスコードに変換
    ↓
404 Not Found, 401 Unauthorized, 400 Bad Request
```

**重要:** Gatewayで**DBエラー → ドメインエラー**に変換する！

```go
// ❌ NG: pgx.ErrNoRowsをそのまま返す
if err != nil {
    return nil, err
}

// ✅ OK: ドメインエラーに変換
if errors.Is(err, pgx.ErrNoRows) {
    return nil, domainerr.ErrNotFound
}
```

---

### Q6: HTTPとgRPCを両方サポートする場合の設計は？

**A:** **このプロジェクトではHTTPとgRPCの両方をサポート**しています。

```
UseCase (ビジネスロジック) - 1つだけ
    ↑                    ↑
    │                    │
HTTP Presenter      gRPC Presenter
HTTP Controller     gRPC Controller
    ↓                    ↓
REST API            gRPC API
```

**ポイント:**
- UseCaseは1つだけ（共通のビジネスロジック）
- Presenter/Controllerをプロトコルごとに用意
- Factoryもプロトコルごとに分ける（`factory/http/`、`factory/grpc/`）

```go
// HTTP Presenterの例
// internal/adapter/http/presenter/account_presenter.go
func (p *AccountPresenter) PresentAccount(ctx context.Context, acc *account.Account) error {
    // OpenAPI型に変換
    p.account = toOpenAPIAccount(acc)
    return nil
}

// gRPC Presenterの例
// internal/adapter/grpc/presenter/account_presenter.go
func (p *AccountPresenter) PresentAccount(ctx context.Context, acc *account.Account) error {
    // Protobuf型に変換
    p.response = toProtobufAccount(acc)
    return nil
}
```

**メリット:**
- ビジネスロジックを重複させない
- プロトコル固有の変換は各Adapterで吸収
- テストも共通化できる

---

### Q7: ORMを切り替えるにはどうする？

**A:** **このプロジェクトはORM切り替えに対応**しています。

現在の構成:
```
internal/adapter/gateway/db/
├── sqlc/              # sqlc実装（現在使用中）
│   ├── account_repository.go
│   ├── note_repository.go
│   └── template_repository.go
└── gorm/              # GORM実装（準備済み）
    └── account_repository.go
```

**切り替え手順（3ステップ）:**

1. **import文の切り替え**
   ```go
   // internal/driver/factory/repository_factory.go
   import (
       "github.com/jackc/pgx/v5/pgxpool"

       "immortal-architecture-cqrs/backend/internal/adapter/gateway/db/sqlc"
       // "immortal-architecture-cqrs/backend/internal/adapter/gateway/db/gorm"  ← コメント外す
       "immortal-architecture-cqrs/backend/internal/port"
   )
   ```

2. **Factory関数の修正**
   ```go
   func NewAccountRepoFactory(pool *pgxpool.Pool) func() port.AccountRepository {
       return func() port.AccountRepository {
           // Current: sqlc実装
           return sqlc.NewAccountRepository(pool)

           // To switch to GORM, replace above with:
           // return gorm.NewAccountRepository(db)
       }
   }
   ```

3. **DB接続の変更**（sqlc用のpoolからGORM用のdbへ）
   ```go
   // internal/driver/initializer/api/initializer.go
   // pool, err := driverdb.NewPool(ctx, cfg.DatabaseURL)  // sqlc用
   db, err := gorm.Open(postgres.Open(cfg.DatabaseURL))   // GORM用
   ```

**重要:**
- **Domain、UseCase、Controller、Presenterは変更不要**
- Gateway層だけ差し替えればOK
- これがクリーンアーキテクチャの変更可能性！

---

## ✅ チェックリスト: コードを書く前に

新しい機能を追加する前に、このチェックリストを確認しよう！

### ✅ Domain層

- [ ] エンティティや値オブジェクトに**フレームワークの型**を使っていないか？
- [ ] ビジネスルール（検証、状態遷移）を**ドメイン層**に書いているか？
- [ ] `import "github.com/..."` などの外部依存がないか？

### ✅ UseCase層

- [ ] Port（インターフェース）経由で依存しているか？
- [ ] 具体的な実装（`&db.Repository{}`）に依存していないか？
- [ ] トランザクション境界を適切に決めているか？
- [ ] 1メソッド = 1仕事になっているか？

### ✅ Controller層

- [ ] HTTPリクエスト → ドメインDTOに変換しているか？
- [ ] UseCaseを呼んで、Presenterから結果を取得しているか？
- [ ] ビジネスロジックを書いていないか？（UseCaseに書く）

### ✅ Presenter層

- [ ] ドメインモデル → OpenAPI型に変換しているか？
- [ ] 結果を内部で保存して、Controllerが取り出せるようにしているか？

### ✅ Gateway層

- [ ] DB行 → ドメインモデルに変換しているか？
- [ ] DBエラー → ドメインエラーに変換しているか？
- [ ] `pgtype.UUID` → `string` などの型変換をしているか？

### ✅ Port層

- [ ] InputPort, OutputPort, Repositoryのインターフェースを定義しているか？
- [ ] 入力DTOを定義しているか？

---

## 📁 付録: ディレクトリ構成の全体像

```
backend-clean/
├── cmd/
│   └── api/
│       └── main.go                      # エントリーポイント
│
├── internal/
│   ├── domain/                          # ❤️ ビジネスルール
│   │   ├── note/
│   │   │   ├── entity.go                # Note, Section
│   │   │   ├── types.go                 # NoteStatus
│   │   │   ├── logic.go                 # 検証ロジック
│   │   │   ├── aggregate.go             # WithMeta
│   │   │   └── *_test.go
│   │   ├── template/
│   │   ├── account/
│   │   ├── service/                     # ドメインサービス
│   │   │   ├── note_lifecycle.go        # BuildNote
│   │   │   ├── status_transition.go     # CanPublish
│   │   │   └── build_sections_from_template.go
│   │   └── errors/
│   │       └── errors.go                # ドメインエラー定義
│   │
│   ├── usecase/                         # 🎯 アプリケーションロジック
│   │   ├── note_interactor.go
│   │   ├── template_interactor.go
│   │   ├── account_interactor.go
│   │   └── mock/
│   │
│   ├── port/                            # 📝 インターフェース
│   │   ├── note_port.go
│   │   ├── template_port.go
│   │   ├── account_port.go
│   │   └── tx.go
│   │
│   ├── adapter/                         # 🔌 外部との接続
│   │   ├── http/
│   │   │   ├── controller/              # HTTPハンドラ
│   │   │   │   ├── note_controller.go
│   │   │   │   ├── template_controller.go
│   │   │   │   ├── account_controller.go
│   │   │   │   ├── server.go            # ルーティング
│   │   │   │   └── mock/
│   │   │   ├── presenter/               # レスポンス変換
│   │   │   │   ├── note_presenter.go
│   │   │   │   ├── template_presenter.go
│   │   │   │   └── account_presenter.go
│   │   │   └── generated/
│   │   │       └── openapi/             # OpenAPI生成物
│   │   │           └── server.gen.go
│   │   ├── grpc/
│   │   │   ├── controller/              # gRPCハンドラ
│   │   │   │   └── account_controller.go
│   │   │   ├── presenter/               # gRPCレスポンス変換
│   │   │   │   └── account_presenter.go
│   │   │   └── generated/
│   │   │       └── accountpb/           # protobuf生成物
│   │   └── gateway/
│   │       ├── db/                      # DB Repository
│   │       │   ├── sqlc/                # sqlc実装
│   │       │   │   ├── note_repository.go
│   │       │   │   ├── template_repository.go
│   │       │   │   ├── account_repository.go
│   │       │   │   ├── generated/       # sqlc生成物
│   │       │   │   ├── queries/         # SQLクエリ
│   │       │   │   └── mock/
│   │       │   └── gorm/                # GORM実装（準備済み）
│   │       │       └── account_repository.go
│   │       └── externalapi/             # 外部API (将来用)
│   │
│   └── driver/                          # 🔧 配線・初期化
│       ├── config/                      # 設定
│       ├── db/                          # DB接続
│       │   ├── pool.go
│       │   └── tx.go
│       ├── factory/                     # Factory関数
│       │   ├── usecase_factory.go
│       │   ├── repository_factory.go    # ORM切り替えポイント
│       │   ├── tx_factory.go
│       │   ├── http/                    # HTTP専用Factory
│       │   │   └── presenter_factory.go
│       │   └── grpc/                    # gRPC専用Factory
│       │       └── presenter_factory.go
│       └── initializer/
│           ├── api/
│           │   └── initializer.go       # HTTP API組み立て
│           └── grpc/
│               └── initializer.go       # gRPCサーバー組み立て
│
├── migrations/                          # DBマイグレーション
├── docs/                                # ドキュメント
└── tests/                               # E2Eテスト (将来用)
```

---

## 🚀 次のステップ

この処理の流れを理解したら、次は実際のコードを読んでみましょう：

1. **[01_why_clean_architecture.md](./01_why_clean_architecture.md)** - なぜ必要か理解する
2. **実際のコードを読む** - `internal/adapter/http/controller/note_controller.go` から始める
3. **[04_testing_guide.md](./04_testing_guide.md)** - テストの書き方を学ぶ
4. **デバッガーで追いかける** - 実際に動かして確認する

**Happy Learning!** 🎉
