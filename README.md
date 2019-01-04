# httpclient

httpclient is a rich implementation of the concepts provided in the talk "Embrace the Interface" by Tomas Senart.

## why

Because I implemented this in too many projects, and I want a go-to implementation that works out of the box.

## what does it do

httpclient is a set of decorators that compose on top of a `http.Client`.

Example:

```go
    // everytime this client is used, the X-Sent-By header will be set with httpclient1.0
    client := httpclient.New(http.DefaultClient, httpclient.Header("X-Sent-By", "httpclient1.0"))
    // error handling omitted for brevity
    req, _ := http.NewRequest("GET", "http://example.com", nil)
    res, _ := client.Do(req)
```

You can configure the root `http.Client` as well:

```go
    rootClient := &http.Client{
        Transport: &http.Transport{}, // or whatever you want here
    }
    client := httpclient.New(rootClient,
        httpclient.Header("X-Sent-By", "httpclient1.0"),
        httpclient.BasicAuthorization("my-user", "my-password"),
    )
    // error handling omitted for brevity
    req, _ := http.NewRequest("GET", "http://example.com", nil)
    res, err := client.Do(req)
    // ...
```

You can also create your own Decorators:

```go
    func MyDecorator() httpclient.Decorator {
        return func(client httpclient.Client) httpclient.Client {
            return httpclient.ClientFunc(func(r *http.Request) (*http.Response, error) {
                // do whatever you want with the request *before* it is made
                res, err := client.Do(r)
                // do whatever you want with the response now it is made
                return res, err
            })
        }
    }

    // and to use it:
    client := httpclient.New(http.DefaultClient,
        MyDecorator(),
        // you can use it alongside other decorators as well
        httpclient.Header("X-Sent-By", "httpclient1.0"),
        httpclient.BasicAuthorization("my-user", "my-password"),
    )
    req, err := http.NewRequest(...)
    // error handling, etc
    res, err := client.Do(req)

```

## active development disclaimer

This is in active development, so the APIs are prone to change.
For the most part though, I believe the client is going to grow in functionality that doesn't affect the core interfaces such as `Client` and `Decorator`.