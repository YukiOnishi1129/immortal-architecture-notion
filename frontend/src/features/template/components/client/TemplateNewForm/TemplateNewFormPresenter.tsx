"use client";

import type { DropResult } from "@hello-pangea/dnd";
import { DragDropContext, Draggable, Droppable } from "@hello-pangea/dnd";
import { GripVertical, Plus, Trash2 } from "lucide-react";
import type { Route } from "next";
import Link from "next/link";
import type { FieldArrayWithId, UseFormReturn } from "react-hook-form";
import { Breadcrumb } from "@/shared/components/ui/breadcrumb";
import { Button } from "@/shared/components/ui/button";
import { Card } from "@/shared/components/ui/card";
import { Checkbox } from "@/shared/components/ui/checkbox";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/shared/components/ui/form";
import { Input } from "@/shared/components/ui/input";
import type { TemplateNewFormData } from "./schema";

interface TemplateNewFormPresenterProps {
  form: UseFormReturn<TemplateNewFormData>;
  fields: FieldArrayWithId<TemplateNewFormData, "fields", "id">[];
  isCreating?: boolean;
  onSubmit: (e?: React.BaseSyntheticEvent) => Promise<void>;
  onCancel?: () => void;
  onRemoveField: (index: number) => void;
  onDragEnd: (result: DropResult) => void;
  onAddField: () => void;
}

export function TemplateNewFormPresenter({
  form,
  fields,
  isCreating = false,
  onSubmit,
  onCancel,
  onRemoveField,
  onDragEnd,
  onAddField,
}: TemplateNewFormPresenterProps) {
  const breadcrumbItems = [
    {
      label: "テンプレート",
      href: "/templates" as Route,
    },
    {
      label: "新規作成",
    },
  ];

  return (
    <div className="container mx-auto p-6 max-w-4xl space-y-4">
      <Breadcrumb items={breadcrumbItems} />
      <Card className="p-6">
        <h1 className="text-2xl font-bold mb-6">テンプレート新規作成</h1>

        <Form {...form}>
          <form onSubmit={onSubmit} className="space-y-6">
            <FormField
              control={form.control}
              name="name"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>テンプレート名</FormLabel>
                  <FormControl>
                    <Input
                      {...field}
                      placeholder="テンプレート名を入力"
                      maxLength={100}
                    />
                  </FormControl>
                  <FormMessage />
                </FormItem>
              )}
            />

            <div>
              <div className="flex justify-between items-center mb-4">
                <FormLabel>項目</FormLabel>
                <Button
                  type="button"
                  onClick={onAddField}
                  size="sm"
                  variant="outline"
                >
                  <Plus className="w-4 h-4 mr-1" />
                  項目を追加
                </Button>
              </div>

              {fields.length === 0 ? (
                <p className="text-gray-500 text-center py-8 border-2 border-dashed rounded-lg">
                  項目を追加してください
                </p>
              ) : (
                <DragDropContext onDragEnd={onDragEnd}>
                  <Droppable droppableId="fields">
                    {(provided) => (
                      <div
                        {...provided.droppableProps}
                        ref={provided.innerRef}
                        className="space-y-2"
                      >
                        {fields.map((field, index) => (
                          <Draggable
                            key={field.id}
                            draggableId={field.id}
                            index={index}
                          >
                            {(provided, snapshot) => (
                              <div
                                ref={provided.innerRef}
                                {...provided.draggableProps}
                                className={`flex items-center gap-2 p-3 bg-gray-50 rounded-lg ${
                                  snapshot.isDragging ? "shadow-lg" : ""
                                }`}
                              >
                                <div
                                  {...provided.dragHandleProps}
                                  className="cursor-move"
                                >
                                  <GripVertical className="w-5 h-5 text-gray-400" />
                                </div>

                                <FormField
                                  control={form.control}
                                  name={`fields.${index}.label`}
                                  render={({ field }) => (
                                    <FormItem className="flex-1">
                                      <FormControl>
                                        <Input
                                          {...field}
                                          placeholder="項目名を入力"
                                        />
                                      </FormControl>
                                      <FormMessage />
                                    </FormItem>
                                  )}
                                />

                                <FormField
                                  control={form.control}
                                  name={`fields.${index}.isRequired`}
                                  render={({ field }) => (
                                    <FormItem className="flex items-center space-x-2">
                                      <FormControl>
                                        <Checkbox
                                          checked={field.value}
                                          onCheckedChange={field.onChange}
                                        />
                                      </FormControl>
                                      <FormLabel className="m-0 cursor-pointer">
                                        必須
                                      </FormLabel>
                                    </FormItem>
                                  )}
                                />

                                <Button
                                  type="button"
                                  onClick={() => onRemoveField(index)}
                                  size="sm"
                                  variant="ghost"
                                  className="text-red-500 hover:text-red-700"
                                >
                                  <Trash2 className="w-4 h-4" />
                                </Button>
                              </div>
                            )}
                          </Draggable>
                        ))}
                        {provided.placeholder}
                      </div>
                    )}
                  </Droppable>
                </DragDropContext>
              )}
              {form.formState.errors.fields && (
                <p className="text-red-500 text-sm mt-2">
                  {form.formState.errors.fields.message}
                </p>
              )}
            </div>

            <div className="flex justify-between">
              <Button
                type="button"
                variant="outline"
                onClick={onCancel}
                disabled={isCreating}
                asChild
              >
                <Link href={"/templates"}>キャンセル</Link>
              </Button>
              <Button type="submit" disabled={isCreating}>
                {isCreating ? "作成中..." : "作成"}
              </Button>
            </div>
          </form>
        </Form>
      </Card>
    </div>
  );
}
