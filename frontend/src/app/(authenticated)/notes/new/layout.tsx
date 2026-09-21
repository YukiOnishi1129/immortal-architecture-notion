import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "新しいノート作成 | Mini Notion",
  description: "設計メモを構造化して残すミニノートアプリ",
};

export default function NewNoteLayout({ children }: LayoutProps<"/notes/new">) {
  return <>{children}</>;
}
