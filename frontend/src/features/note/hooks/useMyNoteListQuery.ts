import { useQuery } from "@tanstack/react-query";
import { listMyNoteQueryAction } from "@/external/handler/note/note.query.action";
import { noteKeys } from "@/features/note/queries/keys";
import type { NoteFilters } from "@/features/note/types";

export function useMyNoteListQuery(filters: NoteFilters) {
  return useQuery({
    queryKey: noteKeys.myList(filters),
    queryFn: () => listMyNoteQueryAction(filters),
  });
}
