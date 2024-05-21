# go-tinder

Lightweld tinder SDK for Golang

## Install

```
go get github.com/dvcrn/go-tinder
```

## Usage

Check the godoc to see what's available. Currently only the stuff I needed for myself is implemented.

### Authentication

To use this library, you need a tinder auth token. Login and retrieving such a token is not part of this library yet, so you'll have to get it through other means. A easy way is to copy the Auth header of a request to the tinder API.

### Get Matches

```go
c := tinder.NewClient("token")
m1, err := c.GetMatches(20, tinder.WithHasMessagesFilter(true))
if err != nil {
    log.Fatal(err)
}

fmt.Println("matches with 1 or more messages:", len(m1))

m2, err = c.GetMatches(20, tinder.WithHasMessagesFilter(false))
if err != nil {
    log.Fatal(err)
}
fmt.Println("matches with 0:", len(m2))
```

### Get Messages

```golang
c := tinder.NewClient("token")
msgs, err := c.GetMessages("matchID", 20)
if err != nil {
    log.Fatal(err)
}

fmt.Println("messages", len(msgs))
```

### Subscribe to updates (websocket)

```golang
c := tinder.NewClient("token")

// get updates for the past 5 minutes
res, err := c.GetUpdates(time.Now().Add(-5 * time.Minute))
if err != nil {
    fmt.Println("Error fetching updates: ", err)
}

fmt.Println(res)
```
