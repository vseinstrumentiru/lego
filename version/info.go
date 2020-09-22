package version

import (
	"fmt"
	"os"
	"runtime"

	"github.com/vseinstrumentiru/lego/config"
)

var (
	Version    = "undefined"
	CommitHash = "undefined"
	BuildDate  = "undefined"
)

type Info struct {
	Version    string `json:"version"`
	CommitHash string `json:"commit_hash"`
	BuildDate  string `json:"build_date"`
	GoVersion  string `json:"go_version"`
	Os         string `json:"os"`
	Arch       string `json:"arch"`
	Compiler   string `json:"compiler"`
	Host       string `json:"host"`
	DataCenter string `json:"data_center"`
}

func New(config *config.Application) Info {
	host, _ := os.Hostname()
	return Info{
		Version:    Version,
		CommitHash: CommitHash,
		BuildDate:  BuildDate,
		GoVersion:  runtime.Version(),
		Os:         runtime.GOOS,
		Arch:       runtime.GOARCH,
		Compiler:   runtime.Compiler,
		Host:       host,
		DataCenter: config.DataCenter,
	}
}

func (i Info) Print() {
	fmt.Printf("version %s (%s) built on %s (%s+%s+%s)\n", i.Version, i.CommitHash, i.BuildDate, i.GoVersion, i.Compiler, i.Arch)
}

func (i Info) Fields() map[string]interface{} {
	return map[string]interface{}{
		"version":     i.Version,
		"commit_hash": i.CommitHash,
		"build_date":  i.BuildDate,
		"go_version":  i.GoVersion,
		"os":          i.Os,
		"arch":        i.Arch,
		"compiler":    i.Compiler,
		"host":        i.Host,
		"data_center": i.DataCenter,
	}
}
