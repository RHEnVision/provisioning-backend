# 8. Template

Authors: Lukas Zapletal

## Status

Accepted

Supersedes [3. Use viper](003-use-viper.md)

## Problem Statement

The viper tool is quite complex and not very customizable.
It tries to achieve configuration from many sources, but we only use file in dev and environment variables everywhere else.
It is hard to manage Viper for Environment variables which are the main source of configuration.

## Goals

* Simplify our config management


## Non-goals

* Implement Viper features ourselves


## Current Architecture

* Using Viper


## Proposed Architecture

Complete rewrite of our config handling.
Configuration is now done via `config/api.env` (pbapi), `config/worker.env` (pbworker) and configs/test.env (integration tests).
There is also configs/example.env file which is generated via a new makefile target `generate-example-config`.

All config management is now handled in a single file [`config.go`](internal/config/config.go).
The whole structure, types, env variable name, default value and also optionally description.

We also decided to use singular for config directory name renaming it from `configs` to `config`.

## Challenges

* New configuration needs to be deployed for development setups. There should be no changes needed for Clowder, except if we have a typo in an environmental variable.


## Alternatives Considered

* Hacking Viper, but we decided quickly environment files are better option.

## Dependencies

* _None_


## Stakeholders

* EnVision developers as the change only affects dev environments.

## Consequences

* lighter dependency tree
* All devs need to change from yaml configuration file to env file
  * This should be straight forward as we do not have many configuration variables.
