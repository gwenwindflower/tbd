with

source as (

    select * from {{ ref('METRICFLOW_TIME_SPINE') }}

),

renamed as (

    select
        DATE_DAY as DATE_DAY,
        

    from source

)

select * from renamed
