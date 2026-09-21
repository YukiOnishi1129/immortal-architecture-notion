import { dehydrate, HydrationBoundary } from "@tanstack/react-query";
import { notFound, redirect } from "next/navigation";
import { getNoteByIdQuery } from "@/external/handler/note/note.query.server";
import { getSessionServer } from "@/features/auth/servers/auth.server";
import { NoteDetail } from "@/features/note/components/client/NoteDetail";
import { noteKeys } from "@/features/note/queries/keys";
import { getQueryClient } from "@/shared/lib/query-client";

type MyNoteDetailPageTemplateProps = {
  noteId: string;
};

export async function MyNoteDetailPageTemplate({
  noteId,
}: MyNoteDetailPageTemplateProps) {
  const [session, note] = await Promise.all([
    getSessionServer(),
    getNoteByIdQuery({ id: noteId }),
  ]);

  if (!note) {
    notFound();
  }

  // Check if the current user is the owner
  const isOwner = !!session?.account?.id && session.account.id === note.ownerId;

  // マイノート画面なので、自分のノートでなければリダイレクト
  if (!isOwner) {
    redirect("/my-notes");
  }

  const queryClient = getQueryClient();
  await queryClient.prefetchQuery({
    queryKey: noteKeys.detail(noteId),
    queryFn: () => note,
  });

  return (
    <HydrationBoundary state={dehydrate(queryClient)}>
      <div className="container mx-auto py-6">
        <NoteDetail noteId={noteId} isOwner={true} backTo="/my-notes" />
      </div>
    </HydrationBoundary>
  );
}
