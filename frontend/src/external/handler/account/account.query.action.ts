"use server";

import type { GetAccountByIdRequest } from "../../dto/account.dto";
import {
  getAccountByIdQuery,
  getCurrentAccountQuery,
} from "./account.query.server";

export async function getCurrentAccountQueryAction() {
  return getCurrentAccountQuery();
}

export async function getAccountByIdQueryAction(
  request: GetAccountByIdRequest,
) {
  return getAccountByIdQuery(request);
}
