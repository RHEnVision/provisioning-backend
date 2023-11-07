# 10. Simplify job queue

Authors: Lukas Zapletal


## Status

Accepted

Supersedes [6. Tasking queue](006-tasking-queue.md)


## Problem Statement

Since we now have more information and understanding of the domain and environment, we can simplify our job queue.
We no longer think we need too complex tasks for our solution.
Our solution should be as simple as possible and implemented internally.

## Goals

* Simplify tasking system as much as possible
* Choose one synchronization backend


## Non-goals

* Implement fully featured tasking system

## Current Architecture

* dejq library as wrapper around three synchronization backends.
* Ability to choose backend and switch between them.


## Proposed Architecture

Move dejq code from the external repository into pkg/worker and simplifies the API to only cover what our project needs.
The job queue backend is now Redis (with additional memory impl for development).

Simplify interface to:

```go
type Job struct {
	// Random UUID for logging and tracing. It is set by Enqueue function.
	ID uuid.UUID

    // Associated account.
    AccountID int64

    // Associated identity
    Identity identity.Principal

	// Job type or "queue".
	Type JobType

	// Job arguments.
	Args any
}
```

Jobs are super simple, just args, no output, no return value or error, no dependencies.
For this reason, AWS launch job is merged into a single job (with two functions/steps).

Uses Go's gob encoding instead of JSON which gives more type safety.
Improve logging since there is no need for a logging facade now.
All logs now have the following fields:
```
account_id
account_number
org_id
job_id
job_type
reservation_id
```

Finally, statistics now also contains "in flight" jobs, meaning jobs which are currently being processed by the worker.


## Challenges

* There will be no dependencies for jobs, so in case we need them, we will need to redesign again.

## Alternatives Considered

* [taskq](https://taskq.uptrace.dev/) library
  * It is much more code than we actually need to have.
  * Most of its complexity comes from the delaying feature and rescheduling of failed jobs.
  * None of that we actually need

## Dependencies

* _None_


## Stakeholders

* EnVision developers

## Consequences

* dejq will be deprecated
* all our code will live in single repository
