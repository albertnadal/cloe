//go:generate go run assets_generate.go

package runtime

import "github.com/albertnadal/cloe/femto"

var Files = femto.NewRuntimeFiles(files)
