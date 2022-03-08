# jox-algorithm

This repository contains 5 different implementations of the [jox](https://github.com/keelefi/jox/) algorithm. The
repository contains a standardized set of test cases. All implementations are tested using the same tests.

## Algorithm

To run, `jox` needs to know all the pairs of `after` and `before` ordering of jobs. As the user of `jox` needs to only
state either `after` or `before` for each pair, `jox` needs to build all the missing directives.

## Job

Every `job` has a name, a list of jobs that are `after` and a list of jobs that are `before`. Job names are their unique
identifiers. The lists of jobs `after` and `before` address the other jobs by their unique identifiers, i.e. their
names.

The JSON representation of a job is as follows:

```
"<name>": {
    "after": ["<after1>","<after2>",...,"<afterN>"],
    "before": ["<before1>","<before2>",...,"<beforeN>"]
}
```

## Errors and Warnings

When running the algorithm, errors or warnings can occur. The exact errors and warnings are documented in `ERRORS.md`.

The JSON representation of a error is as follows:

```
{
    "type": <"ERROR"|"WARNING">,
    "error": "<error id>",
    "message": "<error message>"
}
```

## Tests

The test case data can be found under directory `tests`. The directory contains JSON files containing data for one test
case each. Each file contains `input` which denotes the starting point for the jobs. Furthermore, each file contains
`output` which denotes the expected result after the algorithm has been run.

Each file contains also `targets` to denote the jobs that are required to run. Finally, each file contains `errors` to
denote the optional errors that are expected to be encountered. The file structure is as follows:

```
{
  "input": { "<job1>": {...}, "<job2>": {...}, ... , "<jobN>": {...} },
  "target": ["<job1>", "<job2>",...,"<jobN>"],
  "output": { "<job1>": {...}, "<job2>": {...}, ... , "<jobN>": {...} },
  "errors": []
}
```

For a test case to pass, both expected `output` and `errors` need to match with the actual result, respectively.

## Status

Here's the current status of the various implementations:

| Language | Tests Created | Tests Passing |
| --- | --- | --- |
| Python | :heavy_check_mark: | ![Python](https://github.com/keelefi/jox-algorithm/actions/workflows/python.yml/badge.svg) |
| JavaScript | :heavy_check_mark: | ![javascript](https://github.com/keelefi/jox-algorithm/actions/workflows/javascript.yml/badge.svg) |
| C++ | :heavy_check_mark: | ![C++](https://github.com/keelefi/jox-algorithm/actions/workflows/cpp.yml/badge.svg) |
| Golang | :heavy_check_mark: | ![golang](https://github.com/keelefi/jox-algorithm/actions/workflows/golang.yml/badge.svg) |
| Scheme | :heavy_check_mark: | ![Scheme](https://github.com/keelefi/jox-algorithm/actions/workflows/scheme.yml/badge.svg) |

