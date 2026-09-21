import type { Route } from "next";
import { redirect } from "next/navigation";

export default async function NoteEditRedirectPage({
  params,
  searchParams,
}: PageProps<"/notes/[id]/edit">) {
  const { id } = await params;
  const from = (await searchParams)?.from as string | undefined;

  // 編集画面は /my-notes 側にのみ存在するため、必ずそちらへ誘導する。
  const query = from ? `?from=${from}` : "";
  redirect(`/my-notes/${id}/edit${query}` as Route);
}
