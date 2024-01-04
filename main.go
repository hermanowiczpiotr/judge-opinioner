package main

import (
	"os"

	"judge-opinioner/internal/application/command_handler"
	"judge-opinioner/internal/domain/service"
	"judge-opinioner/internal/infrastructure"
	"judge-opinioner/internal/infrastructure/client"
	"judge-opinioner/internal/infrastructure/logs"
	"judge-opinioner/internal/ui"

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