# 🔥 なぜクリーンアーキテクチャが必要？ - backend-badとの比較

> 💡 **このドキュメントのゴール**
> 「クリーンアーキテクチャって面倒くさそう...」と思っているあなたに、
> 「あ、これやらないとヤバいんだ！」と腑に落ちてもらうためのドキュメントです。

---

## 🏠 2つの家を比較してみよう

このプロジェクトには、**2つのバックエンド**があります。どちらも同じAPIを提供しますが、中身は全然違います。

```
immortal-architecture-cqrs/
├── backend-bad/    ← ❌ グチャグチャの家（わざと悪く建てた）
└── backend-clean/  ← ✅ きれいに整理された家（このプロジェクト）
```

### 🏚️ backend-bad = ワンルームアパート（全部1部屋）

```
┌──────────────────────────────────────┐
│  ごちゃ混ぜの部屋                    │
│                                      │
│  🛏️ ベッド  🍳 キッチン  🚿 シャワー │
│  📚 本棚   🍽️ 食器     💻 PC      │
│  👔 服     🧴 洗剤     📺 TV      │
│                                      │
│  全部が1つの部屋！                   │
└──────────────────────────────────────┘

問題:
・どこに何があるかわからない
・掃除しにくい（全部ぐちゃぐちゃ）
・模様替えが超大変（全部動かす必要）
```

### 🏡 backend-clean = 一戸建て（部屋が分かれてる）

```
┌─────────┬─────────┬─────────┬─────────┐
│ 寝室    │ キッチン│ 浴室    │ 書斎    │
│         │         │         │         │
│ 🛏️ ベッド│🍳 コンロ│🚿シャワー│💻 PC   │
│ 👔 服   │🍽️ 食器  │🧴 洗剤  │📚 本棚 │
│         │         │         │         │
└─────────┴─────────┴─────────┴─────────┘

メリット:
・どこに何があるか一目瞭然
・掃除しやすい（部屋ごとに独立）
・模様替えが簡単（1部屋だけ変えればいい）
```

**同じ機能を提供**していますが、**整理の仕方が全然違います**。

これからbackend-badで起きる9つの問題を見ていきましょう！

---

## 📊 backend-bad で起きている9つの問題

| # | 何が悪いの？ | どう痛いの？ |
|---|-------------|-------------|
| **1** | 🔗 **ドメイン層がない** | API/DBが変わると全部書き直し |
| **2** | 🐘 **Serviceが何でも屋** | 1つの関数が100行超え |
| **3** | 🔒 **DBを直接呼ぶ** | テストできない |
| **4** | 🌳 **if文地獄** | 読めない、変更できない |
| **5** | 🗺️ **ルールが点在** | 「どこに書いてあるの？」 |
| **6** | 📝 **トランザクションコピペ** | Begin/Commit/Rollbackを毎回書く |
| **7** | ⚠️ **エラーが雑** | 「何が悪いの？」がわからない |
| **8** | 🔌 **設定を直参照** | テストで差し替えできない |
| **9** | 👻 **未実装に気づけない** | エンドポイント作り忘れ |

---

## 問題1: ドメイン層がない → ジェンガ崩壊

### ❌ 何が問題？

```go
// ❌ ServiceでTypeSpec型とsqlc型を直接扱う
func (s *TemplateService) Create(req openapi.CreateTemplateRequest) error {
    // OpenAPI の型をそのまま使ってる
    name := req.Name
    ownerID := req.OwnerID

    // sqlc の DB型に直接変換してる
    _, err := s.queries.CreateTemplate(ctx, sqldb.CreateTemplateParams{
        Name:        name,
        Description: req.Description,
        OwnerID:     uuid.MustParse(ownerID),
    })

    return err
}
```

**🔥 変更の連鎖反応が起きる！**

```
例: API仕様で「Name」を「Title」に変更したい

影響範囲:
├─ Controller (10箇所)
├─ Service (30箇所)
├─ DB型変換 (20箇所)
└─ テストコード (50箇所)

👉 全ファイル修正！半日作業！バグ混入リスク大！
```

### ✅ backend-clean の解決策

**戦略: ドメイン型という「緩衝材」を真ん中に挟む！**

