package api

import (
	"database/sql"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	db "github.com/e-commerce/db/sqlc"
	"github.com/e-commerce/util"
	"github.com/gin-gonic/gin"
)

func (server *Server) createProduct(ctx *gin.Context) {
	// imput by user in the form of string
	name := ctx.PostForm("name")
	description := ctx.PostForm("description")
	strStockQuantity := ctx.PostForm("stock_quantity")
	strCategoryID := ctx.PostForm("category_id")
	floatPrice := ctx.PostForm("price")

	// checking if user is not sending empty fields
	if name == "" || description == "" || strStockQuantity == "" ||
		strCategoryID == "" || floatPrice == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "name, description, category_id, stock_quantity & price are required",
		})
		return
	}

	// convert user string into db required fields
	stockQuantity, err := strconv.ParseInt(strStockQuantity, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}
	categoryID, err := strconv.ParseInt(strCategoryID, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}
	price, err := strconv.ParseFloat(floatPrice, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	// product image being stored in cloudinary
	productImage := ""
	fileHeader, _ := ctx.FormFile("product_image")
	if fileHeader != nil {
		productImage, err = uploadToCloud(ctx) // cloudinary
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errResponse(err))
			return
		}
	}

	arg := db.CreateProductParams{
		Name:          name,
		Description:   description,
		ProductImage:  productImage,
		Price:         price,
		StockQuantity: stockQuantity,
		CategoryID:    categoryID,
	}

	// save the product to database
	product, err := server.store.CreateProduct(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "product created successfully",
		"product": product,
	})
}

// Page 1 (first 50 products): LIMIT 50 OFFSET 0
// Page 2 (next 50):→ LIMIT 50 OFFSET 50
// Page 3 (next 50):→ LIMIT 50 OFFSET 100
// Offset = limit * (page - 1)
func (server *Server) listProducts(ctx *gin.Context) {
	// get query params (e.g., ?page=1&limit=50)
	pageStr := ctx.DefaultQuery("page", "1")
	limitStr := ctx.DefaultQuery("limit", "50")

	page, err := strconv.ParseInt(pageStr, 10, 32)
	if err != nil || page < 1 {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	limit, err := strconv.ParseInt(limitStr, 10, 32)
	if err != nil || limit < 1 || limit > 100 {
		ctx.JSON(http.StatusBadRequest,
			gin.H{"error": "limit must be between 1 and 100"})
		return
	}

	offset := (page - 1) * limit

	arg := db.ListProductsParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	}

	products, err := server.store.ListProducts(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":  "product listed",
		"products": products,
	})
}

type categoryIDStruct struct {
	CategoryID int64 `uri:"id" binding:"required,min=1"`
}

// get products by category id
func (server *Server) listProductsByCategoryID(ctx *gin.Context) {
	var req categoryIDStruct

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	products, err := server.store.ListProductsByCategory(ctx, req.CategoryID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":  "product listed by category_id",
		"products": products,
	})
}

type productIDStruct struct {
	ProductID int64 `uri:"id" binding:"required,min=1"`
}

// get the product by its ID
func (server *Server) getProductByID(ctx *gin.Context) {
	var req productIDStruct

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	product, err := server.store.GetProduct(ctx, req.ProductID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusInternalServerError, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "product by product_id",
		"product": product,
	})
}

