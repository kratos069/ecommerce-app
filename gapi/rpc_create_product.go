package gapi

import (
	"bytes"
	"io"

	db "github.com/e-commerce/db/sqlc"
	"github.com/e-commerce/pb"
	"github.com/e-commerce/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateProduct(
	stream pb.Ecommerce_CreateProductServer,
) error {
	// authorization
	ctx := stream.Context()
	authPayload, err := server.authorizeUser(ctx, []string{util.AdminRole})
	if err != nil {
		return unauthenticatedError(err)
	}

	if authPayload.Role != util.AdminRole {
		return status.Errorf(
			codes.PermissionDenied, "cannot update other user's info")
	}

	// Get custom metadata EARLY in the stream
	var contentType string
	if md, ok := metadata.FromIncomingContext(stream.Context()); ok {
		// look for our custom header
		if ct := md.Get("x-content-type"); len(ct) > 0 {
			contentType = ct[0]
			// log.Printf("Received x-content-type: %s", contentType)
		}
		// Debug: log all headers
		// log.Printf("All headers: %+v", md)
	}

	// Receive first message which must be ProductInfo
	firstReq, err := stream.Recv()
	if err != nil {
		return status.Errorf(codes.Internal,
			"failed to receive initial request: %v", err)
	}

	info := firstReq.GetInfo()
	if info == nil {
		return status.Error(codes.InvalidArgument,
			"failed to get product info")
	}

	// Validating required fields
	if info.Name == "" || info.Description == "" ||
		info.StockQuantity == 0 || info.CategoryId == 0 ||
		info.Price == 0 {
		return status.Error(
			codes.InvalidArgument,
			"name, description, category_id, stock_quantity & price are required",
		)
	}

	// Stream in image chunks (if any)
	var imgBuffer bytes.Buffer
	totalBytes := 0
	receivedImage := false

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return status.Errorf(codes.Internal,
				"failed to receive data: %v", err)
		}

		if chunk := req.GetImageChunk(); chunk != nil {
			receivedImage = true
			n, err := imgBuffer.Write(chunk.Data)
			if err != nil {
				return status.Errorf(codes.Internal,
					"failed to write chunk: %v", err)
			}
			totalBytes += n

			// Enforce size limit during streaming
			if totalBytes > 5<<20 {
				return status.Error(codes.InvalidArgument,
					"image exceeds 5MB limit")
			}
		}
	}

	// If image was sent, validate header + upload
	var imageURL string
	if receivedImage {
		if !containsValidFormat(contentType) {
			return status.Error(codes.InvalidArgument,
				"invalid x-content-type (image/png, image/jpeg, image/jpg, image/gif)")
		}

		cloudService, err := util.NewCloudinaryService()
		if err != nil {
			return status.Errorf(codes.Internal,
				"failed to initialize cloudinary: %v", err)
		}

		imageData := imgBuffer.Bytes()
		imageURL, err = cloudService.UploadBytes(stream.Context(),
			imageData)
		if err != nil {
			return status.Errorf(codes.Internal,
				"failed to upload image: %v", err)
		}
	}

	// Persist product in DB
	arg := db.CreateProductParams{
		Name:          info.Name,
		Description:   info.Description,
		ProductImage:  imageURL,
		Price:         info.Price,
		StockQuantity: info.StockQuantity,
		CategoryID:    info.CategoryId,
	}
	product, err := server.store.CreateProduct(stream.Context(), arg)
	if err != nil {
		return status.Errorf(codes.Internal,
			"failed to create product: %v", err)
	}

	return stream.SendAndClose(&pb.CreateProductResponse{
		Message: "product created successfully",
		Product: convertProduct(product),
	})
}

// convert db product to pb product
func convertProduct(dbProduct db.Product) *pb.Product {
	return &pb.Product{
		Id:            dbProduct.ProductID,
		Name:          dbProduct.Name,
		Description:   dbProduct.Description,
		ProductImage:  dbProduct.ProductImage,
		Price:         dbProduct.Price,
		StockQuantity: dbProduct.StockQuantity,
		CategoryId:    dbProduct.CategoryID,
	}
}

// Helper to validate image content-type
func containsValidFormat(item string) bool {
	valid := []string{"image/png", "image/jpeg", "image/jpg", "image/gif"}
	for _, s := range valid {
		if s == item {
			return true
		}
	}
	return false
}
