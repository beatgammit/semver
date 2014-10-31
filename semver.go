package semver

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var semverReg = regexp.MustCompile(`^v?(\d+).(\d+).(\d+)(?:-([0-9A-Za-z-.]+))?(?:\+([0-9A-Za-z-.]+))?$`)

type Semver struct {
	Major      int
	Minor      int
	Patch      int
	Prerelease string
	Build      string
}

// Parse parses semver into a Semver. A leading v may be included.
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
	err = v.Validate()
	return
}

// MustParse parses semver into a semver. It will panic if there is an error in parsing.
func MustParse(semver string) Semver {
	if ver, err := Parse(semver); err != nil {
		panic(err)
	} else {
		return ver
	}
}

// String produces a Semver string.
func (v Semver) String() string {
	s := fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	if v.Prerelease != "" {
		s += "-" + v.Prerelease
	}
	if v.Build != "" {
		s += "+" + v.Build
	}
	return s
}

// Validate checks a semver for appropriate values.
// The outputs of String(), MarshalJSON and MarshalText are only guaranteed to be valid
// semvers if this function does not return an error.
func (v Semver) Validate() error {
	if v.Major < 0 || v.Minor < 0 || v.Patch < 0 {
		return fmt.Errorf("Major, minor and patch version numbers must be non-negative")
	} else if v.Major == 0 && v.Minor == 0 && v.Patch == 0 {
		return fmt.Errorf("Must supply at least one of: major, minor, patch")
	}
	return nil
}

func (ver Semver) MarshalJSON() ([]byte, error) {
	b, err := ver.MarshalText()
	return []byte(fmt.Sprintf(`"%s"`, string(b))), err
}

func (ver *Semver) UnmarshalJSON(arr []byte) error {
	for _, c := range arr {
		// TODO: this is completely gross (backwards compatibility for older version)
		if !unicode.IsSpace(rune(c)) {
			if c == '{' {
				// we can't just unmarshal into a Semver because
				// we'd end up with infinite recursion
				type semver Semver
				var sem semver
				if err := json.Unmarshal(arr, &sem); err != nil {
					return err
				}
				*ver = Semver(sem)
				return ver.Validate()
			}
			break
		}
	}
	return ver.UnmarshalText(arr[1 : len(arr)-1])
}

func (ver *Semver) UnmarshalText(arr []byte) error {
	v, err := Parse(string(arr))
	if err == nil {
		*ver = v
	}
	return err
}

func (ver Semver) MarshalText() ([]byte, error) {
	return []byte(ver.String()), nil
}

// Cmp compares two semantic versions:
// - < 0 if a < b
// - > 0 if a > b
// - == 0 if a == b
//
// In order of importance: Major > Minor > Patch > Prerelease (Build ignored)
//
// Major, Minor and Patch are compared numerically.
// Prerelease is compared by splitting on the . and:
// - comparing identifiers lexically (in ASCII sort order)
// - comparing numeric identifiers numerically
// Numeric identifiers have lower precedence
func (a Semver) Cmp(b Semver) int {
	if a.Major != b.Major {
		return a.Major - b.Major
	}
	if a.Minor != b.Minor {
		return a.Minor - b.Minor
	}
	if a.Patch != b.Patch {
		return a.Patch - b.Patch
	}

	if a.Prerelease == "" {
		if b.Prerelease == "" {
			return 0
		}
		return 1
	} else if b.Prerelease == "" {
		// a.Prerelease != ""
		return -1
	}

	partsA := strings.Split(a.Prerelease, ".")
	partsB := strings.Split(b.Prerelease, ".")
	total := len(partsA)
	if len(partsB) < total {
		total = len(partsB)
	}
	for i := 0; i < total; i++ {
		sa, sb := partsA[i], partsB[i]
		ai, errA := strconv.Atoi(sa)
		bi, errB := strconv.Atoi(sb)

		if errA != nil && errB != nil {
			if sa < sb {
				return -1
			} else if sa > sb {
				return 1
			}
		} else if errA == nil && errB == nil {
			if ai != bi {
				return ai - bi
			}
		} else if errA != nil {
			return 1
		} else if errB != nil {
			return -1
		}
	}

	return len(partsA) - len(partsB)
}
