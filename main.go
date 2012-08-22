// The scExplorer server processes commands sent as JSON messages over a
// WebSocket interface. It accepts the following JSON objects as commands:
//
// {cmd: "solve", env: ([]Environment), type: (TempZone)}
//
// {cmd: "plot", envs: ([]Environment), graphParams: (GraphParams)}
//
// The response to "solve" will be a JSON object:
//
// {solved: ([]Environment), errs: ([]string)}
//
// The response to "plot" is a series of messages. The first is a JSON object
// which serves as a directory for the messages to follow. The remaining
// messages are binary data encoding plots as png files. To ensure that
// messages are not received out-of-order, after each message is fully received
// the client must send an acknowledgment to the server (TODO: make sure this
// explicit syncing is needed).
package main

func main() {
}
