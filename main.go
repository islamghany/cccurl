package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
)

// urlOptions holds the parsed components of a URL
type urlOptions struct {
	Protocol string
	Host     string
	Port     string
	Path     string
	Query    string
	Fragment string
}

// parseURL parses the input URL string and returns its components
func parseURL(urlstr string) (urlOptions, error) {
	parsedURL, err := url.Parse(urlstr)
	if err != nil {
		return urlOptions{}, err
	}

	port := parsedURL.Port()
	if port == "" {
		if parsedURL.Scheme == "http" {
			port = "80"
		} else if parsedURL.Scheme == "https" {
			port = "443"
		}
	}

	path := parsedURL.Path
	if path == "" {
		path = "/"
	}
	if parsedURL.RawQuery != "" {
		path += "?" + parsedURL.RawQuery
	}

	return urlOptions{
		Protocol: parsedURL.Scheme,
		Host:     strings.Split(parsedURL.Host, ":")[0],
		Port:     port,
		Path:     path,
		Fragment: parsedURL.Fragment,
	}, nil
}

// headerList is a custom flag type to allow multiple -H flags
type headerList []string

// String returns the string representation of the headerList
func (h *headerList) String() string {
	return strings.Join(*h, ", ")
}

// Set appends a new header to the headerList
func (h *headerList) Set(value string) error {
	*h = append(*h, value)
	return nil
}

// requestOptions holds all the configurations for the HTTP request
type requestOptions struct {
	Method  string
	Data    string
	Headers headerList
	URL     string
}

// parseFlags parses and validates the command-line flags and arguments
func parseFlags() (requestOptions, error) {
	var opts requestOptions

	// Define command-line flags
	flag.StringVar(&opts.Method, "X", "GET", "HTTP method")
	flag.StringVar(&opts.Data, "d", "", "HTTP payload")
	flag.Var(&opts.Headers, "H", "HTTP header")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options] <URL>\n", os.Args[0])
		flag.PrintDefaults()
	}

	// Parse flags
	flag.Parse()

	// Ensure that exactly one positional argument (the URL) is provided
	if flag.NArg() != 1 {
		flag.Usage()
		return opts, fmt.Errorf("error: exactly one URL must be provided")
	}

	opts.URL = flag.Arg(0)
	opts.Method = strings.ToUpper(opts.Method)

	return opts, nil
}

// buildHeaders constructs the headers map, incorporating default and user-provided headers
func buildHeaders(options urlOptions, userHeaders headerList, data string) (map[string]string, error) {
	headersMap := make(map[string]string)

	// Set default headers
	headersMap["Host"] = options.Host
	headersMap["Accept"] = "*/*"
	headersMap["Connection"] = "close"

	// Parse and add user-provided headers
	for _, header := range userHeaders {
		parts := strings.SplitN(header, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid header format: %s. Expected 'Key: Value'", header)
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		headersMap[key] = value
	}

	// If data is provided, set the Content-Length header
	if data != "" {
		headersMap["Content-Length"] = fmt.Sprintf("%d", len(data))
		// If Content-Type is not set, default to application/x-www-form-urlencoded
		if _, exists := headersMap["Content-Type"]; !exists {
			headersMap["Content-Type"] = "application/x-www-form-urlencoded"
		}
	}

	return headersMap, nil
}

// constructHTTPRequest builds the full HTTP request string
func constructHTTPRequest(method string, path string, headers map[string]string, body string) string {
	var requestBuilder strings.Builder

	// Request line
	requestBuilder.WriteString(fmt.Sprintf("%s %s HTTP/1.1\r\n", method, path))

	// Headers
	for k, v := range headers {
		requestBuilder.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}

	// Blank line to indicate end of headers
	requestBuilder.WriteString("\r\n")

	// Body (if any)
	if body != "" {
		requestBuilder.WriteString(body)
	}

	return requestBuilder.String()
}

// sendHTTPRequest sends the HTTP request over a TCP connection and returns the response
func sendHTTPRequest(address string, request string) (string, error) {
	// Establish TCP connection
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return "", fmt.Errorf("error connecting to %s: %v", address, err)
	}
	defer conn.Close()

	// Send HTTP request
	_, err = conn.Write([]byte(request))
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}

	// Read HTTP response
	var responseBuilder strings.Builder
	respReader := bufio.NewReader(conn)
	for {
		line, err := respReader.ReadString('\n')
		if err != nil {
			break // EOF is expected when the server closes the connection
		}
		responseBuilder.WriteString(line)
	}

	return responseBuilder.String(), nil
}

func main() {
	// Parse command-line flags and arguments
	requestOpts, err := parseFlags()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Parse the URL
	options, err := parseURL(requestOpts.URL)
	if err != nil {
		fmt.Printf("Error parsing URL: %v\n", err)
		os.Exit(1)
	}

	// Ensure the protocol is supported
	if options.Protocol != "http" {
		fmt.Println("Error: Only HTTP protocol is supported")
		os.Exit(1)
	}

	// Build headers map
	headersMap, err := buildHeaders(options, requestOpts.Headers, requestOpts.Data)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Display connection details and request components
	fmt.Printf("Connecting to %s\n", options.Host)
	fmt.Printf("Sending request %s %s HTTP/1.1\n", requestOpts.Method, options.Path)
	for key, value := range headersMap {
		fmt.Printf("%s: %s\n", key, value)
	}
	fmt.Println()

	/*
		HTTP Request Anatomy
		GET /path?query#fragment HTTP/1.1
		Host: example.com
		Content-Type: application/json

		Body (optional)
	*/

	// Construct the HTTP request
	request := constructHTTPRequest(requestOpts.Method, options.Path, headersMap, requestOpts.Data)

	// Establish TCP connection address
	address := net.JoinHostPort(options.Host, options.Port)

	// Send HTTP request and receive response
	response, err := sendHTTPRequest(address, request)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	/*
		HTTP Response Anatomy
		HTTP/1.1 200 OK
		Content-Type: application/json
		Content-Length: 123
		Connection: close

		Body
	*/

	// Print the HTTP response
	fmt.Print(response)
}
