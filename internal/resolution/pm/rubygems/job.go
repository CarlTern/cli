package rubygems

import (
	"regexp"
	"strings"

	"github.com/debricked/cli/internal/resolution/job"
	"github.com/debricked/cli/internal/resolution/pm/util"
)

const (
	bundle                     = "bundle"
	executableNotFoundErrRegex = `executable file not found`
	gemNotFoundRegex           = `Could not find gem '([^']+)'`
	versionConflictRegex       = `Bundler could not find compatible versions for gem "([^"]+)"`
	networkErrorRegex          = `Could not reach rubygems repository|Network is unreachable|Failed to connect`
	rubyVersionRegex           = `requires Ruby version ([^ ]+)`
	gemfileParseErrorRegex     = `There was an error parsing|syntax error`
)

type Job struct {
	job.BaseJob
	install       bool
	bundleCommand string
	cmdFactory    ICmdFactory
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

		j.SendStatus("installing dependencies")
		_, err := j.runInstallCmd()
		if err != nil {
			cmdErr := util.NewPMJobError(err.Error())
			j.handleError(cmdErr)

			return
		}
	}

}

func (j *Job) runInstallCmd() ([]byte, error) {
	j.bundleCommand = bundle
	installCmd, err := j.cmdFactory.MakeInstallCmd(j.bundleCommand, j.GetFile())
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
		gemNotFoundRegex,
		versionConflictRegex,
		networkErrorRegex,
		rubyVersionRegex,
		gemfileParseErrorRegex,
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
		documentation = j.GetExecutableNotFoundErrorDocumentation("Bundler")
	case gemNotFoundRegex:
		documentation = j.addGemNotFoundErrorDocumentation(matches)
	case versionConflictRegex:
		documentation = j.addVersionConflictErrorDocumentation(matches)
	case networkErrorRegex:
		documentation = j.addNetworkUnreachableErrorDocumentation()
	case rubyVersionRegex:
		documentation = j.addRubyVersionErrorDocumentation(matches)
	case gemfileParseErrorRegex:
		documentation = j.addGemfileParseErrorDocumentation()
	}

	cmdErr.SetDocumentation(documentation)

	return cmdErr
}

func (j *Job) addGemNotFoundErrorDocumentation(matches [][]string) string {
	message := ""
	if len(matches) > 0 && len(matches[0]) > 1 {
		message = matches[0][1]
	}

	return strings.Join(
		[]string{
			"Couldn't find gem",
			message,
			", please make sure it is spelt correctly and exists in the RubyGems repository:\n",
		}, " ")
}

func (j *Job) addVersionConflictErrorDocumentation(matches [][]string) string {
	message := ""
	if len(matches) > 0 && len(matches[0]) > 1 {
		message = matches[0][1]
	}

	return strings.Join(
		[]string{
			"Couldn't resolve version conflict for gem",
			message,
			", please check your Gemfile for conflicting version requirements:\n",
		}, " ")
}

func (j *Job) addNetworkUnreachableErrorDocumentation() string {
	return strings.Join(
		[]string{
			"We weren't able to retrieve one or more dependencies.",
			"Please check your Internet connection and try again.",
		}, " ")
}

func (j *Job) addRubyVersionErrorDocumentation(matches [][]string) string {
	message := ""
	if len(matches) > 0 && len(matches[0]) > 1 {
		message = matches[0][1]
	}

	return strings.Join(
		[]string{
			"This project requires Ruby version",
			message,
			"or newer. Please update your Ruby installation.",
		}, " ")
}

func (j *Job) addGemfileParseErrorDocumentation() string {
	return strings.Join(
		[]string{
			"Failed to parse Gemfile.",
			"Please check that your Gemfile has valid Ruby syntax.",
		}, " ")
}
