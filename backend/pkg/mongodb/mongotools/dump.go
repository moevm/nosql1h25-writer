package mongotools

import (
	"fmt"
	"os/exec"
)

type MongoDumper interface {
	Restore(filePath string) error
	Dump(filePath string) error
}

type dumper struct {
	uri string
}

func NewDumper(uri string) MongoDumper {
	return &dumper{uri: uri}
}

func (d *dumper) Restore(filepath string) error {
	if err := exec.Command("mongorestore",
		"--uri="+d.uri,
		"--archive="+filepath,
		"--gzip",
		"--drop",
	).Run(); err != nil {
		return fmt.Errorf("mongorestore: %w", err)
	}

	return nil
}

func (d *dumper) Dump(filepath string) error {
	if err := exec.Command("mongodump",
		"--uri="+d.uri,
		"--archive="+filepath,
		"--gzip",
	).Run(); err != nil {
		return fmt.Errorf("mongodump: %w", err)
	}

	return nil
}
