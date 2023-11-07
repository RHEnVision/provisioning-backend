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
* Use zerolog writer to send all events with log level Error or higher to Sentry
  * This approach makes sure we send all errors when logging is set up correctly
  * Puts more incentive on logging correctly as it is even more useful to do so now

## Challenges

* Make sure we send all relevant information to sentry.
  * Can evolve over time


## Alternatives Considered

* _None_

## Dependencies

* Glitchtip/Sentry instance and project to send data to
* DSN secret for the project in Clowdapp
