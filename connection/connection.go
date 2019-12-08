package connection

import (
	"fmt"
	"os"
	"sort"

	"github.com/jdrivas/gafw/config"
	"github.com/jdrivas/gafw/term"
	t "github.com/jdrivas/gafw/term"
	"github.com/juju/ansiterm"
	"github.com/spf13/viper"
)

//
// Public API
//

//
//  Viper config file constants
//

// This is intended to support multiple configurations
// read in through a viper config file.
// Sttructured as (e.g. using yaml)

// defaultConnection: connection-name-1
// connections:
//       default:
//             serviceURL: http://127.0.0.1
//						 authToken: XXX-YYY-ZZZ
//             heaeders:
//                   X-APP-PARAM:  some-param
//       connection-name-1:
//             serviceURL: http://localhost
//						 authToken: XXX-YYY-ZZZ
//             heaeders:
//                   X-APP-PARAM:  some-param
//       connection-name-2:
//             serviceURL: http://localhost
//						 authToken: XXX-YYY-ZZZ
//             heaeders:
//                   X-APP-PARAM:  some-param
//
// DefaultConnection
//  If the config paramater defaultConnection is set, then this name is used as a default,
// if there is connection with that name deflined.
// If defaultConnection is not set, or not found, then the connection named DefaultConnectionNameKey
//  is used.
// If that is not defined, then the list of connetions is sorted lexographically and the first
// connection is used (I would rather have it be the first one in the connection list, but viper
// manages nested configurations as maps and they are randomly ordered).
// If not connections are defined then there is a default connection named DefaultConnectionNameValue
// and with ServiceURL set by DefaultServiceURL.

// initConnections sets up the first current Connection,
// initializes the ShowTokens state, and should be called whenever the Viper config file gets reloaded.
// Since we need at least a URL to break and/or let us know that no token has been set.

const (
	ConnectionsKey             = "connections"       // string
	DefaultConnectionNameKey   = "defaultConnection" // string
	DefaultConnectionNameValue = "default"           // value is a string
	ConnectionFlagKey          = "connection"        //string
	ServiceURLKey              = "serviceURL"        // string
	AuthTokenKey               = "authToken"         //string
	HeadersKey                 = "headers"           // map[string]string
)

const DefaultServiceURL = "http://127.0.0.1:80"

// Connection contains information for connecting to a service endpoint.
type Connection struct {
	Name       string
	ServiceURL string
	AuthToken  string
	Headers    map[string]string
}

// ConnectionsList is a colleciton of connections.
type ConnectionList []*Connection

// Sort by name
type byName ConnectionList

func (a byName) Len() int           { return len(a) }
func (a byName) Swap(i, j int)      { a[j], a[i] = a[i], a[j] }
func (a byName) Less(i, j int) bool { return a[i].Name < a[j].Name }

// Manage a stack of connections
var currentConnections = make([]*Connection, 0)

// GetCurrentConnection is the primary interface for getting the connection to use.
func GetCurrentConnection() *Connection {
	if len(currentConnections) == 0 {
		return nil
	}

	// Get the top of the stack
	return currentConnections[len(currentConnections)-1]
}

// GetConnection by name (from configuration).
func GetConnection(name string) (*Connection, bool) {
	return getConnectionFromConfig(name)
}

// SetConnection sets the current connection by name
// That is, replace the stop of the stack with the onen given.
func SetConnection(name string) (err error) {
	conn, ok := GetConnection(name)
	if ok {
		PopCurrentConnection()
		PushCurrentConnection(conn)
	} else {
		err = fmt.Errorf("couldn't find connection \"%s\"", name)
	}
	return err
}

// GetAllConnections returns a list of known connections
func GetAllConnections() ConnectionList {
	return getAllConnectionsFromConfig()
}

//
// The current connection is implemented with a stack
// to facilitate changing the connection via a flag.

// PushCurrentConnection  adds the connection to the stack, setting the new current connection.
func PushCurrentConnection(conn *Connection) {
	// fmt.Printf("Push - Size of conn stack: %d\n", len(currentConnections))
	// Don't push a dupilicate on the stack.
	if len(currentConnections) == 0 || GetCurrentConnection().Name != conn.Name {
		currentConnections = append(currentConnections, conn)
	}
	// fmt.Printf("Push done - Size of conn stack: %d\n", len(currentConnections))
}

// PopCurrentConnection returns the top of the stack and sets the new
// current connection to the next topmost stack item.
// Returns nil on empty stack.
func PopCurrentConnection() (c *Connection) {
	// fmt.Printf("Pop - Size of conn stack: %d\n", len(currentConnections))
	if len(currentConnections) > 0 {
		currentConnections, c = currentConnections[:len(currentConnections)-1], currentConnections[len(currentConnections)-1]
	}
	// fmt.Printf("Pop Done - Size of conn stack: %d\n", len(currentConnections))
	return c
}

