package rbac

import "path"

func matchesPattern(pattern, value string) bool {
	if pattern == "" || pattern == "*" {
		return true
	}

	match, err := path.Match(pattern, value)
	if err != nil {
		return pattern == value
	}
	return match
}
