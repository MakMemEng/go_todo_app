package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/MakMemEng/go_todo_app/config"
	"golang.org/x/sync/errgroup"
)

// os.Args変数を使用して実行時の引数でポート番号を指定
// ex: go run . 18080
func main() {
	if err := run(context.Background()); err != nil {
		log.Printf("failed to terminate server: %v", err)
		os.Exit(1)
	}
}

// context.Context型の値を引数にとり，
// 外部からのキャンセル操作を受け取った際にサーバーを終了するように実装
// 異常時にはerror型の値を返す,func run(ctx context.Context) error関数を実装

// 割り当てたいポート番号が既に利用されている場合，競合してエラーが発生
// 動的にポート番号を変更してrun関数を起動する
func run(ctx context.Context) error {
	cfg, err := config.New()
	if err != nil {
		return err
	}
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatalf("failed to listen port %d: %v", cfg.Port, err)
	}
	url := fmt.Sprintf("http://%s", l.Addr().String())
	log.Printf("start with: %v", url)
	// *http.Server型を経由してHTTPサーバーを起動
	s := &http.Server{
		// 引数で受け取ったnet.Listenerを利用するので
		// Addrフィールドは指定しない
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
		}),
	}
	eg, ctx := errgroup.WithContext(ctx)
	// 別ゴルーチンでHTTPサーバーを起動
	eg.Go(func() error {
		// ListenAndServerメソッドではなく，Serveメソッドへ変更
		if err := s.Serve(l); err != nil &&
		// http.ErrServerClosedは
		// http.Server.Shutdown()が正常に終了したことを示すので異常ではない
		err != http.ErrServerClosed {
			log.Printf("failed to close: %+v", err)
			return err
		}
		return nil
	})

	// チャネルからの終了通知を待機
	<-ctx.Done()
	if err := s.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown: %+v", err)
	}
	// Goメソッドで起動した別ゴルーチンの終了待機
	return eg.Wait()
}