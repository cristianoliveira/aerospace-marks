# aerospace-client

A Unix socket client for the aerospace window manager.

## Description

This package allows one to interact with the aerospace window manager using a Unix socket. 
Usually located in `/tmp/\(aeroSpaceAppId)-\(unixUserName).sock` ([see](https://github.com/nikitabobko/AeroSpace/blob/f12ee6c9d914f7b561ff7d5c64909882c67061cd/Sources/AppBundle/server.swift#L9))

## Implentation details

As of now the payload for communicating with the socket is a JSON object.

```json
{
    "command": "__deprecated__",
    "args": ["list-windows", "--focused"],
    "stdin": ""
}
```
This is the same as `aerospace list-windows --focused` command. The `stdin` field is not used in this case, but it is included for future extensibility.
