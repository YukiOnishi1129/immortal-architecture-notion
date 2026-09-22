"use client";

import { HelpCircle } from "lucide-react";
import { useState } from "react";
import { Button } from "@/shared/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/shared/components/ui/dialog";

const STEPS = [
  {
    title: "Notionでインテグレーションを作る",
    body: "Notionの「設定」→「コネクト」→「インテグレーションを開発または管理する」から新規作成します。作成後に表示されるトークンを、アプリの環境変数に設定します。",
  },
  {
    title: "置き場所にするページを用意する",
    body: "ノートを書き出したいページをNotionで作ります。このページの下に、公開したノートが1件ずつ作られます。",
  },
  {
    title: "そのページをインテグレーションに共有する",
    body: "ページ右上の「•••」→「コネクト」から、作成したインテグレーションを追加します。この操作を忘れると、トークンが正しくてもページが見つからないエラーになります。",
  },
  {
    title: "ページのURLをテンプレートに設定する",
    body: "ページのURLをコピーし、テンプレートの編集画面にある「NotionのページURL」に貼り付けて保存します。",
  },
];

/**
 * Notion連携の準備手順を表示する。
 *
 * 手順3（ページの共有）を飛ばすと、トークンが正しくても404になる。
 * 画面から気づけないため、設定する場所のそばに置いている。
 */
export function NotionSetupGuide({ className }: { className?: string }) {
  const [open, setOpen] = useState(false);

  return (
    <>
      <Button
        type="button"
        variant="ghost"
        size="sm"
        className={className}
        onClick={() => setOpen(true)}
        aria-label="Notion連携の設定方法"
      >
        <HelpCircle className="h-4 w-4" />
      </Button>

      <Dialog open={open} onOpenChange={setOpen}>
        <DialogContent className="sm:max-w-lg">
          <DialogHeader>
            <DialogTitle>Notion連携の設定方法</DialogTitle>
            <DialogDescription>
              設定すると、ノートを公開したときにNotionへページが作られます。
            </DialogDescription>
          </DialogHeader>

          <ol className="space-y-4">
            {STEPS.map((step, index) => (
              <li key={step.title} className="flex gap-3">
                <span className="flex h-6 w-6 shrink-0 items-center justify-center rounded-full bg-muted text-xs font-medium">
                  {index + 1}
                </span>
                <div className="space-y-1">
                  <p className="text-sm font-medium">{step.title}</p>
                  <p className="text-sm text-muted-foreground">{step.body}</p>
                </div>
              </li>
            ))}
          </ol>
        </DialogContent>
      </Dialog>
    </>
  );
}
