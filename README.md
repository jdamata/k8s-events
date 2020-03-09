# k8s-event
Logs kubernetes events as they are generated. Useful for having a history of k8s-events when paired with a log aggregation system such as ELK, datadog, splunk, etc..

Sample generated log:
```json
{"FirstTimestamp":"2020-03-09T04:58:36Z","LastTimestamp":"2020-03-09T05:13:47Z","count":61,"event_level":"Warning","level":"info","msg":"Invalid metrics (1 invalid out of 1), last error was: failed to get cpu utilization: unable to get metrics for resource cpu: unable to fetch metrics from resource metrics API: the server could not find the requested resource (get pods.metrics.k8s.io)","namespace":"default","pod":"testapi-sample.15fa89f7d2fcab0c","reason":"FailedComputeMetricsReplicas","time":"2020-03-09T01:33:56-04:00"}
```

## Deployment 
Use the [helm chart](helm_chart)

## Running locally
```bash
make run
```
