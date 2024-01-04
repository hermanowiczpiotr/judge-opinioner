package ui

import (
	"net/http"

	"judge-opinioner/application/command_handler"
	"judge-opinioner/infrastructure/server"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"
)

type Controller struct {
	getJudgeOpinionCommandHandler command_handler.GetJudgeOpinionCommandHandler
}

func NewController(
	getJudgeOpinionCommandHandler command_handler.GetJudgeOpinionCommandHandler,
) Controller {
	return Controller{
		getJudgeOpinionCommandHandler: getJudgeOpinionCommandHandler,
	}
}

func (c Controller) GetJudgmentsList(ctx echo.Context, params server.GetJudgmentsListParams) error {
	list, err := c.getJudgeOpinionCommandHandler.Handle(
		command_handler.GetJudgeOpinionCommand{
			JudgeName: params.JudgeName,
		})

	if err != nil {
		log.Error(err)
		return ctx.String(http.StatusBadRequest, "Bad request")
	}

	return ctx.JSONPretty(http.StatusOK, list, "    ")
}
