"use client";

import { notFound } from "next/navigation";
import { TemplateDetailPresenter } from "./TemplateDetailPresenter";
import { TemplateDetailSkeleton } from "./TemplateDetailSkeleton";
import { useTemplateDetail } from "./useTemplateDetail";

interface TemplateDetailContainerProps {
  templateId: string;
  isOwner: boolean;
}

export function TemplateDetailContainer({
  templateId,
  isOwner,
}: TemplateDetailContainerProps) {
  const {
    template,
    isLoading,
    isDeleting,
    showDeleteModal,
    handleEdit,
    handleCreateNote,
    handleDeleteClick,
    handleDeleteCancel,
    handleDeleteConfirm,
  } = useTemplateDetail(templateId);

  if (isLoading) {
    return <TemplateDetailSkeleton />;
  }

  if (!template) {
    notFound();
  }

  return (
    <TemplateDetailPresenter
      template={template}
      isOwner={isOwner}
      isDeleting={isDeleting}
      showDeleteModal={showDeleteModal}
      onEdit={handleEdit}
      onCreateNote={handleCreateNote}
      onDeleteClick={handleDeleteClick}
      onDeleteCancel={handleDeleteCancel}
      onDeleteConfirm={handleDeleteConfirm}
    />
  );
}
