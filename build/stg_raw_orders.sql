with

source as (

    select * from {{ ref('RAW_ORDERS') }}

),

renamed as (

    select
        ID as ID,
        ORDER_TOTAL as ORDER_TOTAL,
        SUBTOTAL as SUBTOTAL,
        TAX_PAID as TAX_PAID,
        STORE_ID as STORE_ID,
        ORDERED_AT as ORDERED_AT,
        CUSTOMER as CUSTOMER,
        

    from source

)

select * from renamed
