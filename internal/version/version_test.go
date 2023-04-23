package version

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetVersion(t *testing.T) {
	BuildVersion = "0.0.0-test"
	BuildDate = "now"
	GitCommit = "abc"

	actual := GetVersion()
	expected := ProviderVersion{
		BuildVersion: BuildVersion,
		GitCommit:    GitCommit,
		BuildDate:    BuildDate,
		GoVersion:    runtime.Version(),
		OsArch:       runtime.GOOS + "/" + runtime.GOARCH,
	}

	require.EqualValues(t, expected, actual)
}
