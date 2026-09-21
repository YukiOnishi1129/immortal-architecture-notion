import { dehydrate, HydrationBoundary } from "@tanstack/react-query";
import { notFound } from "next/navigation";
import { getNoteByIdQuery } from "@/external/handler/note/note.query.server";
import { getSessionServer } from "@/features/auth/servers/auth.server";
import { NoteDetail } from "@/features/note/components/client/NoteDetail";
import { noteKeys } from "@/features/note/queries/keys";
import { getQueryClient } from "@/shared/lib/query-client";

type NoteDetailPageTemplateProps = {
  noteId: string;
};

export async function NoteDetailPageTemplate({
  noteId,
}: NoteDetailPageTemplateProps) {
  const [session, note] = await Promise.all([
    getSessionServer(),
    getNoteByIdQuery({ id: noteId }),
  ]);

  if (!note) {
    notFound();
  }

  // Check if the current user is the owner
  const isOwner = !!session?.account?.id && session.account.id === note.ownerId;

  const queryClient = getQueryClient();
  await queryClient.prefetchQuery({
    queryKey: noteKeys.detail(noteId),
    queryFn: () => note,
  });

  return (
    <HydrationBoundary state={dehydrate(queryClient)}>
      <div className="container mx-auto py-6">
        <NoteDetail noteId={noteId} isOwner={isOwner} />
      </div>
    </HydrationBoundary>
  );
}
