package semver

import (
	"github.com/go-git/go-git/v5/plumbing"
	"testing"
)

type bumpWithSpecTest struct {
	version string
	bump    string
	expect  string
}

func TestNewEmptyVersion(t *testing.T) {
	vsn := NewEmptyVersion()
	switch true {
	case vsn.Major != 0:
		t.Fatalf(`NewEmptyVersion() = %v, expected major to be 0, got %d`, vsn, vsn.Major)
	case vsn.Minor != 0:
		t.Fatalf(`NewEmptyVersion() = %v, expected minor to be 0, got %d`, vsn, vsn.Minor)
	case vsn.Patch != 0:
		t.Fatalf(`NewEmptyVersion() = %v, expected patch to be 0, got %d`, vsn, vsn.Patch)
	case vsn.preRelease.Tag != "":
		t.Fatalf(`NewEmptyVersion() = %v, expected pre-release tag to be "", got "%s"`, vsn, vsn.preRelease.Tag)
	case vsn.preRelease.Number != 0:
		t.Fatalf(`NewEmptyVersion() = %v, expected pre-release number to be 0, got %d`, vsn, vsn.preRelease.Number)
	case vsn.Metadata != "":
		t.Fatalf(`NewEmptyVersion() = %v, expected build-metadata to be "", got "%s"`, vsn, vsn.Metadata)
	}
}

func TestNewVersion(t *testing.T) {
	src := "1.22.3-foo.13+bar"
	vsn, err := NewVersion(src)
	switch true {
	case err != nil:
		t.Fatalf(`NewVersion("%s") = (%v, %v), expected error to be <nil>, got %v`, src, vsn, err, err)
	case vsn.Major != 1:
		t.Fatalf(`NewVersion("%s") = (%v, %v), expected major to be 1, got %d`, src, vsn, err, vsn.Major)
	case vsn.Minor != 22:
		t.Fatalf(`NewVersion("%s") = (%v, %v), expected minor to be 22, got %d`, src, vsn, err, vsn.Minor)
	case vsn.Patch != 3:
		t.Fatalf(`NewVersion("%s") = (%v, %v), expected patch to be 3, got %d`, src, vsn, err, vsn.Patch)
	case vsn.preRelease.Tag != "foo":
		t.Fatalf(`NewVersion("%s") = (%v, %v), expected pre-release tag to be "foo", got "%s"`, src, vsn, err, vsn.preRelease.Tag)
	case vsn.preRelease.Number != 13:
		t.Fatalf(`NewVersion("%s") = (%v, %v), expected pre-release number to be 13, got %d`, src, vsn, err, vsn.preRelease.Number)
	case vsn.Metadata != "bar":
		t.Fatalf(`NewVersion("%s") = (%v, %v), expected build-metadata to be "bar", got "%s"`, src, vsn, err, vsn.Metadata)
	case vsn.prefix != false:
		t.Fatalf(`NewVersion("%s") = (%v, %v), expected prefix to be false, got %v`, src, vsn, err, vsn.prefix)
	}

	src = "v1.0.0-foo"
	vsn, err = NewVersion(src)
	switch true {
	case err != nil:
		t.Fatalf(`NewVersion("%s") = (%v, %v), expected error to be <nil>, got %v`, src, vsn, err, err)
	case vsn.preRelease.Tag != "foo":
		t.Fatalf(`NewVersion("%s") = (%v, %v), expected pre-release tag to be "foo", got "%s"`, src, vsn, err, vsn.preRelease.Tag)
	case vsn.preRelease.Number != 0:
		t.Fatalf(`NewVersion("%s") = (%v, %v), expected pre-release number to be 0, got "%d"`, src, vsn, err, vsn.preRelease.Number)
	case vsn.prefix != true:
		t.Fatalf(`NewVersion("%s") = (%v, %v), expected prefix to be true, got %v`, src, vsn, err, vsn.prefix)
	}

	src = "1.0.0-999"
	vsn, err = NewVersion(src)
	switch true {
	case err != nil:
		t.Fatalf(`NewVersion("%s") = (%v, %v), expected error to be <nil>, got %v`, src, vsn, err, err)
	case vsn.preRelease.Tag != "":
		t.Fatalf(`NewVersion("%s") = (%v, %v), expected pre-release tag to be "", got "%s"`, src, vsn, err, vsn.preRelease.Tag)
	case vsn.preRelease.Number != 999:
		t.Fatalf(`NewVersion("%s") = (%v, %v), expected pre-release number to be 999, got "%d"`, src, vsn, err, vsn.preRelease.Number)
	}

	src = "1.0.0-foo.bar"
	vsn, err = NewVersion("1.0.0-foo.bar")
	switch true {
	case vsn != nil:
		t.Fatalf(`NewVersion("%s") = (%v, %v), expected <nil>, got %v`, src, vsn, err, vsn)
	case err == nil:
		t.Fatalf(`NewVersion("%s") = (%v, %v), expected error to NOT be <nil>`, src, vsn, err)
	}

	vsn, err = NewVersion("foo")
	switch true {
	case vsn != nil:
		t.Fatalf(`NewVersion("%s") = (%v, %v), expected <nil>, got %v`, src, vsn, err, vsn)
	case err == nil:
		t.Fatalf(`NewVersion("%s") = (%v, %v), expected error to NOT be <nil>`, src, vsn, err)
	}
}

