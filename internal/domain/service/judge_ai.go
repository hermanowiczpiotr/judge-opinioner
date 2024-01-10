package service

type JudgeAI interface {
	AskAboutJudge(judgeName string, caseDescription string, judgment *Judgment) (string, error)
}
