package remediate

import (
	"fmt"

	vulnerabilities "github.com/debricked/cli/internal/vulnerabilities"
	remediate "github.com/debricked/cli/internal/vulnerabilities/remediate"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var commitId string
var repositoryId string
var vulnerabilityId string

const CommitFlag = "commit"
const RepositoryFlag = "repository"
const VulnerabilityFlag = "vulnerability"

func NewRemediateCmd(remediator vulnerabilities.IVulnerabilities) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remediate",
		Short: "Show vulnerability remediation advice.",
		Long:  `Show vulnerability remediation advice for a specific vulnerability and repository.`,
		PreRun: func(cmd *cobra.Command, _ []string) {
			_ = viper.BindPFlags(cmd.Flags())
		},
		RunE: RunE(remediator),
	}

	cmd.Flags().StringVarP(&vulnerabilityId, VulnerabilityFlag, "v", "", "The ID of the vulnerability you wish to receive remediation advice for")
	_ = cmd.MarkFlagRequired(VulnerabilityFlag)
	viper.MustBindEnv(VulnerabilityFlag)

	cmd.Flags().StringVarP(&repositoryId, RepositoryFlag, "r", "", "The repository containing the vulnerabilities you want remediation advice for")
	_ = cmd.MarkFlagRequired(RepositoryFlag)
	viper.MustBindEnv(RepositoryFlag)

	cmd.Flags().StringVarP(&commitId, CommitFlag, "c", "", "The commit containing the vulnerabilities you want remediation advice for")
	viper.MustBindEnv(CommitFlag)

	return cmd
}

func RunE(r vulnerabilities.IVulnerabilities) func(_ *cobra.Command, args []string) error {
	return func(_ *cobra.Command, _ []string) error {
		orderArgs := remediate.OrderArgs{
			RepositoryID:    viper.GetString(RepositoryFlag),
			CommitID:        viper.GetString(CommitFlag),
			VulnerabilityID: viper.GetString(VulnerabilityFlag),
		}

		if _, err := r.Order(orderArgs); err != nil {
			return fmt.Errorf("%s %s", color.RedString("⨯"), err.Error())
		}

		return nil
	}
}
