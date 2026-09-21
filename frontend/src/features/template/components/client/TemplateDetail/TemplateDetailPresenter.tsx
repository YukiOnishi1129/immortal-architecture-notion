"use client";

import { Calendar, Check, Edit, FileText, Trash2, User } from "lucide-react";
import type { Route } from "next";
import Link from "next/link";
import type { Template } from "@/features/template/types";
import { ConfirmDialog } from "@/shared/components/dialog";
import {
  Avatar,
  AvatarFallback,
  AvatarImage,
} from "@/shared/components/ui/avatar";
import { Badge } from "@/shared/components/ui/badge";
import { Breadcrumb } from "@/shared/components/ui/breadcrumb";
import { Button } from "@/shared/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/shared/components/ui/card";

interface TemplateDetailPresenterProps {
  template: Template;
  isOwner: boolean;
  isDeleting: boolean;
  showDeleteModal: boolean;
  onEdit: () => void;
  onCreateNote: () => void;
  onDeleteClick: () => void;
  onDeleteCancel: () => void;
  onDeleteConfirm: () => void;
}

export function TemplateDetailPresenter({
  template,
  isOwner,
  isDeleting,
  showDeleteModal,
  onEdit,
  onCreateNote,
  onDeleteClick,
  onDeleteCancel,
  onDeleteConfirm,
}: TemplateDetailPresenterProps) {
  const canDelete = isOwner && !template.isUsed;

  const breadcrumbItems = [
    {
      label: "テンプレート",
      href: "/templates" as Route,
    },
    {
      label: template.name,
    },
  ];

  return (
    <>
      <div className="space-y-4">
        <Breadcrumb items={breadcrumbItems} />
        <Card>
          <CardHeader>
            <div className="flex justify-between items-start">
              <div className="space-y-2">
                <CardTitle className="text-2xl">{template.name}</CardTitle>
                {template.owner && (
                  <div className="flex items-center gap-3">
                    <Avatar className="w-6 h-6">
                      {template.owner.thumbnail ? (
                        <AvatarImage
                          src={template.owner.thumbnail}
                          alt={`${template.owner.firstName} ${template.owner.lastName}`}
                        />
                      ) : null}
                      <AvatarFallback className="text-xs">
                        <User className="w-3 h-3" />
                      </AvatarFallback>
                    </Avatar>
                    <span className="text-sm text-muted-foreground">
                      {template.owner.firstName} {template.owner.lastName}
                    </span>
                  </div>
                )}
                {template.updatedAt && (
                  <CardDescription className="flex items-center gap-4 text-sm">
                    <span className="flex items-center gap-1">
                      <Calendar className="w-4 h-4" />
                      更新日:{" "}
                      {new Date(template.updatedAt).toLocaleDateString("ja-JP")}
                    </span>
                  </CardDescription>
                )}
              </div>
              {isOwner && (
                <div className="flex gap-2">
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={onEdit}
                    disabled={template.isUsed}
                    title={
                      template.isUsed
                        ? "ノートで使用中のため編集できません"
                        : ""
                    }
                  >
                    <Edit className="w-4 h-4 mr-1" />
                    編集
                  </Button>
                  <Button
                    variant="destructive"
                    size="sm"
                    onClick={onDeleteClick}
                    disabled={!canDelete || isDeleting}
                    title={
                      template.isUsed
                        ? "ノートで使用中のため削除できません"
                        : ""
                    }
                  >
                    <Trash2 className="w-4 h-4 mr-1" />
                    削除
                  </Button>
                </div>
              )}
            </div>
          </CardHeader>

          <CardContent className="space-y-4">
            <div>
              <h3 className="text-sm font-medium text-muted-foreground mb-3">
                項目一覧
              </h3>
              <div className="space-y-2">
                {template.fields.map((field, index) => (
                  <div
                    key={field.id}
                    className="flex items-center justify-between p-3 bg-muted/50 rounded-lg"
                  >
                    <div className="flex items-center gap-3">
                      <span className="text-sm text-muted-foreground font-medium">
                        {index + 1}
                      </span>
                      <span className="font-medium">{field.label}</span>
                    </div>
                    {field.isRequired && (
                      <Badge variant="secondary" className="text-xs">
                        必須
                      </Badge>
                    )}
                  </div>
                ))}
              </div>
            </div>

            {template.isUsed && isOwner && (
              <div className="bg-blue-50 border border-blue-200 rounded-lg p-3 flex items-start gap-2">
                <Check className="w-4 h-4 text-blue-600 mt-0.5" />
                <div className="text-sm text-blue-800">
                  <p className="font-medium">
                    このテンプレートはノートで使用されています
                  </p>
                  <p className="mt-1">編集・削除はできません。</p>
                </div>
              </div>
            )}
          </CardContent>

          <CardFooter className="flex justify-between">
            <Button variant="outline" asChild>
              <Link href={"/templates"}>戻る</Link>
            </Button>
            <Button onClick={onCreateNote}>
              <FileText className="w-4 h-4 mr-1" />
              このテンプレートでノート作成
            </Button>
          </CardFooter>
        </Card>
      </div>

      <ConfirmDialog
        open={showDeleteModal}
        onOpenChange={onDeleteCancel}
        title="テンプレートの削除"
        description="このテンプレートを削除してもよろしいですか？この操作は取り消せません。"
        confirmLabel="削除"
        cancelLabel="キャンセル"
        onConfirm={onDeleteConfirm}
        onCancel={onDeleteCancel}
        isLoading={isDeleting}
        variant="destructive"
      />
    </>
  );
}
