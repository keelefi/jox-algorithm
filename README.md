# Algorithm

This directory structure contains tests and implementations of the main algorithm needed for `jox` to run. We have a
standardized set of test cases that each language implementation needs to pass before being accepted. To keep all the
test cases across different languages in sync, we store the test case data as JSON.

## Requirements

Before running the actual algorithm, the initial state of the graph structure of jobs needs to be loaded. Thereafter the
algorithm runs.

The algorithm needs to traverse all jobs. For each job the algorithm needs to add the counterpart for each `after` and
`before`, respectively. For example, if job `B` is `after` job `A`, then job `A` must be `before` job `B`.

## Job

As we are developing the job dependency graph algorithm, our definition of a job can be far simplified. We need only job
`name`, `after` and `before`. This is expressed in JSON as follows:

```
"<name>": {"after": ["<after1>","<after2>",...,"<afterN>"], "before": ["<before1>","<before2>",...,"<beforeN>"]}
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

## Running

Each langauge implementation comes with a shell script `run_tests.sh` which runs the tests and returns 0 if pass and 1
if fail. In this directory the `run_tests.sh` will go through all implementations and run their tests, respectively. If
any fail, it returns 1, otherwise it returns 0.

## Status

Here's the current status of the various implementations:

| Language | Tests Created | Tests Passing |
| --- | --- | --- |
| Python | :heavy_check_mark: | ![Python](https://github.com/keelefi/jox-algorithm/actions/workflows/python.yml/badge.svg) |
| JavaScript | :heavy_check_mark: | ![javascript](https://github.com/keelefi/jox-algorithm/actions/workflows/javascript.yml/badge.svg) |
| C++ | :heavy_check_mark: | ![C++](https://github.com/keelefi/jox-algorithm/actions/workflows/cpp.yml/badge.svg) |
| Golang | :heavy_check_mark: | ![golang](https://github.com/keelefi/jox-algorithm/actions/workflows/golang.yml/badge.svg) |
| Scheme | :heavy_check_mark: | ![Scheme](https://github.com/keelefi/jox-algorithm/actions/workflows/scheme.yml/badge.svg) |

