import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "ノート新規作成 | Mini Notion",
  description: "設計メモを構造化して残すミニノートアプリ",
};

export default function MyNoteNewLayout({
  children,
}: LayoutProps<"/my-notes/new">) {
  return <>{children}</>;
}
