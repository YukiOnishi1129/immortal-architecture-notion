import { useQuery } from "@tanstack/react-query";
import { getNoteByIdQueryAction } from "@/external/handler/note/note.query.action";
import { noteKeys } from "@/features/note/queries/keys";

export function useNoteDetailQuery(noteId: string) {
  return useQuery({
    queryKey: noteKeys.detail(noteId),
    queryFn: () => getNoteByIdQueryAction({ id: noteId }),
  });
}
