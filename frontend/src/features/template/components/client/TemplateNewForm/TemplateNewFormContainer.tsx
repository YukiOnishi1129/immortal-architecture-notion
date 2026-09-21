"use client";

import { TemplateNewFormPresenter } from "./TemplateNewFormPresenter";
import { useTemplateNewForm } from "./useTemplateNewForm";

export function TemplateNewFormContainer() {
  const {
    form,
    fields,
    isCreating,
    handleSubmit,
    handleCancel,
    handleDragEnd,
    handleAddField,
    remove,
  } = useTemplateNewForm();

  return (
    <TemplateNewFormPresenter
      form={form}
      fields={fields}
      isCreating={isCreating}
      onSubmit={handleSubmit}
      onCancel={handleCancel}
      onRemoveField={remove}
      onDragEnd={handleDragEnd}
      onAddField={handleAddField}
    />
  );
}
