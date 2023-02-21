# 5. Migrate with Tern

Authors: Lukas Zapletal


## Status

Accepted

Amends [2. Backend main tech stack](002-backend-main-tech-stack.md)

## Problem Statement

Not every migration is done in a transaction. Even when each SQL file is enclosed with BEGIN/COMMIT this does not work as expected - when something fails, golang-migrate leaves the database in a "dirty" state asking the user to manually fix this.

## Goals

* Run migration in transaction
* Improve migrations error reporting

## Current Architecture

Using golang-migrate, which has following challenges:
* Biggest concern is that not every migration is done in a transaction. Even when each SQL file is enclosed with BEGIN/COMMIT this does not work as expected - when something fails, golang-migrate leaves the database in a "dirty" state asking the user to manually fix this.
* The second big concern are errors - error reporting in golang-migrate is generic and not very helpful. Tern supports native PostgreSQL error parsing with very detailed error messages including line and column extraction.
* Golang-migrate has no validation of file naming, specifically for sequenced migrations. We currently have some code that performs some validations, with Tern this all can be dropped.


## Proposed Architecture

Replace golang-migrate with tern. We found golang migrate problematic for couple of reasons which are all solved in tern: https://github.com/jackc/tern

* Biggest concern is that not every migration is done in a transaction. Tern is designed by default to work with transactions, migration either succeeds or fails, there is no dirty state.
* Tern supports native PostgreSQL error parsing with very detailed error messages including line and column extraction.
* We currently have some code that performs some validations, with Tern this all can be dropped.
* We agreed not to write down migrations for now. In Tern there are no up and down migrations as separate files, its a single file and down part is optional via a separator.
* Tern is more simple, less code to worry about, yet powerful and its from the respected PGX PostgreSQL driver author.

## Challenges

* We need to remigrate all our databases, but it is minor issue at this point


## Alternatives Considered

* None

## Dependencies

* PROV00001, which is adjusted by this ADR


## Stakeholders

* EnVision developers


## Consequences

* Migration management is simplified and more reliable