func TestSelectLatest(t *testing.T) {
	var vsnA, vsnB, latest *Version
	vsnA, _ = NewVersion("1.0.0")
	vsnB, _ = NewVersion("1.0.0-beta")

	latest = SelectLatest(vsnA, vsnB)
	if latest != vsnA {
		t.Fatalf(`SelectLatest(%v, %v), expected %v, got %v`, vsnA, vsnB, vsnA, latest)
	}

	vsnB, _ = NewVersion("1.1.0+foo")
	latest = SelectLatest(vsnA, vsnB)
	if latest != vsnB {
		t.Fatalf(`SelectLatest(%v, %v), expected %v, got %v`, vsnA, vsnB, vsnB, latest)
	}

	vsnB = nil
	latest = SelectLatest(vsnA, vsnB)
	if latest != vsnA {
		t.Fatalf(`SelectLatest(%v, %v), expected %v, got %v`, vsnA, vsnB, vsnA, latest)
	}
}

func TestSelectGreaterBumpSpec(t *testing.T) {
	switch true {
	case SelectGreaterBumpSpec(MAJOR, MAJOR) != MAJOR:
		t.Fatalf(`SelectGreaterBumpSpec("%s", "%s"), expected "%s"`, MAJOR, MAJOR, MAJOR)
	case SelectGreaterBumpSpec(MAJOR, MINOR) != MAJOR:
		t.Fatalf(`SelectGreaterBumpSpec("%s", "%s"), expected "%s"`, MAJOR, MINOR, MAJOR)
	case SelectGreaterBumpSpec(MAJOR, PATCH) != MAJOR:
		t.Fatalf(`SelectGreaterBumpSpec("%s", "%s"), expected "%s"`, MAJOR, PATCH, MAJOR)
	case SelectGreaterBumpSpec(MAJOR, NONE) != MAJOR:
		t.Fatalf(`SelectGreaterBumpSpec("%s", "%s"), expected "%s"`, MAJOR, NONE, MAJOR)
	case SelectGreaterBumpSpec(MINOR, MAJOR) != MAJOR:
		t.Fatalf(`SelectGreaterBumpSpec("%s", "%s"), expected "%s"`, MINOR, MAJOR, MAJOR)
	case SelectGreaterBumpSpec(PATCH, MAJOR) != MAJOR:
		t.Fatalf(`SelectGreaterBumpSpec("%s", "%s"), expected "%s"`, PATCH, MAJOR, MAJOR)
	case SelectGreaterBumpSpec(NONE, MAJOR) != MAJOR:
		t.Fatalf(`SelectGreaterBumpSpec("%s", "%s"), expected "%s"`, NONE, MAJOR, MAJOR)

	case SelectGreaterBumpSpec(MINOR, MINOR) != MINOR:
		t.Fatalf(`SelectGreaterBumpSpec("%s", "%s"), expected "%s"`, MAJOR, MAJOR, MAJOR)
	case SelectGreaterBumpSpec(MINOR, PATCH) != MINOR:
		t.Fatalf(`SelectGreaterBumpSpec("%s", "%s"), expected "%s"`, MINOR, PATCH, MINOR)
	case SelectGreaterBumpSpec(MINOR, NONE) != MINOR:
		t.Fatalf(`SelectGreaterBumpSpec("%s", "%s"), expected "%s"`, MINOR, NONE, MINOR)
	case SelectGreaterBumpSpec(PATCH, MINOR) != MINOR:
		t.Fatalf(`SelectGreaterBumpSpec("%s", "%s"), expected "%s"`, PATCH, MINOR, MINOR)
	case SelectGreaterBumpSpec(NONE, MINOR) != MINOR:
		t.Fatalf(`SelectGreaterBumpSpec("%s", "%s"), expected "%s"`, NONE, MINOR, MINOR)

	case SelectGreaterBumpSpec(PATCH, PATCH) != PATCH:
		t.Fatalf(`SelectGreaterBumpSpec("%s", "%s"), expected "%s"`, MAJOR, MAJOR, MAJOR)
	case SelectGreaterBumpSpec(PATCH, NONE) != PATCH:
		t.Fatalf(`SelectGreaterBumpSpec("%s", "%s"), expected "%s"`, PATCH, NONE, PATCH)
	case SelectGreaterBumpSpec(NONE, PATCH) != PATCH:
		t.Fatalf(`SelectGreaterBumpSpec("%s", "%s"), expected "%s"`, NONE, PATCH, PATCH)

	case SelectGreaterBumpSpec(NONE, NONE) != NONE:
		t.Fatalf(`SelectGreaterBumpSpec("%s", "%s"), expected "%s"`, NONE, NONE, NONE)
	}
}

