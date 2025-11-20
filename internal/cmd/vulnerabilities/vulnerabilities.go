package vulnerabilities

import (
	"github.com/debricked/cli/internal/cmd/vulnerabilities/list"
	"github.com/debricked/cli/internal/cmd/vulnerabilities/remediate"
	lister "github.com/debricked/cli/internal/vulnerabilities/list"
	remediator "github.com/debricked/cli/internal/vulnerabilities/remediate"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewVulnerabilitiesCmd(
	vulnerabilityRemediator remediator.Vulnerabilities,
	vulnerabilityLister lister.Vulnerabilities,
) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vulnerabilities",
		Short: "Show vulnerabilities advice",
		Long:  "Analyze vulnerabilities and show remediation advice.",
		PreRun: func(cmd *cobra.Command, _ []string) {
			_ = viper.BindPFlags(cmd.Flags())
		},
	}

	cmd.AddCommand(remediate.NewRemediateCmd(vulnerabilityRemediator))
	cmd.AddCommand(list.NewListCmd(vulnerabilityLister))

	return cmd
}
