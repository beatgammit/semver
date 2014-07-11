// Package semver implements a semantic version parser according to
// this spec: http://semver.org
//
// This package allows prefixing version numbers with a v.
//
// Example:
//
//   ver, _ := semver.Parse("v1.2.3")
//   fmt.Println(ver) // prints 1.2.3
package semver
