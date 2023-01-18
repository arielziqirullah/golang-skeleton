# Simple Golang Skeleton

## About <a name = "about"></a>

Golang Rest API sample with MySQL integration using Framework Gin and Gorm.

## Getting Started <a name = "getting_started"></a>


### Installing

First step need to copy .env.exmplae into .env file to create envirionment variable.

```
cp .env.example .env
```

Next step fill in variable on .env file

```
#APP
APP_RUN=
APP_DEBUG=
RUNNING_PORT=

#DATABASE
DB_HOST=
DB_PORT=
DB_NAME=
DB_USER=
DB_PASSWORD=

#DATABASE CONFIG
MIGRATION=
MAX_IDLE_CONNECTION=
MAX_OPEN_CONNECTION=
MAX_LIFETIME_CONNECTION=

#JWT
TOKEN_TYPE="Bearer"
ACCESS_TOKEN_MAXAGE=
REFRESH_TOKEN_MAXAGE=
JWT_SECRET=
```

