# Architecture Decision Records

This folder keeps all important architecture decisions for Provisioning service backend.
ADR is a well established process in software engineering and more details can be found at [https://adr.github.io/](https://adr.github.io/).

### ADR lifecycle management

* **Key decider** - person driving the decision
* **Stakeholder** - everyone affected by the decision

1. ADR is based on [Template](000-template.md) by key decider
   1. _Tip:_ draft can be in Google doc, if it helps authoring
2. _Optional_ ADR is merged in proposed state
   1. Useful if implementation spans multiple PRs and we are not sure about final implementation
3. PR implementing the decision will come with ADR in Accepted status.
   1. _Note_ or change the status to Accepted if it already exists
   2. Merging a PR means accepting the ADR
