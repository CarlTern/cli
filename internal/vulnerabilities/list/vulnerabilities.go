package list

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/debricked/cli/internal/client"
	internalIO "github.com/debricked/cli/internal/io"
	"github.com/debricked/cli/internal/tui"
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
	DebClient      client.IDebClient
	FileWriter     internalIO.IFileWriter
	spinnerManager tui.ISpinnerManager
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
	// Spinner logic. TODO: Simplify error handling, so we don't need multiple stops.
	if v.spinnerManager == nil {
		v.spinnerManager = tui.NewSpinnerManager("Listing vulnerabilities", "Fetching scan results")
	}
	v.spinnerManager.Start()
	spinnerMessage := ""
	spinner := v.spinnerManager.AddSpinner(spinnerMessage)

	startTime := time.Now()
	rowsPerPage := 2000
	page := 1
	allResults := []string{}
	endpoint := fmt.Sprintf("/api/1.0/open/vulnerabilities/get-vulnerabilities?rowsPerPage=%d&sortColumn=cvss&order=asc", rowsPerPage)
	if orderArgs.RepositoryID != "" {
		endpoint += fmt.Sprintf("&repositoryId=%s", orderArgs.RepositoryID)
	}
	if orderArgs.DependencyID != "" {
		endpoint += fmt.Sprintf("&dependencyId=%s", orderArgs.DependencyID)
	}
	if orderArgs.CommitID != "" {
		endpoint += fmt.Sprintf("&commitId=%s", orderArgs.CommitID)
	}
	for {
		endpointWithPage := endpoint + fmt.Sprintf("&page=%d", page)

		response, err := (v.DebClient).Get(
			endpointWithPage,
			"application/json",
		)
		if err != nil {
			v.spinnerManager.Stop()
			return "", err
		}
		defer response.Body.Close()
		if response.StatusCode == http.StatusPaymentRequired {
			v.spinnerManager.Stop()
			return "", ErrSubscription
		} else if response.StatusCode != http.StatusOK {
			v.spinnerManager.Stop()
			return "", fmt.Errorf("failed to get vulnerabilities due to status code %d", response.StatusCode)
		}

		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			v.spinnerManager.Stop()
			return "", err
		}

		vulnerabilitiesList, err := v.parseVulnerabilities(bodyBytes)
		if err != nil {
			v.spinnerManager.Stop()
			return "", err
		}

		// Print results for this page as they are fetched
		fmt.Println(vulnerabilitiesList)
		allResults = append(allResults, vulnerabilitiesList)

		// Check if we need to fetch another page
		// parseVulnerabilities returns a \n-separated string, so count lines
		numEntries := 0
		if vulnerabilitiesList != "" {
			numEntries = len(strings.Split(vulnerabilitiesList, "\n"))
		}
		if numEntries < rowsPerPage {
			break
		}
		page++
	}

	v.spinnerManager.SetSpinnerMessage(spinner, "", "done")
	spinner.Complete()
	v.spinnerManager.Stop()

	duration := time.Since(startTime)
	fmt.Printf("Endpoint call took: %s\n", duration)

	return strings.Join(allResults, "\n"), nil
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
