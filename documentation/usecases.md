### onConnect: receives the connection id and user data in req query 
 create user if not exists
 if user exists, update connection id
 set user state to connected

###  onSearch: receives the connection id 
 set user state to searching
 if user is in room, leave room
 search for users with state searching
 request to join room with first user found
 if user accepts, this user will emit onMatch
 if user rejects, return message to frontend to search again

### onMatch: receives user id and the id of the user who has triggered onSearch 
 search for the existence of some room with the user id and the user who has triggered onSearch and state open
 create room with user and matched user
 send to both users the room data

### onMessage: receives the connection id and the message data 
 search for the user with the connection id
 search for the room with the user id
 verify if the room created_at is less than 1 minute, if not, close the room and return false
 if the room is not open, return false
 create message with the room id, user id, message content and message type
 send to all users in the room the message data