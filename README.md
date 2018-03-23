# OTR-CHAT

An end-to-end encrypted chat system based on the OTR protocol.

[简体中文](README-CN.md)


## Language

* Golang ![golang](http://i.imgur.com/UEdZpr4.png)

## Features

- [x] UDP protocol
- [x] Communication with Json format between C & S
- [x] Local configuration with Json format
- [x] end-to-end encryption based on OTR protocol

## OTR protocol

Off-the-Record Messaging (OTR) is a cryptographic protocol that provides encryption for instant messaging conversations. OTR uses a combination of AES symmetric-key algorithm with 128 bits key length, the Diffie–Hellman key exchange with 1536 bits group size, and the SHA-1 hash function. In addition to authentication and encryption, OTR provides forward secrecy and malleable encryption.

For more information, please visit my [blog](http://blog.yfgeek.com/2016/12/06/OTR%E6%8A%80%E6%9C%AF%E6%8E%A2%E8%AE%A8/)

## Snapshots

[](img/otr-0.png)

[](img/otr-1.png)

## Client

The client configuration file located at `~\chat-config.json`. 

Sample
```json
{
	"listen": ":52915",
	"remote": "127.0.0.1",
}
```
The user should set his ID and nickname when running the client.

``
./client
``

## Server 

The client configuration file located at`~\chat-config.json`.

Sample
```json
{
	"listen": ":52915",
	"remote": "", //can be empty
}
```

``
./server
``
