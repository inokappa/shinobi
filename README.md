# shinobi

## これなに

* Cognito User Pool のユーザー操作 (追加, 削除, 一覧表示) を行うツールです
* direnv や jq 等と組み合わせて利用してください

## 使い方

### インストール

https://github.com/inokappa/shinobi/releases から環境に応じたバイナリをダウンロードしてください.

```
wget https://github.com/inokappa/shinobi/releases/download/v0.0.1/shinobi_darwin_amd64 -O ~/bin/shinobi
chmod +x ~/bin/shinobi
```

### 事前準備

* Cognito User Pool を作成
* アプリクライアントを作成 (`ADMIN_NO_SRP_AUTH` を有効にすること)

### ヘルプ

```sh
$ shinobi -h
Usage of shinobi:
  -create
        ユーザーを作成.
  -delete
        ユーザーを削除.
  -email string
        User Pool に作成するユーザーのメールアドレスを指定.
  -endpoint string
        AWS API のエンドポイントを指定.
  -list
        ユーザー一覧を取得.
  -nickname string
        User Pool に作成するユーザーのニックネームを指定.
  -password string
        User Pool に作成するユーザーのパスワードを指定.
  -profile string
        Profile 名を指定.
  -region string
        Region 名を指定. (default "ap-northeast-1")
  -role string
        Role ARN を指定.
  -username string
        User Pool に作成するユーザー名を指定.
  -version
        バージョンを出力.
```

### 環境変数の設定

```sh
export COGNITO_USER_POOL_ID=ap-northeast-1_XyzXyzXyz
export COGNITO_CLIENT_ID=12abc12abc12abc12abc12abc1
```

### 登録ユーザーの一覧を取得

```sh
$ shinobi
User Pool Name: sample-user-pool User Pool ID: ap-northeast-1_XyzXyzXyz
+-----------+----------------+-----------------------+------------+---------------------+----------------------+
| USERNAME  |    NICKNAME    |         EMAIL         | USERSTATUS |   USERCREATEDATE    | USERLASTMODIFIEDDATE |
+-----------+----------------+-----------------------+------------+---------------------+----------------------+
| testtest2 | かっぱほげほげ | hogehoge1@example.com | CONFIRMED  | 2018-11-03 11:32:45 | 2018-11-03 11:32:46  |
| testtest3 | かっぱほげほげ | hogehoge1@example.com | CONFIRMED  | 2018-11-03 11:33:15 | 2018-11-03 11:33:16  |
| testtest4 | かっぱほげほげ | hogehoge1@example.com | CONFIRMED  | 2018-11-03 11:46:42 | 2018-11-03 11:46:43  |
+-----------+----------------+-----------------------+------------+---------------------+----------------------+
```

### 登録ユーザーの追加

```sh
$ shinobi -create -username=testtest5 -email=hogehoge1@example.com -nickname=かっぱほげほげ5
パスワードを入力して下さい:
パスワードをもう一度入力して下さい:
ユーザー testtest5 を作成しました.
User Pool Name: sample-user-pool User Pool ID: ap-northeast-1_XyzXyzXyz
+-----------+-----------------+-----------------------+------------+---------------------+----------------------+
| USERNAME  |    NICKNAME     |         EMAIL         | USERSTATUS |   USERCREATEDATE    | USERLASTMODIFIEDDATE |
+-----------+-----------------+-----------------------+------------+---------------------+----------------------+
| testtest5 | かっぱほげほげ5 | hogehoge1@example.com | CONFIRMED  | 2018-11-03 12:29:12 | 2018-11-03 12:29:13  |
+-----------+-----------------+-----------------------+------------+---------------------+----------------------+
```

### 登録ユーザーの削除

```sh
$ shinobi -delete -username=testtest2
User Pool Name: sample-user-pool User Pool ID: ap-northeast-1_XyzXyzXyz
+-----------+----------------+-----------------------+------------+---------------------+----------------------+
| USERNAME  |    NICKNAME    |         EMAIL         | USERSTATUS |   USERCREATEDATE    | USERLASTMODIFIEDDATE |
+-----------+----------------+-----------------------+------------+---------------------+----------------------+
| testtest2 | かっぱほげほげ | hogehoge1@example.com | CONFIRMED  | 2018-11-03 11:32:45 | 2018-11-03 11:32:46  |
+-----------+----------------+-----------------------+------------+---------------------+----------------------+
上記のユーザーを削除しますか?(y/n): y
ユーザーを削除します.
ユーザー testtest2 を削除しました.
```

## todo

* 色々
