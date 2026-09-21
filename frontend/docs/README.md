# フロントエンド技術仕様書

このディレクトリには、フロントエンド開発に関する技術仕様書が含まれています。

## ドキュメント一覧

1. **[技術スタック](./01_tech_stack.md)**
   - 使用技術とライブラリの一覧
   - インフラ構成
   - 開発環境要件

2. **[アーキテクチャ設計](./02_architecture.md)**
   - 全体構成
   - レイヤー設計
   - データフロー

3. **[App Router設計ガイド](./03_app_router.md)**
   - ルーティング戦略
   - 認証別ルートグループ
   - ページコンポーネントパターン

4. **[Featuresディレクトリ設計](./04_features_structure.md)**
   - 機能モジュールの構成
   - Container/Presenterパターン
   - Server Componentsテンプレート

5. **[External Layer設計](./05_external_layer.md)**
   - 外部連携層の責務
   - CQRSパターン
   - バックエンド変更への対応力

6. **[TanStack Query実装ガイド](./06_tanstack_query.md)**
   - サーバーサイドプリフェッチ
   - クライアントサイド実装
   - パフォーマンス最適化

7. **[開発ガイド](./07_development_guide.md)**
   - 新規画面追加フロー
   - コーディング規約
   - デバッグテクニック

## クイックスタート

新しい機能を実装する場合は、以下の順序でドキュメントを参照してください：

1. アーキテクチャ設計を理解する → [02_architecture.md](./02_architecture.md)
2. 該当する機能の実装場所を決める → [04_features_structure.md](./04_features_structure.md)
3. 必要に応じて外部連携を実装 → [05_external_layer.md](./05_external_layer.md)
4. 開発ガイドラインに従って実装 → [07_development_guide.md](./07_development_guide.md)