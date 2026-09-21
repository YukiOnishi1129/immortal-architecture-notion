"use server";

import type {
  GetTemplateByIdRequest,
  ListTemplateRequest,
} from "@/external/dto/template.dto";
import { withAuth } from "@/features/auth/servers/auth.guard";
import {
  getTemplateByIdQuery,
  listMyTemplatesQuery,
  listTemplatesQuery,
} from "./template.query.server";

export async function getTemplateByIdQueryAction(
  request: GetTemplateByIdRequest,
) {
  return getTemplateByIdQuery(request);
}

export async function listTemplatesQueryAction(request?: ListTemplateRequest) {
  return listTemplatesQuery(request);
}

export async function listMyTemplatesQueryAction() {
  return withAuth(({ accountId }) => listMyTemplatesQuery(accountId));
}
