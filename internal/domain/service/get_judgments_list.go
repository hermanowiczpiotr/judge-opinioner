package service

import (
	"encoding/json"
	"errors"
	"os"
	"strconv"

	"judge-opinioner/internal/infrastructure/client"
)

type JudgmentService struct {
	httpClient client.HttpClientInterface
}

type Judgment struct {
	Id          int    `json:"id"`
	TextContent string `json:"textContent"`
	JudgeName   string `json:"judgeName"`
	CourtCases  []struct {
		CaseNumber string `json:"caseNumber"`
	} `json:"courtCases"`
}

type JudgmentsList struct {
	Items []*Judgment `json:"items"`
}

type JudgmentDataWrapper struct {
	Data *Judgment `json:"data"`
}

func NewJudgmentService(httpClient client.HttpClientInterface) JudgmentService {
	return JudgmentService{httpClient: httpClient}
}

func (service JudgmentService) GetListOfJudgments(judgeName string) ([]*Judgment, error) {
	queryParams := make(map[string]string)

	queryParams["judgeName"] = judgeName
	queryParams["sortingField"] = "JUDGMENT_DATE"
	queryParams["sortingDirection"] = "DESC"
	body, err := service.httpClient.GET(os.Getenv("SAOS_URL")+"/api/search/judgments", queryParams)
	var list JudgmentsList

	if err != nil {
		return list.Items, err
	}

	err = json.Unmarshal(body, &list)
	if err != nil {
		return list.Items, err

	}

	if len(list.Items) == 0 {
		return list.Items, errors.New("empty list")
	}

	return list.Items, nil
}

func (service JudgmentService) GetSpecificJudgment(judgment *Judgment) (*Judgment, error) {
	queryParams := make(map[string]string)

	type specificJudgmentResponse struct {
		Data *Judgment `json:"data"`
	}

	var sp specificJudgmentResponse

	body, err := service.httpClient.GET(os.Getenv("SAOS_URL")+"/api/judgments/"+strconv.Itoa(judgment.Id), queryParams)

	if err != nil {
		return nil, err

	}

	err = json.Unmarshal(body, &sp)

	if err != nil {
		return nil, err
	}

	return sp.Data, nil
}
