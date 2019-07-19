# 4sq-exports

Simple CLI-tool to export your Foursquare checkins into a KML-file.



## Usage

1. First you need to run the authorise command to get an access token by logging into your Foursquare account. The app listens on localhost:12345/4sq for the callback so you need to allow it in your firewall.
```bash
$ ./4sq-exports authorise
Your Foursquare authentication detais. Please make a note of them. You will need them for other commands.
- Authorisation code: XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
- Access token: XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
```

2. Then you use the access token with the checkins command to get your checkins list from the API.
```bash
$ ./4sq-exports checkins --accessToken XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX --output outputfile.kml
Output file: outputfile.kml
- Total number of check-ins: 24894
- 100 pages of 250 check-ins
Getting check-ins 0 to 250 (page 0 of 100)
...
Getting check-ins 24750 to 25000 (page 99 of 100)
```
