with

source as (

    select * from {{ ref('LOCATIONS') }}

),

renamed as (

    select
        TAX_RATE as TAX_RATE,
        LOCATION_ID as LOCATION_ID,
        OPENED_DATE as OPENED_DATE,
        LOCATION_NAME as LOCATION_NAME,
        

    from source

)

select * from renamed
