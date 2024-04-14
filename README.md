# 🏁tbd🏎️✨

## A sweet and speedy code generator for dbt

> [!IMPORTANT]
> This project is still in its _very_ early stages. It should be relatively safe to use, but don't be surprised by bugs, breaking changes, or missing features. Please feel free to open issues or PRs if you have any feedback! Most importantly, be sure you're always pointing the output at an empty directory, or it may overwrite files with the same names.

## Installation

For the time being this project ideally requires `go`. When I've gotten test coverage up to a reasonable level and covered another dbt adapter or two, I'll set up a Homebrew tap. In the meantime, you can install it with the following command:

```bash
go install github.com/gwenwindflower/tbd@latest
```

You can also download the binary from the [releases page](https://github.com/gwenwindflower/tbd/releases) and add it to your PATH if you're comfortable with that.

## Usage

The tool has a lovely TUI interface that you can use by running the following command:

```bash
tbd
```

It will walk you through inputting the necessary information to generate your `_sources.yml` and staging models for any schema you point it at. The idea is you point it at your `raw` unmodeled schema(s), and it will generate the necessary files to get those models up and running in dbt.

## To Do

- [ ] Add more test coverage
- [x] Add support for Snowflake
- [x] Add support for BigQuery
- [ ] Add support for Redshift
- [ ] Add support for Databricks
