"use client";

import { format } from "date-fns";
import { ja } from "date-fns/locale";
import {
  Edit,
  ExternalLink,
  Eye,
  EyeOff,
  Loader2,
  Trash2,
  User,
} from "lucide-react";
import type { Route } from "next";
import type { Note } from "@/features/note/types";
import { ConfirmDialog } from "@/shared/components/dialog";
import { NotionSetupGuide } from "@/shared/components/notion";
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
  CardHeader,
  CardTitle,
} from "@/shared/components/ui/card";
import { Skeleton } from "@/shared/components/ui/skeleton";

type NoteDetailPresenterProps = {
  note?: Note | null;
  /** テンプレートにNotionの親ページが設定されているか。未設定なら公開できない。 */
  hasNotionParentPage?: boolean;
  isSyncingToNotion?: boolean;
  onSyncToNotion?: () => void;
  isLoading: boolean;
  isDeleting: boolean;
  isTogglingPublish: boolean;
  showDeleteDialog: boolean;
  showPublishDialog: boolean;
  isOwner?: boolean;
  backTo?: Route;
  onEdit: () => void;
  onDelete: () => void;
  onConfirmDelete: () => void;
  onCancelDelete: () => void;
  onTogglePublish: () => void;
  onConfirmPublish: () => void;
  onCancelPublish: () => void;
};

