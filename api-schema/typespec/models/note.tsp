import "@typespec/http";
import "@typespec/openapi3";
import "./account.tsp";

using TypeSpec.Http;

namespace MiniNotion.Models;

/** ノートのステータス */
enum NoteStatus {
  /** 下書き */
  Draft: "Draft",

  /** 公開 */
  Publish: "Publish",
}

/** セクション（ノートの各項目） */
model Section {
  /** セクションID */
  id: string;

  /** フィールドID */
  fieldId: string;

  /** フィールドラベル */
  fieldLabel: string;

  /** 内容 */
  content: string;

  /** 必須項目かどうか */
  isRequired: boolean;
}

/** セクション作成リクエスト */
model CreateSectionRequest {
  /** フィールドID */
  fieldId: string;

  /** 内容 */
  content: string;
}

/** セクション更新リクエスト */
model UpdateSectionRequest {
  /** セクションID */
  id: string;

  /** 内容 */
  content: string;
}

/** ノート作成リクエスト */
model CreateNoteRequest {
  /** タイトル */
  @minLength(1)
  @maxLength(100)
  title: string;

  /** テンプレートID */
  @format("uuid")
  templateId: string;

  /** 所有者ID */
  @format("uuid")
  ownerId: string;

  /** セクション（オプション） */
  sections?: CreateSectionRequest[];
}

/** ノート更新リクエスト */
model UpdateNoteRequest {
  /** ノートID */
  id: string;

  /** タイトル */
  @minLength(1)
  @maxLength(100)
  title: string;

  /** セクション */
  sections: UpdateSectionRequest[];
}

/** ノートレスポンス */
model NoteResponse {
  /** ノートID */
  id: string;

  /** タイトル */
  title: string;

  /** テンプレートID */
  templateId: string;

  /** テンプレート名 */
  templateName: string;

  /** 所有者ID */
  ownerId: string;

  /** 所有者情報 */
  owner: AccountSummary;

  /** ステータス */
  status: NoteStatus;

  /** セクション */
  sections: Section[];

  /** 作成日時 */
  createdAt: utcDateTime;

  /** 更新日時 */
  updatedAt: utcDateTime;
}

/** ノートフィルター（クエリパラメータ） */
model NoteFilters {
  /** タイトルキーワード検索 */
  @query
  q?: string;

  /** ステータスフィルター */
  @query
  status?: NoteStatus;

  /** テンプレートIDフィルター */
  @query
  templateId?: string;

  /** 所有者IDフィルター */
  @query
  ownerId?: string;
}
