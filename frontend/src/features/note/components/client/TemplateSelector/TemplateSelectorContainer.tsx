"use client";

import { TemplateSelectorPresenter } from "./TemplateSelectorPresenter";
import { useTemplateSelector } from "./useTemplateSelector";

type TemplateSelectorContainerProps = {
  value: string;
  onChange: (value: string) => void;
  disabled?: boolean;
};

export function TemplateSelectorContainer({
  value,
  onChange,
  disabled,
}: TemplateSelectorContainerProps) {
  const { templates, isLoading } = useTemplateSelector();

  return (
    <TemplateSelectorPresenter
      templates={templates}
      isLoading={isLoading}
      value={value}
      onChange={onChange}
      disabled={disabled}
    />
  );
}
