The protocol consists of two parts:
- Code
- Message

The code is an uint of 4 bytes (32 bits) that represents the type of message.

After a code is sent, the behaviour is set, the code with their messages are:

### ECHO_MESSAGE
- Code: 1
After the code is sent, then another 4 bytes are sent, which are the length of the message, finally the message is sent.

The server will respond with the same message.
