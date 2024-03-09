with

source as (

    select * from {{ ref('RAW_PRODUCTS') }}

),

renamed as (

    select
        PRICE as PRICE,
        SKU as SKU,
        DESCRIPTION as DESCRIPTION,
        TYPE as TYPE,
        NAME as NAME,
        

    from source

)

select * from renamed
