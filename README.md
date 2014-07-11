semver [![Build Status](https://travis-ci.org/beatgammit/semver.png)](https://travis-ci.org/beatgammit/semver)
======

semver is a [semantic versioning](http://semver.org/) library for Go.

example
-------

    package main

    import "github.com/beatgammit/semver"

    func main() {
        ver, err := semver.Parse("1.2.3-beta+jp")
        if err != nil {
            panic(err)
        }
        println(ver.String())
    }

api
===

[godoc](https://godoc.org/github.com/beatgammit/semver)

`semver.Semver` implements the following interfaces from the standard libary:

* `encoding.TextMarshaler`
* `encoding.TextUnmarshaler`
* `json.Marshaler`
* `json.Unmarshaler`
