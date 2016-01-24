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

Sql generation
--------------

After request parameters are validated, they are used in sql template to generate sql request. Example of sql template:

```
select * from products where id={{.id}}
```

In this example id is expected to be integer (make sure that that's the case using schema validation). If you want to insert string into sql statement, it has to be quoted:

```
select * from products where status={{.status | quote}}
```

If you are performing insert/update/delete operation, don't forget `returning` statement to get response data. Example:

```
insert into products (name, price) values ({{.name | quote}}, {{.price}}) returning *
```

or

```
insert into products (name, price) values ({{.name | quote}}, {{.price}}) returning true as success
```

Sql generation is quite powerful and you can use all the power of `text/template` package.

Response
--------

After sql query is executed, resulting data is serialized as json array or object (that depends if route is for collection or not). And then it's returned back to user.

Possible improvements
---------------------

Right now one of the biggest shortcomings of this project is that there is no ability to handle authentication and user sessions. That might come in future if project will turn out to be useful for people.
