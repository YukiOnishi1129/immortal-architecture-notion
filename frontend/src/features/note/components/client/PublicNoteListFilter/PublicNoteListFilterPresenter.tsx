"use client";

import { Search } from "lucide-react";
import { Button } from "@/shared/components/ui/button";
import { Input } from "@/shared/components/ui/input";

interface PublicNoteListFilterPresenterProps {
  searchQuery: string;
  isPending: boolean;
  onSearchQueryChange: (value: string) => void;
  onSearchSubmit: (e: React.FormEvent) => void;
}

export function PublicNoteListFilterPresenter({
  searchQuery,
  isPending,
  onSearchQueryChange,
  onSearchSubmit,
}: PublicNoteListFilterPresenterProps) {
  return (
    <form onSubmit={onSearchSubmit} className="w-full">
      <div className="relative">
        <Input
          type="text"
          value={searchQuery}
          onChange={(e) => onSearchQueryChange(e.target.value)}
          placeholder="ノートを検索..."
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
  );
}
