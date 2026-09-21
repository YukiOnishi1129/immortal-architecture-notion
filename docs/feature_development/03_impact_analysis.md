# ③ 影響範囲とデグレ分析

> **この章のゴール**
> 「デグレしません」を、**感覚ではなく根拠**で言えるようになる。

---

## なぜ影響範囲を調べるのか

```
レビュアー「これ、既存の機能壊れない？」

❌ ダメな回答:
   「たぶん大丈夫です」
   「動作確認したので大丈夫です」

✅ 良い回答:
   「既存のnotesテーブルは1列も変更していません。
     既存メソッドは ChangeStatus と Update の2つを変更します。
     この2つにはテストが無かったので、変更前に先にテストを書きました。
     既存テストは全部パスしています」
```

**根拠を持って言い切るために、事前に調べます。**

---

## 影響範囲の調べ方

### 手順

```
1. 変更するものを「追加」と「変更」に分ける
2. 層ごとに何を触るか洗い出す
3. 「変更」するものについて、それを使っている箇所を全部探す
4. 既存テストがどこまで守ってくれるかを確認する
5. 守ってくれない範囲を、リスクとして明記する
```

### 最重要の原則: 「追加」と「変更」を分ける

```
┌─────────────────────────────────────────┐
│ 追加（新規ファイル・新規テーブル）          │
│   → デグレリスク: ほぼゼロ                 │
│   → 既存コードは何も知らないまま動き続ける   │
├─────────────────────────────────────────┤
│ 変更（既存ファイル・既存テーブル）          │
│   → デグレリスク: ここに集中する           │
│   → 使っている箇所を全部調べる必要がある    │
└─────────────────────────────────────────┘
```

**良い設計とは、「変更」を最小にして「追加」に寄せた設計です。**

---

## 今回の変更一覧

第2章で決めた方針に基づいて洗い出します。

### ✅ 追加するもの（デグレリスクほぼゼロ）

| 層 | ファイル / 対象 | 内容 |
|---|---|---|
| **DB** | `migrations/xxx_add_notion_sync.up.sql` | 新規テーブル2つ |
| **Domain** | `internal/domain/notion/entity.go` | 連携状態のエンティティ |
| **Domain** | `internal/domain/notion/logic.go` | 連携可否の判定ロジック |
| **Port** | `internal/port/notion_port.go` | 連携のIn/Outポート、外部APIのポート |
| **UseCase** | `internal/usecase/notion_sync_interactor.go` | 連携のユースケース |
| **Gateway** | `internal/adapter/gateway/notion/client.go` | Notion APIクライアント |
| **Gateway** | `internal/adapter/gateway/db/sqlc/notion_sync_repository.go` | 連携状態の永続化 |
| **Controller** | `internal/adapter/http/controller/notion_controller.go` | エンドポイント |
| **Presenter** | `internal/adapter/http/presenter/notion_presenter.go` | レスポンス整形 |
| **Config** | 環境変数 `NOTION_API_KEY` の読み込み | ※configは変更 |

### ⚠️ 変更するもの（ここを重点的に調べる）

| ファイル | 変更内容 | 危険度 |
|---|---|---|
| **`internal/usecase/note_command_interactor.go`** | **`ChangeStatus` にNotion連携を追加** | 🔴 **高** |
| **`internal/usecase/note_command_interactor.go`** | **`Update` にNotion更新を追加** | 🔴 **高** |
| `migrations/` | `templates` に `notion_destination_id` 列を追加 | 🟡 中 |
| `internal/driver/config/config.go` | `NotionAPIKey` フィールド追加 | 🟢 低 |
| `internal/driver/initializer/api/initializer.go` | DI配線を追加 | 🟡 中 |
| `internal/driver/factory/*.go` | ファクトリ追加 | 🟢 低 |
| `api-schema/typespec/` | API定義追加 → 生成物が変わる | 🟡 中 |

### 変更しないもの

