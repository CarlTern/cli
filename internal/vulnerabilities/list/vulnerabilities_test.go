package list

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/debricked/cli/internal/client/testdata"
	"github.com/stretchr/testify/assert"
)

func TestOrderBadArgs(t *testing.T) {
	debClientMock := &testdata.DebClientMock{}
	vulnerabilities := Vulnerabilities{DebClient: debClientMock}
	args := struct{}{}

	_, err := vulnerabilities.Order(args)

	assert.ErrorIs(t, err, ErrHandleArgs)
}

func TestOrderUnauthorized(t *testing.T) {
	debClientMock := testdata.NewDebClientMock()
	errorAssertion := errors.New("unauthorized")
	debClientMock.AddMockResponse(testdata.MockResponse{Error: errorAssertion})
	vulnerabilities := Vulnerabilities{DebClient: debClientMock}
	args := OrderArgs{RepositoryID: ""}

	_, err := vulnerabilities.Order(args)

	assert.ErrorIs(t, err, errorAssertion)
}

func TestOrderForbidden(t *testing.T) {
	debClientMock := testdata.NewDebClientMock()
	debClientMock.AddMockResponse(testdata.MockResponse{StatusCode: http.StatusForbidden})
	vulnerabilities := Vulnerabilities{DebClient: debClientMock}
	args := OrderArgs{RepositoryID: ""}

	_, err := vulnerabilities.Order(args)

	assert.ErrorIs(t, err, ErrSubscription)
}

func TestOrderNotOkResponse(t *testing.T) {
	debClientMock := testdata.NewDebClientMock()
	debClientMock.AddMockResponse(testdata.MockResponse{StatusCode: http.StatusTeapot})
	vulnerabilities := Vulnerabilities{DebClient: debClientMock}
	args := OrderArgs{RepositoryID: ""}

	_, err := vulnerabilities.Order(args)

	assert.ErrorContains(t, err, fmt.Sprintf("failed to get vulnerabilities due to status code %d", http.StatusTeapot))
}

func TestOrder(t *testing.T) {
	debClientMock := testdata.NewDebClientMock()
	debClientMock.AddMockResponse(testdata.MockResponse{StatusCode: http.StatusOK})
	vulnerabilities := Vulnerabilities{DebClient: debClientMock}
	args := OrderArgs{RepositoryID: ""}

	_, err := vulnerabilities.Order(args)

	assert.NoError(t, err)
}
