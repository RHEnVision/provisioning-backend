# 11. AWS Permissions Check

Authors: Anna Vitova


## Status

Accepted


## Problem Statement

We provide documentation on how to set up Amazon AWS roles. For customers using these suggested policies, the application should not fail due to insufficient permissions. In case of errors because of the missing permissions, it would be nice to have a way to list them.
As explained below, checking for missing permissions with 100 % accuracy might be impossible. False positives, though, could be avoided.

## Goals

* Check if permissions are set as recommended for the used role.
* Not generate false positives of missing permissions.


## Non-goals

* Check the permissions if they are set in other ways than recommended.


## Current Architecture

* None.


## Proposed Architecture

* The list of recommended permissions would be hardcoded into the backend to be easily accessible.
* First, we fetch all standalone policies for the used role. There could be more policies than the one we recommended, so we will go through every one of them. In case the customer splits out our recommended policies into more, we will look through every policy declared as "Effect": "Allow" and with the field "Action". Both Statement and Action fields could be defined as either array or an object, both should be supported by the checking endpoint.
* On the other hand, what will not be supported are complementary policies (NotAction, Effect: Deny; explained below). In these cases, adding our recommended policies might not have any effect, and we will not show that. Also, checking for the Resource value will not be supported, which could lead to false positives.
* For potentially missing permissions, check if these are not in the inline policy.


## Challenges

* AWS has a variety of options to set policies for roles. These unique ways to add policies end up with different structures of the response. The same field can be of [various types](https://docs.aws.amazon.com/IAM/latest/UserGuide/reference_policies_grammar.html). Thus, it is impossible to determine the structure beforehand and unmarshal the policies into a predefined struct.
* When choosing a resource for which the permissions are to be applied, it is possible to use wildcards. A check of the policy's relevancy when the resource is more specific is complicated.
* Wildcards could also be used for defining actions. The following rule can be applied to obtain permission to delete anything: "Action": "ec2:*Delete*. This adds all the Delete permissions for ec2.
* Some policies can also be described as a complement to stated permissions. It is possible to express permissions as either "Action": "ec2:StopInstances" or "NotAction": "ec2:StopInstances". One adds the stated permission. The other adds everything except the one. The obstacle in this and the previous case is the absence of listing all possible existing permissions in AWS.
* There are two contradictory approaches when describing what to do with permissions. The permissions could be defined with either "Effect": "Allow" or "Effect:" "Deny". We could find a relevant policy with all required permissions, but there could be a different one denying some of the permissions. Detection of this and determining the priority of these conflicting permissions would be needed to cover the permission check fully.
* Other provided way to add policy is the one called inline policy. Our recommended way, a standalone policy, has its own ARN, whereas an inline policy is linked right to the role.

## Alternatives Considered

* Searching for one specific policy to compare with the recommended one, i.e. searching for one based on Sid/policy name, was considered. This approach could lead to false positives, as the permissions missing in this specific policy might be stated elsewhere.
* Accepting only one possible structure of the JSON - the one we recommend - would allow us to map the JSON right away on a struct. This approach could create false positives if the structure is different, but the required permissions are defined.

## Stakeholders

* EnVision developers


## Consequences

* After changing the required permissions, the hard-coded permissions need to change.
