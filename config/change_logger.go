package config

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

var sensitiveKeys = []string{"password", "secret", "token", "key"}

func LogConfigDiff(source string, oldCfg, newCfg map[string]interface{}) {
	diff := diffConfig(oldCfg, newCfg)
	if len(diff) == 0 {
		return
	}

	for key, val := range diff {
		oldVal := maskIfSensitive(key, fmt.Sprintf("%v", val[0]))
		newVal := maskIfSensitive(key, fmt.Sprintf("%v", val[1]))

		var changeType string
		switch {
		case val[0] == nil:
			changeType = "ADD"
		case val[1] == nil:
			changeType = "DELETE"
		default:
			changeType = "UPDATE"
		}

		entry := map[string]interface{}{
			"type":      "config_change",
			"source":    source,
			"key":       key,
			"old":       oldVal,
			"new":       newVal,
			"change":    changeType,
			"timestamp": time.Now().Format(time.RFC3339),
		}

		out, _ := json.Marshal(entry)
		log.Println(string(out))
	}
}

func diffConfig(old, new map[string]interface{}) map[string][2]interface{} {
	result := make(map[string][2]interface{})

	for key, newVal := range new {
		oldVal, ok := old[key]
		if !ok || !deepEqual(oldVal, newVal) {
			result[key] = [2]interface{}{oldVal, newVal}
		}
	}

	for key, oldVal := range old {
		if _, ok := new[key]; !ok {
			result[key] = [2]interface{}{oldVal, nil}
		}
	}

	return result
}

// deepEqual performs deep comparison of two values, handling maps and slices
func deepEqual(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	// Handle map[string]interface{} comparison
	if mapA, okA := a.(map[string]interface{}); okA {
		if mapB, okB := b.(map[string]interface{}); okB {
			return deepEqualMaps(mapA, mapB)
		}
		return false
	}

	// Handle []interface{} comparison
	if sliceA, okA := a.([]interface{}); okA {
		if sliceB, okB := b.([]interface{}); okB {
			return deepEqualSlices(sliceA, sliceB)
		}
		return false
	}

	// For primitive types, use direct comparison
	return a == b
}

// deepEqualMaps compares two map[string]interface{} recursively
func deepEqualMaps(a, b map[string]interface{}) bool {
	if len(a) != len(b) {
		return false
	}

	for key, valA := range a {
		valB, exists := b[key]
		if !exists || !deepEqual(valA, valB) {
			return false
		}
	}

	return true
}

// deepEqualSlices compares two []interface{} recursively
func deepEqualSlices(a, b []interface{}) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if !deepEqual(a[i], b[i]) {
			return false
		}
	}

	return true
}

func maskIfSensitive(key, val string) string {
	lowerKey := strings.ToLower(key)
	for _, s := range sensitiveKeys {
		if strings.Contains(lowerKey, s) {
			return "******"
		}
	}
	return val
}
