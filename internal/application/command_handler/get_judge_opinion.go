package command_handler

import "judge-opinioner/internal/domain/service"

type GetJudgeOpinionCommand struct {
	JudgeName string
}

type GetJudgeOpinionCommandHandler struct {
	judgmentService service.JudgmentService
	openaiService   service.JudgeAI
}

func NewGetJudgeOpinionCommandHandler(
	judgmentService service.JudgmentService,
	openaiService service.JudgeAI,
) GetJudgeOpinionCommandHandler {
	return GetJudgeOpinionCommandHandler{
		judgmentService: judgmentService,
		openaiService:   openaiService,
	}
}

func (h GetJudgeOpinionCommandHandler) Handle(command GetJudgeOpinionCommand) (string, error) {
	listOfJudgments, err := h.judgmentService.GetListOfJudgments(command.JudgeName)

	if err != nil {
		return "", err
	}

	opinion, err := h.openaiService.AskAboutJudge(command.JudgeName, listOfJudgments)

	if err != nil {
		return "", err
	}

	return opinion, nil
}
