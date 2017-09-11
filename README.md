# SQLpack 

[![Build Status](https://travis-ci.org/dmush/sqlpack.svg?branch=master)](https://travis-ci.org/dmush/sqlpack)

SQLpack is an attempt to speed up a PostgreSQL database development. It is a simple utility that can bundle sql files and execute them in a database.

## Intro

If you're looking for some tools that can help you with development of your PostgreSQL database then you might want to try the SQLpack.

You use your preferred editor or IDE while SQLpack watches for file changes, bundles them and executes in database. 

It doesn't provide any GUI. It's just a bunlder and a script runner with a few useful extensions.

## Getting started

Set up the Go (see [Go's official instructions](https://golang.org/doc/install)) and install SQLpack.

```bash
go install github.com/dmush/sqlpack
```

Check that installation was successfull.
```bash
sqlpack -h
```

Clone the repo and open a folder containing it.

```bash
git clone git://github.com/dmush/sqlpack && cd sqlpack
```

Install Docker (see [Docker's official instructions](https://docs.docker.com)) and start PostgreSQL server.

```bash
docker-compose up
```

You can use any other PostgreSQL server instead. Just specify the right database connection in the next step.

Run the SQLpack.

```bash
sqlpack -pg "user=postgres sslmode=disable" example/example.sql
```

This will bundle the `example.sql` file and execute it in PostgreSQL.

You also might want to try a testing mode. Execute the command below and try to edit `example/struct.sql` or `example/get_example_name.sql` (actually try to break something to fail test).

```bash
sqlpack test -pg "user=postgres" -w example
```

Testing mode is assumed to be the main mode while you are developing.

## Testing

SQLpack provides `test` command to run tests. Tests are files that have name mathing the `*_test.sql` pattern and are in the same folder as the original file.

```bash
sqlpack test -pg "user=postgres" example
```

This command will bundle and execute all tests in `example` folder and all child folders.

If you specify `-w` flag SQLpack will watch for file changes. If changed file is a test file then it bundles and executes the file. If not then SQLpack tries to find an appropriate test file.

```bash
sqlpack test -pg "user=postgres" -w example
```

## Bundling

SQLpack uses Go's [text/template](https://golang.org/pkg/text/template). All the basic functionality of Go's templates is available. A few additional functions provided by SQLpack.

### Include

Include allows to include contents of a certain file into specified place. It includes the file each time it appears in the code. If you have two includes of the same file then the file contents will be included twice.

Please take a look at `example/example.sql` file of this repo. Files `example/schema.sql` and `example/get_example_name.sql` are included into it.

```sql
{{ include "./schema" "./get_example_name" }}
```

### Import

Currently in plans only. It should include the file contents only once.

