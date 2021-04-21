# Fibonacci By Memory

Golang and PostgreSQL solution for finding the nth value in the Fibonacci sequence using
memorization. 

### Install & Execute

#### Prerequisites
- Git
- Make
- Docker
- Docker Compose
- Golang ( optional if user wished to build the binary locally )

#### Run
- Clone the repository
```shell
git clone https://github.com/lfordyce/fibonacci-by-memory.git && cd fibonacci-by-memory 
```
- Run using the Makefile
```shell
make compose
```
- Additional Makefile commands can be viewed with:
```shell
make help
✓ usage: make [target]

build                          Build program binary
clean                          Cleanup everything
compose                        Run docker-compose
docker-pull                    Pull latest Docker images in preparation for build
docker                         Build docker image
fmt                            Run gofmt on all source files
get                            Run go get for dependencies
help                           - Show help message
lint                           Run golint
```

### API:
- Get the nth value in the Fibonacci sequence by providing the ordinal value.
```shell
curl http://localhost:8000/v1/api/fib/{ordinal}
# example
curl http://localhost:8000/v1/api/fib/11
```
- Fetch the number of memoized results less than a given value (e.g. there are 12 intermediate results less than 120)
```shell
curl http://localhost:8000/v1/api/fib/records/{count}
# example
curl http://localhost:8000/v1/api/fib/records/12
```
- Clear the data store.
```shell
curl -XDELETE http://localhost:8000/v1/api/fib/
```