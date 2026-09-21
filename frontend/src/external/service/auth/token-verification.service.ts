import { getGoogleOAuth2Client } from "@/external/client/google-auth/client";

export interface TokenPayload {
  userId: string;
  email?: string | null;
  emailVerified?: boolean;
  name?: string;
  picture?: string;
  isValid: boolean;
}

export class TokenVerificationService {
  /**
   * Verify Google ID token and return payload
   */
  async verifyIdToken(idToken: string): Promise<TokenPayload> {
    try {
      const ticket = await getGoogleOAuth2Client().verifyIdToken({
        idToken: idToken,
        audience: process.env.GOOGLE_CLIENT_ID || "",
      });

      const payload = ticket.getPayload();
      if (!payload) {
        throw new Error("Invalid token payload");
      }

      return {
        userId: payload.sub,
        email: payload.email,
        emailVerified: payload.email_verified,
        name: payload.name,
        picture: payload.picture,
        isValid: true,
      };
    } catch (error) {
      console.error("Token verification failed:", error);
      throw new Error("Invalid ID token");
    }
  }

  /**
   * Use refresh token to get new tokens
   */
  async refreshTokens(refreshToken: string) {
    try {
      getGoogleOAuth2Client().setCredentials({ refresh_token: refreshToken });
      const { credentials } =
        await getGoogleOAuth2Client().refreshAccessToken();

      return {
        accessToken: credentials.access_token,
        idToken: credentials.id_token,
        expiryDate: credentials.expiry_date,
      };
    } catch (error) {
      console.error("Token refresh failed:", error);
      throw new Error("Failed to refresh tokens");
    }
  }
}
