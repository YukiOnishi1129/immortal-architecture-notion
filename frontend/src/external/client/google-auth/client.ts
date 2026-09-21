import { OAuth2Client } from "google-auth-library";

let oAuth2Client: OAuth2Client | null = null;

export const getGoogleOAuth2Client = () => {
  if (!oAuth2Client) {
    const baseUrl =
      process.env.NEXTAUTH_URL ??
      process.env.NEXT_PUBLIC_APP_URL ??
      "http://localhost:3000";

    oAuth2Client = new OAuth2Client(
      process.env.GOOGLE_CLIENT_ID,
      process.env.GOOGLE_CLIENT_SECRET,
      `${baseUrl.replace(/\/$/, "")}/api/auth/callback/google`,
    );
  }
  return oAuth2Client;
};
