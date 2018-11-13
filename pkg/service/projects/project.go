package projects

import (
	"fmt"
	"strings"
	"unicode"
)

func checkJenkinsGoodName(name string) error {
	if strings.TrimSpace(name) == "." || strings.TrimSpace(name) == ".." {
		return fmt.Errorf(". or .. is not a valid jenkins name")
	}

	if strings.IndexAny(name, "?*/\\%!@#$^&|<>[]:;") != -1 {
		return fmt.Errorf("unsafe char in jenkins name")
	}

	for _, r := range name {
		if r > unicode.MaxASCII || !unicode.IsPrint(r) || unicode.IsSpace(r) {
			return fmt.Errorf("name [%s] should be printable ascii code", name)
		}
	}
	return nil
}
