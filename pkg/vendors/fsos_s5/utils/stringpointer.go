package utils

import "fmt"

func PtrFromString(v string) *string {
	return &v
}
func StringFromPtr(v *string) string {
	if v == nil {
		return "\"\""
	}
	return fmt.Sprintf("%q", *v)
}
