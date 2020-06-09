package build

import "runtime"

var (
	Version    string
	CommitHash string
	BuildDate  string
)

type Info struct {
	Version    string `json:"version"`
	CommitHash string `json:"commit_hash"`
	BuildDate  string `json:"build_date"`
	GoVersion  string `json:"go_version"`
	Os         string `json:"os"`
	Arch       string `json:"arch"`
	Compiler   string `json:"compiler"`
}

func New() Info {
	return Info{
		Version:    Version,
		CommitHash: CommitHash,
		BuildDate:  BuildDate,
		GoVersion:  runtime.Version(),
		Os:         runtime.GOOS,
		Arch:       runtime.GOARCH,
		Compiler:   runtime.Compiler,
	}
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
	}
}
