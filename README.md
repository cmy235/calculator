
# Golang web app abacus
#### Connor Yanz

## Code
##### The server first serves static files at localhost:8080 which gives the user the option to input the operations via a form. The server handles 4 routes for addition, subtraction, multiplication & division and each function routes to the compute function, passing the appropriate action.
##### When the main function runs a datastore is instantiated, and the timer on the datastore begins ticking. Every second this timer checks for operations that have been in the cache for longer than 1 minute.
##### The sync package was utilized to put mutual exclusion locks on the cache while clearing expired values.  This is good practice to prevent situations where deletions or writes to the cache are happening simultaneously.  The mutex ensures that for the brief period where the cache is locked, the only changes that are happening are the ones you are making.
##### The result of the operation is written to the response in JSON, and displays a boolean based on whether or not that operation was previously cached.

##### To run the web app execute the binary and navigate to localhost:8080
##### ``` $ ./main ```

## Testing
##### Test files are currently being refactored...
##### The app can be tested via curl or wget. To run unit tests individually
##### ``` $ go test -v main_test.go main.go operations.go cache.go ```
##### or run the test binary
##### ``` $ ./go_abacus.test ```


