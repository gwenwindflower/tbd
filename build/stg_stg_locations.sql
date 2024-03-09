with

source as (

    select * from {{ ref('STG_LOCATIONS') }}

),

renamed as (

    select
        OPENED_DATE as OPENED_DATE,
        TAX_RATE as TAX_RATE,
        LOCATION_NAME as LOCATION_NAME,
        LOCATION_ID as LOCATION_ID,
        

    from source

)

select * from renamed
