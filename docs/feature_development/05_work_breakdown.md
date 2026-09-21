# ⑤ 段取りとPR分割

> **この章のゴール**
> 1つの機能を、**レビュー可能な単位**に割れるようになる。

---

## ❌ 巨大PRで起きること

第3章で変更するファイルが15個ほど、第4章でQAケースが18個あるとわかりました。これを1つのPRで出すとどうなるか。

```
「Notion連携機能」PR
  変更ファイル: 15個
  追加行数: +1200行
```

```
レビュアーの心の声:
  「うっ…あとで見よう」
     ↓ 3日後
  「まだ見てない…」
     ↓ 5日後
  「LGTM」（ちゃんと見ていない）
     ↓
  バグがすり抜ける
```

さらに悪いこと。

```
マイグレーション + API + 画面 が1つのPRに入っている
     ↓
「画面の文言だけ直したい」
     ↓
でもrevertするとテーブルも消える
     ↓
部分的に戻せない
```

---

## ✅ PR分割の3原則

```
┌────────────────────────────────────────────┐
│ 原則1: 常にmainが動く状態を保つ              │
│   → どのPRをマージしても、アプリは壊れない    │
├────────────────────────────────────────────┤
│ 原則2: 1PRは1つの関心事だけ                 │
│   → レビュアーが「何を見ればいいか」わかる    │
├────────────────────────────────────────────┤
│ 原則3: 下の層から積み上げる                  │
│   → DB → Domain → UseCase → API → 画面      │
└────────────────────────────────────────────┘
```

### 原則1の実現方法: 使われないコードを先に入れる

```
❓ 「まだ使われないコードをマージしていいの？」
✅ いいです。それが安全な進め方です。

PR①: テーブルを作る（誰も使わない）
      → マージしてもアプリの挙動は1ミリも変わらない
      → デグレリスク: ゼロ

PR②: ドメインモデルを作る（誰も呼ばない）
      → マージしても既存コードは何も知らない
      → デグレリスク: ゼロ

PR⑤: 最後に配線する
      → ここで初めて機能が動き出す
```

> **これが「追加に寄せる」設計の威力**
> 第3章で「追加」と「変更」を分けたのは、このためです。
> 追加だけのPRは、いくつマージしても安全です。

---

## PR分割案（7本）

### 依存関係

```
PR⓪  既存テストの追加（ChangeStatus / Update）  ← 安全網を先に張る
  │
  ▼
PR①  DBマイグレーション
  │
  ├──────────────┐
  ▼              ▼
PR②  ドメイン    PR③  Notion APIクライアント
  │              │
  └──────┬───────┘
         ▼
PR④  リポジトリ + 送り先の管理
         │
         ▼
PR⑤  ChangeStatus / Update への組み込み  ← ここで機能が動く
         │
         ▼
PR⑥  フロントエンド + 手動連携ボタン
```

> **PR② と PR③ は並行作業できます**
> 2人いれば同時に進められる、という情報も段取りの一部です。

---

### PR⓪ 既存テストの追加 ⚠️

```
目的: これから変更するコードに、安全網を張る
変更: internal/usecase/note_command_interactor_test.go を新規作成
行数: +250行程度
依存: なし（最初にやる）
```

**なぜ最初にやるか**

第3章で判明したとおり、`NoteCommandInteractor` にはテストがありません。
PR⑤でこのファイルの `ChangeStatus` と `Update` を変更します。

```
テストのないコードを変更する
   ↓
壊しても気づけない
```

**内容**

現状の振る舞いを固定するテストを書きます。**Notion連携は一切含めません。**

```
□ ChangeStatus: Draft → Publish ができる
□ ChangeStatus: Publish → Draft ができる
□ ChangeStatus: 他人のノートは ErrUnauthorized
□ ChangeStatus: 不正な遷移は ErrInvalidStatusChange
□ ChangeStatus: ReadModel が更新される
□ Update: タイトルとセクションが更新される
□ Update: 他人のノートは ErrUnauthorized
□ Update: 空タイトルは ErrTitleRequired
□ Create / Delete も同様に
```

