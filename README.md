# DEVICE2AD
## What is this?

[Pumping Station: One](www.pumpingstationone.org) has a number of tools that _really_ shouldn't be operated without the proper authorization/training. To eliminate casual, unuathorized use, an RFID system has been put in place on some tools so that the machine cannot even be powered on without the user being properly authorized to use it.

## Files
### `main.go`
This is the main file that begins the webservice and listens on port 8080.
### `mqtt.go`
Responsible for sending any access attempts to an MQTT topic for logging, etc.
### `adaccess.go`
This file does most of the heavy lifting. In a nutshell, the process of determining whether someone has access is as follows:
1. Get all the users in the specific OU; the name of the OU is stored in the code of whatever device invokes this web service (e.g. the controller that is attached to the LeBlond lathe invokes the service with "LeBlond Lathe")
2. For each user in the group, check the RFID tags (stored in the "otherPager" array) and if there's a match, return 1 to the caller to indicate they're good to. If there is not a match for anyone, return 0 and the machine should not allow it to be used.
