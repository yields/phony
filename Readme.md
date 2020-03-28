
# phony

  Tiny command line program that accepts a template and outputs fake data.

  ![](https://cldup.com/RZoAhReDqN.gif)

## Examples

```bash
# publish email to nsq every 1ms.
echo '{"email":"{{email}}", "subject": "welcome!"}' \
  | phony --tick 1ms \
  | json-to-nsq --topic users

# add users to FoundationDB.
echo "'set {{username}} {{avatar}}'" \
  | phony \
  | xargs -L1 -n3 fdbcli --exec

# add users to MongoDB.
echo "'db.users.insert({ name: \"{{name}}\" })'" \
  | phony \
  | xargs -L1 -n1 mongo --eval

# add users to Redis.
echo "set {{username}} {{avatar}}" \
  | phony \
  | xargs -L1 -n3 redis-cli

# send a single request using curl.
echo 'country={{country}}' \
  | phony --max 1 \
  | curl -d @- httpbin.org/post
```

## Installation

```bash
$ go get github.com/yields/phony
```

## Usage

```text

Usage: phony
  [--tick d]
  [--max n]
  [--list]

  phony -h | --help
  phony -v | --version

Options:
  --list          list all available generators
  --max n         generate data up to n [default: -1]
  --tick d        generate data every d [default: 10ms]
  -v, --version   show version information
  -h, --help      show help information

```

## Generators

```text
  avatar
  color
  country
  country.code
  domain
  domain.name
  domain.tld
  double
  email
  event.action
  http.method
  id
  ipv4
  ipv6
  ksuid
  latitude
  longitude
  mac.address
  name
  name.first
  name.last
  product.category
  product.name
  state
  state.code
  timezone
  unixtime
  username
  uuid
```

## License

  (MIT), 2014 Amir Abu Shareb.
