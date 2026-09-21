# 認証システム実装ガイド

## 概要

Better Auth を使用した認証システム。Google OAuth 2.0 による認証と、stateless セッション管理を採用。

## 技術スタック

| 項目             | 技術                                  |
| ---------------- | ------------------------------------- |
| 認証ライブラリ   | Better Auth                           |
| OAuth プロバイダ | Google OAuth 2.0                      |
| セッション管理   | Stateless（Cookie ベース）            |
| ユーザーデータ   | PostgreSQL（accounts テーブル）       |
| キャッシュ       | Next.js unstable_cache + Cookie Cache |
| トークン検証     | google-auth-library                   |

## アーキテクチャ

```
┌─────────────────────────────────────────────────────────────────┐
│                        認証フロー                                │
├─────────────────────────────────────────────────────────────────┤
│                                                                 │
│  ユーザー                                                        │
│     │                                                           │
│     ▼                                                           │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────┐         │
│  │ ログイン    │───▶│ Google     │───▶│ コールバック │         │
│  │ ボタン     │    │ OAuth 2.0  │    │ /api/auth/* │         │
│  └─────────────┘    └─────────────┘    └──────┬──────┘         │
│                                               │                 │
│                                               ▼                 │
│                                        ┌─────────────┐         │
│                                        │ Better Auth │         │
│                                        │ (stateless) │         │
│                                        └──────┬──────┘         │
│                                               │                 │
│                      ┌────────────────────────┼────────────┐   │
│                      │                        │            │   │
│                      ▼                        ▼            ▼   │
│               ┌─────────────┐         ┌─────────────┐  ┌─────┐│
│               │ onSuccess   │         │customSession│  │Cookie││
│               │ (初回のみ)  │         │ (毎回実行)  │  │保存 ││
│               └──────┬──────┘         └──────┬──────┘  └─────┘│
│                      │                       │                 │
│                      ▼                       ▼                 │
│               ┌─────────────┐         ┌─────────────┐         │
│               │ handler     │         │unstable_cache│         │
│               │ (command)   │         │ (5分キャッシュ)│        │
│               └──────┬──────┘         └──────┬──────┘         │
│                      │                       │                 │
│                      ▼                       ▼                 │
│               ┌─────────────────────────────────────┐         │
│               │         accounts テーブル           │         │
│               │    (id, email, provider, etc.)     │         │
│               └─────────────────────────────────────┘         │
│                                                                 │
└─────────────────────────────────────────────────────────────────┘
```

## Stateless セッションとは

従来の認証（Stateful）との違い：

| 項目             | Stateful                | Stateless（採用）        |
| ---------------- | ----------------------- | ------------------------ |
| セッション保存   | DB（sessions テーブル） | Cookie（署名付き）       |
| DB アクセス      | 毎リクエスト            | 不要                     |
| スケーラビリティ | サーバー間で共有必要    | 各サーバーで独立処理可能 |
| ログアウト       | DB から削除             | Cookie 削除のみ          |

### メリット

- DB へのセッション問い合わせが不要（高速）
- 水平スケーリングが容易
- sessions テーブルが不要

### デメリット

- サーバー側から強制ログアウトが難しい
- セッションデータの即時更新が難しい

## 処理フロー詳細

### 1. ログインボタンクリック

```typescript
// features/auth/components/client/LoginPageClient/useLoginClient.ts
export function useLoginClient() {
  const handleGoogleLogin = useCallback(async () => {
    await signIn.social({
      provider: "google",
      callbackURL: "/notes", // ログイン後のリダイレクト先
    });
  }, []);

  return { handleGoogleLogin };
}
```

### 2. Google OAuth 認証

1. ユーザーが Google のログイン画面にリダイレクト
2. Google アカウントでログイン
3. 認可コードがコールバック URL に返される
4. Better Auth がアクセストークンを取得

### 3. Better Auth コールバック処理

```typescript
// app/api/auth/[...all]/route.ts
import { toNextJsHandler } from "better-auth/next-js";
import { auth } from "@/features/auth/lib/better-auth";

export const { GET, POST } = toNextJsHandler(auth);
```

### 4. onSuccess コールバック（初回ログイン時）

