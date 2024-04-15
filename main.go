package main

import (
	"fmt"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	// CORSヘッダーを設定
	w.Header().Set("Access-Control-Allow-Origin", "*")             // すべてのオリジンからのアクセスを許可
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS") // GETとOPTIONSメソッドを許可
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type") // 特定のヘッダーの使用を許可

	// リクエストメソッドがOPTIONSの場合は、プリフライトリクエストとして扱う
	if r.Method == "OPTIONS" {
		return // プリフライトリクエストにはステータス200で応答して、処理を終了する
	}

	// hello worldという文字列をレスポンスとして返す
	fmt.Fprintf(w, "hello world")
}

func main() {
	// /パスにアクセスがあった場合に、helloHandler関数を実行するように設定
	http.HandleFunc("/", helloHandler)

	// 8080ポートでサーバーを起動
	fmt.Println("HTTPサーバを起動しました。ポート: 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("HTTPサーバの起動に失敗しました: ", err)
	}
}
