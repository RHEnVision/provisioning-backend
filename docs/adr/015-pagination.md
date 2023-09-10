# 11. Pagination

**Authors:** Adi Abramovich

## Status

**Status:** Accepted

## Problem Statement

The Launch Service consists of several endpoints designed to provision custom images to the cloud. In certain scenarios, it becomes necessary to list resources such as Pubkeys, Reservations, Instance/Machine types, sources, and launch templates. To ensure proper presentation and enable pagination on the frontend, backend pagination implementation is required.

## Goals

* Implement pagination support for all list endpoints: Launch Templates, Sources, Pubkeys, Reservations, and Instance Types.

## Current Architecture

There is currently no specific architecture in place for pagination.

## Proposed Architecture

Pagination will be implemented through middleware, and this middleware will be selectively used on specific endpoints listed above.

Two use cases need to be considered:

1. **Postgres, Internal Services (e.g., Sources):**
   For this case, limit and offset information will be added as query parameters. If not provided, default values will be used as mutually decided: offset equals 0, and limit equals 100. To leverage pagination, the following mechanisms will be employed:
   - PostgreSQL built-in pagination functionality: Usage of Limit and Offset built-in variables.
   - For sources, optional parameters (Limit and Offset) will be utilized when using their client.

2. **External Services (e.g., AWS, GCP):**
   The Limit and Token parameters will be added as query parameters. The token is provided by external services and leads to the next page when given. When the token is empty, it indicates the need to list the first page.

Adding the Metadata Structure:
When querying one of the list endpoints, a metadata structure will appear in the response. It holds the total number of records for the first option and links for the next and previous pages based on the current limit provided.

Frontend: 
To handle pagination on the frontend, we will utilize React Query along with the wizard context / useState to manage the offset and Limit. We will leverage React Query's built-in pagination mechanism. [Learn more about React Query pagination](https://tanstack.com/query/v4/docs/react/guides/paginated-queries).

## Challenges

- The necessity of making breaking changes in the API by nesting all the list results under a data root variable.
- Additional features required to provide proper pagination, such as ordering and filtering.

## Alternatives Considered

- Implementation of pagination using Redis.

## Stakeholders

* EnVision developers

## Consequences

By implementing pagination, we aim to enhance the user experience and ensure efficient resource listing. It will enable users to navigate through large datasets with ease, providing a more manageable and user-friendly interface.
