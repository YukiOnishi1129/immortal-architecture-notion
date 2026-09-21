import { accountsApiClient } from "@/external/client/api/config";
import type { AccountsApi } from "@/external/client/api/generated/apis/AccountsApi";
import type { ModelsAccountResponse } from "@/external/client/api/generated/models/ModelsAccountResponse";
import type {
  AccountResponse,
  CreateOrGetAccountRequest,
  UpdateAccountRequest,
} from "@/external/dto/account.dto";
import { AccountResponseSchema } from "@/external/dto/account.dto";
import { isNotFoundError } from "../http-error";

function formatDate(value?: Date | null): string | null {
  return value ? value.toISOString() : null;
}

function toAccountResponse(model: ModelsAccountResponse): AccountResponse {
  return AccountResponseSchema.parse({
    id: model.id,
    email: model.email,
    firstName: model.firstName,
    lastName: model.lastName,
    fullName: model.fullName,
    thumbnail: model.thumbnail ?? null,
    lastLoginAt: formatDate(model.lastLoginAt),
    createdAt: formatDate(model.createdAt),
    updatedAt: formatDate(model.updatedAt),
  });
}

export class AccountService {
  constructor(private readonly api: AccountsApi) {}

  async createOrGet(
    input: CreateOrGetAccountRequest,
  ): Promise<AccountResponse> {
    const account = await this.api.accountsCreateOrGetAccount({
      modelsCreateOrGetAccountRequest: input,
    });
    return toAccountResponse(account);
  }

  async getAccountById(id: string): Promise<AccountResponse | null> {
    try {
      const account = await this.api.accountsGetAccountById({
        accountId: id,
      });
      return toAccountResponse(account);
    } catch (error) {
      if (isNotFoundError(error)) {
        return null;
      }
      throw error;
    }
  }

  async getCurrentAccountByEmail(
    email: string,
  ): Promise<AccountResponse | null> {
    try {
      const account = await this.api.accountsGetAccountByEmail({
        email: email,
      });
      return toAccountResponse(account);
    } catch (error) {
      if (isNotFoundError(error)) {
        return null;
      }
      throw error;
    }
  }

  async update(
    _id: string,
    _input: UpdateAccountRequest,
  ): Promise<AccountResponse> {
    throw new Error(
      "Account update API is not implemented on the backend yet.",
    );
  }
}

export const accountService = new AccountService(accountsApiClient);
