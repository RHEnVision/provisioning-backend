# 2. Initial backend (Go) dependencies

Authors: Lukas Zapletal, Ondrej Ezr


## Status

Accepted

Amended by [5. Tern migraions](005-tern-migrations.md)
Amended by [9. Use scany and pgx in DAO](009-dao-use-scany.md)

## Problem Statement

We want to select our main tech stack and main dependencies so we know clearly what we depend on and why we’ve chosen each of them.


## Goals

* Select initial stack and dependencies we start the backend project with
* Put the dependencies in [https://github.com/RHEnVision/provisioning-backend](https://github.com/RHEnVision/provisioning-backend)
* Take personal preferences into account [https://docs.google.com/forms/d/1TJnym_YyHdV2oxdsCAEwvUcYpg_HkfYkkQa1EhU_s5o/edit#responses](https://docs.google.com/forms/d/1TJnym_YyHdV2oxdsCAEwvUcYpg_HkfYkkQa1EhU_s5o/edit#responses)


## Non-goals

* These dependencies are not meant to persist forever
* It is not meant to be a full list, just the main dependencies


## Current Architecture

* None


## Proposed Architecture

As the main language we’ve chosen Golang as we like it best from the proposed choices Python, Java and Golang. Go is cloud native, fastest and feels easiest to code in and debug.

We want to have generated types for OpenAPI specs, but maintain OpenAPI schema and take schema first approach for the rest.

Following Go dependencies have been chosen:

* **go 1.16**. While 1.18 with its big feature ("generics") is stable currently, SRE team requires us to use a UBI8 container with Go from RHEL8. That is currently 1.16, we will upgrade as soon as possible.
* **sqlx for database**. Creating a DAO (Database Access Objects) layer is considered useful [even when using Gorm ORM](https://github.com/RedHatInsights/sources-api-go/tree/main/dao). Therefore this is our opportunity to keep things simple, sqlx is extremely easy to learn as it only provides row scanning (mapping to Go structs). Also we might want to use some advanced features of PostgreSQL (e.g. work queue via LISTEN/NOTIFY and locking). If we ever find any need for full ORM, we can switch to Gorm later and thanks to the DAO package it will be easier.
* **pgx as database driver**. Instead of using the standard pg driver, pgx is a native port and it is used by some consoledot projects. It provides more features than the standard driver, it provides both the standard library interface and extended interface which can be used as well for additional features (LISTEN/NOTIFY, JSON/JSONB etc).
* **golang-migrate for migrations**. The most popular migration library supports plain SQL migrations as well as both timestamp-based and sequential naming schemes. More importantly, it implements migration locking via PostgreSQL advisory locks.
* **chi for routing**. It is a very simple and easy to understand routing library with binding support, there are not any features we will be missing out from bigger frameworks as most of the current console projects do write their own logging and telemetry middleware anyway. Since chi is compatible with the standard library, if we choose to change a framework we can do it gradually.
* **environmental variables for configuration**. Most of consoledot projects use viper for parsing of the environmental variables, however, this can be seen as an overkill as only small portion of the functionality of the library is used. Parsing of environment variables and .env files can be done with two tiny libraries while keeping a good developer experience.
* **zerolog for logging**. Zerolog is a robust logging library with strictly typed structured logging and very good AWS CloudWatch batching integration.
* **prometheus for monitoring**. This is a hard requirement by the SRE team. The initial application implements a basic request duration and counter monitors.
* **kafka for message queue**. This is a hard requirement by the SRE team. Not implemented in the initial application yet.
* **GNU GPL v3 license**. Different consoledot projects use different licenses, this is the most restrictive license which made Linux kernel possible.


## Challenges

* All the team members need to learn Go and get familiar with the development style it proposes and it may take some time before the team gets productive


## Alternatives Considered

* **Gorm** for ORM seemed too heavy and not necessary at this point as we expect lightweight DB. We want to revisit this decision once we hit 5 or more tables
* **Echo** for routing - it has more mature community, but it is bit more heavy dependency and it is easy to switch to Echo from Chi then the other way around. We had taken a look at [https://brunoscheufler.com/blog/2019-04-26-choosing-the-right-go-web-framework](https://brunoscheufler.com/blog/2019-04-26-choosing-the-right-go-web-framework) and didn't find strong incentive against Chi


## Dependencies

* go.mod in [RHEnVision/provisioning-backend](https://github.com/RHEnVision/provisioning-backend)


## Stakeholders

* EnVision dev team
* Provisioning service QE


## Consequences

* This will set style for our backend app development and get all devs on the same page in the tech stack.
* We ease up the transition from monolithic approach as the dev team relearns tech stack from scratch