```typescript
// features/auth/lib/better-auth.ts
socialProviders: {
  google: {
    clientId: process.env.GOOGLE_CLIENT_ID || "",
    clientSecret: process.env.GOOGLE_CLIENT_SECRET || "",
    async onSuccess(ctx) {
      // accounts テーブルにユーザー情報を保存
      await createOrGetAccountCommand({
        email: ctx.user.email,
        firstName: ctx.user.name?.split(" ")[0] || ctx.user.email,
        lastName: ctx.user.name?.split(" ").slice(1).join(" ") || "",
        provider: "google",
        providerAccountId: ctx.user.id,
        thumbnail: ctx.user.image || undefined,
      });
    }
  }
}
```

### 5. customSession プラグイン（毎回実行）

```typescript
// features/auth/lib/better-auth.ts
plugins: [
  customSession(async ({ user, session }) => {
    // unstable_cache でDBアクセスをキャッシュ（5分）
    let account = await getCachedAccount(user.email);

    // アカウントが存在しない場合は作成（フォールバック）
    if (!account) {
      await createOrGetAccountCommand({
        email: user.email,
        firstName: user.name?.split(" ")[0] || user.email,
        lastName: user.name?.split(" ").slice(1).join(" ") || "",
        provider: "google",
        providerAccountId: user.id,
        thumbnail: user.image || undefined,
      });
      account = await getAccountByEmailQuery(user.email);
    }

    return { user, session, account };
  }),
];
```

### 6. セッション取得

```typescript
// features/auth/servers/auth.server.ts
export async function getSessionServer(): Promise<Session | null> {
  return await auth.api.getSession({ headers: await headers() });
}
```

## キャッシュ戦略

### サーバーサイド（unstable_cache）

```typescript
const getCachedAccount = unstable_cache(
  async (email: string): Promise<Account | null> => {
    return await getAccountByEmailQuery(email);
  },
  ["account-by-email"],
  {
    revalidate: 300, // 5分間キャッシュ
    tags: ["account"],
  }
);
```

| 設定       | 値          | 説明                       |
| ---------- | ----------- | -------------------------- |
| revalidate | 300 秒      | キャッシュの有効期間       |
| tags       | ["account"] | revalidateTag で無効化可能 |

### クライアントサイド（Cookie Cache）

```typescript
session: {
  cookieCache: {
    enabled: true,
    maxAge: 5 * 60, // 5分間
  },
}
```

セッション情報を Cookie にキャッシュし、毎回の DB/API 呼び出しを削減。

## 型定義（Module Augmentation）

```typescript
// features/auth/types/better-auth.d.ts
declare module "better-auth" {
  interface Session {
    account?: Account;
    error?: "RefreshTokenMissing" | "RefreshAccessTokenError";
  }

  interface User {
    id: string;
    account?: Account;
  }
}
```

これにより：

- `session.account` でアカウント情報にアクセス可能
- `session.error` でトークン関連エラーをハンドリング可能

## 認証ガード

### サーバーコンポーネント

```typescript
// features/auth/servers/redirect.server.ts

// 認証が必須 - セッションがなければログインページへリダイレクト
export const requireAuthServer = async () => {
  const session = await getSessionServer();
  if (!session?.account || session.error) {
    redirect("/login");
  }
};

// 認証済みセッションを取得 - なければリダイレクト
export const getAuthenticatedSessionServer = async () => {
  const session = await getSessionServer();
  if (!session?.account || session.error) {
    redirect("/login");
  }
  return session;
};

// 認証済みならリダイレクト - ログインページでの使用
export const redirectIfAuthenticatedServer = async () => {
  const session = await getSessionServer();
  if (session?.account && !session.error) {
    redirect("/notes");
  }
};
```

### 認証チェック（リダイレクトなし）

```typescript
// features/auth/servers/auth-check.server.ts
export async function checkAuthAndRefreshServer(): Promise<boolean> {
  const account = await getSessionServer();
  return Boolean(account);
}
```

### レイアウトラッパー

```typescript
// shared/components/layout/server/AuthenticatedLayoutWrapper/AuthenticatedLayoutWrapper.tsx
export async function AuthenticatedLayoutWrapper({ children }) {
  await requireAuthServer();

  return (
    <div className="flex h-screen bg-background">
      <Sidebar />
      <div className="flex flex-1 flex-col">
        <Header />
        <main className="flex-1 overflow-y-auto bg-gray-50 p-6">
          {children}
        </main>
      </div>
    </div>
  );
}
```

```typescript
// shared/components/layout/server/GuestLayoutWrapper/GuestLayoutWrapper.tsx
export const GuestLayoutWrapper = async ({ children }) => {
  await redirectIfAuthenticatedServer();
  return <>{children}</>;
};
```

