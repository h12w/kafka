h12.me/kafka
============

A modular & idiomatic Kafka client in Go
----------------------------------------

Compatible with Kafka 0.8.2 or later.

Install
-------

```
go get -u h12.me/kafka
```

Design
------

`(-)` means to be done.

### Modular

* **proto**: Kafka Wire Protocol (automatically generated types and
  marshal/unmarshal functions from
  [the HTML spec](https://cwiki.apache.org/confluence/display/KAFKA/A+Guide+To+The+Kafka+Protocol)).
* **broker**: client that talks to a single Kafka broker (concurrent,
  synchronous API wraps inside asynchronous request/response IO).
* **client**: client that talks to a Kafka cluster (leader/coordinator management).
* **producer**: fault tolerant high-level producer (batching and partitioning strategy).
* **consumer**: fault tolerant high-level consumer (consumer group and offset commit).
* **log**: replaceable global logger

### Error Handling

* broker
  + fail fast: timeout returns error immediately
  + release resources carefully
  + reconnect when requested next time
* client
  + metadata reload lazily (only when a leader/coordinator cannot be found in cache)
  + leader/coordinator should be deleted on error
* producer
  + when a leader down, try another partition immediately
  + failed partition will be retried again after a period of time
  + partition expand (-)
* consumer
  + just loop & wait on error
  + partition expand (-)
* graceful shutdown (-)

### Efficiency

* efficiency
  + batching
    - consumer response
    - consumer request (-)
    - producer (-)
  + compression (-)
