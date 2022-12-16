# structpb-conv

`structpb-conv` is a simple repo for interacting with the `structpb` family of protos.

This library provides utilities to:

- flatten Structs
- minify Structs
- set and get nested values from a Struct

## Installation

To download the latest version of the library run:

    go get -u github.com/structpb-conv/structpb-conv

To download a specific version use:

    go get github.com/structpb-conv/structpb-conv@<version>

## Minify usage

Minify reduces the structpb value to the smallest possible set.

- Empty lists are removed
- Empty structs are removed
- Nil pointers & null values are removed

The result of minifying is that you get the smallest possible set of data to reduce the overhead of storing, processing, and maintaining the data.

MinifyStruct:

```go
const input = `
{
    "key": "value",
    "nil": null,
    "empty-list": [],
    "nested-struct": {
        "nested": "value"
    }
}
`

data := map[string]interface{}{}
err := json.Unmarshal([]byte(input), &data)
require.NoError(t, err)
require.NotNil(t, data)

structpbData, _ := structpb.NewStruct(data)
output := MinifyStruct(structpbData)

outputJson, _ := json.Marshal(output)

// {"key":"value","nested-struct":{"nested":"value"}}
fmt.Println(outputJson)
```
