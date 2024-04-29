# goexpert-lab-cloud-run

This project is an implementation of a simple application deployed using Google Clound Run.

The goal of this application is to check the current temperatures in Celsius, Kevin and Farenheit for a given location based on its CEP (Brazilian Zip Code).


## Commands to start the application locally

1. Just run the command below to start 

```sh
docker compose up -d
```

2. The server will be available at the url below:

```sh
http://localhost:8080
```

3. The main route is accessed by informing a valid CEP as a path parameter. 
http://localhost:8080/{cep}

  For example:
```sh
http://localhost:8080/01153000
```

## Informations regarding the Google Cloud run deploy

The web application deployed using Google Cloud Run can be accessed using the address https://lab-cloud-run-7s5babtwoa-uc.a.run.app/{cep}

For example:
```sh
https://lab-cloud-run-7s5babtwoa-uc.a.run.app/01153000
```

## Tests

On the `/api` folder there are four http files with examples of requests that can be made to the Rest server

If you are using VSCode as your IDE downloading the [Rest Client Plugin](https://marketplace.visualstudio.com/items?itemName=humao.rest-client) is an excellent way execute these requests.

The other tests for the project can be run using the command below:
```sh
go test -v
```
