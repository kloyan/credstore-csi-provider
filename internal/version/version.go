package version

import (
	"runtime"
)

var (
	BuildVersion string
	BuildDate    string
	GoVersion    string
	GitCommit    string
)

type ProviderVersion struct {
	BuildVersion string `json:"Version"`
	BuildDate    string `json:"BuildDate"`
	GitCommit    string `json:"GitCommit"`
	GoVersion    string `json:"GoVersion"`
	OsArch       string `json:"OsArch"`
}

func GetVersion() ProviderVersion {
	return ProviderVersion{
		BuildVersion: BuildVersion,
		BuildDate:    BuildDate,
		GitCommit:    GitCommit,
		GoVersion:    runtime.Version(),
		OsArch:       runtime.GOOS + "/" + runtime.GOARCH,
	}
}
