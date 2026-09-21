# 🏗️ CQRSアーキテクチャガイド - 読み書き分離の実装

> 💡 **このドキュメントのゴール**
> Clean Architecture（CRUD）を **CQRSにどう組み替えたか** が
> 具体的にわかること。実際のコードを見ながら読んでね。

---

## 🔄 Before / After を一枚絵で見よう

### Before: Clean Architecture（CRUD）

```
📱 画面
 │
 ▼
┌─────────────────────────────────────────────────┐
│ Controller（HTTP）                                │
│   NoteController                                 │
│     .List()   ← 読む                             │
│     .GetByID()← 読む                             │
│     .Create() ← 書く    ← 全部同じControllerに    │
│     .Update() ← 書く      まとまってる           │
│     .Delete() ← 書く                             │
└──────────┬──────────────────────────────────────┘
           │
           ▼
┌─────────────────────────────────────────────────┐
│ UseCase（NoteInteractor）                        │
│   .List()    ← 読む                              │
│   .Get()     ← 読む                              │
│   .Create()  ← 書く    ← 全部同じInteractorに     │
│   .Update()  ← 書く      まとまってる            │
│   .Delete()  ← 書く                              │
└──────────┬──────────────────────────────────────┘
           │
           ▼
┌─────────────────────────────────────────────────┐
│ Repository（NoteRepository）                     │
│   .List()    ← 読む（3テーブルJOIN）              │
│   .Get()     ← 読む（3テーブルJOIN）              │
│   .Create()  ← 書く                              │
│   .Update()  ← 書く    ← 全部同じRepositoryに     │
│   .Delete()  ← 書く      まとまってる            │
└──────────┬──────────────────────────────────────┘
           │
           ▼
┌─────────────────────────────────────────────────┐
│ PostgreSQL                                       │
│   notes / templates / accounts / sections / fields│
│   ← 読むのも書くのも同じテーブル                   │
└─────────────────────────────────────────────────┘
```

**問題：全部が1本の道を通ってる。読むも書くもごちゃ混ぜ。**

---

### After: CQRS（読み書き分離）

```
📱 画面
 │
 ├──── 読む（Query）──────────────── 書く（Command）────┐
 ▼                                    ▼                │
┌──────────────────────┐  ┌────────────────────────┐   │
│ Controller           │  │ Controller             │   │
│  .List()             │  │  .Create()             │   │
│  .GetByID()          │  │  .Update()             │   │
│                      │  │  .Delete()             │   │
│  読む専用！           │  │  .Publish()            │   │
└────────┬─────────────┘  └────────┬───────────────┘   │
         │                         │                    │
         ▼                         ▼                    │
┌──────────────────────┐  ┌────────────────────────┐   │
│ NoteQueryInteractor  │  │ NoteCommandInteractor  │   │
│  .List()             │  │  .Create()             │   │
│  .Get()              │  │  .Update()             │   │
│                      │  │  .Delete()             │   │
│  バリデーションなし    │  │  .ChangeStatus()       │   │
│  ただ取るだけ         │  │                        │   │
│                      │  │  バリデーションあり      │   │
│                      │  │  ドメインロジック実行    │   │
│                      │  │  + Read Model同期      │   │
└────────┬─────────────┘  └───┬────────┬───────────┘   │
         │                    │        │                │
         ▼                    │        ▼                │
┌──────────────────────┐  │  ┌────────────────────┐    │
│ NoteReadModel        │  │  │ NoteRepository     │    │
│ Repository           │  │  │（既存のまま）        │    │
│  .List()             │  │  │  .Create()         │    │
│  .Get()              │  │  │  .Update()         │    │
│  .Upsert() ←同期用   │◀─┘  │  .Delete()         │    │
│  .Delete() ←同期用   │     │  .ReplaceSections() │    │
└────────┬─────────────┘     └────────┬────────────┘   │
         │                            │                 │
         ▼                            ▼                 │
┌──────────────────────┐  ┌─────────────────────────┐  │
│ note_read_models     │  │ notes / templates /      │  │
│（キャッシュテーブル）  │  │ accounts / sections /    │  │
│                      │  │ fields                   │  │
│ JOINの結果が         │  │（元のテーブル = マスター） │  │
│ 事前に入ってる！      │  │                          │  │
└──────────────────────┘  └─────────────────────────┘  │
         PostgreSQL（同じDB内）                           │
```

