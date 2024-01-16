# go-counter-server

A simple server that takes a request and returns the number of requests in the last 60 seconds to the client


## DEMO

Here is a short demo of the program, here we can see that the request count is persisted even after terminating the program and spinning it up again,
this is achieved by a persistence file "data.txt" that is updated in disk with each request   

https://github.com/HenBK/go-counter-server/assets/42653917/b62e7a1a-4ef1-4683-a8cc-ca24de67325f
