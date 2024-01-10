package command_handler

import (
	"sync"

	"judge-opinioner/internal/domain/service"

	log "github.com/sirupsen/logrus"
)

type GetJudgeOpinionCommand struct {
	JudgeName       string
	caseDescription string
}

type GetJudgeOpinionCommandHandler struct {
	judgmentService service.JudgmentService
	openaiService   service.JudgeAI
}

type opinionItem struct {
	CaseNr  string `json:"case_nr"`
	Opinion string `json:"opinion"`
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

func (h GetJudgeOpinionCommandHandler) Handle(command GetJudgeOpinionCommand) ([]opinionItem, error) {

	listOfJudgments, err := h.judgmentService.GetListOfJudgments(command.JudgeName)

	var listOfOpinions []opinionItem

	if err != nil {
		return listOfOpinions, err
	}

	OpinionChan := make(chan opinionItem)
	errorChan := make(chan error)

	var wg sync.WaitGroup

	for _, judgment := range listOfJudgments {

		if len(judgment.TextContent) == 0 {
			continue
		}

		wg.Add(1)

		go func(judgment *service.Judgment) {
			defer wg.Done()
			specificJudgment, err := h.judgmentService.GetSpecificJudgment(judgment)
			if err != nil {
				errorChan <- err
			}

			opinionText, err := h.openaiService.AskAboutJudge(command.JudgeName, command.caseDescription, specificJudgment)

			if len(opinionText) == 0 {
				return
			}

			OpinionChan <- opinionItem{
				CaseNr:  judgment.CourtCases[0].CaseNumber,
				Opinion: opinionText,
			}
		}(judgment)
	}

	go func() {
		wg.Wait()
		close(OpinionChan)
		close(errorChan)
	}()

	for {
		select {
		case opinion, ok := <-OpinionChan:
			if !ok {
				OpinionChan = nil
			} else {
				listOfOpinions = append(listOfOpinions, opinion)
			}
		case err, ok := <-errorChan:
			if !ok {
				errorChan = nil
			} else {
				log.Error(err)
			}
		}

		if OpinionChan == nil && errorChan == nil {
			break // Break the loop when both channels are closed
		}
	}
	return listOfOpinions, nil

}
