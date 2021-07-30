package project

type config struct {
	Name          string
	CamelCaseName string
	ModuleName    string
	Verbose       bool
	DryRun        bool

	HasGit        bool
	GitRemotePath string

	UseProtobuf bool
	UseGraphql  bool
	UseOpenAPI  bool

	workDir   string
	dirExists bool
}
