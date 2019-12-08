package connection

import (
	"fmt"
	"testing"

	"github.com/jdrivas/gafw/config"
	"github.com/spf13/viper"
)

func Test_Stack_PushPop(t *testing.T) {

	c := &Connection{
		Name: "test-1",
	}

	PushCurrentConnection(c)

	expectedSize := 1
	if len(currentConnections) != expectedSize {
		t.Errorf("Stack is wrong size. Should be %d, got %d", expectedSize, len(currentConnections))
		fmt.Printf("Stack: %v\n", currentConnections)
	}

	nc := PopCurrentConnection()
	if len(currentConnections) != 0 {
		t.Errorf("Stack is wrong size. Should be 0, got %d", len(currentConnections))
		fmt.Printf("Stack: %v\n", currentConnections)
	}
	if c != nc {
		t.Errorf("Popped the wrong elmeent. \nt\tShould be %v\n\tgot %v", c, nc)
		fmt.Printf("Stack: %v\n", currentConnections)
	}
}

func Test_Stack_Bounds(t *testing.T) {

	nameBase := "Test-"
	var testSize = 10
	for i := 0; i < testSize; i++ {
		c := &Connection{
			Name: fmt.Sprintf("%s%d", nameBase, i),
		}
		PushCurrentConnection(c)
	}

	if len(currentConnections) != testSize {
		t.Errorf("Pushed %d connections, but length was %d.", testSize, len(currentConnections))
		fmt.Printf("Stack: %v\n", currentConnections)
	}

	for i := testSize - 1; i >= 0; i-- {
		expectName := fmt.Sprintf("%s%d", nameBase, i)
		c := PopCurrentConnection()
		if c.Name != expectName {
			t.Errorf("Checking names, got: %s, expected %s", c.Name, expectName)
			fmt.Printf("Stack: %#+v\n", currentConnections)
		}
	}

}

// This sets up viper with Connection configuration.
func setConnectionConfig(name, url string) {
	viper.Set(fmt.Sprintf("%s.%s.%s", config.ConnectionsKey, name, config.ServiceURLKey), url)
}

func Test_Stack_Growth(t *testing.T) {
	viper.Reset()

	type sTest struct {
		name string
		url  string
	}

	sTests := []sTest{
		sTest{name: "Test-1", url: "htpps://localhost"},
		sTest{name: "Test-2", url: "htpps://forest-service.qcs.rigetti.com"},
		sTest{name: "Test-3", url: "htpps://127.0.0.1"},
	}

	for _, st := range sTests {
		setConnectionConfig(st.name, st.url)
	}

	// Run throuh  the connections ...
	for _, st := range sTests {

		// .. set them ..
		SetConnection(st.name)

		// ... don't grow the stack ...
		expectedLen := 1
		if len(currentConnections) != expectedLen {
			t.Errorf("Length of stack is: %d, expect it to be: %d.", len(currentConnections), expectedLen)
			fmt.Printf("Stack: %#+v\n", currentConnections)
		}

		//  ... make sure that the right one is set.
		gcName := GetCurrentConnection().Name
		if gcName != st.name {
			t.Errorf("Pushed connection %s, GetConnection returnd %s.", st.name, gcName)
			fmt.Printf("Stack: %#+v\n", currentConnections)
		}
	}

	viper.Reset()
}
