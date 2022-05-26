# Interstellar Cab Service

Simple REST API that gives users the possibility to reserve a starship that is supplied from an external API. 

The API used is:
https://swapi.dev/documentation

# Requirements

- Docker 20.10.12+

## Endpoints

### List ships

Should return a list of available ships, in json format:

```sh
  curl -i 'localhost:3000/ships'
```
#### Output
The output is a json array with objects with the following fields:
- Id
- Name
- Model
- Cost

### User registration

Basic user registration for authorized calls. It returns a token to be used on authorized calls.

#### Input
The input data is as follows:

- Email
- PasswordHash
- Vip (whether the user is premium or not)

```sh
curl -d '{"email":"examle@email.com", "passwordHash":"123", "vip":true}' -H "Content-Type: application/json" -X POST 'http://localhost:3000/signup'
```

#### Output
It will return a json object with the fields:
- Email
- Token

### User login

Basic user login for authorized calls. The usage and purpose is exactly the same as the [User registration](###User-registration). 

#### Input
The input will be:

- Email
- PasswordHash

```sh
curl -d '{"email":"examle@email.com", "passwordHash":"123"}' -H "Content-Type: application/json" -X POST 'http://localhost:3000/login'
```

#### Output
It will return the same json object as [User registration](###User-registration).


### Ship reservation (authorized)

Enpoint where the user is able to make reservations from the available ships. The user will need to provide a valid ship id, if not it will be prompted with an error.

#### Input
The input parameters are:

- Id (string)
- date_from
- date_to
- Token (as a header)

```sh
curl -i -d '{"id": "2","date_from": "2022-05-20T15:00:00.000Z","date_to": "2022-05-26T16:00:00.000Z"}' -H "Content-Type: application/json" -H "Token: $TOKEN" -X POST 'http://localhost:3000/reservations'
```

#### Output
The output will be a 200 code if the reservations has been completed succesfully.

There are some restrictions for this endpoint:

- Ship cannot be reserved by two users in the same day
- User cannot make a reservation for longer than 15 days
- Non-premium users cannot make a reservation for ships with cost over 250000

### User reservations (authorized)

List current user reservations. It will return a json array with the reservation information for that user.

#### Input
 The input parameters are:

- Token (as a header)

```sh
curl -i -H "Token: $TOKEN" 'http://localhost:3000/reservations'
```

#### Output
The output will be a json array with objects with fields:

- Reservation Id
- Date from
- Date to
- Email
- Creation date

## Installation

To install it, you will need docker installed in your machine. There's a Dockerfile responsible for building a lightweight image.

First step is to build the docker image:

```sh
docker build . --tag interestellar
```

Then, you will be able to start the docker server and call the API endpoints to your desired port:

```sh
docker run -d -p $PORT:3000 --env JWT_SECRET='secret' interestellar
```