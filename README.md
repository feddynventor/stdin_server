# ruddr

HTTP Server for caching and streaming bytes from `stdin`

A fixed size buffer stores some amount of previous data

The client can easily seek from the `start` byte in the buffer and then immediately syncs to the tail of the stream

`http://:8080/stream?start=0`

### options

GET param `start` refers to the starting byte in the buffer

`const bufferSize = 5 * 1024 * 1024 // 5MB`
