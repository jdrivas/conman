package conman

//
//  Viper config file constants
//

// This is intended to support multiple configurations
// read in through a viper config file.
// Sttructured as (e.g. using yaml)

// defaultConnection: connection-name-1
// connections:
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
// If the config paramater defaultConnection is set, then this name is used as a default,
// if there is connection with that name deflined.
// If that is not defined, then the list of connetions is sorted lexographically and the first
// connection is used (I would rather have it be the first one in the connection list, but viper
// manages nested configurations as maps and they are randomly ordered).
// If not connections are defined then there is a default connection named DefaultConnectionNameValue
// and with ServiceURL set by DefaultServiceURL.

const (
	ConnectionsKey           = "connections"       // string
	DefaultConnectionNameKey = "defaultConnection" // string
	ServiceURLKey            = "serviceURL"        // string
	AuthTokenKey             = "authToken"         //string
	HeadersKey               = "headers"           // map[string]string
)

// ConnectionFlagKey          = "connection"        //string
