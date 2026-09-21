import { dehydrate, HydrationBoundary } from "@tanstack/react-query";
import { redirect } from "next/navigation";
import { getNoteByIdQuery } from "@/external/handler/note/note.query.server";
import { getSessionServer } from "@/features/auth/servers/auth.server";
import { NoteEditForm } from "@/features/note/components/client/NoteEditForm";
import { noteKeys } from "@/features/note/queries/keys";
import { getQueryClient } from "@/shared/lib/query-client";

type MyNoteEditPageTemplateProps = {
  noteId: string;
};

export async function MyNoteEditPageTemplate({
  noteId,
}: MyNoteEditPageTemplateProps) {
  const queryClient = getQueryClient();

  const [session, note] = await Promise.all([
    getSessionServer(),
    getNoteByIdQuery({ id: noteId }),
  ]);

  if (!note) {
    redirect("/my-notes");
  }

  // Check if the current user is the owner
  if (!session?.account?.id || session.account.id !== note.ownerId) {
    // 他人のノートは編集できない
    redirect("/my-notes");
  }

  await queryClient.prefetchQuery({
    queryKey: noteKeys.detail(noteId),
    queryFn: () => note,
  });

  return (
    <HydrationBoundary state={dehydrate(queryClient)}>
      <NoteEditForm noteId={noteId} backTo="/my-notes" />
    </HydrationBoundary>
  );
}
