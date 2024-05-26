# 仕様
基礎機能 - エンドポイントは適宜変更するかも
- `GET /` : アップロードされてるファイルの一覧
- `GET /files/:id` : アップロードされたファイルの詳細
- `GET /files/:id/download` : アップロードされたファイルのダウンロード
- `DELETE /files/:id` : アップロードされたファイルの削除
- `GET /files/new` : ファイルの新規アップロード画面
- `POST /files` : ファイルの新規アップロード

細かい仕様
- データストア : MySQL 8.4
- ファイルは静的フォルダに格納
- ファイルのメタデータを DB に保存
- 1つのファイルは最大 1MB
- 対応する MIME タイプ（ファイル種別） - [画像ファイルの種類と形式ガイド](https://developer.mozilla.org/ja/docs/Web/Media/Formats/Image_types#apng_animated_portable_network_graphics)
  - `image/jpeg` : `.png`
  - `image/png` : `.jpg`, `.jpeg`
  - `image/gif` : `.git`
- Google アカウントでユーザ認証
  - [Firebase Authentication](https://firebase.google.com/docs/auth?hl=ja)
- 各 API ではユーザ認証済みか確認する
- ファイル一覧は 1ページあたり 10件。ページャーつける

Web・モバイル用の画像をオブジェクトストレージにアップロードする管理ツールのようなイメージ
