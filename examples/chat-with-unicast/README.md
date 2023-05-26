# Chat with Unicast

## Build

```shell
go build -o chat
```

## Run

Open three terminal, in the 1st, run:

```shell
./chat
2023-05-26T10:30:21.185+0800	INFO	chat-with-unicast/main.go:55	host created	{"desc": ["/ip4/192.168.64.1/tcp/49266/p2p/12D3KooWLE1UPeYKzAWpkajMc7GASXPVMB3p9PPgnuDbSMFtWrxS", "/ip4/127.0.0.1/tcp/49266/p2p/12D3KooWLE1UPeYKzAWpkajMc7GASXPVMB3p9PPgnuDbSMFtWrxS"]}
2023-05-26T10:30:21.185+0800	INFO	chat-with-unicast/main.go:57	run as bootstrap node
```

In the 2nd, run:

```shell
./chat -boots /ip4/127.0.0.1/tcp/49266/p2p/12D3KooWLE1UPeYKzAWpkajMc7GASXPVMB3p9PPgnuDbSMFtWrxS
2023-05-26T10:30:41.989+0800	INFO	DHT	go-easy-p2p/dht.go:58	connected to bootstrap peer	{"peer": "{12D3KooWLE1UPeYKzAWpkajMc7GASXPVMB3p9PPgnuDbSMFtWrxS: [/ip4/127.0.0.1/tcp/49266]}"}
2023-05-26T10:30:41.990+0800	INFO	chat-with-unicast/main.go:55	host created	{"desc": ["/ip4/192.168.64.1/tcp/49267/p2p/12D3KooWExXqHPT5t3A8TTLc5dr5EWEMtj6ECRpyMne5KbQRjLDm", "/ip4/127.0.0.1/tcp/49267/p2p/12D3KooWExXqHPT5t3A8TTLc5dr5EWEMtj6ECRpyMne5KbQRjLDm"]}
2023-05-26T10:30:41.990+0800	INFO	chat-with-unicast/main.go:73	run as receiver node
```

In the 3rd, run:

```shell
./chat -boots /ip4/127.0.0.1/tcp/49266/p2p/12D3KooWLE1UPeYKzAWpkajMc7GASXPVMB3p9PPgnuDbSMFtWrxS -dst 12D3KooWExXqHPT5t3A8TTLc5dr5EWEMtj6ECRpyMne5KbQRjLDm
2023-05-26T10:31:43.428+0800	INFO	DHT	go-easy-p2p/dht.go:58	connected to bootstrap peer	{"peer": "{12D3KooWLE1UPeYKzAWpkajMc7GASXPVMB3p9PPgnuDbSMFtWrxS: [/ip4/127.0.0.1/tcp/49266]}"}
2023-05-26T10:31:43.428+0800	INFO	chat-with-unicast/main.go:55	host created	{"desc": ["/ip4/192.168.64.1/tcp/49292/p2p/12D3KooWPUNCwJj7yEBUCGHp3AEinfj7wC4XgeJFuKeJXwrChPNp", "/ip4/127.0.0.1/tcp/49292/p2p/12D3KooWPUNCwJj7yEBUCGHp3AEinfj7wC4XgeJFuKeJXwrChPNp"]}
2023-05-26T10:31:43.428+0800	INFO	chat-with-unicast/main.go:77	run as sender node
2023-05-26T10:31:43.428+0800	INFO	chat-with-unicast/main.go:83	peers connected before send	{"peers": ["12D3KooWLE1UPeYKzAWpkajMc7GASXPVMB3p9PPgnuDbSMFtWrxS"]}
```