package app

import (
	"fmt"
	"strings"
)

func (a *APITest) iGenerateFakeData(dataSpec string) error {
	pairs := splitByCommaOutsideBrackets(dataSpec)

	for _, pair := range pairs {
		parts := strings.SplitN(strings.TrimSpace(pair), "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid data specification format: %s", pair)
		}

		varName := strings.TrimSpace(parts[0])
		pattern := strings.TrimSpace(parts[1])

		value, err := generateFromTag(pattern)
		if err != nil {
			return fmt.Errorf("error generating fake data for %s: %w", varName, err)
		}

		a.store[varName] = value
	}

	return nil
}
