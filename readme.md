# Synapsis
This project is a Golang-based application for **Synapsis**.
It uses Golang 1.24.5, echo framework, and Postgres

---

## 🚀 Requirements
Before starting, ensure you have:
- Golang >= 1.24.5
- Postgres

---

## ⚙ Manual Setup instruction
### 1. clone the repository

```shell
git clone https://github.com/oktopriima/synapsis
```
### 2. Configure the local environment
Create the database for `orders` and `inventory` services on your local postgres.
Adjust the environment for both of the service
```shell
postgres:
  host: "127.0.0.1"
  database: "orders"
  password: "sapi"
  port: "5432"
  user: "root"
```
update the credentials ` env.yaml` both on `order` and `inventory` on this part.

### 3. Update modules
Since this project separated on two different services, you have to update modules on both of them and **gRPC** Proto Definitions.

#### ORDER Service
From root folder
```shell
cd order
go mod tidy
go mod vendor
```

#### INVENTORY Service
From root folder
```shell
cd inventory
go mod tidy
go mod vendor
```

#### GRPC Proto Definition
From root folder
```shell
cd proto-definitions
go mod tidy
go mod vendor
```

### 4. Running the migration and seeder
The migration already provided both on `order` and `inventory` service. For running the migration, you can run this command
```shell
# order service migration
go run order/database/migration/main.go
```
```shell
# inventory service migration
go run inventory/database/migration/main.go
```
This project also provided a seeder on `inventory` service for data example.

```shell
go run inventory/database/seeder/main.go
```
It will give you some `product` with `stock` data example.

### 5. Running the application
Since this project have two different services, you should run both of them on different terminal.
#### ORDER service
From your root folder
```shell
go run order/main.go
```
#### INVENTORY service
From your root folder
```shell
go run inventory/main.go
```

### 6. Access the application
- For `order` HTTP service will run on: 👉 http://localhost:8000
- For `inventory` HTTP service will run on: 👉 http://localhost:8001
- For `inventory` RPC service will run on: localhost:5000

---

## 📝 Notes
This project provide some example unit test on `order` service.
```shell
## move to order service folder
cd order

## running the unit test
go test ./...
```

## 📂 Project Structure
```
synapsis/
├── order/
├─────── app/
├─────── bootstrap/
├─────── config/
├─────── database/
├─────── router/
├─────── test/
├─────── go.mod
├─────── go.sum
├─────── env.yaml
├── inventory/
├─────── app/
├─────── bootstrap/
├─────── config/
├─────── database/
├─────── router/
├─────── go.mod
├─────── go.sum
├─────── env.yaml
├── proto-definitions/
├─────── inventory/
├─────── go.mod
└─────── go.sum
```