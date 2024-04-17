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

### 5. APIの起動
以下のコマンドでAPIを起動する

```shell
cd cloudtech-reservation-api
nohup go run main.go &
```

### 6. 再起動時に起動されるように設定
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
WorkingDirectory=/home/ec2-user/my-go-api
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
