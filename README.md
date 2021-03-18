# Pubpipe

An implementation of the [patchbay.pub](https://patchbay.pub) server format that implements queues, pubsub, and JS EventSource events.

# Building

Pubpipe has no external dependencies.

Clone the repository and run `go build` to build a binary.

Tested with Go 1.16, but other versions may also work.

# Usage

Run the binary and a server should be hosted automatically.

To change the port, change the `PUBPIPE_PORT` environment variable to the port to run on. The default port is `8041`.

There are three types of methods:

## Queue

Send:

    > curl localhost:8041/queue/channel -d "Hello World!"

Receive:

    > curl localhost:8041/queue/channel
    Hello World!

`POST` requests will queue up messages one by one, and `GET` requests will receive those messages one by one. Multiple `POST`/`GET` requests will stall until they have a place to receive or transmit their data.

In the example, `channel` can be anything, as long as the path starts with `/queue/`.

## Pubsub

Send:

    > curl localhost:8041/sub/channel -d "Hello World!"

Receive:

    > curl localhost:8041/sub/channel
    Hello World!

Similar format, however the main difference is that `POST` requests will never stall, and each `GET` request will simultaneously receive the same data when a `POST` is fired on the channel. This allows for "event listeners" that can receive any event posted to the channel.

## Events

Send:
```
> curl localhost:8041/event/channel?id=msg1&type=message -d "Hello World!"
```

Receive:
```
> curl localhost:8041/event/channel
: connected
id: ID
event: message
data: Hello World!

```

This is a special format used for JavaScript [EventSources](https://developer.mozilla.org/en-US/docs/Web/API/EventSource) to capture events. Each `GET` request never terminates on own, and data is received in the format that the EventSource requires. `POST` requests will send an event to each listener and then terminate. The `id` parameter controls the event id, and the `type` parameter controls the event type (Both are optional). When `id` is not specified it is `undefined` in JS, and `type` is 'message' if not set.

Using in JS:
```JavaScript
let event_type = "message";

let es = new EventSource("localhost:8041/event/channel");
es.addEventListener(event_type, (ev) => {
    console.log("ID: ", ev.lastEventId)
    console.log("Data: ", ev.data);
});
```
Output:
```
ID: msg1
Data: Hello World!
```