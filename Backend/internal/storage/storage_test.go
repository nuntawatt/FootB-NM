package storage

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestStorageMinioConnection(t *testing.T) {
	cfg := Config{
		Endpoint:  "localhost:9000",
		AccessKey: "minioadmin",
		SecretKey: "minioadmin",
		UseSSL:    false,
		Bucket:    "football-avatars",
		PublicURL: "http://localhost:9000/football-avatars",
	}

	store, err := New(cfg)
	if err != nil {
		t.Skipf("skipping test: minio not available: %v", err)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Simulate multipart image upload
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("avatar", "test-avatar.png")
	if err != nil {
		t.Fatalf("failed to create form file: %v", err)
	}
	part.Write([]byte("fake-png-data"))
	writer.Close()

	req := httptest.NewRequest("POST", "/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	file, header, err := req.FormFile("avatar")
	if err != nil {
		t.Fatalf("failed to get form file: %v", err)
	}
	defer file.Close()

	userID := uuid.New()
	url, err := store.UploadAvatar(ctx, header, userID)
	if err != nil {
		t.Fatalf("failed to upload avatar: %v", err)
	}

	t.Logf("Upload succeeded: %s", url)
	if url == "" {
		t.Errorf("expected non-empty avatar URL")
	}
}
