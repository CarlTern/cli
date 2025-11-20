package rubygems

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeInstallCmd(t *testing.T) {
	bundleCommand := "bundle"
	cmd, err := CmdFactory{
		execPath: ExecPath{},
	}.MakeInstallCmd(bundleCommand, "file")
	assert.NoError(t, err)
	assert.NotNil(t, cmd)
	args := cmd.Args
	assert.Contains(t, args, "bundle")
	assert.Contains(t, args, "install")
	assert.Contains(t, args, "--quiet")
}
