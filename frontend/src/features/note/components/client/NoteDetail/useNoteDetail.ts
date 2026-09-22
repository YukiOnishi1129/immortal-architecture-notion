"use client";

import { useMutation, useQueryClient } from "@tanstack/react-query";
import type { Route } from "next";
import { useRouter } from "next/navigation";
import { useCallback, useState } from "react";
import { toast } from "sonner";
import {
  deleteNoteCommandAction,
  publishNoteCommandAction,
  syncNoteToNotionCommandAction,
  unpublishNoteCommandAction,
} from "@/external/handler/note/note.command.action";
import { useNoteDetailQuery } from "@/features/note/hooks/useNoteDetailQuery";
import { noteKeys } from "@/features/note/queries/keys";
import type { Note } from "@/features/note/types";
import { useTemplateQuery } from "@/features/template/hooks/useTemplateQuery";

type UseNoteDetailOptions = {
  backTo?: Route;
};

export function useNoteDetail(
  noteId: string,
  options: UseNoteDetailOptions = {},
) {
  const { backTo } = options;
  const router = useRouter();
  const queryClient = useQueryClient();
  const [showDeleteDialog, setShowDeleteDialog] = useState(false);
  const [showPublishDialog, setShowPublishDialog] = useState(false);

  const { data: note, isLoading } = useNoteDetailQuery(noteId);

  const deleteMutation = useMutation({
    mutationFn: () => deleteNoteCommandAction({ id: noteId }),
    onSuccess: () => {
      toast.success("ノートを削除しました");
      queryClient.invalidateQueries({ queryKey: noteKeys.lists() });
      router.push(backTo ?? "/notes");
    },
    onError: () => {
      toast.error("ノートの削除に失敗しました");
    },
  });

  const publishMutation = useMutation({
    mutationFn: () => publishNoteCommandAction({ noteId }),
    onSuccess: (updatedNote: Note) => {
      toast.success("ノートを公開しました");
      queryClient.setQueryData(noteKeys.detail(noteId), updatedNote);
      queryClient.invalidateQueries({ queryKey: noteKeys.lists() });
    },
    onError: (error: unknown) => {
      // Notion連携の失敗など、理由はサーバーが返す。握りつぶすと
      // 「なぜ公開できないのか」が画面から分からなくなる。
      toast.error(errorMessage(error, "ノートの公開に失敗しました"));
    },
  });

  const unpublishMutation = useMutation({
    mutationFn: () => unpublishNoteCommandAction({ noteId }),
    onSuccess: (updatedNote: Note) => {
      toast.success("ノートを下書きに戻しました");
      queryClient.setQueryData(noteKeys.detail(noteId), updatedNote);
      queryClient.invalidateQueries({ queryKey: noteKeys.lists() });
    },
    onError: (error: unknown) => {
      toast.error(errorMessage(error, "ノートの非公開に失敗しました"));
    },
  });

  const syncMutation = useMutation({
    mutationFn: () => syncNoteToNotionCommandAction({ noteId }),
    onSuccess: (updatedNote: Note) => {
      toast.success("Notionに連携しました");
      queryClient.setQueryData(noteKeys.detail(noteId), updatedNote);
      queryClient.invalidateQueries({ queryKey: noteKeys.lists() });
    },
    onError: (error: unknown) => {
      toast.error(errorMessage(error, "Notionへの連携に失敗しました"));
    },
  });

  const handleSyncToNotion = useCallback(() => {
    syncMutation.mutate();
  }, [syncMutation]);

  const handleEdit = useCallback(() => {
    const editPath = backTo
      ? `/my-notes/${noteId}/edit`
      : `/notes/${noteId}/edit`;
    router.push(editPath as Route);
  }, [backTo, noteId, router]);

  const handleDelete = useCallback(() => {
    setShowDeleteDialog(true);
  }, []);

  const handleConfirmDelete = useCallback(() => {
    deleteMutation.mutate();
    setShowDeleteDialog(false);
  }, [deleteMutation]);

  const handleCancelDelete = useCallback(() => {
    setShowDeleteDialog(false);
  }, []);

  const handleTogglePublish = useCallback(() => {
    if (note?.status === "Publish") {
      unpublishMutation.mutate();
    } else {
      setShowPublishDialog(true);
    }
  }, [note?.status, unpublishMutation]);

  const handleConfirmPublish = useCallback(() => {
    publishMutation.mutate();
    setShowPublishDialog(false);
  }, [publishMutation]);

  const handleCancelPublish = useCallback(() => {
    setShowPublishDialog(false);
  }, []);

  // 公開するとNotionへ連携するため、テンプレートに親ページが要る。
  // ノートのレスポンスには含まれないので、テンプレートを引いて判定する。
  const { data: template } = useTemplateQuery(note?.templateId ?? "");
  const hasNotionParentPage = Boolean(template?.notionParentPageUrl);

  return {
    note,
    template,
    hasNotionParentPage,
    isSyncingToNotion: syncMutation.isPending,
    handleSyncToNotion,
    isLoading,
    isDeleting: deleteMutation.isPending,
    isTogglingPublish: publishMutation.isPending || unpublishMutation.isPending,
    showDeleteDialog,
    showPublishDialog,
    handleEdit,
    handleDelete,
    handleConfirmDelete,
    handleCancelDelete,
    handleTogglePublish,
    handleConfirmPublish,
    handleCancelPublish,
  };
}

/** サーバーが返した理由を使う。取り出せないときだけ既定の文言にする。 */
function errorMessage(error: unknown, fallback: string): string {
  if (error instanceof Error && error.message !== "") {
    return error.message;
  }
  return fallback;
}
