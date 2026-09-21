"use client";

import type { Route } from "next";
import { NoteDetailPresenter } from "./NoteDetailPresenter";
import { useNoteDetail } from "./useNoteDetail";

type NoteDetailContainerProps = {
  noteId: string;
  isOwner?: boolean;
  backTo?: Route;
};

export function NoteDetailContainer({
  noteId,
  isOwner = true,
  backTo,
}: NoteDetailContainerProps) {
  const {
    note,
    isLoading,
    isDeleting,
    isTogglingPublish,
    showDeleteDialog,
    showPublishDialog,
    handleEdit,
    handleDelete,
    handleConfirmDelete,
    handleCancelDelete,
    handleTogglePublish,
    handleConfirmPublish,
    handleCancelPublish,
  } = useNoteDetail(noteId, { backTo });

  return (
    <NoteDetailPresenter
      note={note}
      isLoading={isLoading}
      isDeleting={isDeleting}
      isTogglingPublish={isTogglingPublish}
      showDeleteDialog={showDeleteDialog}
      showPublishDialog={showPublishDialog}
      isOwner={isOwner}
      backTo={backTo}
      onEdit={handleEdit}
      onDelete={handleDelete}
      onConfirmDelete={handleConfirmDelete}
      onCancelDelete={handleCancelDelete}
      onTogglePublish={handleTogglePublish}
      onConfirmPublish={handleConfirmPublish}
      onCancelPublish={handleCancelPublish}
    />
  );
}
