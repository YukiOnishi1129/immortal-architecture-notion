"use client";

import { format } from "date-fns";
import { ja } from "date-fns/locale";
import { Calendar, Clock, FileText } from "lucide-react";
import type { Route } from "next";
import Link from "next/link";
import { NoteListFilter } from "@/features/note/components/client/NoteListFilter";
import type { Note, NoteFilters } from "@/features/note/types";
import {
  Avatar,
  AvatarFallback,
  AvatarImage,
} from "@/shared/components/ui/avatar";
import { Button } from "@/shared/components/ui/button";
import { Card } from "@/shared/components/ui/card";

interface MyNoteListPresenterProps {
  notes: Note[];
  isLoading?: boolean;
  filters: NoteFilters;
}

export function MyNoteListPresenter({
  notes,
  isLoading,
  filters,
}: MyNoteListPresenterProps) {
  return (
    <div className="space-y-6">
      <div className="bg-white p-6 rounded-lg shadow-sm border">
        <div className="flex items-center justify-between mb-6">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">マイノート</h1>
            <p className="mt-2 text-gray-600">あなたが作成したノートの一覧</p>
          </div>

          <Button asChild>
            <Link href={"/my-notes/new"} className="flex items-center">
              <FileText className="mr-2 h-4 w-4" />
              新規作成
            </Link>
          </Button>
        </div>
        <NoteListFilter filters={filters} />
      </div>

      {isLoading ? (
        <div className="flex items-center justify-center h-64">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-gray-900" />
        </div>
      ) : notes.length === 0 ? (
        <div className="text-center py-12">
          <p className="text-gray-500 mb-4">ノートがありません</p>
          <Button asChild>
            <Link href={"/my-notes/new"}>新しいノートを作成</Link>
          </Button>
        </div>
      ) : (
        <div className="grid gap-4 grid-cols-1 md:grid-cols-2 lg:grid-cols-3">
          {notes.map((note) => (
            <Card
              key={note.id}
              className="overflow-hidden hover:shadow-lg transition-shadow duration-200 relative"
            >
              <div
                className={`absolute top-0 left-0 right-0 h-1 ${
                  note.status === "Publish" ? "bg-green-500" : "bg-gray-400"
                }`}
              />
              <Link
                href={`/my-notes/${note.id}` as Route}
                className="block p-6"
              >
                <div className="space-y-3">
                  <div className="flex items-center justify-between mb-2">
                    <span
                      className={`inline-flex items-center px-3 py-1 rounded-full text-xs font-medium ${
                        note.status === "Publish"
                          ? "bg-green-100 text-green-800"
                          : "bg-gray-100 text-gray-800"
                      }`}
                    >
                      {note.status === "Publish" ? "公開" : "下書き"}
                    </span>
                  </div>
                  <h3 className="text-lg font-semibold text-gray-900 line-clamp-2">
                    {note.title}
                  </h3>
                  <div className="flex items-center gap-2 mb-2">
                    <Avatar className="w-6 h-6">
                      {note.owner.thumbnail ? (
                        <AvatarImage
                          src={note.owner.thumbnail}
                          alt={`${note.owner.firstName} ${note.owner.lastName}`}
                        />
                      ) : null}
                      <AvatarFallback className="text-xs">
                        {note.owner.firstName[0]}
                        {note.owner.lastName[0]}
                      </AvatarFallback>
                    </Avatar>
                    <span className="text-sm text-gray-600">
                      {note.owner.firstName} {note.owner.lastName}
                    </span>
                  </div>
                  <p className="text-sm text-gray-600">
                    テンプレート: {note.templateName}
                  </p>
                  <div className="flex items-center gap-3 text-xs text-gray-500 pt-2">
                    <div className="flex items-center gap-1">
                      <Calendar className="w-3 h-3" />
                      <span>
                        作成:{" "}
                        {format(new Date(note.createdAt), "yyyy/MM/dd", {
                          locale: ja,
                        })}
                      </span>
                    </div>
                    <div className="flex items-center gap-1">
                      <Clock className="w-3 h-3" />
                      <span>
                        更新:{" "}
                        {format(new Date(note.updatedAt), "yyyy/MM/dd", {
                          locale: ja,
                        })}
                      </span>
                    </div>
                  </div>
                </div>
              </Link>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}
