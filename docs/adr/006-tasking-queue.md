# 6. Tasking queue

Authors: Lukas Zapletal


## Status

Accepted

Superseded by [10. Simplify job queue](010-simplify-jobqueue.md)

## Glossary

* Image builder's *jobqueue* library - internal code at https://github.com/osbuild/osbuild-composer/tree/09f57b6c2f8199a0a8a9365dfb56b3e52e7a60e4/internal/jobqueue that image builder uses to manage their tasks.
* **Task chaining** Running multiple tasks synchronously one after another in predefined order

## Problem Statement

We do not want to block requests on instance deployments and provisioning steps.
These should be handled asynchronously as those are potentially long running.
Thus we want to define a background processing library that will make handling longer running tasks easier.

## Goals

* Define our background processing engine and its dependencies
* Simple and easy to use solution
* Fast enough to not block us even under bit of a pressure
* Robust enough to support task chaining


## Non-goals

* Not aiming at Process orchestration
* Not necessary to have complex task management


## Current Architecture

None


## Proposed Architecture

We decided to ask the Image buider team to extract their *jobqueue* library
to an external library and abstract logging and connection pool. Then
use the *jobqueue* and all its features for our tasking system.

We have decided following in accordance to usage of this tasking system:
* We will name all items as "jobs" and not "tasks".
* We do not want to expose secret information through the API.
* We should not have a resource named "jobs" but rather "launch" or "provision".
* We will not utilize channels and job dependencies features in the first version.
* We will define an abstraction API on top of *jobqueue* that will make it easier to switch to different implementation.

## Challenges

* Solution developed to support one team's needs might be hard to adapt to.
* There is no documentation and we'll need to maintain knowledge about the solution ourselves.


## Alternatives Considered

We have explored AWS SQS managed service which can used on ConsoleDot.
It runs on AWS and is supported by AppSRE team.

Building a task system on SQS is straightforward, we could utilize an
existing library as well. It is possible to implement heartbeat but
not so easy to implement task dependencies in the same way as in
*jobqueue*. It should be possible to create task chains (a task schedule
another task).

However, we do not believe that SQS brings any value to the table.
Since we do not expect heavy load, PostgreSQL NOTIFY will suffice,
the same applies to costs (not a big difference for small load we
expect). We do keep SQS in mind as a backup plan if we find PostgreSQL
not suitable at any moment.


## Dependencies

* Extraction of Image builder's *jobqueue* in separate library


## Stakeholders

* EnVision engineering manager - needs to ack this cross team collaboration
  * Marek Hulan
* Image builder's PO - needs to ack *jobqueue* extraction
  * Sanne Raymaekers
* EnVision developers - need to ack the challenges
  * Adi Abramovich
  * Amir Feferkuchen
  * Lukas Zapletal
  * Ondrej Ezr


## Consequences

* We end up partially maintaining Image builder's *jobqueue*
