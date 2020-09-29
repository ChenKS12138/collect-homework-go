# collect-homework-go

![testing](https://github.com/ChenKS12138/collect-homework-go/workflows/testing/badge.svg)

作业提交平台

## Table of Contents

- [collect-homework-go](#collect-homework-go)
  - [Table of Contents](#table-of-contents)
  - [How To Run](#how-to-run)
  - [Environment Variables](#environment-variables)
  - [Support](#support)
  - [Contributing](#contributing)

## How To Run

```sh
vim .env
make install
make run
```

## Environment Variables

```sh
PORT=3000

DB_NETWORK=tcp
DB_ADDR=127.0.0.1:5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_DATABASE=postgres
DB_DEBUG=false

EMAIL_USER=user@example.com
EMAIL_PASSWORD=secret
EMAIL_PREVENT=false

STORAGE_PATH_PREFIX=./tmp

SUPER_USER_NAME=admin
SUPER_USER_EMAIL=admin@example.com
SUPER_USER_PASSWORD=password

```

## Support

Please [open an issue](https://github.com/fraction/readme-boilerplate/issues/new) for support.

## Contributing

Please contribute using [Github Flow](https://guides.github.com/introduction/flow/). Create a branch, add commits, and [open a pull request](https://github.com/fraction/readme-boilerplate/compare/).
