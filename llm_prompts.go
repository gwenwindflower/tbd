package main

const (
	DESC_PROMPT = `Generate a description for a column in a specific table in a data warehouse,
  the table is called %s and the column is called %s. The description should be concise, 1 to 3 sentences,
  and inform both business users and technical data analyts about the purpose and contents of the column.
  Avoid using the column name in the description, as it is redundant — put another way do not use tautological
  descriptions, for example on an 'order_id' column saying "This is the id of an order". Don't do that. A good
  example for an 'order_id' column would be something like "This is the primary key of the orders table,
  each distinct order has a unique 'order_id'". Another good example for an orders table would be describing
  'product_type' as "The category of product, the bucket that a product falls into, for example 'electronics' or 'clothing'".
  Avoid making assumptions about the data, as you don't have access to it. Don't make assertions about data that you 
  haven't seen, just use business context, the table name, and the column to generate the description. The description.
  There is no need to add a title just the sentences that compose the description, it's being put onto a field in a YAML file, 
so again, no title, no formatting, just 1 to 3 sentences.`
	TESTS_PROMPT = `Generate a list of tests that can be run on a column in a specific table in a data warehouse,
the table is called %s and the column is called %s. The tests are YAML config, there are 2 to choose from.
They have the following structure, follow this structure exactly:
  - unique
  - not_null
Return only the tests that are applicable to the column, for example, a column that is a primary key should have 
both unique and not_null tests, while a column that is a foreign key should only have the not_null test. If a 
column is potentially optional, then it should have neither test. Return only the tests that are applicable to the column.
They will be nested under a 'tests' key in a YAML file, so no need to add a title or key, just the list of tests by themselves.
  For example, a good response for a 'product_type' column in an 'orders' table would be:
  - not_null

  A good response for an 'order_id' column in an 'orders' table would be:
  - unique
  - not_null

  A good response for a 'product_sku' column in an 'orders' table would be:
  - not_null
`
)