**受け入れ基準**

```
□ usecase パッケージのカバレッジが 37.2% → 70%以上
□ プロダクションコードを1行も変更していない
□ 全テストがパスする
```

> **これは「ついでの改善」ではありません**
> 今回変更するコードの安全確保なので、**今回のスコープ**です。
> これを別チケットにすると、PR⑤で無防備な変更をすることになります。

---

### PR① DBマイグレーション

```
目的: テーブルを用意する
変更: migrations/ に2ファイル追加
行数: +60行程度
```

**内容**

[02_design_policy.md](./02_design_policy.md) で決めた3つの変更を入れます。

```
CREATE TABLE notion_destinations   … 送り先マスタ
CREATE TABLE note_notion_syncs     … 連携状態（notes削除時にCASCADE）
ALTER TABLE templates              … notion_destination_id を追加（NULL許容）
```

**受け入れ基準**

```
□ up を実行してもアプリが正常に動く
□ down を実行して元に戻せる
□ 既存の templates 行がエラーにならない（NULL許容の確認）
□ sqlc generate を実行し、生成物の差分をコミット
□ 既存テストが全てパスする
```

> ⚠️ **このPRで一番危ないのは `ALTER TABLE templates`**
> sqlc の再生成が必要です。忘れると後続PRでビルドが通りません。

---

### PR② ドメイン層

```
目的: 連携という概念をコードで表現する
変更: internal/domain/notion/ を新規作成
行数: +150行程度
依存: PR①（不要。並行可能）
```

**内容**

```
internal/domain/notion/
  ├── entity.go     … SyncStatus, Sync エンティティ
  ├── logic.go      … 連携可否の判定、ページID有無の判定
  └── logic_test.go … ドメインロジックのテスト
```

```go
type Sync struct {
    NoteID        string
    NotionPageID  *string   // nil なら未連携
    NotionPageURL *string
    Status        SyncStatus  // pending / synced / failed
    ErrorMessage  *string
    LastSyncedAt  *time.Time
}

// IsFirstSync returns true if the note has never been synced.
func (s *Sync) IsFirstSync() bool {
    return s == nil || s.NotionPageID == nil
}
```

**受け入れ基準**

```
□ ドメイン層が外部パッケージに依存していない
   （net/http も database/sql も import しない）
□ テストカバレッジ 90%以上
   （既存の domain/service は100%。基準を揃える）
□ 既存コードから一切参照されていない（＝影響ゼロ）
```

> **なぜドメイン層を先に作るか**
> ここが決まると、上の層の形が自動的に決まります。
> 逆に、ここが曖昧だと全部やり直しになります。

---

### PR③ Notion APIクライアント

```
目的: 外部APIを叩く部品を作る
変更: internal/adapter/gateway/notion/ を新規作成
行数: +250行程度
依存: PR②（ポート定義が必要なら）
```

**内容**

```
internal/port/notion_port.go          … インターフェース定義
internal/adapter/gateway/notion/
  ├── client.go       … Notion APIクライアント実装
  ├── client_test.go  … httptest でモックサーバーを立ててテスト
  └── converter.go    … ノート → Notionブロック の変換
```

```go
// port/notion_port.go
type NotionClient interface {
    CreatePage(ctx context.Context, parentPageID, title string, blocks []Block) (*PageResult, error)
    UpdatePage(ctx context.Context, pageID, title string, blocks []Block) (*PageResult, error)
}
```

**受け入れ基準**

```
□ タイムアウトが設定されている（無限に待たない）
□ リトライがある（429, 5xx のみ。4xxはリトライしない）
□ テストで実際のNotion APIを呼んでいない（httptest を使う）
□ APIキーがログに出ない
□ 既存コードから一切参照されていない
```

> 🚨 **タイムアウトを必ず設定する**
> ```go
> client := &http.Client{Timeout: 10 * time.Second}
> ```
> デフォルトの `http.Client` は**タイムアウトが無限**です。
> Notionが応答しないと、リクエストが永久に詰まります。

---

### PR④ リポジトリ + 送り先の管理

