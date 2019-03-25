# WhatDash
Experimental application Whatsapp messaging based on server-side communication bridge.

<img src="https://i.imgur.com/6M0knpk.png" align="right">

This is an experimental project which covering service of whatsapp messaging, so do the development with a caution.

## Web GUI
### Homepage Preview
<img src="https://i.imgur.com/xN0d6vt.png" align="center">

### Register new Number Preview
<img src="https://i.imgur.com/JMP0hsS.png" align="center">

### Chat Window Preview
<img src="https://i.imgur.com/do5wxgK.png" align="center">

## API Documentation
### Registering New Number
Register a new number before starting to communicate by this service

**Request**: POST `/wa/session/create`

- `number` - the string phone number with area code without "+" eg: 6285716116666

```
  {
    "number": "6285716116666" <string|required>
  }
```

### Checking Number
Checking the phone number if has been registered to the service

**Request**: POST `/wa/session/check`

- `number` - the string phone number with area code without "+" eg: 6285716116666

```
  {
    "number": "6285716116666" <string|required>
  }
```

### Send Text Message
Sending text message

**Request**: POST `/wa/send/text`

- `from` the sender phone number
- `to` receiver phone number
- `message` the string message to send

```
  {
    "from": "6285716116666", <string|required>
    "to": "6285716117777", <string|required> 
    "message": "Hello from whatsapp"
  }
```

### Terminate Socket Connection
Terminate existing socket connection

**Request**: POST `/wa/connection/terminate`

- `number` the phone number to terminate

```
  {
    "number": "6285716116666" <string|required>
  }
```