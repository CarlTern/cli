package cargo

import (
	"errors"
	"testing"

	jobTestdata "github.com/debricked/cli/internal/resolution/job/testdata"
	"github.com/debricked/cli/internal/resolution/pm/cargo/testdata"
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
			name:  "Cargo not found",
			error: "        |exec: \"cargo\": executable file not found in $PATH",
			doc:   "Cargo wasn't found. Please check if it is installed and accessible by the CLI.",
		},
		{
			name:  "Network error - failed to download",
			error: "error: failed to download `serde v1.0.152`\n\nCaused by:\n  unable to get packages from source",
			doc:   "We weren't able to retrieve one or more dependencies. Please check your Internet connection and try again, or run 'cargo fetch' to populate the local cache before running the Debricked CLI.",
		},
		{
			name:  "Network error - unable to update",
			error: "error: Unable to update registry `crates-io`\n\nCaused by:\n  failed to fetch `https://github.com/rust-lang/crates.io-index`",
			doc:   "We weren't able to retrieve one or more dependencies. Please check your Internet connection and try again, or run 'cargo fetch' to populate the local cache before running the Debricked CLI.",
		},
		{
			name:  "Invalid crate name",
			error: "error: no matching package named `nonexistent-crate` found\nlocation searched: registry `crates-io`",
			doc:   "Couldn't find crate nonexistent-crate , please make sure it is spelt correctly and exists in the crates.io registry:\n",
		},
		{
			name:  "Version resolution error",
			error: "error: failed to select a version for the requirement `serde = \"^1.0\"`\ncandidate versions found which didn't match: 0.9.15, 0.9.14, 0.9.13, ...",
			doc:   "Couldn't resolve version requirement for serde = \"^1.0\" , please check your Cargo.toml for conflicting version requirements:\n",
		},
		{
			name:  "Incompatible Rust version",
			error: "error: package `tokio v1.25.0` cannot be built because it requires rustc 1.56 or newer, while the currently active rustc version is 1.55.0",
			doc:   "This project requires Rust version 1.56 or newer. Please update your Rust installation with 'rustup update'.",
		},
		{
			name:  "TOML parse error",
			error: "error: could not parse input as TOML\n\nCaused by:\n  TOML parse error at line 5, column 1",
			doc:   "Failed to parse Cargo.toml. Please check that your Cargo.toml file is valid TOML syntax.",
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
