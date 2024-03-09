with

source as (

    select * from {{ ref('STG_CUSTOMERS') }}

),

renamed as (

    select
        CUSTOMER_ID as CUSTOMER_ID,
        CUSTOMER_NAME as CUSTOMER_NAME,
        

    from source

)

select * from renamed
