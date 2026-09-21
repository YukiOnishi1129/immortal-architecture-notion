import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "ログイン | Mini Notion",
  description: "Mini Notionにログインして設計メモを管理しましょう",
};

export default function LoginLayout({ children }: LayoutProps<"/login">) {
  return <>{children}</>;
}
