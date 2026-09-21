"use server";

import type {
  CreateTemplateRequest,
  DeleteTemplateRequest,
  UpdateTemplateByIdRequest,
} from "@/external/dto/template.dto";
import { withAuth } from "@/features/auth/servers/auth.guard";
import {
  createTemplateCommand,
  deleteTemplateCommand,
  updateTemplateCommand,
} from "./template.command.server";

export async function createTemplateCommandAction(
  request: CreateTemplateRequest,
) {
  return withAuth(({ accountId }) => createTemplateCommand(request, accountId));
}

export async function updateTemplateCommandAction(
  request: UpdateTemplateByIdRequest,
) {
  return withAuth(({ accountId }) => updateTemplateCommand(request, accountId));
}

export async function deleteTemplateCommandAction(
  request: DeleteTemplateRequest,
) {
  return withAuth(({ accountId }) => deleteTemplateCommand(request, accountId));
}
