import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "テンプレート新規作成 | Mini Notion",
  description: "設計メモを構造化して残すミニノートアプリ",
};

export default function TemplateNewLayout({
  children,
}: LayoutProps<"/templates/new">) {
  return <>{children}</>;
}
