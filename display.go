package conman

import (
	"fmt"
	"os"

	t "github.com/jdrivas/termtext"
	"github.com/juju/ansiterm"
)

// List displpays the list of connections and notes the current one.
func (conns ConnectionList) List() {
	if len(conns) > 0 {
		cn := ""
		if c, err := GetCurrentConnection(); err == nil {
			cn = c.Name
		} // eat the error if we
		w := ansiterm.NewTabWriter(os.Stdout, 4, 4, 3, ' ', 0)
		fmt.Fprintf(w, "%s\n", t.Title("\tName\tURL"))
		for _, c := range conns {
			name := t.Text(c.Name)
			current := ""
			if c.Name == cn {
				name = t.Highlight("%s", c.Name)
				current = t.Highlight("%s", "*")
			}
			fmt.Fprintf(w, "%s\t%s\t%s\n", current, name, t.Text("%s", c.ServiceURL))
		}
		w.Flush()
	} else {
		fmt.Printf("%s\n", t.Title("There were no connections."))
	}
}

func (conns ConnectionList) Describe() {
	if len(conns) > 0 {
		w := ansiterm.NewTabWriter(os.Stdout, 4, 4, 3, ' ', 0)
		fmt.Fprintf(w, describeHeader())
		for _, c := range conns {
			fmt.Fprintf(w, c.describeBody())
		}
		w.Flush()

	} else {
		fmt.Printf("%s\n", t.Title("There were no connections."))
	}

}

func (conn *Connection) Describe() {
	w := ansiterm.NewTabWriter(os.Stdout, 4, 4, 3, ' ', 0)
	fmt.Fprintf(w, describeHeader())
	fmt.Fprintf(w, conn.describeBody())
	w.Flush()

}

func describeHeader() string {
	return t.Title("\tName\tServiceURL\tAuthToken\tHeaders\n")
}

const currentDisplay = "*"

func (conn *Connection) describeBody() (rv string) {
	// First header
	headers := getHeadersDisplay(conn.Headers)
	current := ""
	name := t.Text(conn.Name)
	if cn, err := GetCurrentConnection(); err == nil {
		if conn.Name == cn.Name {
			current = t.Highlight(currentDisplay)
			name = t.Highlight(conn.Name)
		}
	}
	rv += fmt.Sprintf("%s\t%s\t%s\t",
		current, name,
		t.Text("%s\t%s\t%s\n",
			conn.ServiceURL, conn.AuthToken, headers[0]))
	for i := 1; i < len(headers); i++ {
		rv += fmt.Sprintf("\t\t\t%s\n", headers[i])
	}
	return rv
}

// TODO: make the length limit a parameter for viper?
const lengthLimit = 40
const emptyHeader = "<empty>"

// Will always return a list with a first element in it,
// either the actual first header, or the emptyHeader string.
func getHeadersDisplay(hm map[string]string) (hl []string) {
	if len(hm) > 0 {
		for k, v := range hm {
			if len(v) > lengthLimit {
				v = v[:lengthLimit/2-5] + " ... " + v[len(v)-lengthLimit/2+5:]
			}
			hl = append(hl, fmt.Sprintf("%s: %s", t.Title(k), t.Text(v)))
		}
	} else {
		hl = append(hl, emptyHeader)
	}
	return hl
}