```
目的: 永続化と、送り先を設定できるようにする
変更: gateway/db と 送り先の CRUD を追加
行数: +300行程度
依存: PR①②③
```

**内容**

```
internal/adapter/gateway/db/sqlc/notion_sync_repository.go
internal/adapter/gateway/db/sqlc/notion_destination_repository.go
internal/usecase/notion_destination_interactor.go
internal/adapter/http/controller/notion_destination_controller.go
```

送り先（notion_destinations）の登録・一覧・テンプレートへの割り当てを作ります。

**受け入れ基準**

```
□ 送り先を登録・一覧・編集できる
□ テンプレートに送り先を割り当てられる
□ テンプレート作成時に送り先が必須になっている
□ 既存コードの振る舞いは変わっていない（PR⓪のテストが通る）
```

> **PR⑤より先に送り先を設定できるようにする**
> 順序が逆だと、リリース前の移行作業（全テンプレートへの設定）が
> できません。

---

### PR⑤ ChangeStatus / Update への組み込み ⚡

```
目的: 公開と連携を連動させる
変更: note_command_interactor.go, config.go, initializer.go
行数: +250行程度
依存: PR④
```

**⚡ このPRで機能が動き出します。同時に、最もリスクが高いPRです。**

**内容**

```go
// ChangeStatus の変更イメージ

// 公開する場合
if input.Status == note.StatusPublish {
    dest, err := u.destinations.GetByTemplate(ctx, current.Note.TemplateID)
    if err != nil {
        return err   // 送り先未設定なら公開しない
    }
    // 先にNotionを呼ぶ（トランザクション外）
    result, err := u.notion.PublishPage(ctx, dest, current)
    if err != nil {
        return err   // 失敗したらDBを触らない = ロールバック不要
    }
    // 成功したのでDBを更新
    ...
}
```

**受け入れ基準**

```
□ 🚨 PR⓪で追加した既存テストが全てパスする（デグレなし）
□ 🚨 トランザクションの中でNotion APIを呼んでいない
□ 🚨 Notionを先に呼び、成功後にDBを更新している
□ 公開でNotionにページができる（QA-01）
□ 公開中の編集でNotionも更新される（QA-02）
□ 非公開でNotionがゴミ箱に入る（QA-03）
□ 再公開で同じページが戻る（QA-04）
□ Notion失敗時に公開されない（QA-20）
□ 送り先未設定で公開できない（QA-30）
□ NOTION_API_KEY 未設定でもアプリが起動する（QA-22）
```

> ⚠️ **レビューを最も厳しくするPR**
> 既存の公開機能を変更します。
> PR⓪のテストが通ることを、必ず確認してください。

---

### PR⑥ フロントエンド + 手動連携ボタン

```
目的: ユーザーが使えるようにする
変更: frontend/src/features/note/ と template/
行数: +300行程度
依存: PR⑤
```

**内容**

```
- テンプレート編集画面に送り先の選択欄
- ノート詳細に「Notionで開く」リンク（公開中のみ）
- 公開/非公開の失敗時にエラー表示
- 公開済み & 未連携のノートに「連携する」ボタン（既存データ救済）
```

**受け入れ基準**

```
□ 公開失敗時にエラー内容が表示される
□ もう一度押せば再試行できる
□ 公開済み未連携のノートにボタンが出る（QA-23）
□ 処理中はボタンが無効化される（QA-34）
```

---

## PR一覧まとめ

| PR | 内容 | 行数 | 依存 | デグレリスク | レビュー重点 |
|:---:|---|---:|:---:|:---:|---|
| ⓪ | **既存テスト追加** | ~250 | - | 🟢 ゼロ | カバレッジが上がったか |
| ① | マイグレーション | ~60 | - | 🟢 低 | NULL許容、sqlc再生成 |
| ② | ドメイン層 | ~150 | - | 🟢 ゼロ | 依存の方向 |
| ③ | Notionクライアント | ~250 | ② | 🟢 ゼロ | タイムアウト、リトライ |
| ④ | リポジトリ+送り先管理 | ~300 | ①②③ | 🟡 中 | テンプレート編集への影響 |
| ⑤ | **既存メソッドへの組み込み** | ~250 | ④ | 🔴 **高** | **PR⓪のテストが通るか** |
| ⑥ | フロントエンド | ~300 | ⑤ | 🟡 中 | エラー表示、二重送信防止 |

