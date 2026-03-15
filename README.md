# HTTP from TCP

HTTP library written on top of TCP module (net) in Golang.

## Motivation

The goal of this project was not to create the best client for working with HTTP traffic, but to understand the principles behind HTTP by implementing it yourself. During this project I read RFC specs for HTTP 1.1, which helped me understand the semantics behind HTTP that are similar across all versions, wrote my own parser, and implemented trailer headers as well as basic server functionality.

## Background

HTTP (Hypertext Transfer Protocol) is an application-layer protocol built on top of TCP. By implementing HTTP from scratch using raw TCP sockets, we gain a deeper understanding of:
- How request/response lines are formatted
- How headers are structured and parsed
- How chunked transfer encoding works
- How trailer headers enable streaming metadata

This project demonstrates the core concepts without relying on Go's `net/http` package.

## Project Structure

```
.
├── cmd/
│   ├── httpserver/       # HTTP server with example handlers
│   ├── tcplistener/     # TCP listener utility
│   └── udpsender/       # UDP sender utility
├── internal/
│   ├── headers/         # HTTP headers parsing and handling
│   ├── request/        # HTTP request parser
│   ├── response/       # HTTP response writer
│   └── server/         # TCP server implementation
├── go.mod
└── Makefile
```

- **internal/request** - Parses HTTP/1.1 request line, headers, and body
- **internal/response** - Writes responses with chunked encoding and trailers
- **internal/headers** - Generic header parsing utilities
- **internal/server** - TCP listener with handler pattern
- **cmd/httpserver** - Example server demonstrating proxy and binary handlers

## Features

- Custom HTTP/1.1 request parser (RFC 7230)
- Request line parsing (method, path, version)
- Headers parsing with Content-Length support
- Body reading with Content-Length validation
- Chunked transfer encoding (Transfer-Encoding: chunked)
- Trailer headers support for streaming metadata
- TCP server with graceful shutdown (SIGINT/SIGTERM)
- Example proxy handler to httpbin.org

## Architecture

```
TCP Connection
      │
      ▼
┌─────────────────┐
│   request.Parse │  ← Parses request line, headers, body
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│    Handler      │  ← User-defined (e.g., ProxyHandler)
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ response.Writer │  ← Writes status, headers, body/chunked
└─────────────────┘
```

## Installation

```bash
git clone https://github.com/evok02/httpfromtcp.git
cd httpfromtcp
go mod download
```

## Usage

```bash
# Run the HTTP server
make run

# Run tests
make test
```

The server starts on port 42069 by default.

### Example Handlers

**Proxy Handler** - Proxies requests to httpbin.org with chunked encoding:
```
http://localhost:42069/httpbin/get
http://localhost:42069/httpbin/post
```

**Binary Handler** - Serves binary data:
```
http://localhost:42069/video
```

## How It Works

1. **Request Parsing**: The server accepts a TCP connection and reads raw bytes into a buffer. It parses the request line (method, path, HTTP version), then headers, then body (if Content-Length is present).

2. **Handler Execution**: The parsed `*Request` and `response.Writer` are passed to a user-defined handler function. The handler sets headers and writes the response.

3. **Response Writing**: The response writer handles three modes:
   - Regular body with Content-Length
   - Chunked body (streaming)
   - Trailer headers (appended after chunks)

4. **Graceful Shutdown**: The server listens for SIGINT/SIGTERM and closes the TCP listener cleanly.

## Limitations

- HTTP/1.1 only (no HTTP/2 or HTTP/3)
- No keep-alive support (Connection: close always)
- No TLS/HTTPS support
- No request routing (single handler)
- Limited status code support (200, 400, 500)

## References

- [RFC 7230 - HTTP/1.1 Message Syntax and Routing](https://tools.ietf.org/html/rfc7230)
- [RFC 7231 - HTTP/1.1 Semantics and Content](https://tools.ietf.org/html/rfc7231)

## License

MIT
