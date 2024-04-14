# 🏁tbd🏎️✨

## A sweet and speedy code generator for dbt

> [!IMPORTANT]
> This project is still in its _very_ early stages. It should be relatively safe to use, but don't be surprised by bugs, breaking changes, or missing features. Please feel free to open issues or PRs if you have any feedback! Most importantly, be sure you're always pointing the output at an empty directory, or it may overwrite files with the same names.

_Disclaimer: This project is not affiliated with dbt Labs in any way. It is a personal project and is not officially supported by dbt Labs. I work at dbt Labs, but I develop this project in my own time._

## Supported warehouses

- [x] BigQuery
- [x] Snowflake
- [ ] Redshift
- [ ] Databricks
- [ ] Postgres
- [ ] DuckDB

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

The output will be a directory with the following structure:

```
your_build_dir/
├── _sources.yml
├── stg_model_a.sql
└── etc...
```

### LLM features

`tbd` has some neat alpha features that are still in development. One of these is the ability to generate documentation and tests for your sources via LLM. It uses [Groq](https://groq.com) running `mixtral-8x7b-32768` to do its inference. It's definitely not perfect, but it's pretty good! It requires setting an environment variable `GROQ_API_KEY` with your Groq API key.

The biggest thing to flag is that while Groq is in free beta, they have a very low rate limit on their API: 30 requests per minute. The actual inference on Groq is _super_ fast, but for now I've had to rate limit the API calls so it will take a minute or several depending on your schema size. Once Groq is out of beta, I'll remove the rate limit, but you'll of course have to pay for the API calls via your Groq account.

I will _definitely_ be adding other LLM providers in the future, probably Anthropic Claude 3 Opus as the next one so you can choose between maximum quality (Claude) or maximum speed (Groq, when I can remove the rate limit).

## To Do

- [ ] Add more test coverage
- [ ] Add Claude 3 Opus option
- [x] Add support for Snowflake
- [x] Add support for BigQuery
- [ ] Add support for Redshift
- [ ] Add support for Databricks
