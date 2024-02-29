package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"cloud.google.com/go/vertexai/genai"
	"github.com/labstack/gommon/log"
	"google.golang.org/api/option"
)

type GoogleAiService struct {
	generativeModel *genai.GenerativeModel
}

func NewGoogleAiService(projectID string, location string, model string, tokenPath string) (GoogleAiService, error) {

	serviceAccountKey, err := os.ReadFile(tokenPath)
	if err != nil {
		log.Fatalf("Unable to read service account key file: %v", err)
	}

	gc, err := genai.NewClient(
		context.TODO(),
		projectID,
		location,
		option.WithCredentialsJSON(serviceAccountKey),
	)

	if err != nil {
		log.Fatal(err)
	}

	modelClient := gc.GenerativeModel(model)

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
		log.Info(err)
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
