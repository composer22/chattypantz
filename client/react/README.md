## ChattyPantz Client Demo - React.js

This folder contains a simple demonstration of a client connection to the server using javascript and html.

### Dependencies

* You must be connected to the internet to load CDN react.js, semnatic-ui.js, and jquery.js.
* The server should be running on localhost:6660

### Instructions

Assumption: server is up and running on localhost:6660.

If needed, you can change the API host and port endpoint in scripts/services/api.js:
```
var server = "ws://127.0.0.1:6660/v1.0/chat";
```
1. Load index.html in your browser.
2. Enter a nickname.
3. Connect to the server.

You will be placed into room "Demo". A list of user nicknames from the room will also be displayed.
NOTE:  If your nickname already exists when logging into the room, you will be warned and disconnected.

4. Now, send your messages.
5. When you are done, disconnect from the server.

### Enhancements

The application doesn't demonstrate multi-room management, nor does it demonstrate the following request types:
* GET_NICKNAME: 102
* LIST_ROOMS: 103
* HIDE: 106
* UNHIDE: 107
