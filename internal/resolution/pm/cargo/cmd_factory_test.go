package cargo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeInstallCmd(t *testing.T) {
	cargoCommand := "cargo"
	cmd, err := CmdFactory{
		execPath: ExecPath{},
	}.MakeInstallCmd(cargoCommand, "file")
	assert.NoError(t, err)
	assert.NotNil(t, cmd)
	args := cmd.Args
	assert.Contains(t, args, "cargo")
	assert.Contains(t, args, "generate-lockfile")
}
