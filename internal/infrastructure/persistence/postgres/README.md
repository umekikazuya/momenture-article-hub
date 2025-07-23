# Article Hub 永続化テスト

## テストケース概要

### 1. 記事作成 (Create)
- **MinimalFields**: 必須フィールドのみでの記事作成
- **AllFields**: 全オプションフィールドを含む記事作成
- **InvalidData**: 不正データでの作成試行

### 2. 記事取得 (FindByID)
- **Success**: 既存記事の取得
- **NotFound**: 存在しない記事の取得試行

### 3. 記事更新 (Update)
- **Success**: 基本的な記事更新
- **OptionalFieldsSet**: オプションフィールドの設定
- **OptionalFieldsCleared**: オプションフィールドのクリア
- **NotFound**: 存在しない記事の更新試行

### 4. 記事削除 (Delete)
- **Success**: 記事の論理削除
- **LogicalDeletion**: 論理削除の詳細確認
- **NotFound**: 存在しない記事の削除試行

### 5. 記事検索 (FindByCriteria)
- **DefaultConditions**: デフォルト条件での全件取得
- **FilterByStatus**: ステータスでのフィルタリング
- **FilterByProviderType**: プロバイダタイプでのフィルタリング
- **SortByUpdatedAt**: 更新日時でのソート
- **Pagination**: ページネーション
- **IncludeDeleted**: 削除済み記事を含む検索
- **EmptyResult**: 条件に合致しない場合の空結果

### 6. 全件取得 (FindAll)
- **Success**: 全記事の取得

## 次のステップ

1. **PostgreSQL接続設定**: 実際のデータベース接続を実装
2. **GORM設定**: データベースモデルとマッピングを実装
3. **Repository実装**: 各メソッドの実際の実装
4. **テスト実行**: 実装後にテストを実行してGREEN状態を確認
5. **リファクタリング**: 必要に応じてコードの改善

## テストユーティリティ

### ヘルパー関数
- `stringPtr(s string) *string`: 文字列のポインタを生成
- `timePtr(t time.Time) *time.Time`: 時刻のポインタを生成
- `setupTestDB(t *testing.T)`: テスト用DB接続の設定（要実装）
- `cleanupTestData(t *testing.T, repo)`: テストデータのクリーンアップ（要実装）

## 実装時の注意点

1. **NULL許容フィールド**: Body, ProviderType, Linkは全てNULL許容
2. **論理削除**: DeletedAtフィールドを使用した論理削除
3. **タイムスタンプ**: CreatedAt, UpdatedAtの自動管理
4. **トランザクション**: テストでの適切なトランザクション管理

## テスト実行コマンド

```bash
# 永続化層テストのみ実行
go test ./internal/infrastructure/persistence/postgres/ -v

# 全テスト実行
go test ./... -v
```
