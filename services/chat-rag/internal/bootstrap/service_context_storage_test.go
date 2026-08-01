package bootstrap

import (
	"strings"
	"testing"

	"github.com/zgsm-ai/chat-rag/internal/config"
)

func TestInitializeStorage_NoneDisablesBackend(t *testing.T) {
	svc := &ServiceContext{Config: config.Config{}}
	svc.Config.Log.StorageType = " none "

	if err := svc.initializeStorage(); err != nil {
		t.Fatalf("initializeStorage() error = %v", err)
	}
	if svc.StorageBackend != nil {
		t.Fatal("expected no storage backend for storageType none")
	}
}

func TestInitializeStorage_UnknownTypeListsNone(t *testing.T) {
	svc := &ServiceContext{Config: config.Config{}}
	svc.Config.Log.StorageType = "unknown"

	err := svc.initializeStorage()
	if err == nil {
		t.Fatal("expected error for unknown storage type")
	}
	want := `supported: none, disk, s3`
	if got := err.Error(); !strings.Contains(got, want) {
		t.Fatalf("initializeStorage() error = %q, want it to contain %q", got, want)
	}
}
