import "@typespec/http";
import "@typespec/openapi3";
import "../models/account.tsp";
import "../models/common.tsp";

using TypeSpec.Http;
using MiniNotion.Models;

namespace MiniNotion.Routes;

@route("/api/accounts")
@tag("Accounts")
interface Accounts {
  /** 現在のアカウント取得 */
  @get
  @route("/me")
  @summary("Get current account")
  getCurrentAccount(): AccountResponse | UnauthorizedError;

  /** アカウント詳細取得 */
  @get
  @route("/{accountId}")
  @summary("Get account by ID")
  getAccountById(
    @path accountId: string
  ): AccountResponse | NotFoundError | UnauthorizedError;

  /** メールアドレスでアカウント取得 */
  @get
  @route("/by-email")
  @summary("Get account by email")
  getAccountByEmail(
    @query email: string
  ): AccountResponse | NotFoundError | UnauthorizedError;

  /** OAuth認証（内部処理） */
  @post
  @route("/auth")
  @summary("Create or get account via OAuth")
  createOrGetAccount(
    @body request: CreateOrGetAccountRequest
  ): AccountResponse | BadRequestError;
}
