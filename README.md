# Redis

I wrote my own version of Redis using Go.

## Installation

Still need to set up go.mod and go.sum files.

```bash
...
```

## Run

You must have the Redis CLI installed. After cloning the repo...

- First run the code via
```bash
go run ./src
```

- In a separate terminal run
```bash
redis-cli
```

## Features

My version supports: **SET**, **GET**, **HSET**, **HGET**, **HGETALL**. The commands are __not__ case sensitive.

## Enhancements

- Include more detailed error handling
- Data persistence when the server stops running

## License

[MIT](https://choosealicense.com/licenses/mit/)