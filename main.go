package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/sync/errgroup"
)

func main() {
	if err := run(context.Background()); err != nil {
	log.Printf("failed to terminate server: %v", err)
	}
}

// context.Context型の値を引数にとり，外部からのキャンセル操作を受け取った際にサーバーを終了するように実装
// 異常時にはerror型の値を返す,func run(ctx context.Context) error関数を実装
func run(ctx context.Context) error {
	// *http.Server型を経由してHTTPサーバーを起動
	s := &http.Server{
		Addr: ":18080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, %s!", r.URL.Path[:1])
		}),
	}
	eg, ctx := errgroup.WithContext(ctx)
	// 別ゴルーチンでHTTPサーバーを起動
	eg.Go(func() error {
		if err := s.ListenAndServe(); err != nil &&
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