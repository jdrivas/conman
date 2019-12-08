package conman

import (
	"fmt"
	"testing"
)

// Test to make sure that when InitConnections is called
// the correct connection is set.
func Test_InitialStack(t *testing.T) {
	InitConnections() // Set up.
	setConnectionConfig("test-1", "htpps://localhost")
	c := GetCurrentConnection()
	if c.Name != DefaultConnectionNameValue {
		t.Errorf("Checking names, got: %s, expected %s", c.Name, DefaultConnectionNameValue)
		fmt.Printf("Stack: %#+v\n", currentConnections)
	}

	resetConnections()
}

func Test_FlagConig(t *testing.T) {
	// vconfig.SetDebug(true)
	InitConnections() // Set up.
	// Let's get noisy
	setName := "Test-1"
	setConnectionConfig(DefaultConnectionNameValue, "http://localhost")
	setConnectionConfig(setName, "http:127.0.0.1")
	ConnectionFlagValue = setName

	// Initalize the connections
	InitConnections()
	cn := GetCurrentConnection().Name
	if setName != cn {
		t.Errorf("Checking names, got: %s, expected %s", cn, setName)
	}

	// This is the protocol for manging this in an interactive mode
	// situation (ie. a loop to read commands one after the other.)
	ConnectionFlagValue = "" // reset as in cobra.Reset()
	ResetConnection()        // Reset the connection for the single flag.
	InitConnections()        // read in default connections as needed.

	cn = GetCurrentConnection().Name
	if cn != DefaultConnectionNameValue {
		t.Errorf("Checking names, got: %s, expected %s", cn, DefaultConnectionNameValue)
	}

	resetConnections()
}
