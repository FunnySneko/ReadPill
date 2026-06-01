package trash

import "fmt"

func WrapError(text string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", text, err)
}
