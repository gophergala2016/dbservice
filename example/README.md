Example application
===================

This is example `dbservice` application. It is used to showcase functionality that can be provided by dbservice.

Database setup
--------------
Simplest way of setting up databaes is to use Docker:

* Start Postgrs on 5434 (in case you already have one running locally)
```
sudo docker run -e POSTGRES_PASSWORD=secret123 -d --publish=5434:5432 postgres:9.5
```
* If you currently have psql tool (it goes with postgres), get into postgres console:
```
PGPASSWORD=secret123 psql -U postgres -h 127.0.0.1 -p 5434
```

Otherwise,
```
sudo docker exec -it <postgres container id> /bin/bash
psql -U postgres
```

* Create database and table:
```
create database dbservice_example;
\c dbservice_example

create table products(
  id serial,
  name text not null,
  price integer not null,
  status text
);
```
Try different requests
----------------------

Get all the products
```
curl http://localhost:8080/products
```

Create a new product:
```
curl -H "Content-Type: application/json" -X POST -d '{"name":"Brush","price":10}' http://localhost:8080/products
```

Create product with status:
```
curl -H "Content-Type: application/json" -X POST -d '{"name":"Fancy Brush","price":20, "status": "available"}' http://localhost:8080/products
```

Show products with `available` status
```
curl http://localhost:8080/products-with/available
```

Update product name and price: (make sure that you specify existing product id)
```
curl -H "Content-Type: application/json" -X PUT -d '{"name":"Brush++","price":21}' http://localhost:8080/products/2
```

Get specific product:
```
curl http://localhost:8080/products/2
```
