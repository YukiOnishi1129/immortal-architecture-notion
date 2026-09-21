import { notesApiClient } from "@/external/client/api/config";
import type { NotesApi } from "@/external/client/api/generated/apis/NotesApi";
import type { ModelsNoteResponse } from "@/external/client/api/generated/models/ModelsNoteResponse";
import type { ModelsNoteStatus } from "@/external/client/api/generated/models/ModelsNoteStatus";
import {
  type CreateNoteRequest,
  type NoteResponse,
  NoteResponseSchema,
  type UpdateNoteRequest,
} from "@/external/dto/note.dto";
import type { NoteFilters } from "@/features/note/types";
import { isNotFoundError } from "../http-error";

function toNoteResponse(model: ModelsNoteResponse): NoteResponse {
  return NoteResponseSchema.parse({
    id: model.id,
    title: model.title,
    templateId: model.templateId,
    templateName: model.templateName,
    ownerId: model.ownerId,
    owner: {
      id: model.owner.id,
      firstName: model.owner.firstName,
      lastName: model.owner.lastName,
      thumbnail: model.owner.thumbnail ?? null,
    },
    status: model.status,
    sections: model.sections.map((section) => ({
      id: section.id,
      fieldId: section.fieldId,
      fieldLabel: section.fieldLabel,
      content: section.content,
      isRequired: section.isRequired,
    })),
    createdAt: model.createdAt.toISOString(),
    updatedAt: model.updatedAt.toISOString(),
  });
}

function toQueryStatus(
  status?: NoteFilters["status"],
): ModelsNoteStatus | undefined {
  if (!status) {
    return undefined;
  }
  return status as ModelsNoteStatus;
}

export class NoteService {
  constructor(private readonly api: NotesApi) {}

  async getNoteById(id: string): Promise<NoteResponse | null> {
    try {
      const note = await this.api.notesGetNoteById({ noteId: id });
      return toNoteResponse(note);
    } catch (error) {
      if (isNotFoundError(error)) {
        return null;
      }
      throw error;
    }
  }

  async getNotes(
    filters?: NoteFilters & { ownerId?: string },
  ): Promise<NoteResponse[]> {
    const list = await this.api.notesListNotes({
      q: filters?.q ?? filters?.search,
      status: toQueryStatus(filters?.status),
      templateId: filters?.templateId,
      ownerId: filters?.ownerId,
    });
    return list.map((note) => toNoteResponse(note));
  }

  async createNote(
    ownerId: string,
    request: CreateNoteRequest,
  ): Promise<NoteResponse> {
    const note = await this.api.notesCreateNote({
      modelsCreateNoteRequest: {
        title: request.title,
        templateId: request.templateId,
        ownerId,
        sections:
          request.sections?.map((section) => ({
            fieldId: section.fieldId,
            content: section.content,
          })) ?? [],
      },
    });
    return toNoteResponse(note);
  }

  async updateNote(
    id: string,
    ownerId: string,
    request: UpdateNoteRequest,
  ): Promise<NoteResponse> {
    const note = await this.api.notesUpdateNote({
      noteId: id,
      ownerId,
      modelsUpdateNoteRequest: {
        id,
        title: request.title,
        sections: request.sections.map((section) => ({
          id: section.id,
          content: section.content,
        })),
      },
    });
    return toNoteResponse(note);
  }

  async publishNote(noteId: string, ownerId: string): Promise<NoteResponse> {
    const note = await this.api.notesPublishNote({ noteId, ownerId });
    return toNoteResponse(note);
  }

  async unpublishNote(noteId: string, ownerId: string): Promise<NoteResponse> {
    const note = await this.api.notesUnpublishNote({ noteId, ownerId });
    return toNoteResponse(note);
  }

  async deleteNote(id: string, ownerId: string): Promise<void> {
    await this.api.notesDeleteNote({ noteId: id, ownerId });
  }
}

export const noteService = new NoteService(notesApiClient);