```
┌─────────────┐
│  OpenAPI型  │  ← API設計が変わっても...
└──────┬──────┘
       │ 変換（Controllerで）
       ↓
┌─────────────┐
│  ドメイン型  │  ← ここは変わらない（安定）
└──────┬──────┘
       │ 変換（Gatewayで）
       ↓
┌─────────────┐
│  DB型       │  ← ここも独立
└─────────────┘

影響範囲: Controllerの変換処理だけ（1箇所、5分）
```

---

## 問題2: Serviceが何でも屋 → 100行関数

### ❌ 何が問題？

```go
// ❌ 1つの関数に全部詰め込む
func (s *NoteService) UpdateStatus(noteID, newStatus, actorID string) error {
    // 80行以上のメソッド...
    // - バリデーション
    // - DB取得
    // - ビジネスルール検証
    // - 状態遷移チェック
    // - 更新処理
    // - トランザクション
    // ...全部ここ！
}
```

**🔥 読めない！テストできない！変更できない！**

### ✅ backend-clean の解決策

**戦略: 責務を分ける！**

```
Domain層（ビジネスルール）:
  └─ CanPublish() - 10行（公開可能かチェック）

UseCase層（手順）:
  └─ Publish() - 30行（手順書みたいに読める）

Gateway層（DB）:
  └─ Update() - 15行（DB更新だけ）

👉 各関数10-30行！読みやすい！テストしやすい！
```

---

## 問題3: DBを直接呼ぶ → テスト不能

### ❌ 何が問題？

```go
// ❌ sqlcを直接呼ぶ
func (s *NoteService) Get(id string) (*Note, error) {
    return s.queries.GetNote(ctx, id)  // 👈 sqlcに直接依存
}
```

**🔥 テストするには本物のDBが必要！**

```
テストの課題:
├─ PostgreSQL起動が必要（遅い）
├─ データ準備が面倒
├─ 並列実行できない（データが競合）
└─ CIが不安定

👉 テストが遅い（1分）、誰も書かない
```

### ✅ backend-clean の解決策

**戦略: Interfaceで抽象化！**

```go
// ✅ Interfaceに依存
type NoteInteractor struct {
    notes port.NoteRepository  // ← Interface
}

func (u *NoteInteractor) Get(ctx context.Context, id string) error {
    n, err := u.notes.Get(ctx, id)  // ← どの実装か知らない
    // ...
}
```

**テスト時はMockを注入:**

```go
// テスト用のMock（DB不要）
mockRepo := &MockNoteRepository{
    notes: map[string]*note.Note{
        "123": {ID: "123", Title: "テスト"},
    },
}

interactor := &NoteInteractor{notes: mockRepo}

👉 テストが爆速（0.1秒）、DBなし、安定！
```

---

## 問題4: if文地獄 → 迷路

### ❌ 何が問題？

```go
// ❌ if文が5段ネスト
func (s *NoteService) Publish(...) error {
    if note.Status == "Draft" {
        if note.OwnerID == actorID {
            if len(note.Sections) > 0 {
                for _, section := range note.Sections {
                    if field.IsRequired && section.Content == "" {
                        // さらに続く...
                    }
                }
            }
        }
    }
    // 150行...誰も読めない
}
```

**🔥 迷路！読めない！変更できない！**

### ✅ backend-clean の解決策

**戦略: 関数に分ける！**

```go
// ✅ Domain層（ビジネスルール）
func CanPublish(note Note) error {
    if note.Status != StatusDraft {
        return domainerr.ErrInvalidStatus
    }
    if len(note.Sections) == 0 {
        return domainerr.ErrNoSections
    }
    // 10行、読みやすい
}

// ✅ UseCase層（手順書）
func (u *NoteInteractor) Publish(...) error {
    // 1. 取得
    note, _ := u.notes.Get(ctx, id)

    // 2. 検証（関数呼び出し）
    if err := service.CanPublish(note); err != nil {
        return err
    }

    // 3. 更新
    note.Status = StatusPublish
    return u.notes.Update(ctx, note)

    // 30行、手順書みたい！
}
```

---

## 問題5-9: その他の問題まとめ

### 5. ルールが点在 → 「どこに書いてあるの？」

```
❌ backend-bad: ビジネスルールが全ファイルに散らばってる
✅ backend-clean: domain/note/logic.go に集約
```

### 6. トランザクションコピペ

```
❌ backend-bad: Begin/Commit/Rollback を毎回コピペ
✅ backend-clean: TxManager で一元管理
```

### 7. エラーが雑

