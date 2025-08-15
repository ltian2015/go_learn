package usecontext

import (
	"context"
	"testing"
	"time"

	"com.example/golearn/concurrent/usecontext/google"
)

func TestHttpRequest(t *testing.T) {
	ctx := context.Background()
	ctxCancel, cancelFunc := context.WithDeadline(ctx, <-time.After(30*time.Second))
	defer cancelFunc()
	chResult := make(chan google.Results, 1)
	chErr := make(chan error, 1)
	go func() {
		results, err := google.Search(ctxCancel, "golang")
		if err != nil {
			chErr <- err
			return
		}
		chResult <- results
	}()
	select {
	case results := <-chResult:
		t.Logf("Search results: %v", results)
	case err := <-chErr:
		t.Errorf("Search failed: %v", err)
	case <-ctx.Done():
		t.Error("Context was cancelled before search completed")
	}
}
