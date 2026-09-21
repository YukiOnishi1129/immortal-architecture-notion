import "server-only";
import { betterAuth } from "better-auth";
import { customSession } from "better-auth/plugins";
import { unstable_cache, updateTag } from "next/cache";
import { createOrGetAccountCommand } from "@/external/handler/account/account.command.server";
import { getAccountByEmailQuery } from "@/external/handler/account/account.query.server";
import type { Account } from "@/features/account/types";

// customSessionは毎回実行されるため、Next.jsのunstable_cacheでキャッシング
// キャッシュ期間: 5分（セッションの有効期間と同等）
// NOTE: unstable_cacheは関数の引数も自動的にキャッシュキーに含まれる
const getCachedAccount = unstable_cache(
  async (email: string): Promise<Account | null> => {
    return await getAccountByEmailQuery({ email });
  },
  ["account-by-email"], // 引数emailは自動的にキャッシュキーに含まれる
  {
    revalidate: 300, // 5分間キャッシュ
    tags: ["account"],
  },
);

export const auth = betterAuth({
  // データベース設定なし = stateless mode
  session: {
    cookieCache: {
      enabled: true,
      maxAge: 5 * 60, // 5分間キャッシュ
    },
  },
  socialProviders: {
    google: {
      clientId: process.env.GOOGLE_CLIENT_ID || "",
      clientSecret: process.env.GOOGLE_CLIENT_SECRET || "",
      // OAuthコールバック後の処理を追加
      async onSuccess(ctx: {
        user: {
          id: string;
          email: string;
          name: string;
          image?: string;
        };
      }) {
        try {
          await createOrGetAccountCommand({
            email: ctx.user.email,
            name: ctx.user.name || ctx.user.email,
            provider: "google",
            providerAccountId: ctx.user.id,
            thumbnail: ctx.user.image || undefined,
          });
          // アカウント作成/更新後にキャッシュを無効化して、customSessionで最新データを取得
          updateTag("account");
        } catch (error) {
          console.error(
            "[better-auth] Failed to save account in onSuccess:",
            error,
          );
          throw error; // エラーを再スローして認証を失敗させる
        }
      },
    },
  },
  // ベースURLの設定
  baseURL: process.env.NEXTAUTH_URL || "http://localhost:3000",
  // セッション設定 - statelessモードではデフォルトでクッキーベース
  plugins: [
    // カスタムセッション: DBからaccount情報を取得してセッションに追加
    // NOTE: customSessionは毎回実行されるため、unstable_cacheでキャッシングしている
    customSession(async ({ user, session }) => {
      let account = await getCachedAccount(user.email);

      // accountが存在する場合は、そのまま返す
      if (account) {
        return { user, session, account };
      }

      // accountが存在しない場合は、DB保存を試みる（初回ログイン時）
      try {
        account = await createOrGetAccountCommand({
          email: user.email,
          name: user.name || user.email,
          provider: "google",
          providerAccountId: user.id,
          thumbnail: user.image || undefined,
        });

        return { user, session, account };
      } catch (error) {
        console.error("[better-auth] Failed to create account:", error);
        throw new Error("Failed to create account");
      }
    }),
  ],
});
