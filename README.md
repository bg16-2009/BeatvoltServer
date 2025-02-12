# Beatvolt server
Implementation of the server for a reward system that encourages battery recycling

## Dictionary
- Server -> this app
- Client -> the battery collector that displays codes
- User -> Person recycling and claiming rewards
- Admin -> Person that validates the users claims and gives them the rewards

## Features:
- Provides an API for the clients to authenticate to the server, generate valid prize codes and recieve a random *song
- Provides a web interface for users to claim codes from recycled batteries and get a qr code that can be used by an admin to validate their claimed batteries and give them a reward.
- Provides a web interface for admins to see all users names, emails and how many batteries they recycled
- Provides a web interface for admins to see a list of how many batteries each client recycled
- Provides a web interface for admins to validate a user's QR code and remove batteries for future reward validations

* A song is represented as a json array of numbers. Two numbers are used for every note. The first one represents the note frequency in hertz and the second one represents the duration (1 - full note, 2 - half note, 4 - quarter note etc.). If the duration is negative, the note is dotted(1.5 times the normal duration). If the note frequency is 0, the note is a rest note(silence). If the note frequency is -1 it is treated as an instruction to change the tempo to the value provided in the duration of that note. The default tempo is 140.


## The client API
The protected routes need either a http header: `Authorization: Bearer <token>` or a cookie called `jwt` with the token as a value

### `POST /robot/login` (not protected): 
- Params
    - `username`: the username of the client
    - `password`: the password of the client
- Returns JSON object with a boolean property called success. If success is false, an error will be sent in the `error` property. Otherwise the jwt token that will be used for other requests will be sent in the `token` property
- Gives the client a jwt token for easier further requests 

### `GET /robot/generate_code` (protected): 
- Returns a six character random code that can be used by a user to claim a battery

### `GET /robot/random_song` (protected): 
- Returns a song
