# 🧪 テストの書き方 - backend-badとの比較

> 💡 **このドキュメントのゴール**
> backend-badとbackend-cleanのテストを比較して、
> 「あ、クリーンアーキテクチャだとテストがこんなに楽なんだ！」と腑に落ちてもらうためのドキュメントです。

---

## 🏚️ backend-bad のテスト = 地獄

### ❌ 問題1: DBがないとテストできない

```go
// ❌ backend-bad: ServiceがDBに直接依存
func (s *NoteService) Publish(noteID string) error {
    // sqlcを直接呼んでる
    note, err := s.queries.GetNote(ctx, noteID)
    // ...
}
```

**テストコード:**

```go
// ❌ テストするには本物のDBが必要
func TestNoteService_Publish(t *testing.T) {
    // 1. PostgreSQL起動（遅い）
    db := setupTestDB(t)
    defer db.Close()

    // 2. テストデータ投入（面倒）
    _, err := db.Exec("INSERT INTO notes ...")
    _, err = db.Exec("INSERT INTO templates ...")
    _, err = db.Exec("INSERT INTO accounts ...")

    // 3. Service作成
    service := &NoteService{queries: sqlc.New(db)}

    // 4. テスト実行（やっと！）
    err = service.Publish("123")

    // 👉 面倒、不安定
}
```

**🔥 何が問題？**

```
問題:
├─ PostgreSQL起動が必要（遅い）
├─ テストデータ投入が面倒（3テーブル、10行以上）
├─ 並列実行できない（データが競合する）
├─ CIで不安定（DB起動失敗、タイムアウト）
└─ 誰もテスト書かない（面倒すぎて）

👉 結果: テストカバレッジが低い、バグ多発
```

---

### ❌ 問題2: テストケースが増やせない

```go
// ❌ backend-bad: DB使うからパターン増やせない
func TestNoteService_Publish(t *testing.T) {
    db := setupTestDB(t)  // 遅い
    defer db.Close()

    // 正常系だけテスト
    // 異常系を追加したら時間がかかりすぎる...
    // 👉 誰も書かない
}
```

**🔥 何が問題？**

```
テストしたいケース:
  ✅ 正常系
  ❌ ノートが見つからない
  ❌ 既に公開済み
  ❌ 所有者が違う
  ❌ セクションが空

👉 DB準備が面倒で誰も書かない
👉 異常系がテストされない
👉 本番で異常系バグが発覚
```

---

### ❌ 問題3: モックが作れない

```go
// ❌ backend-bad: Interfaceがない
type NoteService struct {
    queries *sqldb.Queries  // ← 具体的な型に直接依存
}

// モックを作りたいけど...
// Interface がないから不可能！
```

---

## 🏡 backend-clean のテスト = 楽園

### ✅ 解決策1: モックでDB不要（爆速）

```go
// ✅ backend-clean: Interfaceに依存
type NoteInteractor struct {
    notes port.NoteRepository  // ← Interface
}

func (u *NoteInteractor) Publish(...) error {
    note, err := u.notes.Get(ctx, id)  // ← どの実装か知らない
    // ...
}
```

**テストコード:**

```go
// ✅ DBなし、モックだけでテスト
func TestNoteInteractor_Publish(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    // 1. モック作成
    mockRepo := mock.NewMockNoteRepository(ctrl)
    mockOutput := presenter.NewNotePresenter()

    // 2. モックの挙動を設定（簡単）
    mockRepo.EXPECT().
        Get(gomock.Any(), "123").
        Return(&note.WithMeta{
            Note: note.Note{
                ID:     "123",
                Status: note.StatusDraft,
            },
        }, nil)

    // 3. UseCase作成（モック注入）
    interactor := usecase.NewNoteInteractor(mockRepo, nil, nil, mockOutput)

    // 4. テスト実行
    err := interactor.Publish(ctx, "123", "owner-123")

    // 👉 簡単、安定
}
```

**🎉 何が良い？**

```
✅ DBなし（速い）
✅ テストデータ不要（モックで設定）
✅ 並列実行可能（メモリだけ）
✅ CIで安定（DB起動不要）
✅ みんなテスト書く（簡単だから）

👉 結果: テストカバレッジが高い、バグ激減
```

---

### ✅ 解決策2: テストケースが簡単に増やせる

