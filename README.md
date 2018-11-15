# Usage

make and run 

## Example 

```bash

 curl -X POST -H "Content-Type: application/json" -d '{"type": "confirmation"}' localhost:9911/ 


 Note: Unnecessary use of -X or --request, POST is already inferred.
*   Trying 127.0.0.1...
* TCP_NODELAY set
* Connected to localhost (127.0.0.1) port 9911 (#0)
> POST / HTTP/1.1
> Host: localhost:9911
> User-Agent: curl/7.58.0
> Accept: */*
> Content-Type: application/json
> Content-Length: 24
> 
* upload completely sent off: 24 out of 24 bytes
< HTTP/1.1 200 OK
< Content-Type: text/plain; charset=UTF-8
< Date: Thu, 15 Nov 2018 17:54:20 GMT
< Content-Length: 5
< 
* Connection #0 to host localhost left intact
coded

```