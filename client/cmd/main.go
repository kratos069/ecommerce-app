package main

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/e-commerce/pb"
	"github.com/e-commerce/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

// Returns the absolute path to a file in source directory
func getLocalPath(filename string) string {
	_, currentFile, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(currentFile), filename)
}

func main() {
	// Connect to gRPC server
	conn, err := grpc.NewClient("localhost:9090",
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Connection failed: %v", err)
	}
	defer conn.Close()
	client := pb.NewEcommerceClient(conn)

	// Prepare context with a 30s timeout
	ctx, cancel := context.WithTimeout(context.Background(),
		30*time.Second)
	defer cancel()

	// send login request (from admin user)
	loginRes, err := client.LoginUser(ctx, &pb.LoginUserRequest{
		Email: "leo2@messi.com",
		Password: "secret",
	})
	if err != nil {
		log.Fatalf("cannot login: %v", err)
	}
	accessToken := loginRes.AccessToken

	// Attach our custom image content-type
	md := metadata.New(map[string]string{
		"authorization":  "Bearer " + accessToken,
		"x-content-type": "image/jpeg",
	})
	ctx = metadata.NewOutgoingContext(ctx, md)

	// Start the CreateProduct stream
	stream, err := client.CreateProduct(ctx)
	if err != nil {
		log.Fatalf("CreateProduct stream error: %v", err)
	}

	// 5) Send product info
	productInfo := &pb.ProductInfo{
		Name:          "Wireless Earbuds",
		Description:   util.RandomDescription(),
		StockQuantity: 100,
		CategoryId:    3,
		Price:         129.99,
	}
	if err := stream.Send(&pb.CreateProductRequest{
		Data: &pb.CreateProductRequest_Info{Info: productInfo},
	}); err != nil {
		log.Fatalf("Failed to send product info: %v", err)
	}

	// Open and send image in chunks
	imagePath := getLocalPath("test.jpg")
	log.Printf("Using image path: %s", imagePath)

	file, err := os.Open(imagePath)
	if err != nil {
		log.Fatalf("Cannot open image: %v", err)
	}
	defer file.Close()

	stat, _ := file.Stat()
	log.Printf("preparing upload for %s (%d KB)",
		stat.Name(), stat.Size()/1024)

	const chunkSize = 64 << 10 // 64 kb
	buf := make([]byte, chunkSize)
	chunks := 0

	for {
		n, err := file.Read(buf)
		if err != nil {
			break // EOF
		}
		if err := stream.Send(&pb.CreateProductRequest{
			Data: &pb.CreateProductRequest_ImageChunk{
				ImageChunk: &pb.ImageChunk{Data: buf[:n]},
			},
		}); err != nil {
			log.Fatalf("Failed to send chunk: %v", err)
		}
		chunks++
		if chunks%10 == 0 {
			log.Printf("Sent %d chunks (%d KB)", chunks,
				(chunks*chunkSize)/1024)
		}
	}

	// Close & receive server response
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("CreateProduct failed: %v", err)
	}

	log.Printf("Product created! ID=%d URL=%s", res.Product.Id, res.Product.ProductImage)
	log.Printf("Full response: %+v", res)
}