// List displpays the list of connections and notes the current one.
func (conns ConnectionList) List() {
	sort.Sort(byName(conns))
	if len(conns) > 0 {
		currentName := GetCurrentConnection().Name
		w := ansiterm.NewTabWriter(os.Stdout, 4, 4, 3, ' ', 0)
		fmt.Fprintf(w, "%s\n", t.Title("\tName\tURL"))
		for _, c := range conns {
			name := term.Text(c.Name)
			current := ""
			if c.Name == currentName {
				name = term.Highlight("%s", c.Name)
				current = term.Highlight("%s", "*")
			}
			fmt.Fprintf(w, "%s\t%s\t%s\n", current, name, t.Text("%s", c.ServiceURL))
		}
		w.Flush()
	} else {
		fmt.Printf("%s\n", t.Title("There were no connections."))
	}
}

// Private API
// Read in the config to get all the named connections
func getAllConnectionsFromConfig() (conns ConnectionList) {
	connectionsMap := viper.GetStringMap(ConnectionsKey) // map[string]interface{}
	for name := range connectionsMap {
		conn, ok := getConnectionFromConfig(name)
		if ok {
			conns = append(conns, conn)
		} else {
			fmt.Printf(term.Error(fmt.Errorf("couldn't create a config for connection \"%s\"", name)))
		}
	}
	return conns
}

func getConnectionFromConfig(name string) (conn *Connection, ok bool) {
	connKey := fmt.Sprintf("%s.%s", ConnectionsKey, name)
	if viper.IsSet(connKey) {
		conn = &Connection{
			Name:       name,
			ServiceURL: viper.GetString(fmt.Sprintf("%s.%s", connKey, ServiceURLKey)),
			AuthToken:  viper.GetString(fmt.Sprintf("%s.%s", connKey, AuthTokenKey)),
			Headers:    viper.GetStringMapString(fmt.Sprintf("%s.%s", connKey, HeadersKey)),
		}
		ok = true
	}
	return conn, ok
}

// ConnectionFlagValue his is where command line  flag will store a conenction value to use.
var ConnectionFlagValue string
var previouslySetByFlag bool

// InitConnections initializes a default connection. Needs to happen after we've read in the viper configuration file.
// TODO: It's probably best if init is idempotent.
func InitConnections() {
	if config.Debug() {
		fmt.Printf("Initializing Connections\n")
	}

	// Flag will overide all so:
	var ok bool
	var conn *Connection
	if ConnectionFlagValue != "" {
		if config.Debug() {
			fmt.Printf("Using flag value.\n")
		}
		conn, ok = GetConnection(ConnectionFlagValue)
		if !ok {
			fmt.Printf("Couldn't find the connection: \"%s\"\n", ConnectionFlagValue)
			// Yes, now we will have no connection set if there wsas not one already set.
		} else {
			previouslySetByFlag = true
		}
	} else {

		// Current conenction should be durable during interactive mode
		// reset it to the default ...
		if len(currentConnections) == 0 {
			if config.Debug() {
				fmt.Printf("No current Connection.\n")
			}
			// If there is a connection named default, use it ....
			conn, ok = GetConnection(DefaultConnectionNameValue)
			if !ok {
				// .. Otherwise, see if there is a _name_ of a defined connection to use as default ...
				defaultName := viper.GetString(DefaultConnectionNameKey)
				conn, ok = GetConnection(defaultName)
				if !ok {

					// ... next look for _any_ defined connections.
					// Rather than pick a random connection (maps don't have a determined order.
					// and we get connections from the config file as a map), pick the first lexographic one.
					conns := getAllConnectionsFromConfig()
					if len(conns) > 0 {
						sort.Sort(byName(conns))
						conn = conns[0]
					} else {
						// ... As a last resort set up a broken empty connection.
						// We won't panic here as we can set it during interactive
						// mode and it will otherwise error.
						if config.Debug() {
							fmt.Printf("Using a 'broken' default connection.\n")
						}
						conn = &Connection{
							Name:       DefaultConnectionNameValue,
							ServiceURL: DefaultServiceURL,
						}
					}
				}
			}
		} else { // use the current connection whatever it is.
			conn = GetCurrentConnection()
		}
	}
	if config.Debug() {
		fmt.Printf("Using connection: %s[%s]\n", conn.Name, conn.ServiceURL)
	}
	PushCurrentConnection(conn)
}

// ResetConnection is called to unset the conneciton set by a flag.
// Most usefull in an interactive mode where you want the flag
// to be a one time effect.
func ResetConnection() {
	// If the last time through, we were set by a flag
	// get the old connetion back and decide what to do.
	if previouslySetByFlag {
		if config.Debug() {
			fmt.Printf("Reseting connection to pre-flag.\n")
		}

		// Don't reset to an empty connection.
		// This can happen if we provide a connection
		// by flag on the command line but the command line
		// sends us into interactive mode.
		previouslySetByFlag = false
		if len(currentConnections) > 1 {
			PopCurrentConnection()
		}
	}

}
