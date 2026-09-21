"use client";

import { SidebarPresenter } from "./SidebarPresenter";
import { useSidebar } from "./useSidebar";

export function SidebarContainer() {
  const { isOpen, pathname, handleToggle } = useSidebar();

  return (
    <SidebarPresenter
      isOpen={isOpen}
      pathname={pathname}
      onToggle={handleToggle}
    />
  );
}
