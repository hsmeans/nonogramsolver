# Nonogram Solver

A nonogram solver written in golang. Written for educational purposes.

### Prerequisites

- Ensure Go is installed: https://go.dev/
- A JSON file with the puzzle's hints in the following format (see test/data/example.json):

```
{
  "columns": [...]
  "rows": [...]
}
```

### Installing using `go install`

- run `go install github.com/hsmeans/nonogramsolver@latest`
- run `nonogramsolver [INPUT_FILE]`

### Building and running from source

- Clone the repository
- Run `go build`
- Run `./nonogramsolver test/data/example.json`
- The solution (or partial solution if not solvable) will be displayed in the output
