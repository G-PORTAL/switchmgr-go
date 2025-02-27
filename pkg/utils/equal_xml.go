package utils

import (
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

// AssertXMLEquals assets that 2 xml strings are the same by
// stripping away all whitespaces between xml tags
func AssertXMLEquals(t *testing.T, a, b string) {
	re := regexp.MustCompile(`>\s+<`)
	a = re.ReplaceAllString(a, "><")
	b = re.ReplaceAllString(b, "><")

	assert.Equal(t, a, b)
}
