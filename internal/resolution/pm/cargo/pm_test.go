package cargo

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPm(t *testing.T) {
	pm := NewPm()
	assert.Equal(t, Name, pm.name)
}

func TestName(t *testing.T) {
	pm := NewPm()
	assert.Equal(t, Name, pm.Name())
}

func TestManifests(t *testing.T) {
	pm := Pm{}
	manifests := pm.Manifests()
	assert.Len(t, manifests, 1)
	manifest := manifests[0]
	assert.Equal(t, `Cargo\.toml$`, manifest)
	_, err := regexp.Compile(manifest)
	assert.NoError(t, err)

	cases := map[string]bool{
		"Cargo.toml":        true,
		"Cargo.lock":        false,
		"package.json":      false,
		"MyCargo.toml":      false,
		"path/to/Cargo.toml": true,
	}
	for file, isMatch := range cases {
		t.Run(file, func(t *testing.T) {
			matched, _ := regexp.MatchString(manifest, file)
			assert.Equal(t, isMatch, matched)
		})
	}
}
