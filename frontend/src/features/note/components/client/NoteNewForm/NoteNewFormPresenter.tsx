"use client";

import { Loader2 } from "lucide-react";
import type { Route } from "next";
import type { UseFormReturn } from "react-hook-form";
import type { Template } from "@/features/template/types";
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
import { TemplateSelector } from "../TemplateSelector";
import type { NoteNewFormData } from "./schema";

type NoteNewFormPresenterProps = {
  form: UseFormReturn<NoteNewFormData>;
  selectedTemplate: Template | null | undefined;
  isLoadingTemplate: boolean;
  isCreating: boolean;
  backTo?: Route;
  onSubmit: (e: React.FormEvent) => void;
  onCancel: () => void;
  onSectionContentChange: (index: number, content: string) => void;
};

export function NoteNewFormPresenter({
  form,
  selectedTemplate,
  isLoadingTemplate,
  isCreating,
  backTo,
  onSubmit,
  onCancel,
  onSectionContentChange,
}: NoteNewFormPresenterProps) {
  const sections = form.watch("sections");

  const listPath = backTo ?? "/notes";
  const listLabel = backTo === "/my-notes" ? "マイノート" : "みんなのノート";

  const breadcrumbItems = [
    {
      label: listLabel,
      href: listPath,
    },
    {
      label: "新規作成",
    },
  ];

  return (
    <div className="space-y-4">
      <Breadcrumb items={breadcrumbItems} />
      <Form {...form}>
        <form onSubmit={onSubmit} className="space-y-6">
          <Card>
            <CardHeader>
              <CardTitle>新規ノート作成</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <FormField
                control={form.control}
                name="title"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>タイトル</FormLabel>
                    <FormControl>
                      <Input
                        {...field}
                        placeholder="ノートのタイトルを入力"
                        disabled={isCreating}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="templateId"
                render={({ field }) => (
                  <FormItem>
                    <FormLabel>テンプレート</FormLabel>
                    <FormControl>
                      <TemplateSelector
                        value={field.value}
                        onChange={field.onChange}
                        disabled={isCreating}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />

              {isLoadingTemplate && (
                <div className="flex items-center justify-center p-4">
                  <Loader2 className="h-6 w-6 animate-spin" />
                </div>
              )}

              {selectedTemplate && sections.length > 0 && (
                <div className="space-y-4">
                  <h3 className="text-lg font-semibold">テンプレート項目</h3>
                  {sections.map((section, index) => (
                    <FormItem key={section.fieldId}>
                      <FormLabel>
                        {section.fieldLabel}
                        {section.isRequired && (
                          <span className="ml-1 text-destructive">*</span>
                        )}
                      </FormLabel>
                      <FormControl>
                        <Textarea
                          value={section.content}
                          onChange={(e) =>
                            onSectionContentChange(index, e.target.value)
                          }
                          placeholder={`${section.fieldLabel}を入力`}
                          disabled={isCreating}
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
              )}

              {form.formState.errors.root && (
                <p className="text-sm text-destructive">
                  {form.formState.errors.root.message}
                </p>
              )}

              <div className="flex gap-2">
                <Button
                  type="submit"
                  disabled={isCreating || !selectedTemplate}
                >
                  {isCreating ? (
                    <>
                      <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                      作成中...
                    </>
                  ) : (
                    "作成"
                  )}
                </Button>
                <Button
                  type="button"
                  variant="outline"
                  onClick={onCancel}
                  disabled={isCreating}
                >
                  キャンセル
                </Button>
              </div>
            </CardContent>
          </Card>
        </form>
      </Form>
    </div>
  );
}
