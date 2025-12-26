# Homebrewでのリリース手順

## 前提条件

1. GitHubリポジトリが公開されている
2. Homebrew Tap用の別リポジトリを作成（例: `homebrew-tap`）

## 初回セットアップ

### 1. Homebrew Tap リポジトリを作成

```bash
# GitHubで新しいリポジトリを作成
# リポジトリ名: homebrew-tap (必ず homebrew- で始める)
```

### 2. GitHub Personal Access Token を作成

1. GitHub Settings → Developer settings → Personal access tokens → Tokens (classic)
2. "Generate new token (classic)" をクリック
3. 以下の権限を付与：
   - `repo` (全て)
   - `write:packages`
4. トークンを生成してコピー

### 3. GitHub Secrets に登録

specs-viewer リポジトリの Settings → Secrets and variables → Actions で：

- Secret名: `HOMEBREW_TAP_GITHUB_TOKEN`
- Value: 上記で作成したトークン

### 4. .goreleaser.yaml の確認

すでに `KOU050223` で設定済みです：

```yaml
brews:
  - name: specs-viewer
    repository:
      owner: KOU050223
      name: homebrew-tap
```

## リリース方法

### タグをプッシュするだけ！

```bash
# バージョンタグを作成
git tag -a v0.1.0 -m "First release"

# タグをプッシュ
git push origin v0.1.0
```

GitHub Actionsが自動的に：
1. マルチプラットフォームバイナリをビルド
2. GitHub Releasesを作成
3. Homebrew Formulaを `homebrew-tap` リポジトリに追加

## ユーザーがインストールする方法

### Go install（最も簡単）

```bash
go install github.com/KOU050223/specs-viewer@latest
specs-viewer
```

### Homebrew

```bash
# Tapを追加
brew tap KOU050223/tap

# インストール
brew install specs-viewer

# 使用
specs-viewer
```

## ローカルでテスト

```bash
# GoReleaserをインストール
brew install goreleaser

# リリースをテスト（実際にはリリースしない）
goreleaser release --snapshot --clean
```

## トラブルシューティング

### GoReleaser のエラー

```bash
# 設定ファイルをチェック
goreleaser check
```

### Homebrew Formula が作成されない

1. `HOMEBREW_TAP_GITHUB_TOKEN` が正しく設定されているか確認
2. トークンに `repo` 権限があるか確認
3. `homebrew-tap` リポジトリが存在するか確認

## バージョニング

セマンティックバージョニングに従う：

- `v0.1.0`: 初回リリース
- `v0.1.1`: バグフィックス
- `v0.2.0`: 新機能追加
- `v1.0.0`: 安定版
