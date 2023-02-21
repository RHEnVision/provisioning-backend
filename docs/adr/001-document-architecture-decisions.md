# 1. Document Architecture decisions using Architecture decision records

Authors: Ondrej Ezr


## Status

Proposed


## Problem Statement

We want to capture our architecture decisions to convey them to all developers.
This capture should make it easy for people to find a given decision and refer to it.
The decision context should be captured, so we can question that decision when changes arise.


## Goals

* Capture key decisions
* Let the team rely on these decisions
* Capture decision context so we can revisit when context changes
* Be easily discoverable
* Align with https://github.com/operate-first/sre best practices.


## Non-goals

* Does not aim to achieve perfection from the get go


## Current Architecture

* None


## Proposed Architecture

* Embrace a well formed template for each Architecture Decision
  * [[Template] Architecture Decision Record](000-template.md)
* Store each significant Architecture Decision Record (ADR) in `adr` folder.
  * if decisions involve just specific team, those are shared in subfolders.
* Each Architecture Decision Record file name is formatted as `###-<adr-title>`
  * `###` is a quasi monotonic number assigned manually at ADR creation time
* ADRs are immutable once their status moves to accepted
* ADRs can supersede preceding ADRs: this is how a decision is changed.
* No ADR once moved past Draft can be ‘deleted’ or ‘removed’.
* Drafts can be developed as pull requests or using gdocs.
  However once they reach Proposed, they must become PRs and be commited to this repository.


## Challenges

* Since ADRs are numbered at creation time, it is possible that an ADR with a higher number is in effect before an ADR with a lower number. Shouldn't be a problem as long as they don’t contradict (nor supersede) one another.


## Alternatives Considered

* Document at gdocs only, but that is not easily discoverable and is hard to maintain in time.


## Dependencies

* [[Template] Architecture Decision Record](000-template.md)
* Architecture Decision Records directory `adr/`
* Some materials on the concept of ADR
    * [https://cognitect.com/blog/2011/11/15/documenting-architecture-decisions.html](https://cognitect.com/blog/2011/11/15/documenting-architecture-decisions.html)
    * [https://adr.github.io/](https://adr.github.io/)
    * [https://adr.github.io/madr/](https://adr.github.io/madr/)
    * [https://engineering.atspotify.com/2020/04/14/when-should-i-write-an-architecture-decision-record/](https://engineering.atspotify.com/2020/04/14/when-should-i-write-an-architecture-decision-record/)
    * Note that we are deviating from these descriptions in some areas (including the template)

## Stakeholders

* Product owners of management services


## Consequences

Decisions will be recorded and documented at cost of slowing down the process of making them.
