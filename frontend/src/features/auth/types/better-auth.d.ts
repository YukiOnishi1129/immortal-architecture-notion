import type { Account } from "@/features/account/types";

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
