import { redirectIfAuthenticatedServer } from "@/features/auth/servers/redirect.server";

type GuestLayoutWrapperProps = {
  children: React.ReactNode;
};

export const GuestLayoutWrapper = async ({
  children,
}: GuestLayoutWrapperProps) => {
  await redirectIfAuthenticatedServer();

  return <>{children}</>;
};