export function NoteDetailPresenter({
  note,
  hasNotionParentPage = false,
  isSyncingToNotion = false,
  onSyncToNotion,
  isLoading,
  isDeleting,
  isTogglingPublish,
  showDeleteDialog,
  showPublishDialog,
  isOwner = true,
  backTo,
  onEdit,
  onDelete,
  onConfirmDelete,
  onCancelDelete,
  onTogglePublish,
  onConfirmPublish,
  onCancelPublish,
}: NoteDetailPresenterProps) {
  const listPath = backTo ?? "/notes";
  const listLabel = backTo === "/my-notes" ? "マイノート" : "みんなのノート";

  const breadcrumbItems = [
    {
      label: listLabel,
      href: listPath,
    },
    {
      label: note?.title ?? "ノート詳細",
    },
  ];

  if (isLoading) {
    return <NoteDetailSkeleton />;
  }

  if (!note) {
    return (
      <Card>
        <CardContent className="p-8 text-center">
          <p className="text-muted-foreground">ノートが見つかりません</p>
        </CardContent>
      </Card>
    );
  }

  // 公開はNotion連携を伴うので、親ページが無いと実行できない。
  // 下書きに戻す操作は連携が不要なので止めない。
  const cannotPublish = note.status !== "Publish" && !hasNotionParentPage;

  // この機能より前に公開されたノートはNotionページを持たない。公開操作が
  // 連携のきっかけなので、状態を変えずに連携だけ実行できるようにする。
  const needsManualSync =
    note.status === "Publish" && !note.notionPageUrl && hasNotionParentPage;

  const statusBadgeVariant =
    note.status === "Publish" ? "default" : "secondary";
  const statusText = note.status === "Publish" ? "公開" : "下書き";

  return (
    <div className="space-y-4">
      <Breadcrumb items={breadcrumbItems} />
      <Card>
        <CardHeader>
          <div className="flex items-start justify-between">
            <div className="space-y-2">
              <CardTitle className="text-2xl">{note.title}</CardTitle>
              <div className="flex items-center gap-3">
                <Avatar className="w-6 h-6">
                  {note.owner.thumbnail ? (
                    <AvatarImage
                      src={note.owner.thumbnail}
                      alt={`${note.owner.firstName} ${note.owner.lastName}`}
                    />
                  ) : null}
                  <AvatarFallback className="text-xs">
                    <User className="w-3 h-3" />
                  </AvatarFallback>
                </Avatar>
                <span className="text-sm text-muted-foreground">
                  {note.owner.firstName} {note.owner.lastName}
                </span>
              </div>
              <div className="flex items-center gap-2 text-sm text-muted-foreground">
                <span>テンプレート: {note.templateName}</span>
                <Badge variant={statusBadgeVariant}>{statusText}</Badge>
              </div>

              {note.notionPageUrl && (
                <a
                  href={note.notionPageUrl}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="inline-flex items-center gap-1 text-sm text-blue-600 hover:underline"
                >
                  Notionで開く
                  <ExternalLink className="h-3.5 w-3.5" />
                </a>
              )}
            </div>
            {isOwner && (
              <div className="flex gap-2">
                <Button
                  onClick={onTogglePublish}
                  size="sm"
                  variant={note.status === "Publish" ? "secondary" : "default"}
                  disabled={isTogglingPublish || cannotPublish}
                  title={
                    cannotPublish
                      ? "テンプレートにNotionのページURLが設定されていません"
                      : undefined
                  }
                >
                  {isTogglingPublish ? (
                    <>
                      <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                      処理中...
                    </>
                  ) : note.status === "Publish" ? (
                    <>
                      <EyeOff className="h-4 w-4 mr-2" />
                      下書きに戻す
                    </>
                  ) : (
                    <>
                      <Eye className="h-4 w-4 mr-2" />
                      公開する
                    </>
                  )}
                </Button>
                <Button onClick={onEdit} size="sm" variant="outline">
                  <Edit className="h-4 w-4 mr-2" />
                  編集
                </Button>
                <Button
                  onClick={onDelete}
                  size="sm"
                  variant="outline"
                  className="text-destructive"
                >
                  <Trash2 className="h-4 w-4 mr-2" />
                  削除
                </Button>
              </div>
            )}
          </div>

          {isOwner && needsManualSync && (
            <div className="mt-4 flex items-center justify-between gap-2 rounded-md border bg-muted/40 px-3 py-2">
              <p className="text-sm text-muted-foreground">
                このノートはまだNotionに連携されていません。
              </p>
              <Button
                onClick={onSyncToNotion}
                size="sm"
                variant="outline"
                disabled={isSyncingToNotion}
              >
                {isSyncingToNotion ? (
                  <>
                    <Loader2 className="h-4 w-4 mr-2 animate-spin" />
                    連携中...
                  </>
                ) : (
                  "連携する"
                )}
              </Button>
            </div>
          )}

          {isOwner && cannotPublish && (
            <div className="mt-4 flex items-start gap-2 rounded-md border border-amber-200 bg-amber-50 px-3 py-2">
              <p className="text-sm text-amber-900">
                このノートのテンプレートに、NotionのページURLが設定されていません。
                設定するまで公開できません。
              </p>
              <NotionSetupGuide className="-my-1 shrink-0 text-amber-900 hover:bg-amber-100" />
            </div>
          )}
        </CardHeader>
        <CardContent className="space-y-6">
          {note.sections.map((section) => (
            <div key={section.id} className="space-y-3">
              <div>
                <h2 className="text-lg font-bold text-gray-900">
                  {section.fieldLabel}
                  {section.isRequired && (
                    <span className="ml-1 text-destructive">*</span>
                  )}
                </h2>
              </div>
              <div className="relative">
                <div className="bg-gray-50 border border-gray-200 rounded-md p-3 min-h-[100px] font-mono text-sm">
                  <p className="whitespace-pre-wrap leading-relaxed">
                    {section.content || (
                      <span className="text-gray-400 italic font-sans">
                        未入力
                      </span>
                    )}
                  </p>
                </div>
              </div>
            </div>
          ))}

          <div className="pt-4 border-t text-sm text-muted-foreground space-y-1">
            <div>
              作成日:{" "}
              {format(new Date(note.createdAt), "yyyy年MM月dd日 HH:mm", {
                locale: ja,
              })}
            </div>
            <div>
              更新日:{" "}
              {format(new Date(note.updatedAt), "yyyy年MM月dd日 HH:mm", {
                locale: ja,
              })}
            </div>
          </div>
        </CardContent>
      </Card>

      <ConfirmDialog
        open={showDeleteDialog}
        onOpenChange={onCancelDelete}
        title="ノートを削除しますか？"
        description="この操作は取り消すことができません。本当に削除してよろしいですか？"
        confirmLabel="削除"
        cancelLabel="キャンセル"
        onConfirm={onConfirmDelete}
        onCancel={onCancelDelete}
        isLoading={isDeleting}
        variant="destructive"
      />

      <ConfirmDialog
        open={showPublishDialog}
        onOpenChange={onCancelPublish}
        title="ノートを公開しますか？"
        description="このノートを公開すると、他のユーザーからも閲覧可能になります。"
        confirmLabel="公開"
        cancelLabel="キャンセル"
        onConfirm={onConfirmPublish}
        onCancel={onCancelPublish}
        isLoading={isTogglingPublish}
        variant="default"
      />
    </div>
  );
}

function NoteDetailSkeleton() {
  return (
    <Card>
      <CardHeader>
        <div className="space-y-2">
          <Skeleton className="h-8 w-2/3" />
          <Skeleton className="h-4 w-1/3" />
        </div>
      </CardHeader>
      <CardContent className="space-y-4">
        <div className="space-y-2">
          <Skeleton className="h-6 w-1/4" />
          <Skeleton className="h-24 w-full" />
        </div>
        <div className="space-y-2">
          <Skeleton className="h-6 w-1/4" />
          <Skeleton className="h-24 w-full" />
        </div>
      </CardContent>
    </Card>
  );
}
