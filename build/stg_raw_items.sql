with

source as (

    select * from {{ ref('RAW_ITEMS') }}

),

renamed as (

    select
        SKU as SKU,
        ID as ID,
        ORDER_ID as ORDER_ID,
        

    from source

)

select * from renamed
