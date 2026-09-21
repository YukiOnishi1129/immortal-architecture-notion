import { z } from "zod";

export const templateFieldSchema = z.object({
  id: z.string(),
  label: z.string().min(1, "項目名は必須です"),
  isRequired: z.boolean(),
  order: z.number(),
});

export const templateEditFormSchema = z.object({
  name: z
    .string()
    .min(1, "テンプレート名は必須です")
    .max(100, "テンプレート名は100文字以内で入力してください"),
  fields: z.array(templateFieldSchema).min(1, "少なくとも1つの項目が必要です"),
});

export type TemplateField = z.infer<typeof templateFieldSchema>;
export type TemplateEditFormData = z.infer<typeof templateEditFormSchema>;
