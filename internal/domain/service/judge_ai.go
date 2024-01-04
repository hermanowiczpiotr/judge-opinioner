package service

type JudgeAI interface {
	AskAboutJudge(judgeName string, judgments []*Judgment) (string, error)
}
