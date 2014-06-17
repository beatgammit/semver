package semver

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var semverReg = regexp.MustCompile("^v?(\\d+).(\\d+).(\\d+)(?:-([0-9A-Za-z-.]+))?(?:\\+([0-9A-Za-z-.]+))?$")

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

func (v *Semver) Validate() error {
	if v.Major < 0 || v.Minor < 0 || v.Patch < 0 {
		return fmt.Errorf("Major, minor and patch version numbers must be non-negative")
	}
	return nil
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
