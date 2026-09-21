import "server-only";
import type { Session } from "better-auth";
import { headers } from "next/headers";
import { auth } from "@/features/auth/lib/better-auth";

export async function getSessionServer(): Promise<Session | null> {
  return await auth.api.getSession({
    headers: await headers(),
  });
}
