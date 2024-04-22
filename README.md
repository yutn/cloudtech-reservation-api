# APIサーバの設定
## 概要
EC2インスタンスに対してAPIサーバとして起動させる設定を行います。

## 前提
- EC2インスタンスssh接続などでログインしていること

## 手順

### 1. yumのアップデート
各種インストールを行うために使用するyumパッケージを最新化する
```shell
sudo yum update -y
```

### 2. Gitのインストール
ソースコードをEC2インスタンスにダウンロードするため、gitをインストールする
```shell
sudo yum install -y git
```

### 3. goのインストール
APIとして起動させるためのGoのインストールを行う
```shell
sudo yum install -y golang
```

インストール後、下記コマンドでバージョン確認を行う
```shell
go version
```

### 4. ソースコードのダウンロード
gitコマンドで、ソースコードをダウンロードする
```shell
git clone https://github.com/CloudTechOrg/cloudtech-reservation-api.git
```

### 5. 再起動時に起動されるように設定
EC2インスタンスの再起動時にAPIが起動するように設定する

まずは以下のコマンドで、systemdにファイルを作る
```shell
sudo vi /etc/systemd/system/goserver.service
```

viエディターが開かれるので、以下のコードを貼り付ける
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

ESCボタンをクリックし、:wqにて保存してファイルを閉じた後、以下のコマンドで設定を反映させる

```shell
sudo systemctl daemon-reload
sudo systemctl enable goserver.service
sudo systemctl start goserver.service
```

以上で、EC2インスタンスの再起動にAPIが起動されるようになる

### 6. リバースプロキシの設定
GoのAPIは8080ポートで起動しているため、80ポートで起動させるためにnginxのリバースプロキシを利用する。

以下のコマンドでnginxをインストールする
```shell
sudo yum install nginx
```

インストール後、下記コマンドでNginxを起動する
```shell
sudo systemctl start nginx
```

再起動時に起動されるように設定する
```
sudo systemctl enable nginx
```

`http://[api-serber-01のIPアドレス]`をブラウザに入力することで、`Weblome to nginx!`のページが開かれることを確認する

以下のコマンドで、nginxの設定ファイルを開く
```shell
sudo vi /etc/nginx/nginx.conf
```

`server { ・・・ }` の部分を、下記内容に変更
```
server {
        listen 80;
        server_name _;

        location / {
            proxy_pass http://localhost:8080;
            proxy_http_version 1.1;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection 'upgrade';
            proxy_set_header Host $host;
            proxy_cache_bypass $http_upgrade;
        }
    }
```

下記コマンドで、nginxを再起動する
```shell
sudo systemctl restart nginx
```


## 起動方法
起動したAPIをローカル（同じEC2内）で呼び出したい場合、下記のURLにアクセスする

```
http://localhost:8080
```

外部から呼び出したい場合、以下のURLにアクセスする（セキュリティグループなどのアクセス権限が許可されている前提）

```
http://[IPアドレス]:8080
```

# データベース接続
## 概要
APIサーバから、RDSに接続できるように設定する

## 前提
- APIサーバのEC2インスタンスにsshなどでログインしていること

## 手順

### 1. mysqlのインストール

パッケージマネージャーdefのインストール
```shell
sudo dnf update -y
```

MSQLのリポジトリ設定
```shell
sudo rpm -Uvh https://dev.mysql.com/get/mysql80-community-release-el8-1.noarch.rpm
```

MySQLサーバのインストール
```shell
sudo dnf install mysql-community-server -y
```

MySQLサービスの起動
```shell
sudo systemctl start mysqld
```

システム起動時にMySQLが自動起動するように設定
```shell
sudo systemctl enable mysqld
shell
```

### 2. RDSに接続

以下のコマンドで、RDSのMySQLに接続
```
mysql -h 【エンドポイント】 -P 3306 -u admin -p
```

エンドポイントは、RDSのコンソールから確認可能


### 3. データベースとテーブルの作成

以下のコマンドで、reservation_dbというデータベースを作る
```sql
create database reservation_db;
```

以下のコマンドで、Reservationsテーブルを作成する
```sql
CREATE TABLE reservation_db.Reservations (
    ID INT AUTO_INCREMENT PRIMARY KEY,
    company_name VARCHAR(255) NOT NULL,
    reservation_date DATE NOT NULL,
    number_of_people INT NOT NULL
);
```

以下のコマンドで、Reservationsテーブルに1件追加する
```sql
INSERT INTO reservation_db.Reservations (company_name, reservation_date, number_of_people)
VALUES ('株式会社テスト', '2024-04-21', 5);
```

### 5. 動作確認
以下のCURLコマンドで、データベース接続が正しく行われていることを確認する

```shell
curl http://localhost:8080/test
```