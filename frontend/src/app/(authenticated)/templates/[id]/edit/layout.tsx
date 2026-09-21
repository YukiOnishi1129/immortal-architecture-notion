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

  // If template is used, don't generate edit metadata
  if (template?.isUsed) {
    return {
      title: `${template.name} | Mini Notion`,
      description: "設計メモを構造化して残すミニノートアプリ",
    };
  }

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
