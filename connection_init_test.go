package conman

import (
	"testing"

	"github.com/spf13/viper"
)

func TestInitialConfig(t *testing.T) {

	type config struct {
		name string
		url  string
	}

	type e struct {
		conName  string
		numConns int
	}

	// Some configs to use
	a := config{"a", "http://127.0.0.1"}
	b := config{"b", "http://127.0.0.1"}

	cases := []struct {
		name          string
		defaultConfig string
		setConfig     string
		configs       []config
		expected      e
	}{
		{
			name:          "Empty Configuration, empty default",
			defaultConfig: "",
			configs:       []config{},
			expected:      e{conName: defaultConnectionName, numConns: 1},
		},
		{
			name:          "One Configuration, empty default",
			defaultConfig: "",
			configs:       []config{a},
			expected:      e{conName: "a", numConns: 1},
		},
		{
			name:          "Two Configurations empty default",
			defaultConfig: "",
			configs:       []config{a, b},
			expected:      e{conName: "a", numConns: 2},
		},
		{
			name:          "One Configuration, with default",
			defaultConfig: "a",
			configs:       []config{a},
			expected:      e{conName: "a", numConns: 1},
		},
		{
			name:          "Two Configurations, with default",
			defaultConfig: "b",
			configs:       []config{a, b},
			expected:      e{conName: "b", numConns: 2},
		},
		{
			name:          "Two Configurations empty default, and setting to the other.",
			defaultConfig: "",
			configs:       []config{a, b},
			setConfig:     "b",
			expected:      e{conName: "a", numConns: 2},
		},
		{
			name:          "Two Configurations, with default, and setting to the other",
			defaultConfig: "b",
			configs:       []config{a, b},
			setConfig:     "a",
			expected:      e{conName: "b", numConns: 2},
		},
		{
			name:          "Two Configurations, with default, and setting to the other (swapped)",
			defaultConfig: "a",
			configs:       []config{a, b},
			setConfig:     "b",
			expected:      e{conName: "a", numConns: 2},
		},
	}

	for _, c := range cases {

		// Set up configuraiton
		for _, cfg := range c.configs {
			setConnectionConfig(cfg.name, cfg.url)
		}

		// Set default
		if c.defaultConfig != "" {
			viper.Set(DefaultConnectionNameKey, c.defaultConfig)
		}

		InitConnections()

		t.Run(c.name, func(t *testing.T) {
			for {

				if cn, err := GetCurrentConnection(); err == nil {
					if cn.Name != c.expected.conName {
						t.Errorf("Checking names, got: %s, expected %s", cn.Name, c.expected.conName)
					}
				} else {
					t.Errorf("Failed to get a current conneciton: %v", err)
				}

				cl := GetAllConnections()
				if len(cl) != c.expected.numConns {
					t.Errorf("Got %d connections, expected %d connections: %#v", len(cl), c.expected.numConns, cl)
				}

				for _, cfg := range c.configs {
					if cn, ok := GetConnection(cfg.name); !ok {
						t.Errorf("Couldn't find expected connection %s with GetConnection()", cfg.name)
						if cn.ServiceURL != cfg.url {
							t.Errorf("Service URL got corrupted in tranist. Got: %s, expected %s", cn.ServiceURL, cfg.url)
						}
					}
				}

				// If we're not doing a set then we're done
				if c.setConfig == "" {
					break
				}

				// Set up for the set and then test.
				t.Logf("Setting new default connection to %s", c.setConfig)
				viper.Set(DefaultConnectionNameKey, c.setConfig)
				c.expected.conName = c.setConfig // and we expect a different return from get curreßnt connection.
				c.setConfig = ""                 // we're done the next time through
			}
		})

		resetConfig() // We can do this either at the top or the bottom.
	}

}
