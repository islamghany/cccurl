# cccurl

**cccurl** is a lightweight command-line tool built in Go, designed to mimic the basic functionalities of the widely-used `curl` utility. It allows users to send HTTP requests with customizable methods, headers, and payloads directly from the terminal.

## Features

- **Custom HTTP Methods:** Specify HTTP methods such as GET, POST, DELETE, etc., using the `-X` flag.
- **Custom Headers:** Add one or multiple HTTP headers to your requests using the `-H` flag.
- **Data Payloads:** Send data with your requests (e.g., JSON payloads) using the `-d` flag.
- **Automatic Header Management:** Automatically handles essential headers like `Content-Length` and defaults `Content-Type` when sending data.
- **Simple and Intuitive:** Designed for ease of use with clear command-line options.

## Installation

### Prerequisites

- **Go:** Ensure that you have Go installed on your system. You can download it from [https://golang.org/dl/](https://golang.org/dl/).

### Building from Source

1. **Clone the Repository:**

   ```bash
   git clone https://github.com/islamghany/cccurl.git
   cd cccurl
   ```

2. **Build the Executable:**

   ```bash
   go build -o cccurl main.go
   ```

   This command compiles the Go source code and generates an executable named `cccurl` in the current directory.

3. **Move to a Directory in Your PATH (Optional):**

   To use `cccurl` from anywhere in your terminal, move the executable to a directory that's included in your system's `PATH`, such as `/usr/local/bin`.

   ```bash
   sudo mv cccurl /usr/local/bin/
   ```

## Usage

The basic syntax for using `cccurl` is as follows:

```bash
cccurl [options] <URL>
```

### Options

- `-X <method>`: Specify the HTTP method to use (e.g., GET, POST, DELETE). Defaults to `GET` if not provided.
- `-d <data>`: Send data payload with the request. Commonly used with POST requests to send JSON or form data.
- `-H "<Header>: <Value>"`: Add a custom HTTP header to the request. This option can be used multiple times to include multiple headers.

### Examples

#### 1. Sending a GET Request (Default Method)

```bash
cccurl http://eu.httpbin.org/get
```

**Output:**

```
Connecting to eu.httpbin.org
Sending request GET /get HTTP/1.1
Host: eu.httpbin.org
Accept: */*
Connection: close

HTTP/1.1 200 OK
Date: Fri, 15 Dec 2023 14:29:23 GMT
Content-Type: application/json
Content-Length: 227
Connection: close
Server: gunicorn/19.9.0
Access-Control-Allow-Origin: *
Access-Control-Allow-Credentials: true

{
  "args": {},
  "headers": {
    "Accept": "*/*",
    "Host": "eu.httpbin.org",
    "X-Amzn-Trace-Id": "Root=1-657c62c3-26068fd12f977c810ce87090"
  },
  "url": "http://eu.httpbin.org/get"
}
```

#### 2. Sending a DELETE Request

```bash
cccurl -X DELETE http://eu.httpbin.org/delete
```

**Output:**

```
Connecting to eu.httpbin.org
Sending request DELETE /delete HTTP/1.1
Host: eu.httpbin.org
Accept: */*
Connection: close

HTTP/1.1 200 OK
Date: Fri, 15 Dec 2023 14:29:23 GMT
Content-Type: application/json
Content-Length: 227
Connection: close
Server: gunicorn/19.9.0
Access-Control-Allow-Origin: *
Access-Control-Allow-Credentials: true

{
  "args": {},
  "data": "",
  "files": {},
  "form": {},
  "headers": {
    "Accept": "*/*",
    "Host": "eu.httpbin.org",
    "X-Amzn-Trace-Id": "Root=1-657c68d7-7b7a96900d27d3a952f99f65"
  },
  "json": null,
  "url": "http://eu.httpbin.org/delete"
}
```

#### 3. Sending a POST Request with JSON Payload

```bash
cccurl -X POST http://eu.httpbin.org/post \
  -d '{"key": "value"}' \
  -H "Content-Type: application/json"
```

**Output:**

```
Connecting to eu.httpbin.org
Sending request POST /post HTTP/1.1
Host: eu.httpbin.org
Accept: */*
Connection: close
Content-Type: application/json
Content-Length: 16

HTTP/1.1 200 OK
Date: Fri, 15 Dec 2023 14:29:23 GMT
Content-Type: application/json
Content-Length: 227
Connection: close
Server: gunicorn/19.9.0
Access-Control-Allow-Origin: *
Access-Control-Allow-Credentials: true

{
  "args": {},
  "data": "{\"key\": \"value\"}",
  "files": {},
  "form": {},
  "headers": {
    "Accept": "*/*",
    "Content-Length": "16",
    "Content-Type": "application/json",
    "Host": "eu.httpbin.org",
    "X-Amzn-Trace-Id": "Root=1-657c69ae-6ea3b1ea7084a25843f4814c"
  },
  "json": {
    "key": "value"
  },
  "url": "http://eu.httpbin.org/post"
}
```

#### 4. Sending a GET Request with Custom Headers

```bash
cccurl -H "User-Agent: cccurl/1.0" http://eu.httpbin.org/get
```

**Output:**

```
Connecting to eu.httpbin.org
Sending request GET /get HTTP/1.1
Host: eu.httpbin.org
Accept: */*
Connection: close
User-Agent: cccurl/1.0

HTTP/1.1 200 OK
Date: Fri, 15 Dec 2023 14:29:23 GMT
Content-Type: application/json
Content-Length: 227
Connection: close
Server: gunicorn/19.9.0
Access-Control-Allow-Origin: *
Access-Control-Allow-Credentials: true

{
  "args": {},
  "headers": {
    "Accept": "*/*",
    "Host": "eu.httpbin.org",
    "User-Agent": "cccurl/1.0",
    "X-Amzn-Trace-Id": "Root=1-657c62c3-26068fd12f977c810ce87090"
  },
  "url": "http://eu.httpbin.org/get"
}
```

## Error Handling

- **Invalid Header Format:**

  If a header is not in the correct `Key: Value` format, `cccurl` will display an error message.

  ```bash
  cccurl -X POST http://eu.httpbin.org/post -H "InvalidHeaderFormat"
  ```

  **Output:**

  ```
  invalid header format: InvalidHeaderFormat. Expected 'Key: Value'
  ```

- **Missing URL Argument:**

  If no URL is provided, `cccurl` will display usage instructions.

  ```bash
  cccurl -X DELETE -H "Content-Type: application/json"
  ```

  **Output:**

  ```
  Usage: cccurl [options] <URL>
    -H value
          HTTP header
    -X string
          HTTP method (default "GET")
    -d string
          HTTP payload
  error: exactly one URL must be provided
  ```

- **Unsupported Protocol:**

  Currently, `cccurl` only supports the HTTP protocol. Attempting to use HTTPS will result in an error.

  ```bash
  cccurl -X POST https://eu.httpbin.org/post -d '{"key": "value"}' -H "Content-Type: application/json"
  ```

  **Output:**

  ```
  Error: Only HTTP protocol is supported
  ```
