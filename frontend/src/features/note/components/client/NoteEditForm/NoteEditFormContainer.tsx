"use client";

import type { Route } from "next";
import { NoteEditFormPresenter } from "./NoteEditFormPresenter";
import { useNoteEditForm } from "./useNoteEditForm";

type NoteEditFormContainerProps = {
  noteId: string;
  backTo?: Route;
};

export function NoteEditFormContainer({
  noteId,
  backTo,
}: NoteEditFormContainerProps) {
  const {
    form,
    note,
    isLoading,
    isSubmitting,
    handleSubmit,
    handleCancel,
    handleSectionContentChange,
  } = useNoteEditForm(noteId, { backTo });

  if (isLoading) {
    return <NoteEditFormPresenter isLoading />;
  }

  if (!note) {
    return null;
  }

  return (
    <NoteEditFormPresenter
      form={form}
      note={note}
      isSubmitting={isSubmitting}
      backTo={backTo}
      onSubmit={handleSubmit}
      onCancel={handleCancel}
      onSectionContentChange={handleSectionContentChange}
    />
  );
}
