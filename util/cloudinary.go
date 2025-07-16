package util

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryService struct {
	cloudinary *cloudinary.Cloudinary
}

// Initialize a CloudinaryService instance
func NewCloudinaryService() (*CloudinaryService, error) {
	config, err := LoadConfig("../")
	if err != nil {
		log.Fatalln("cannot load config file:", err)
	}

	cld, err := cloudinary.NewFromParams(
		config.CloudName, config.CloudApiKey, config.CloudApiSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Cloudinary: %w", err)
	}

	return &CloudinaryService{cloudinary: cld}, nil
}

// UploadBytes uploads raw image data by wrapping it in an io.Reader.
func (cs *CloudinaryService) UploadBytes(
	ctx context.Context,
	data []byte,
) (string, error) {
	params := uploader.UploadParams{
		Folder: "ecommerce", // folder in cloudinary cloud
	}

	// Wrap byte slice in a reader
	reader := bytes.NewReader(data)

	result, err := cs.cloudinary.Upload.Upload(ctx, reader, params)
	if err != nil {
		return "", fmt.Errorf("failed to upload bytes: %w", err)
	}
	if result.SecureURL == "" {
		return "", fmt.Errorf("uploaded image URL is empty")
	}

	return result.SecureURL, nil
}

// UploadImage uploads an image from a multipart.FileHeader
func (cs *CloudinaryService) UploadImage(ctx context.Context,
	file *multipart.FileHeader) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Read file content into bytes
	data, err := io.ReadAll(src)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	return cs.UploadBytes(ctx, data)
}

// DeleteImage deletes an image by its Public ID
func (cs *CloudinaryService) DeleteImage(
	ctx context.Context,
	publicID string,
) error {
	result, err := cs.cloudinary.Upload.Destroy(ctx,
		uploader.DestroyParams{
			PublicID: publicID,
		})
	if err != nil {
		return fmt.Errorf("failed to delete image from Cloudinary: %w", err)
	}
	if result.Result == "not found" {
		return fmt.Errorf("image with publicID %s not found", publicID)
	}
	return nil
}
