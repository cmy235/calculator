
# Teltech coding challenge
# Connor Yanz

## Code
### The server first serves static files at localhost:8080 which gives the user the option to input the operations via the gui. 4 routes are available routes for addition, subtraction, multiplication & division and each routes to the compute function passing the appropriate action.
### When the main function runs a datastore is instantiated, and the timer on the datastore begins ticking. Every second this timer checks for operations that have been in the cache for longer than 1 minute.

### The sync package was utilized to put mutual exclusion locks on the key/value store while clearing expired values.  This is good practice to prevent any situation where deletions or writes were happening simultaneously.  The mutex ensures that for the brief period where the cache is locked, the only changes that are happening are the ones you are making.

### To run the web app execute the binary
### $ ./main

## Testing
### The app can be tested via curl or wget. To run all unit tests:
### $ go test -v main_test.go main.go operations.go cache.go

