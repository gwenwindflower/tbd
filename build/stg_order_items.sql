with

source as (

    select * from {{ ref('ORDER_ITEMS') }}

),

renamed as (

    select
        IS_DRINK_ITEM as IS_DRINK_ITEM,
        PRODUCT_ID as PRODUCT_ID,
        PRODUCT_PRICE as PRODUCT_PRICE,
        ORDER_ID as ORDER_ID,
        ORDER_ITEM_ID as ORDER_ITEM_ID,
        ORDERED_AT as ORDERED_AT,
        PRODUCT_NAME as PRODUCT_NAME,
        IS_FOOD_ITEM as IS_FOOD_ITEM,
        SUPPLY_COST as SUPPLY_COST,
        

    from source

)

select * from renamed
