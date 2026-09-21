"use client";

import { FileText, FolderOpen, Menu, Users, X } from "lucide-react";
import type { Route } from "next";
import Link from "next/link";
import { Button } from "@/shared/components/ui/button";
import { cn } from "@/shared/lib/utils";

type NavigationItem = {
  name: string;
  href: Route;
  icon: React.ComponentType<React.SVGProps<SVGSVGElement>>;
};

const navigation: NavigationItem[] = [
  {
    name: "みんなのノート",
    href: "/notes",
    icon: Users,
  },
  {
    name: "マイノート",
    href: "/my-notes",
    icon: FileText,
  },
  {
    name: "テンプレート一覧",
    href: "/templates",
    icon: FolderOpen,
  },
];

interface SidebarPresenterProps {
  isOpen: boolean;
  pathname: string;
  onToggle: () => void;
}

export function SidebarPresenter({
  isOpen,
  pathname,
  onToggle,
}: SidebarPresenterProps) {
  return (
    <>
      {/* モバイル用メニューボタン */}
      <Button
        variant="ghost"
        size="icon"
        className="fixed top-3 left-3 z-50 md:hidden"
        onClick={onToggle}
      >
        {isOpen ? <X className="h-6 w-6" /> : <Menu className="h-6 w-6" />}
        <span className="sr-only">メニューを開く</span>
      </Button>

      {/* オーバーレイ */}
      {isOpen && (
        <button
          type="button"
          className="fixed inset-0 bg-black/50 z-40 md:hidden"
          onClick={onToggle}
          aria-label="メニューを閉じる"
        />
      )}

      {/* サイドバー */}
      <aside
        className={cn(
          "fixed top-0 left-0 z-40 h-screen w-64 transform bg-white border-r transition-transform md:relative md:translate-x-0",
          isOpen ? "translate-x-0" : "-translate-x-full",
        )}
      >
        <div className="flex h-14 items-center px-6 bg-gray-900">
          <Link href={"/notes"} className="flex items-center space-x-2">
            <span className="font-bold text-lg text-white">Mini Notion</span>
          </Link>
        </div>
        <nav className="flex flex-col h-[calc(100%-3.5rem)]">
          <div className="flex-1 space-y-1 px-3 py-4">
            {navigation.map((item) => {
              const Icon = item.icon;
              const isActive = pathname.startsWith(item.href);

              return (
                <Link
                  key={item.name}
                  href={item.href}
                  onClick={onToggle}
                  className={cn(
                    "flex items-center space-x-3 px-3 py-2 rounded-md text-sm font-medium transition-colors",
                    isActive
                      ? "bg-gray-200 text-gray-900 font-semibold"
                      : "text-gray-700 hover:bg-gray-100 hover:text-gray-900",
                  )}
                >
                  <Icon className="h-4 w-4 shrink-0" />
                  <span>{item.name}</span>
                </Link>
              );
            })}
          </div>
        </nav>
      </aside>
    </>
  );
}
