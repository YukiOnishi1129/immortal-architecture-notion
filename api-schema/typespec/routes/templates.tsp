import "@typespec/http";
import "@typespec/openapi3";
import "../models/template.tsp";
import "../models/common.tsp";

using TypeSpec.Http;
using MiniNotion.Models;

namespace MiniNotion.Routes;

@route("/api/templates")
@tag("Templates")
interface Templates {
  /** テンプレート一覧取得 */
  @get
  @summary("Get templates list")
  listTemplates(
    /** テンプレート名のキーワード検索 */
    @query q?: string,

    /** 所有者IDフィルター */
    @query ownerId?: string
  ): TemplateResponse[] | UnauthorizedError;

  /** テンプレート詳細取得 */
  @get
  @route("/{templateId}")
  @summary("Get template by ID")
  getTemplateById(
    @path templateId: string
  ): TemplateResponse | NotFoundError | UnauthorizedError;

  /** テンプレート作成 */
  @post
  @summary("Create template")
  createTemplate(
    @body request: CreateTemplateRequest
  ): TemplateResponse | BadRequestError | UnauthorizedError;

  /** テンプレート更新 */
  @put
  @route("/{templateId}")
  @summary("Update template")
  updateTemplate(
    @path templateId: string,
    @query ownerId: string,
    @body request: UpdateTemplateRequest
  ): TemplateResponse | NotFoundError | ForbiddenError | BadRequestError | UnauthorizedError;

  /** テンプレート削除 */
  @delete
  @route("/{templateId}")
  @summary("Delete template")
  deleteTemplate(
    @path templateId: string,
    @query ownerId: string
  ): SuccessResponse | NotFoundError | ForbiddenError | BadRequestError | UnauthorizedError;
}
