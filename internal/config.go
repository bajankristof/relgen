package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bajankristof/relgen/internal/conventionalcommits"
	"github.com/bajankristof/relgen/internal/semver"
	"os"
	"regexp"
)

var DefaultChangeSpec = []ChangeSpec{
	{&TypeSpec{regexp.MustCompile("^feat$")}, semver.MINOR, "Features"},
	{&TypeSpec{regexp.MustCompile("^fix$")}, semver.PATCH, "Fixes"},
	{&TypeSpec{regexp.MustCompile("^build|chore|ci|docs|style|refactor|perf|test$")}, "NONE", "Other"},
}

type Config struct {
	PreRelease    string            `json:"preRelease"`
	BuildMetadata string            `json:"buildMetadata"`
	VersionPrefix bool              `json:"versionPrefix"`
	ChangeSpec    []ChangeSpec      `json:"changeSpec"`
	Outputs       OutputWriterGroup `json:"outputs"`
}

type ChangeSpec struct {
	Type     *TypeSpec `json:"type"`
	Bump     string    `json:"bump"`
	Category string    `json:"category"`
}

type TypeSpec struct {
	*regexp.Regexp
}

func ReadConfig(path string) (*Config, error) {
	_, err := os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		cfg := &Config{ChangeSpec: DefaultChangeSpec, Outputs: DefaultOutputGroup}
		return cfg, nil
	}

	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg *Config
	err = json.Unmarshal(bytes, &cfg)
	if err != nil {
		return nil, err
	}

	err = cfg.Check()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (cfg *Config) Check() error {
	if len(cfg.Outputs) < 1 {
		cfg.Outputs = DefaultOutputGroup
	}

	if len(cfg.ChangeSpec) < 1 {
		cfg.ChangeSpec = DefaultChangeSpec
		return nil
	}

	for _, change := range cfg.ChangeSpec {
		if err := change.Check(); err != nil {
			return err
		}
	}

	return nil
}

func (cfg *Config) FindChangeSpec(cc *conventionalcommits.ConventionalCommit) (int, *ChangeSpec) {
	for i, spec := range cfg.ChangeSpec {
		if spec.Type.MatchString(cc.Type) {
			return i, &spec
		}
	}

	return 0, nil
}

func (spec *ChangeSpec) Check() error {
	switch spec.Bump {
	case
		semver.NONE,
		semver.MAJOR,
		semver.MINOR,
		semver.PATCH:
		return nil
	default:
		return fmt.Errorf("unrecognized bump spec \"%s\"", spec.Bump)
	}
}

func (spec *TypeSpec) MarshalJSON() ([]byte, error) {
	str := spec.String()
	return json.Marshal(str[4:])
}

func (spec *TypeSpec) UnmarshalJSON(bytes []byte) error {
	str := ""
	err := json.Unmarshal(bytes, &str)
	if err != nil {
		return err
	}

	regex, err := regexp.Compile("(?i)" + str)
	if err != nil {
		return err
	}

	spec.Regexp = regex
	return nil
}
