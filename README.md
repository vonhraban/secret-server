## Secret server

This application allows a user to store a secret key and retireve it. Additionally every key can have a maxinum number of allowed views after which it key is not available anymore. An optional expiry time for the secret can also be specified to assign TTL to it.

### Usage

#### Running

`go run .` if not already compiled, `./secret-server` if compiled already

#### Saving the key

The following curl requests adds a secret `abc` which allows a maximum of 5 views and has an expiry time of NOW + 5 minutes

```
curl --request POST \
  --url http://localhost/v1/secret \
  --header 'content-type: application/x-www-form-urlencoded' \
  --data secret=abc \
  --data expireAfterViews=5 \
  --data expireAfter=5
```

`secret` is the actual secret text we are storing
`expireAfterViews` specifies the maximum number of allowed views.
`expireAfter` specified number of minutes the secret should be available for, starting from NOW. Setting the value to 0 will mean that the secret never expires

#### Retrieving the key

An example of retieving the token stored under `176bd4fd-c312-43bf-91f7-03bde6326e6e` hash:

```
curl --request GET --url http://localhost/v1/secret/176bd4fd-c312-43bf-91f7-03bde6326e6e
```

#### Content type

This app is a polyglot, different content types can be asked for when talking to the server. Currently, the following are supported:

- `application/json`
- `application/xml`

### Testing

The usual `go test ./...` will do the trick. Test suites were written using Gingko testing framework.

### Considerations

Due to time contraints, I did not fully implements CQRS pattern here. In the ideal world I would support domain level events so that the app can react on changes of the state. For instance, I could not find a way to segregate the entire business logic into domain layer, and some business rules are duplicated in the application layer (http handlers). This is because it would be an antipattern to have GetSecret query execute a DecreaseRemainingViews command. Queries are not supposed to change any state, even if the change is not a direct manipulation but a side effect. I would prefer to have SecretViewed event emitted, so that the Event Listener could decrease the number of views left for the secret. This would make the code decoupled and easy to maintain; although it would increase the complexity of testing.

Yet, the commands and events are clearly separated, making it easy to test and understand the business logic in the Secret domain.

Also, some bits are not tested, e.g. user input validation in the http handlers and xml/json responses.

OpenAPI spec did not specify error code for the unrecoverable runtime error, but I took the liberty to return properly formatted 500 error in such a case.

No e2e testing for mongo
Some tests missing from validation
Stale view is returned and it is not tested if views were deducted since it is async