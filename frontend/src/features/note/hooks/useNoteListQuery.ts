import { useQuery } from "@tanstack/react-query";
import { listNoteQueryAction } from "@/external/handler/note/note.query.action";
import { noteKeys } from "@/features/note/queries/keys";
import type { NoteFilters } from "@/features/note/types";

export function useNoteListQuery(filters: NoteFilters) {
  return useQuery({
    queryKey: noteKeys.list(filters),
    queryFn: () => listNoteQueryAction(filters),
  });
}
