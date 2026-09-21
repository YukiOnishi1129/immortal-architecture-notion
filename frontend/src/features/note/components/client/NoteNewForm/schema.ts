import { z } from "zod";

export const sectionInputSchema = z.object({
  fieldId: z.string().uuid(),
  fieldLabel: z.string(),
  content: z.string(),
  isRequired: z.boolean(),
});

export const noteNewFormSchema = z.object({
  title: z
    .string()
    .min(1, "タイトルは必須です")
    .max(100, "タイトルは100文字以内で入力してください"),
  templateId: z.string().uuid("テンプレートを選択してください"),
  sections: z.array(sectionInputSchema),
});

export type SectionInput = z.infer<typeof sectionInputSchema>;
export type NoteNewFormData = z.infer<typeof noteNewFormSchema>;
