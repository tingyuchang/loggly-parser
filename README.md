# loggly-parser
convenient tool to parse Loggly's JOSN file

We want to analysis the below JSON log files (from loggly)

```
{
  "total_events": 1,
  "page": 0,
  "events": [
    {
      "id": "6839bf3f-e7ae-11ec-8033-12c72a2782d6",
      "timestamp": 1654749680058,
      "event": {
        "json": {
          "duration": 307.40000000596046,
          "method": "GET",
          "size": 2,
          "resource": "/profiles/mattchang19?expand=%7B%22currentCompany%22%3A%7B%22_fields%22%3A%5B%22name%22%2C%22identifier%22%5D%7D%7D",
          "tag": "response-time",
          "user": "matt@t1.co"
        },
        "http": {
          "clientHost": "1.160.142.28",
          "contentType": "text/plain;charset=UTF-8"
        }
      },
      "logtypes": [
        "json"
      ],
      "tags": [
        "http"
      ]
    }
  ]
}
```

## How to use it

this program needs input PATH to file as arguement

example:

```
go run main.go ~/PATH_TO_JSON_FILE
```

or 

```
go build .
./loggly-parser  ~/PATH_TO_JSON_FILE
```

