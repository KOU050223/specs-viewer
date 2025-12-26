# specs-viewer

specで作成したMarkdownファイルを見やすく表示するWebビューアー

## 特徴

- 📁 ディレクトリツリー表示でファイル構造を把握しやすい
- 🔄 ファイル変更時の自動リロード（WebSocket使用）
- 🎨 見やすいMarkdownレンダリング
- ⚡ シングルバイナリで配布・実行が簡単（Go製）

## インストール

### Go install（推奨）

Go 1.21以上がインストールされている場合：

```bash
go install github.com/KOU050223/specs-viewer@latest
```

これで`specs-viewer`コマンドが使えるようになります。

### Homebrew

```bash
# Tap を追加
brew tap KOU050223/tap

# インストール
brew install specs-viewer
```

### ソースからビルド

```bash
git clone https://github.com/KOU050223/specs-viewer.git
cd specs-viewer
go build -o specs-viewer .
```

## 使い方

```bash
specs-viewer [オプション] [specディレクトリのパス...]
```

パスを省略した場合、カレントディレクトリから以下を自動検出します：
- `./specs`
- `./.specify`

**複数のディレクトリを同時に表示可能**：両方が存在する場合、サイドバーに別々のルートとして表示されます。

### オプション

- `-port <ポート番号>`: サーバーのポート番号を指定（デフォルト: 4829）

### 例

```bash
# 自動検出（./specs と ./.specify を探して両方表示）
specs-viewer

# 単一ディレクトリを指定
specs-viewer ./my-specs

# 複数のディレクトリを指定
specs-viewer ./specs ./docs ./design

# カスタムポートで起動
specs-viewer -port 3000

# ディレクトリとポートを両方指定
specs-viewer -port 3000 ./specs ./.specify
```

ブラウザで `http://localhost:4829` （または指定したポート）にアクセスすると、Markdownファイルが表示されます。

ファイルを編集すると、ブラウザが自動的に更新されます。

## 技術スタック

- **Go**: サーバーサイド
- **goldmark**: Markdownパーサー
- **fsnotify**: ファイル監視
- **gorilla/websocket**: WebSocketによるリアルタイム通信
- **embed**: HTMLテンプレートの組み込み