**ポイント：Repositoryは分離しすぎない。UseCaseレベルで分けるのがCQRSの基本。**
- `NoteRepository` は既存のまま（Command側が使う）
- `NoteReadModelRepository` だけ新規追加（Read Model専用）

---

## 🧩 ファイル構成 - 何がどこにある？

### CQRS後のファイル構成（実際のコード）

```
internal/
├── port/
│   ├── note_port.go               ← NoteRepository（既存のまま）
│   │                                 + Input型（NoteCreateInput等）
│   ├── note_command_port.go       ← NoteCommandInputPort
│   │                                 NoteCommandOutputPort
│   └── note_read_model_port.go    ← NoteQueryInputPort
│                                     NoteQueryOutputPort
│                                     NoteReadModelRepository
│
├── usecase/
│   ├── note_command_interactor.go ← Command（書く + Read Model同期）
│   ├── note_query_interactor.go   ← Query（キャッシュテーブルから取るだけ）
│   └── note_helpers.go            ← 共通ヘルパー（buildSections等）
│
├── adapter/
│   ├── http/
│   │   ├── controller/
│   │   │   └── note_controller.go ← Command/Queryを使い分ける
│   │   └── presenter/
│   │       ├── note_command_presenter.go ← Command用プレゼンター
│   │       ├── note_query_presenter.go   ← Query用プレゼンター
│   │       └── note_helpers.go           ← 共通ヘルパー（toNoteResponse）
│   └── gateway/db/sqlc/
│       ├── note_repository.go            ← NoteRepository（既存のまま）
│       ├── note_read_model_repository.go ← NoteReadModelRepository（新規）
│       └── queries/
│           ├── notes.sql                 ← 既存SQL
│           └── note_read_models.sql      ← Read Model用SQL（新規）
│
├── domain/
│   └── note/
│       └── read_model.go          ← ReadModel, SectionReadModel（新規）
│
└── driver/
    ├── factory/
    │   ├── repository_factory.go  ← NoteRepoFactory + NoteReadModelRepoFactory
    │   ├── usecase_factory.go     ← NoteCommandInputFactory + NoteQueryInputFactory
    │   └── http/
    │       └── presenter_factory.go ← NoteCommandOutputFactory + NoteQueryOutputFactory
    └── initializer/
        └── api/initializer.go     ← 全部ここで配線

migrations/
├── 20250210000000_add_note_read_models.up.sql   ← キャッシュテーブル作成
└── 20250210000001_seed_note_read_models.up.sql  ← 既存データの初回同期
```

---

## 📐 各レイヤーの実装を見ていこう

### 1️⃣ ドメイン層 - Read Modelの型を追加

> 📂 `internal/domain/note/read_model.go`

```go
// ReadModel はキャッシュテーブル note_read_models に対応する構造体。
// JOINの結果を事前にまとめたもの。
type ReadModel struct {
    ID             string
    Title          string
    Status         NoteStatus
    TemplateID     string
    TemplateName   string         // ← templates テーブルの情報が入ってる！
    OwnerID        string
    OwnerFirstName string         // ← accounts テーブルの情報が入ってる！
    OwnerLastName  string
    OwnerThumbnail *string
    Sections       []SectionReadModel  // ← sections + fields の情報が入ってる！
    CreatedAt      time.Time
    UpdatedAt      time.Time
}

type SectionReadModel struct {
    ID         string `json:"id"`
    FieldID    string `json:"field_id"`
    FieldLabel string `json:"field_label"`
    FieldOrder int    `json:"field_order"`
    IsRequired bool   `json:"is_required"`
    Content    string `json:"content"`
}
```

