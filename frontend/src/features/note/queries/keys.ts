import type { NoteFilters } from "../types";

export const noteKeys = {
  all: ["notes"] as const,
  lists: () => [...noteKeys.all, "list"] as const,
  list: (filters: NoteFilters) => [...noteKeys.lists(), filters] as const,
  myLists: () => [...noteKeys.all, "myList"] as const,
  myList: (filters: NoteFilters) => [...noteKeys.myLists(), filters] as const,
  details: () => [...noteKeys.all, "detail"] as const,
  detail: (id: string) => [...noteKeys.details(), id] as const,
};
