import { z } from "zod";

export const RefreshGoogleTokensRequestSchema = z.object({
  refreshToken: z.string().min(1),
});

export type RefreshGoogleTokensRequest = z.infer<
  typeof RefreshGoogleTokensRequestSchema
>;

export const RefreshGoogleTokensResponseSchema = z.object({
  accessToken: z.string().nullable(),
  idToken: z.string().nullable(),
  accessTokenExpires: z.number().int().nullable(),
});

export type RefreshGoogleTokensResponse = z.infer<
  typeof RefreshGoogleTokensResponseSchema
>;
