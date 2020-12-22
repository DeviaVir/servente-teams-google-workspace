# servente-teams-google-workspace

The `servente-teams-[name]` pattern is used by providers that enrich Servente,
in this case, `servente-teams-google-workspace` trails Google Workspace's admin
directory API and returns a list of all teams, and is capable of returning a
list of all groups (teams) a member (email) is part of.


## Docker releases

https://hub.docker.com/r/deviavir/servente-teams-google-workspace

Example run:

```
$ docker run -it --rm \
  -v /home/chase/Credentials/credentials.json:/secrets/credentials.json:ro \
  -v /home/chase/Credentials/tls:/secrets/tls:ro \
  deviavir/servente-teams-google-workspace:0.1.0 /servente-teams-google-workspace \
    --credentialsPath=/secrets/credentials.json \
    --userEmail=user@domain.ext \
    --tls-cert-path=/secrets/tls/cert.pem \
    --tls-key-path=/secrets/tls/key.pem
```

## Setup

Find instructions for creating a service account configured for access to Google
Workspaces [here](https://developers.google.com/admin-sdk/directory/v1/guides/delegation).

Make sure to enable the Admin SDK API on the project where you created the
service account [here](https://console.developers.google.com/apis/library/admin.googleapis.com?project=blockstream-source).

Once configured, you can start this binary with the following options:

```
$ ./servente-teams-google-workspace \
  --accessKey=my-secret-password \
  --credentialsPath=/some/path/credentials.json \
  --userEmail=user@domain.ext
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
