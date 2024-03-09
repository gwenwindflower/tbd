with

source as (

    select * from {{ ref('STG_ORDER_ITEMS') }}

),

renamed as (

    select
        ORDER_ITEM_ID as ORDER_ITEM_ID,
        PRODUCT_ID as PRODUCT_ID,
        ORDER_ID as ORDER_ID,
        

    from source

)

select * from renamed
