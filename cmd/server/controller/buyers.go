package controller

import (
	"net/http"
	"strconv"

	"github.com/GuiTadeu/mercado-fresh-panic/internal/buyers"
	"github.com/GuiTadeu/mercado-fresh-panic/pkg/web"
	"github.com/gin-gonic/gin"
)

type createBuyersRequest struct {
	CardNumberId string `json:"card_number_id" binding:"required"`
	FirstName    string `json:"first_name" binding:"required"`
	LastName     string `json:"last_name" binding:"required"`
}

type updateBuyersRequest struct {
	CardNumberId string `json:"card_number_id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
}

type buyerController struct {
	buyerService buyers.BuyerService
}

func NewBuyerController(s buyers.BuyerService) *buyerController {
	return &buyerController{
		buyerService: s,
	}
}

func (c buyerController) Create() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req createBuyersRequest

		err := ctx.ShouldBindJSON(&req)
		if err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, web.NewResponse(http.StatusUnprocessableEntity, nil, err.Error()))
			return
		}

		buyer, err := c.buyerService.Create(req.CardNumberId, req.FirstName, req.LastName)

		if err != nil {
			status := buyerErrorHandler(err)
			ctx.JSON(status, web.NewResponse(status, nil, err.Error()))
			return
		}
		ctx.JSON(http.StatusCreated, web.NewResponse(http.StatusCreated, buyer, ""))
	}

}

func (c *buyerController) GetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		buyers, err := c.buyerService.GetAll()

		if err != nil {
			status := buyerErrorHandler(err)
			ctx.JSON(status, web.NewResponse(status, nil, err.Error()))
			return
		}

		ctx.JSON(http.StatusOK, web.NewResponse(http.StatusOK, buyers, ""))
	}

}

func (c *buyerController) Get() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)

		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		buyer, err := c.buyerService.Get(id)

		if err != nil {
			status := buyerErrorHandler(err)
			ctx.JSON(status, web.NewResponse(status, nil, err.Error()))
			return
		}

		ctx.JSON(http.StatusOK, web.NewResponse(http.StatusOK, buyer, ""))
	}
}

func (c *buyerController) CountPurchaseOrdersByBuyers() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, ok := ctx.GetQuery("id")
		if ok {
			id, _ := strconv.ParseUint(id, 10, 64)
			buyers, err := c.buyerService.CountPurchaseOrdersByBuyer(id)
			if err != nil {
				status := buyerErrorHandler(err)
				ctx.JSON(status, web.NewResponse(status, nil, err.Error()))
				return
			}
			ctx.JSON(http.StatusOK, web.NewResponse(http.StatusOK, buyers, ""))
		} else {

			buyers, err := c.buyerService.CountPurchaseOrdersByBuyers()

			if err != nil {
				status := buyerErrorHandler(err)
				ctx.JSON(status, web.NewResponse(status, nil, err.Error()))
				return
			}

			ctx.JSON(http.StatusOK, web.NewResponse(http.StatusOK, buyers, ""))
		}
	}
}

func (c *buyerController) Delete() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)

		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		err = c.buyerService.Delete(id)
		if err != nil {
			status := buyerErrorHandler(err)
			ctx.JSON(status, web.NewResponse(status, nil, err.Error()))
			return
		}

		ctx.JSON(http.StatusNoContent, nil)
	}
}

func (c *buyerController) Update() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		var req updateBuyersRequest
		err = ctx.ShouldBindJSON(&req)
		if err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"error": err.Error(),
			})
			return
		}

		buyer, err := c.buyerService.Update(id, req.CardNumberId, req.FirstName, req.LastName)
		if err != nil {
			status := buyerErrorHandler(err)
			ctx.JSON(status, web.NewResponse(status, nil, err.Error()))
			return
		}

		ctx.JSON(http.StatusOK, web.NewResponse(http.StatusOK, buyer, ""))
	}
}

func buyerErrorHandler(err error) int {
	switch err {

	case buyers.BuyerNotFoundError:
		return http.StatusNotFound

	case buyers.ExistsBuyerCardNumberIdError:
		return http.StatusConflict

	default:
		return http.StatusInternalServerError
	}
}
