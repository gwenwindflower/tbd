with

source as (

    select * from {{ ref('SUPPLIES') }}

),

renamed as (

    select
        PRODUCT_ID as PRODUCT_ID,
        SUPPLY_COST as SUPPLY_COST,
        SUPPLY_UUID as SUPPLY_UUID,
        SUPPLY_NAME as SUPPLY_NAME,
        IS_PERISHABLE_SUPPLY as IS_PERISHABLE_SUPPLY,
        SUPPLY_ID as SUPPLY_ID,
        

    from source

)

select * from renamed
