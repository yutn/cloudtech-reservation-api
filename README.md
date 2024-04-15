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
sudo yum update -y
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

### 4. ソースコードのインストール
gitコマンドで、ソースコードをインストールする
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

ToDo

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
