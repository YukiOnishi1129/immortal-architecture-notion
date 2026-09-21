"use client";

import { Search } from "lucide-react";
import { Button } from "@/shared/components/ui/button";
import { Input } from "@/shared/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/shared/components/ui/select";

interface NoteListFilterPresenterProps {
  searchQuery: string;
  statusFilter: string;
  isPending: boolean;
  onSearchQueryChange: (value: string) => void;
  onSearchSubmit: (e: React.FormEvent) => void;
  onStatusChange: (value: string) => void;
}

export function NoteListFilterPresenter({
  searchQuery,
  statusFilter,
  isPending,
  onSearchQueryChange,
  onSearchSubmit,
  onStatusChange,
}: NoteListFilterPresenterProps) {
  return (
    <div className="flex flex-col gap-4 sm:flex-row sm:items-center">
      <form onSubmit={onSearchSubmit} className="flex-1">
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

      <div className="sm:w-48">
        <Select
          value={statusFilter}
          onValueChange={onStatusChange}
          disabled={isPending}
        >
          <SelectTrigger className="w-full h-12 text-base">
            <SelectValue placeholder="すべてのステータス" />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">すべてのステータス</SelectItem>
            <SelectItem value="Draft">下書き</SelectItem>
            <SelectItem value="Publish">公開</SelectItem>
          </SelectContent>
        </Select>
      </div>
    </div>
  );
}
