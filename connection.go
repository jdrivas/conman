package conman

import (
	"fmt"
	"sort"

	t "github.com/jdrivas/termtext"
	"github.com/jdrivas/vconfig"
	"github.com/spf13/viper"
)

//
// Public API
//

// Connection contains information for connecting to a service endpoint.
type Connection struct {
	Name       string
	ServiceURL string
	AuthToken  string
	Headers    map[string]string
}

// ConnectionList for handling our set of connections.
type ConnectionList []*Connection

type byName ConnectionList

func (a byName) Len() int           { return len(a) }
func (a byName) Swap(i, j int)      { a[j], a[i] = a[i], a[j] }
func (a byName) Less(i, j int) bool { return a[i].Name < a[j].Name }

// GetCurrentConnection is the primary interface for obtaining a connection.
func GetCurrentConnection() (c *Connection, err error) {
	if viper.IsSet(DefaultConnectionNameKey) {
		cn := viper.GetString(DefaultConnectionNameKey)
		var ok bool
		if c, ok = GetConnection(cn); !ok {
			err = fmt.Errorf("couldn't find connection: \"%s\"", cn)
		}
	} else {
		err = fmt.Errorf("defualt connection not set")
	}
	return c, err
}

// GetConnection by name (from configuration).
func GetConnection(name string) (*Connection, bool) {
	return getConnectionFromConfig(name)
}

// SetConnection sets a new default.
func SetConnection(name string) (ok bool) {
	var conn *Connection
	if conn, ok = GetConnection(name); ok {
		viper.Set(DefaultConnectionNameKey, conn.Name)
	}
	return ok
}

// GetAllConnections returns a list of known connections
func GetAllConnections() ConnectionList {
	conns := getAllConnectionsFromConfig()
	sort.Sort(byName(conns))
	return conns
}

// FindConnection returns a connection if it's in the list
// otherwise nil.
func (cl ConnectionList) FindConnection(name string) (conn *Connection) {
	for _, c := range cl {
		if c.Name == name {
			conn = c
			break
		}
	}
	return conn
}

// Private API
// Read in the config to get all the named connections

func getAllConnectionsFromConfig() (cl ConnectionList) {
	cm := viper.GetStringMap(ConnectionsKey) // map[string]interface{}
	for name := range cm {
		if c, ok := getConnectionFromConfig(name); ok {
			cl = append(cl, c)
		} else {
			panic(fmt.Sprintf("Couldn't find connection name \"%s\" in configuration.", name))
		}
	}
	return cl
}

func getConnectionFromConfig(name string) (c *Connection, ok bool) {
	ck := fmt.Sprintf("%s.%s", ConnectionsKey, name)
	if viper.IsSet(ck) {
		c = &Connection{
			Name:       name,
			ServiceURL: viper.GetString(fmt.Sprintf("%s.%s", ck, ServiceURLKey)),
			AuthToken:  viper.GetString(fmt.Sprintf("%s.%s", ck, AuthTokenKey)),
			Headers:    viper.GetStringMapString(fmt.Sprintf("%s.%s", ck, HeadersKey)),
		}
		ok = true
	}
	return c, ok
}

// ConnectionFlagValue this is where command line flag must store a conenction value to use.
var ConnectionFlagValue string
var previouslySetByFlag bool

// InitConnections initializes a default connection.
// Needs to happen after we've read in the viper configuration file.
const defaultServiceURL = "http://127.0.0.1:80"
const defaultConnectionName = "broken-default"

var defaultConn = &Connection{
	Name:       defaultConnectionName,
	ServiceURL: defaultServiceURL,
}

func InitConnections() {
	if vconfig.Debug() {
		t.Pef()
		defer t.Pxf()
	}

	var conn *Connection
	var err error
	var ok bool

	// Get the current connection (this looks up the viper variable for the currrent default connection.)
	if conn, err = GetCurrentConnection(); err != nil {
		// .. Otherwise, see if there is a _name_ of a defined connection to use as default ...
		defaultName := viper.GetString(DefaultConnectionNameKey)
		if conn, ok = GetConnection(defaultName); !ok {

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
				if vconfig.Debug() {
					fmt.Printf("Using a 'broken' default connection.\n")
				}
				conn = defaultConn
				// Add this to the configuraiton so we find it in a any latter GetCurrentConnection.
				viper.Set(fmt.Sprintf("%s.%s.%s",
					ConnectionsKey, conn.Name, ServiceURLKey), conn.ServiceURL)
			}
			viper.Set(DefaultConnectionNameKey, conn.Name)
		} else {
			viper.Set(DefaultConnectionNameKey, conn.Name)
		}
	}
	if vconfig.Debug() {
		fmt.Printf("Using connection: %s[%s]\n", conn.Name, conn.ServiceURL)
	}
}
