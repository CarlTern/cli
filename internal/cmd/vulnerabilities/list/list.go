package list

import (
	"fmt"

	vulnerabilities "github.com/debricked/cli/internal/vulnerabilities"
	list "github.com/debricked/cli/internal/vulnerabilities/list"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var commitId string
var repositoryId string
var dependencyId string

const CommitFlag = "commit"
const RepositoryFlag = "repository"
const DependencyFlag = "dependency"

func NewListCmd(vulnerabilities vulnerabilities.IVulnerabilities) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List vulnerabilities.",
		Long:  `Show vulnerabilities affecting chosen repository.`,
		PreRun: func(cmd *cobra.Command, _ []string) {
			_ = viper.BindPFlags(cmd.Flags())
		},
		RunE: RunE(vulnerabilities),
	}

	cmd.Flags().StringVarP(&repositoryId, RepositoryFlag, "r", "", "The repository containing the vulnerabilities you wish to list")
	viper.MustBindEnv(RepositoryFlag)

	cmd.Flags().StringVarP(&commitId, CommitFlag, "c", "", "The commit containing the vulnerabilities you wish to list")
	viper.MustBindEnv(CommitFlag)

	cmd.Flags().StringVarP(&dependencyId, DependencyFlag, "d", "", "The dependency containing the vulnerabilities you wish to list")
	viper.MustBindEnv(DependencyFlag)

	return cmd
}

func RunE(r vulnerabilities.IVulnerabilities) func(_ *cobra.Command, args []string) error {
	return func(_ *cobra.Command, _ []string) error {
		orderArgs := list.OrderArgs{
			RepositoryID: viper.GetString(RepositoryFlag),
			CommitID:     viper.GetString(CommitFlag),
			DependencyID: viper.GetString(DependencyFlag),
		}

		if _, err := r.Order(orderArgs); err != nil {
			return fmt.Errorf("%s %s", color.RedString("⨯"), err.Error())
		}

		return nil
	}
}
