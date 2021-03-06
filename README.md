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
select * from products where id={{.params.id}}
```

In this example id is expected to be integer (make sure that that's the case using schema validation). If you want to insert string into sql statement, it has to be quoted:

```
select * from products where status={{.params.status | quote}}
```

If you are performing insert/update/delete operation, don't forget `returning` statement to get response data. Example:

```
insert into products (name, price) values ({{.params.name | quote}}, {{.params.price}}) returning *
```

or

```
insert into products (name, price) values ({{.params.name | quote}}, {{.params.price}}) returning true as success
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

Plugins
=======

Plugins can register and then add hook before request execution (possibly setting some data that will be accessible in sql templates). That is optional. Also plugins can modify response json as well as to add additional response headers/cookies.

JWT plugin
==========

Configuration
-------------

Create `plugins/jwt.toml` configuration file. Example:

```
secret = "secret123"
issuer = "issuer"
expiration = "4h"
rotation_deadline = "2h"
```

This will set secret for encoding token, issuer claim (optional), expiration time for token and rotation_deadline (optional). Rotation deadline states how close to expiration should it be for server to transparantly issue new token with updated expiration time.

Routes
------

```
post /login, name: 'login' | jwt
```

Once you put jwt.toml file into plugins folder, jwt plugin will be enabled for dbservice. If you add jwt plugin to one of the actions, then response json of plugin will need to contain `__jwt` key that will contain payload data. Example of login json response:

```
{"__jwt": {"user_id": 5, "admin": true}, "success": true}
```

Based on content of `__jwt` value, jwt token will be created (or not if it's not present). After this `__jwt` key will be removed from response. And in this case user wil get `{"success": true}` as response with additional header that has jwt token.

Reading payload
---------------

If jwt plugin is enabled, you can use content from jwt payload in sql templates by using following syntax:

```
{{.jwt.<key name>}}
```

This will insert value from the payload into sql query.

TODO:
- Email sending plugin
- Browser detection plugin
- Serving static html plugin (likely should include A/B testing ability)
- Validation of files and files upload to s3
- Delayed jobs
- Testing endpoints
- Automatic reation of documentation (possibly client libraries in future)
- Code generators
- Ability to hold migration files and run migrations up/down
- Plugin for region (country) detection (possibly setting up redirect or serve different content)
