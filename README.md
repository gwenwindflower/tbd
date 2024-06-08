# 🏁 tbd 🏎️✨

## 💽 A sweet and speedy code generator for dbt

_**Disclaimer**: This project is not affiliated with dbt Labs in any way. It is a personal project and is not officially supported by dbt Labs. I work at dbt Labs, but I develop this project in my own time._

`tbd` quickly generates code for new and existing dbt projects with a variety of features:

- generates dbt sources and staging models from your raw schemas' metadata
- scaffolds out a complete dbt project (_optional_)
- saves connection details as dbt profiles for future use (_optional_)
- uses LLMs to infer documentation and tests for your sources (_optional_)

It's designed to be super fast and easy to use with a friendly TUI that fast forwards you to writing meaningful dbt models as quickly as possible.

Still confused? No worries! Check out the [example output diagrams](#without-scaffolding-a-dbt-project) and it will make sense.

### It's the **_easy button_** for dbt projects.

> [!NOTE]
> `tbd` does not depend on or utilize dbt directly, although it can leverage its `profiles.yml` configuration files. If you're using dbt Cloud CLI, you'll have to either create a `profiles.yml` it can access or enter connection details manually. It can't access the connections configured in dbt Cloud.

#### Quickstart

```bash
brew tap gwenwindflower/homebrew-tbd
brew install tbd
tbd
```

If you're new to dbt, [check out the wiki](https://github.com/gwenwindflower/tbd/wiki) for some great learning resources and tips on setting up a cozy coding environment!

## 🔌 Supported warehouses

- [x] BigQuery
- [x] Snowflake
- [ ] Redshift
- [x] Databricks
- [x] Postgres
- [x] DuckDB

If you don't have a cloud warehouse, but want to spin up a dbt project with `tbd` I recommend either:

- **BigQuery** — they have a generous free tier, authenticating with `gcloud` CLI is super easy, and `tbd` requires very few manual configurations. They also have a ton of great public datasets you can model.
- **DuckDB** — you can work completely locally and skip the cloud altogether. You will need to find some data, but DuckDB can _very_ easily ingest CSVs, JSON, or Parquet, so if you have some raw data you want to work with, this is a great option as well.

## 💾 Installation

For the time being this project is **only compatible with MacOS**. Linux and Windows support are definitely on the roadmap, just have to wait for a day when I can dive deep into CGO and understand the intricacies of building for those platforms. The easiest way to install is via Homebrew:

```bash
brew tap gwenwindflower/homebrew-tbd
brew install tbd
```

If you have Go installed, you can also install it via `go install`:

```bash
go get -u github.com/gwenwindflower/tbd@latest
go install github.com/gwenwindflower/tbd@latest
```

That's it! It's a single binary and has no dependencies on `dbt` itself, for maximum speed it operates directly with your warehouse, so you don't even need to have `dbt` installed to use it. That said, it _can_ leverage the profiles in your `~/.dbt/profiles.yml` file if you have them set up, so you can use the same connection information to save yourself some typing.

## 🔐 Warehouse-specific setup

`tbd` at present, for security, only supports SSO methods of authentication. Please check out the below guides for your target warehouse before using `tbd` to ensure a smooth experience.

### ❄️ Snowflake

Snowflake uses `externalbrowser` SSO authentication. It requires that you have SSO set up in your warehouse, it will then open a browser tab to authenticate and refresh a token in your **local** keychain. You'll be prompted to enter _your computer's user login_ (not your Snowflake password) to retrieve the token locally from your keychain.

Here's an example of a profile that works with `tbd` on Snowflake:

```yaml
lothlorien_labs:
  target: dev
  outputs:
    dev:
      type: snowflake
      account: galadriel123
      user: winnie@lothlorienlabs.com
      role: transformer
      authenticator: externalbrowser
      database: lothlorien
      schema: mallorn_trees
      threads: 16
```

### 🌯 BigQuery

BigQuery requires that you have the `gcloud` [CLI installed](https://cloud.google.com/sdk/docs/install) and authenticated for whatever projects you target.

```bash
gcloud auth application-default login
```

Here's an example of a profile that works with `tbd` on BigQuery:
```yaml
moria_mining:
  target: dev
  outputs:
    dev:
      type: bigquery
      method: oauth
      project: moria-mining
      dataset: dbt_gimli
      threads: 16
```

I will likely bring in some other authentication options soon, but this is the easiest and most secure.

### 🦆 DuckDB

Using local DuckDB doesn't require authentication, just an existing DuckDB database to query against and the path to that file. Note that when passing the path manually (i.e. _not_ via a dbt profile) the path should be relative to the current working directory. So if you're in a DuckDB dbt project that has, as is typical, the DuckDB file at the root of the project, you can just pass the name of the file (including the extension, as duckdb files can be either `.duckdb` or `.db`).

> [!NOTE]
> I've built-in support for [MotherDuck](https://motherduck.com/), you just need to set an env var called `MOTHERDUCK_TOKEN` with your service token, then pass the path as `md:`, **but** until MotherDuck upgrades to v10 this requires you to use DuckDB 0.9.2 locally for compatibility. MotherDuck says the upgrade will happen any day now so hopefully this note will be removed soon!

Here's an example of a profile that works with `tbd` on DuckDB:

```yaml
legolas_analytics:
  target: dev
  outputs:
    dev:
      type: duckdb
      path: /Users/winnie/dev/jaffle_shop_duckdb/mirkwood.duckdb
      database: mirkwood
      schema: archers
      threads: 16
```


### 🐘 Postgres

Postgres setup is fairly normal! It just needs a host, port, user, password, and of course a database and schema that you want to model. The default values point to a localhost instance on port 5432, but you can change that to whatever you need and it should work with cloud-hosted Postgres instances as well (though I haven't yet tested this myself). Tests have been running on Postgres 16, but `tbd` makes very standard queries to `INFORMATION_SCHEMA` and doesn't do anything fancy with the connection, so it should work with any relatively modern Postgres version probably.

Here's an example of a profile that works on Postgres with `tbd`:

```yaml
gondor_security:
  target: dev
  outputs:
    dev:
      type: postgres
      host: localhost
      port: 5432
      user: aragorn
      password: arw3n1sBa3
      database: gondor
      schema: minas_tirith
      threads: 8
      connect_timeout: 10
```

### 🧱 Databricks

Databricks for now only works with Personal Access Token authentication. You can generate one in your Databricks account settings under the Developer tab. Put this in an environment variable called something like `DATABRICKS_TOKEN` (the default option). `tbd` will _not_ let you pass the value directly.

You'll need the information found in your SQL Warehouse's 'Connection details' tab which has the Hostname and Http Path you'll need, unless you already have it in a dbt profile, then you can choose that profile, pick which catalog and schema you want to target, and you're good to go!

If you're using a dbt profile, it will need to be configured for PAT authentication and referencing the token via an environment variable. The profile should look something like this:

```yaml
gandalf_wizardy_co:
  target: dev
  outputs:
    dev:
      type: databricks
      catalog: maiar
      schema: wizards
      host: 1234.gandalf.databricks.com
      http_path: /sql/1.0/warehouses/1234
      token: "{{ env_var('DATABRICKS_TOKEN') }}"
      threads: 16
```

## Usage

The tool has a lovely TUI interface that will walk you through the necessary steps. You can run it with the following command:

```bash
tbd
```

It will guide you through inputting the necessary information to generate your `_sources.yml` and staging models for any schema you point it at. The idea is you point it at your `raw` unmodeled schema(s), and it will generate the necessary files to get those models up and running in dbt.

The output will be a directory with the following structure:

### Without scaffolding a dbt project

```terminal
your/build/dir/
├── _sources.yml
├── stg_raw_customers.sql
├── stg_raw_items.sql
├── stg_raw_orders.sql
├── stg_raw_products.sql
├── stg_raw_stores.sql
└── stg_raw_supplies.sql
```

### With scaffolding a dbt project

```terminal
your/build/dir
├── analyses
├── data-tests
├── dbt_project.yml
├── macros
├── models
│   ├── marts
│   └── staging
│       ├── _sources.yml
│       ├── stg_raw_customers.sql
│       ├── stg_raw_items.sql
│       ├── stg_raw_orders.sql
│       ├── stg_raw_products.sql
│       ├── stg_raw_stores.sql
│       └── stg_raw_supplies.sql
├── seeds
└── snapshots
```

### 🦙 LLM features

`tbd` has some neat alpha features that infer documentation and tests for your columns. There are multiple supported LLMs via API: Groq running Llama 3 70B, Anthropic Claude 3 Opus, and OpenAI GPT-4 Turbo. They have very different rate limits (these are limits in the APIs that `tbd` respects):

| Model                      | Requests per minute  |
|----------------------------|----------------------|
| Groq (Running Llama 3 70B) | 30                   |
| Claude 3 Opus              | 5                    |
| GPT-4 Turbo                | 500                  |

As you can see, if you have anything but a very smol schema, you should stick with OpenAI. When Groq ups their rate limit after they're out of beta, that will be the fastest option, but for now, OpenAI is the best bet. The good news is that GPT-4 Turbo is _really_ good at this task (honestly better than Claude Opus) and pretty dang fast! The results are great in my testing.

I'm going to experiment very soon with using structured output conformed to dbt's JSON schema and passing entire tables, rather than iterating through columns, and see how it does with that. If it works that will be significantly faster as it can churn out entire files (and perhaps improve quality through having more context) and the rate limits will be less of a factor.

### 🌊 Example workflows

`tbd` is designed to be self-explanatory, but just in case you get stuck, we have a blossoming wiki [with various example workflows](https://github.com/gwenwindflower/tbd/wiki/Example-workflows) that take you through step-by-step.

## 😅 To Do

### Tier 1
- [ ] Get to 100% test coverage
- [x] Add Claude 3 Opus option
- [x] Add OpenAI GPT-4 Turbo option
- [x] Add support for Snowflake
- [x] Add support for BigQuery
- [ ] Add support for Redshift
- [x] Add support for Databricks
- [x] Add support for Postgres
- [x] Add support for DuckDB
- [x] Add support for MotherDuck
- [ ] Build on Linux
- [ ] Build on Windows

### Tier 2

- [ ] Add `tbd rerun` option, loading a previous form run from a config file in `~/.config/tbd` to 'rerun' on a different schema

## 🤗 Contributing

I welcome Discussions, Issues, and PRs! This is pre-release software and without folks using it and opening Issues or Discussions I won't be able to find the rough edges and smooth them out. So please if you get stuck open an Issue and let's figure out how to fix it!

If you're a dbt user and aren't familiar with Go, but interested in learning a bit of it, I'm also happy to help guide you through opening a PR, just let me know 💗.
