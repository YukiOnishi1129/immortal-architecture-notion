"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import type { Route } from "next";
import { useRouter } from "next/navigation";
import { useEffect, useTransition } from "react";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import { createNoteCommandAction } from "@/external/handler/note/note.command.action";
import { useTemplateQuery } from "@/features/template/hooks/useTemplateQuery";
import { type NoteNewFormData, noteNewFormSchema } from "./schema";

type UseNoteNewFormProps = {
  backTo?: Route;
  initialTemplateId?: string;
};

export function useNoteNewForm({
  backTo,
  initialTemplateId,
}: UseNoteNewFormProps = {}) {
  const router = useRouter();
  const [isPending, startTransition] = useTransition();

  const form = useForm<NoteNewFormData>({
    resolver: zodResolver(noteNewFormSchema),
    defaultValues: {
      title: "",
      templateId: initialTemplateId ?? "",
      sections: [],
    },
  });

  const templateId = form.watch("templateId");

  // TanStack QueryでテンプレートをフェッチするHookを使用
  const { data: selectedTemplate, isLoading: isLoadingTemplate } =
    useTemplateQuery(templateId);

  // テンプレートが取得されたらsectionsを初期化
  useEffect(() => {
    if (!templateId) {
      form.setValue("sections", []);
      return;
    }

    if (selectedTemplate && !isLoadingTemplate) {
      // テンプレートのフィールドを基に sections を初期化
      const sections = selectedTemplate.fields.map((field) => ({
        fieldId: field.id,
        fieldLabel: field.label,
        content: "",
        isRequired: field.isRequired,
      }));
      form.setValue("sections", sections);
    }
  }, [selectedTemplate, templateId, isLoadingTemplate, form]);

  const handleSubmit = form.handleSubmit((data: NoteNewFormData) => {
    // Validate required sections
    const hasEmptyRequiredFields = data.sections.some(
      (section) => section.isRequired && !section.content.trim(),
    );

    if (hasEmptyRequiredFields) {
      toast.error("必須項目を入力してください");
      return;
    }

    startTransition(async () => {
      try {
        const result = await createNoteCommandAction({
          title: data.title,
          templateId: data.templateId,
          sections: data.sections.map(({ fieldId, content }) => ({
            fieldId,
            content,
          })),
        });

        if (result?.id) {
          toast.success("ノートを作成しました");
          // backToが指定されている場合はbackToに戻る、それ以外はマイノート一覧へ（新規作成は下書きのため）
          const redirectPath = backTo ?? "/my-notes";
          router.push(redirectPath);
        }
      } catch (error) {
        console.error("ノートの作成に失敗しました:", error);
        toast.error("ノートの作成に失敗しました");
        form.setError("root", {
          message: "ノートの作成に失敗しました。もう一度お試しください。",
        });
      }
    });
  });

  const handleCancel = () => {
    router.push(backTo ?? "/my-notes");
  };

  const handleSectionContentChange = (index: number, content: string) => {
    const sections = form.getValues("sections");
    sections[index].content = content;
    form.setValue("sections", sections);
  };

  return {
    form,
    selectedTemplate,
    isLoadingTemplate,
    isCreating: isPending,
    handleSubmit,
    handleCancel,
    handleSectionContentChange,
  };
}
