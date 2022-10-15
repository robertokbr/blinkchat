### Functional requirements
- User should be able to match with some random user that is searching for some chat as well
- User should be disconnected and connected with another user after a timeout of one minute 
- User should be redirected to the searching chat state if the other user left the chat before the timeout
- User should be able to send and receive messages
- User should be able to send pictures

### Non Functional requirements
- Should use some key value database
- Should use some bucket

### Business rules
- Should not be able to send message after one minute of the chat match
