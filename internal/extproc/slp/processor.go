package slp

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/kaz/pprotein/internal/collect"
)

type (
	processor struct {
		confPath string
	}
)

const (
	dbTypeMySQL      = "my"
	dbTypePostgreSQL = "pg"
)

func (p *processor) Cacheable() bool {
	return true
}

func (p *processor) Process(snapshot *collect.Snapshot) (io.ReadCloser, error) {
	bodyPath, err := snapshot.BodyPath()
	if err != nil {
		return nil, fmt.Errorf("failed to find snapshot body: %w", err)
	}

	var dbType string
	dbTypeEnv, ok := os.LookupEnv("DB_TYPE")
	if !ok {
		dbType = dbTypeMySQL
	} else {
		if strings.Contains(dbTypeEnv, "p") {
			dbType = dbTypePostgreSQL
		} else {
			dbType = dbTypeMySQL
		}
	}

	cmd := exec.Command("slp", dbType, "--config", p.confPath, "--output", "standard", "--format", "tsv", "--file", bodyPath)

	res, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("external process aborted: %w", err)
	}

	return io.NopCloser(bytes.NewBuffer(res)), nil
}
