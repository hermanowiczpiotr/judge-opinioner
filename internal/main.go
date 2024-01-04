package main

import (
	"os"

	"judge-opinioner/application/command_handler"
	"judge-opinioner/domain/judgment/service"
	"judge-opinioner/infrastructure"
	"judge-opinioner/infrastructure/client"
	"judge-opinioner/infrastructure/logs"
	"judge-opinioner/ui"

	log "github.com/sirupsen/logrus"
)

func main() {
	logs.Init()
	log.Info("starting JUDGE-OPINIONER")

	openaiService, err := service.NewGoogleAiService(
		os.Getenv("GOOGLE_AI_PROJECT_ID"),
		os.Getenv("GOOGLE_AI_LOCATION"),
		os.Getenv("GOOGLE_AI_MODEL"),
	)

	controller := ui.NewController(command_handler.NewGetJudgeOpinionCommandHandler(
		service.NewJudgmentService(client.HTTPClient{}),
		openaiService,
	))

	router := infrastructure.NewRouter(controller)

	err = router.Start(":8001")

	if err != nil {
		log.Fatal(err)
	}

	log.Info("JUDGE-OPINIONER started")
}