// update an exisitng product using product ID
func (server *Server) updateProduct(ctx *gin.Context) {
	var req productIDStruct

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	product, err := server.store.GetProduct(ctx, req.ProductID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusInternalServerError, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	name := ctx.PostForm("name")
	description := ctx.PostForm("description")
	strStockQuantity := ctx.PostForm("stock_quantity")
	strCategoryID := ctx.PostForm("category_id")
	floatPrice := ctx.PostForm("price")

	if name == "" {
		name = product.Name
	}

	if description == "" {
		description = product.Description
	}

	var stockQuantity int64
	if strStockQuantity == "" {
		stockQuantity = product.StockQuantity
	} else {
		parsedQuantiy, err := strconv.ParseInt(strStockQuantity, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errResponse(err))
			return
		}
		stockQuantity = parsedQuantiy
	}

	var categoryID int64
	if strCategoryID == "" {
		categoryID = product.CategoryID
	} else {
		parsedID, err := strconv.ParseInt(strCategoryID, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errResponse(err))
			return
		}
		categoryID = parsedID
	}

	var price float64
	if floatPrice == "" {
		price = product.Price
	} else {
		parsedPrice, err := strconv.ParseFloat(floatPrice, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, errResponse(err))
			return
		}
		price = parsedPrice
	}

	// Optional product_image upload
	_, err = ctx.FormFile("product_image")
	var productImageURL string
	if err == nil {
		productImageURL, err = uploadToCloud(ctx)
		if err != nil {
			// uploadToCloud already handles JSON error response
			return
		}
	} else {
		// No file uploaded → use old URL
		productImageURL = product.ProductImage
	}

	arg := db.UpdateProductParams{
		ProductID:     product.ProductID,
		Name:          name,
		Description:   description,
		ProductImage:  productImageURL,
		StockQuantity: stockQuantity,
		Price:         price,
		CategoryID:    categoryID,
	}

	updatedProduct, err := server.store.UpdateProduct(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "updated product",
		"product": updatedProduct,
	})

}

// delete a product
func (server *Server) deleteProduct(ctx *gin.Context) {
	var req productIDStruct

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	product, err := server.store.GetProduct(ctx, req.ProductID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusInternalServerError, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	cloudService, err := util.NewCloudinaryService()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	// Delete image from Cloudinary
	// extract cloud public ID from URL
	if product.ProductImage != "" {
		publicID := extractPublicID(product.ProductImage)
		// fmt.Printf("Extracted Public ID: %s\n", publicID)

		if err := cloudService.DeleteImage(ctx, publicID); err != nil {
			log.Printf("Failed to delete image from Cloudinary: %v\n", err)
			ctx.JSON(http.StatusInternalServerError, errResponse(err))
			return
		}
	}

	err = server.store.DeleteProduct(ctx, product.ProductID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK,
		gin.H{"message": "product deleted"})

}

// ------------------------------------------------------------------//
// ------------------------------Helper Funcs------------------------//
// ------------------------------------------------------------------//

func uploadToCloud(ctx *gin.Context) (string, error) {
	file, err := ctx.FormFile("product_image")
	if err != nil {
		// checking file size (less than 5 mb)
		if file.Size > 5<<20 {
			ctx.JSON(http.StatusBadRequest, errResponse(err))
			return "", err
		}

		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return "", err
	}

	// check format validity
	if !containsValidFormat(file.Header.Get("Content-Type")) {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return "", err
	}

	// Upload the image locally
	err = ctx.SaveUploadedFile(file, "uploads/"+file.Filename)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return "", err
	}

	cloudService, err := util.NewCloudinaryService()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return "", err
	}

	imageUrl, err := cloudService.UploadImage(ctx, file)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return "", err
	}

	return imageUrl, nil
}

// Helper function to check if a string exists in a slice
func containsValidFormat(item string) bool {
	slice := []string{"image/png", "image/jpeg", "image/jpg", "image/gif"}

	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// Helper function to extract public ID from a Cloudinary URL
func extractPublicID(url string) string {
	parts := strings.Split(url, "/")
	lastPart := parts[len(parts)-1]
	publicID := strings.TrimSuffix(lastPart, filepath.Ext(lastPart)) // Remove file extension

	// Extract folder path if it exists
	if len(parts) > 7 { // Cloudinary path structure
		folderPath := strings.Join(parts[7:len(parts)-1], "/") // Preserve folder structure
		publicID = folderPath + "/" + publicID
	}

	return publicID
}
