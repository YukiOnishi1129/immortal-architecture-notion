import type { Metadata } from "next";
import { getTemplateByIdQuery } from "@/external/handler/template/template.query.server";

interface LayoutProps {
  children: React.ReactNode;
  params: Promise<{ id: string }>;
}

export async function generateMetadata({
  params,
}: LayoutProps): Promise<Metadata> {
  const id = (await params).id;
  const template = await getTemplateByIdQuery({ id });

  return {
    title: template
      ? `${template.name}を編集 | Mini Notion`
      : "テンプレート編集 | Mini Notion",
    description: "設計メモを構造化して残すミニノートアプリ",
  };
}

export default function TemplateEditLayout({ children }: LayoutProps) {
  return <>{children}</>;
}
