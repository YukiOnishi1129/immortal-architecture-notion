import "server-only";

import { getSessionServer } from "@/features/auth/servers/auth.server";

export async function checkAuthAndRefreshServer(): Promise<boolean> {
  const account = await getSessionServer();
  return Boolean(account);
}
