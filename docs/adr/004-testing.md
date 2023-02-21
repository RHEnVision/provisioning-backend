# 4. Unit and itegration testing

Authors: Ondrej Ezr


## Status

Accepted

## Glossary

* **Unit** test - tests only one function at a time, tests live in the same repository as code
* **Integration** test - tests one component as a whole, tests live in the same repository as code
* **System** test - tests the whole system (all components together), tests live in separate repository


## Problem Statement

The application should simplify refactoring and adding features by covering existing code by tests.
These tests should make sure the application is working properly.
Tests also help document the code for people who do not know what the code is supposed to do (new people, people who have not touched the code area yet).


## Goals

* Define test strategy
* Cover at least happy path and corner cases with unit tests
* Incentivize good interfaces definition
* Fast feedback loop even if we grow into a medium sized application


## Non-goals

* System tests are out of scope of this ADR as those are worked on by QE


## Current Architecture

* None


## Proposed Architecture

* Use [testify](https://github.com/stretchr/testify/#suite-package) that is a wrapper around basic go test tooling
    * These wrappers simplify mocking and asserting
    * It allows to define test suites
* Stub out database and other external dependencies that the function use, but are not subject of the test
* Aim for TDD to get the best coverage
    * Try to follow best practices for go TDD: [https://quii.gitbook.io/learn-go-with-tests](https://quii.gitbook.io/learn-go-with-tests)
* Cover the app with integration tests
    * Just integration of the backend application
    * These should cover only happy path and still have stubbed out external dependencies


## Challenges

* Consent on full unit test coverage is not easy as it is hard to see the value of it.


## Alternatives Considered

* Only Integration tests, unit test only for corner cases
    * Does not document the code and thus does not make it easier to refactor code
    * Integration tests tend to be slow and thus might force testing only the main cases
    * Full coverage is usually not achieved here


## Dependencies

* [CONTRIBUTING.md](https://github.com/RHEnVision/provisioning-backend/blob/main/CONTRIBUTING.md)


## Stakeholders

* EnVision developers


## Consequences

* Tests will document intended use of the code written by us
* Proper mocking incentivize good interface definitions
