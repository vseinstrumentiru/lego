{{ if .Versions -}}
<a name="unreleased"></a>
## [Unreleased]

{{- if .Unreleased.CommitGroups }}
{{ $prevTitle := "" -}}
{{ range .Unreleased.CommitGroups -}}
{{ if not (eq $prevTitle .Title) -}}
{{ $prevTitle = .Title -}}
### {{ .Title }}
{{ end -}}
{{ range .Commits -}}
- [{{ .Hash.Short }}] {{ if .Scope }}**{{ .Scope }}:** {{ end }}{{ .Subject }}
{{ end -}}
{{ end -}}
{{ end -}}
{{ end -}}

{{- range .Versions }}
<a name="{{ .Tag.Name }}"></a>
## {{ if .Tag.Previous }}[{{ .Tag.Name }}]{{ else }}{{ .Tag.Name }}{{ end }} - {{ datetime "2006-01-02" .Tag.Date }}
{{ $prevTitle := "" -}}
{{ range .CommitGroups -}}
{{ if not (eq $prevTitle .Title) -}}
{{ $prevTitle = .Title -}}
### {{ .Title }}
{{ end -}}
{{ range .Commits -}}
- [{{ .Hash.Short }}] {{ if .Scope }}**{{ .Scope }}:** {{ end }}{{ .Subject }}
{{ end -}}
{{ end -}}

{{- if .RevertCommits -}}
### Reverts
{{ range .RevertCommits -}}
- {{ .Revert.Header }}
{{ end -}}
{{ end -}}

{{- if .MergeCommits -}}
### Pull Requests
{{ range .MergeCommits -}}
- {{ .Header }}
{{ end -}}
{{ end -}}

{{- if .NoteGroups -}}
{{ range .NoteGroups -}}
### {{ .Title }}
{{ range .Notes }}
{{ .Body }}
{{ end -}}
{{ end -}}
{{ end -}}
{{ end -}}

{{- if .Versions }}
[Unreleased]: {{ .Info.RepositoryURL }}/compare/{{ $latest := index .Versions 0 }}{{ $latest.Tag.Name }}...HEAD
{{ if .Unreleased.CommitGroups -}}
{{ range .Unreleased.CommitGroups -}}
{{ range .Commits -}}
[{{ .Hash.Short }}]: {{ $.Info.RepositoryURL }}/commit/{{ .Hash.Long }}
{{ end -}}
{{ end -}}
{{ end -}}
{{ range .Versions -}}
{{ if .Tag.Previous -}}
[{{ .Tag.Name }}]: {{ $.Info.RepositoryURL }}/compare/{{ .Tag.Previous.Name }}...{{ .Tag.Name }}
{{ end -}}
{{ range .CommitGroups -}}
{{ range .Commits -}}
[{{ .Hash.Short }}]: {{ $.Info.RepositoryURL }}/commit/{{ .Hash.Long }}
{{ end -}}
{{ end -}}
{{ end -}}
{{ end -}}
