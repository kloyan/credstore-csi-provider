package version

import (
	"encoding/json"
	"runtime"
)

var (
	BuildVersion string
	BuildDate    string
	GoVersion    string
	GitCommit    string
)

type providerVersion struct {
	BuildVersion string `json:"Version"`
	BuildDate    string `json:"BuildDate"`
	GitCommit    string `json:"GitCommit"`
	GoVersion    string `json:"GoVersion"`
	OsArch       string `json:"OsArch"`
}

func GetVersion() (string, error) {
	ver := providerVersion{
		BuildVersion: BuildVersion,
		BuildDate:    BuildDate,
		GitCommit:    GitCommit,
		GoVersion:    runtime.Version(),
		OsArch:       runtime.GOOS + "/" + runtime.GOARCH,
	}

	out, err := json.Marshal(ver)
	if err != nil {
		return "", err
	}

	return string(out), nil
}
