package main

import (
	"errors"
	"fmt"
	"strings"
)

func Split(s string, d string) []string {
	arr := strings.Split(s, d)
	return arr
}

// https://go.dev/play/p/5Hiajt4H2Z5
func getFunctions() map[string]any {
	funcs := map[string]any{
		"map": func(pairs ...any) (map[string]any, error) {
			if len(pairs)%2 != 0 {
				return nil, errors.New("misaligned map")
			}

			m := make(map[string]any, len(pairs)/2)

			for i := 0; i < len(pairs); i += 2 {
				key, ok := pairs[i].(string)

				if !ok {
					return nil, fmt.Errorf("cannot use type %T as map key", pairs[i])
				}
				m[key] = pairs[i+1]
			}
			return m, nil
		},
		"split": Split,
	}
	return funcs
}
