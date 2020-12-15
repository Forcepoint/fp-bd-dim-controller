# Module Endpoints
All service endpoints are prefixed by the servcies specific `inbound_route` which is specified in the metadata the service uses to register with the controller. 

An example would be a service for pushing updates to the Forcepoint NGFW with an `inbound_route` of `/fpngfw`, this means that to access the config endpoint of that particular module the path would be `/api/fpngfw/config`
### Config
#### This is a required endpoint.
The `/config` endpoint is used for inspecting the config of a modules if there is one, updating the config and also for pulling the template for the dynamic UI which allows for creating the UI for the configuration page dynamically for each service.

The config endpoint should support both `POST` and `GET` requests, an exception being a service that has no config, then the `POST` method could be removed from the metadata.

The response to be `POST`-ed to the `/config` endpoint is created from the metadata received from the `GET` request to `/config`


```
const (
	Text ElementType = iota + 1
	Select
	Radio
	Number
	Password
	Disabled
	Info
)

{

    // To be retrieved
    "fields": [
        {
            "label":"text",
            "type": 1
            "expected_json_name":"aws_account_number",
            "rationale":"Text describing what this field is for"
            "value":"text",
            "possible_values": [
                "value1",
                "value2",
                "value3"
            ],
            "required": true
        },
    ],
    
    // To be sent
    "values": {
        "aws_account_number": "123456789",
        ....
    }
}
```
### Icon
#### This is a required endpoint.
The `/icon` endpoint should return a 200x200 `png`.

This endpoint only supports `GET` requests.

### Health
#### This is a required endpoint.
The `/health` endpoint is just used by the controller to hit the service and check the http status code.

This endpoint only supports `GET` requests.

### Run
#### This is a required endpoint.
The `/run` endpoint is what is used to run the functionality of a particular service, i.e. hitting `/run` with a `POST` request with blocklist items for the `fp-smc` module will start the process of pushing them to the NGFW.

This endpoint can support `POST` requests

## Data Structure
### Module Metadata
```
type ModuleMetadata struct {
	ModuleServiceName string           `json:"module_service_name"`
	ModuleDisplayName string           `json:"module_display_name"`
	ModuleType        string           `json:"module_type"`
	ModuleDescription string           `json:"module_description"`
	InboundRoute      string           `json:"inbound_route"`
	InternalIP        string           `json:"internal_ip"`
	InternalPort      string           `json:"internal_port"`
	Configured        bool             `json:"configured"`
	InternalEndpoints []ModuleEndpoint `json:"internal_endpoints"`
}

type ModuleEndpoint struct {
	Secure      bool         `json:"secure"`
	Endpoint    string       `json:"endpoint"`
	HttpMethods []HttpMethod `json:"http_methods"`
}

type HttpMethod struct {
	Method string `json:"method"`
}
```
### Config Structs
```
type ElementType int

const (
	Text ElementType = iota + 1
	Select
	Radio
	Number
	Password
	Disabled
	Info
)

type ModuleConfig struct {
	Fields []Element `json:"fields"`
}

type Element struct {
	Label            string      `json:"label"`
	Type             ElementType `json:"type"`
	ExpectedJsonName string      `json:"expected_json_name"`
	Rationale        string      `json:"rationale"`
	Value            string      `json:"value"`
	PossibleValues   []string    `json:"possible_values"`
	Required         bool        `json:"required"`
}
```
### Log Event
```
type LogEvent struct {
	ModuleName string    `json:"module_name"`
	Level      string    `json:"level"`
	Message    string    `json:"message"`
	Caller     string    `json:"caller"`
	Time       time.Time `json:"time"`
}
```
### Intelligence Item
These are what are pushed to and received from the controller by the modules.  

```
type ProcessedItem struct {
	Source        string `json:"source"`
	ServiceName   string `json:"service_name"`
	Type          string `json:"type"`
	Value         string `json:"value"`
	UpdateBatchId uint   `json:"batch_number"`
}
```
They are sent as `Lists`.  

```
type ProcessedItemsWrapper struct {
	Items []ProcessedItem `json:"items"`
}
```