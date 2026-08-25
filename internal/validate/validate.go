package validate

import (
	"fmt"
	"strings"
)

func NonEmpty(name, v string) error {
	if strings.TrimSpace(v) == "" {
		return fmt.Errorf("%s: empty", name)
	}
	return nil
}
