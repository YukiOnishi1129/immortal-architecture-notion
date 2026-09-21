import { redirect } from "next/navigation";
import { getSessionServer } from "@/features/auth/servers/auth.server";

export default async function Home() {
  const session = await getSessionServer();
  if (session) {
    redirect("/notes");
  }

  redirect("/login");
}
