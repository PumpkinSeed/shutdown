# shutdown

Graceful shutdown package for Go services with built-in healthcheck for persistent dependencies.

# Usage

### Shutdown

Everything what depends on the graceful shutdown (ex.: long-time file writes, database connection, http requests), most implement the `ServiceDescriptor` interface.

```go
package example
type ServiceDescriptor interface {
	Stop() error
	Ping() error
	Reconnect() error
}
```

Also the constructor requires a logging instance with the following interface implementation.

```go
package example
type Log interface {
	Errorf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
}
```

Create a new instance of the shutdown and adding new services and call the `GracefulExit` handler.

```go
package example

handler := shutdown.NewHandler(&testLog{t})
handler.Add("service1", "", shutdown.Init, &service1)
wg := &sync.WaitGroup{}
shutdown.GracefuleExit(handler, wg)

wg.Wait()
```

Determine the sequence of the shutdown when we add them into the list.

```go
package example

handler := shutdown.NewHandler(&testLog{t})
handler.Add("service1", "", shutdown.Init, &service1)
handler.Add("service2", "service1", shutdown.Before, &service2)
handler.Add("service3", "service2", shutdown.After, &service3)
```

### Important
Healthcheck not fully implemented yet!



