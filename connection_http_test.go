package conman

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_smoke(t *testing.T) {

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"test": "value"}`)
	}))
	defer ts.Close()

	conn := Connection{
		ServiceURL: ts.URL,
	}

	effect, resp, err := conn.Get("/", nil)
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	if effect == nil {
		t.Errorf("Got nil side effect.")
	}

	if resp == nil {
		t.Errorf("Got nil response.")
	}
	// fmt.Printf("effect: %#+v\nresponse: %#+v\n", effect, resp)

}
