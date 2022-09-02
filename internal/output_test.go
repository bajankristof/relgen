package internal

import (
	"fmt"
	"github.com/bajankristof/relgen/internal/semver"
	"os"
	"path"
	"testing"
	"text/template"
)

func TestOutputWriter_Execute(t *testing.T) {
	writer := &OutputWriter{
		Path:     path.Join(t.TempDir(), "foo/bar/test.txt"),
		Template: template.Must(template.New("test.txt").Parse(`Version: {{.Version | print}}`)),
	}

	rel := &Release{Version: semver.NewEmptyVersion()}
	err := writer.Execute(rel)

	if err != nil {
		t.Fatalf("(*OutputWriter(%v)).Execute(%v) = %v, expected error to be <nil>, got %v", writer, rel, err, err)
	}

	data, _ := os.ReadFile(writer.Path)
	got := string(data)
	expect := "Version: 0.0.0"
	if got != expect {
		t.Fatalf(`(*OutputWriter(%v)).Execute(%v) = %v, expected to write "%s", got "%s"`, writer, rel, err, expect, got)
	}
}

func TestOutputWriter_UnmarshalJSON(t *testing.T) {
	tpl := path.Join(t.TempDir(), "test.tpl")
	txt := path.Join(t.TempDir(), "test.txt")
	err := os.WriteFile(tpl, []byte(`Version: {{.Version}}`), 0777)
	if err != nil {
		panic(err)
	}

	writer := &OutputWriter{}
	data := fmt.Sprintf(`{"path":"%s","template":"%s"}`, txt, tpl)
	err = writer.UnmarshalJSON([]byte(data))
	switch true {
	case err != nil:
		t.Fatalf(`(*OutputWriter(%v)).Execute(%v), expected error to be <nil>, got %v`, writer, data, err)
	case writer.Path != txt:
		t.Fatalf(`(*OutputWriter(%v)).Execute(%v), expected path to be "%s", got "%s"`, writer, data, txt, writer.Path)
	case writer.Template == nil:
		t.Fatalf(`(*OutputWriter(%v)).Execute(%v), expected template NOT to be <nil>"`, writer, data)
	}
}

func TestOutputWriter_UnmarshalJSON_WithType(t *testing.T) {
	p := path.Join(t.TempDir(), "test.txt")
	writer := &OutputWriter{}
	data := fmt.Sprintf(`{"path":"%s","type":"version.txt"}`, p)
	err := writer.UnmarshalJSON([]byte(data))
	switch true {
	case err != nil:
		t.Fatalf(`(*OutputWriter(%v)).Execute(%v), expected error to be <nil>, got %v`, writer, data, err)
	case writer.Path != p:
		t.Fatalf(`(*OutputWriter(%v)).Execute(%v), expected path to be "%s", got "%s"`, writer, data, p, writer.Path)
	case writer.Template != DefaultVersionOutput.Template:
		expect := DefaultVersionOutput.Template
		got := writer.Template
		t.Fatalf(`(*OutputWriter(%v)).Execute(%v), expected template to be %v, got %v"`, writer, data, expect, got)
	}

	p = path.Join(t.TempDir(), "test.md")
	writer = &OutputWriter{}
	data = fmt.Sprintf(`{"path":"%s","type":"changelog-entry.md"}`, p)
	err = writer.UnmarshalJSON([]byte(data))
	switch true {
	case err != nil:
		t.Fatalf(`(*OutputWriter(%v)).Execute(%v), expected error to be <nil>, got %v`, writer, data, err)
	case writer.Path != p:
		t.Fatalf(`(*OutputWriter(%v)).Execute(%v), expected path to be "%s", got "%s"`, writer, data, p, writer.Path)
	case writer.Template != DefaultChangelogOutput.Template:
		expect := DefaultChangelogOutput.Template
		got := writer.Template
		t.Fatalf(`(*OutputWriter(%v)).Execute(%v), expected template to be %v, got %v"`, writer, data, expect, got)
	}
}

func TestOutputWriterGroup_Execute(t *testing.T) {
	group := OutputWriterGroup{
		&OutputWriter{
			Path:     path.Join(t.TempDir(), "foo.txt"),
			Template: template.Must(template.New("foo.txt").Parse("foo")),
		},
		&OutputWriter{
			Path:     path.Join(t.TempDir(), "bar.txt"),
			Template: template.Must(template.New("bar.txt").Parse("bar")),
		},
	}

	rel := &Release{Version: semver.NewEmptyVersion()}
	err := group.Execute(rel)
	if err != nil {
		t.Fatalf(`(*OutputWriterGroup(%v)).Execute(%v), expected error to be <nil>, got %v`, group, rel, err)
	}
}
