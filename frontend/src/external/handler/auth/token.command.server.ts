import "server-only";

import {
  type RefreshGoogleTokensRequest,
  RefreshGoogleTokensRequestSchema,
  type RefreshGoogleTokensResponse,
  RefreshGoogleTokensResponseSchema,
} from "@/external/dto/auth.dto";
import { TokenVerificationService } from "@/external/service/auth/token-verification.service";

const tokenVerificationService = new TokenVerificationService();

export async function refreshGoogleTokenCommand(
  request: RefreshGoogleTokensRequest,
): Promise<RefreshGoogleTokensResponse> {
  const { refreshToken } = RefreshGoogleTokensRequestSchema.parse(request);
  const refreshed = await tokenVerificationService.refreshTokens(refreshToken);

  return RefreshGoogleTokensResponseSchema.parse({
    accessToken: refreshed.accessToken ?? null,
    idToken: refreshed.idToken ?? null,
    accessTokenExpires: refreshed.expiryDate ?? null,
  });
}
