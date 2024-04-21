package main

import (
	"database/sql"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
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
	fmt.Fprintf(w, "API接続テストが成功しました")
}

func TestHandler(w http.ResponseWriter, r *http.Request) {
	// CORSヘッダーを設定
	w.Header().Set("Access-Control-Allow-Origin", "*")             // すべてのオリジンからのアクセスを許可
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS") // GETとOPTIONSメソッドを許可
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type") // 特定のヘッダーの使用を許可

	// リクエストメソッドがOPTIONSの場合は、プリフライトリクエストとして扱う
	if r.Method == "OPTIONS" {
		return // プリフライトリクエストにはステータス200で応答して、処理を終了する
	}

	reservation_count, err := database_test()

	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 正しいフォーマットでレスポンスを返す
	fmt.Fprintf(w, "データベース接続テストが成功しました（Reservationsの件数：%d）", reservation_count)
}

func database_test() (int, error) {
	// データベースに接続
	connection, err := sql.Open(
		"mysql",
		"root:CloudTech2024!@tcp(localhost:3306)/testdb")
	if err != nil {
		return 0, err
	}
	defer connection.Close()

	// SQLの実行
	rows, err := connection.Query("SELECT COUNT(*) AS reservation_count FROM Reservations")
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	// 結果の読み取り
	var reservation_count int
	for rows.Next() {
		err := rows.Scan(&reservation_count)
		if err != nil {
			return 0, err
		}
	}

	return reservation_count, nil
}

func main() {
	// /パスにアクセスがあった場合に、helloHandler関数を実行するように設定
	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/test", TestHandler)

	// 8080ポートでサーバーを起動
	fmt.Println("HTTPサーバを起動しました。ポート: 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("HTTPサーバの起動に失敗しました: ", err)
	}

}
