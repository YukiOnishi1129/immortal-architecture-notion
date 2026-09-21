"use client";

import { Loader2 } from "lucide-react";
import type { Route } from "next";
import Link from "next/link";
import type { UseFormReturn } from "react-hook-form";
import type { Note } from "@/features/note/types";
import { Breadcrumb } from "@/shared/components/ui/breadcrumb";
import { Button } from "@/shared/components/ui/button";
import {
  Card,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/shared/components/ui/card";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/shared/components/ui/form";
import { Input } from "@/shared/components/ui/input";
import { Textarea } from "@/shared/components/ui/textarea";
import { NoteEditFormSkeleton } from "./NoteEditFormSkeleton";
import type { NoteEditFormData } from "./schema";

interface NoteEditFormPresenterProps {
  form?: UseFormReturn<NoteEditFormData>;
  note?: Note;
  isLoading?: boolean;
  isSubmitting?: boolean;
  backTo?: Route;
  onSubmit?: (e: React.FormEvent) => void;
  onSectionContentChange?: (index: number, content: string) => void;
  onCancel?: () => void;
}

export function NoteEditFormPresenter({
  form,
  note,
  isLoading,
  isSubmitting,
  backTo,
  onSubmit,
  onSectionContentChange,
  onCancel,
}: NoteEditFormPresenterProps) {
  if (isLoading) {
    return <NoteEditFormSkeleton />;
  }

  if (!note || !form) {
    return (
      <Card>
        <CardContent className="p-8 text-center">
          <p className="text-muted-foreground">ノートが見つかりません</p>
        </CardContent>
      </Card>
    );
  }

  const listPath = backTo ?? "/notes";
  const listLabel = backTo === "/my-notes" ? "マイノート" : "みんなのノート";
  const detailPath = backTo ? `/my-notes/${note.id}` : `/notes/${note.id}`;

  const breadcrumbItems = [
    {
      label: listLabel,
      href: listPath,
    },
    {
      label: note.title,
      href: detailPath as Route,
    },
    {
      label: "編集",
    },
  ];

  const sections = form.watch("sections");

  return (
    <div className="space-y-4">
      <Breadcrumb items={breadcrumbItems} />
      <Card>
        <CardHeader>
          <CardTitle>ノート編集</CardTitle>
        </CardHeader>
        <CardContent>
          <Form {...form}>
            <form onSubmit={onSubmit} className="space-y-6">
              <FormField
                control={form.control}
                name="title"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>タイトル</FormLabel>
                    <FormControl>
                      <Input
                        {...field}
                        placeholder="タイトルを入力"
                        maxLength={100}
                        disabled={isSubmitting}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <div className="space-y-4">
                {sections.map((section, index) => (
                  <FormItem key={section.id}>
                    <FormLabel>
                      {section.fieldLabel}
                      {section.isRequired && (
                        <span className="text-destructive ml-1">*</span>
                      )}
                    </FormLabel>
                    <FormControl>
                      <Textarea
                        value={section.content}
                        onChange={(e) =>
                          onSectionContentChange?.(index, e.target.value)
                        }
                        placeholder={`${section.fieldLabel}を入力`}
                        disabled={isSubmitting}
                        className="min-h-[100px]"
                      />
                    </FormControl>
                    {section.isRequired && !section.content && (
                      <p className="text-sm text-destructive">
                        この項目は必須です
                      </p>
                    )}
                  </FormItem>
                ))}
              </div>

              <div className="flex justify-between">
                <Button
                  type="button"
                  variant="outline"
                  onClick={onCancel}
                  disabled={isSubmitting}
                  asChild
                >
                  <Link href={`/notes/${note.id}`}>キャンセル</Link>
                </Button>
                <Button type="submit" disabled={isSubmitting}>
                  {isSubmitting ? (
                    <>
                      <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                      保存中...
                    </>
                  ) : (
                    "保存"
                  )}
                </Button>
              </div>
            </form>
          </Form>
        </CardContent>
      </Card>
    </div>
  );
}
