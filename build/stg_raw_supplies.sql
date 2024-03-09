with

source as (

    select * from {{ ref('RAW_SUPPLIES') }}

),

renamed as (

    select
        SKU as SKU,
        PERISHABLE as PERISHABLE,
        COST as COST,
        NAME as NAME,
        ID as ID,
        

    from source

)

select * from renamed
