"use client";

import type { Route } from "next";
import { usePathname, useRouter, useSearchParams } from "next/navigation";
import { useCallback, useState, useTransition } from "react";

export function useTemplateListFilter() {
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();
  const [isPending, startTransition] = useTransition();

  const [searchQuery, setSearchQuery] = useState(searchParams.get("q") || "");
  const [onlyMyTemplates, setOnlyMyTemplates] = useState(
    searchParams.get("onlyMyTemplates") === "true",
  );

  const updateSearchParams = useCallback(
    (key: string, value: string) => {
      const params = new URLSearchParams(searchParams.toString());

      if (value) {
        params.set(key, value);
      } else {
        params.delete(key);
      }

      // Reset to page 1 when filters change
      if (key !== "page") {
        params.delete("page");
      }

      startTransition(() => {
        router.push(`${pathname}?${params.toString()}` as Route);
      });
    },
    [searchParams, pathname, router],
  );

  const handleSearch = useCallback(
    (e: React.FormEvent) => {
      e.preventDefault();
      updateSearchParams("q", searchQuery);
    },
    [searchQuery, updateSearchParams],
  );

  const handleOnlyMyTemplatesChange = useCallback(
    (checked: boolean) => {
      setOnlyMyTemplates(checked);
      updateSearchParams("onlyMyTemplates", checked ? "true" : "");
    },
    [updateSearchParams],
  );

  return {
    searchQuery,
    isPending,
    onlyMyTemplates,
    setSearchQuery,
    handleSearch,
    handleOnlyMyTemplatesChange,
  };
}
