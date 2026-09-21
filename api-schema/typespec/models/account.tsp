import "@typespec/http";
import "@typespec/openapi3";

using TypeSpec.Http;

namespace MiniNotion.Models;

/** アカウント情報 */
model Account {
  /** アカウントID */
  id: string;

  /** メールアドレス */
  email: string;

  /** 名前 */
  firstName: string;

  /** 苗字 */
  lastName: string;

  /** フルネーム */
  fullName: string;

  /** プロフィール画像URL */
  thumbnail?: string;

  /** 最終ログイン日時 */
  lastLoginAt: utcDateTime;

  /** 作成日時 */
  createdAt: utcDateTime;

  /** 更新日時 */
  updatedAt: utcDateTime;
}

/** OAuth認証リクエスト */
model CreateOrGetAccountRequest {
  /** メールアドレス */
  email: string;

  /** 名前 */
  name: string;

  /** プロバイダー（例: google） */
  provider: string;

  /** プロバイダーのアカウントID */
  providerAccountId: string;

  /** プロフィール画像URL */
  thumbnail?: string;
}

/** アカウントレスポンス */
model AccountResponse {
  /** アカウントID */
  id: string;

  /** メールアドレス */
  email: string;

  /** 名前 */
  firstName: string;

  /** 苗字 */
  lastName: string;

  /** フルネーム */
  fullName: string;

  /** プロフィール画像URL */
  thumbnail?: string;

  /** 最終ログイン日時 */
  lastLoginAt: utcDateTime;

  /** 作成日時 */
  createdAt: utcDateTime;

  /** 更新日時 */
  updatedAt: utcDateTime;
}

/** 簡易アカウント情報（他のレスポンスに埋め込まれる） */
model AccountSummary {
  /** アカウントID */
  id: string;

  /** 名前 */
  firstName: string;

  /** 苗字 */
  lastName: string;

  /** プロフィール画像URL */
  thumbnail?: string;
}
