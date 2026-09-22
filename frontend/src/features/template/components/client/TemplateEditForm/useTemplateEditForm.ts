"use client";

import type { DropResult } from "@hello-pangea/dnd";
import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter } from "next/navigation";
import { useEffect, useTransition } from "react";
import { useFieldArray, useForm } from "react-hook-form";
import { toast } from "sonner";
import { updateTemplateCommandAction } from "@/external/handler/template/template.command.action";
import { useTemplateQuery } from "@/features/template/hooks/useTemplateQuery";
import { type TemplateEditFormData, templateEditFormSchema } from "./schema";

export function useTemplateEditForm(templateId: string) {
  const router = useRouter();
  const { data: template, isLoading } = useTemplateQuery(templateId);
  const [isPending, startTransition] = useTransition();

  const form = useForm<TemplateEditFormData>({
    resolver: zodResolver(templateEditFormSchema),
    defaultValues: {
      name: "",
      fields: [],
      notionParentPageUrl: "",
    },
  });

  const { fields, append, remove, move } = useFieldArray({
    control: form.control,
    name: "fields",
  });

  // テンプレートデータが読み込まれたらフォームに設定
  useEffect(() => {
    if (template && !isLoading) {
      form.reset({
        name: template.name,
        fields: template.fields.map((field) => ({
          id: field.id,
          label: field.label,
          isRequired: field.isRequired,
          order: field.order,
        })),
        notionParentPageUrl: template.notionParentPageUrl ?? "",
      });
    }
  }, [template, isLoading, form]);

  const handleSubmit = form.handleSubmit((data: TemplateEditFormData) => {
    startTransition(async () => {
      try {
        // 既存項目の id はそのまま送る。落とすとサーバーからは
        // 「消して作り直す」リクエストに見え、ノートが使っている項目を
        // 削除しようとして失敗する。
        const fields = data.fields.map(({ id, label, isRequired }, index) => ({
          id,
          label,
          isRequired,
          order: index + 1,
        }));

        await updateTemplateCommandAction({
          id: templateId,
          name: data.name,
          fields,
          // 空文字は「連携を解除する」という意味でそのまま送る。
          notionParentPageUrl: data.notionParentPageUrl.trim(),
        });

        toast.success("テンプレートを更新しました");
        router.push(`/templates/${templateId}`);
      } catch (error) {
        console.error("テンプレートの更新に失敗しました:", error);

        let errorMessage = "テンプレートの更新に失敗しました";

        if (error instanceof Error) {
          errorMessage = error.message;

          // 特定のエラーメッセージの場合は、わかりやすい通知を表示
          if (errorMessage.includes("ノートで使用されています")) {
            toast.error(errorMessage, {
              duration: 5000, // 5秒間表示
            });
          } else {
            toast.error(errorMessage);
          }
        } else {
          toast.error(errorMessage);
        }

        form.setError("root", {
          message: errorMessage,
        });
      }
    });
  });

  const handleDragEnd = (result: DropResult) => {
    if (!result.destination) return;

    move(result.source.index, result.destination.index);
  };

  const handleAddField = () => {
    append({
      id: crypto.randomUUID(),
      label: "",
      isRequired: false,
      order: fields.length + 1,
    });
  };

  const handleCancel = () => {
    router.push(`/templates/${templateId}`);
  };

  return {
    form,
    fields,
    template,
    isLoading,
    isSubmitting: isPending,
    handleSubmit,
    handleCancel,
    handleDragEnd,
    handleAddField,
    remove,
  };
}
