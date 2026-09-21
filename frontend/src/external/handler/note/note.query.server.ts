import "server-only";

import { requireAuthServer } from "@/features/auth/servers/redirect.server";
import {
  type GetNoteByIdRequest,
  GetNoteByIdRequestSchema,
  type ListMyNoteRequest,
  ListMyNoteRequestSchema,
  type ListNoteRequest,
  ListNoteRequestSchema,
  NoteResponseSchema,
} from "../../dto/note.dto";
import { noteService } from "../../service/note/note.service";

export async function getNoteByIdQuery(request: GetNoteByIdRequest) {
  await requireAuthServer();

  const validated = GetNoteByIdRequestSchema.parse(request);
  const note = await noteService.getNoteById(validated.id);

  if (!note) {
    return null;
  }

  return NoteResponseSchema.parse(note);
}

export async function listNoteQuery(request?: ListNoteRequest) {
  await requireAuthServer();

  const validated = request ? ListNoteRequestSchema.parse(request) : undefined;
  const notes = await noteService.getNotes(validated);
  return notes.map((note) => NoteResponseSchema.parse(note));
}

// NOTE: 認証チェック（withAuth）は .action.ts で行う
export async function listMyNoteQuery(
  request: ListMyNoteRequest | undefined,
  accountId: string,
) {
  const validated = request ? ListMyNoteRequestSchema.parse(request) : {};

  const notes = await noteService.getNotes({
    ...validated,
    ownerId: accountId,
  });

  return notes.map((note) => NoteResponseSchema.parse(note));
}
