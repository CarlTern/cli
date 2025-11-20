package remediate

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
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
	VulnerabilityID string
}

type Vulnerabilities struct {
	DebClient  client.IDebClient
	FileWriter internalIO.IFileWriter
}

type FileInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type DepTreeResp struct {
	Trees []struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	} `json:"trees"`
}

func (r Vulnerabilities) Order(args vulnerabilities.IOrderArgs) (string, error) {
	orderArgs, ok := args.(OrderArgs)
	//	var err error
	if !ok {
		return "", ErrHandleArgs
	}

	_, nil := findDependencyLinesInFile("example_file.txt", "example_dependency") // Example usage of the new function

	return r.remediationAdvice(orderArgs)
}

// fetchFileInformation calls the files endpoint for a given vulnerability, repository, and commit.
func (r Vulnerabilities) fetchFileInformation(orderArgs OrderArgs, vulnerabilityId string) ([]FileInfo, error) {
	endpoint := fmt.Sprintf(
		"/api/1.0/open/vulnerability/%s/files?repositoryId=%s",
		vulnerabilityId,
		orderArgs.RepositoryID,
	)
	if orderArgs.CommitID != "" {
		endpoint += fmt.Sprintf("&commitId=%s", orderArgs.CommitID)
	}
	response, err := (r.DebClient).Get(
		endpoint,
		"application/json",
	)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusPaymentRequired {
		return nil, ErrSubscription
	} else if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get file information due to status code %d", response.StatusCode)
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	// Parse the JSON array and extract 'id' and 'name' values
	var files []FileInfo
	err = json.Unmarshal(bodyBytes, &files)
	if err != nil {
		return nil, err
	}

	return files, nil
}

// getDependencyTrees fetches dependency trees for each fileId and extracts top-level name and version values.
func (r Vulnerabilities) getDependencyTrees(orderArgs OrderArgs, vulnerabilityId string) (map[string][][2]string, error) {
	files, err := r.fetchFileInformation(orderArgs, vulnerabilityId)
	if err != nil {
		return nil, err
	}

	results := make(map[string][][2]string)
	for _, file := range files {
		endpoint := fmt.Sprintf(
			"/api/1.0/open/vulnerability/%s/files/%d/dependency-tree?repositoryId=%s",
			vulnerabilityId,
			file.ID,
			orderArgs.RepositoryID,
		)
		if orderArgs.CommitID != "" {
			endpoint += fmt.Sprintf("&commitId=%s", orderArgs.CommitID)
		}
		response, err := (r.DebClient).Get(
			endpoint,
			"application/json",
		)
		if err != nil {
			return nil, err
		}
		defer response.Body.Close()
		if response.StatusCode == http.StatusPaymentRequired {
			return nil, ErrSubscription
		} else if response.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("failed to get dependency tree for fileId %d due to status code %d", file.ID, response.StatusCode)
		}

		bodyBytes, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		var depTreeResp DepTreeResp
		err = json.Unmarshal(bodyBytes, &depTreeResp)
		if err != nil {
			return nil, err
		}

		var nameVersionPairs [][2]string
		for _, tree := range depTreeResp.Trees {
			nameVersionPairs = append(nameVersionPairs, [2]string{tree.Name, tree.Version})
		}
		results[file.Name] = nameVersionPairs
	}

	return results, nil
}

// Sends remediation advice based on the orderArgs
func (r Vulnerabilities) remediationAdvice(orderArgs OrderArgs) (string, error) {
	// Get the internal vulnerability ID from the CVE ID
	vulnerabilityId, err := r.getCveID(orderArgs)
	if err != nil {
		return "", err
	}
	// Get dependency trees (file.Name -> [][2]string{Name, Version})
	depTrees, err := r.getDependencyTrees(orderArgs, vulnerabilityId)
	if err != nil {
		return "", fmt.Errorf("failed to get dependency trees: %w", err)
	}

	endpoint := fmt.Sprintf(
		"/api/1.0/open/vulnerability/%s/repositories/%s/root-fixes",
		vulnerabilityId,
		orderArgs.RepositoryID,
	)
	response, err := (r.DebClient).Get(
		endpoint,
		"application/json",
		//      bytes.NewBuffer(body), Not yet supported
	)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusPaymentRequired {
		return "", ErrSubscription
	} else if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get remediation advice due to status code %d", response.StatusCode)
	}
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	// Parse the JSON response for root fixes
	type RemediationResponse struct {
		RootFixesCount int               `json:"rootFixesCount"`
		Fixes          map[string]string `json:"fixes"`
		IsReady        bool              `json:"isReady"`
	}

	var parsed RemediationResponse
	err = json.Unmarshal(bodyBytes, &parsed)
	if err != nil {
		return "", err
	}

	//output := fmt.Sprintf("Root fixes count: %d\n", parsed.RootFixesCount)

	var adviceList []string
	for key, fix := range parsed.Fixes {
		// key format: dependency#currentVersion
		depName := key
		currVersion := ""
		if idx := strings.LastIndex(key, "#"); idx != -1 {
			depName = key[:idx]
			currVersion = key[idx+1:]
		}
		/*
				majorChange := ""
				// Check for major version change if both versions are semantic
				if fix != "reinstall_dependency" && currVersion != "" {
					currParts := strings.Split(currVersion, ".")
					fixParts := strings.Split(fix, ".")
					if len(currParts) > 0 && len(fixParts) > 0 && currParts[0] != fixParts[0] {
						majorChange = " (major version change!)"
					}
				}
				output += fmt.Sprintf("Dependency: %s, Current: %s, Suggested: %s%s\n", depName, currVersion, fix, majorChange)
			}

			print(output)
		*/

		// For each file, check if depName matches any tree.Name
		for fileName, trees := range depTrees {
			for _, pair := range trees {
				treeName := pair[0]
				if treeName == depName {
					adviceList = append(adviceList,
						fmt.Sprintf("Update %s in %s from %s to %s", depName, fileName, currVersion, fix),
					)
				}
			}
		}
	}
	if len(adviceList) == 0 {
		fmt.Println("No remediation advice found matching dependencies in files.")
	} else {
		fmt.Println("Remediation Advice:")
		for _, advice := range adviceList {
			fmt.Println(advice)
		}
	}

	return string(bodyBytes), err
}

func (r Vulnerabilities) getCveID(orderArgs OrderArgs) (string, error) {
	endpoint := fmt.Sprintf(
		"/api/1.0/public/vulnerability-database/%s/get-vulnerability-id",
		orderArgs.VulnerabilityID,
	)
	response, err := (r.DebClient).Get(
		endpoint,
		"application/json",
	)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(bodyBytes), err
}

// findDependencyLinesInFile searches for depName in fileName in the current directory and returns matching lines as JSON.
func findDependencyLinesInFile(fileName, depName string) (string, error) {
	type Match struct {
		FileName string `json:"fileName"`
		Line     string `json:"line"`
		LineNum  int    `json:"lineNum"`
	}
	var matches []Match

	file, err := os.Open(fileName)
	if err != nil {
		return "", fmt.Errorf("could not open file %s: %w", fileName, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 1
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, depName) {
			matches = append(matches, Match{
				FileName: fileName,
				Line:     line,
				LineNum:  lineNum,
			})
		}
		lineNum++
	}
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading file %s: %w", fileName, err)
	}

	for _, match := range matches {
		fmt.Printf("%s:%d: %s\n", match.FileName, match.LineNum, match.Line)
	}
	jsonBytes, err := json.MarshalIndent(matches, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error marshaling matches to JSON: %w", err)
	}
	return string(jsonBytes), nil
}
