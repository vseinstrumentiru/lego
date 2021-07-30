package project

import (
	"embed"
	"fmt"
	"net/url"
	"os"

	"emperror.dev/errors"
	"github.com/iancoleman/strcase"

	"github.com/vseinstrumentiru/lego/v2/gen/internal/pkg"
)

//nolint:gochecknoglobals
//go:embed stubs/*
var stubs embed.FS

func create(c *config) ([]string, error) {
	var err error
	if c.Name, c.workDir, err = pkg.GetPath(c.Name); err != nil {
		return nil, err
	}
	c.CamelCaseName = strcase.ToCamel(c.Name)

	if c.GitRemotePath != "" {
		u, err := url.Parse(c.GitRemotePath)
		if err != nil {
			return nil, errors.New("wrong git remote url")
		}

		c.ModuleName = fmt.Sprintf("%s%s", u.Host, u.Path)
		c.HasGit = true
	} else {
		c.ModuleName = c.Name
	}

	if _, err = os.Stat(c.workDir); os.IsNotExist(err) {
		return nil, err
	}

	path := c.workDir + "/" + c.Name
	cmd := pkg.NewStubExecutor(stubs, path, c.Verbose, c.DryRun)
	// create project directory
	if c.dirExists, err = cmd.Mkdir(fmt.Sprintf("../%s", c.Name)); err != nil {
		return nil, err
	}

	// create root files
	if err = cmd.CopyStub("stubs/gitignore.tmpl", "./.gitignore", c); err != nil {
		return nil, err
	}
	if err = cmd.CopyStub("stubs/README.md.tmpl", "./README.md", c); err != nil {
		return nil, err
	}

	if err = cmd.CopyStub("stubs/go.mod.tmpl", "./go.mod", c); err != nil {
		return nil, err
	}
	if err = cmd.CopyStub("stubs/golangci.yml", "./golangci.yml"); err != nil {
		return nil, err
	}

	if err = cmd.CopyStub("stubs/Dockerfile", "./Dockerfile"); err != nil {
		return nil, err
	}
	if err = cmd.CopyStub("stubs/dockerignore.tmpl", "./.dockerignore", c); err != nil {
		return nil, err
	}

	if err = cmd.CopyStub("stubs/doc.go.tmpl", "./doc.go", c); err != nil {
		return nil, err
	}

	if err = cmd.CopyStub("stubs/Makefile.tmpl", "./Makefile", c); err != nil {
		return nil, err
	}
	if _, err = cmd.Mkdir("./.make"); err != nil {
		return nil, err
	}
	if err = cmd.CopyStub("stubs/.make/build.mk", "./.make/build.mk"); err != nil {
		return nil, err
	}
	if err = cmd.CopyStub("stubs/.make/mga.mk", "./.make/mga.mk"); err != nil {
		return nil, err
	}
	if err = cmd.CopyStub("stubs/.make/graphql.mk", "./.make/graphql.mk"); err != nil {
		return nil, err
	}
	if err = cmd.CopyStub("stubs/.make/protobuf.mk", "./.make/protobuf.mk"); err != nil {
		return nil, err
	}

	if c.UseProtobuf {
		if err = cmd.CopyStub("stubs/buf.yaml", "./buf.yaml"); err != nil {
			return nil, err
		}

		apiPath := fmt.Sprintf("./api/%s/v1", c.Name)

		if _, err = cmd.Mkdir(apiPath); err != nil {
			return nil, err
		}

		if err = cmd.CopyStub("stubs/api/protobuf/main.proto.tmpl", fmt.Sprintf("%s/%s.proto", apiPath, c.Name), c); err != nil {
			return nil, err
		}
	}

	if c.UseGraphql {
		if err = cmd.CopyStub("stubs/gqlgen.yml.tmpl", "./gqlgen.yml", c); err != nil {
			return nil, err
		}

		apiPath := fmt.Sprintf("./api/%s/v1", c.Name)

		if _, err = cmd.Mkdir(apiPath); err != nil {
			return nil, err
		}

		if err = cmd.CopyStub("stubs/api/graphql/main.graphql", fmt.Sprintf("%s/main.graphql", apiPath), c); err != nil {
			return nil, err
		}
	}

	if _, err = cmd.Mkdir("./cmd/server"); err != nil {
		return nil, err
	}
	if err = cmd.CopyStub("stubs/cmd/server/app.go", "./cmd/server/app.go"); err != nil {
		return nil, err
	}
	if err = cmd.CopyStub("stubs/cmd/server/config.go", "./cmd/server/config.go"); err != nil {
		return nil, err
	}
	if err = cmd.CopyStub("stubs/cmd/server/main.go", "./cmd/server/main.go"); err != nil {
		return nil, err
	}
	if err = cmd.CopyStub("stubs/cmd/server/config.yaml.tmpl", "./cmd/server/config.yaml"); err != nil {
		return nil, err
	}

	if _, err = cmd.Mkdir("./internal/app"); err != nil {
		return nil, err
	}

	if _, err = cmd.Mkdir("./internal/pkg"); err != nil {
		return nil, err
	}

	if _, err = cmd.Mkdir("./pkg"); err != nil {
		return nil, err
	}

	// init git
	// add git remote

	return cmd.Result(), nil
}
