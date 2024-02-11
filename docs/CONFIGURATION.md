# Configuration

```json
{
    "providers": [],
    "rules": []
}
```
- **providers**: (repeated *Provider*) Configuration of JWT validation methods
- **rules**: (repeated *Rule*) Configuration of path matching and validation rules

## Provider
```json
{
    "name": "string",
    "from_header": {},
    "from_param": "string",
    "validate": {}
}
```
- **name**: (string, REQUIRED) Name of the validation method
- **from_header**: (*Header*) Obtain JWT from the request header
- **from_param**: (string) Obtain JWT from query parameter
- **validate**: (*Validate*) Validation method

## Header
```json
{
    "name": "string",
    "value_prefix": "string"
}
```
- **name**: (string, REQUIRED) Key of the request header
- **value_prefix**: (string) Prefix of the value of the request header

## Validate
```json
{
    "issuer": "string",
    "audiences": [],
    "jwk": "string"
}
```
- **issuer**: (string) Match the specified issuer
- **audiences**: (repeated string) Match the specified audiences
- **jwk**: (string, REQUIRED) Key

## Rule
```json
{
    "match": {},
    "requires": {}
}
```
- **match**: (*Match*, REQUIRED) Specify the matching path
- **requires**: (*Requires*) Specify how to perform validation

## Match
```json
{
    "prefix": "string",
    "path": "string"
}
```
- **prefix**: (string) Specify the required prefix for matching
- **path**: (string) Specify the matched path

## Requires
```json
{
    "name": "string",
    "requires_any": [],
    "requires_all": []
}
```
- **name**: (string) Name of the validation method, indicating the method to be used for validation
- **requires_any**: (repeated *Requires*) If any one of them is true, the validation condition is satisfied
- **requires_all**: (repeated *Requires*) All of them must be true to satisfy the validation condition
