with

source as (

    select * from {{ ref('ORDERS') }}

),

renamed as (

    select
        ORDER_ID as ORDER_ID,
        CUSTOMER_ID as CUSTOMER_ID,
        ORDER_COST as ORDER_COST,
        IS_DRINK_ORDER as IS_DRINK_ORDER,
        ORDERED_AT as ORDERED_AT,
        ORDER_TOTAL as ORDER_TOTAL,
        LOCATION_ID as LOCATION_ID,
        IS_FOOD_ORDER as IS_FOOD_ORDER,
        COUNT_FOOD_ITEMS as COUNT_FOOD_ITEMS,
        COUNT_DRINK_ITEMS as COUNT_DRINK_ITEMS,
        TAX_PAID as TAX_PAID,
        

    from source

)

select * from renamed
