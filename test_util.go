package conman

import (
	"github.com/spf13/viper"
)

// the rest of the tests will mock out viper and
// assum ean empty stack.
// Tests here are intending to call InitConnections to
// see if it works properly.
func resetConnections() {
	viper.Reset()
	currentConnections = make([]*Connection, 0)
}
