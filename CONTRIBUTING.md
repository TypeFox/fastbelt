# Contributing to Fastbelt

Issues and pull requests are welcome.

## Performance Benchmarking

Except for changes where the expected performance impact is negligible, contributors should run benchmarks both with and without their changes (for example by comparing their branch with `main`) and include the results in the Pull Request description.
Since some of the benchmarks are susceptible to runtime/system noise, it can be beneficial to increase the benchmark duration via the `-benchtime` flag.

A GitHub Action will automatically run all benchmarks and post them as a summary to the respective action run. If a benchmark exceeds the baseline by 1.5x, a warning comment will be generated on the PR.
Refer to [this workflow](./.github/workflows/benchmark.yml) for more information.
