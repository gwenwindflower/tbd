with

source as (

    select * from {{ "{{" }} ref('{{.Name}}') {{ "}}" }}

),

renamed as (

    select
        {{ range .Columns -}}
            {{ .Name }} as {{ .Name }},
        {{ end }}

    from source

)

select * from renamed
