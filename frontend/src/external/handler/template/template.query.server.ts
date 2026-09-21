import "server-only";

import { getSessionServer } from "@/features/auth/servers/auth.server";
import { requireAuthServer } from "@/features/auth/servers/redirect.server";
import {
  type GetTemplateByIdRequest,
  GetTemplateByIdRequestSchema,
  type ListTemplateRequest,
  ListTemplateRequestSchema,
  TemplateDetailResponseSchema,
  TemplateResponseSchema,
} from "../../dto/template.dto";
import { templateService } from "../../service/template/template.service";

export async function getTemplateByIdQuery(request: GetTemplateByIdRequest) {
  const validated = GetTemplateByIdRequestSchema.parse(request);
  const template = await templateService.getTemplateById(validated.id);

  if (!template) {
    return null;
  }

  return TemplateDetailResponseSchema.parse(template);
}

export async function listTemplatesQuery(request?: ListTemplateRequest) {
  await requireAuthServer();

  // Get current user for onlyMyTemplates filter
  const session = await getSessionServer();
  if (!session?.account?.id) {
    throw new Error("Unauthorized: No active session");
  }

  const validated = request ? ListTemplateRequestSchema.parse(request) : {};

  const ownerFilter =
    validated?.onlyMyTemplates && session?.account.id
      ? session.account.id
      : validated?.ownerId;

  const templates = await templateService.getTemplates({
    ownerId: ownerFilter,
    q: validated?.q,
  });

  return templates.map((template) => TemplateResponseSchema.parse(template));
}

// NOTE: 認証チェック（withAuth）は .action.ts で行う
export async function listMyTemplatesQuery(accountId: string) {
  const templates = await templateService.getTemplates({
    ownerId: accountId,
  });

  return templates.map((template) => TemplateResponseSchema.parse(template));
}
