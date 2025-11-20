package cocoapods

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeInstallCmd(t *testing.T) {
	podCommand := "pod"
	cmd, err := CmdFactory{
		execPath: ExecPath{},
	}.MakeInstallCmd(podCommand, "file")
	assert.NoError(t, err)
	assert.NotNil(t, cmd)
	args := cmd.Args
	assert.Contains(t, args, "pod")
	assert.Contains(t, args, "install")
	assert.Contains(t, args, "--no-repo-update")
}
