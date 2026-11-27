package storage

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Storage struct {
	client    *minio.Client
	bucket    string
	publicURL string
}

type Config struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	UseSSL    bool
	Bucket    string
	PublicURL string
}

// New สร้าง Storage client ที่รองรับทั้ง MinIO ในเครื่อง และ Cloud (AWS S3, Cloudflare R2)
func New(cfg Config) (*Storage, error) {
	if cfg.Endpoint == "" {
		cfg.Endpoint = "localhost:9000"
	}
	if cfg.Bucket == "" {
		cfg.Bucket = "football-avatars"
	}
	if cfg.PublicURL == "" {
		scheme := "http"
		if cfg.UseSSL {
			scheme = "https"
		}
		cfg.PublicURL = fmt.Sprintf("%s://%s/%s", scheme, cfg.Endpoint, cfg.Bucket)
	}

	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize storage client: %w", err)
	}

	s := &Storage{
		client:    client,
		bucket:    cfg.Bucket,
		publicURL: strings.TrimRight(cfg.PublicURL, "/"),
	}

	// ตรวจสอบและสร้าง Bucket อัตโนมัติ พร้อมตั้งค่านโยบาย Public Read
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.ensureBucket(ctx); err != nil {
		log.Printf("warning: storage ensureBucket: %v", err)
	}

	return s, nil
}

func (s *Storage) ensureBucket(ctx context.Context) error {
	exists, err := s.client.BucketExists(ctx, s.bucket)
	if err != nil {
		return err
	}
	if !exists {
		if err := s.client.MakeBucket(ctx, s.bucket, minio.MakeBucketOptions{}); err != nil {
			return err
		}
		log.Printf("storage: created bucket %q", s.bucket)
	}

	// กำหนด Public Read Policy เพื่อให้ Frontend/Browser เปิดดูรูปผ่าน URL ได้โดยตรง
	policy := fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": {"AWS": ["*"]},
				"Action": ["s3:GetObject"],
				"Resource": ["arn:aws:s3:::%s/*"]
			}
		]
	}`, s.bucket)

	if err := s.client.SetBucketPolicy(ctx, s.bucket, policy); err != nil {
		log.Printf("warning: storage set bucket policy: %v", err)
	}

	return nil
}

// UploadAvatar อัปโหลดไฟล์รูปภาพโปรไฟล์ พร้อมตรวจสอบขนาดและสกุลไฟล์
func (s *Storage) UploadAvatar(ctx context.Context, file *multipart.FileHeader, userID uuid.UUID) (string, error) {
	if file == nil {
		return "", fmt.Errorf("no file provided")
	}

	// จำกัดขนาดไฟล์ไม่เกิน 5MB
	const maxFileSize = 5 * 1024 * 1024
	if file.Size > maxFileSize {
		return "", fmt.Errorf("file size exceeds maximum limit of 5MB")
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExts := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".webp": "image/webp",
		".gif":  "image/gif",
	}

	contentType, ok := allowedExts[ext]
	if !ok {
		return "", fmt.Errorf("invalid image format: only jpg, jpeg, png, webp, and gif are allowed")
	}

	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// ตั้งชื่อไฟล์: avatars/<userID>-<timestamp><ext>
	objectName := fmt.Sprintf("avatars/%s-%d%s", userID.String(), time.Now().Unix(), ext)

	_, err = s.client.PutObject(ctx, s.bucket, objectName, src, file.Size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload image to storage: %w", err)
	}

	fileURL := fmt.Sprintf("%s/%s", s.publicURL, objectName)
	return fileURL, nil
}
