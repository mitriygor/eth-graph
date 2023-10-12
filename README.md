# The Graph API

This is a REST API developed in Golang and utilizes The Graph’s GraphQL API to retrieve and provide information related
to Uniswap v3 upon user request.

## Overview

The API doesn't use any frameworks, and it has very few dependencies.

The API has versioning, — e.g. v1, — and it exposes four endpoints:
- /v1/tokens/{tokenID}/pools — based on token ID, it returns a list of pools that include the token;
- /v1/tokens/{tokenID}/volume?from={from}&to={to} — based on token ID, it returns the total volume of the token swapped in the time range;
- /v1/blocks/{blockNumber}/swaps  — based on a block number, it returns what swaps occurred during the block;
- /v1/blocks/{blockNumber}/swaps/tokens — based on a block number, it returns a list of tokens swapped during the block;

In the assets directory, there is a Postman collection that can be used for testing the API.

The entire API has logging-middleware. In addition, it has a rate-limiting middleware for the endpoints that are using The Graph's API (presently, it is one request per second)

Furthermore, all requests have timeout (presently, it is 5 seconds) in order to prevent long waiting if the third-party API is slow.

For development purposes, the API accepts both HTTPS and HTTP requests.

Just of the presentation purposes, the code contains excessive-level of comments.

### Dependencies

- [Chi](https://github.com/go-chi/chi) — lightweight, idiomatic and composable router for building Go HTTP services;
- [Air](https://github.com/cosmtrek/air) — live-reloading command line utility for developing Go applications;
- [Zap](https://github.com/uber-go/zap) — fast, structured, leveled logging in Go;

## Structure

```
eth-graph-api/                        # Root directory of the project
|-- cmd/                              
|   |-- api/                          
|   |   |-- init.go                   # File for initializing repositories, services, and handlers
|   |   |-- main.go                   
|   |   |-- routes.go                 
|
|-- internal/                         # Business logic that needs to be isolated
|   |-- token/                        # The token package directory
|   |   |-- token.go                  # Definitions of the token-related structs
|   |   |-- token_handler.go          # HTTP handler related to "token"
|   |   |-- token_service.go          # Service layer related to "token"
|   |   |-- token_repository.go       # Repository layer related to "token"
|   |   |-- token_repository_test.go  
|   |   |-- token_service_test.go     
|   |   |-- token_handler_test.go     
|   |
|   |-- block/                        # "block" package directory
|   |   |-- block.go                  # Definitions of structs related to "block"
|   |   |-- block_handler.go          # HTTP handler related to "block"
|   |   |-- block_service.go          # Service layer related to "block"
|   |   |-- block_repository.go       # Repository layer related to "block"
|   |   |-- block_repository_test.go  
|   |   |-- block_service_test.go     
|   |   |-- block_handler_test.go     
|
|-- pkg/                              # Pieces of functionality meant to be shared across the project
|   |-- calc/                         # "calc" package directory
|   |   |-- calc.go                   # Calculates big numbers presented in string-format
|   |   |-- calc_test.go              
|   |
|   |-- formatter/                    # "formatter" package directory
|   |   |-- formatter.go              # Formatting numbers in USD-currency format
|   |   |-- formatter_test.go         
|   |
|   |-- json_helper/                  # "json_helper" package directory
|   |   |-- json_helper.go            # Processes JSON data
|   |   |-- json_helper_test.go       
|   |
|   |-- validator/                    # "validator" package directory
|   |   |-- validator.go              # Validates different types of data, e.g. timestamp, token ID, etc.
|   |   |-- validator_test.go         
|
|-- .air.toml                         # Configuration file for Air, for the live-reloading
|-- eth-graph-api.dockerfile          # Dockerfile for building the image
|-- go.mod                            # Go module file, defining module path and dependency versions
|-- go.sum                            # Checksums of dependency files
                     

```


## Requests:

### Tokens

#### GET: /v1/tokens/{tokenID}/pools

Based on given a token ID, it returns a list of pools that include that specific token. By default, it returns first 5 pools.

**Token ID example:** 0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2

```
[
    {
        "id": "0x0000d36ab86d213c14d93cd5ae78615a20596505"
    },
    {
        "id": "0x0001fcbba8eb491c3ccfeddc5a5caba1a98c4c28"
    },
    {
        "id": "0x0005647722f0f8233d64063dcae20ebdb99fdc02"
    },
    {
        "id": "0x000c0d31f6b7cecde4645eef0c4ec6a492659d62"
    },
    {
        "id": "0x000ea4a83acefdd62b1b43e9ccc281f442651520"
    }
]
```

#### GET: /v1/tokens/{tokenID}/volume?from={from}&to={to}

Based on given a token ID, it returns what is the total volume of that token swapped in a given time range. The 'from' and 'to'
are UNIX timestamps.
If the 'from' and 'to' aren't provided, there will be returned the volume for last 24 hours.

**Token ID example:** 0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2
**Example of from or to:** 1632960000

**Response example:**

```
{
    "volume": "8039962301.3240150269"
}
```

### Block

#### GET: /v1/blocks/{blockID}/swaps

Based on given a block number, it returns what swaps occurred during that specific block. By default, it returns first 5 swaps.

**Block Number example:** 18319881

**Response example:**

```
[
    {
        "id": "0x0000006cfab44a8db8060b58bb9f7c261ce2c2554cd954bbec2bc5024d191720#2345455",
        "pool": {
            "token0": {
                "symbol": "USDC"
            },
            "token1": {
                "symbol": "WETH"
            }
        }
    },
    {
        "id": "0x00000101018767cd08b5909373a88eaba5b673915ab1562f47ab3d75ca9bc638#750832",
        "pool": {
            "token0": {
                "symbol": "DAI"
            },
            "token1": {
                "symbol": "WETH"
            }
        }
    },
    {
        "id": "0x0000017229c4e5d849d31c303e59efe3e600d1082cea5ba0d9a8222002472503#54483",
        "pool": {
            "token0": {
                "symbol": "SPELL"
            },
            "token1": {
                "symbol": "WETH"
            }
        }
    },
    {
        "id": "0x000001ac8cbe4303d766b166e6122f413aa6c1ce65ad05e1cd478d996de695f3#181101",
        "pool": {
            "token0": {
                "symbol": "DAI"
            },
            "token1": {
                "symbol": "USDC"
            }
        }
    },
    {
        "id": "0x0000030c8c6599a34beadf79593f29a015a38e4cbd1bb618d6462d0ca2c8d965#2370551",
        "pool": {
            "token0": {
                "symbol": "WETH"
            },
            "token1": {
                "symbol": "USDT"
            }
        }
    }
]
```

#### GET: /v1/blocks/{blockID}/swaps/tokens

Based on given a block number, it returns a list of all tokens swapped during that specific block. By default, it gets first 5
swaps.

**Block Number example:** 18319881

**Response example:**

```
[
    {
        "symbol": "USDC"
    },
    {
        "symbol": "WETH"
    },
    {
        "symbol": "DAI"
    },
    {
        "symbol": "SPELL"
    },
    {
        "symbol": "USDT"
    }
]
```

## Running and testing

```
eth-graph-builder/                    # Directory containing the build files
|-- .env                              # Environment variables configuration file
|-- docker-compose.yaml               # Defining and running the project images. Presently, the project contains just one
|-- Makefile                          # A file containing a set of directives used by make build automation tool
          
```

In the eth-graph-build, run:

```
make run
```

To test:

```
make test
```
