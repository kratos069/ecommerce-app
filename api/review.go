package api

import (
	"net/http"

	db "github.com/e-commerce/db/sqlc"
	"github.com/e-commerce/token"
	"github.com/gin-gonic/gin"
)

type createReviewInput struct {
	ProductID int64  `json:"product_id" binding:"required,min=1"`
	Rating    int64  `json:"rating" binding:"required,min=1,max=5"`
	Comment   string `json:"comment" binding:"required"`
}

func (server *Server) createReview(ctx *gin.Context) {
	var input createReviewInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	product, err := server.store.GetProduct(ctx, input.ProductID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	// Step 3: Create the review
	arg := db.CreateReviewParams{
		UserID:    authPayload.UserID,
		ProductID: product.ProductID,
		Rating:    input.Rating,
		Comment:   input.Comment,
	}

	review, err := server.store.CreateReview(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "review added successfully",
		"review":  review,
	})
}

type productIDUri struct {
	ProductID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getReviewsForProduct(ctx *gin.Context) {
	var uri productIDUri

	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	reviews, err := server.store.GetReviewsForProduct(ctx, uri.ProductID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "reviews for product",
		"reviews": reviews,
	})
}