**ポイント：JOINしないと取れなかった情報が、1つの構造体にぜんぶ入ってる。**

ドメイン層で変更したのは **このファイル1つだけ**。
`entity.go`、`logic.go`、`aggregate.go` は一切触ってない。

---

### 2️⃣ ポート層 - CommandとQueryを分ける

#### Command側のポート

> 📂 `internal/port/note_command_port.go`

```go
// NoteCommandInputPort は書く専用のユースケース。
type NoteCommandInputPort interface {
    Create(ctx context.Context, input NoteCreateInput) error
    Update(ctx context.Context, input NoteUpdateInput) error
    ChangeStatus(ctx context.Context, input NoteStatusChangeInput) error
    Delete(ctx context.Context, id, ownerID string) error
}

// NoteCommandOutputPort は書く専用のプレゼンター。
type NoteCommandOutputPort interface {
    PresentNote(ctx context.Context, note *note.WithMeta) error
    PresentNoteDeleted(ctx context.Context) error
}
```

#### Query側 + Read Modelのポート

> 📂 `internal/port/note_read_model_port.go`

```go
// NoteQueryInputPort は読む専用のユースケース。
type NoteQueryInputPort interface {
    List(ctx context.Context, filters note.Filters) error
    Get(ctx context.Context, id string) error
}

// NoteQueryOutputPort は読む専用のプレゼンター。
type NoteQueryOutputPort interface {
    PresentNoteList(ctx context.Context, notes []note.ReadModel) error
    PresentNote(ctx context.Context, note *note.ReadModel) error
}

// NoteReadModelRepository はRead Modelテーブルの操作全般。
// 読み取り（List/Get）と同期（Upsert/Delete）の両方を持つ。
type NoteReadModelRepository interface {
    List(ctx context.Context, filters note.Filters) ([]note.ReadModel, error)
    Get(ctx context.Context, id string) (*note.ReadModel, error)
    Upsert(ctx context.Context, model note.ReadModel) error  // ← Command側から呼ばれる
    Delete(ctx context.Context, id string) error              // ← Command側から呼ばれる
}
```

#### 既存のNoteRepository（変更なし）

> 📂 `internal/port/note_port.go`

```go
// NoteRepository は既存のリポジトリ。変更なし。
// Command側がそのまま使う。
type NoteRepository interface {
    List(ctx context.Context, filters note.Filters) ([]note.WithMeta, error)
    Get(ctx context.Context, id string) (*note.WithMeta, error)
    Create(ctx context.Context, n note.Note) (*note.Note, error)
    Update(ctx context.Context, n note.Note) (*note.Note, error)
    UpdateStatus(ctx context.Context, id string, status note.NoteStatus) (*note.Note, error)
    Delete(ctx context.Context, id string) error
    ReplaceSections(ctx context.Context, noteID string, sections []note.Section) error
}
```

**Repositoryを過剰に分離しないのがポイント。** UseCaseレベルで分離すればCQRSは成立する。

---

### 3️⃣ ユースケース層 - QueryとCommandの中身

#### Query側 — めっちゃシンプル

> 📂 `internal/usecase/note_query_interactor.go`

```go
type NoteQueryInteractor struct {
    readModelRepo port.NoteReadModelRepository
    output        port.NoteQueryOutputPort
}

// List はキャッシュテーブルから一覧を取得する。JOINなし！
func (u *NoteQueryInteractor) List(ctx context.Context, filters note.Filters) error {
    notes, err := u.readModelRepo.List(ctx, filters)
    if err != nil {
        return err
    }
    return u.output.PresentNoteList(ctx, notes)
}
```

**バリデーションもドメインロジックもない。ただ取って返すだけ。これがQuery側の強み。**

#### Command側 — 書く＋Read Model同期

> 📂 `internal/usecase/note_command_interactor.go`

