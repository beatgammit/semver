package semver

import (
	"encoding/json"
	"testing"
)

type goodParseTest struct {
	given  string
	exp    Semver
	reason string
}

type badParseTest struct {
	given  string
	reason string
}

func TestParseValid(t *testing.T) {
	tests := []goodParseTest{
		{"1.0.0", Semver{Major: 1, Minor: 0, Patch: 0}, "no prerelease or build"},
		{"1.0.0-test", Semver{Major: 1, Minor: 0, Patch: 0, Prerelease: "test"}, "prerelease but no build"},
		{"1.0.0-test+5334", Semver{Major: 1, Minor: 0, Patch: 0, Prerelease: "test", Build: "5334"}, "prerelease and build"},
		{"1.0.0+5334", Semver{Major: 1, Minor: 0, Patch: 0, Build: "5334"}, "build, but no prerelease"},
	}

	for _, test := range tests {
		test.exp.Semver = test.given
		v, err := Parse(test.given)
		if err != nil {
			t.Errorf("%s: error parsing: %s; given: %s", test.reason, err, test.given)
		} else if v != test.exp {
			t.Errorf("%s: %+v != %+v", test.reason, v, test.exp)
		}
	}
}

func TestParseInvalid(t *testing.T) {
	tests := []badParseTest{
		{"empty", "invalid semver"},
		{"a.0.0", "non-numeric major version"},
		{"0.b.0", "non-numeric minor version"},
		{"0.0.c", "non-numeric patch version"},
		{"-1.0.0", "negative major version"},
		{"0.-1.0", "negative minor version"},
		{"0.0.-1", "negative patch version"},
	}

	for _, test := range tests {
		v, err := Parse(test.given)
		if err == nil {
			t.Errorf("%s: expected error, returned: %+v", test.reason, v)
		}
	}
}

type validateTest struct {
	given  Semver
	reason string
}

func TestValidate(t *testing.T) {
	bad := []validateTest{
		{Semver{"-1.0.0", -1, 0, 0, "", ""}, "negative major version"},
		{Semver{"0.-1.0", 0, -1, 0, "", ""}, "negative minor version"},
		{Semver{"0.0.-1", 0, 0, -1, "", ""}, "negative patch version"},
	}

	for _, test := range bad {
		if err := test.given.Validate(); err == nil {
			t.Errorf("%s: given: %+v", test.reason, test.given)
		}
	}

	good := []validateTest{
		{Semver{"1.0.0", 1, 0, 0, "", ""}, "basic"},
	}

	for _, test := range good {
		if err := test.given.Validate(); err != nil {
			t.Errorf("%s: %s: given: %+v", test.reason, err, test.given)
		}
	}
}

type stringTest struct {
	given       Semver
	exp, reason string
}

func TestString(t *testing.T) {
	tests := []stringTest{
		{Semver{"", 1, 0, 0, "", ""}, "1.0.0", "basic"},
		{Semver{"", 1, 0, 0, "test", ""}, "1.0.0-test", "prerease, no build"},
		{Semver{"", 1, 0, 0, "", "test"}, "1.0.0+test", "build, no prerease"},
		{Semver{"", 1, 0, 0, "blah", "test"}, "1.0.0-blah+test", "prerelease and build"},
	}

	for _, test := range tests {
		s := test.given.String()
		if s != test.exp {
			t.Errorf("%s: %s != %s", test.reason, s, test.exp)
		}
	}
}

type goodJsonTest struct {
	given  string
	exp    Semver
	reason string
}

type badJsonTest struct {
	given, reason string
}

func TestUnmarshalJson(t *testing.T) {
	good := []goodJsonTest{
		{`{"semver": "1.0.0"}`, Semver{Semver: "1.0.0", Major: 1}, "basic"},
		{`{"semver": "1.2.3", "major": 1, "minor": 2, "patch": 3}`, Semver{Semver: "1.0.0", Major: 1, Minor: 2, Patch: 3}, "ints parse"},
		{`{"semver": "1.2.3", "major": "1", "minor": "2", "patch": "3"}`, Semver{Semver: "1.0.0", Major: 1, Minor: 2, Patch: 3}, "strings parse"},
		{`{"semver": "1.0.0-test"}`, Semver{Semver: "1.0.0", Major: 1, Prerelease: "test"}, "prerelease"},
		{`{"semver": "1.0.0+test"}`, Semver{Semver: "1.0.0", Major: 1, Build: "build"}, "build"},
		{`{"semver": "1.0.0-blah+test"}`, Semver{Semver: "1.0.0", Major: 1, Prerelease: "blah", Build: "build"}, "prerelease and build"},
	}

	for _, test := range good {
		var ver Semver
		if err := json.Unmarshal([]byte(test.given), &ver); err != nil {
			t.Errorf("%s: %s; given: %s", test.reason, err, test.given)
		}
	}

	bad := []badJsonTest{
		{`[]`, "invalid json"},
		{`{}`, "empty (no semver)"},
		{`{"semver": {}}`, "wrong type for semver"},
		{`{"semver": "1.0.0", "major": 2}`, "semver doesn't match"},
		{`{"semver": "-1.0.0"}`, "fails validation"},
		{`{"semver": "1.0.0", "major": "a"}`, "major not a number"},
		{`{"semver": "1.0.0", "minor": "b"}`, "minor not a number"},
		{`{"semver": "1.0.0", "patch": "c"}`, "patch not a number"},
	}

	for _, test := range bad {
		var ver Semver
		if err := json.Unmarshal([]byte(test.given), &ver); err == nil {
			t.Errorf("%s: given: %s", test.reason, test.given)
		}
	}
}
