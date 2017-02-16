nats_metrics is a simple metrics library that will push multidimensional metrics onto a nats subject.

It supports three types of metrics:

- Counters ~ good for counting when things happen
- Gauges ~ good for measuring the level when things happen
- Timers ~ good for timing how long things take

This library will do a nats publish each time you measure something. It does NOT do any worker pooling or fancy-ness to make that non-blocking. That is all delegated to the nats client.

## getting started

To start it is important that you initialize the library. It requires a subject and a nats connection.

``` go
  func main() {
	  nc, _ := nats.Connect("nats://localhost:4222")
    nats_metrics.Init(nc, "metrics")

    // push metrics
    counter := nats_metrics.NewCounter("metric.name", nil)
    counter.Count(DimMap{"response_code": 200})
  }
```

## dimensions

Dimensions are done on 3 levels:

- Global   `nats_metrics.AddDimension("key", "value")`
- Metric   `counter.AddDimension("key", "value")`
- Instance `counter.Count(DimMap{"key": "value"})`

Global dimensions will be on all metrics being sent. Metric level will be on all emissions of that metric. Instance level will be on *only* that emission. If a dimension exists in multiple levels, the overrides go like this: Instance > Metric > global.

## output

It publishes values as JSON:

```
{
  "name": "metric.name",
  "type": "gauge",
  "value": 1,
  "timestamp": "2009-11-10T23:00:00Z",
  "dimensions": {
    "app": "example",
    "key": "value"
  }
}
```


## timestamps

A gauge, counter, and timer each have a `SetTimestamp` method that will let you override the system generated timestamp. Note _this is for all calls on that metric going forward_. That is to say, if you set it that timestamp will always be sent. To reset it (aka let the timestamp be system generated) you can just reset it to `time.Zero`.

``` go
ts := time.Now()
counter := metrics.NewCounter("metric.name", nil)
counter.SetTimestamp(ts)
counter.Count(nil) // 'timestamp' will == ts

<- time.After(1 * time.Second)
counter.Count(nil) // 'timestamp' will == ts

counter.SetTimestamp(time.Time{})
counter.Count(nil) // 'timestamp' will != ts, it will be set by the system (time.Now())
```
