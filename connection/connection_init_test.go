package connection

import (
	"fmt"
	"testing"

	"github.com/jdrivas/gafw/config"
)

// Test to make sure that when InitConnections is called
// the correct connection is set.
func Test_InitialStack(t *testing.T) {
	InitConnections()
	setConnectionConfig("test-1", "htpps://localhost")
	c := GetCurrentConnection()
	if c.Name != config.DefaultConnectionNameValue {
		t.Errorf("Checking names, got: %s, expected %s", c.Name, config.DefaultConnectionNameValue)
		fmt.Printf("Stack: %#+v\n", currentConnections)
	}

	resetConnections()
}

func Test_FlagConig(t *testing.T) {
	config.SetDebug(true)
	setName := "Test-1"
	setConnectionConfig(config.DefaultConnectionNameValue, "http://localhost")
	setConnectionConfig(setName, "http:127.0.0.1")
	ConnectionFlagValue = setName

	// Initalize the connections
	InitConnections()
	cn := GetCurrentConnection().Name
	if setName != cn {
		t.Errorf("Checking names, got: %s, expected %s", cn, setName)
	}

	ConnectionFlagValue = "" // rest as in cobra.Reset()
	ResetConnection()
	InitConnections()

	cn = GetCurrentConnection().Name
	if cn != config.DefaultConnectionNameValue {
		t.Errorf("Checking names, got: %s, expected %s", cn, config.DefaultConnectionNameValue)
	}

	resetConnections()
}
