# 11. Use sentry error monitoring

Authors: Ondrej Ezr


## Status

Accepted


## Problem Statement

We do have logs and prometheus metrics.
We are missing error dashboard with clear list of errors that happened in our app.
It should be a simple place where we can always start with error resolving.

## Goals

* Errors overview from runtime


## Non-goals

* Replace Jira bugs, these errors should be linked to bugs


## Current Architecture

* We need to figure out manually from Grafana chart or Kibana logs if something went wrong 


## Proposed Architecture

* Use sentry to send all errors to [Sentry](https://sentry.io/) or [Glitchtip](https://glitchtip.com/)
  * all codes 5xx responses
  * all recovered panics

## Challenges

* We need to identify all places we error out and convey all necessary info


## Alternatives Considered

* _None_

## Dependencies

* Glitchtip/Sentry instance and project to send data to
* DSN secret for the project in Clowdapp
