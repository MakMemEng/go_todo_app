package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"golang.org/x/sync/errgroup"
)

func main() {
	if err := run(context.Background()); err != nil {
	log.Printf("failed to terminate server: %v", err)
	}
}

// 割り当てたいポート番号が既に利用されている場合，競合してエラーが発生
// 動的にポート番号を変更してrun関数を起動する
func run(ctx context.Context, l net.Listener) error {
	s := &http.Server{
		// 引数で受け取ったnet.Listenerを利用するので
		// Addrフィールドは指定しない
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
		}),
	}
	eg, ctx := errgroup.WithContext(ctx)
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

	<-ctx.Done()
	if err := s.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown: %+v", err)
	}
	return eg.Wait()
}