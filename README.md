# Week2：APIサーバの設定
## 概要
この説明は、ハンズオン課題における`Week2：基本サービス`にて必要となる、APIサーバに対する設定方法を説明しています。

## 前提
- APIサーバのEC2インスタンスにssh接続などでログインしていること

## 手順

### 1. yumのアップデート
システムを最新の状態に保つためにyumパッケージをアップデートします。
```shell
sudo yum update -y
```

### 2. Gitのインストール
EC2インスタンスにソースコードをダウンロードするために、Gitをインストールします。
```shell
sudo yum install -y git
```

### 3. Goのインストール
APIサーバとして機能するGo言語をインストールします。
```shell
sudo yum install -y golang
```
インストール後、Goのバージョンを確認します。
```shell
go version
```

### 4. ソースコードのダウンロード
Gitを使用してソースコードをダウンロードします。
```shell
git clone https://github.com/CloudTechOrg/cloudtech-reservation-api.git
```

### 5. サービスの自動起動設定
システムの再起動時にもAPIが自動で起動するようにsystemdを設定します。

まずはviエディターを使用し、サービス起動時の設定ファイルを作成します。
```shell
sudo vi /etc/systemd/system/goserver.service
```
以下の内容をファイルに追記し、保存を行います。
```
[Unit]
Description=Go Server

[Service]
WorkingDirectory=/home/ec2-user/cloudtech-reservation-api
ExecStart=/usr/bin/go run main.go
User=ec2-user
Restart=always

[Install]
WantedBy=multi-user.target
```
設定を有効にし、サービスを開始します。
```shell
sudo systemctl daemon-reload
sudo systemctl enable goserver.service
sudo systemctl start goserver.service
```

### 6. リバースプロキシの設定
8080ポートで動作するGoのAPIを80ポートで利用できるように、Nginxをリバースプロキシとして設定します。
```shell
sudo yum install nginx
sudo systemctl start nginx
sudo systemctl enable nginx
```
Nginxの設定ファイルを編集し、適切なリバースプロキシ設定を行います。
```shell
sudo vi /etc/nginx/nginx.conf
```
設定を更新した後、Nginxを再起動します。
```shell
sudo systemctl restart nginx
```

## 起動方法
下記は、ローカルまたは外部からAPIを呼び出す方法です。
```
# ローカルからのアクセス
http://localhost:8080

# 外部からのアクセス
http://[IPアドレス]:8080
```

# Week3：データベース接続
## 概要
この説明は、ハンズオン課題における`Week3：データベースとストレージ`にて必要となる、APIサーバ（EC2インスタンス）からDBサーバ（RDS）に接続するための設定を説明しています。

## 前提
- APIサーバのEC2インスタンスにsshなどでログインしていること

## 手順

### 1. mysqlのインストール

パッケージマネージャー（DNF）の更新を行います。
```shell
sudo dnf update -y
```

MySQLの公式リポジトリをシステムに追加します。
```shell
sudo rpm -Uvh https://dev.mysql.com/get/mysql80-community-release-el8-1.noarch.rpm
```

MySQLサーバのインストールを行います。
```shell
sudo dnf install mysql-community-server -y
```

MySQLサービスの起動を行います。
```shell
sudo systemctl start mysqld
```

システム起動時にMySQLが自動起動するように設定します。
```shell
sudo systemctl enable mysqld
shell
```

### 2. RDSに接続

以下のコマンドで、RDSのMySQLに接続
```
mysql -h 【エンドポイント】 -P 3306 -u 【ユーザ名】 -p
```

AWS RDSインスタンスのMySQLデータベースに接続します。以下のコマンドを実行前に、適切なエンドポイントとユーザ名に置き換えてください。また、パスワードの入力を求められるため、はRDSインスタンス作成時に指定したものを入力してください。


### 3. データベースとテーブルの作成

`reservation_db`データベースを作る
```sql
create database reservation_db;
```

`Reservations`テーブルを作成する
```sql
CREATE TABLE reservation_db.Reservations (
    ID INT AUTO_INCREMENT PRIMARY KEY,
    company_name VARCHAR(255) NOT NULL,
    reservation_date DATE NOT NULL,
    number_of_people INT NOT NULL
);
```

`Reservations`テーブルにサンプルデータを1件追加する
```sql
INSERT INTO reservation_db.Reservations (company_name, reservation_date, number_of_people)
VALUES ('株式会社テスト', '2024-04-21', 5);
```

### 5. 動作確認
APIがデータベースに正しく接続できているか確認するため、以下のCURLコマンドを使用します。

```shell
curl http://localhost:8080/test
```

# Week9:コンテナー

# ローカル環境での実行方法
## Goのインストール
```shell
brew install go
```

## MySQLのインストール
```shell
brew install mysql
```

## テストデータを作成

MySQLを起動する
```shell
mysql.server start
```

MySQLにログインする
```shell
mysql -u root -p
```

`reservation_db`データベースを作る
```sql
create database reservation_db;
```

`Reservations`テーブルを作成する
```sql
CREATE TABLE reservation_db.Reservations (
    ID INT AUTO_INCREMENT PRIMARY KEY,
    company_name VARCHAR(255) NOT NULL,
    reservation_date DATE NOT NULL,
    number_of_people INT NOT NULL
);
```

## Goのモジュール作成とパッケージインストール

モジュール作成
```shell
go mod init cloudtech-reservation-api
```

パッケージのインストール
```shell
go get github.com/go-sql-driver/mysql
```

## 環境変数の作成

.envファイルの作成

ファイルの作成
```shell
touch .env
```

下記の内容を記載
```
DB_USERNAME=root
DB_PASSWORD= #MySQLのパスワードに置き換える
DB_SERVERNAME=localhost
DB_PORT=3306
DB_NAME=reservation_db
```

## ローカル環境での実行方法

Goのビルド
```shell
go run main.go
```

curlコマンドで実行
```shell
curl http://localhost:8080
curl http://localhost:8080/test
```

## Dockerによる実行方法

Dockerのインストール
https://matsuand.github.io/docs.docker.jp.onthefly/desktop/mac/install/

ビルドの実施

```shell
docker build -t cloudtech-reservation-api . 
```

イメージの起動

```shell
docker run -e DB_SERVERNAME=host.docker.internal -p 80:80 cloudtech-reservation-api
```

curlコマンドの実行
```shell
curl http://localhost
curl http://localhost/test
```

## ECRのプッシュ
ECRへのログイン
```shell
aws ecr get-login-password --region ap-northeast-1 | docker login --username AWS --password-stdin 【アカウントID】.dkr.ecr.ap-northeast-1.amazonaws.com
```

Dockerイメージのビルド（linux/x86_64環境向けにビルド）
```shell
docker build --platform linux/x86_64 -t cloudtech-reservation-api .
```

Dockerイメージのタグ付け
```shell
docker tag cloudtech-reservation-api:latest 【アカウントID】.dkr.ecr.ap-northeast-1.amazonaws.com/cloudtech-reservation-api:latest
```

Dockerイメージのプッシュ
```shell
docker push 【アカウントID】.dkr.ecr.ap-northeast-1.amazonaws.com/cloudtech-reservation-api:latest
```