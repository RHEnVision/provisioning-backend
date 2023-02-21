# 9. Use scany and pgx in DAO

Authors: Lukas Zapletal

## Status

Accepted

Amends [2. Backend main tech stack](002-backend-main-tech-stack.md)

## Problem Statement

Current solution does not offer solution for goals mentioned below.


## Goals

* Zerolog logging of DB
* OpenTelemetry tracing for DB
* Add support for PostgreSQL's native types like JSONB
* We want to improve errors from SQLs


## Non-goals

* Make our own SQL errors


## Current Architecture

* DAO uses database/sql and sqlx


## Proposed Architecture

Replace database/sql and sqlx with pgx driver and scany library.

* The code is cleaner and more simple since pgx driver automatically prepares SQL statements and cache them.
* Better errors (error messages do appear to be more detailed).
* Native logging via zerolog and tracing via opentelemetry.
* The driver natively supports many custom types, including JSONB which we use (which should be faster to encode/decode).
* The driver of choice for image builder, this allows us to directly use jobqueue and we could drop our copy from the codebase.

Uses very recent version of pgx (v5) and scanny (v2) and tern migration library (v2).
While pgx/v5 has been out for few weeks, both mentioned libraries are currently in v2-beta stage.

## Challenges

* Uses very recent version of pgx (v5) and scanny (v2) and tern migration library (v2)

## Dependencies

* https://github.com/jackc/tern/pull/69
