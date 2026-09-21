"use client";

import { zodResolver } from "@hookform/resolvers/zod";
import { useQueryClient } from "@tanstack/react-query";
import type { Route } from "next";
import { useRouter } from "next/navigation";
import { useEffect, useTransition } from "react";
import { useForm } from "react-hook-form";
import { toast } from "sonner";
import { updateNoteCommandAction } from "@/external/handler/note/note.command.action";
import { useNoteDetailQuery } from "@/features/note/hooks/useNoteDetailQuery";
import { noteKeys } from "@/features/note/queries/keys";
import { type NoteEditFormData, NoteEditFormSchema } from "./schema";

type UseNoteEditFormOptions = {
  backTo?: Route;
};

export function useNoteEditForm(
  noteId: string,
  options: UseNoteEditFormOptions = {},
) {
  const { backTo } = options;
  const router = useRouter();
  const queryClient = useQueryClient();
  const { data: note, isLoading } = useNoteDetailQuery(noteId);
  const [isPending, startTransition] = useTransition();

  const form = useForm<NoteEditFormData>({
    resolver: zodResolver(NoteEditFormSchema),
    defaultValues: {
      title: "",
      sections: [],
    },
  });

  // フォームにノートデータをセット
  useEffect(() => {
    if (note) {
      form.reset({
        title: note.title,
        sections: note.sections,
      });
    }
  }, [note, form]);

  const handleSubmit = form.handleSubmit((data) => {
    // 必須項目のバリデーション
    const hasEmptyRequiredSection = data.sections.some(
      (section) => section.isRequired && !section.content.trim(),
    );

    if (hasEmptyRequiredSection) {
      toast.error("必須項目を入力してください");
      return;
    }

    startTransition(async () => {
      try {
        const sections = data.sections.map((section) => ({
          id: section.id,
          content: section.content,
        }));

        const updatedNote = await updateNoteCommandAction({
          id: noteId,
          title: data.title,
          sections,
        });

        toast.success("ノートを更新しました");

        // キャッシュを直接更新
        queryClient.setQueryData(noteKeys.detail(noteId), updatedNote);

        // 一覧のキャッシュも無効化
        await queryClient.invalidateQueries({
          queryKey: noteKeys.lists(),
        });

        const detailPath = backTo ? `/my-notes/${noteId}` : `/notes/${noteId}`;
        router.push(detailPath as Route);
      } catch (error) {
        console.error("Failed to update note:", error);
        toast.error("ノートの更新に失敗しました");
      }
    });
  });

  const handleCancel = () => {
    const detailPath = backTo ? `/my-notes/${noteId}` : `/notes/${noteId}`;
    router.push(detailPath as Route);
  };

  const handleSectionContentChange = (index: number, content: string) => {
    form.setValue(`sections.${index}.content`, content);
  };

  return {
    form,
    note,
    isLoading,
    isSubmitting: isPending,
    handleSubmit,
    handleCancel,
    handleSectionContentChange,
  };
}
