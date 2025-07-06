package api

import (
	"database/sql"
	"net/http"

	db "github.com/e-commerce/db/sqlc"
	"github.com/e-commerce/token"
	"github.com/gin-gonic/gin"
)

type addToCartInputParams struct {
	ProductID int64 `json:"product_id" binding:"required,min=1"`
	Quantity  int64 `json:"quantity" binding:"required,min=1"`
}

// add items to the cart
func (server *Server) addToCart(ctx *gin.Context) {
	var req addToCartInputParams

	if err := ctx.ShouldBindJSON(&req); err != nil {
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

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.AddToCartTxParams{
		UserID:    authPayload.UserID,
		ProductID: product.ProductID,
		Quantity:  req.Quantity,
	}

	cartItem, err := server.store.AddToCartTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":   "added to cart",
		"cart_item": cartItem,
	})
}

// get cart items by user
func (server *Server) getCartItemsByUser(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	cartItems, err := server.store.GetCartItemsByUser(
		ctx, authPayload.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusInternalServerError, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":    "items in user's cart",
		"cart_items": cartItems,
	})
}

type cartItemIDStruct struct {
	CartItemID int64 `uri:"id" binding:"required,min=1"`
}

// delete items in user's cart
func (server *Server) deleteItemInCart(ctx *gin.Context) {
	var req cartItemIDStruct

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	_, err := server.store.GetCartItemsByUser(
		ctx, authPayload.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusInternalServerError, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	err = server.store.DeleteCartItem(ctx, req.CartItemID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "deleted successfully",
	})
}
