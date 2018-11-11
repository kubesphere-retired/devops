package projects

import (
	"fmt"
	"strings"
)

func checkJenkinsGoodName(name string) error {
	if strings.TrimSpace(name) == "." || strings.TrimSpace(name) == ".." {
		return fmt.Errorf(". or .. is not a valid jenkins name")
	}

	if strings.IndexAny(name, "?*/\\%!@#$^&|<>[]:;") != -1 {
		return fmt.Errorf("unsafe char in jenkins name")
	}
	return nil
}
