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
		{`{"major": 1, "minor": 2, "patch": 3}`, Semver{Semver: "1.0.0", Major: 1, Minor: 2, Patch: 3}, "without semver in json"},
	}

	for _, test := range good {
		var ver Semver
		if err := json.Unmarshal([]byte(test.given), &ver); err != nil {
			t.Errorf("%s: %s; given: %s - %#v", test.reason, err, test.given, ver)
		}
	}

	bad := []badJsonTest{
		{`[]`, "invalid json"},
		{`{}`, "empty (no semver)"},
		{`{"semver": {}}`, "wrong type for semver"},
		//{`{"semver": "1.0.0", "major": 2}`, "semver doesn't match"},
		//{`{"semver": "-1.0.0"}`, "fails validation"},
		{`{"semver": "1.0.0", "major": "a"}`, "major not a number"},
		{`{"semver": "1.0.0", "minor": "b"}`, "minor not a number"},
		{`{"semver": "1.0.0", "patch": "c"}`, "patch not a number"},
	}

	for _, test := range bad {
		var ver Semver
		if err := json.Unmarshal([]byte(test.given), &ver); err == nil {
			t.Errorf("%s: given: %s - %#v", test.reason, test.given, ver)
		}
	}
}

type compareTest struct {
	a, b   Semver
	exp    int
	reason string
}

func TestCmp(t *testing.T) {
	tests := []compareTest{
		{Semver{Major: 1, Minor: 0, Patch: 0}, Semver{Major: 0, Minor: 1, Patch: 1}, 1, "major version left"},
		{Semver{Major: 0, Minor: 1, Patch: 1}, Semver{Major: 1, Minor: 0, Patch: 0}, -1, "major version right"},
		{Semver{Major: 1, Minor: 1, Patch: 0}, Semver{Major: 1, Minor: 0, Patch: 1}, 1, "minor version left"},
		{Semver{Major: 1, Minor: 0, Patch: 1}, Semver{Major: 1, Minor: 1, Patch: 0}, -1, "minor version right"},
		{Semver{Major: 1, Minor: 1, Patch: 1}, Semver{Major: 1, Minor: 1, Patch: 0}, 1, "patch version left"},
		{Semver{Major: 1, Minor: 1, Patch: 0}, Semver{Major: 1, Minor: 1, Patch: 1}, -1, "patch version right"},
		{Semver{Major: 1, Minor: 1, Patch: 1}, Semver{Major: 1, Minor: 1, Patch: 1}, 0, "equal (no prerelease)"},

		{Semver{}, Semver{Prerelease: "a"}, 1, "no prerelease trumps prerelease (left)"},
		{Semver{Prerelease: "a"}, Semver{}, -1, "no prerelease trumps prerelease (right)"},

		{Semver{Prerelease: "a"}, Semver{Prerelease: "a"}, 0, "equal prerelease strings"},
		{Semver{Prerelease: "1"}, Semver{Prerelease: "1"}, 0, "equal prerelease numbers"},

		{Semver{Prerelease: "b"}, Semver{Prerelease: "a"}, 1, "string compare (left)"},
		{Semver{Prerelease: "a"}, Semver{Prerelease: "b"}, -1, "string compare (right)"},

		{Semver{Prerelease: "1"}, Semver{Prerelease: "0"}, 1, "number compare (left)"},
		{Semver{Prerelease: "0"}, Semver{Prerelease: "1"}, -1, "number compare (right)"},
		{Semver{Prerelease: "02"}, Semver{Prerelease: "1"}, 1, "number compare two digits (left)"},
		{Semver{Prerelease: "1"}, Semver{Prerelease: "02"}, -1, "number compare two digits (right)"},

		{Semver{Prerelease: "b.1"}, Semver{Prerelease: "a.1"}, 1, "multiple; first (left)"},
		{Semver{Prerelease: "a.1"}, Semver{Prerelease: "b.1"}, -1, "multiple; first (right)"},
		{Semver{Prerelease: "a.2"}, Semver{Prerelease: "a.1"}, 1, "multiple; secnond (left)"},
		{Semver{Prerelease: "a.1"}, Semver{Prerelease: "a.2"}, -1, "multiple; second (right)"},

		{Semver{Prerelease: "a.1"}, Semver{Prerelease: "a"}, 1, "length mismatch (left)"},
		{Semver{Prerelease: "a"}, Semver{Prerelease: "a.1"}, -1, "length mismatch (right)"},

		{Semver{Prerelease: "a"}, Semver{Prerelease: "1"}, 1, "mismatch type (left)"},
		{Semver{Prerelease: "1"}, Semver{Prerelease: "a"}, -1, "mismatch type (right)"},
	}

	for _, test := range tests {
		val := test.a.Cmp(test.b)
		if val != test.exp {
			t.Errorf("%s: %d != %d; a=%s,b=%s", test.reason, val, test.exp, test.a, test.b)
		}
	}
}
