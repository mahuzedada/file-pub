package common

import (
	"fmt"
	"reflect"
)

// RequireNonNil panics if a dependency is nil
func RequireNonNil(dependency interface{}, name string) {
	if dependency == nil || reflect.ValueOf(dependency).IsNil() {
		panic(fmt.Sprintf("%s is required", name))
	}
}

// PanicOnInvalidDependencies validates all dependencies and panics with descriptive error
func PanicOnInvalidDependencies(serviceName string, dependencies map[string]interface{}) {
	for name, dep := range dependencies {
		if dep == nil || reflect.ValueOf(dep).IsNil() {
			panic(fmt.Sprintf("%s: dependency '%s' is required", serviceName, name))
		}
	}
}
