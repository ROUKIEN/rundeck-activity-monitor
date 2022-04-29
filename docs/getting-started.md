# Getting started

## Starting ram

You can easily start a ram from the `ghcr.io/roukien/ram` docker image

```shell
# check the available tags on repo packages section
docker pull ghcr.io/roukien/ram:v0.0.7
```

> You can check the available tags on the [ram container page](https://github.com/ROUKIEN/rundeck-activity-monitor/pkgs/container/ram)

## Persistence

`ram` requires a postgresql database to work. Start one the way you want, then run the following command to configure your database:

```shell
# configure the RAM_DB_DSN env variable
export RAM_DB_DSN=postgres://youknowthe:drill@database:5432/ram
# configure the database
ram database update
```

This command will create the required tables in postgres. Once that step is done, you can start defining your rundeck instances.

## Configuring Rundeck instances

### Generating a Rundeck API token

For each rundeck instance that you want to monitor you'll need an API token that is at least able to list projects & project executions. [Checkout the rundeck documentation](https://docs.rundeck.com/docs/manual/10-user.html#user-api-tokens) about how to generate one.

### Configuration file

Create a yml file with the following content:

```yaml
# you can place any number of rundeck instances under that section
instances:
  # the rundeck instance identifier in ram
  rundeck_foo:
    # the rundeck URL
    url: http://rundeck.foo
    # the rundeck API token
    token: my.rundeckT0k3n
    # the rundeck API version to use
    apiversion: 41
    # an API timeout in milliseconds
    timeout: 5000
  rundeck_bar:
    url: http://rundeck.bar
    token: my.rundeckT0k3n.bar
    apiversion: 38
    timeout: 5000
```

## Running RAM

RAM is composed of two main components: the scraper and the webserver.

The scraper purpose is to retrieve the rundeck activity and persist it into ram database.

The Webserver aims to visualize rundeck executions and offer an intuitive UI to search for executions across all managed instances.

### Scraping

You can start a scraping process of all your rundeck instances with the following command. It will scrape in parallel every projects of every rundeck instance that you configured.

> You can either scrape your instances with a one-shot process (that will exit once the scraping is over) or as a daemon process, which scrapes the instances on a regular interval

As an example

```shell
export RAM_DB_DSN=postgres://...
ram scrape --newer-than=24h
```

> The format of the option `newer-than` accept any value that can be parsed by [time.ParseDuration() in golang](https://pkg.go.dev/time#ParseDuration)

You can also decide to scrape every executions that ended in a given timeframe with the `begin` and `end` options:

```shell
export RAM_DB_DSN=postgres://...
ram scrape --begin="2022-04-21T00:00:00.000Z" --end="2022-04-24T00:00:00.000Z"
```

If you want to scrape your instances every five minutes, you can use the `interval` option

```shell
export RAM_DB_DSN=postgres://...
ram scrape --newer-than=5m --interval=5m
```

### Serving

You can serve the UI by starting the webserver

```shell
export RAM_DB_DSN=postgres://...
ram serve
```

By default, the webserver will start on port 4000.

## Best practices

### Scrape a few but scrape often

RAM will query the executions that _ended_ in the given timeframe. That means that if you run the `scrape` command every minutes, you only need to scrape for the last past 5minutes (with the `--newer-than=5m` option). Scraping will be way quicker and will have a smaller footprint on your rundeck server.

If you need to scrape jobs in a given timeframe (because something went wrong during a scrape), you can decide to scrape that timeframe by using the `begin` & `end` options.