```
❌ backend-bad: return errors.New("error") （何のエラー？）
✅ backend-clean: return domainerr.ErrNotFound （意味がわかる）
```

### 8. 設定を直参照

```
❌ backend-bad: config.Get() を直呼び
✅ backend-clean: Factoryで注入、テストで差し替え可能
```

### 9. 未実装に気づけない

```
❌ backend-bad: エンドポイント実装忘れても気づけない
✅ backend-clean: TypeSpec Interfaceを実装、コンパイルエラーで気づける
```

---

## 📈 開発速度の変化: backend-bad vs backend-clean

### 🏚️ backend-bad の世界 = 負のスパイラル

```
開発速度の推移:
┌────────────────────────────────┐
│ 初期:   機能追加 1日  😊       │
│   ↓                            │
│ 3ヶ月:  機能追加 3日  😰       │
│   ↓   （if文が増えて複雑化）   │
│ 6ヶ月:  機能追加 1週間 😱      │
│   ↓   （誰も触りたくない）     │
│ 1年後:  機能追加 2週間 💀      │
│        「リプレイスしたい...」 │
└────────────────────────────────┘

💸 コスト:
- 新機能開発: 遅い（2週間）
- バグ修正: 多い（週10件）
- 離職率: 高い（半年で3人辞める）
```

### 🏡 backend-clean の世界 = 正のスパイラル

```
開発速度の推移:
┌────────────────────────────────┐
│ 初期:   機能追加 2日  🚶       │
│   ↓   （設計に時間かかる）     │
│ 3ヶ月:  機能追加 1日  🚀       │
│   ↓   （パターンが見えてくる） │
│ 6ヶ月:  機能追加 半日 🚀🚀     │
│   ↓   （安心して変更できる）   │
│ 1年後:  機能追加 2時間 🚀🚀🚀  │
│        「最高に楽しい！」      │
└────────────────────────────────┘

💰 コスト:
- 新機能開発: 速い（2時間）
- バグ修正: 少ない（月1-2件）
- 離職率: 低い（みんな楽しく働く）
```

---

## 💡 結論: 初期コストは投資、長期的には圧倒的にお得

### ⚖️ コスト比較

```
┌──────────────────────────────────────────┐
│  backend-bad: 最初は速いが、後で地獄     │
├──────────────────────────────────────────┤
│  初期:  機能追加 1日   ← 速い！          │
│  3ヶ月: 機能追加 3日   ← 遅くなる        │
│  6ヶ月: 機能追加 1週間 ← さらに遅い      │
│  1年後: 機能追加 2週間 ← 地獄            │
│                                          │
│  累計コスト: 超高い（リプレイス必要）    │
└──────────────────────────────────────────┘

┌──────────────────────────────────────────┐
│  backend-clean: 最初は遅いが、後で爆速   │
├──────────────────────────────────────────┤
│  初期:  機能追加 2日   ← ちょっと遅い    │
│  3ヶ月: 機能追加 1日   ← 速くなる        │
│  6ヶ月: 機能追加 半日  ← さらに速い      │
│  1年後: 機能追加 2時間 ← 爆速            │
│                                          │
│  累計コスト: 安い（ずっと成長し続ける）  │
└──────────────────────────────────────────┘
```

### 🎯 クリーンアーキテクチャを選ぶべき理由

```
✅ 変更に強い
   → 影響範囲が限定的（5分で修正）

✅ テストしやすい
   → DBなし、爆速（0.1秒）

✅ 読みやすい
   → 各ファイル10-30行、手順書みたい

✅ AI時代に最適
   → 構造が明確、AIが正確にコード生成

✅ チーム開発しやすい
   → 新メンバーがすぐ理解できる

✅ 長期運用に強い
   → ライブラリ変更、仕様変更が楽

結論:
  👉 初期コストは投資
  👉 長期的には圧倒的にお得
  👉 開発が楽しい、チームが幸せ
```

---

## 🚀 次のステップ

動機づけは完了！次は実際の処理の流れを理解しましょう：

**👉 [02_clean_architecture_guide.md](./02_clean_architecture_guide.md) へ進む**

このドキュメントでは、以下を学びます：
- 同心円図の意味
- リクエストがどう処理されるか
- Interface、Factoryとは何か
- FAQ（よくある質問）
- チェックリスト（コードを書く前に）

**Happy Learning!** 🎉
