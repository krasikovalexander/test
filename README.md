**Варианты деплоя**: как монолит; каждый эндпоинт по отдельности; на AWS Lambda (понадобится AWS аккаунт, dep и [Serverless Framework](https://serverless.com/framework/docs/providers/aws/guide/quick-start/)).

**Локально:**
cd src/service

dep ensure

go run main.go

POST http://localhost:3000/list

POST http://localhost:3000/rank

Content-Type: multipart/form-data



data                  xml file

source                string

destination           string

max_flights_in_route  int [optional]



POST http://localhost:3000/compare

Content-Type: multipart/form-data



data_a                  xml file

data_b                  xml file



POST http://localhost:3000/compare/routes

Content-Type: multipart/form-data



data_a                xml file

data_b                xml file

source                string

destination           string

max_flights_in_route  int [optional]