"use client";

import { LoginClientPresenter } from "./LoginClientPresenter";
import { useLoginClient } from "./useLoginClient";

export function LoginClientContainer() {
  const { handleGoogleLogin } = useLoginClient();

  return <LoginClientPresenter onGoogleLogin={handleGoogleLogin} />;
}
