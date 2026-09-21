"use server";

import { withAuth } from "@/features/auth/servers/auth.guard";
import type {
  CreateOrGetAccountRequest,
  CreateOrGetAccountResponse,
  UpdateAccountByIdRequest,
  UpdateAccountResponse,
} from "../../dto/account.dto";
import {
  createOrGetAccountCommand,
  updateAccountCommand,
} from "./account.command.server";

// NOTE: createOrGetAccountCommandはOAuth認証時に呼ばれるため、withAuthを適用しない
export async function createOrGetAccountCommandAction(
  request: CreateOrGetAccountRequest,
): Promise<CreateOrGetAccountResponse> {
  return createOrGetAccountCommand(request);
}

export async function updateAccountCommandAction(
  request: UpdateAccountByIdRequest,
): Promise<UpdateAccountResponse> {
  return withAuth(({ accountId }) => updateAccountCommand(request, accountId));
}
