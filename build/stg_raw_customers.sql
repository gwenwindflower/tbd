with

source as (

    select * from {{ ref('RAW_CUSTOMERS') }}

),

renamed as (

    select
        ID as ID,
        NAME as NAME,
        

    from source

)

select * from renamed
