package service

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	"cloud.google.com/go/vertexai/genai"
	"github.com/labstack/gommon/log"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
)

type GoogleAiService struct {
	generativeModel *genai.GenerativeModel
}

func NewGoogleAiService(projectID string, location string, model string, token string) (GoogleAiService, error) {
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: token,
	})
	gc, err := genai.NewClient(
		context.TODO(),
		projectID,
		location,
		option.WithTokenSource(tokenSource),
	)

	if err != nil {
		log.Fatal(err)
	}

	modelClient := gc.GenerativeModel(model)

	log.Error(modelClient)
	return GoogleAiService{generativeModel: modelClient}, nil
}

func (openai GoogleAiService) AskAboutJudge(judgeName string, caseDescription string, judgment *Judgment) (string, error) {
	jsonData, err := json.Marshal(judgment.TextContent)
	text := removeHTML(string(jsonData))
	if err != nil {
		return "", err
	}

	prompt := genai.Text(fmt.Sprintf(
		"Przeanalizuj poniższy tekst i odpowiedz na pytania poniżej w podpunktach: %v"+
			""+
			""+
			""+
			"1.Jakie argumenty i zarzuty przedstawiały strony przeciwne?"+
			"2.Czy argumentacja i zarzuty przedstawiana przez te strony “trafiła” do sędziego i jak uzasadniał i na co powoływał się przy ocenie danego argumentu?", text))

	resp, err := openai.generativeModel.GenerateContent(context.Background(), prompt)

	if err != nil {
		return "", err
	}
	log.Info(resp)

	return string(resp.Candidates[0].Content.Parts[0].(genai.Text)), nil
}

func removeHTML(input string) string {
	// Regular expression to match HTML tags
	re := regexp.MustCompile("<[^>]*>")
	// Replace HTML tags with an empty string
	return re.ReplaceAllString(input, "")
}
