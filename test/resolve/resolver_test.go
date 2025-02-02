package resolve

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/debricked/cli/internal/cmd/resolve"
	"github.com/debricked/cli/internal/wire"
	"github.com/stretchr/testify/assert"
)

func TestResolves(t *testing.T) {
	cases := []struct {
		name           string
		manifestFile   string
		lockFileName   string
		packageManager string
	}{
		{
			name:           "basic package.json",
			manifestFile:   "testdata/npm/package.json",
			lockFileName:   "yarn.lock",
			packageManager: "npm",
		},
		{
			name:           "basic requirements.txt",
			manifestFile:   "testdata/pip/requirements.txt",
			lockFileName:   "requirements.txt.pip.debricked.lock",
			packageManager: "pip",
		},
		{
			name:           "basic .csproj",
			manifestFile:   "testdata/nuget/csproj/basic.csproj",
			lockFileName:   "packages.lock.json",
			packageManager: "nuget",
		},
		{
			name:           "basic packages.config",
			manifestFile:   "testdata/nuget/packagesconfig/packages.config",
			lockFileName:   "packages.config.nuget.debricked.lock",
			packageManager: "nuget",
		},
		{
			name:           "basic go.mod",
			manifestFile:   "testdata/gomod/go.mod",
			lockFileName:   "gomod.debricked.lock",
			packageManager: "gomod",
		},
		{
			name:           "basic pom.xml",
			manifestFile:   "testdata/maven/pom.xml",
			lockFileName:   "maven.debricked.lock",
			packageManager: "maven",
		},
		{
			name:           "basic build.gradle",
			manifestFile:   "testdata/gradle/build.gradle",
			lockFileName:   "gradle.debricked.lock",
			packageManager: "gradle",
		},
	}

	for _, cT := range cases {
		c := cT
		t.Run(c.name, func(t *testing.T) {
			resolveCmd := resolve.NewResolveCmd(wire.GetCliContainer().Resolver())
			lockFileDir := filepath.Dir(c.manifestFile)
			lockFile := filepath.Join(lockFileDir, c.lockFileName)
			// Remove the lock file if it exists
			os.Remove(lockFile)

			err := resolveCmd.RunE(resolveCmd, []string{c.manifestFile})
			assert.NoError(t, err)

			lockFileContents, fileErr := os.ReadFile(lockFile)
			assert.NoError(t, fileErr)

			actualString := string(lockFileContents)

			assert.Greater(t, len(actualString), 0)

		})
	}
}
