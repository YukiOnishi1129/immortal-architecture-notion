import "server-only";

import { redirect } from "next/navigation";

import type { Account } from "@/features/account/types";
import { getSessionServer } from "@/features/auth/servers/auth.server";

// 認証済みセッションの型（accountが必ず存在する）
export type AuthenticatedSession = {
  account: Account;
  user: {
    id: string;
    email: string;
    name: string;
    image?: string;
  };
  session: {
    id: string;
    userId: string;
    expiresAt: Date;
  };
};

export const requireAuthServer = async () => {
  const session = await getSessionServer();
  if (!session?.account || session.error) {
    redirect("/login");
  }
};

export const getAuthenticatedSessionServer =
  async (): Promise<AuthenticatedSession> => {
    const session = await getSessionServer();
    if (!session?.account || session.error) {
      redirect("/login");
    }
    // redirectはneverを返すため、ここに到達した場合はaccountが存在する
    return session as AuthenticatedSession;
  };

export const redirectIfAuthenticatedServer = async () => {
  const session = await getSessionServer();
  if (session?.account && !session.error) {
    redirect("/notes");
  }
};