```
✅ notes テーブル             … カラムは1列も変更なし
✅ sections テーブル          … 変更なし
✅ note_read_models テーブル  … 変更なし
✅ NoteQueryInteractor        … 1行も変更なし
✅ Create（ノート作成）        … 下書きで作られるので連携不要
✅ Delete（ノート削除）        … Notionページは残す仕様のため
✅ 既存の全ドメインロジック     … 変更なし
```

> **論点Cの成果**
> 連携状態を別テーブルにしたので、`notes` のカラムは無変更です。
> ReadModel の再構築も不要になりました。

---

## 🔴 最大のリスク: テストのないコードを変更する

第2章で「公開に連動させる」と決めた結果、
**`NoteCommandInteractor` の2メソッドを変更**することになりました。

そして、このファイルには**テストがありません**。

```bash
ls internal/usecase/*_test.go
```

```
account_interactor_test.go       ✅ ある
template_interactor_test.go      ✅ ある
note_command_interactor_test.go  ❌ ない
note_query_interactor_test.go    ❌ ない
```

```
テストのないコードを変更する
    ↓
壊しても気づけない
    ↓
「公開できなくなった」が本番で発覚
```

### 対策: 変更する前にテストを書く

**これは今回のスコープに含めます。** 「別チケットで」では間に合いません。

```
PR②（仮）: ChangeStatus / Update の既存振る舞いのテストを追加
   ↓ ここで現状の動作が固定される
PR⑥: Notion連携を追加
   ↓ 既存テストが通れば、壊していないと言える
```

> **「ついでに直す」との違い**
> 第3章の後半で「ついでに直すな」と書きますが、これは別です。
>
> ```
> ついでに直す   … 今回触らない箇所の改善 → 別チケット
> 先にテストを書く … 今回触る箇所の安全確保 → 今回のスコープ
> ```
>
> **触るコードにテストが無いなら、それは今回の仕事です。**

---

## 「変更」するものを深掘りする

### 変更① templates への列追加

```sql
ALTER TABLE templates
  ADD COLUMN notion_destination_id UUID
    REFERENCES notion_destinations(id);
```

#### 使っている箇所を全部探す

```bash
grep -rn "templates" backend-clean/internal --include="*.go" | grep -i "select\|query"
```

**調査結果:**

| 影響先 | 内容 | 判定 |
|---|---|---|
| `template_repository.go` | `SELECT * FROM templates` を使っているか？ | ⚠️ 要確認 |
| `sqlc/generated/` | sqlc がスキーマから型を再生成する | ⚠️ 再生成必要 |
| `Template` エンティティ | 列を足しても、ドメインに足さなければ影響なし | ✅ 制御可能 |

> 🚨 **`SELECT *` を使っていたら要注意**
> 列を足すと、返ってくるカラム数が変わります。
> sqlc は型を再生成するので、コンパイルエラーで気づけます。
> **コンパイルエラーになるのは良いこと**です（黙って壊れるより安全）。

#### NULL許容にする理由

```sql
-- ✅ NULL許容にする
ADD COLUMN notion_destination_id UUID NULL

-- ❌ NOT NULL にすると既存行が全部エラー
ADD COLUMN notion_destination_id UUID NOT NULL  -- 既存テンプレートに値がない！
```

**既存データを壊さないマイグレーションの鉄則です。**

### 変更② config.go

```go
type Config struct {
	DatabaseURL    string
	ServerPort     int
	AllowedOrigins []string
	NotionAPIKey   string  // ← 追加
}
```

#### 危険なパターン

```go
// ❌ これをやると既存環境が起動しなくなる
if os.Getenv("NOTION_API_KEY") == "" {
    return nil, errors.New("NOTION_API_KEY is not set")
}
```

**なぜダメか:**

```
本番環境に NOTION_API_KEY を設定し忘れる
    ↓
アプリが起動しない
    ↓
Notion連携どころか、全機能が停止
```

