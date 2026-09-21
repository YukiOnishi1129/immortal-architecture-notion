# 🤖 AI駆動開発ガイド - このドキュメントを使ってAIに実装させよう

> 💡 **このドキュメントのゴール**
> このドキュメント群をAIに読み込ませて、
> 「タグ機能を追加して」と言うだけで、
> Clean Architectureに従った実装+テストコードができることを体験する

---

## 🎯 なぜAI駆動開発なのか？

```
従来の開発:
├─ 設計書を書く（1時間）
├─ 実装する（3時間）
├─ テスト書く（2時間）
└─ レビュー修正（1時間）
  👉 合計: 7時間

AI駆動開発（Clean Architecture）:
├─ AIにドキュメント読ませる（1分）
├─ 「タグ機能追加して」と指示（10秒）
├─ AIが実装+テスト生成（1分）
└─ レビュー修正（30分）
  👉 合計: 32分（13倍速い！）
```

**なぜClean Architectureだと速いの？**
- 構造が明確 → AIが「どこに書くべきか」を理解しやすい
- パターンが統一 → AIが「同じパターンで書けばいい」と判断しやすい
- 責務が分離 → AIが「この機能はこの層」と正確に配置

---

## 📖 ステップ1: AIにドキュメントを読み込ませる

### 1-1. 必要なドキュメントをAIに渡す

```
以下のドキュメントを読んでください。
このプロジェクトのClean Architecture設計方針を理解してください。

- backend-clean/docs/01_why_clean_architecture.md
- backend-clean/docs/02_clean_architecture_guide.md
- backend-clean/docs/03_testing_guide.md
```

### 1-2. ドメイン設計・テーブル設計・API仕様を読ませる

```
以下のドキュメントを読んで、システム全体の設計を理解してください。

【全体設計ドキュメント】
- docs/global_design/05_domain_design.md
  → ドメインモデルの全体像
  → エンティティの関係性
  → ビジネスルール

- docs/global_design/06_database_design.md
  → テーブル設計
  → ER図
  → カラム定義

- docs/global_design/07_api_design.md
  → API設計方針
  → エンドポイント一覧

【TypeSpec定義（詳細）】
- api-schema/typespec/models/note.tsp
  → Note, Section, NoteStatus などのドメインエンティティ
  → CreateNoteRequest, UpdateNoteRequest などのリクエスト型
  → NoteResponse などのレスポンス型
  → バリデーションルール（@minLength, @maxLength など）

- api-schema/typespec/models/template.tsp
  → Template, Field などのエンティティ

- api-schema/typespec/models/account.tsp
  → Account, AccountSummary などのエンティティ

- api-schema/typespec/routes/*.tsp
  → GET /api/notes, POST /api/notes などのエンドポイント
  → リクエスト/レスポンスの型マッピング

これらを読むことで、AIは以下を理解します:
✅ ドメインモデルの全体像と関係性
✅ テーブル設計とカラム定義
✅ ビジネスルール（バリデーション、ステータス遷移）
✅ API仕様（どんなリクエスト/レスポンスか）
```

### 1-3. 既存コードの実装パターンを理解させる

```
以下の既存コードを読んで、実装パターンを理解してください。

【Domain層】
- backend-clean/internal/domain/note/note.go
- backend-clean/internal/domain/note/status.go
- backend-clean/internal/domain/service/status_transition.go

【Port層】
- backend-clean/internal/port/note_repository.go

【UseCase層】
- backend-clean/internal/usecase/note_interactor.go

【Adapter層】
- backend-clean/internal/adapter/gateway/note_gateway.go
- backend-clean/internal/adapter/controller/note_controller.go
```

---

## 💬 ステップ2: AIに機能追加を指示する

### 例: タグ機能を追加

