package utils

import (
	"encoding/xml"
	"reflect"
)

// CompareXMLIgnoreWhitespace compares two XML strings and ignores whitespace
func CompareXMLIgnoreWhitespace(xml1, xml2 string) (bool, error) {
	var obj1 interface{}
	var obj2 interface{}

	err := xml.Unmarshal([]byte(xml1), &obj1)
	if err != nil {
		return false, err
	}

	err = xml.Unmarshal([]byte(xml2), &obj2)
	if err != nil {
		return false, err
	}

	return reflect.DeepEqual(obj1, obj2), nil
}
