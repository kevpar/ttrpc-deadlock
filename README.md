Run server:

```
./ttrpc-deadlock.exe -conn-type npipe server
```

Run client:

```
./ttrpc-deadlock.exe -conn-type npipe -workers 4 client
```

Client will repeatedly send requests to the server and print out
the worker goroutine index as each one completes. At some point
you should see the output stop as the client/server will deadlock.

Increasing the client worker count will make the deadlock happen sooner.