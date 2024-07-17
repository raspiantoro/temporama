package info

import (
	"fmt"
	"runtime"
)

var (
	//Server the name of the application (temporama)
	Server string

	//Version the number of version that are running at the moment.
	Version string

	//GitCommit compiled from which commit.
	GitCommit string

	//Environment the environment that is used by this version.
	Environment string

	//BuildDate date of the application was build.
	BuildDate = ""

	//GoVersion version of go.
	GoVersion = runtime.Version()

	//OsArch os architecture that used to run this application.
	OsArch = fmt.Sprintf("%s %s", runtime.GOOS, runtime.GOARCH)
)
