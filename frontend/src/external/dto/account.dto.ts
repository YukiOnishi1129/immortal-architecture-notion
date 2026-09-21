import { z } from "zod";

// Query request schemas
export const GetAccountByIdRequestSchema = z.object({
  id: z.uuid(),
});

export const GetAccountByEmailRequestSchema = z.object({
  email: z.email(),
});

// Command request schemas
export const CreateAccountRequestSchema = z.object({
  email: z.email(),
  name: z.string(),
  provider: z.string(),
  providerAccountId: z.string(),
  thumbnail: z.string().optional(),
});

export const CreateOrGetAccountRequestSchema = CreateAccountRequestSchema;

export const UpdateAccountRequestSchema = z.object({
  firstName: z.string().min(1).optional(),
  lastName: z.string().min(1).optional(),
  thumbnail: z.string().nullable().optional(),
});

export const UpdateAccountByIdRequestSchema = z.object({
  id: z.uuid(),
  firstName: z.string().min(1).optional(),
  lastName: z.string().min(1).optional(),
  thumbnail: z.string().nullable().optional(),
});

// Response schemas
export const AccountResponseSchema = z.object({
  id: z.uuid(),
  email: z.email(),
  firstName: z.string(),
  lastName: z.string(),
  fullName: z.string(),
  thumbnail: z.string().nullable(),
  lastLoginAt: z.iso.datetime().nullable(),
  createdAt: z.iso.datetime(),
  updatedAt: z.iso.datetime(),
});

// Type exports
export type GetAccountByIdRequest = z.infer<typeof GetAccountByIdRequestSchema>;
export type GetAccountByEmailRequest = z.infer<
  typeof GetAccountByEmailRequestSchema
>;
export type CreateAccountRequest = z.infer<typeof CreateAccountRequestSchema>;
export type CreateOrGetAccountRequest = z.infer<
  typeof CreateOrGetAccountRequestSchema
>;
export type UpdateAccountRequest = z.infer<typeof UpdateAccountRequestSchema>;
export type UpdateAccountByIdRequest = z.infer<
  typeof UpdateAccountByIdRequestSchema
>;
export type AccountResponse = z.infer<typeof AccountResponseSchema>;
export type CreateOrGetAccountResponse = AccountResponse;
export type UpdateAccountResponse = AccountResponse;
