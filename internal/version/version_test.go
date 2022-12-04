package version

import (
	"encoding/json"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetVersion(t *testing.T) {
	BuildVersion = "0.0.0-test"
	BuildDate = "now"
	GitCommit = "abc"

	out, err := GetVersion()
	require.NoError(t, err)

	actual := providerVersion{}
	err = json.Unmarshal([]byte(out), &actual)
	require.NoError(t, err)

	expected := providerVersion{
		BuildVersion: BuildVersion,
		GitCommit:    GitCommit,
		BuildDate:    BuildDate,
		GoVersion:    runtime.Version(),
		OsArch:       runtime.GOOS + "/" + runtime.GOARCH,
	}

	require.EqualValues(t, expected, actual)
}
