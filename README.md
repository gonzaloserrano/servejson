# servejson

Serve a json file via HTTP with CORS and content-type headers in a route.

    # execute server
    ❯ servejson -file fixtures/test.json -route api/file.json -port 8080
    2018/03/12 20:38:44 Serving fixtures/test.json at /api/file.json in port 8080

    # execute client
    ❯ http localhost:8080/api/file.json | jq .
    {
      "foo": "bar"
    }

This is useful to test HTTP requests from a browser (uses OPTIONS, needs CORS) to a real server.

## Install or update

    go get -u github.com/gonzaloserrano/servejson

## Why

I wanted to serve a json easily and quickly, not to download the whole internet or configure something complex.

## TODO

- [ ] check is a valid json file
- [ ] enable/disable CORS flag
