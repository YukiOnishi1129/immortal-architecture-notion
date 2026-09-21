"use client";

import { useCallback } from "react";
import { signIn } from "@/features/auth/lib/better-auth-client";

export function useLoginClient() {
  const handleGoogleLogin = useCallback(async () => {
    await signIn.social({
      provider: "google",
      callbackURL: "/notes",
    });
  }, []);

  return {
    handleGoogleLogin,
  };
}
