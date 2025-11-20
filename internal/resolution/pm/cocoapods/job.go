package cocoapods

import (
	"regexp"
	"strings"

	"github.com/debricked/cli/internal/resolution/job"
	"github.com/debricked/cli/internal/resolution/pm/util"
)

const (
	pod                        = "pod"
	executableNotFoundErrRegex = `executable file not found`
	specsRepoNotFoundRegex     = `Unable to find a specification`
	networkErrorRegex          = `Failed to connect to|Connection refused|Network is unreachable`
	invalidPodNameRegex        = `Unable to find a pod with name.*\x60([^\x60]+)\x60`
	versionConflictRegex       = `Unable to satisfy the following requirements.*\x60([^\x60]+)\x60`
	xcodeVersionRegex          = `requires a higher minimum deployment target`
)

type Job struct {
	job.BaseJob
	install    bool
	podCommand string
	cmdFactory ICmdFactory
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
	j.podCommand = pod
	installCmd, err := j.cmdFactory.MakeInstallCmd(j.podCommand, j.GetFile())
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
		specsRepoNotFoundRegex,
		networkErrorRegex,
		invalidPodNameRegex,
		versionConflictRegex,
		xcodeVersionRegex,
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
		documentation = j.GetExecutableNotFoundErrorDocumentation("CocoaPods")
	case specsRepoNotFoundRegex:
		documentation = j.addSpecsRepoNotFoundErrorDocumentation()
	case networkErrorRegex:
		documentation = j.addNetworkUnreachableErrorDocumentation()
	case invalidPodNameRegex:
		documentation = j.addInvalidPodNameErrorDocumentation(matches)
	case versionConflictRegex:
		documentation = j.addVersionConflictErrorDocumentation(matches)
	case xcodeVersionRegex:
		documentation = j.addXcodeVersionErrorDocumentation()
	}

	cmdErr.SetDocumentation(documentation)

	return cmdErr
}

func (j *Job) addSpecsRepoNotFoundErrorDocumentation() string {
	return strings.Join(
		[]string{
			"Failed to find pod specification.",
			"Try running 'pod repo update' to update your local specs repository,",
			"or check that the pod name and version are correct in your Podfile.",
		}, " ")
}

func (j *Job) addNetworkUnreachableErrorDocumentation() string {
	return strings.Join(
		[]string{
			"We weren't able to retrieve one or more dependencies.",
			"Please check your Internet connection and try again.",
		}, " ")
}

func (j *Job) addInvalidPodNameErrorDocumentation(matches [][]string) string {
	message := ""
	if len(matches) > 0 && len(matches[0]) > 1 {
		message = matches[0][1]
	}

	return strings.Join(
		[]string{
			"Couldn't find pod",
			message,
			", please make sure it is spelt correctly and exists in the CocoaPods repository:\n",
		}, " ")
}

func (j *Job) addVersionConflictErrorDocumentation(matches [][]string) string {
	message := ""
	if len(matches) > 0 && len(matches[0]) > 1 {
		message = matches[0][1]
	}

	return strings.Join(
		[]string{
			"Couldn't resolve version conflict for",
			message,
			", please check your Podfile for conflicting version requirements:\n",
		}, " ")
}

func (j *Job) addXcodeVersionErrorDocumentation() string {
	return strings.Join(
		[]string{
			"One or more pods require a higher minimum deployment target.",
			"Please update the platform version in your Podfile or update the affected pods.",
		}, " ")
}