```go
type NoteCommandInteractor struct {
    notes         port.NoteRepository          // ← 既存のRepository
    readModelRepo port.NoteReadModelRepository // ← Read Model同期用
    templates     port.TemplateRepository
    tx            port.TxManager
    output        port.NoteCommandOutputPort
}

// Create はノートを作成し、Read Modelも同期する。
func (u *NoteCommandInteractor) Create(ctx context.Context, input port.NoteCreateInput) error {
    // ... バリデーション ...

    err = u.tx.WithinTransaction(ctx, func(txCtx context.Context) error {
        // 1. 元のテーブルに書く（マスターデータ）
        nn, err := u.notes.Create(txCtx, newNote)
        // ... sections も作成 ...

        // 2. Read Modelも同期する（キャッシュ）
        created, err := u.notes.Get(txCtx, noteID)
        return u.readModelRepo.Upsert(txCtx, toReadModel(created))
    })
    // ...
}
```

**ポイント：同じトランザクション内で元テーブルもキャッシュも更新する。**

```
トランザクション開始
  │
  ├─ notes テーブルに書く            ← マスターデータ
  ├─ sections テーブルに書く
  ├─ note_read_models に同期         ← キャッシュ
  │
  └─ 全部成功 → COMMIT
     1つでも失敗 → ROLLBACK（全部元に戻る）
```

---

### 4️⃣ アダプター層 - キャッシュテーブルとリポジトリ

#### キャッシュテーブル

> 📂 `migrations/20250210000000_add_note_read_models.up.sql`

```sql
CREATE TABLE note_read_models (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('Draft', 'Publish')),
    template_id UUID NOT NULL,
    template_name TEXT NOT NULL,         -- ← JOINしなくていい！
    owner_id UUID NOT NULL,
    owner_first_name TEXT NOT NULL,      -- ← JOINしなくていい！
    owner_last_name TEXT NOT NULL,
    owner_thumbnail TEXT,
    sections_json JSONB NOT NULL DEFAULT '[]',  -- ← 全セクション入り！
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);
```

#### Read Model用SQLクエリ

> 📂 `internal/adapter/gateway/db/sqlc/queries/note_read_models.sql`

```sql
-- JOINなし！1テーブルから取るだけ！
-- name: ListNoteReadModels :many
SELECT * FROM note_read_models
WHERE (NULLIF($1::text, '') IS NULL OR status = $1)
  AND ($2::uuid IS NULL OR template_id = $2)
  AND ($3::uuid IS NULL OR owner_id = $3)
  AND (NULLIF($4::text, '') IS NULL OR title ILIKE '%' || $4 || '%')
ORDER BY updated_at DESC;

-- name: UpsertNoteReadModel :exec
INSERT INTO note_read_models (...) VALUES (...)
ON CONFLICT (id) DO UPDATE SET ...;  -- ← 同じIDがあれば更新
```

**Before/After比較：**

```sql
-- 【Before】 3テーブルJOIN + セクション別クエリ
SELECT n.*, t.name, a.first_name, a.last_name, a.thumbnail
FROM notes n
JOIN templates t ON t.id = n.template_id
JOIN accounts a ON a.id = n.owner_id

-- 【After】 JOINなし！
SELECT * FROM note_read_models WHERE ...
```

#### Read Modelリポジトリ

> 📂 `internal/adapter/gateway/db/sqlc/note_read_model_repository.go`

```go
type NoteReadModelRepository struct {
    pool    *pgxpool.Pool
    queries *generated.Queries
}

// List はキャッシュテーブルから取得。JOINなし。
func (r *NoteReadModelRepository) List(ctx context.Context, filters note.Filters) ([]note.ReadModel, error) {
    rows, err := queriesForContext(ctx, r.queries).ListNoteReadModels(ctx, params)
    // ... rows を ReadModel に変換して返す
}

// Upsert はキャッシュテーブルを更新（Command側から呼ばれる）。
func (r *NoteReadModelRepository) Upsert(ctx context.Context, model note.ReadModel) error {
    sectionsJSON, _ := json.Marshal(model.Sections)
    return queriesForContext(ctx, r.queries).UpsertNoteReadModel(ctx, ...)
}
```