#### 安全なパターン

```go
// ✅ 未設定でも起動する。連携機能だけ無効になる
notionKey := os.Getenv("NOTION_API_KEY")
// 空でもエラーにしない
```

そして連携実行時にチェックします。

```go
if cfg.NotionAPIKey == "" {
    return domainerr.ErrNotionNotConfigured  // この機能だけ失敗
}
```

> **原則: 新機能の設定不備で、既存機能を止めない**
> 既存の `DatabaseURL` は必須チェックがあります（無いと動かないので正しい）。
> 新機能の設定は、**任意**にします。

### 変更③ initializer.go（DI配線）

```go
// 追加する配線
notionClient := notiongateway.NewClient(cfg.NotionAPIKey)
notionSyncRepoFactory := factory.NewNotionSyncRepoFactory(pool)
notionInputFactory := factory.NewNotionSyncInputFactory()
notionOutputFactory := httpfactory.NewNotionOutputFactory()

nsc := httpcontroller.NewNotionController(...)
server := httpcontroller.NewServer(ac, nc, tc, nsc)  // ← 引数が増える
```

#### 影響先

| 影響先 | 内容 |
|---|---|
| `NewServer()` のシグネチャ | 引数が増える → 呼び出し箇所を全部直す |
| `initializer_test.go` | 既存テストがある → 修正必要 |
| `grpc/initializer.go` | gRPC側も `NewServer` を使っているか確認 |

```bash
# 呼び出し箇所を探す
grep -rn "NewServer(" backend-clean/internal
```

---

## 既存テストはどこまで守ってくれるか

**ここが最重要セクションです。**

### 実測データ

このリポジトリで実際に計測した結果です。

```bash
go test ./internal/... -cover
```

| パッケージ | カバレッジ | 評価 |
|---|---:|---|
| `domain/service` | **100.0%** | ✅ 完璧 |
| `domain/template` | 97.4% | ✅ 手厚い |
| `driver/config` | 95.2% | ✅ 手厚い |
| `domain/account` | 92.6% | ✅ 手厚い |
| `domain/note` | 86.3% | ✅ 良好 |
| `adapter/http/controller` | 84.7% | ✅ 良好 |
| `adapter/http/presenter` | 50.9% | ⚠️ 中程度 |
| **`usecase`** | **37.2%** | 🚨 **低い** |
| `driver/factory` | 0.0% | ⚠️ テストなし |
| `adapter/gateway/notion`（新規） | - | 🚨 これから書く |

### 守ってくれる範囲・くれない範囲

```
┌────────────────────────────────────────────┐
│ ✅ 既存テストが守ってくれる                   │
├────────────────────────────────────────────┤
│ ・ドメインロジック（公開可否・バリデーション）  │
│   → domain/service が100%                   │
│ ・HTTPコントローラーの既存エンドポイント        │
│   → 84.7%                                   │
│ ・config の読み込み                          │
│   → 95.2%（ただし新フィールドは自分で書く）    │
└────────────────────────────────────────────┘

┌────────────────────────────────────────────┐
│ ❌ 既存テストが守ってくれない                 │
├────────────────────────────────────────────┤
│ ・NoteCommandInteractor の振る舞い           │
│   → テストが存在しない（今回書く）           │
│ ・DI配線（initializer）                      │
│   → カバレッジ 0.0%                          │
│ ・マイグレーションの正しさ                     │
│   → テストで検証されていない                  │
│ ・Notion APIとの実際の通信                    │
│   → 当然ない（これから書く）                  │
└────────────────────────────────────────────┘
```

---

## ⚠️ リスク一覧と対策

