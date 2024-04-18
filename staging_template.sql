with

source as (

    select * from {{ "{{" }} source('{{.Schema}}', '{{.Name}}') {{ "}}" }}

),

renamed as (

    select
{{- range $group, $columns := .DataTypeGroups -}}
    {{- if eq $group "text" }}
        -- text
        {{ range $columns -}}
            {{- .Name }} as {{ .Name | lower }},
        {{ end -}}
    {{- end -}}
    {{- if eq $group "numbers" }}
        -- numbers
        {{ range $columns -}}
            {{- .Name }} as {{ .Name | lower }},
        {{ end -}}
    {{- end -}}
    {{- if eq $group "booleans" }}
        -- booleans
        {{ range $columns -}}
            {{- .Name }} as {{ .Name | lower }},
        {{ end -}}
    {{- end -}}
    {{- if eq $group "datetimes" }}
        -- datetimes
        {{ range $columns -}}
            {{- .Name }} as {{ .Name | lower }},
        {{ end -}}
    {{- end -}}
    {{- if eq $group "timestamps" }}
        -- timestamps
        {{ range $columns -}}
            {{- .Name }} as {{ .Name | lower }},
        {{ end -}}
    {{- end -}}
{{ end }}
    from source
)

select * from renamed
