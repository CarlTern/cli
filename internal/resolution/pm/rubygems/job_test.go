package rubygems

import (
	"errors"
	"testing"

	jobTestdata "github.com/debricked/cli/internal/resolution/job/testdata"
	"github.com/debricked/cli/internal/resolution/pm/rubygems/testdata"
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
			name:  "Bundler not found",
			error: "        |exec: \"bundle\": executable file not found in $PATH",
			doc:   "Bundler wasn't found. Please check if it is installed and accessible by the CLI.",
		},
		{
			name:  "Gem not found",
			error: "Could not find gem 'nonexistent-gem (>= 0)' in rubygems repository https://rubygems.org/ or installed locally.",
			doc:   "Couldn't find gem nonexistent-gem (>= 0) , please make sure it is spelt correctly and exists in the RubyGems repository:\n",
		},
		{
			name:  "Version conflict",
			error: "Bundler could not find compatible versions for gem \"rails\":\n  In Gemfile:\n    rails (~> 6.0)\n\n    rails (= 5.2.3)",
			doc:   "Couldn't resolve version conflict for gem rails , please check your Gemfile for conflicting version requirements:\n",
		},
		{
			name:  "Network error",
			error: "Could not reach rubygems repository https://rubygems.org/\nRetrying fetcher due to error (2/4): Bundler::HTTPError Network is unreachable",
			doc:   "We weren't able to retrieve one or more dependencies. Please check your Internet connection and try again.",
		},
		{
			name:  "Ruby version requirement",
			error: "Your Ruby version is 2.6.0, but your Gemfile specified ~> 2.7.0\nruby_dep-1.5.0 requires Ruby version >= 2.2.5, ~> 2.2",
			doc:   "This project requires Ruby version >= 2.2.5, ~> 2.2 or newer. Please update your Ruby installation.",
		},
		{
			name:  "Gemfile parse error",
			error: "There was an error parsing `Gemfile`: syntax error, unexpected end-of-input, expecting keyword_end. Bundler cannot continue.",
			doc:   "Failed to parse Gemfile. Please check that your Gemfile has valid Ruby syntax.",
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
