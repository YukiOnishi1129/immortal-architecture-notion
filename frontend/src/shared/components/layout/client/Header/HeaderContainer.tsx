"use client";

import { HeaderPresenter } from "./HeaderPresenter";
import { useHeader } from "./useHeader";

export function HeaderContainer() {
  const { userName, userEmail, userImage, handleSignOut } = useHeader();

  return (
    <HeaderPresenter
      userName={userName}
      userEmail={userEmail}
      userImage={userImage}
      onSignOut={handleSignOut}
    />
  );
}