func TestVersion_PreRelease(t *testing.T) {
	vsn := NewEmptyVersion()
	preRelease := vsn.PreRelease()

	if preRelease != *vsn.preRelease {
		t.Fatalf(`(*Version(%v)).PreRelease(), expected %v, got %v`, vsn, *vsn.preRelease, preRelease)
	}
}

func TestVersion_IsPreRelease(t *testing.T) {
	vsn := NewEmptyVersion()

	if vsn.IsPreRelease() != false {
		t.Fatalf(`(*Version(%v)).IsPreRelease(), expected false, got true`, vsn)
	}

	vsn.preRelease = &PreRelease{Tag: "beta"}
	if vsn.IsPreRelease() != true {
		t.Fatalf(`(*Version(%v)).IsPreRelease(), expected true, got false`, vsn)
	}

	vsn.preRelease = &PreRelease{Tag: "", Number: 99}
	if vsn.IsPreRelease() != true {
		t.Fatalf(`(*Version(%v)).IsPreRelease(), expected true, got false`, vsn)
	}
}

func TestVersion_BumpPreRelease(t *testing.T) {
	vsn := NewEmptyVersion()
	vsn.BumpPreRelease()

	if vsn.preRelease.Number != 1 {
		t.Fatalf(`(*Version(%v)).BumpPreRelease(), expected pre-release number to be 1, got %d`, vsn, vsn.preRelease.Number)
	}
}

func TestVersion_MatchPreReleaseTag(t *testing.T) {
	vsn := NewEmptyVersion()

	if vsn.MatchPreReleaseTag("") != true {
		t.Fatalf(`(*Version(%v)).MatchPreReleaseTag(""), expected true, got false`, vsn)
	}

	if vsn.MatchPreReleaseTag("foo") != true {
		t.Fatalf(`(*Version(%v)).MatchPreReleaseTag("foo"), expected true, got false`, vsn)
	}

	if vsn.MatchPreReleaseTag("bar") != true {
		t.Fatalf(`(*Version(%v)).MatchPreReleaseTag("bar"), expected true, got false`, vsn)
	}

	vsn.preRelease = &PreRelease{Tag: "foo"}

	if vsn.MatchPreReleaseTag("") != false {
		t.Fatalf(`(*Version(%v)).MatchPreReleaseTag(""), expected false, got true`, vsn)
	}

	if vsn.MatchPreReleaseTag("foo") != true {
		t.Fatalf(`(*Version(%v)).MatchPreReleaseTag("foo"), expected true, got false`, vsn)
	}

	if vsn.MatchPreReleaseTag("bar") != false {
		t.Fatalf(`(*Version(%v)).MatchPreReleaseTag("bar"), expected false, got true`, vsn)
	}
}

