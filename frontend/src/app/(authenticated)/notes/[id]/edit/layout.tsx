import type { Metadata } from "next";
import { getNoteByIdQuery } from "@/external/handler/note/note.query.server";

export async function generateMetadata({
  params,
}: LayoutProps<"/notes/[id]/edit">): Promise<Metadata> {
  const { id } = await params;
  const note = await getNoteByIdQuery({ id });

  return {
    title: note
      ? `${note.title}を編集 | Mini Notion`
      : "ノート編集 | Mini Notion",
    description: "設計メモを構造化して残すミニノートアプリ",
  };
}

export default function NoteEditRedirectLayout({
  children,
}: LayoutProps<"/notes/[id]/edit">) {
  return <>{children}</>;
}
