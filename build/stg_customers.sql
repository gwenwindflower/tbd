with

source as (

    select * from {{ ref('CUSTOMERS') }}

),

renamed as (

    select
        CUSTOMER_ID as CUSTOMER_ID,
        LAST_ORDERED_AT as LAST_ORDERED_AT,
        CUSTOMER_NAME as CUSTOMER_NAME,
        LIFETIME_SPEND as LIFETIME_SPEND,
        CUSTOMER_TYPE as CUSTOMER_TYPE,
        COUNT_LIFETIME_ORDERS as COUNT_LIFETIME_ORDERS,
        LIFETIME_SPEND_PRETAX as LIFETIME_SPEND_PRETAX,
        FIRST_ORDERED_AT as FIRST_ORDERED_AT,
        

    from source

)

select * from renamed