func TestVersion_WithPreReleaseTag(t *testing.T) {
	tag := "foo"
	vsn := NewEmptyVersion().WithPreReleaseTag(tag)

	if vsn.preRelease.Tag != tag {
		t.Fatalf(`(*Version(%v)).WithPreReleaseTag("%s"), expected pre-release tag to be "%s", got "%s"`, vsn, tag, tag, vsn.preRelease.Tag)
	}
}

func TestVersion_WithPreReleaseNumber(t *testing.T) {
	var num int64 = 99
	vsn := NewEmptyVersion().WithPreReleaseNumber(num)

	if vsn.preRelease.Number != num {
		t.Fatalf(`(*Version(%v)).WithPreReleaseNumber(%d), expected pre-release number to be %d, got %d`, vsn, num, num, vsn.preRelease.Number)
	}
}

func TestVersion_Reference(t *testing.T) {
	vsn := NewEmptyVersion()
	vsn.reference = &plumbing.Reference{}
	ref := vsn.Reference()

	if ref != vsn.reference {
		t.Fatalf(`(*Version(%v)).Reference(), expected %v, got %v`, vsn, vsn.reference, ref)
	}
}

func TestVersion_WithReference(t *testing.T) {
	ref := &plumbing.Reference{}
	vsn := NewEmptyVersion().WithReference(ref)

	if ref != vsn.reference {
		t.Fatalf(`(*Version(%v)).WithReference(%v), expected reference to be %v, got %v`, vsn, ref, ref, vsn.reference)
	}
}

func TestVersion_IsReference(t *testing.T) {
	hash := plumbing.NewHash("123")
	ref := plumbing.NewHashReference("foo", hash)
	vsn := NewEmptyVersion()

	if vsn.IsReference(hash) {
		t.Fatalf(`(*Version(%v)).IsReference(%v), expected false, got true`, vsn, hash)
	}

	if !vsn.WithReference(ref).IsReference(ref.Hash()) {
		t.Fatalf(`(*Version(%v)).IsReference(%v), expected true, got false`, vsn, ref.Hash())
	}

	other := plumbing.NewHash("456")
	if vsn.IsReference(other) {
		t.Fatalf(`(*Version(%v)).IsReference(%v), expected false, got true`, vsn, other)
	}
}

func TestVersion_Prefix(t *testing.T) {
	vsn := NewEmptyVersion()
	vsn.prefix = true

	if !vsn.Prefix() {
		t.Fatalf(`(*Version(%v)).Prefix(), expected true, got false`, vsn)
	}
}

func TestVersion_WithPrefix(t *testing.T) {
	vsn := NewEmptyVersion()
	vsn.WithPrefix(true)

	if !vsn.prefix {
		t.Fatalf(`(*Version(%v)).WithPrefix(true), expected prefix to be true, got false`, vsn)
	}
}

func TestVersion_BumpWithSpec(t *testing.T) {
	tests := []bumpWithSpecTest{
		// REGULAR cases
		{"1.0.0", NONE, "1.0.0"},
		{"1.0.0", MAJOR, "2.0.0"},
		{"1.0.0", MINOR, "1.1.0"},
		{"1.0.0", PATCH, "1.0.1"},
		// PRE-MAJOR cases
		{"0.1.2", MAJOR, "0.2.0"},
		{"0.2.3", MINOR, "0.2.4"},
		{"0.3.4", PATCH, "0.3.5"},
		// PRE-RELEASE cases
		{"1.0.0-alpha", MAJOR, "1.0.0-alpha.1"},
		{"1.1.0-alpha", MAJOR, "2.0.0-alpha"},
		{"1.0.1-alpha", MAJOR, "2.0.0-alpha"},
		{"1.0.0-alpha", MINOR, "1.0.0-alpha.1"},
		{"1.1.0-alpha", MINOR, "1.1.0-alpha.1"},
		{"1.0.1-alpha", MINOR, "1.1.0-alpha"},
		{"1.0.0-alpha", PATCH, "1.0.0-alpha.1"},
		{"1.1.0-alpha", PATCH, "1.1.0-alpha.1"},
		{"1.0.1-alpha", PATCH, "1.0.1-alpha.1"},
		// PRE-MAJOR PRE-RELEASE cases
		{"0.1.0-beta", MAJOR, "1.0.0-beta"},
		{"0.2.3-beta", MINOR, "0.3.0-beta"},
		{"0.3.4-beta", PATCH, "0.3.4-beta.1"},
	}

	var test bumpWithSpecTest
	for _, test = range tests {
		var start, vsn, expect *Version
		start, _ = NewVersion(test.version)
		vsn, _ = NewVersion(test.version)
		expect, _ = NewVersion(test.expect)

		if !vsn.BumpWithSpec(test.bump).Equal(*expect.version) {
			t.Fatalf(`(*Version(%v)).BumpWithSpec("%s"), expected to equal %v, got %v`, start, test.bump, expect, vsn)
		}
	}
}

