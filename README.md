# mkw-server

mkw-server is the WIP client server component of MKW-Server. Clients send [Race Packets](https://wiki.tockdom.com/wiki/Network_Protocol#Race_Packets) to mkw-server, which broadcasts them to other clients. mkw-server instances are managed by [wfc-server](https://github.com/mkw-server/wfc-server).

## Current Support

- Transmitting Race Packets between clients

## Setup

1. Set up [wfc-server](https://github.com/mkw-server/wfc-server). Refer to its README for setup instructions.
2. Run `go build` to build the `mkw-server` executable.
3. Add the path to the `mkw-server` executable in [wfc-server's config.xml](https://github.com/mkw-server/wfc-server). This allows wfc-server to locate and launch mkw-server instances.
