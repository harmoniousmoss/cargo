# Cargo - Server Performance Testing Tool
Cargo is load testing solution that helps businesses ensure their websites can handle high traffic conditions. Built with Go and Fiber framework, it simulates thousands of concurrent users to identify performance bottlenecks before they impact real customers.

## Features
- **Dynamic Domain Configuration**: Easily update the target domain using a `.env` file.
- **Randomized Requests**: Generate and send requests from a list of 1000 unique paths.
- **High Concurrency**: Utilizes Go's goroutines to simulate real-world load.
- **Performance Insights**: Evaluate server RPS performance and identify bottlenecks.
- **Extensible**: Easily modify for additional HTTP methods, headers, or payloads.

## How It Works

1. **Load Configuration**:
   - The domain is loaded dynamically from the `.env` file.

2. **Path Generation**:
   - A file (`urls.txt`) contains 1000 unique paths (e.g., `/news1.php`, `/api/endpoint`, etc.).
   - Paths are read and randomized for each request.

3. **Concurrency**:
   - Requests are sent concurrently using goroutines, simulating real-world high traffic scenarios.

4. **Performance Evaluation**:
   - Measures the time taken to complete all requests, helping determine the server's RPS handling capacity.

## Example `.env` File

```env
DOMAIN=https://example.com
```

## Example `urls.txt`

```txt
/news1.php
/blog/random1.php
/api/data?query=test
/page123.html
/products/item99
```

# Terms and Conditions
Cargo is designed strictly for testing and evaluating the performance of servers that you own or have explicit permission to test. By using this tool, you agree to the following terms:

## Authorized Use Only:
You must only use Cargo on servers you own or have written consent to test. Unauthorized use on third-party servers is strictly prohibited.

## No Malicious Intent:
Cargo must not be used for malicious activities such as Distributed Denial-of-Service (DDoS) attacks or any other actions intended to disrupt, harm, or degrade systems.

## Compliance with Laws:
Users must comply with all applicable local, national, and international laws and regulations when using GoRPS.

## Liability Disclaimer:
The developers of Cargo are not responsible for any misuse, illegal activity, or damage caused by the tool. Use it at your own risk.

By using this tool, you acknowledge and agree to these terms. Violations may lead to legal consequences.
