import type { Metadata } from "next";
import { getNoteByIdQuery } from "@/external/handler/note/note.query.server";

interface LayoutProps {
  children: React.ReactNode;
  params: Promise<{ id: string }>;
}

export async function generateMetadata({
  params,
}: LayoutProps): Promise<Metadata> {
  const id = (await params).id;
  const note = await getNoteByIdQuery({ id });

  return {
    title: note
      ? `${note.title}を編集 | Mini Notion`
      : "ノート編集 | Mini Notion",
    description: "設計メモを構造化して残すミニノートアプリ",
  };
}

export default function MyNoteEditLayout({ children }: LayoutProps) {
  return <>{children}</>;
}
