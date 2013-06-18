semver [![Build Status](https://travis-ci.org/beatgammit/semver.png)](https://travis-ci.org/beatgammit/semver)
======

semver is a [semantic versioning](http://semver.org/) library for Go. It provides a lenient JSON parser that accepts strings or numbers for major, minor and patch.

api
===

package functions
-----------------

* `Parse(string) (Semver, error)`- parses and validates a Semver string

Semver
------

* `Compare(Semver) int`- compares two semvers
* `String() string`- constructs a semver string
* `UnmarshalJSON([]byte) error`- for `encoding/json` compatibility
* `Validate() error`- checks that the Semver struct is sane
