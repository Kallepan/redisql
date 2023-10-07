# RediSQL

## Description

A simple abstraction layer to simulate an SQL database on top of Redis written in GoLang. Supports the creation of a table with a primary key, insertion of rows, deletion of rows, and querying of rows by primary key. The table is a hashmap identified by: table_name:primary_key. Columns and values are stored as key-value pairs in the hashmap.

## Usage

```bash
# Open in devcontainer
# Run tests
bash hack/test.sh

# Build
TODO
```

## TODO

- [ ] Add support for secondary indexes
- [ ] Add support for querying by secondary indexes
