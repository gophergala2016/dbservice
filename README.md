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

After request parameters are validated, they are used in sql template to generate sql request. Template should be in `sql/<route_name>.sql` file. Example of sql template:

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

Database connection configuration
---------------------------------

Database connection configuration is in `config.toml` file. Example:

```
user = "postgres"
password = "secret123"
database = "dbservice_example"
host = "127.0.0.1"
port = 5434
sslmode = "disable"
```

User, host, port and sslmode are optional. Defaults are '127.0.0.1', 5432 and 'disbale'.

Response
--------

After sql query is executed, resulting data is serialized as json array or object (that depends if route is for collection or not). And then it's returned back to user.

Versioning
==========

If you want to introduce versioning to your api, you need to can couple of parameters to your routes file. All of them are optional. Example:

```
api_version: 10
deprecated_api_version: [1-2, 4, 6-9]
min_api_version: 1
```

`api_version` - allows to specify current api version

`deprecated_api_version` - here you can specify which api versions will get additional deprecated header in response. `X-Api-Deprecated: true`. You can specify both individual values and ranges.

`min_api_version` - specifies minimal api version. Versions before that will return an error and won't be served.

By default all api versions will be served by the same sql template. But if you want to customize sql template for a particular version, you can do that by adding `.v{version_number}.sql` extension. Let's say you want to customize `get_product.sql` templat for version 2. In this case you would have another file with `get_product.v2.sql` name. And it would only be used for version 2 and all the versions below (unless one of this versions have been customized by another file). Same thing about changing request schema for particular version. Just rename schema file in a similar fasion. `get_product.v2.schema` in our case. By default it will use `get_product.schema` if it exists or will just omit schema check if file has not been supplied.

Possible improvements
---------------------

Right now one of the biggest shortcomings of this project is that there is no ability to handle authentication and user sessions. That might come in future if project will turn out to be useful for people.

Install
=======

```
go get github.com/gophergala2016/dbservice
```

Usage
=====

To start `dbservice` server, just launch `dbservice` executable when you are in your dbservice root folder. Example of files structure:

```
├── config.toml
├── routes
├── schemas
│   ├── create_product.schema
│   ├── get_products_by_status.schema
│   ├── get_product.schema
│   └── update_product.schema
└── sql
    ├── create_product.sql
    ├── get_products_by_status.sql
    ├── get_product.sql
    ├── get_products.sql
    └── update_product.sql
```

By default web server will be running on port `8080`. To customize port, just specify it as command line argument:

```
dbservice 3000
```

You can try `example` project that is located in [example folder](https://github.com/gophergala2016/dbservice/tree/master/example). It has README.
