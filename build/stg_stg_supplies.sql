with

source as (

    select * from {{ ref('STG_SUPPLIES') }}

),

renamed as (

    select
        SUPPLY_NAME as SUPPLY_NAME,
        SUPPLY_UUID as SUPPLY_UUID,
        IS_PERISHABLE_SUPPLY as IS_PERISHABLE_SUPPLY,
        SUPPLY_COST as SUPPLY_COST,
        SUPPLY_ID as SUPPLY_ID,
        PRODUCT_ID as PRODUCT_ID,
        

    from source

)

select * from renamed