---

### 5️⃣ コントローラー層 - CommandとQueryの使い分け

> 📂 `internal/adapter/http/controller/note_controller.go`

```go
type NoteController struct {
    // Command用（書く）
    commandInputFactory  func(noteRepo port.NoteRepository, readModelRepo port.NoteReadModelRepository, ...) port.NoteCommandInputPort
    commandOutputFactory func() *presenter.NoteCommandPresenter

    // Query用（読む）
    queryInputFactory  func(readModelRepo port.NoteReadModelRepository, ...) port.NoteQueryInputPort
    queryOutputFactory func() *presenter.NoteQueryPresenter

    // 共通
    noteRepoFactory      func() port.NoteRepository
    readModelRepoFactory func() port.NoteReadModelRepository
    tplRepoFactory       func() port.TemplateRepository
    txFactory            func() port.TxManager
}

// List は Query側を使う（読む）
func (c *NoteController) List(ctx echo.Context, params ...) error {
    queryInput, p := c.newQueryIO()  // ← Query用のUseCaseを生成
    // ...
}

// Create は Command側を使う（書く）
func (c *NoteController) Create(ctx echo.Context) error {
    commandInput, p := c.newCommandIO()  // ← Command用のUseCaseを生成
    // ...
}
```

**どのエンドポイントがどっちを使うか：**

| エンドポイント | Query / Command | メソッド |
|--------------|----------------|---------|
| `GET /notes` | **Query** | `newQueryIO()` |
| `GET /notes/:id` | **Query** | `newQueryIO()` |
| `POST /notes` | **Command** | `newCommandIO()` |
| `PUT /notes/:id` | **Command** | `newCommandIO()` |
| `DELETE /notes/:id` | **Command** | `newCommandIO()` |
| `POST /notes/:id/publish` | **Command** | `newCommandIO()` |
| `POST /notes/:id/unpublish` | **Command** | `newCommandIO()` |

---

### 6️⃣ ファクトリー＆初期化 - 配線

> 📂 `internal/driver/initializer/api/initializer.go`

```go
// リポジトリ
noteRepoFactory := factory.NewNoteRepoFactory(pool)              // 既存
noteReadModelRepoFactory := factory.NewNoteReadModelRepoFactory(pool)  // 新規

// プレゼンター
noteCommandOutputFactory := httpfactory.NewNoteCommandOutputFactory()
noteQueryOutputFactory := httpfactory.NewNoteQueryOutputFactory()

// ユースケース
noteCommandInputFactory := factory.NewNoteCommandInputFactory()
noteQueryInputFactory := factory.NewNoteQueryInputFactory()

// コントローラーに全部渡す
nc := httpcontroller.NewNoteController(
    noteCommandInputFactory, noteCommandOutputFactory,
    noteQueryInputFactory, noteQueryOutputFactory,
    noteRepoFactory, noteReadModelRepoFactory,
    templateRepoFactory, txFactory,
)
```

**将来Redisに移行したいとき：**

```go
// ここを変えるだけ！UseCase層は一切変更なし！
func NewNoteReadModelRepoFactory(...) func() port.NoteReadModelRepository {
    return func() port.NoteReadModelRepository {
        return redis.NewNoteReadModelRepository(client)  // ← ここだけ差し替え
    }
}
```

---

## 🔄 データの流れを追ってみよう

### 📖 読む場合（GET /notes）

```
1. 📱 画面: GET /notes?status=Publish

2. 🎮 NoteController: List()
   → queryInput, p := c.newQueryIO()

3. 📋 NoteQueryInteractor: List()
   → readModelRepo.List(ctx, filters)

4. 💾 NoteReadModelRepository: List()
   → SELECT * FROM note_read_models WHERE ...
   → JOINなし！1テーブルから一発取得！

5. 🎨 NoteQueryPresenter: PresentNoteList()
   → ReadModel → APIレスポンスに変換

6. 📱 画面: ノート一覧が表示される（速い！）
```