### 使用例

```typescript
// app/(authenticated)/layout.tsx
export default function AuthenticatedPageLayout({ children }) {
  return <AuthenticatedLayoutWrapper>{children}</AuthenticatedLayoutWrapper>;
}

// app/(guest)/layout.tsx
export default function GuestPageLayout({ children }) {
  return <GuestLayoutWrapper>{children}</GuestLayoutWrapper>;
}
```

## 環境変数

```env
# API
API_BASE_URL=http://localhost:8080
NEXT_PUBLIC_APP_URL=http://localhost:3000

# Better Auth
BETTER_AUTH_URL=http://localhost:3000
BETTER_AUTH_SECRET=your-secret-key

# Google OAuth
GOOGLE_CLIENT_ID=your-client-id
GOOGLE_CLIENT_SECRET=your-client-secret
```

| 変数名                 | 説明                                  |
| ---------------------- | ------------------------------------- |
| `NEXT_PUBLIC_APP_URL`  | アプリケーションの公開 URL            |
| `BETTER_AUTH_URL`      | Better Auth の認証 URL                |
| `BETTER_AUTH_SECRET`   | セッション署名用のシークレットキー    |
| `GOOGLE_CLIENT_ID`     | Google OAuth クライアント ID          |
| `GOOGLE_CLIENT_SECRET` | Google OAuth クライアントシークレット |

## データベーススキーマ

マイグレーションは backend-clean 側で管理されています。

```sql
-- backend-clean/migrations/20250209000000_init_schema.up.sql
CREATE TABLE accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT NOT NULL UNIQUE,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    provider TEXT NOT NULL,
    provider_account_id TEXT NOT NULL,
    thumbnail TEXT,
    last_login_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT provider_account_unique UNIQUE (provider, provider_account_id)
);
```

### ユニーク制約

| 制約                    | カラム                          | 目的                     |
| ----------------------- | ------------------------------- | ------------------------ |
| email (UNIQUE)          | email                           | メールアドレスの重複防止 |
| provider_account_unique | (provider, provider_account_id) | 同一プロバイダの重複防止 |

## 重複アカウント処理

同じユーザーが再ログインした場合の処理：

```typescript
// external/repository/account.repository.ts
.onConflictDoUpdate({
  target: [accounts.provider, accounts.providerAccountId],
  set: {
    email: data.email,
    firstName: data.firstName,
    lastName: data.lastName,
    thumbnail: data.thumbnail,
    lastLoginAt: new Date(),
    updatedAt: new Date(),
  },
})
```

これにより：

- 新規ユーザー → INSERT
- 既存ユーザー → UPDATE（lastLoginAt 更新）

## ファイル構成

```
features/auth/
├── lib/
│   ├── better-auth.ts          # サーバー側設定（auth, customSession）
│   └── better-auth-client.ts   # クライアント側設定（authClient, signIn, signOut, useSession）
├── servers/
│   ├── auth.server.ts          # getSessionServer
│   ├── auth-check.server.ts    # checkAuthAndRefreshServer
│   └── redirect.server.ts      # requireAuthServer, getAuthenticatedSessionServer, redirectIfAuthenticatedServer
├── types/
│   └── better-auth.d.ts        # 型定義（Module Augmentation）
└── components/
    ├── client/
    │   └── LoginPageClient/    # ログインUI
    │       ├── useLoginClient.ts
    │       ├── LoginClientContainer.tsx
    │       └── LoginClientPresenter.tsx
    └── server/
        └── LoginPageTemplate/

features/account/
└── types/
    └── index.ts                # Account 型定義

shared/components/layout/
├── client/
│   └── Header/
│       ├── Header.tsx
│       └── useHeader.ts        # ログアウト処理、ユーザー情報取得
└── server/
    ├── AuthenticatedLayoutWrapper/  # 認証済みユーザー用レイアウト
    └── GuestLayoutWrapper/          # 未認証ユーザー用レイアウト

external/
├── handler/
│   ├── auth/
│   │   └── token.command.server.ts   # refreshGoogleTokenCommand
│   └── account/
│       ├── account.command.server.ts # createOrGetAccountCommand, updateAccountCommand
│       └── account.query.server.ts   # getCurrentAccountQuery, getAccountByEmailQuery
├── service/auth/
│   └── token-verification.service.ts # TokenVerificationService（トークン検証・リフレッシュ）
└── client/google-auth/
    └── client.ts                     # Google OAuth2 クライアント

app/
├── api/auth/[...all]/
│   └── route.ts              # Better Auth API ハンドラー
├── (authenticated)/          # 認証済みユーザー用ルート
│   ├── layout.tsx
│   ├── notes/
│   ├── my-notes/
│   └── templates/
└── (guest)/                  # 未認証ユーザー用ルート
    ├── layout.tsx
    └── login/
```

