# Log

simplified golang logger

## Example

```golang
func main() {

	opts := []log.LoggerOption{
		log.WithLevel("DEBUG"),
		log.WithName("unit_test"),
		log.WithWriter(os.Stdout),
	}

	logger := log.New(opts...)

	logger.Debug("hello world")
}
```