# Chat with GossipSub

## Build

```shell
go build -o chat
```

## Run

Open three terminal, in the 1st, run:

```shell
./chat
2023-05-26T10:25:00.337+0800	INFO	chat-with-gossipsub/main.go:54	host created	{"desc": ["/ip4/192.168.64.1/tcp/49155/p2p/12D3KooWDjSU1KJCZ1repqeofgGG9NDeDZdxX4HePTLDG2Zm9wLt", "/ip4/127.0.0.1/tcp/49155/p2p/12D3KooWDjSU1KJCZ1repqeofgGG9NDeDZdxX4HePTLDG2Zm9wLt"]}
```

In the 2nd and 3rd, run:

```shell
./chat -peers /ip4/127.0.0.1/tcp/49155/p2p/12D3KooWDjSU1KJCZ1repqeofgGG9NDeDZdxX4HePTLDG2Zm9wLt
```
