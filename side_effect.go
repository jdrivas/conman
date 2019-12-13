package conman

import "time"

// Meta information about a request
// TODO: Consider a better name.
type SideEffect struct {
	ElapsedTime time.Duration
}
