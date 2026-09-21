"use client";

import { Search } from "lucide-react";
import { Button } from "@/shared/components/ui/button";
import { Checkbox } from "@/shared/components/ui/checkbox";
import { Input } from "@/shared/components/ui/input";

interface TemplateListFilterPresenterProps {
  searchQuery: string;
  isPending: boolean;
  onlyMyTemplates: boolean;
  onSearchQueryChange: (value: string) => void;
  onSearchSubmit: (e: React.FormEvent) => void;
  onOnlyMyTemplatesChange: (checked: boolean) => void;
}

export function TemplateListFilterPresenter({
  searchQuery,
  isPending,
  onlyMyTemplates,
  onSearchQueryChange,
  onSearchSubmit,
  onOnlyMyTemplatesChange,
}: TemplateListFilterPresenterProps) {
  return (
    <div className="space-y-4">
      <form onSubmit={onSearchSubmit} className="w-full">
        <div className="relative">
          <Input
            type="text"
            value={searchQuery}
            onChange={(e) => onSearchQueryChange(e.target.value)}
            placeholder="テンプレートを検索..."
            disabled={isPending}
            className="w-full pl-10 pr-20 h-12 text-base"
          />
          <div className="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none">
            <Search className="w-5 h-5 text-gray-400" />
          </div>
          <Button
            type="submit"
            disabled={isPending}
            variant="default"
            size="sm"
            className="absolute right-1.5 top-1.5 h-9"
          >
            検索
          </Button>
        </div>
      </form>

      <div className="flex items-center space-x-2">
        <Checkbox
          id="only-my-templates"
          checked={onlyMyTemplates}
          onCheckedChange={onOnlyMyTemplatesChange}
          disabled={isPending}
        />
        <label
          htmlFor="only-my-templates"
          className="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70"
        >
          自分のテンプレートのみ表示
        </label>
      </div>
    </div>
  );
}
