# Controller Endpoints
### Login
The `/login` endpoint allows the user to `POST` their username and email combination in a JSON object to be checked against the user record in the system, if successsful, the JSON response will contain the authenticated user and a JWT the user may use for authorization.
##### POST body
```
{
	"email":"user.name@forcepoint.com",
	"password":"password"
}
```
##### Response body
```
{
    "message": "logged in",
    "status": true,
    "token": "<Json-Web-Token>"",
    "user": {
        "ID": 1,
        "CreatedAt": "2020-06-04T14:38:19Z",
        "UpdatedAt": "2020-06-04T14:38:19Z",
        "DeletedAt": null,
        "Name": "User Name",
        "Email": "user.name@forcepoint.com",
        "Password": "<bcrypt-password>"
    }
}
```
##### Request Authorization
The JWT returned from successful authentication should be added to the `x-access-token` header on each request, if it is not specified you will recieve a `403` error and an error message.

##### All requests require authentication except for `/login`, every other external endpoint is prefixed with `/api` and requires the `x-access-token` header.
### Register (Internal)
The `/register` endpoint allows services to announce themselves to the controller and also to push a list of their endpoints to it so that it may create reverse proxy routes to allow for configuration, pulling service icons, pushing data to the service etc.

This endpoint supports `POST` requests.

Register is an `internal` endpoint and as such uses the `/internal` prefix. These endpoints require the use of the `x-internal-token` header which can be retrieved upon first run of the controller or from the `/api/keys` endpoint (jwt authenticated).
All module-specific internal endpoints (register,update, queue) use internal auth.  

##### POST body
```
 {
        "module_service_name": "fp-dep1",
        "module_display_name": "Forcepoint DEP",
        "module_type": "egress",
        "module_description": "Etiam at sodales urna. Morbi vulputate sollicitudin massa, eget mattis odio vulputate vitae.",
        "inbound_route": "/fpdep",
        "internal_ip": "172.25.0.4",
        "internal_port": "8080",
        "configured": false,
        "internal_endpoints": [
            {
                "endpoint": "/run",
                "http_methods": [
                    {
                        "method": "GET"
                    }
                ]
            },
            {
                "endpoint": "/health",
                "http_methods": [
                    {
                        "method": "GET"
                    }
                ]
            },
            {
                "endpoint": "/icon",
                "http_methods": [
                    {
                        "method": "GET"
                    }
                ]
            },
            {
                "endpoint": "/config",
                "http_methods": [
                    {
                        "method": "GET"
                    },
                    {
                        "method": "POST"
                    }
                ]
            }
        ]
}
```

Every service is required to `POST` this metadata if it wishes to use data from the system.

### Queue (Internal)
The `/queue` endpoint is used by data sources to push blocklist items to so that the controller may create records for them in the database for persistence but also to allow tracking of completeness of export service records.

This endpoint supports `POST` requests.

Queue is an `internal` endpoint and as such uses the `/internal` prefix. These endpoints require the use of the `x-internal-token` header which can be retrieved upon first run of the controller or from the `/api/keys` endpoint (jwt authenticated).
All module-specific internal endpoints (register,update, queue) use internal auth.  

##### POST body
```
{
	"items": [
		{
			"source":"AWS Guard Duty",
			"service_name":"aws-gd1",
			"type":"Domain",
			"value":"china.com.cn"
		},
		{
			"source":"AWS Guard Duty",
			"service_name":"aws-gd1",
			"type":"URL",
			"value":"www.china.com.cn/malware"
		}
		]
}
```
##### Response Body
The response from the `POST` is basically the same as the data sent but with the `batch_number` added to show that the data was persisted.
```
{
    "items": [
        {
            "source": "AWS Guard Duty",
            "service_name": "aws-gd1",
            "type": "Domain",
            "value": "china.com.cn",
            "batch_number": 2
        },
        {
            "source": "AWS Guard Duty",
            "service_name": "aws-gd1",
            "type": "URL",
            "value": "www.china.com.cn/malware",
            "batch_number": 2
        }
    ]
}
```


### Update (Internal)
The `/update` endpoint is used by the connected services to signal that they have successfully processed a batch of blocklist items from one of the sources, this is then persisted by the controller as an audit record and also to make sure that all consumers of the blocklist data are completely up to date.

This endpoint supports `POST` requests.

Update is an `internal` endpoint and as such uses the `/internal` prefix. These endpoints require the use of the `x-internal-token` header which can be retrieved upon first run of the controller or from the `/api/keys` endpoint (jwt authenticated).
All module-specific internal endpoints (register,update, queue) use internal auth.  

##### POST body

```
{
	"status":"success",
	"service_name":"fp-fba1",
	"update_batch_id":3
}
```
The possible statuses are: `success` and `failed`

The `batch_id` refers to a record in the table created when the data source pushed its updates to the system to allow for tracking how up to date services are.
### Stats
The `/stats` endpoint is a `GET` request to return the blocklist statistics for the current installation. You can retrieve the number of separate sources for the blocklist and also a breakdown of the numbers of each blocked type.

The last time the blocklist was updated is also returned here.

