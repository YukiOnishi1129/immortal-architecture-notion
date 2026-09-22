import "@typespec/http";
import "@typespec/openapi3";
import "./account.tsp";

using TypeSpec.Http;

namespace MiniNotion.Models;

/** テンプレートフィールド */
model Field {
  /** フィールドID */
  id: string;

  /** ラベル */
  @minLength(1)
  label: string;

  /** 表示順序 */
  @minValue(1)
  order: int32;

  /** 必須フラグ */
  isRequired: boolean;
}

/** テンプレートフィールド作成リクエスト */
model CreateFieldRequest {
  /** ラベル */
  @minLength(1)
  label: string;

  /** 表示順序 */
  @minValue(1)
  order: int32;

  /** 必須フラグ */
  isRequired: boolean;
}

/** テンプレート作成リクエスト */
model CreateTemplateRequest {
  /** テンプレート名 */
  @minLength(1)
  @maxLength(100)
  name: string;

  /** 所有者ID */
  @format("uuid")
  ownerId: string;

  /** フィールド一覧 */
  fields: CreateFieldRequest[];

  /** Notionの置き場所となる親ページのURL（未設定なら連携しない） */
  notionParentPageUrl?: string;
}

/** テンプレート更新リクエスト */
model UpdateTemplateRequest {
  /** テンプレートID */
  id: string;

  /** テンプレート名 */
  @minLength(1)
  @maxLength(100)
  name: string;

  /** フィールド一覧 */
  fields: UpdateFieldRequest[];

  /** Notionの置き場所となる親ページのURL（未設定なら連携しない） */
  notionParentPageUrl?: string;
}

/** フィールド更新リクエスト */
model UpdateFieldRequest {
  /** フィールドID（既存フィールドの場合は必須） */
  id?: string;

  /** ラベル */
  @minLength(1)
  label: string;

  /** 表示順序 */
  @minValue(1)
  order: int32;

  /** 必須フラグ */
  isRequired: boolean;
}

/** テンプレートレスポンス */
model TemplateResponse {
  /** テンプレートID */
  id: string;

  /** テンプレート名 */
  name: string;

  /** 所有者ID */
  ownerId: string;

  /** 所有者情報 */
  owner: AccountSummary;

  /** フィールド一覧 */
  fields: Field[];

  /** 更新日時 */
  updatedAt: utcDateTime;

  /** 使用中フラグ */
  isUsed: boolean;

  /** Notionの置き場所となる親ページのURL */
  notionParentPageUrl?: string;
}
