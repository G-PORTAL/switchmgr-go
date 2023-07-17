package iosconfig

import (
	"regexp"
	"strconv"
	"strings"
)

type ConfigValues []string

func (cv ConfigValues) GetStringValue(config string, def string) string {
	for _, line := range cv {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, config) {
			return strings.TrimSpace(strings.TrimPrefix(line, config))
		}
	}

	return def
}

// Exists checks if a config value exists
func (cv ConfigValues) Exists(config string, exact bool) bool {
	for _, line := range cv {
		line = strings.TrimSpace(line)

		if exact && line == config {
			return true
		} else if !exact && strings.HasPrefix(line, config) {
			return true
		}
	}

	return false
}

func (cv ConfigValues) GetIntValue(config string, def int) int {
	value := cv.GetStringValue(config, strconv.Itoa(def))
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return def
	}

	return intValue
}

func (cv ConfigValues) GetInt32Value(config string, def int32) int32 {
	return int32(cv.GetIntValue(config, int(def)))
}

func (cv ConfigValues) GetStringValues(config string, def []string) []string {
	for _, line := range cv {
		if strings.HasPrefix(line, config) {
			return strings.Split(strings.TrimSpace(strings.TrimPrefix(line, config)), ",")
		}
	}

	return def
}

func (cv ConfigValues) GetIntValues(config string, def []int) []int {
	defaults := make([]string, len(def))
	for i, d := range def {
		defaults[i] = strconv.Itoa(d)
	}

	values := make([]int, 0)
	stringValues := cv.GetStringValues(config, defaults)
	for _, value := range stringValues {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return def
		}

		values = append(values, intValue)
	}

	return values
}

func (cv ConfigValues) GetInt32Values(config string, def []int32) []int32 {
	defaults := make([]int, len(def))
	for i, d := range def {
		defaults[i] = int(d)
	}

	intValues := cv.GetIntValues(config, defaults)
	values := make([]int32, len(intValues))
	for i, value := range intValues {
		values[i] = int32(value)
	}

	return values
}

type Config map[string]ConfigValues

var parentRgx = regexp.MustCompile(`^[a-z]+`)

func Parse(config string) (Config, error) {
	parts := strings.Split(config, "!")

	cfg := Config{}

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		parent := ""
		lines := strings.Split(part, "\n")
		for _, line := range lines {
			if parentRgx.MatchString(line) {
				parent = strings.TrimSpace(line)
				cfg[parent] = nil
			}

			if !parentRgx.MatchString(line) && parent != "" {
				if cfg[parent] == nil {
					cfg[parent] = make(ConfigValues, 0)
				}

				cfg[parent] = append(cfg[parent], strings.TrimSpace(line))
			}
		}
	}

	// TODO: sort the config by key value of the map

	return cfg, nil
}
