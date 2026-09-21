"use client";

import { useRouter } from "next/navigation";
import { useState, useTransition } from "react";
import { toast } from "sonner";
import { deleteTemplateCommandAction } from "@/external/handler/template/template.command.action";
import { useTemplateQuery } from "@/features/template/hooks/useTemplateQuery";

export function useTemplateDetail(templateId: string) {
  const router = useRouter();
  const { data: template, isLoading, error } = useTemplateQuery(templateId);
  const [isPending, startTransition] = useTransition();
  const [showDeleteModal, setShowDeleteModal] = useState(false);

  const handleEdit = () => {
    router.push(`/templates/${templateId}/edit`);
  };

  const handleCreateNote = () => {
    router.push(`/notes/new?templateId=${templateId}`);
  };

  const handleDeleteClick = () => {
    setShowDeleteModal(true);
  };

  const handleDeleteCancel = () => {
    setShowDeleteModal(false);
  };

  const handleDeleteConfirm = () => {
    startTransition(async () => {
      try {
        await deleteTemplateCommandAction({ id: templateId });
        toast.success("テンプレートを削除しました");
        router.push("/templates");
      } catch (error) {
        console.error("テンプレートの削除に失敗しました:", error);
        toast.error("テンプレートの削除に失敗しました");
      } finally {
        setShowDeleteModal(false);
      }
    });
  };

  return {
    template,
    isLoading,
    error,
    isDeleting: isPending,
    showDeleteModal,
    handleEdit,
    handleCreateNote,
    handleDeleteClick,
    handleDeleteCancel,
    handleDeleteConfirm,
  };
}
