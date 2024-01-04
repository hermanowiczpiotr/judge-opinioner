package service

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/vertexai/genai"
)

const (
	promptPattern = "Przygotowuję się do rozprawy. Sędzia: %s. Na podstawie podanych orzeczeń sądu, zdefiniuj: " +
		"Czego dotyczyły sprawy, skup się szczególnie na tendencja decyzji, Ocena stylu orzekania, Wnioski " +
		"Oto json z orzeczeniami: %v"
)

type GoogleAiService struct {
	generativeModel *genai.GenerativeModel
}

func NewGoogleAiService(projectID string, location string, model string) (*GoogleAiService, error) {

	gc, err := genai.NewClient(
		context.TODO(),
		projectID,
		location,
	)

	if err != nil {
		return nil, err
	}

	modelClient := gc.GenerativeModel(model)

	return &GoogleAiService{generativeModel: modelClient}, nil
}

func (openai GoogleAiService) AskAboutJudge(judgeName string, judgments []*Judgment) (string, error) {
	jsonData, err := json.Marshal(judgments)

	if err != nil {
		return "", err
	}

	prompt := genai.Text(fmt.Sprintf(promptPattern, judgeName, string(jsonData)))

	resp, err := openai.generativeModel.GenerateContent(context.Background(), prompt)

	return string(resp.Candidates[0].Content.Parts[0].(genai.Text)), nil
}
