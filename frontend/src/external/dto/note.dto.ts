import { z } from "zod";

// Query request schemas
export const GetNoteByIdRequestSchema = z.object({
  id: z.uuid(),
});

export const ListNoteRequestSchema = z.object({
  status: z.enum(["Draft", "Publish"]).optional(),
  templateId: z.uuid().optional(),
  q: z.string().optional(),
  page: z.number().int().positive().optional(),
  ownerId: z.uuid().optional(),
});

export const ListMyNoteRequestSchema = z.object({
  status: z.enum(["Draft", "Publish"]).optional(),
  templateId: z.uuid().optional(),
  q: z.string().optional(),
  page: z.number().int().positive().optional(),
});

// Section schemas
export const SectionResponseSchema = z.object({
  id: z.uuid(),
  fieldId: z.uuid(),
  fieldLabel: z.string(),
  content: z.string(),
  isRequired: z.boolean(),
});

export const SectionInputSchema = z.object({
  fieldId: z.uuid(),
  content: z.string(),
});

// Owner schema
export const OwnerSchema = z.object({
  id: z.uuid(),
  firstName: z.string(),
  lastName: z.string(),
  thumbnail: z.string().nullable(),
});

// Note schemas
export const NoteResponseSchema = z.object({
  id: z.uuid(),
  title: z.string(),
  templateId: z.uuid(),
  templateName: z.string(),
  ownerId: z.uuid(),
  owner: OwnerSchema,
  status: z.enum(["Draft", "Publish"]),
  sections: z.array(SectionResponseSchema),
  createdAt: z.iso.datetime(),
  updatedAt: z.iso.datetime(),
});

export const CreateNoteRequestSchema = z.object({
  title: z.string().min(1).max(100),
  templateId: z.uuid(),
  sections: z.array(SectionInputSchema),
});

export const UpdateNoteRequestSchema = z.object({
  title: z.string().min(1).max(100),
  sections: z.array(
    z.object({
      id: z.uuid(),
      content: z.string(),
    }),
  ),
});

export const UpdateNoteByIdRequestSchema = z.object({
  id: z.uuid(),
  title: z.string().min(1).max(100),
  sections: z.array(
    z.object({
      id: z.uuid(),
      content: z.string(),
    }),
  ),
});

export const DeleteNoteRequestSchema = z.object({
  id: z.uuid(),
});

export const PublishNoteRequestSchema = z.object({
  noteId: z.uuid(),
});

export const UnpublishNoteRequestSchema = z.object({
  noteId: z.uuid(),
});

// Type exports
export type GetNoteByIdRequest = z.infer<typeof GetNoteByIdRequestSchema>;
export type ListNoteRequest = z.infer<typeof ListNoteRequestSchema>;
export type ListMyNoteRequest = z.infer<typeof ListMyNoteRequestSchema>;
export type Owner = z.infer<typeof OwnerSchema>;
export type SectionResponse = z.infer<typeof SectionResponseSchema>;
export type SectionInput = z.infer<typeof SectionInputSchema>;
export type NoteResponse = z.infer<typeof NoteResponseSchema>;
export type CreateNoteRequest = z.infer<typeof CreateNoteRequestSchema>;
export type UpdateNoteRequest = z.infer<typeof UpdateNoteRequestSchema>;
export type UpdateNoteByIdRequest = z.infer<typeof UpdateNoteByIdRequestSchema>;
export type DeleteNoteRequest = z.infer<typeof DeleteNoteRequestSchema>;
export type PublishNoteRequest = z.infer<typeof PublishNoteRequestSchema>;
export type UnpublishNoteRequest = z.infer<typeof UnpublishNoteRequestSchema>;
