package conman

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
	"time"

	t "github.com/jdrivas/termtext"
	"github.com/jdrivas/vconfig"
)

var httpClient = http.DefaultClient

//
// Public API
//

// Send performs an HTTP request on the URL with the Token in the connection, using
// the HTTP provided by method.
// If content is non-nil, it's marshalled into the body  as a json string.
// If content is a string, it's written directly into the booy assuming its correct
// JSON: this string is not validated as correct JSON.
// If content is non-nil the Content-Type header is set to application/json.
// If result is non-nil Send umarshalls the response body,
// aasumed to be JSON encoded, into the result object passed in.
// If result is a []map[string]interface{}, you'll get a map of the JSON object.
func (conn Connection) Send(method, cmd string, content interface{}, result interface{}) (effect *SideEffect, resp *http.Response, err error) {

	if content == nil {
		req := conn.newRequest(method, cmd, nil)
		effect, resp, err = sendReq(req, result)
	} else {
		var b []byte
		switch c := content.(type) {
		case string:
			// If we marshall the string, it escapes the quotes: "foo" => \"foo\".
			// This makes for bad JSON.
			b = []byte(c)
		default:
			b, err = json.Marshal(c)
		}
		if err == nil {
			buff := bytes.NewBuffer(b)
			req := conn.newRequest(method, cmd, buff)
			req.Header.Add("Content-Type", "application/json")
			effect, resp, err = sendReq(req, result)
		}
	}
	return effect, resp, err
}

// Get works like Send with the GET verb,  but doesn't require a content object.
func (conn Connection) Get(cmd string, result interface{}) (effect *SideEffect, resp *http.Response, err error) {
	return conn.Send(http.MethodGet, cmd, nil, result)
}

// GetWithContent works like Send with GET verb. (This is provided for compatilibty with non-standard REST apis)
func (conn Connection) GetWithContent(cmd string, content, result interface{}) (effect *SideEffect, resp *http.Response, err error) {
	return conn.Send(http.MethodGet, cmd, content, result)
}

// Post works like Send using the POST verb.
func (conn Connection) Post(cmd string, content, result interface{}) (effect *SideEffect, resp *http.Response, err error) {
	return conn.Send(http.MethodPost, cmd, content, result)
}

// Delete works like Send using the Delte verb.
func (conn Connection) Delete(cmd string, content, result interface{}) (effect *SideEffect, resp *http.Response, err error) {
	return conn.Send(http.MethodDelete, cmd, content, result)
}

// Patch works like Send using the Patch verb.
func (conn Connection) Patch(cmd string, content, result interface{}) (effect *SideEffect, resp *http.Response, err error) {
	return conn.Send(http.MethodPatch, cmd, content, result)
}

//
// Private API
//

// sendReq sends along the request with some logging along the way.
func sendReq(req *http.Request, result interface{}) (effect *SideEffect, resp *http.Response, err error) {

	switch {
	// TODO: This wil dump the authorization token. Which it probably shouldn't do.
	case vconfig.Debug():
		reqDump, dumpErr := httputil.DumpRequestOut(req, true)
		reqStr := string(reqDump)
		if dumpErr != nil {
			fmt.Printf("Error dumping request (display as generic object): %v\n", dumpErr)
			reqStr = fmt.Sprintf("%v", req)
		}
		fmt.Printf("%s %s\n", t.Title("Request"), t.Text(reqStr))
		fmt.Println()
	case vconfig.Verbose():
		fmt.Printf("%s %s\n", t.Title("Request:"), t.Text("%s %s", req.Method, req.URL))
		// fmt.Println()
	}

	// Send the request
	start := time.Now()
	resp, err = httpClient.Do(req)
	effect = &SideEffect{
		ElapsedTime: time.Since(start),
	}
	if vconfig.Verbose() {
		fmt.Printf("%s %s\n", t.Title("Elapsed request time:"), t.Text("%d milliseconds", effect.ElapsedTime.Milliseconds()))
	}

	// Process
	if err == nil {

		if vconfig.Debug() {
			respDump, dumpErr := httputil.DumpResponse(resp, true)
			respStr := string(respDump)
			if dumpErr != nil {
				fmt.Printf("Error dumping response (display as generic object): %v\n", dumpErr)
				respStr = fmt.Sprintf("%v", resp)
			}
			fmt.Printf("%s\n%s\n", t.Title("Respose:"), t.Text(respStr))
			fmt.Println()
		}

		// Do this after the Dump, the dump reads out the response for reprting and
		// replaces the reader with anothe rone that has the data.
		// TODO: Figure out how to do the same replacement here so
		// th unmarshal doesn't eat the Response body.
		err = checkReturnCode(*resp)
		if result != nil {
			if err == nil {
				err = unmarshal(resp, result)
			}
		}

	}
	return effect, resp, err
}

// newRequest creates a request as usual prepending the connections ServiceURL to the cmd.
func (conn Connection) newRequest(method, cmd string, body io.Reader) *http.Request {
	// fmt.Printf("Generating request for connection: %#+v\n", conn)
	req, err := http.NewRequest(method, conn.ServiceURL+cmd, body)
	if err != nil {
		panic(fmt.Sprintf("Coulnd't generate HTTP request - %s\n", err.Error()))
	}

	for k, v := range conn.Headers {
		req.Header.Add(k, v)
	}

	return req
}

// This eats the body in the response, but returns the body in
//  obj passed in. They must match of course.
func unmarshal(resp *http.Response, obj interface{}) (err error) {
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err == nil {

		if vconfig.Debug() {
			// TODO: Printf-ing the output of json.Indent through the bytes.Buffer.String
			// produces cruft. However writting directly to it, works o.k.
			// prettyJSON := bytes.Buffer{}
			var prettyJSON bytes.Buffer
			fmt.Fprintf(&prettyJSON, t.Title("Pretty print response body:\n"))
			indentErr := json.Indent(&prettyJSON, body, "", " ")
			if indentErr == nil {
				// fmt.Printf("%s %s\n", t.Title("Response body is:"), t.Text("%s\n", prettyJSON))
				prettyJSON.WriteTo(os.Stdout)
				fmt.Println()
				fmt.Println()
			} else {
				fmt.Printf("%s\n", t.Fail("Error indenting JSON - %s", indentErr.Error()))
				fmt.Printf("%s %s\n", t.Title("Body:"), t.Text(string(body)))
			}
		}

		json.Unmarshal(body, &obj)
		if vconfig.Debug() {
			fmt.Printf("%s %s\n", t.Title("Unmarshaled object: "), t.Text("%#v", obj))
			fmt.Println()
		}
	}
	return err
}

// Returns an "informative" error if not 200
func checkReturnCode(resp http.Response) (err error) {
	err = nil
	if resp.StatusCode >= 300 {
		switch resp.StatusCode {
		case http.StatusNotFound:
			err = httpErrorMesg(resp, "Check for valid argument (user, group etc).")
		case http.StatusUnauthorized:
			err = httpErrorMesg(resp, "Check for valid token.")
		case http.StatusForbidden:
			err = httpErrorMesg(resp, "Check for valid token and token user must be an admin")
		default:
			err = httpError(resp)
		}
	}
	return err
}

func httpErrorMesg(resp http.Response, message string) error {
	return fmt.Errorf("HTTP Request %s:%s, HTTP Response: %s. %s",
		resp.Request.Method, resp.Request.URL, resp.Status, message)
}

func httpError(resp http.Response) error {
	return httpErrorMesg(resp, "")
}