### ✏️ 書く場合（POST /notes）

```
1. 📱 画面: POST /notes {title: "新しいノート", ...}

2. 🎮 NoteController: Create()
   → commandInput, p := c.newCommandIO()

3. 📋 NoteCommandInteractor: Create()
   → バリデーション実行
   → tx.WithinTransaction 開始

4.    💾 NoteRepository（既存）:
      → notes テーブルに INSERT（マスターデータ）
      → sections テーブルに INSERT

5.    💾 NoteReadModelRepository:
      → note_read_models テーブルに UPSERT（キャッシュ同期）

      → トランザクション COMMIT
        （4と5が両方成功して初めて確定）

6. 🎨 NoteCommandPresenter: PresentNote()
   → ドメインモデル → APIレスポンスに変換

7. 📱 画面: 作成されたノートが表示される
```

---

## 🧱 変更しなかったところ（超重要）

```
┌─────────────────────────────────────────────────┐
│  ✅ 変更なし                                     │
│                                                  │
│  domain/note/entity.go      ← Note, Section     │
│  domain/note/logic.go       ← バリデーション     │
│  domain/note/aggregate.go   ← 集約操作           │
│  domain/note/types.go       ← Filters, WithMeta  │
│  domain/errors/             ← エラー定義         │
│  domain/service/            ← ドメインサービス    │
│  domain/template/           ← テンプレート全般    │
│  domain/account/            ← アカウント全般      │
│                                                  │
│  port/note_port.go          ← NoteRepository     │
│  port/account_port.go       ← アカウント系       │
│  port/template_port.go      ← テンプレート系     │
│  port/tx_port.go            ← トランザクション    │
│                                                  │
│  adapter/gateway/db/sqlc/                        │
│    note_repository.go       ← 既存Repository     │
│                                                  │
│  adapter/http/generated/    ← OpenAPI生成コード   │
│  adapter/grpc/              ← gRPC全般           │
│  driver/config/             ← 設定               │
│  driver/db/                 ← DB接続/TxManager   │
│                                                  │
│  → ドメイン層はまったく触ってない！               │
│  → 既存のNoteRepositoryもそのまま！              │
│  → Account, Template も変更なし！                │
└─────────────────────────────────────────────────┘
```

---

## 🏫 学校で言うと

```
【Before: CRUD】

  先生（UseCase）が1人で
  「テストを作る」「テストを配る」「テストを採点する」を全部やってた

【After: CQRS】

  「テストを作る先生」と「テストを配る先生」に分けた

  作る先生（Command）:
    → テスト問題を作る
    → 正式な記録に残す（notes テーブル = マスター）
    → ついでに配布用のコピーも作る（note_read_models = キャッシュ）

  配る先生（Query）:
    → 事前にコピーされた配布物を配るだけ
    → 超速い！
    → 作る先生の仕事を邪魔しない！
```

---

## 🎯 まとめ

| 質問 | 答え |
|-----|------|
| 何を分けた？ | UseCaseとPresenterを Command / Query に分離 |
| Repositoryは？ | 既存の `NoteRepository` はそのまま。`NoteReadModelRepository` だけ追加 |
| ドメイン層は変わった？ | `read_model.go` を1つ追加しただけ |
| 同期はどうやる？ | Command側で元テーブルに書いた後、同じトランザクション内でキャッシュも更新 |
| 既存データは？ | 初回マイグレーションでキャッシュテーブルに流し込み |
| 全員のデータが同期される？ | はい。誰がAPIを叩いても、Command経由で必ずキャッシュが更新される |
| Clean Architectureとの関係は？ | ポート（インターフェース）があるから、Read Modelの実装を自由に差し替えられる |
