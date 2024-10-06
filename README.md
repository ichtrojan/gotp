# GOTP

GOTP is a simple and efficient Go package for generating and verifying one-time passwords (OTPs) using Redis for storage. It provides support for different formats of tokens including alphabetic, alphanumeric, and numeric tokens.

>**NOTE**<br/>
> * This package is named after the GOAT of Formula one, none other than [Goatifi](https://en.wikipedia.org/wiki/Nicholas_Latifi) himself
> * Not named after Go + OTP 🤣

## Features
* Generate one-time passwords with specified length and format.
* Store tokens in Redis with an expiration time.
* Verify tokens against stored values in Redis.
* Automatic token deletion upon successful verification.

## Installation

To install the GOTP package, use the following command:

```bash
go get github.com/ichtrojan/gotp
```

## Prerequisites
* Go 1.16 or later.
* A running instance of Redis.

## Configuration

You need to create a configuration that includes a Redis client before using the package.

## Example Configuration

```go
package main

import (
    "fmt"
    "log"
    "time"
	"github.com/go-redis/redis/v9"
	"github.com/ichtrojan/gotp"
)

func main() {
    rdb := redis.NewClient(&redis.Options{
        Addr: "localhost:6379",
    })

	otp, err := gotp.New(gotp.Config{Redis: rdb})
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	// Continue with token generation and verification...
}
```

## Usage

### Generating a Token

You can generate a token by creating a Generate payload and calling the Generate method on your Config instance.

```go
package main

import (
    "fmt"
    "log"
    "time"

	"github.com/go-redis/redis/v9"
	"github.com/ichtrojan/gotp"
)

func main() {
    // (Assuming Redis client setup as above)

	payload := gotp.Generate{
		Format:     gotp.ALPHA, // or gotp.ALPHA_NUMERIC, gotp.NUMERIC
		Length:     6,
		Identifier: "testIdentifier",
		Expires:    10 * time.Minute,
	}

	token, err := otp.Generate(payload)
	if err != nil {
		log.Fatalf("Error generating token: %v", err)
	}

	fmt.Printf("Generated Token: %s\n", token)
}
```

### Verifying a Token

To verify a token, create a Verify payload and call the Verify method.

```go
payload := gotp.Verify{
    Token:      token, // The token you want to verify
    Identifier: "testIdentifier",
}

valid, err := otp.Verify(payload)

if err != nil {
    log.Fatalf("Error verifying token: %v", err)
}

if valid {
    fmt.Println("Token is valid!")
} else {
    fmt.Println("Token is invalid or expired.")
}
```

## Token Formats

The Generate struct has a Format field which can take the following values:

* `ALPHA`: Generates a token with alphabetic characters only.
* `ALPHA_NUMERIC`: Generates a token with both alphabetic and numeric characters.
* `NUMERIC`: Generates a token with numeric characters only.

## License

This package is licensed under the MIT License. See the LICENSE file for details.