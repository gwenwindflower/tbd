with

source as (

    select * from {{ ref('PRODUCTS') }}

),

renamed as (

    select
        PRODUCT_NAME as PRODUCT_NAME,
        PRODUCT_DESCRIPTION as PRODUCT_DESCRIPTION,
        PRODUCT_TYPE as PRODUCT_TYPE,
        PRODUCT_PRICE as PRODUCT_PRICE,
        PRODUCT_ID as PRODUCT_ID,
        IS_FOOD_ITEM as IS_FOOD_ITEM,
        IS_DRINK_ITEM as IS_DRINK_ITEM,
        

    from source

)

select * from renamed
