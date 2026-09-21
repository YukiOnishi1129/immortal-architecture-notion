"use client";

import type { Route } from "next";
import { usePathname, useRouter, useSearchParams } from "next/navigation";
import { useCallback, useState, useTransition } from "react";
import type { NoteFilters } from "@/features/note/types";

export function useNoteListFilter(initialFilters?: NoteFilters) {
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();
  const [isPending, startTransition] = useTransition();

  const [searchQuery, setSearchQuery] = useState(
    searchParams.get("q") || initialFilters?.q || "",
  );
  const statusFilter =
    searchParams.get("status") || initialFilters?.status || "all";

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

  const handleStatusChange = useCallback(
    (value: string) => {
      updateSearchParams("status", value === "all" ? "" : value);
    },
    [updateSearchParams],
  );

  return {
    searchQuery,
    statusFilter,
    isPending,
    setSearchQuery,
    handleSearch,
    handleStatusChange,
  };
}
