package semver

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var semverReg = regexp.MustCompile("^(\\d+).(\\d+).(\\d+)(?:-([0-9A-Za-z-.]+))?(?:\\+([0-9A-Za-z-.]+))?$")

type Semver struct {
	Semver     string `json:"semver"`
	Major      int    `json:"major"`
	Minor      int    `json:"minor"`
	Patch      int    `json:"patch"`
	Prerelease string `json:"prerelease,omitempty"`
	Build      string `json:"build,omitempty"`
}

func Parse(semver string) (v Semver, err error) {
	pieces := semverReg.FindStringSubmatch(semver)
	if pieces == nil {
		err = fmt.Errorf("Invalid semver string: %s", semver)
		return
	}
	// will always be a number, but we're explicitly not checking for out of bounds errors
	v.Major, _ = strconv.Atoi(pieces[1])
	v.Minor, _ = strconv.Atoi(pieces[2])
	v.Patch, _ = strconv.Atoi(pieces[3])
	v.Prerelease = pieces[4]
	v.Build = pieces[5]
	v.Semver = semver
	err = v.Validate()
	return
}

func (v *Semver) String() string {
	s := fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	if v.Prerelease != "" {
		s += "-" + v.Prerelease
	}
	if v.Build != "" {
		s += "+" + v.Build
	}
	return s
}

func (v *Semver) Validate() error {
	if v.Major < 0 || v.Minor < 0 || v.Patch < 0 {
		return fmt.Errorf("Major, minor and patch version numbers must be non-negative")
	}
	return nil
}

func (ver *Semver) UnmarshalJSON(arr []byte) (err error) {
	var tmap map[string]interface{}
	if err = json.Unmarshal(arr, &tmap); err != nil {
		return
	}

	fmt.Println(tmap)

	rVal := reflect.ValueOf(ver)
	for k, v := range tmap {
		field := rVal.Elem().FieldByName(strings.Title(k))
		valType := reflect.TypeOf(v)
		if valType.AssignableTo(field.Type()) {
			field.Set(reflect.ValueOf(v))
		} else if valType.ConvertibleTo(field.Type()) {
			field.Set(reflect.ValueOf(v).Convert(field.Type()))
		} else {
			// we'll only get here for Major, Minor & Patch
			if valType.Kind() == reflect.String {
				var val int
				val, err = strconv.Atoi(v.(string))
				if err != nil {
					return
				}
				field.SetInt(int64(val))
			}
		}
	}

	if ver.Semver == "" {
		return fmt.Errorf("semver must not be empty")
	}

	if ver.Major == 0 && ver.Minor == 0 && ver.Patch == 0 {
		*ver, err = Parse(ver.Semver)
	}

	if ver.String() != ver.Semver {
		return fmt.Errorf("semver must match parsed version")
	}

	return ver.Validate()
}
