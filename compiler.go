package micro

import "regexp"

// Used to strip the formatting stuff when running directly through 'go run'.
var stripFormatting = regexp.MustCompile("\\$\\{[^\\}]+\\}")

// errorMessageRe is a regex to find lines that look like they're specifying a file.
var errorMessageRe = regexp.MustCompile(`^([^ ]+\.[^: /]+):([0-9]+):(?:([0-9]+):)? *(?:([a-z-_ ]+):)? (.*)$`)
