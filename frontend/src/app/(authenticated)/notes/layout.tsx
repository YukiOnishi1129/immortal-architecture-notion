import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "ノート一覧 | Mini Notion",
  description: "設計メモを構造化して残すミニノートアプリ",
};

export default function NotesLayout({ children }: LayoutProps<"/notes">) {
  return <>{children}</>;
}
