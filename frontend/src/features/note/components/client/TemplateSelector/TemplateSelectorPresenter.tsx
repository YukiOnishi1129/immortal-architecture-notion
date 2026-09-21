"use client";

import { Loader2 } from "lucide-react";
import type { Template } from "@/features/template/types";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/shared/components/ui/select";

type TemplateSelectorPresenterProps = {
  templates: Template[];
  isLoading: boolean;
  value: string;
  onChange: (value: string) => void;
  disabled?: boolean;
};

export function TemplateSelectorPresenter({
  templates,
  isLoading,
  value,
  onChange,
  disabled,
}: TemplateSelectorPresenterProps) {
  if (isLoading) {
    return (
      <div className="flex items-center gap-2 h-10 px-3 py-2 text-sm rounded-md border border-input bg-background">
        <Loader2 className="h-4 w-4 animate-spin" />
        <span>読み込み中...</span>
      </div>
    );
  }

  return (
    <Select value={value} onValueChange={onChange} disabled={disabled}>
      <SelectTrigger>
        <SelectValue placeholder="テンプレートを選択してください" />
      </SelectTrigger>
      <SelectContent>
        {templates.map((template) => (
          <SelectItem key={template.id} value={template.id}>
            {template.name}
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  );
}
