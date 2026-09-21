import Link from "next/link";
import { Button } from "@/shared/components/ui/button";
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/shared/components/ui/card";

export default function NotFound() {
  return (
    <div className="container mx-auto px-4 py-16 max-w-2xl">
      <Card className="text-center">
        <CardHeader>
          <CardTitle className="text-4xl font-bold">404</CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <h1 className="text-2xl font-semibold">ページが見つかりません</h1>
          <p className="text-muted-foreground">
            お探しのページは存在しないか、移動した可能性があります。
          </p>
        </CardContent>
        <CardFooter className="justify-center">
          <Button asChild>
            <Link href={"/"}>ホームに戻る</Link>
          </Button>
        </CardFooter>
      </Card>
    </div>
  );
}
