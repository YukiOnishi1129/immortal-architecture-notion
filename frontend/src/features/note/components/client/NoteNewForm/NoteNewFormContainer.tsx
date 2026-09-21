"use client";

import type { Route } from "next";
import { NoteNewFormPresenter } from "./NoteNewFormPresenter";
import { useNoteNewForm } from "./useNoteNewForm";

type NoteNewFormContainerProps = {
  backTo?: Route;
  initialTemplateId?: string;
};

export function NoteNewFormContainer({
  backTo,
  initialTemplateId,
}: NoteNewFormContainerProps) {
  const {
    form,
    selectedTemplate,
    isLoadingTemplate,
    isCreating,
    handleSubmit,
    handleCancel,
    handleSectionContentChange,
  } = useNoteNewForm({ backTo, initialTemplateId });

  return (
    <NoteNewFormPresenter
      form={form}
      selectedTemplate={selectedTemplate}
      isLoadingTemplate={isLoadingTemplate}
      isCreating={isCreating}
      backTo={backTo}
      onSubmit={handleSubmit}
      onCancel={handleCancel}
      onSectionContentChange={handleSectionContentChange}
    />
  );
}
