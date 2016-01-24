DbService
=========

Intension of `dbservice` is to be able create simple RESTful services that respond with json and use `PostgreSQL` as data storage. Here is how `dbservice` works:

Routes
------

First you need to define routes. Here is an example of routes file:

```
get /products, name: 'get_products', collection: true
get /product/:id, name: 'get_product'
post /product, name: 'create_product'
put /product, name: 'update_product'
```

Here you define you web service routes by specifying request method, request url and route name. There are couple optional parameters: `collection` and `custom`. Set `collection: true` if your endpoint returns array of values and `custom: true` if your sql request will already return formatted json data (otherwise custom sql will be added to make sure that PostgreSQL returns query results in json format).

Parameteters validation
-----------------------

Every request made to dbservice has parameters that can come from multiple sources. Parameters can come as part of url (`:<parameter_name>` section in url), they can come from query string (`?parameter_name=parameter_value` in url) or they can come as part of body. You can submit form (it has some limitations as you are only able to set key- value parameters) or use `application/json` Content-Type to supply arbitrary data structures.

After all parameters are merge, they are validated with json schema (if json schema file is present for a particular route). It should be located in `schemas/<route_name>.schema` file. If there were validation errors during parameters validation, json with field names as keys and error messages as values will be returned back (status code will be 400). If schema is missing for a particular route, no validation of parameters will occur.
