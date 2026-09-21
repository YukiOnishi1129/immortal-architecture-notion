"use client";

import type { Route } from "next";
import { useSearchParams } from "next/navigation";
import type { BreadcrumbItem } from "@/shared/components/ui/breadcrumb";

export function useNoteNavigation() {
  const searchParams = useSearchParams();
  const from = searchParams.get("from");

  const isFromMyNotes = from === "my-notes";

  const getListPath = (): Route => {
    return isFromMyNotes ? "/my-notes" : "/notes";
  };

  const getListLabel = (): string => {
    return isFromMyNotes ? "マイノート" : "みんなのノート";
  };

  const getBreadcrumbItems = (
    noteTitle?: string,
    isEdit = false,
  ): BreadcrumbItem[] => {
    const items: BreadcrumbItem[] = [
      {
        label: getListLabel(),
        href: getListPath(),
      },
    ];

    if (noteTitle) {
      if (isEdit) {
        items.push({
          label: noteTitle,
          href: `/notes/${noteTitle}?from=${from || ""}` as Route,
        });
        items.push({
          label: "編集",
        });
      } else {
        items.push({
          label: noteTitle,
        });
      }
    }

    return items;
  };

  return {
    isFromMyNotes,
    getListPath,
    getListLabel,
    getBreadcrumbItems,
  };
}
