import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "テンプレート一覧 | Mini Notion",
  description: "設計メモを構造化して残すミニノートアプリ",
};

export default function TemplateListLayout({
  children,
}: LayoutProps<"/templates">) {
  return <>{children}</>;
}
