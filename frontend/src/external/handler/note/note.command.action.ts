"use server";

import { withAuth } from "@/features/auth/servers/auth.guard";
import type {
  CreateNoteRequest,
  DeleteNoteRequest,
  PublishNoteRequest,
  UnpublishNoteRequest,
  UpdateNoteByIdRequest,
} from "../../dto/note.dto";
import {
  createNoteCommand,
  deleteNoteCommand,
  publishNoteCommand,
  unpublishNoteCommand,
  updateNoteCommand,
} from "./note.command.server";

export async function createNoteCommandAction(request: CreateNoteRequest) {
  return withAuth(({ accountId }) => createNoteCommand(request, accountId));
}

export async function updateNoteCommandAction(request: UpdateNoteByIdRequest) {
  return withAuth(({ accountId }) => updateNoteCommand(request, accountId));
}

export async function publishNoteCommandAction(request: PublishNoteRequest) {
  return withAuth(({ accountId }) => publishNoteCommand(request, accountId));
}

export async function unpublishNoteCommandAction(
  request: UnpublishNoteRequest,
) {
  return withAuth(({ accountId }) => unpublishNoteCommand(request, accountId));
}

export async function deleteNoteCommandAction(request: DeleteNoteRequest) {
  return withAuth(({ accountId }) => deleteNoteCommand(request, accountId));
}