| # | リスク | 影響度 | 対策 |
|---|---|:---:|---|
| R1 | **テストのない ChangeStatus / Update を変更して壊す** | 🔴 **高** | **変更前に既存振る舞いのテストを書く**（今回のスコープ） |
| R2 | **Notion障害中にノートを公開できない** | 🔴 **高** | 設計上の割り切り（論点A-2）。**ビジネス側の合意が前提** |
| R3 | `NOTION_API_KEY` 未設定で全体が起動しない | 🔴 高 | 必須チェックを**入れない**。連携時のみ検証 |
| R4 | Notion成功後のDB失敗で状態がずれる | 🟡 中 | 許容（論点D）。監視で検知し、手動対応 |
| R5 | `templates` の列追加で sqlc 再生成漏れ | 🟡 中 | CIで `make generate` の差分チェック |
| R6 | `NewServer` 引数追加で gRPC側が壊れる | 🟡 中 | 呼び出し箇所を grep で全部確認 |
| R7 | 外部API呼び出しでトランザクション長期化 | 🔴 高 | 設計で回避済み（論点D）。レビューで確認 |
| R8 | Notion APIのレート制限でエラー多発 | 🟡 中 | リトライとバックオフを実装 |
| R9 | 既存の公開済みノートが連携されないまま | 🟡 中 | 詳細ページに手動連携ボタン（論点A） |

> **R2 が他と性質が違う**
> これは「対策する」ものではなく、**設計として選んだトレードオフ**です。
> 消せるリスクではないので、**合意を取ることが対策**になります。
>
> レビューで「Notionが落ちたら公開できません」と明示し、
> それでよいかを確認します。

## ロールバック手順

「戻せます」と言えることも、レビューで問われます。

### コードのロールバック

```
✅ 容易
理由: 既存コードへの変更が最小限のため、
      revert してもコンフリクトしにくい
```

### DBのロールバック

```sql
-- down マイグレーション
ALTER TABLE templates DROP COLUMN notion_destination_id;
DROP TABLE note_notion_syncs;
DROP TABLE notion_destinations;
```

```
⚠️ 注意
連携済みのデータ（notion_page_id）が消えます。
戻したあと再度適用すると、
Notionに重複ページが作られる可能性があります。
```

### 機能だけ止める（推奨）

```
NOTION_API_KEY を空にする
    ↓
連携機能だけが無効になる
    ↓
他の機能は動き続ける
    ↓
DBは触らないのでデータも無事
```

> **設定で止められる設計にしておく**
> ロールバックは「コードを戻す」だけではありません。
> **環境変数1つで止まる**なら、深夜の障害対応が楽になります。

---

## レビュアー向けチェックポイント

PRを出すときに、この観点を明記すると親切です。

```markdown
## レビュー観点

### 必ず見てほしい
- [ ] トランザクションの中でNotion APIを呼んでいないか
      （internal/usecase/notion_sync_interactor.go）
- [ ] NOTION_API_KEY 未設定でもアプリが起動するか
      （internal/driver/config/config.go）
- [ ] マイグレーションが NULL許容になっているか

### 確認済み（報告）
- [x] notes テーブルのカラムは無変更
- [x] ChangeStatus / Update の既存振る舞いテストを先に追加した
- [x] 既存テストは全てパス（go test ./internal/...）
- [x] NewServer の呼び出し箇所は全て更新済み

### 合意済みのトレードオフ
- Notion障害中はノートを公開できない（論点A-2でビジネス側と合意済み）
```

---

## ✅ この章のチェックリスト

```
□ 「追加」と「変更」を分けて整理したか
□ 「変更」するものの使用箇所を grep で全部探したか
□ 既存テストのカバレッジを実測したか（推測していないか）
□ 守ってくれない範囲を明記したか
□ 新機能の設定不備で既存機能が止まらない設計か
□ ロールバック手順を書いたか
□ 「ついでに直す」を我慢して別チケットにしたか
```

---

## 次の章へ

何を触り、どこが壊れうるかがわかりました。

次は「どう動けば合格か」を決めます。
仕様が決まった今なら、実装を待たずに書けます。

👉 [04_test_strategy.md](./04_test_strategy.md)