```go
// ✅ backend-clean: テーブル駆動テスト
func TestNoteInteractor_Publish(t *testing.T) {
    tests := []struct {
        name      string
        noteID    string
        setupMock func(*mock.MockNoteRepository)
        wantErr   error
    }{
        {
            name:   "[Success] 公開成功",
            noteID: "123",
            setupMock: func(repo *mock.MockNoteRepository) {
                repo.EXPECT().
                    Get(gomock.Any(), "123").
                    Return(&note.WithMeta{
                        Note: note.Note{Status: note.StatusDraft},
                    }, nil)
            },
            wantErr: nil,
        },
        {
            name:   "[Fail] ノートが見つからない",
            noteID: "999",
            setupMock: func(repo *mock.MockNoteRepository) {
                repo.EXPECT().
                    Get(gomock.Any(), "999").
                    Return(nil, domainerr.ErrNotFound)
            },
            wantErr: domainerr.ErrNotFound,
        },
        {
            name:   "[Fail] 既に公開済み",
            noteID: "123",
            setupMock: func(repo *mock.MockNoteRepository) {
                repo.EXPECT().
                    Get(gomock.Any(), "123").
                    Return(&note.WithMeta{
                        Note: note.Note{Status: note.StatusPublish},  // 既に公開
                    }, nil)
            },
            wantErr: domainerr.ErrInvalidStatus,
        },
        // ケース追加が簡単！（3行足すだけ）
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // モック作成
            ctrl := gomock.NewController(t)
            defer ctrl.Finish()
            mockRepo := mock.NewMockNoteRepository(ctrl)
            tt.setupMock(mockRepo)

            // テスト実行
            interactor := usecase.NewNoteInteractor(mockRepo, nil, nil, nil)
            err := interactor.Publish(ctx, tt.noteID, "owner-123")

            // 検証
            if err != tt.wantErr {
                t.Errorf("want %v, got %v", tt.wantErr, err)
            }
        })
    }

    // 👉 3ケースを簡単にテスト
}
```

**🎉 何が良い？**

```
✅ ケース追加が簡単（3行足すだけ）
✅ 速い（DBなし）
✅ 正常系・異常系を全部テスト
✅ テーブルで一覧が見える

👉 backend-bad: DB準備が面倒で誰も書かない
👉 backend-clean: ケース追加が簡単で書きやすい
```

---

## 📊 比較: backend-bad vs backend-clean

### テストの書きやすさの比較

```
backend-bad:
┌─────────────────────────────┐
│ カバレッジ: 低い            │
│                             │
│ 理由:                       │
│ - テスト書くのが面倒        │
│ - DB起動が必要              │
│ - 遅い                      │
│ - 不安定（CI失敗）          │
│                             │
│ 👉 バグ多発                  │
└─────────────────────────────┘


backend-clean:
┌─────────────────────────────┐
│ カバレッジ: 高い            │
│                             │
│ 理由:                       │
│ - テスト書くのが簡単        │
│ - DB不要                    │
│ - 速い                      │
│ - 安定（CI成功）            │
│                             │
│ 👉 バグ激減                  │
└─────────────────────────────┘
```

---

## 💡 レイヤー別テストガイド

### 1️⃣ Domain層のテスト

**backend-bad:**

```go
// ❌ ドメイン層がない！
// ビジネスルールがServiceに埋まってる
func (s *NoteService) Publish(...) error {
    // 80行のメソッドの中にルールが点在
    if note.Status == "Publish" {
        return errors.New("already published")
    }
    // ...
    // 👉 ドメイン層のテストが書けない
    // 👉 Serviceごとテストする必要がある（DB必須）
}
```

**backend-clean:**

```go
// ✅ ビジネスルールが独立してる
// domain/service/status_transition.go
func CanPublish(note Note) error {
    if note.Status != StatusDraft {
        return domainerr.ErrInvalidStatus
    }
    if len(note.Sections) == 0 {
        return domainerr.ErrNoSections
    }
    return nil
}

// テストが簡単！
func TestCanPublish(t *testing.T) {
    tests := []struct {
        name    string
        note    note.Note
        wantErr error
    }{
        {
            name: "[Success] 下書き→公開",
            note: note.Note{
                Status:   note.StatusDraft,
                Sections: []note.Section{{Content: "内容"}},
            },
            wantErr: nil,
        },
        {
            name: "[Fail] 既に公開済み",
            note: note.Note{Status: note.StatusPublish},
            wantErr: domainerr.ErrInvalidStatus,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := service.CanPublish(tt.note)
            if err != tt.wantErr {
                t.Errorf("want %v, got %v", tt.wantErr, err)
            }
        })
    }
    // 👉 超簡単、速い
}
```

---

### 2️⃣ UseCase層のテスト

**backend-bad:**

```go
// ❌ DBに直接依存、テストできない
func (s *NoteService) Publish(noteID string) error {
    note, err := s.queries.GetNote(ctx, noteID)  // ← DB直呼び
    // ...
    // モック作れない
}
```

**backend-clean:**

