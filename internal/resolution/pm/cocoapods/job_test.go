package cocoapods

import (
	"errors"
	"testing"

	jobTestdata "github.com/debricked/cli/internal/resolution/job/testdata"
	"github.com/debricked/cli/internal/resolution/pm/cocoapods/testdata"
	"github.com/debricked/cli/internal/resolution/pm/util"
	"github.com/stretchr/testify/assert"
)

const (
	badName = "bad-name"
)

func TestNewJob(t *testing.T) {
	j := NewJob("file", false, CmdFactory{
		execPath: ExecPath{},
	})
	assert.Equal(t, "file", j.GetFile())
	assert.False(t, j.Errors().HasError())
}

func TestRunInstall(t *testing.T) {
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	j := NewJob("file", false, cmdFactoryMock)

	_, err := j.runInstallCmd()
	assert.NoError(t, err)

	assert.False(t, j.Errors().HasError())
}

func TestInstall(t *testing.T) {
	j := Job{install: true}
	assert.Equal(t, true, j.Install())

	j = Job{install: false}
	assert.Equal(t, false, j.Install())
}

func TestRunInstallCmdErr(t *testing.T) {
	cases := []struct {
		name  string
		error string
		doc   string
	}{
		{
			name:  "General error",
			error: "cmd-error",
			doc:   util.UnknownError,
		},
		{
			name:  "Pod not found",
			error: "        |exec: \"pod\": executable file not found in $PATH",
			doc:   "CocoaPods wasn't found. Please check if it is installed and accessible by the CLI.",
		},
		{
			name:  "Specs repo not found",
			error: "[!] Unable to find a specification for `NonExistentPod`",
			doc:   "Failed to find pod specification. Try running 'pod repo update' to update your local specs repository, or check that the pod name and version are correct in your Podfile.",
		},
		{
			name:  "Network error",
			error: "Failed to connect to raw.githubusercontent.com port 443: Connection refused",
			doc:   "We weren't able to retrieve one or more dependencies. Please check your Internet connection and try again.",
		},
		{
			name:  "Invalid pod name",
			error: "[!] Unable to find a pod with name `InvalidPodName`",
			doc:   "Couldn't find pod InvalidPodName , please make sure it is spelt correctly and exists in the CocoaPods repository:\n",
		},
		{
			name:  "Version conflict",
			error: "[!] Unable to satisfy the following requirements: `Alamofire (~> 5.0)` required by `Podfile`",
			doc:   "Couldn't resolve version conflict for Alamofire (~> 5.0) , please check your Podfile for conflicting version requirements:\n",
		},
		{
			name:  "Xcode deployment target",
			error: "The iOS deployment target 'IPHONEOS_DEPLOYMENT_TARGET' is set to 8.0, but the range of supported deployment target versions is 9.0 to 16.2.99. requires a higher minimum deployment target",
			doc:   "One or more pods require a higher minimum deployment target. Please update the platform version in your Podfile or update the affected pods.",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			expectedError := util.NewPMJobError(c.error)
			expectedError.SetDocumentation(c.doc)

			cmdErr := errors.New(c.error)
			j := NewJob("file", true, testdata.CmdFactoryMock{InstallCmdName: "echo", MakeInstallErr: cmdErr})

			go jobTestdata.WaitStatus(j)
			j.Run()

			allErrors := j.Errors().GetAll()

			assert.Len(t, allErrors, 1)
			assert.Contains(t, allErrors, expectedError)
		})
	}
}

func TestRunInstallCmdOutputErr(t *testing.T) {
	cmdMock := testdata.NewEchoCmdFactory()
	cmdMock.InstallCmdName = badName
	j := NewJob("file", true, cmdMock)

	go jobTestdata.WaitStatus(j)
	j.Run()

	jobTestdata.AssertPathErr(t, j.Errors())
}
