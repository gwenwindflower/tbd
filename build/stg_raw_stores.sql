with

source as (

    select * from {{ ref('RAW_STORES') }}

),

renamed as (

    select
        TAX_RATE as TAX_RATE,
        NAME as NAME,
        ID as ID,
        OPENED_AT as OPENED_AT,
        

    from source

)

select * from renamed
