"use server";

import { withAuth } from "@/features/auth/servers/auth.guard";
import type {
  GetNoteByIdRequest,
  ListMyNoteRequest,
  ListNoteRequest,
} from "../../dto/note.dto";
import {
  getNoteByIdQuery,
  listMyNoteQuery,
  listNoteQuery,
} from "./note.query.server";

export async function getNoteByIdQueryAction(request: GetNoteByIdRequest) {
  return getNoteByIdQuery(request);
}

export async function listNoteQueryAction(request?: ListNoteRequest) {
  return listNoteQuery(request);
}

export async function listMyNoteQueryAction(request?: ListMyNoteRequest) {
  return withAuth(({ accountId }) => listMyNoteQuery(request, accountId));
}
