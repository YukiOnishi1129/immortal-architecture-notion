"use client";

import { TemplateEditFormPresenter } from "./TemplateEditFormPresenter";
import { TemplateEditFormSkeleton } from "./TemplateEditFormSkeleton";
import { useTemplateEditForm } from "./useTemplateEditForm";

interface TemplateEditFormContainerProps {
  templateId: string;
}

export function TemplateEditFormContainer({
  templateId,
}: TemplateEditFormContainerProps) {
  const {
    form,
    fields,
    template,
    isLoading,
    isSubmitting,
    handleSubmit,
    handleCancel,
    handleDragEnd,
    handleAddField,
    remove,
  } = useTemplateEditForm(templateId);

  if (isLoading) {
    return <TemplateEditFormSkeleton />;
  }

  return (
    <TemplateEditFormPresenter
      form={form}
      fields={fields}
      templateName={template?.name}
      templateId={templateId}
      isSubmitting={isSubmitting}
      onSubmit={handleSubmit}
      onCancel={handleCancel}
      onRemoveField={remove}
      onDragEnd={handleDragEnd}
      onAddField={handleAddField}
    />
  );
}
