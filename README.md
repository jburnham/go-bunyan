go-bunyan
=========

Create a new stdout logger

```
import github.com/jburnham/go-bunyan

stdoutStream := bunyan.NewStdoutStream(bunyan.Info, nil)
log := bunyan.NewLogger("myLoggingApplication", []bunyan.StreamInterface{stdoutStream})
log.Info("My First Bunyan Logging Line")
```