## ログアウト処理

```typescript
// shared/components/layout/client/Header/useHeader.ts
export function useHeader() {
  const { data: session } = useSession();
  const router = useRouter();
  const queryClient = useQueryClient();

  const handleSignOut = useCallback(async () => {
    await signOut(); // Better Auth のログアウト
    queryClient.clear(); // TanStack Query キャッシュクリア
    router.push("/login"); // ログインページへリダイレクト
  }, [router, queryClient]);

  return {
    userName: session?.user?.name,
    userEmail: session?.user?.email,
    userImage: session?.user?.image,
    handleSignOut,
  };
}
```

### ログアウト時の処理フロー

1. `signOut()` - Better Auth のセッションクッキーを削除
2. `queryClient.clear()` - TanStack Query のキャッシュを完全クリア
3. `router.push("/login")` - ログインページへリダイレクト

## クライアントフック

### useSession

Better Auth クライアントから提供される `useSession` フックを使用してセッション情報を取得。

```typescript
// features/auth/lib/better-auth-client.ts
export const authClient = createAuthClient({
  baseURL: process.env.NEXT_PUBLIC_APP_URL || "http://localhost:3000",
  plugins: [
    customSessionClient<typeof auth>(), // 型推論を有効化
  ],
});

export const { signIn, signOut, useSession } = authClient;
```

### 使用例

```typescript
const { data: session, isPending } = useSession();

if (isPending) return <Loading />;
if (!session) return <LoginButton />;

return <div>Welcome, {session.user.name}</div>;
```

## トークン管理

### Google トークン検証サービス

```typescript
// external/service/auth/token-verification.service.ts
export class TokenVerificationService {
  async verifyIdToken(idToken: string): Promise<TokenPayload> {
    const ticket = await getGoogleOAuth2Client().verifyIdToken({
      idToken,
      audience: process.env.GOOGLE_CLIENT_ID || "",
    });
    const payload = ticket.getPayload();
    return {
      userId: payload.sub,
      email: payload.email,
      emailVerified: payload.email_verified,
      name: payload.name,
      picture: payload.picture,
      isValid: true,
    };
  }

  async refreshTokens(refreshToken: string) {
    getGoogleOAuth2Client().setCredentials({
      refresh_token: refreshToken,
    });
    const { credentials } = await getGoogleOAuth2Client().refreshAccessToken();
    return {
      accessToken: credentials.access_token,
      idToken: credentials.id_token,
      expiryDate: credentials.expiry_date,
    };
  }
}
```

### Google OAuth2 クライアント

```typescript
// external/client/google-auth/client.ts
export const getGoogleOAuth2Client = () => {
  if (!oAuth2Client) {
    const baseUrl =
      process.env.NEXTAUTH_URL ??
      process.env.NEXT_PUBLIC_APP_URL ??
      "http://localhost:3000";

    oAuth2Client = new OAuth2Client(
      process.env.GOOGLE_CLIENT_ID,
      process.env.GOOGLE_CLIENT_SECRET,
      `${baseUrl}/api/auth/callback/google`
    );
  }
  return oAuth2Client;
};
```

## トラブルシューティング

### セッションが取得できない

1. Cookie が正しく設定されているか確認
2. `BETTER_AUTH_SECRET` が設定されているか確認
3. `BETTER_AUTH_URL` が正しいか確認

### アカウントが作成されない

1. データベース接続を確認
2. accounts テーブルが存在するか確認（`pnpm db:push`）
3. ユニーク制約違反がないか確認

### customSession でエラー

1. `getAccountByEmailQuery` が null を返していないか確認
2. unstable_cache のタグが正しいか確認
3. handler → service → repository の呼び出し順序を確認

### ログアウト後もデータが残る

1. `queryClient.clear()` が呼ばれているか確認
2. ブラウザの Cookie を手動で削除してテスト
3. `useSession` の `isPending` 状態を確認
