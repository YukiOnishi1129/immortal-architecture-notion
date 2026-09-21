import { AuthenticatedLayoutWrapper } from "@/shared/components/layout/server/AuthenticatedLayoutWrapper";

export default function AuthenticatedPageLayout({
  children,
}: LayoutProps<"/">) {
  return <AuthenticatedLayoutWrapper>{children}</AuthenticatedLayoutWrapper>;
}
