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