package cargo

import (
	"regexp"
	"strings"

	"github.com/debricked/cli/internal/resolution/job"
	"github.com/debricked/cli/internal/resolution/pm/util"
)

const (
	cargo                      = "cargo"
	executableNotFoundErrRegex = `executable file not found`
	networkErrorRegex          = `failed to download|failed to fetch|Unable to update registry`
	invalidDependencyRegex     = `no matching package named \x60([^\x60]+)\x60`
	versionResolveErrorRegex   = `failed to select a version for the requirement \x60([^\x60]+)\x60`
	incompatibleVersionRegex   = `requires rustc ([^ ]+) or newer`
	parseErrorRegex            = `could not parse input as TOML`
)

type Job struct {
	job.BaseJob
	install      bool
	cargoCommand string
	cmdFactory   ICmdFactory
}

func NewJob(
	file string,
	install bool,
	cmdFactory ICmdFactory,
) *Job {
	return &Job{
		BaseJob:    job.NewBaseJob(file),
		install:    install,
		cmdFactory: cmdFactory,
	}
}

func (j *Job) Install() bool {
	return j.install
}

func (j *Job) Run() {
	if j.install {

		j.SendStatus("generating lockfile")
		_, err := j.runInstallCmd()
		if err != nil {
			cmdErr := util.NewPMJobError(err.Error())
			j.handleError(cmdErr)

			return
		}
	}

}

func (j *Job) runInstallCmd() ([]byte, error) {
	j.cargoCommand = cargo
	installCmd, err := j.cmdFactory.MakeInstallCmd(j.cargoCommand, j.GetFile())
	if err != nil {
		return nil, err
	}

	installCmdOutput, err := installCmd.Output()
	if err != nil {
		return nil, j.GetExitError(err, string(installCmdOutput))
	}

	return installCmdOutput, nil
}

func (j *Job) handleError(cmdErr job.IError) {
	expressions := []string{
		executableNotFoundErrRegex,
		networkErrorRegex,
		invalidDependencyRegex,
		versionResolveErrorRegex,
		incompatibleVersionRegex,
		parseErrorRegex,
	}

	for _, expression := range expressions {
		regex := regexp.MustCompile(expression)
		matches := regex.FindAllStringSubmatch(cmdErr.Error(), -1)

		if len(matches) > 0 {
			cmdErr = j.addDocumentation(expression, matches, cmdErr)
			j.Errors().Critical(cmdErr)

			return
		}
	}

	j.Errors().Critical(cmdErr)
}

func (j *Job) addDocumentation(expr string, matches [][]string, cmdErr job.IError) job.IError {
	documentation := cmdErr.Documentation()

	switch expr {
	case executableNotFoundErrRegex:
		documentation = j.GetExecutableNotFoundErrorDocumentation("Cargo")
	case networkErrorRegex:
		documentation = j.addNetworkUnreachableErrorDocumentation()
	case invalidDependencyRegex:
		documentation = j.addInvalidDependencyErrorDocumentation(matches)
	case versionResolveErrorRegex:
		documentation = j.addVersionResolveErrorDocumentation(matches)
	case incompatibleVersionRegex:
		documentation = j.addIncompatibleVersionErrorDocumentation(matches)
	case parseErrorRegex:
		documentation = j.addParseErrorDocumentation()
	}

	cmdErr.SetDocumentation(documentation)

	return cmdErr
}

func (j *Job) addNetworkUnreachableErrorDocumentation() string {
	return strings.Join(
		[]string{
			"We weren't able to retrieve one or more dependencies.",
			"Please check your Internet connection and try again,",
			"or run 'cargo fetch' to populate the local cache before running the Debricked CLI.",
		}, " ")
}

func (j *Job) addInvalidDependencyErrorDocumentation(matches [][]string) string {
	message := ""
	if len(matches) > 0 && len(matches[0]) > 1 {
		message = matches[0][1]
	}

	return strings.Join(
		[]string{
			"Couldn't find crate",
			message,
			", please make sure it is spelt correctly and exists in the crates.io registry:\n",
		}, " ")
}

func (j *Job) addVersionResolveErrorDocumentation(matches [][]string) string {
	message := ""
	if len(matches) > 0 && len(matches[0]) > 1 {
		message = matches[0][1]
	}

	return strings.Join(
		[]string{
			"Couldn't resolve version requirement for",
			message,
			", please check your Cargo.toml for conflicting version requirements:\n",
		}, " ")
}

func (j *Job) addIncompatibleVersionErrorDocumentation(matches [][]string) string {
	message := ""
	if len(matches) > 0 && len(matches[0]) > 1 {
		message = matches[0][1]
	}

	return strings.Join(
		[]string{
			"This project requires Rust version",
			message,
			"or newer. Please update your Rust installation with 'rustup update'.",
		}, " ")
}

func (j *Job) addParseErrorDocumentation() string {
	return strings.Join(
		[]string{
			"Failed to parse Cargo.toml.",
			"Please check that your Cargo.toml file is valid TOML syntax.",
		}, " ")
}