The default behaviour for this endpoint is to return global statistics, but to get granular results you can add the `servicename` keyword as a query parameter to get statistics for a single service.
```
{
    "num_sources": 0,
    "num_blocked_ip": 0,
    "num_blocked_domains": 0,
    "num_blocked_urls": 0,
    "last_update": "2020-06-04 17:44:48 +0000 UTC"
}
```
### Health
The `/health` endpoint supports `GET` requests and returns health and status information for the controller and its MariadDB container.
```
{
    "modules": [
        {
            "module_name": "master-controller",
            "status": 1,
            "status_code": 200,
        },
        {
            "module_name": "master-database",
            "status": -1,
            "status_code": 200,
        }
    ]
}
```

The potential statuses that could be returned are:
```
Down = -1 
Unhealthy = 0
Healthy = 1
```
### Elements
The `/elements` endpoint allows for the searching and viewing of the blocklist data. 
It allows for fuzzy searching using the `searchTerm` query parameter.
This endpoint also allows for `edit` and `delete` functions via `PUT` and `DELETE` methods.
```
{
    "results": [
        {
            "source": "AWS Guard Duty",
            "service_name": "aws-gd1",
            "type": "Domain",
            "value": "china.com.cn",
            "batch_number": 1
        },
        {
            "source": "AWS Guard Duty",
            "service_name": "aws-gd1",
            "type": "URL",
            "value": "www.china.com.cn/malware",
            "batch_number": 1
        }...
```
### Export
The `/export` endpoint allows for the export of the blocklist data in different formats with the default being JSON.
The required format can be specified by using the `format` keyword as a query paramter and choosing from the following formats: `text`, `csv`, and `json`.
```
{
    "results": [
        {
            "source": "AWS Guard Duty",
            "service_name": "aws-gd1",
            "type": "Domain",
            "value": "china.com.cn",
            "batch_number": 1
        },
        {
            "source": "Microsoft Sentinel",
            "service_name": "ms-snt1",
            "type": "IP",
            "value": "1.2.3.4",
            "batch_number": 1
        },
        {
            "source": "AWS Guard Duty",
            "service_name": "aws-gd1",
            "type": "URL",
            "value": "www.china.com.cn/malware",
            "batch_number": 1
        }...
```
### Logs
The `/logs` endpoint allows for the extracting of system logs from the database for displaying in the UI. 

The logs can be filtered by level or service name.

Results from this endpoint are paged as they could become large over time, therefore a `page` value must be added to the query.

The logs can be filtered by using these keywords as query parameters: `level`, and `modulename`.
The acceptable values for level are: `trace`, `debug`, `info`, `warning`, `error`, `fatal`, and `panic`.

N.B. When a level is specified, the endpoint will return everything from that level and up.

```
{
    "events": [
        {
            "module_name": "master-controller",
            "level": "info",
            "message": "Adding module routes from persistence...",
            "caller": "github.com/sirupsen/logrus.(*Logger).Log",
            "time": "2020-06-06T17:37:59Z"
        },
        {
            "module_name": "master-controller",
            "level": "info",
            "message": "Adding new module: Forcepoint DEP",
            "caller": "github.com/sirupsen/logrus.(*Logger).Log",
            "time": "2020-06-06T17:36:14Z"
        },
        {
            "module_name": "master-controller",
            "level": "info",
            "message": "Adding new module: Forcepoint DUP",
            "caller": "github.com/sirupsen/logrus.(*Logger).Log",
            "time": "2020-06-06T17:36:15Z"
        },
        {
            "module_name": "master-controller",
            "level": "info",
            "message": "Adding new module: Microsoft Sentinel",
            "caller": "github.com/sirupsen/logrus.(*Logger).Log",
            "time": "2020-06-06T17:36:15Z"
        },
        {
            "module_name": "master-controller",
            "level": "info",
            "message": "Adding new module: AWS Guard Duty",
            "caller": "github.com/sirupsen/logrus.(*Logger).Log",
            "time": "2020-06-06T17:36:15Z"
        }...
```
### Modules
The `/modules` endpoint supports `GET` requests and returns information related to connected modules/services.

There are two types of service, `ingress` and `egress`, the results from this endpoint can be filtered by type by specifying `moduletype` as a query paramter.
```
[
    {
        "module_service_name": "fp-ngfw1",
        "module_display_name": "Forcepoint NGFW",
        "module_type": "egress",
        "module_description": "Etiam at sodales urna. Morbi vulputate sollicitudin massa, eget mattis odio vulputate vitae.",
        "inbound_route": "/fpngfw",
        "internal_ip": "172.25.0.4",
        "internal_port": "8080",
        "configured": false,
        "internal_endpoints": [
            {
                "endpoint": "/run",
                "http_methods": [
                    {
                        "method": "GET"
                    }
                ]
            },
            {
                "endpoint": "/health",
                "http_methods": [
                    {
                        "method": "GET"
                    }
                ]
            },
            {
                "endpoint": "/icon",
                "http_methods": [
                    {
                        "method": "GET"
                    }
                ]
            },
            {
                "endpoint": "/config",
                "http_methods": [
                    {
                        "method": "GET"
                    },
                    {
                        "method": "POST"
                    }
                ]
            }
        ],
        "module_health": {
            "module_name": "Forcepoint NGFW",
            "status": 1,
            "status_code": 200,
            "last_update": ""
        }
    },
```

