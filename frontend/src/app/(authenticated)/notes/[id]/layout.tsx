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
    title: note ? `${note.title} | Mini Notion` : "ノート詳細 | Mini Notion",
    description: "設計メモを構造化して残すミニノートアプリ",
  };
}

export default function NoteDetailLayout({ children }: LayoutProps) {
  return <>{children}</>;
}
