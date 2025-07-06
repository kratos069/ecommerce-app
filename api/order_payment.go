package api

import (
	"database/sql"
	"net/http"

	db "github.com/e-commerce/db/sqlc"
	"github.com/e-commerce/token"
	"github.com/gin-gonic/gin"
)

type orderAndPaymentTxInput struct {
	PaymentMethod string `json:"payment_method" binding:"required,oneof=stripe payoneer"`
}

func (server *Server) createOrderAndPayment(ctx *gin.Context) {
	var input orderAndPaymentTxInput

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.CreateOrderAndPaymentTxParams{
		UserID:        authPayload.UserID,
		PaymentMethod: input.PaymentMethod,
	}

	result, err := server.store.CreateOrderAndPaymentTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "order placed and payment recorded successfully",
		"order":   result.Order,
		"payment": result.Payment,
	})
}

type getPaymentUri struct {
	OrderID int64 `uri:"id" binding:"required,min=1"`
}

// get payment by order id
func (server *Server) getPaymentByOrder(ctx *gin.Context) {
	var req getPaymentUri

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	payment, err := server.store.GetPaymentByOrder(ctx, req.OrderID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, gin.H{"error": "payment not found for this order"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"payment": payment,
	})
}

type getOrderUri struct {
	OrderID int64 `uri:"id" binding:"required,min=1"`
}

// get order by id
func (server *Server) getOrder(ctx *gin.Context) {
	var uri getOrderUri

	if err := ctx.ShouldBindUri(&uri); err != nil {
		ctx.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	order, err := server.store.GetOrder(ctx, uri.OrderID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, order)
}

// list orders by user id
func (server *Server) listOrdersByUser(ctx *gin.Context) {
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	orders, err := server.store.ListOrdersByUser(ctx, authPayload.UserID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"total":  len(orders),
		"orders": orders,
	})
}
