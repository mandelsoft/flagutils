package flagsets

import (
	"encoding/json"
	"strings"
)

// Config is a generic structured config stored in a string map.
type Config = map[string]interface{}

// UnmarshalConfig uses JSON to configure a target object
// with a given Config.
// The object type must match the ConfigOptionTypeSetHandler
// used to create the Config from an OptionSet.
func UnmarshalConfig(cfg Config, target any) error {
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

// ConfigAdder is used to incorporate a partial config into an existing one.
type ConfigAdder func(options Options, config Config) error

func (c ConfigAdder) ApplyConfig(options Options, config Config) error {
	return c(options, config)
}

// ConfigHandler describes the ConfigAdder functionality.
// It is used to apply Options to a Config.
type ConfigHandler interface {
	ApplyConfig(options Options, config Config) error
}

// ConfigOptionTypeSetHandler describes a OptionTypeSet, which also
// provides the possibility to provide config.
type ConfigOptionTypeSetHandler interface {
	OptionTypeSet
	ConfigHandler
}

type configOptionTypeSetHandler struct {
	adder ConfigAdder
	OptionTypeSet
}

// NewConfigOptionTypeSetHandler creates a new ConfigOptionTypeSetHandler based on a ConfigAdder
// and a set of [OptionType]s.
func NewConfigOptionTypeSetHandler(name string, adder ConfigAdder, types ...OptionType) ConfigOptionTypeSetHandler {
	return &configOptionTypeSetHandler{
		adder:         adder,
		OptionTypeSet: NewOptionTypeSet(name, types...),
	}
}

func (c *configOptionTypeSetHandler) ApplyConfig(options Options, config Config) error {
	if c.adder == nil {
		return nil
	}
	return c.adder(options, config)
}

type nopConfigHandler struct{}

// NopConfigHandler is a dummy config handler doing nothing.
var NopConfigHandler = NewNopConfigHandler()

func NewNopConfigHandler() ConfigHandler {
	return &nopConfigHandler{}
}

func (c *nopConfigHandler) ApplyConfig(options Options, config Config) error {
	return nil
}

func FormatOptions(handler OptionTypeSet) string {
	group := ""
	if handler != nil {
		opts := handler.OptionTypeNames()
		var names []string
		if len(opts) > 0 {
			for _, o := range opts {
				names = append(names, "<code>--"+o+"</code>")
			}
			group = "\nOptions used to configure fields: " + strings.Join(names, ", ") + "\n"
		}
	}
	return group
}
