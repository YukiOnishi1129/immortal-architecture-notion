import { GuestLayoutWrapper } from "@/shared/components/layout/server/GuestLayoutWrapper";

export default function GuestLayout({ children }: LayoutProps<"/">) {
  return <GuestLayoutWrapper>{children}</GuestLayoutWrapper>;
}
