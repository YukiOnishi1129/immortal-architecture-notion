import "server-only";

import {
  type CreateOrGetAccountRequest,
  CreateOrGetAccountRequestSchema,
  type CreateOrGetAccountResponse,
  type UpdateAccountByIdRequest,
  UpdateAccountByIdRequestSchema,
  type UpdateAccountResponse,
} from "../../dto/account.dto";
import { accountService } from "../../service/account/account.service";

// NOTE: createOrGetAccountCommandはOAuth認証時に呼ばれるため、withAuthを適用しない
export async function createOrGetAccountCommand(
  request: CreateOrGetAccountRequest,
): Promise<CreateOrGetAccountResponse> {
  const validated = CreateOrGetAccountRequestSchema.parse(request);
  return accountService.createOrGet(validated);
}

// NOTE: 認証チェック（withAuth）は .action.ts で行う
export async function updateAccountCommand(
  request: UpdateAccountByIdRequest,
  accountId: string,
): Promise<UpdateAccountResponse> {
  const validated = UpdateAccountByIdRequestSchema.parse(request);

  // Check if the user is updating their own account
  if (accountId !== validated.id) {
    throw new Error("Forbidden: Can only update your own account");
  }

  const { id, ...updateData } = validated;
  return accountService.update(id, updateData);
}
