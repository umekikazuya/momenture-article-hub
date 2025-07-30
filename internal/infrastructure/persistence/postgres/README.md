# Article Hub 永続化テスト

## テストケース概要


### 1. 記事作成 (Create)
- MinimalFields: 必須フィールドのみでの記事作成（Body, ProviderType, Link, DeletedAtはnil）
- AllFields: 全オプションフィールドを含む記事作成（各フィールド値の検証）
- InvalidData: 不正データ（例: 空タイトルや不正ステータス）での作成試行

### 2. 記事取得 (FindByID)
- Success: 既存記事のIDで取得し、内容一致を検証
- NotFound: 存在しないIDで取得し、nilが返ることを検証

### 3. 記事更新 (Update)
- Success: タイトル・ステータス等の基本的な更新
- OptionalFieldsSet: Body, ProviderType, Linkの追加・更新
- OptionalFieldsCleared: Body, ProviderType, Linkをnilにクリア（DBのNULL化）
- ChangeTitle: タイトルのみ変更
- ChangeStatus: ステータスのみ変更
- ChangeQiita: Linkのみ変更
- NotFound: 存在しないIDで更新し、エラーとなること

### 4. 記事削除 (Delete)
- Success: 論理削除（DeletedAtの設定）
- LogicalDeletion: 削除後もFindByIDで取得できるがDeletedAtが設定されていること
- NotFound: 存在しないIDで削除し、エラーとなること

### 5. 記事検索 (FindByCriteria)
- DefaultConditions: デフォルト条件で全件取得
- FilterByStatus: ステータスでフィルタリング
- FilterByProviderType: ProviderTypeでフィルタリング
- SortByUpdatedAt: UpdatedAtでソート
- Pagination: ページネーション（Page, Limit指定）
- IncludeDeleted: 削除済み記事も含めて取得
- EmptyResult: 条件に合致しない場合の空結果

### 6. 全件取得 (FindAll)
- Success: 全記事の取得（削除されていない記事のみ）

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
5. **Updateメソッドの注意**: Update時はtoModelを使わず、mapとgorm.Expr("NULL")で明示的にNULL更新すること（GORMの仕様により、構造体のnil値は自動でNULLにならないため）。
6. **テストのポイント**:
   - 各テストはトランザクション単位で独立して実行され、終了時にロールバックされる
   - オプションフィールドのクリア（nil化）はassert.Nilで検証
   - NotFound系はnilやエラーの返却を必ず検証

## テスト実行コマンド

```bash
# 永続化層テストのみ実行
go test ./internal/infrastructure/persistence/postgres/ -v

# 全テスト実行
go test ./... -v
```
