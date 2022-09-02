package internal

import (
	"encoding/json"
	"golang.org/x/sync/errgroup"
	"os"
	"text/template"
)

var DefaultVersionOutput = &OutputWriter{
	Path:     "version.txt",
	Template: template.Must(template.New("version.txt").Parse(`{{.Version | print}}`)),
}

var DefaultChangelogOutput = &OutputWriter{
	Path: "changelog-entry.md",
	Template: template.Must(template.New("changelog-entry.md").Parse(`## {{.Version | print}} ({{.Date.Format "2006-01-02"}}){{range $category, $changes := .Changelog}}
### {{$category}}{{range $cc := $changes}}
* {{$cc.Description}} (#{{printf "%.*s" 8 $cc.Hash}}){{end}}
{{end}}`)),
}

var DefaultOutputGroup = OutputWriterGroup{
	DefaultVersionOutput,
	DefaultChangelogOutput,
}

type OutputWriterGroup []*OutputWriter

type OutputWriter struct {
	Path     string
	Template *template.Template
}

func (group OutputWriterGroup) Execute(rel *Release) error {
	errs := &errgroup.Group{}

	for i := range group {
		writer := group[i]
		errs.Go(func() error {
			return writer.Execute(rel)
		})
	}

	return errs.Wait()
}

func (writer *OutputWriter) Execute(rel *Release) error {
	file, err := os.Create(writer.Path)
	if err != nil {
		return err
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	return writer.Template.Execute(file, rel)
}

func (writer *OutputWriter) UnmarshalJSON(data []byte) error {
	tmp := &struct {
		Path     string `json:"path"`
		Type     string `json:"type"`
		Template string `json:"template"`
	}{}

	err := json.Unmarshal(data, tmp)
	if err != nil {
		return err
	}

	writer.Path = tmp.Path
	switch tmp.Type {
	case "version.txt":
		writer.Template = DefaultVersionOutput.Template
	case "changelog-entry.md":
		writer.Template = DefaultChangelogOutput.Template
	default:
		writer.Template, err = template.ParseFiles(tmp.Template)
	}

	return err
}
