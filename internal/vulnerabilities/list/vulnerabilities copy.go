package list

/*
import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/debricked/cli/internal/client"
	internalIO "github.com/debricked/cli/internal/io"
	"github.com/debricked/cli/internal/vulnerabilities"
)

var (
	ErrHandleArgs   = errors.New("failed to handle args")
	ErrSubscription = errors.New("enterprise feature. Please visit https://debricked.com/pricing/ for more info")
)

type OrderArgs struct {
	RepositoryID    string
	CommitID        string
	Branch          string
	Output          string
	Format          string
	Vulnerabilities bool
	Licenses        bool
	DependencyID    string
}

type Vulnerabilities struct {
	DebClient  client.IDebClient
	FileWriter internalIO.IFileWriter
}

func (v Vulnerabilities) Order(args vulnerabilities.IOrderArgs) (string, error) {
	orderArgs, ok := args.(OrderArgs)
	//	var err error
	if !ok {
		return "", ErrHandleArgs
	}

	return v.listVulnerabilities(orderArgs)
}

func (v Vulnerabilities) listVulnerabilities(orderArgs OrderArgs) (string, error) {
	// Currently capped at 250 rows per page. TODO: implement pagination if needed.
	endpoint := "/api/1.0/open/vulnerabilities/get-vulnerabilities?rowsPerPage=20000&sortColumn=cvss&order=asc&page=1"

	// Add optional arguments if present
	if orderArgs.RepositoryID != "" {
		endpoint += fmt.Sprintf("&repositoryId=%s", orderArgs.RepositoryID)
	}
	if orderArgs.DependencyID != "" {
		endpoint += fmt.Sprintf("&dependencyId=%s", orderArgs.DependencyID)
	}
	if orderArgs.CommitID != "" {
		endpoint += fmt.Sprintf("&commitId=%s", orderArgs.CommitID)
	}

	response, err := (v.DebClient).Get(
		endpoint,
		"application/json",
	)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusPaymentRequired {
		return "", ErrSubscription
	} else if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get vulnerabilities due to status code %d", response.StatusCode)
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	vulnerabilitiesList, err := v.parseVulnerabilities(bodyBytes)
	if err != nil {
		return "", err
	}

	print(vulnerabilitiesList)

	return vulnerabilitiesList, nil
}

// Parse the response to extract cveId and cvss combinations
func (v Vulnerabilities) parseVulnerabilities(bodyBytes []byte) (string, error) {
	type CVSS struct {
		Text float64 `json:"text"`
		Type string  `json:"type"`
	}
	type Vulnerability struct {
		CVEId string `json:"cveId"`
		CVSS  *CVSS  `json:"cvss,omitempty"`
	}
	type VulnerabilitiesResponse struct {
		Vulnerabilities []Vulnerability `json:"vulnerabilities"`
	}

	var parsed VulnerabilitiesResponse
	err := json.Unmarshal(bodyBytes, &parsed)
	if err != nil {
		return "", err
	}

	var results []string
	for _, vulnerability := range parsed.Vulnerabilities {
		if vulnerability.CVEId != "" && vulnerability.CVSS != nil {
			results = append(results, fmt.Sprintf("%s: cvss=%.1f type=%s", vulnerability.CVEId, vulnerability.CVSS.Text, vulnerability.CVSS.Type))
		} else if vulnerability.CVEId != "" {
			results = append(results, fmt.Sprintf("%s: cvss=unknown type=unknown", vulnerability.CVEId))
		}
	}
	return strings.Join(results, "\n"), nil
}
*/
