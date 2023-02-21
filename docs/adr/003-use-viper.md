# 3. Use Viper for configuration

Authors: Ondrej Ezr

## Status

Superseded by [8. Config environment files](008-config-env-files.md)


## Problem Statement

We as a team want to have easy way to manage configuration in all environments.

## Goals

* Make our configuration as understandable as possible
* Have clear list of possible configuration options
* Define hierarchy of value sources
* Define how we are loading configuration in Clowder environment.


## Non-goals

* We do not want to implement new configuration solution


## Current Architecture

* Environment Variables with .env files
* Dependencies defined in PROV00001, where configuration dependency is not considered


## Proposed Architecture

* Use viper as it is tool used in all the teams and can easily support cmd line flags and options
* Use github.com/RedHatInsights/app-common-go package to load Clowder configuration


## Challenges

* Viper is additional dependency that will take a bit time to get familiar with and needs to be maintained.


## Alternatives Considered

* Keeping the environment variables and override manually from Clowder
  * Ruled out because ConsoleDot does not use env variables for configuration, we have no need to use them.

## Dependencies

* Brief configuration documentation in README.md in `github.com/RHEnVision/provisioning-backend`

## Stakeholders

* Lukas Zapletal


## Consequences

* _Nothing of significance_
