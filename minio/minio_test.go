package minio

import (
	"context"
	"os"
	"testing"
)

func TestDownload(t *testing.T) {
	err := Download(context.Background(), "action-pack", "start.tar")
	if err != nil {
		t.Fatal(err)
	}
	os.Remove("start.tar")
}
