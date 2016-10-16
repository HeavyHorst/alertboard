[![Build Status](https://travis-ci.org/HeavyHorst/alertboard.svg?branch=master)](https://travis-ci.org/HeavyHorst/alertboard) [![Go Report Card](https://goreportcard.com/badge/github.com/HeavyHorst/alertboard)](https://goreportcard.com/report/github.com/HeavyHorst/alertboard) [![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/HeavyHorst/alertboard/master/LICENCE)
#alertboard

Alertboard is a very simple Alert-Dashboard with zero dependencies.
It provides a simple RESTful-API to write and read alerts to/from the system.
It's compatiple with [kapacitor](https://github.com/influxdata/kapacitor "kapacitor").

##API
Alertboard has a simple json-RESTful API.
A sample alert would look like:
```json
{
  "id": "cpu/rechner1",
  "message": "hallo",
  "details": "",
  "time": "2016-10-16T20:02:34.79046677+02:00",
  "duration": 0,
  "level": "OK",
  "data": {
    "Series": null,
    "Messages": null,
    "Err": null
  },
  "status": "Open"
}
```

###API-Endpoints

GET ==> /api/alerts

 - get a list of all alerts

GET ==> /api/alerts/{prefix}

- get a list of all alerts with prefix `prefix`


GET ⇒ /api/alert/{alert_id}

- get the alert with the id `alert_id`

DELETE ⇒ /api/alert/{alert_id}

 - delete the alert with the id `alert_id`

GET ⇒ /api/backup

- get a database backup

POST ⇒ /api/alert

- create a new alert

###Example:
Create alert:
```
curl -v -H "Content-Type: application/json" -X POST -d '{"id":"abc123!?$*&()-=@~","message":"hallo", "level":"OK"}' http://localhost:8080/api/alert
...
...
...
* upload completely sent off: 78 out of 78 bytes
< HTTP/1.1 201 Created
< Content-Type: application/json; charset=utf-8
< Location: /api/alert/YWJjMTIzIT8kKiYoKS09QH4=
< Date: Sun, 16 Oct 2016 18:02:34 GMT
< Content-Length: 198
<
* Curl_http_done: called premature == 0
* Connection #0 to host localhost left intact
{"id":"abc123!?$*\u0026()-=@~","message":"hallo","details":"","time":"2016-10-16T20:02:34.79046677+02:00","duration":0,"level":"OK","data":{"Series":null,"Messages":null,"Err":null},"status":"Open"}%
```

You can see the new location of the alert:
The id is simply base64url encoded.

 - Location: /api/alert/YWJjMTIzIT8kKiYoKS09QH4=

We can now get and delete this alert:
```
curl http://localhost:8080/api/alert/YWJjMTIzIT8kKiYoKS09QH4=
curl -X DELETE http://localhost:8080/api/alert/YWJjMTIzIT8kKiYoKS09QH4=
```
