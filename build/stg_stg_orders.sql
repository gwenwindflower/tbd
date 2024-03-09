with

source as (

    select * from {{ ref('STG_ORDERS') }}

),

renamed as (

    select
        LOCATION_ID as LOCATION_ID,
        ORDER_TOTAL as ORDER_TOTAL,
        TAX_PAID as TAX_PAID,
        CUSTOMER_ID as CUSTOMER_ID,
        ORDERED_AT as ORDERED_AT,
        ORDER_ID as ORDER_ID,
        

    from source

)

select * from renamed
