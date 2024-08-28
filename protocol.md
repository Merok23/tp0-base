The protocol consists of two parts:
- Code
- Message

The code is an uint of 4 bytes (32 bits) that represents the type of message.

After a code is sent, the behaviour is set, the code with their messages are:

### ECHO_MESSAGE
- Code: 1
After the code is sent, then another 4 bytes are sent, which are the length of the message, finally the message is sent.

The server will respond with the same message.


### BET
- Code: 2
After the code, the client will send:
- DNI (4 bytes)
- NUMERO (4 bytes)
- NACIMIENTO (10 bytes, a string with the format "YYYY-MM-DD")
- NAME_SIZE (4 bytes)
- NAME (NAME_SIZE bytes)
- LAST_NAME_SIZE (4 bytes)
- LAST_NAME (LAST_NAME_SIZE bytes)

The server will try to store the bet, if it does, it will respond with a success code, otherwise it will respond with an error code.

### SUCCESS
- Code: 200

### ERROR
- Code: 400