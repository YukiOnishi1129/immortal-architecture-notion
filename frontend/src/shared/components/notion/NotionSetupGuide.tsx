"use client";

import { ExternalLink, HelpCircle } from "lucide-react";
import { useState } from "react";
import { Button } from "@/shared/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/shared/components/ui/dialog";

const INTEGRATIONS_URL = "https://www.notion.so/my-integrations";

const STEPS = [
  {
    title: "Notionでコネクトを作る",
    body: "を開き、「New connection」を押して新規作成します。認証方法は「Access token」を選びます。表示された Access token を、アプリの環境変数に設定します。",
    link: { href: INTEGRATIONS_URL, label: "notion.so/my-integrations" },
  },
  {
    title: "置き場所にするページを用意する",
    body: "ノートを書き出したいページをNotionで作ります。このページの下に、公開したノートが1件ずつ作られます。",
  },
  {
    title: "そのページにコネクトを追加する",
    body: "ページ右上の「•••」メニューから接続（Connections）を開き、作成したコネクトを追加します。この操作を忘れると、トークンが正しくてもページが見つからないエラーになります。",
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
                  <p className="text-sm text-muted-foreground">
                    {step.link && (
                      <a
                        href={step.link.href}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="mr-1 inline-flex items-center gap-0.5 text-blue-600 hover:underline"
                      >
                        {step.link.label}
                        <ExternalLink className="h-3 w-3" />
                      </a>
                    )}
                    {step.body}
                  </p>
                </div>
              </li>
            ))}
          </ol>

          <p className="text-xs text-muted-foreground">
            Notion側の画面名は変わることがあります。うまく見つからないときは
            <a
              href={INTEGRATIONS_URL}
              target="_blank"
              rel="noopener noreferrer"
              className="mx-1 text-blue-600 hover:underline"
            >
              notion.so/my-integrations
            </a>
            を開いてください。
          </p>
        </DialogContent>
      </Dialog>
    </>
  );
}
