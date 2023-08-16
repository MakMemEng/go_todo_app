package main

import (
	"fmt"
	"net/http"
	"os"
)

/* テスト容易性が低い
・関数外部から中断操作ができない
・関数の戻り値がないため，出力の検証が困難
・異常状態になった際，os.Exit関数により終了してしまう
・サーバー起動のポート番号が固定されているため，
  動作確認用にサーバー起動したままテスト実行すると18080ポートが使用不可
*/
func main() {
	err := http.ListenAndServe(
		":18080",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
		}),
	)
	if err != nil {
		fmt.Printf("failed to terminate server: %v", err)
		os.Exit(1)
	}
}