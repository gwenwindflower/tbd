with

source as (

    select * from {{ ref('STG_PRODUCTS') }}

),

renamed as (

    select
        PRODUCT_DESCRIPTION as PRODUCT_DESCRIPTION,
        PRODUCT_ID as PRODUCT_ID,
        PRODUCT_TYPE as PRODUCT_TYPE,
        IS_DRINK_ITEM as IS_DRINK_ITEM,
        PRODUCT_PRICE as PRODUCT_PRICE,
        IS_FOOD_ITEM as IS_FOOD_ITEM,
        PRODUCT_NAME as PRODUCT_NAME,
        

    from source

)

select * from renamed
