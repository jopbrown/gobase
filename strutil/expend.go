package strutil

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/jopbrown/gobase/errors"
)

type Expander string
type ExpandHandler func(string) (string, bool)

var globalExpandHandle struct {
	mu      sync.Mutex
	handles []ExpandHandler
}

func init() {
	// RegisterExpandHandler(os.LookupEnv)
}

func RegisterExpandHandler(h ExpandHandler) {
	globalExpandHandle.mu.Lock()
	defer globalExpandHandle.mu.Unlock()
	globalExpandHandle.handles = append(globalExpandHandle.handles, h)
}

func (expander Expander) ExpandByDict(dict map[string]string) string {
	return os.Expand(string(expander), func(key string) string {
		v, ok := dict[key]
		if ok {
			return v
		}

		globalExpandHandle.mu.Lock()
		defer globalExpandHandle.mu.Unlock()
		v, _ = expandByHandlers(key, globalExpandHandle.handles...)
		return v
	})
}

func (expander Expander) ExpandByHandlers(handles ...ExpandHandler) string {
	return os.Expand(string(expander), func(key string) string {
		v, ok := expandByHandlers(key, globalExpandHandle.handles...)
		if ok {
			return v
		}

		globalExpandHandle.mu.Lock()
		defer globalExpandHandle.mu.Unlock()
		v, _ = expandByHandlers(key, globalExpandHandle.handles...)
		return v
	})
}

func (expander Expander) String() string {
	globalExpandHandle.mu.Lock()
	defer globalExpandHandle.mu.Unlock()
	return os.Expand(string(expander), func(key string) string {
		v, _ := expandByHandlers(key, globalExpandHandle.handles...)
		return v
	})
}

func expandByHandlers(key string, handles ...ExpandHandler) (string, bool) {
	for i := len(handles) - 1; i >= 0; i-- {
		v, ok := handles[i](key)
		if ok {
			return v, true
		}
	}

	return "", false
}

func (expander Expander) MarshalYAML() (interface{}, error) {
	return expander.String(), nil
}

func (expander *Expander) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var str string
	err := unmarshal(&str)
	if err != nil {
		return errors.ErrorAt(err)
	}

	*expander = Expander(str)
	return nil
}

func (expander Expander) MarshalJSON() ([]byte, error) {
	return json.Marshal(expander.String())
}

func (expander *Expander) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*expander = Expander(s)
	return nil
}
