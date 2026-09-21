import { dehydrate, HydrationBoundary } from "@tanstack/react-query";
import { listMyNoteQuery } from "@/external/handler/note/note.query.server";
import { getAuthenticatedSessionServer } from "@/features/auth/servers/redirect.server";
import { MyNoteList } from "@/features/note/components/client/MyNoteList";
import { noteKeys } from "@/features/note/queries/keys";
import type { NoteStatus } from "@/features/note/types";
import { getQueryClient } from "@/shared/lib/query-client";

type MyNoteListPageTemplateProps = {
  status?: NoteStatus;
  q?: string;
  page?: number;
};

export async function MyNoteListPageTemplate(
  props: MyNoteListPageTemplateProps = {},
) {
  const session = await getAuthenticatedSessionServer();
  const queryClient = getQueryClient();

  const filters = {
    status: props.status,
    q: props.q,
    page: props.page,
  };

  await queryClient.prefetchQuery({
    queryKey: noteKeys.myList(filters),
    queryFn: () => listMyNoteQuery(filters, session.account.id),
  });

  return (
    <HydrationBoundary state={dehydrate(queryClient)}>
      <MyNoteList initialFilters={filters} />
    </HydrationBoundary>
  );
}