```
このプロジェクトに「タグ機能」を追加してください。

【要件】
- ノートに複数のタグを付けられる
- タグの追加・削除ができる
- タグで検索できる

【指示】
1. Clean Architectureの設計方針に従って実装してください
2. 既存のNoteの実装パターンを参考にしてください
3. 以下の層すべてに必要なコードを生成してください:
   - Domain層（エンティティ、バリデーション）
   - Port層（Interface）
   - UseCase層（ビジネスロジック）
   - Adapter層（Gateway、Controller）
4. テストコードも生成してください（C1カバレッジ達成レベル）
5. ファイルパスも明記してください
```

---

## ✅ ステップ3: AIが生成するコード

AIが以下のような構造でコードを生成します:

### 3-1. Domain層（エンティティとロジック）

**ファイル: `internal/domain/tag/tag.go`**

```go
package tag

import "github.com/yourusername/backend-clean/internal/domain/domainerr"

// Tag represents a tag entity
type Tag struct {
    ID   string
    Name string
}

// Validate validates tag
func (t Tag) Validate() error {
    if t.Name == "" {
        return domainerr.ErrInvalidInput
    }
    if len(t.Name) > 50 {
        return domainerr.ErrInvalidInput
    }
    return nil
}
```

**テストファイル: `internal/domain/tag/tag_test.go`**

```go
package tag_test

import (
    "testing"
    "github.com/yourusername/backend-clean/internal/domain/tag"
    "github.com/yourusername/backend-clean/internal/domain/domainerr"
)

func TestTag_Validate(t *testing.T) {
    tests := []struct {
        name    string
        tag     tag.Tag
        wantErr error
    }{
        {
            name:    "[Success] 有効なタグ",
            tag:     tag.Tag{ID: "1", Name: "Go"},
            wantErr: nil,
        },
        {
            name:    "[Fail] 名前が空",
            tag:     tag.Tag{ID: "1", Name: ""},
            wantErr: domainerr.ErrInvalidInput,
        },
        {
            name:    "[Fail] 名前が長すぎる",
            tag:     tag.Tag{ID: "1", Name: string(make([]byte, 51))},
            wantErr: domainerr.ErrInvalidInput,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := tt.tag.Validate()
            if err != tt.wantErr {
                t.Errorf("want %v, got %v", tt.wantErr, err)
            }
        })
    }
}
```

---

### 3-2. Port層（Interface）

**ファイル: `internal/port/tag_repository.go`**

```go
package port

import (
    "context"
    "github.com/yourusername/backend-clean/internal/domain/tag"
)

type TagRepository interface {
    Create(ctx context.Context, t tag.Tag) error
    FindByNoteID(ctx context.Context, noteID string) ([]tag.Tag, error)
    Delete(ctx context.Context, noteID, tagID string) error
}
```

---

### 3-3. UseCase層（ビジネスロジック）

**ファイル: `internal/usecase/tag_interactor.go`**

```go
package usecase

import (
    "context"
    "github.com/yourusername/backend-clean/internal/domain/tag"
    "github.com/yourusername/backend-clean/internal/port"
)

type TagInteractor struct {
    tags port.TagRepository
}

func NewTagInteractor(tags port.TagRepository) *TagInteractor {
    return &TagInteractor{tags: tags}
}

func (u *TagInteractor) AddTag(ctx context.Context, noteID, tagName string) error {
    t := tag.Tag{
        Name: tagName,
    }

    // バリデーション
    if err := t.Validate(); err != nil {
        return err
    }

    // 保存
    return u.tags.Create(ctx, t)
}

func (u *TagInteractor) GetTags(ctx context.Context, noteID string) ([]tag.Tag, error) {
    return u.tags.FindByNoteID(ctx, noteID)
}
```

**テストファイル: `internal/usecase/tag_interactor_test.go`**

