# Calculator CLI

A very simple Go CLI calculator.

## Build

```bash
go build -o calc-go .
```

## Usage

```bash
./calc-go <operation> <a> <b>

# Operations: add, sub, mul, div
./calc-go add 5 3    # Output: 8.0
./calc-go sub 10 4   # Output: 6.0
./calc-go mul 6 7    # Output: 42.0
./calc-go div 20 4   # Output: 5.0
```

## Testing

```bash
go test ./...
```
