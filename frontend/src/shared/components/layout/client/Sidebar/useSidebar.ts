"use client";

import { usePathname } from "next/navigation";
import { useCallback, useEffect, useState } from "react";

export function useSidebar() {
  const [isOpen, setIsOpen] = useState(false);
  const pathname = usePathname();

  const handleToggle = useCallback(() => {
    setIsOpen((prev) => !prev);
  }, []);

  // パス変更時にサイドバーを閉じる（モバイルのみ）
  // biome-ignore lint/correctness/useExhaustiveDependencies: pathname is needed to close sidebar on route change
  useEffect(() => {
    setIsOpen(false);
  }, [pathname]);

  return {
    isOpen,
    pathname,
    handleToggle,
  };
}
