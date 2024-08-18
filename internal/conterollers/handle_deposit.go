package conterollers

import (
	"exinity/internal/server"
	"exinity/internal/usecases/deposit"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h Handler) PostTransactionsGatewayDeposit(eCtx echo.Context, gateway server.PostTransactionsGatewayDepositParamsGateway) error {
	var request server.PostTransactionsGatewayDepositJSONRequestBody

	if err := eCtx.Bind(&request); err != nil {
		return err
	}

	ctx := eCtx.Request().Context()

	err := h.depositer.Handle(ctx, deposit.Request{
		OrderID:  request.OrderId,
		Amount:   request.Amount,
		Currency: request.Currency,
	})
	if err != nil {
		return err
	}

	return eCtx.JSON(http.StatusOK, "OK")
}
