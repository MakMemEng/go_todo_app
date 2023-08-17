package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"

	"golang.org/x/sync/errgroup"
)

// 期待通りにHTTPサーバーが起動しているか
// テストコードから意図通りに終了するか

// キャンセル可能なcontext.Contextのオブジェクトを作成
// 別ゴルーチンでテスト対象のrun関数を実行してHTTPサーバーを起動
// エンドポイントに対してGETリクエストを送信
// cancel関数を実行
// *errgroup.Group.Waitメソッド経由でrun関数の戻り値を検証
// GETリクエストで取得したレスポンスボディが期待する文字列であることを検証
func TestRun(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		return run(ctx)
	})
	in := "message"
	rsp, err := http.Get("http://localhost:18080/" + in)
	if err != nil {
		t.Errorf("failed to get: %+v", err)
	}
	defer rsp.Body.Close()
	got, err := io.ReadAll(rsp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}
	// HTTPサーバーの戻り値を検証
	want := fmt.Sprintf("Hello, %s!", in)
	if string(got) != want {
		t.Errorf("want %q, but got %q", want, got)
	}

	// run関数に終了通知を送信
	cancel()
	// run関数の戻り値を検証
	if err := eg.Wait(); err != nil {
		t.Fatal(err)
	}
}