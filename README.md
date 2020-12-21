# servente-teams-google-workspace

The `servente-teams-[name]` pattern is used by providers that enrich Servente,
in this case, `servente-teams-google-workspace` trails Google Workspace's admin
directory API and returns a list of all teams, and is capable of returning a
list of all groups (teams) a member (email) is part of.

## Setup

Find instructions for generating the service account credentials
[here](https://developers.google.com/admin-sdk/directory/v1/quickstart/go#step_1_turn_on_the).

Once configured, you can start this binary with the following options:

```
$ ./servente-teams-google-workspace \
  --accessKey=my-secret-password \
  --credentialsPath=/some/path/credentials.json
```

The access key is required for all requests to the API this application exposes,
you can fill in this access key in your Servente organization's settings.

## Development

### Generate self-signed cert

```
$ go run /usr/local/go/src/crypto/tls/generate_cert.go \
  --rsa-bits=2048 \
  --host=localhost
2018/10/16 11:50:14 wrote cert.pem
2018/10/16 11:50:14 wrote key.pem
```

### Testing

```
$ go test -short -race -v ./...
```

### Building

```
$ GPG_KEY=... make dist
```

### Deploying

Bump the release version in the `Makefile`.

```
$ make docker
```

## API

```
$ curl --header "Servente-Access-Key: servente-secret-access-key-001\!" -k "https://127.0.0.1:4001/api/v1/teams/membership?member=chase@email.com | jq
$ curl --header "Servente-Access-Key: servente-secret-access-key-001\!" -k "https://127.0.0.1:4001/api/v1/teams/membership/chase@email.com | jq
$ curl -k https://127.0.0.1:4001/api/v1/teams/list\?servente-access-key\=servente-secret-access-key-001! | jq
$ curl -k https://127.0.0.1:4001/api\?servente-access-key\=servente-secret-access-key-001! | jq
```
