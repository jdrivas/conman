package conman

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

// This sets up viper with Connection configuration.
func setConnectionConfig(name, url string) {
	viper.Set(fmt.Sprintf("%s.%s.%s", ConnectionsKey, name, ServiceURLKey), url)
}

// Use this at the end of a test that either Inits, or sets configuration some how.
func resetConfig() {
	viper.Reset()
}

// print function and file location at depth. Depth = 0 is this level.
// depth > 1 prints the caller depth levels up (e.g. 1 is the caller of
// the function that pl is called from).
func pls(depth int) string {
	fc, fn, l := loc(depth + 1)
	return fmt.Sprintf("%s: %s line %d", fc, fn, l)
}

func loc(d int) (fnc, file string, line int) {
	if pc, fl, l, ok := runtime.Caller(d + 1); ok {
		f := runtime.FuncForPC(pc)
		fnc = filepath.Base(f.Name())
		file = filepath.Base(fl)
		line = l
	}
	return fnc, file, line
}