```go
package usecase_test

import (
    "context"
    "testing"
    "github.com/golang/mock/gomock"
    "github.com/yourusername/backend-clean/internal/domain/tag"
    "github.com/yourusername/backend-clean/internal/domain/domainerr"
    "github.com/yourusername/backend-clean/internal/usecase"
    "github.com/yourusername/backend-clean/internal/port/mock"
)

func TestTagInteractor_AddTag(t *testing.T) {
    tests := []struct {
        name      string
        noteID    string
        tagName   string
        setupMock func(*mock.MockTagRepository)
        wantErr   error
    }{
        {
            name:    "[Success] タグ追加成功",
            noteID:  "note-123",
            tagName: "Go",
            setupMock: func(repo *mock.MockTagRepository) {
                repo.EXPECT().
                    Create(gomock.Any(), gomock.Any()).
                    Return(nil)
            },
            wantErr: nil,
        },
        {
            name:    "[Fail] タグ名が空",
            noteID:  "note-123",
            tagName: "",
            setupMock: func(repo *mock.MockTagRepository) {
                // バリデーションで失敗するのでCreateは呼ばれない
            },
            wantErr: domainerr.ErrInvalidInput,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            ctrl := gomock.NewController(t)
            defer ctrl.Finish()
            mockRepo := mock.NewMockTagRepository(ctrl)
            tt.setupMock(mockRepo)

            interactor := usecase.NewTagInteractor(mockRepo)
            err := interactor.AddTag(context.Background(), tt.noteID, tt.tagName)

            if err != tt.wantErr {
                t.Errorf("want %v, got %v", tt.wantErr, err)
            }
        })
    }
}
```

---

### 3-4. Adapter層（Gateway、Controller）

**ファイル: `internal/adapter/gateway/tag_gateway.go`**

```go
package gateway

import (
    "context"
    "github.com/yourusername/backend-clean/internal/domain/tag"
)

type TagGateway struct {
    // DB接続など
}

func NewTagGateway() *TagGateway {
    return &TagGateway{}
}

func (g *TagGateway) Create(ctx context.Context, t tag.Tag) error {
    // DB保存処理
    return nil
}

func (g *TagGateway) FindByNoteID(ctx context.Context, noteID string) ([]tag.Tag, error) {
    // DB取得処理
    return nil, nil
}

func (g *TagGateway) Delete(ctx context.Context, noteID, tagID string) error {
    // DB削除処理
    return nil
}
```

**ファイル: `internal/adapter/controller/tag_controller.go`**

```go
package controller

import (
    "net/http"
    "github.com/yourusername/backend-clean/internal/usecase"
)

type TagController struct {
    tagUseCase *usecase.TagInteractor
}

func NewTagController(tagUseCase *usecase.TagInteractor) *TagController {
    return &TagController{tagUseCase: tagUseCase}
}

func (c *TagController) AddTag(w http.ResponseWriter, r *http.Request) {
    // リクエストパース → UseCase呼び出し → レスポンス
}
```

---

## 🧪 ステップ4: テストコードもAIに作らせる（C1カバレッジ）

### 4-1. C1カバレッジとは？

```
C1（Branch Coverage）:
すべての分岐（if文、switch文）をテストする

例:
func CanPublish(note Note) error {
    if note.Status != StatusDraft {  // ← 分岐1
        return ErrInvalidStatus
    }
    if len(note.Sections) == 0 {     // ← 分岐2
        return ErrNoSections
    }
    return nil
}

C1達成には最低3ケース必要:
1. Status == Draft && Sections > 0 (成功)
2. Status != Draft (分岐1でエラー)
3. Sections == 0 (分岐2でエラー)
```

### 4-2. AIへの指示（テストコード生成）

```
先ほど生成したタグ機能のテストコードを、
C1カバレッジを達成するレベルで生成してください。

【要件】
1. すべての分岐（if文）をカバーする
2. テーブル駆動テストで書く
3. 正常系・異常系をすべて含める
4. モックを使ってDB不要にする

【生成してほしいテストファイル】
- internal/domain/tag/tag_test.go (Domain層)
- internal/usecase/tag_interactor_test.go (UseCase層)
```

---

## 📊 比較: backend-bad vs backend-clean（AI駆動開発）

### backend-bad でAIに指示した場合

```
指示: 「タグ機能を追加して」

AI:
「どこに書けばいいかわからない...」
「Serviceに追加しますか？それとも新しいファイル？」
「既存のコードが複雑で、どのパターンに従えばいいか不明...」

結果:
❌ AIが混乱する
❌ 生成したコードがバグだらけ
❌ テストコードが書けない（DB直結のため）
❌ 結局、手動で修正が必要（2時間）

👉 AI使えない
```

