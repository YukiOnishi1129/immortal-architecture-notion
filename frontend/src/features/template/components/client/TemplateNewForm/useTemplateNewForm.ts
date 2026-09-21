"use client";

import type { DropResult } from "@hello-pangea/dnd";
import { zodResolver } from "@hookform/resolvers/zod";
import { useRouter } from "next/navigation";
import { useTransition } from "react";
import { useFieldArray, useForm } from "react-hook-form";
import { toast } from "sonner";
import { createTemplateCommandAction } from "@/external/handler/template/template.command.action";
import { type TemplateNewFormData, templateNewFormSchema } from "./schema";

export function useTemplateNewForm() {
  const router = useRouter();
  const [isPending, startTransition] = useTransition();

  const form = useForm<TemplateNewFormData>({
    resolver: zodResolver(templateNewFormSchema),
    defaultValues: {
      name: "",
      fields: [],
    },
  });

  const { fields, append, remove, move } = useFieldArray({
    control: form.control,
    name: "fields",
  });

  const handleSubmit = form.handleSubmit((data: TemplateNewFormData) => {
    startTransition(async () => {
      try {
        const fields = data.fields.map(({ label, isRequired }, index) => ({
          label,
          isRequired,
          order: index + 1,
        }));

        const result = await createTemplateCommandAction({
          name: data.name,
          fields,
        });

        if (result?.id) {
          toast.success("テンプレートを作成しました");
          router.push("/templates");
        }
      } catch (error) {
        console.error("テンプレートの作成に失敗しました:", error);
        toast.error("テンプレートの作成に失敗しました");
        form.setError("root", {
          message: "テンプレートの作成に失敗しました。もう一度お試しください。",
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
    router.push("/templates");
  };

  return {
    form,
    fields,
    isCreating: isPending,
    handleSubmit,
    handleCancel,
    handleDragEnd,
    handleAddField,
    remove,
  };
}
