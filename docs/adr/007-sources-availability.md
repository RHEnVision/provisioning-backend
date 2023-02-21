# 7. Sources availability check

Authors: Lukas Zapletal


## Status

Accepted


## Problem Statement

Sources require applications to implement availability checking.
A REST endpoint will be called every 30/30 minutes from Sources background job.
Apps should send a reply via Kafka for all authentications associated with such source ID.
Requests are being made in parallel up to 3 requests.

To verify an AWS account is active and credentials are correct, an AssumeRole operation must succeed. Other cloud providers (Azure, GCP) do not have the concept of AssumeRole, in these cases we need to do some non-mutating operation (e.g. list all regions) instead on behalf of the customer identity (subscription/project id).

For the record, we plan to only check if we can connect to the cloud using the given ARN/Subscription/Project.
Check all permissions (e.g. create instances) is non-trivial task.
So in these cases when customers do mess around with permissions we cannot easily detect.

It should be possible to use a call like DescribeRole to compare permission JSON policy with what we expect and the availability check could actually report an additional state: source is present, working but permission has an old version of policy. This is an idea which we will not implement in the initial version.

Like everything through AWS API, rate limiting is applied for all operations. For non-mutating actions this is 20 operations per second with an initial bucket size of 100 which allows bursting of operations. Now, it is not clear if the rate limitation is applied to the service account or to the customer account. One answer actually says that AssumeRole operations are counted against the service account and all other assumed operations against the customer account, which makes sense. In any case, this is a problem because a single availability check from sources can limit other operations on the account.

Now, quota limits can be increased in the AWS Console and furthermore increased after contacting AWS Support, however, we want to build a scalable service that will also work with other cloud providers where more strict limits may apply. These kinds of issues can be hard to track as cloud SDKs do provide a client backoff failover features (some enabled by default) therefore some operations can be delayed slowing down services, or the whole platform, until the operations team will finally see some client errors.

## Goals

* Check ability to perform actions on behalf of customer


## Non-goals

* Checking if all necessary permissions are still available


## Current Architecture

* None


## Proposed Architecture

*Provisioning backend REST API -> Kafka -> Background job -> Kafka -> Sources*

Since all results need to be sent over Kafka, one obvious solution is to leverage Kafka itself as a temporary store of source IDs to be checked.

In this scenario, the REST API would accept Source requests similarly to the Redis solution, but instead storing source IDs into Redis, a Kafka message would be created. The only subscriber of such a stream of the data would be a single background job which would be started as a regular service (single instance, not a scheduled job).

The subscriber worker process would process source IDs one by one keeping a map of already seen source IDs and last seen timestamp. This is a common streaming/messaging pattern also known as "rolling window". Every new source ID is checked and the result is sent over Kafka back to sources and a map entry is added. Every already-seen source ID within a specified time window (expiration time of check results) would be returned from the map.

This approach is very easy to implement, already seen records can be simply kept in memory in a map with time expiration, no locking is needed since only single worker process would be always spawned for a deployment and this can be easily monitored as the stream of requests coming in can be immediately seen through Kafka brokers. Similarly to Redis implementation, only dozens of MBs of memory would be expected for millions of source IDs, so no additional storage is required.

The solution is also resilient to background worker failures - if the worker fails, there will be a stream of events that nobody listens to. Kafka will be configured to only keep the data for a short period of time (minutes, hour at maximum) and then throw old data away. When a new worker is started up by k8s, it can be configured to either start listening to new events, or work on a backlog to fill up its map (cache window).

This solution only makes sense, if we could reuse the existing Kafka infrastructure. The backend application would publish into the Kafka topic that no other platform services would be interested in. This could be a precedent on the platform.

This solution does allow immediate (user-created) availability checks which can be implemented as a simple "do not cache" flag in the Kafka message payload. In that case, the background worker would simply ignore the "last seen" state and perform a re-check immediately. Additionally, a second "priority" topic could be created so these events are dispatched even more quickly.


{% uml src="diagrams/source-availability.plantuml" %}{% enduml %}


## Challenges

* Solution requires understanding advanced concepts of go.


## Alternatives Considered

* Batching via Redis
* Batching via Postgres


## Dependencies

* Kafka
* Diagram of communication [`source-availability.plantuml`](diagrams/source-availability.plantuml)


## Stakeholders

* EnVision team (devs and QE)


## Consequences

* We need to implement Kafka handling
* The availability check for sources will be complex, but quite scalable.