### backend-clean でAIに指示した場合

```
指示: 「タグ機能を追加して」

AI:
「はい！ドキュメントを読みました。」
「Noteと同じパターンで実装します。」

生成されたコード:
✅ Domain層: tag.go, tag_test.go
✅ Port層: tag_repository.go
✅ UseCase層: tag_interactor.go, tag_interactor_test.go
✅ Adapter層: tag_gateway.go, tag_controller.go

✅ すべて正しい場所に配置
✅ テストコードもC1カバレッジ達成
✅ レビュー修正のみ（30分）

👉 AI超使える！
```

---

## 🚀 実践: 実際にやってみよう

### ステップ1: AIにドキュメントを読ませる

```
以下のドキュメントを読んでください:
- backend-clean/docs/01_why_clean_architecture.md
- backend-clean/docs/02_clean_architecture_guide.md
- backend-clean/docs/03_testing_guide.md

また、以下の既存コードを読んで、パターンを理解してください:
- backend-clean/internal/domain/note/note.go
- backend-clean/internal/usecase/note_interactor.go
```

### ステップ2: 機能追加を指示

```
【課題】
このプロジェクトに「コメント機能」を追加してください。

【要件】
- ノートにコメントを追加できる
- コメントの削除ができる
- コメント一覧を取得できる

【指示】
1. Clean Architectureの設計方針に従って実装
2. 既存のNoteの実装パターンを参考にする
3. すべての層に必要なコードを生成
4. テストコードも生成（C1カバレッジ達成）
5. ファイルパスを明記
```

### ステップ3: 生成されたコードをレビュー

```
チェックリスト:
□ Domain層: エンティティとバリデーションがある
□ Domain層: テストコードがある（C1カバレッジ）
□ Port層: Interfaceが定義されている
□ UseCase層: ビジネスロジックが実装されている
□ UseCase層: テストコードがある（モック使用）
□ Adapter層: Gateway、Controllerが実装されている
□ すべてのファイルパスが正しい
```

---

## 💡 まとめ

### backend-bad の場合

```
❌ AIが混乱する
   → 構造が不明確、どこに書くべきかわからない

❌ 生成したコードがバグだらけ
   → パターンが不統一、AIが判断できない

❌ テストコードが書けない
   → DB直結、モックが作れない

結果:
  👉 AI駆動開発: 不可能
  👉 手動で7時間
```

### backend-clean の場合

```
✅ AIが正確に理解する
   → 構造が明確、ドキュメントで学習できる

✅ 生成したコードが高品質
   → パターンが統一、同じパターンで書ける

✅ テストコードも自動生成
   → モック使用、C1カバレッジ達成

結果:
  👉 AI駆動開発: 可能
  👉 AIで32分（13倍速い！）
```

### 🎯 AI駆動開発のコツ

```
1. ドキュメントを充実させる
   → AIが学習しやすい

2. パターンを統一する
   → AIが「同じパターンで書けばいい」と判断できる

3. 既存コードを見せる
   → AIが「これを真似すればいい」と理解できる

4. 明確な指示を出す
   → 「C1カバレッジ達成」「モック使用」など具体的に

結果:
  👉 AIがあなたの最高のペアプログラマーになる！
```

---

## 🎯 次のステップ

1. **ローカル環境をセットアップする**
   - [05_local_setup.md](./05_local_setup.md) を参照
   - Docker Compose、マイグレーション、テストまで実行

2. **実際にAIに指示してみる**
   - 「コメント機能を追加して」と指示
   - 生成されたコードをレビュー

3. **自分の機能を追加してみる**
   - 「いいね機能」「共有機能」など
   - AIと一緒に実装

4. **OpenAPI定義を変更してみる**
   - `api-schema/typespec/` の定義を編集
   - `pnpm run generate` でコード再生成
   - `make oapi` でGoコード再生成

**Happy AI-Driven Development!** 🎉
