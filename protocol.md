For now, the protocol after a connection is made is:
- client:
1. The client sends the size of the message in 4 bytes
2. The client sends the message
- server:
1. The server reads the size of the message
2. The server reads the message