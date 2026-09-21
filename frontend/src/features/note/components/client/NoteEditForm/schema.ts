import { z } from "zod";

const SectionSchema = z.object({
  id: z.string(),
  fieldId: z.string(),
  fieldLabel: z.string(),
  content: z.string(),
  isRequired: z.boolean(),
});

export const NoteEditFormSchema = z.object({
  title: z
    .string()
    .min(1, "タイトルは必須です")
    .max(100, "タイトルは100文字以内で入力してください"),
  sections: z.array(SectionSchema),
});

export type NoteEditFormData = z.infer<typeof NoteEditFormSchema>;
