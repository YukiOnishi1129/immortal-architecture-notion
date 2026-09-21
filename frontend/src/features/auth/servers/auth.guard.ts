import "server-only";

import { getAuthenticatedSessionServer } from "@/features/auth/servers/redirect.server";

export async function withAuth<T>(
  handler: (ctx: { accountId: string }) => Promise<T>,
): Promise<T> {
  const session = await getAuthenticatedSessionServer();
  // getAuthenticatedSessionServerはaccountが存在することを保証する
  return handler({ accountId: session.account.id });
}
