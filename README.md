### Notifier
With go installed you can build it with:

```
make build
```

And you can run it with:

```
./notify -dummyserver --interval=1s --url http://localhost:8080/notify < messages.txt
```

If you prefer you can use docker to run it:

```
docker build -t notifier .

docker run -it notifier /bin/bash -c "./notify -dummyserver --interval=1s --url http://localhost:8080/notify < messages.txt"
```

You can tun test with: `make test` and with the flag `-dummyserver` it will start a test server on port 8080 and receives requests to `/notify` and returns a dummy response containing the request body.

