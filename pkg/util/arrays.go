// Copyright 2019 Red Hat, Inc. and/or its affiliates
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import (
	"fmt"
	"strings"
)

const keyPairSeparator = "="

// Contains checks if the s string are within the array
func Contains(s string, array []string) bool {
	if len(s) == 0 {
		return false
	}
	for _, item := range array {
		if s == item {
			return true
		}
	}
	return false
}

// AppendStringMap will append to dest the source map
func AppendStringMap(dest map[string]string, source map[string]string) map[string]string {
	if dest == nil {
		dest = map[string]string{}
	}

	for key, value := range source {
		dest[key] = value
	}
	return dest
}

// FromStringsKeyPairToMap converts a string array in the key/pair format (key=value) to a map. Unconvertable strings will be skipped.
func FromStringsKeyPairToMap(array []string) map[string]string {
	if array == nil || len(array) == 0 {
		return nil
	}
	kp := map[string]string{}
	for _, item := range array {
		spplited := strings.SplitN(item, keyPairSeparator, 2)
		if len(spplited) == 0 {
			break
		}

		if len(spplited[0]) == 0 {
			break
		}

		if len(spplited) == 2 {
			kp[spplited[0]] = spplited[1]
		} else if len(spplited) == 1 {
			kp[spplited[0]] = ""
		}
	}
	return kp
}

// ParseStringsForKeyPair will parse the given string array for a valid key=pair format on each item.
// Returns an error if any item is not in the valid format.
func ParseStringsForKeyPair(array []string) error {
	if array == nil || len(array) == 0 {
		return nil
	}
	for _, item := range array {
		if !strings.Contains(item, keyPairSeparator) {
			return fmt.Errorf("Item %s does not contain the key/pair separator '%s'", item, keyPairSeparator)
		}
		if strings.HasPrefix(item, keyPairSeparator) {
			return fmt.Errorf("Item %s starts with key/pair separator '%s'", item, keyPairSeparator)
		}
	}
	return nil
}

// ArrayToSet converts an array of string to a set
func ArrayToSet(array []string) map[string]bool {
	set := make(map[string]bool, len(array))

	for _, e := range array {
		set[e] = true
	}

	return set
}

// ContainsAll checks if all the elements of the second are in the first array
func ContainsAll(array1 []string, array2 []string) bool {
	set1 := ArrayToSet(array1)

	for _, e := range array2 {
		if !set1[e] {
			return false
		}
	}

	return true
}
