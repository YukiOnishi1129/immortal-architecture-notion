"use client";

import { useMutation, useQueryClient } from "@tanstack/react-query";
import type { Route } from "next";
import { useRouter } from "next/navigation";
import { useCallback, useState } from "react";
import { toast } from "sonner";
import {
  deleteNoteCommandAction,
  publishNoteCommandAction,
  unpublishNoteCommandAction,
} from "@/external/handler/note/note.command.action";
import { useNoteDetailQuery } from "@/features/note/hooks/useNoteDetailQuery";
import { noteKeys } from "@/features/note/queries/keys";
import type { Note } from "@/features/note/types";

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
    onError: () => {
      toast.error("ノートの公開に失敗しました");
    },
  });

  const unpublishMutation = useMutation({
    mutationFn: () => unpublishNoteCommandAction({ noteId }),
    onSuccess: (updatedNote: Note) => {
      toast.success("ノートを下書きに戻しました");
      queryClient.setQueryData(noteKeys.detail(noteId), updatedNote);
      queryClient.invalidateQueries({ queryKey: noteKeys.lists() });
    },
    onError: () => {
      toast.error("ノートの非公開に失敗しました");
    },
  });

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

  return {
    note,
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
