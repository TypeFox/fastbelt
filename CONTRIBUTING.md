# Contributing to Fastbelt

Issues and pull requests are welcome.

## Performance Benchmarking

Use the following command as the standard end-to-end benchmark:

```bash
go test -run '^$' -bench '^BenchmarkWorkspaceCycle$' ./examples/statemachine
```

For more targeted changes (for example in the lexer or parser), other benchmarks may be more suitable.

When running benchmarks, execute them at least 5 times and use the average `ms/resource` value as the standard metric. This helps reduce noise from other processes running on your machine, especially because Fastbelt makes heavy use of multithreading. You should also limit the number of CPU cores to 1 (`go test -cpu=1`) unless you need to compare in a multicore scenario.

For convenience, you can run the internal benchmark helper:

```bash
go run ./internal/bench -runs 30 -cpu 1
```

Except for changes where the expected performance impact is negligible, contributors should run benchmarks both with and without their changes (for example by comparing their branch with `main`) and include the results in the Pull Request description.
