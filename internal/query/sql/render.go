package sql

import (
	_ "embed"
)

//go:embed templates/select.tmpl
var selectBuilderTemplate string
