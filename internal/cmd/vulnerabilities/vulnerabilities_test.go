package vulnerabilities

import (
	"testing"

	"github.com/debricked/cli/internal/vulnerabilities/list"
	"github.com/debricked/cli/internal/vulnerabilities/remediate"
	"github.com/stretchr/testify/assert"
)

func TestNewVulnerabilitiesCmd(t *testing.T) {
	cmd := NewVulnerabilitiesCmd(remediate.Vulnerabilities{}, list.Vulnerabilities{})
	commands := cmd.Commands()
	nbrOfCommands := 2
	assert.Lenf(t, commands, nbrOfCommands, "failed to assert that there were %d sub commands connected", nbrOfCommands)
}

func TestPreRun(t *testing.T) {
	var vulnerabilitiesRemediator remediate.Vulnerabilities
	var vulnerabilitiesLister list.Vulnerabilities
	cmd := NewVulnerabilitiesCmd(vulnerabilitiesRemediator, vulnerabilitiesLister)
	cmd.PreRun(cmd, nil)
}
