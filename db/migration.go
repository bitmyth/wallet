package db

import (
	"embed"
	_ "embed"
	"github.com/pkg/errors"
	"io"
	"io/fs"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func (d DB) Migrate() error {
	err := fs.WalkDir(migrationFiles, ".", func(path string, dir fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if dir.IsDir() {
			return nil
		}
		//println(path)
		var f fs.File
		f, err = migrationFiles.Open(path)
		if err != nil {
			return errors.WithStack(err)
		}
		defer f.Close()

		content, _ := io.ReadAll(f)
		_, err = d.Exec(string(content))
		if err != nil {
			return err
		}

		return nil
	})

	return err
}
