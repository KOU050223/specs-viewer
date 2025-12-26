# specs-viewer

specで作成したMarkdownファイルを見やすく表示するWebビューアー

## 特徴

- 📁 ディレクトリツリー表示でファイル構造を把握しやすい
- 🔄 ファイル変更時の自動リロード（WebSocket使用）
- 🎨 見やすいMarkdownレンダリング
- ⚡ シングルバイナリで配布・実行が簡単（Go製）

## インストール

### ソースからビルド

```bash
git clone <repository-url>
cd specs-viewer
go build -o specs-viewer .
```

## 使い方

```bash
specs-viewer [オプション] <specディレクトリのパス>
```

### オプション

- `-port <ポート番号>`: サーバーのポート番号を指定（デフォルト: 8080）

### 例

```bash
# デフォルトポート（8080）で起動
specs-viewer ./specs

# カスタムポートで起動
specs-viewer -port 3000 ./specs
```

ブラウザで `http://localhost:8080` （または指定したポート）にアクセスすると、Markdownファイルが表示されます。

ファイルを編集すると、ブラウザが自動的に更新されます。

## 技術スタック

- **Go**: サーバーサイド
- **goldmark**: Markdownパーサー
- **fsnotify**: ファイル監視
- **gorilla/websocket**: WebSocketによるリアルタイム通信
- **embed**: HTMLテンプレートの組み込み