```go
// ✅ Interfaceに依存、モックでテスト
func (u *NoteInteractor) Publish(...) error {
    note, err := u.notes.Get(ctx, id)  // ← Interface
    // ...
}

// テスト
func TestNoteInteractor_Publish(t *testing.T) {
    // モック作成
    mockRepo := mock.NewMockNoteRepository(ctrl)

    // 挙動設定（このIDならこれ返す）
    mockRepo.EXPECT().
        Get(gomock.Any(), "123").
        Return(&note.Note{...}, nil)

    // テスト実行
    interactor := usecase.NewNoteInteractor(mockRepo, ...)
    err := interactor.Publish(ctx, "123", "owner")

    // 👉 速い、DB不要
}
```

---

## 🎯 実際のテストコード例

### 例: ノート公開のテスト

**backend-bad（書けない）:**

```go
// ❌ テスト不可能
// - DBが必要
// - モックが作れない
// - 遅い
```

**backend-clean（簡単）:**

```go
func TestNoteInteractor_Publish(t *testing.T) {
    tests := []struct {
        name      string
        noteID    string
        ownerID   string
        setupMock func(*mock.MockNoteRepository)
        wantErr   error
    }{
        {
            name:    "[Success] 公開成功",
            noteID:  "123",
            ownerID: "owner-123",
            setupMock: func(repo *mock.MockNoteRepository) {
                // 1. Get が呼ばれる
                repo.EXPECT().
                    Get(gomock.Any(), "123").
                    Return(&note.WithMeta{
                        Note: note.Note{
                            ID:       "123",
                            OwnerID:  "owner-123",
                            Status:   note.StatusDraft,
                            Sections: []note.Section{{Content: "内容"}},
                        },
                    }, nil)

                // 2. Update が呼ばれる
                repo.EXPECT().
                    Update(gomock.Any(), gomock.Any()).
                    Return(&note.Note{Status: note.StatusPublish}, nil)
            },
            wantErr: nil,
        },
        {
            name:    "[Fail] ノートが見つからない",
            noteID:  "999",
            ownerID: "owner-123",
            setupMock: func(repo *mock.MockNoteRepository) {
                repo.EXPECT().
                    Get(gomock.Any(), "999").
                    Return(nil, domainerr.ErrNotFound)
            },
            wantErr: domainerr.ErrNotFound,
        },
        {
            name:    "[Fail] 所有者が違う",
            noteID:  "123",
            ownerID: "other-999",
            setupMock: func(repo *mock.MockNoteRepository) {
                repo.EXPECT().
                    Get(gomock.Any(), "123").
                    Return(&note.WithMeta{
                        Note: note.Note{
                            ID:      "123",
                            OwnerID: "owner-123",  // 違う所有者
                        },
                    }, nil)
            },
            wantErr: domainerr.ErrUnauthorized,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // モック作成
            ctrl := gomock.NewController(t)
            defer ctrl.Finish()
            mockRepo := mock.NewMockNoteRepository(ctrl)
            mockOutput := presenter.NewNotePresenter()

            // モックの挙動を設定
            tt.setupMock(mockRepo)

            // UseCase作成
            interactor := usecase.NewNoteInteractor(mockRepo, nil, nil, mockOutput)

            // テスト実行
            err := interactor.Publish(context.Background(), tt.noteID, tt.ownerID)

            // 検証
            if err != tt.wantErr {
                t.Errorf("want %v, got %v", tt.wantErr, err)
            }
        })
    }
}
```

---

## 🚀 テストの実行

### backend-bad

```bash
# 1. PostgreSQL起動
docker-compose up -d postgres

# 2. マイグレーション実行
make migrate

# 3. テスト実行（遅い）
go test ./...
# 👉 遅い
```

### backend-clean

```bash
# テスト実行（速い）
go test ./...
# 👉 速い

# DB不要！
# マイグレーション不要！
```

---

## 💡 まとめ

### backend-bad の問題

```
❌ DBがないとテストできない
   → PostgreSQL起動、データ投入が面倒

❌ テストが遅い
   → DB起動、マイグレーション実行が必要

❌ テストケースが増やせない
   → 遅すぎて誰も書かない

❌ モックが作れない
   → Interfaceがない

結果:
  👉 テストカバレッジが低い
  👉 バグ多発
  👉 本番で障害
```

### backend-clean の解決策

```
✅ モックでDB不要
   → メモリだけでテスト（速い）

✅ テストが速い
   → DB不要、マイグレーション不要

✅ テストケースが簡単に増やせる
   → テーブル駆動テストで3行足すだけ

✅ モックが簡単に作れる
   → Interfaceがあるから

結果:
  👉 テストカバレッジが高い
  👉 バグ激減
  👉 安心して変更できる
```

---

## 🎯 次のステップ

1. **実際のテストコードを読む**
   - `internal/domain/service/status_transition_test.go`
   - `internal/usecase/note_interactor_test.go`

2. **簡単なテストを書いてみる**
   - Domainのバリデーションから始める
   - テーブル駆動テストで書く

3. **モックを使ってみる**
   - UseCaseのテストを書く
   - gomockで偽物を作る

**Happy Testing!** 🎉
