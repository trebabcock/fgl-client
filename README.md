# FGL Database

This is the client for the FGL Database, a terminal based communication platform for the Future Gadget Lab. This started as a joke, but I enjoyed developing it, and my friends enjoyed using it.

You can build the client and [server](https://github.com/trebabcock/fgl-backend) yourself to use with your friends. Although it's not very useful.

# Features

### Current

- User accounts (register/login)
- Announcements
- Lab Reports
- Gadget Reports
- Member Reports
- Interactive shell
- Group messaging

### In Development  

- Private messages
- Multiplayer games

# Gallery   

![Login Screen](/_images/fgl-client1.png "Login Screen")
![Register Screen](/_images/fgl-client2.png "Register Screen")
![Main Menu](/_images/fgl-client3.png "Main Menu")
![Interactive Shell](/_images/fgl-client4.png "Interactive Shell")

# Building

### Requirements

- Go 1.14+
- Windows, macOS, or Linux
	+ Some things may not work properly on Windows cmd. Windows Terminal is recommended. 

### Compiling

```
git clone https://github.com/trebabcock/fgl-client.git
cd fgl-client
go build
```

### Configuration

In order to connect to your server, you must ensure that `config/config.json` contains the correct IP and port.
