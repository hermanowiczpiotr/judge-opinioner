package service

import (
	"testing"

	mocks "judge-opinioner/mocks/judge-opinioner/infrastructure/client"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockClientInterface struct {
	*mock.Mock
}

func (m mockClientInterface) GET(path string, queryParams map[string]string) ([]byte, error) {
	return nil, nil
}

func TestNewJudgmentService(t *testing.T) {
	mockItem := mockClientInterface{}
	service := NewJudgmentService(mockItem)
	assert.IsTypef(t, JudgmentService{}, service, "")
}

func TestJudgmentService_GetListOfJudgments(t *testing.T) {
	httpClient := new(mocks.MockHttpClientInterface)

	mockedParams := make(map[string]string)

	mockedParams["judgeName"] = "Judge Williams"
	mockedParams["sortingField"] = "JUDGMENT_DATE"
	mockedParams["sortingDirection"] = "DESC"

	jsonStr := `{
	   "items": [
	       {
	           "id": 1,
	           "textContent": "The court finds the defendant guilty.",
	           "judgeName": "Judge Williams"
	       },
	       {
	           "id": 2,
	           "textContent": "The case is dismissed due to lack of evidence.",
	           "judgeName": "Judge Williams"
	       },
	       {
	           "id": 3,
	           "textContent": "The court orders a retrial in light of new evidence.",
	           "judgeName": "Judge Williams"
	       }
	   ]
	}`
	jsonData := []byte(jsonStr)

	httpClient.On("GET", searchEndpoint, mockedParams).Return(jsonData, nil)

	judgmentService := JudgmentService{httpClient: httpClient}

	result, err := judgmentService.GetListOfJudgments("Judge Williams")

	assert.NoError(t, err)
	var expectedType []*Judgment
	assert.IsType(t, expectedType, result, "result should be of type []*Judgment")
	assert.NotEmpty(t, result)
}
