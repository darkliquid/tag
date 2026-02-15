package commands

import (
	"path/filepath"

	"github.com/adrg/xdg"
)

// Common interfaces and types for all commands
var DBDir = filepath.Join(xdg.DataHome, "tag/index.db")
