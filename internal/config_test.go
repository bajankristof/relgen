package internal

import (
	"github.com/bajankristof/relgen/internal/semver"
	"os"
	"path"
	"regexp"
	"testing"
)

func TestReadConfig(t *testing.T) {
	p := path.Join(t.TempDir(), t.Name()+".json")
	err := os.WriteFile(p, []byte(`{"versionPrefix":true}`), 0777)
	if err != nil {
		panic(err)
	}

	cfg, err := ReadConfig(p)
	switch true {
	case err != nil:
		t.Fatalf(`ReadConfig("%s") = (%v, %v), expected error to be <nil>, got %v`, p, cfg, err, err)
	case !cfg.VersionPrefix:
		t.Fatalf(`ReadConfig("%s") = (%v, %v), expected version prefix to be true, got false`, p, cfg, err)
	}
}

func TestReadConfig_ReadError(t *testing.T) {
	p := path.Join(t.TempDir(), t.Name()+".json")
	err := os.WriteFile(p, []byte(`{"changeSpec":[{"bump":"NOK"}]}`), 0000)
	if err != nil {
		panic(err)
	}

	cfg, err := ReadConfig(p)
	switch true {
	case err == nil:
		t.Fatalf(`ReadConfig("%s") = (%v, %v), expected error NOT to be <nil>`, p, cfg, err)
	case cfg != nil:
		t.Fatalf(`ReadConfig("%s") = (%v, %v), expected config to be <nil>, got %v`, p, cfg, err, cfg)
	}
}

func TestReadConfig_JSONError(t *testing.T) {
	p := path.Join(t.TempDir(), t.Name()+".json")
	err := os.WriteFile(p, []byte(`{"versionPrefix":true`), 0777)
	if err != nil {
		panic(err)
	}

	cfg, err := ReadConfig(p)
	switch true {
	case err == nil:
		t.Fatalf(`ReadConfig("%s") = (%v, %v), expected error NOT to be <nil>`, p, cfg, err)
	case cfg != nil:
		t.Fatalf(`ReadConfig("%s") = (%v, %v), expected config to be <nil>, got %v`, p, cfg, err, cfg)
	}
}

func TestReadConfig_CheckError(t *testing.T) {
	p := path.Join(t.TempDir(), t.Name()+".json")
	err := os.WriteFile(p, []byte(`{"changeSpec":[{"bump":"NOK"}]}`), 0777)
	if err != nil {
		panic(err)
	}

	cfg, err := ReadConfig(p)
	switch true {
	case err == nil:
		t.Fatalf(`ReadConfig("%s") = (%v, %v), expected error NOT to be <nil>`, p, cfg, err)
	case cfg != nil:
		t.Fatalf(`ReadConfig("%s") = (%v, %v), expected config to be <nil>, got %v`, p, cfg, err, cfg)
	}
}

func TestReadConfig_ErrNotExists(t *testing.T) {
	p := path.Join(t.TempDir(), t.Name()+".json")
	cfg, err := ReadConfig(p)
	switch true {
	case err != nil:
		t.Fatalf(`ReadConfig("%s") = (%v, %v), expected error to be <nil>, got %v`, p, cfg, err, err)
	case len(cfg.ChangeSpec) != len(DefaultChangeSpec):
		t.Fatalf(`ReadConfig("%s) = (%v, %v), expected change spec to be the default change spec`, p, cfg, err)
	case len(cfg.Outputs) != len(DefaultOutputGroup):
		t.Fatalf(`ReadConfig("%s) = (%v, %v), expected outputs to be the default output group`, p, cfg, err)
	}
}

func TestConfig_Check(t *testing.T) {
	cfg := &Config{ChangeSpec: []ChangeSpec{{Bump: semver.MINOR}}}
	err := cfg.Check()
	switch true {
	case err != nil:
		t.Fatalf(`(*Config(%v)).Check(), expected error to be <nil>, got %v`, cfg, err)
	case len(cfg.ChangeSpec) != 1:
		t.Fatalf(`(*Config(%v)).Check(), expected change spec to be unchanged`, cfg)
	}
}

func TestConfig_CheckDefault(t *testing.T) {
	cfg := &Config{}
	err := cfg.Check()
	switch true {
	case err != nil:
		t.Fatalf(`(*Config(%v)).Check(), expected error to be <nil>, got %v`, cfg, err)
	case len(cfg.ChangeSpec) != len(DefaultChangeSpec):
		t.Fatalf(`(*Config(%v)).Check(), expected change spec to be the default change spec, got %v`, cfg, cfg.ChangeSpec)
	}
}

func TestConfig_CheckError(t *testing.T) {
	cfg := &Config{ChangeSpec: []ChangeSpec{{Bump: semver.MAJOR}, {Bump: "NASTY"}}}
	err := cfg.Check()
	if err == nil {
		t.Fatalf(`(*Config(%v)).Check(), expected error NOT to be <nil>`, cfg)
	}
}

func TestChangeSpec_Check(t *testing.T) {
	spec := &ChangeSpec{Bump: semver.NONE}
	if err := spec.Check(); err != nil {
		t.Fatalf(`(*ChangeSpec(%v)).Check(), expected error to be <nil>, got %v`, spec, err)
	}

	spec = &ChangeSpec{Bump: "NARLY"}
	if err := spec.Check(); err == nil {
		t.Fatalf(`(*ChangeSpec(%v)).Check(), expected error NOT to be <nil>`, spec)
	}
}

func TestTypeSpec_MarshalJSON(t *testing.T) {
	spec := &TypeSpec{regexp.MustCompile("(?i)^ok$")}
	bytes, err := spec.MarshalJSON()
	switch true {
	case err != nil:
		t.Fatalf(`(*TypeSpec(%v)).MarshalJSON(), expected error to be <nil>, got %v`, spec, err)
	case string(bytes) != `"^ok$"`:
		t.Fatalf(`(*TypeSpec(%v)).MarshalJSON(), expected %s, got %s`, spec, `"^ok$"`, string(bytes))
	}
}

func TestTypeSpec_UnmarshalJSON(t *testing.T) {
	spec := &TypeSpec{}
	err := spec.UnmarshalJSON([]byte(`"^ok$"`))
	switch true {
	case err != nil:
		t.Fatalf(`(*TypeSpec(%v)).UnmarshalJSON("\"^ok$\""), expected error to be <nil>, got %v`, spec, err)
	case spec.Regexp.String() != "(?i)^ok$":
		t.Fatalf(`(*TypeSpec(%v)).UnmarshalJSON("\"^ok$\""), expected regexp to be "(?i)^ok$", got "%s"`, spec, spec.Regexp.String())
	}

	err = spec.UnmarshalJSON([]byte("true"))
	if err == nil {
		t.Fatalf(`(*TypeSpec(%v)).UnmarshalJSON("true"), expected error NOT to be <nil>`, spec)
	}

	err = spec.UnmarshalJSON([]byte(`"["`))
	if err == nil {
		t.Fatalf(`(*TypeSpec(%v)).UnmarshalJSON("\"[\""), expected error NOT to be <nil>`, spec)
	}
}