func TestVersion_BumpWithSpec_BumpError(t *testing.T) {
	vsn := NewEmptyVersion()
	bump := "NOK"

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf(`(*Version(%v)).BumpWithSpec("%s"), expected to panic, got <nil>`, vsn, bump)
		}
	}()

	vsn.BumpWithSpec(bump)
}

func TestVersion_MarshalJSON(t *testing.T) {
	vsn := NewEmptyVersion().WithPrefix(true).BumpWithSpec(MAJOR)
	bytes, err := vsn.MarshalJSON()
	switch true {
	case err != nil:
		t.Fatalf(`(*Version(%v)).MarshalJSON(), expected error to be <nil>, got %v`, vsn, err)
	case string(bytes) != `"v0.1.0"`:
		t.Fatalf(`(*Version(%v)).MarshalJSON(), expected "v0.1.0", got %s`, vsn, string(bytes))
	}
}

func TestVersion_UnmarshalJSON(t *testing.T) {
	vsn := &Version{}
	err := vsn.UnmarshalJSON([]byte(`"v1.0.0-alpha.1"`))
	switch true {
	case err != nil:
		t.Fatalf(`(*Version(%v)).MarshalJSON("\"v1.0.0-alpha.1\""), expected error to be <nil>, got %v`, vsn, err)
	case vsn.String() != "v1.0.0-alpha.1":
		t.Fatalf(`(*Version(%v)).MarshalJSON("\"v1.0.0-alpha.1\""), got %v`, vsn, vsn)
	}

	err = vsn.UnmarshalJSON([]byte("true"))
	if err == nil {
		t.Fatalf(`(*Version(%v)).UnmarshalJSON("true"), expected error NOT to be <nil>`, vsn)
	}

	err = vsn.UnmarshalJSON([]byte(`"foo"`))
	if err == nil {
		t.Fatalf(`(*Version(%v)).UnmarshalJSON("\"foo\""), expected error NOT to be <nil>`, vsn)
	}
}

func TestVersion_String(t *testing.T) {
	var vsn *Version
	expect := "1.22.3-foo.13+bar"
	vsn, _ = NewVersion(expect)
	str := vsn.String()
	if str != expect {
		t.Fatalf(`(*Version(%v)).String(), expected to equal "%s", got "%s"`, vsn, expect, str)
	}

	expect = "v" + expect
	vsn.prefix = true
	str = vsn.String()
	if str != expect {
		t.Fatalf(`(*Version(%v)).String(), expected to equal "%s", got "%s"`, vsn, expect, str)
	}
}

func TestPreRelease_String(t *testing.T) {
	pre := &PreRelease{}
	str := pre.String()
	expect := ""
	if str != expect {
		t.Fatalf(`(*PreRelease(%v)).String(), expected to equal "%s", got "%s"`, pre, expect, str)
	}

	pre = &PreRelease{Tag: "foo"}
	str = pre.String()
	expect = "foo"
	if str != expect {
		t.Fatalf(`(*PreRelease(%v)).String(), expected to equal "%s", got "%s"`, pre, expect, str)
	}

	pre = &PreRelease{Number: 99}
	str = pre.String()
	expect = "99"
	if str != expect {
		t.Fatalf(`(*PreRelease(%v)).String(), expected to equal "%s", got "%s"`, pre, expect, str)
	}

	pre = &PreRelease{Tag: "foo", Number: 99}
	str = pre.String()
	expect = "foo.99"
	if str != expect {
		t.Fatalf(`(*PreRelease(%v)).String(), expected to equal "%s", got "%s"`, pre, expect, str)
	}
}
