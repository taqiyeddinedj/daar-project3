# daar-project3

### How to Start
1. First download all the books, as you see it is ignored in the .gitignore, the library must be created locally

    `cd daar-project3`

    `go run cmd/download_books/main.go`
2. seconde step is to build the index, the main purpose is to speed up keyword search just like database indexing

    `go run cmd/build_index/main.go`

3. Last step, test the search mechanism, the keyword 'inputs' are hardcoded, this is only for test until we start building the webserver (UI)

    `go run cmd/test_search/main.go`