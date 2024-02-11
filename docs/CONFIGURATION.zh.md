# Configuration

```json
{
    "providers": [],
    "rules": []
}
```
- **providers**: (repeated *Provider*) JWT 验证方式的配置
- **rules**: (repeated *Rule*) 路径匹配和验证规则的配置

## Provider
```json
{
    "name": ...,
    "from_header": {...},
    "from_param": ...,
    "validate": {...}
}
```
- **name**: (string, REQUIRED) 验证方式的名称
- **from_header**: (*Header*) 从请求头中获取 JWT
- **from_param**: (string) 从查询参数中获取 JWT
- **validate**: (*Validate*) 验证方式

## Header
```json
{
    "name": ...,
    "value_prefix": ...
}
```
- **name**: (string, REQUIRED) 请求头的键
- **value_prefix**: (string) 请求头的值的前缀

## Validate
```json
{
    "issuer": ...,
    "audiences": [],
    "jwk": ...
}
```
- **issuer**: (string) 匹配指定的颁发者
- **audiences**: (repeated string) 匹配指定的受众
- **jwk**: (string, REQUIRED) 密钥

## Rule
```json
{
    "match": {...},
    "requires": {...}
}
```
- **match**: (*Match*, REQUIRED) 设定匹配的路径
- **requires**: (*Requires*) 设定如何进行验证

## Match
```json
{
    "prefix": ...,
    "path": ...
}
```
- **prefix**: (string) 设定匹配所需的前缀
- **path**: (string) 设定匹配的路径

## Requires
```json
{
    "name": ...,
    "requires_any": [],
    "requires_all": []
}
```
- **name**: (string) 验证方式的名称，表示使用使用指定方式进行验证
- **requires_any**: (repeated *Requires*) 其中任意一个结果为真则满足验证条件
- **requires_all**: (repeated *Requires*) 其中所有的结果为真才满足验证条件