> **PR⓪ を先頭に置いた意味**
> PR⑤ で既存コードを変更するとき、
> **PR⓪ のテストが通るかどうかが唯一の判断材料**になります。
>
> これが無いと「たぶん壊れていません」としか言えません。

## 設定で切り替えられるようにする

PR⑤ で**既存の公開機能そのものを変更**します。
そのため「連携をオフにして従来どおり動かす」手段が要ります。

```go
if cfg.NotionAPIKey == "" {
    // 連携をスキップして、従来どおりの公開処理
    return u.changeStatusWithoutNotion(ctx, input)
}
```

```
NOTION_API_KEY 未設定 → 従来どおり公開できる
NOTION_API_KEY 設定済み → Notion連携が動く
```

これで得られるもの。

```
□ 連携部分だけを切り離して動作確認できる
□ トークンが無い開発者も、既存機能の開発を続けられる
□ 不具合が出たとき、設定を外せば切り分けられる
```

> **第3章の「設定不備で既存機能を止めない」の応用です**
> 未設定時に**エラーにせず従来動作にフォールバック**させることで、
> 切り替えスイッチを兼ねさせています。
>
> ただし**恒久的な分岐にはしないこと**。
> 全テンプレートに送り先を設定したら、このフォールバックは削除します。
> （残すと「連携されない公開」が永久に可能になり、論点Aが崩れます）

## スケジュール例

```
Day 1  │ PR⓪ 既存テスト追加          ████████
Day 2  │ PR① マイグレーション        ████
       │ PR② ドメイン層              ████     （並行）
Day 3  │ PR③ Notionクライアント      ████████
Day 4  │ PR④ リポジトリ+送り先管理    ████████
Day 5  │ PR⑤ 既存メソッドへの組み込み ████████
Day 6  │ PR⑥ フロントエンド          ████████
Day 7  │ 送り先の設定                ████
       │ 本物のNotionと繋いで確認     ████
```

**クリティカルパス: ⓪ → ① → ④ → ⑤ → ⑥ → 確認**

> **PR⓪ が1日増えている**
> 既存テストを書く分、着手が1日遅れます。
> それでも、PR⑤ で壊したときの手戻りより安いと判断しました。

## PRテンプレート

各PRの説明文は、この型で書くとレビューが速くなります。

```markdown
## 何をするPRか
Notion連携のためのテーブルを追加します。

## なぜ必要か
[02_design_policy.md 論点C](./02_design_policy.md) の決定に基づき、
連携状態を notes とは別テーブルで管理します。

## 影響範囲
- 既存テーブル: templates に列を1つ追加（NULL許容）
- 既存コード: 変更なし
- 既存の挙動: 変わりません

## レビュー観点
- [ ] NULL許容になっているか（既存行が壊れないか）
- [ ] down で戻せるか
- [ ] sqlc の生成物がコミットされているか

## 確認したこと
- [x] make migrate-up / migrate-down 両方実行
- [x] go test ./internal/... 全てパス
- [x] 既存のテンプレート一覧APIが正常動作

## 次のPR
PR② ドメイン層の追加
```

---

## ✅ この章のチェックリスト

```
□ どのPRをマージしてもmainが動くか
□ 1PRが1つの関心事に絞られているか
□ 下の層から積み上げる順になっているか
□ 各PRに受け入れ基準を書いたか
□ デグレリスクが高いPRを明示したか
□ 依存関係と並行可能な作業を示したか
□ 1PRの行数が300行程度に収まっているか
```

---

## 次の章へ

設計も段取りも決まりました。
ここまで来ると、**実装はAIに任せられます**。

次章で、ここまでの成果物をAIへの指示に変えます。

👉 [06_implementation.md](./06_implementation.md)
