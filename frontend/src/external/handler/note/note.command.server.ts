import "server-only";

import {
  type CreateNoteRequest,
  CreateNoteRequestSchema,
  type DeleteNoteRequest,
  DeleteNoteRequestSchema,
  NoteResponseSchema,
  type PublishNoteRequest,
  PublishNoteRequestSchema,
  type UnpublishNoteRequest,
  UnpublishNoteRequestSchema,
  type UpdateNoteByIdRequest,
  UpdateNoteByIdRequestSchema,
} from "../../dto/note.dto";
import { noteService } from "../../service/note/note.service";

// NOTE: 認証チェック（withAuth）は .action.ts で行う
export async function createNoteCommand(
  request: CreateNoteRequest,
  accountId: string,
) {
  const validated = CreateNoteRequestSchema.parse(request);
  const note = await noteService.createNote(accountId, validated);
  return NoteResponseSchema.parse(note);
}

// NOTE: 認証チェック（withAuth）は .action.ts で行う
export async function updateNoteCommand(
  request: UpdateNoteByIdRequest,
  accountId: string,
) {
  const validated = UpdateNoteByIdRequestSchema.parse(request);
  const { id, ...updateData } = validated;
  const note = await noteService.updateNote(id, accountId, updateData);
  return NoteResponseSchema.parse(note);
}

// NOTE: 認証チェック（withAuth）は .action.ts で行う
export async function publishNoteCommand(
  request: PublishNoteRequest,
  accountId: string,
) {
  const validated = PublishNoteRequestSchema.parse(request);
  const note = await noteService.publishNote(validated.noteId, accountId);
  return NoteResponseSchema.parse(note);
}

// NOTE: 認証チェック（withAuth）は .action.ts で行う
export async function unpublishNoteCommand(
  request: UnpublishNoteRequest,
  accountId: string,
) {
  const validated = UnpublishNoteRequestSchema.parse(request);
  const note = await noteService.unpublishNote(validated.noteId, accountId);
  return NoteResponseSchema.parse(note);
}

// NOTE: 認証チェック（withAuth）は .action.ts で行う
export async function deleteNoteCommand(
  request: DeleteNoteRequest,
  accountId: string,
) {
  const validated = DeleteNoteRequestSchema.parse(request);
  await noteService.deleteNote(validated.id, accountId);
  return { success: true };
}
