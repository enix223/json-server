# JSON server written by golang

This is a simple http server to generate static json response base on user config.

## 1. Usage

Start the json server:

```shell
json-server -a localhost:7000 -e ./example/demo.json
```

## 2. `json-server` arguments

```
-a  Json server listen address
-e  Handler mapping entry list
```

## 3. Handler mapping format

| Field | Type | Required | Desc | Example value |
|-|-|-|-|-|
| path | string | √ | http request path | `/v1/hello` |
| method | string |  | http method | `GET`, `POST`, `DELETE`, `PATCH`, `PUT`, default is `GET` |
| payload | `object/string` |  | http response body, if `payload_type = json`, then `payload` shoulde be in `object` format, otherwise `string` format | `{"name": "foo"}` |
| payload_type | string | | payload serialize type, possible value: [`json`, `text`] | `json` |
| resposne_headers | object | | http response header | `{"content_type": "application/json"}` |
| status_code | number | | http response status code | 200 |
