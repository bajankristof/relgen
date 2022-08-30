package semver

import (
	"encoding/json"
	"fmt"
	"github.com/coreos/go-semver/semver"
	"github.com/go-git/go-git/v5/plumbing"
	"strconv"
)

const (
	NONE  = "NONE"
	MAJOR = "MAJOR"
	MINOR = "MINOR"
	PATCH = "PATCH"
)

type version = semver.Version

type Version struct {
	*version
	preRelease *PreRelease
	reference  *plumbing.Reference
	prefix     bool
}

type PreRelease struct {
	Tag    string
	Number int64
}

func NewEmptyVersion() *Version {
	return &Version{&version{}, &PreRelease{Tag: ""}, nil, false}
}

func NewVersion(version string) (*Version, error) {
	vsn := &Version{}
	err := vsn.loadVersion(version)
	if err != nil {
		return nil, err
	}

	return vsn, nil
}

func SelectLatest(versionA *Version, versionB *Version) *Version {
	if versionB == nil {
		return versionA
	} else if versionA == nil || versionA.LessThan(*versionB.version) {
		return versionB
	} else {
		return versionA
	}
}

func SelectGreaterBumpSpec(bumpA string, bumpB string) string {
	switch true {
	case bumpA == MAJOR || bumpB == MAJOR:
		return MAJOR
	case bumpA == MINOR || bumpB == MINOR:
		return MINOR
	case bumpA == PATCH || bumpB == PATCH:
		return PATCH
	default:
		return NONE
	}
}

func (vsn *Version) loadVersion(version string) error {
	prefix := version[0] == 'v'
	if prefix {
		vsn.prefix = prefix
		version = version[1:]
	}

	src, err := semver.NewVersion(version)
	if err != nil {
		return err
	}

	vsn.version = src
	vsn.preRelease = &PreRelease{}
	err = vsn.loadPreRelease()
	return err
}

func (vsn *Version) loadPreRelease() error {
	if vsn.version.PreRelease == "" {
		return nil
	}

	preReleaseParts := vsn.version.PreRelease.Slice()
	switch len(preReleaseParts) {
	case 0:
		return nil
	case 1:
		value, err := strconv.Atoi(preReleaseParts[0])
		vsn.preRelease.Number = int64(value)
		if err != nil {
			vsn.preRelease.Tag = preReleaseParts[0]
		}
		return nil
	default:
		vsn.preRelease.Tag = preReleaseParts[0]
		value, err := strconv.Atoi(preReleaseParts[1])
		vsn.preRelease.Number = int64(value)
		return err
	}
}

func (vsn *Version) PreRelease() PreRelease {
	return *vsn.preRelease
}

func (vsn *Version) IsPreRelease() bool {
	return vsn.preRelease.Tag != "" || vsn.preRelease.Number != 0
}

func (vsn *Version) BumpPreRelease() *Version {
	vsn.preRelease.Number++
	vsn.version.PreRelease = semver.PreRelease(vsn.preRelease.String())
	return vsn
}

func (vsn *Version) MatchPreReleaseTag(tag string) bool {
	return vsn.preRelease.Tag == "" || vsn.preRelease.Tag == tag
}

func (vsn *Version) WithPreReleaseTag(tag string) *Version {
	vsn.preRelease.Tag = tag
	vsn.version.PreRelease = semver.PreRelease(vsn.preRelease.String())
	return vsn
}

func (vsn *Version) WithPreReleaseNumber(number int64) *Version {
	vsn.preRelease.Number = number
	vsn.version.PreRelease = semver.PreRelease(vsn.preRelease.String())
	return vsn
}

func (vsn *Version) Reference() *plumbing.Reference {
	return vsn.reference
}

func (vsn *Version) WithReference(reference *plumbing.Reference) *Version {
	vsn.reference = reference
	return vsn
}

func (vsn *Version) IsReference(hash plumbing.Hash) bool {
	if vsn.reference == nil {
		return false
	}

	return vsn.reference.Hash() == hash
}

func (vsn *Version) Prefix() bool {
	return vsn.prefix
}

func (vsn *Version) WithPrefix(prefix bool) *Version {
	vsn.prefix = prefix
	return vsn
}

func (vsn *Version) BumpWithSpec(bump string) *Version {
	pre := vsn.IsPreRelease()
	switch true {
	case bump == NONE:
		break
	case
		!pre && bump == MAJOR && vsn.Major != 0,
		pre && bump == MAJOR && (vsn.Minor != 0 || vsn.Patch != 0):
		vsn.BumpMajor()
	case
		!pre && bump == MAJOR && vsn.Major == 0,
		!pre && bump == MINOR && vsn.Major != 0,
		pre && bump == MINOR && vsn.Patch != 0:
		vsn.BumpMinor()
	case
		!pre && bump == MINOR && vsn.Major == 0,
		!pre && bump == PATCH:
		vsn.BumpPatch()
	case pre:
		vsn.BumpPreRelease()
	default:
		panic(fmt.Sprintf("unrecognized bump spec \"%s\"", bump))
	}

	vsn.version.PreRelease = semver.PreRelease(vsn.preRelease.String())

	return vsn
}

func (vsn *Version) MarshalJSON() ([]byte, error) {
	return []byte("\"" + vsn.String() + "\""), nil
}

func (vsn *Version) UnmarshalJSON(data []byte) error {
	str := ""
	err := json.Unmarshal(data, &str)
	if err != nil {
		return err
	}

	err = vsn.loadVersion(str)
	if err != nil {
		return err
	}

	return nil
}

func (vsn *Version) String() string {
	if vsn.prefix {
		return "v" + vsn.version.String()
	}

	return vsn.version.String()
}

func (preRelease *PreRelease) String() string {
	switch true {
	case preRelease.Tag == "" && preRelease.Number == 0:
		return ""
	case preRelease.Tag != "" && preRelease.Number != 0:
		return preRelease.Tag + "." + strconv.FormatInt(preRelease.Number, 10)
	case preRelease.Number != 0:
		return strconv.FormatInt(preRelease.Number, 10)
	default:
		return preRelease.Tag
	}
}
