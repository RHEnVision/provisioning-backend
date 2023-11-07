# 13. Add Provisioning Saga Workflow

Authors: Adi Abramovich

## Status

Proposed

## Problem Statement

This ADR exists because there is a need of making the build and launch workflow more user-friendly.
This ADR will include four propositions regarding how should we implement this feature and what should we take into consideration.

## Goals

* Make the saga workflow fast, secure, and simple
* Decide on the solution architecture

## Non-goals

## Current Architecture

Currently, the Image Builder wizard and the Provisioning Wizard form disconnected user experience.

After collecting the information needed in order to create an image, the Image Builder's wizard closes and the user is redirected to Image Builder's main page which lists all the images that were already created or pending to be.

Only once the image is created successfully, a **"Launch"** button appears and allows
the user to go through the Provisioning experience.

## Proposed Architecture

I have thought about four architecture options where all four options have the first steps in common:

* Image builder collects all the information needed to create an image on a specific hyper-scaler
* Pressing the **Build and launch** button on the wizard
* Validate the build manifest
* A background job is running in order to create an image
* The Provisioning wizard opens up

Important information to remember:

* The Provisioning is dependent on the image id in order to continue and finish the launch
* The shared sources are given directly from the image variable from IB Wizard

**In the following section I will elaborate on the differences:**

1. Combining Image Builder wizard with the Provisioning wizard and showing the status of the instance on Image Builder main page:
   * Add a polling to Image Builder for image readieness.
   * The customer fills in the data needed in order to launch an image.
   * The wizard closes, the polling continues, and when it is done a job is enqueued with the image id.
   * The status of the job is presented on Image Builder main page. (like the status is being presented in the wizard now)

2. Combining Image Builder wizard with the provisioning wizard and adding a step in the Provisioning wizard creating the image:
   * The customer fills in the data needed in order to launch an image.
   * There is an additional step before all the others in the Provisioning wizard showing the status of the image build.
   * After the image has been created the Provisioning workflow continues(creating a reservation, uploading pubkeys when needed and launching the image)

3. Combining Image Builder wizard with the provisioning wizard using polling and showing the status of the instance on Inventory main page:
   * Add a polling mechanism in provisioning backend questioning Image Builder if the image is ready.
   * The customer fills in the data needed in order to launch an image.
   * The wizard closes, the polling continues, and when it is done a job is enqueued with the image id.
   * The status of the job is presented on Inventory main page. (similar to how is the status presented in the Launch wizard currently, with additional step from the solution `2.`)

4. Combining Image Builder wizard with the provisioning wizard using Kafka and showing the status of the instance on Inventory main page:
   * The customer fills in the data needed in order to launch an image.
   * The wizard closes, the Provisioning is being informed by Image Builder when the image has been created successfully (probably using Kafka), then the job is enqueued with the corresponding image id.
   * The status of the job is presented on Inventory main page. (similar to how is the status presented in the Launch wizard currently, with additional step from the solution `2.`)

## Challenges

* Integration with Inventory (third option)
* Integration with Image builder which will require many adjustments and changes on Image builder side and on ours as well. (first option)
* Integration Kafka/Notifications with the Provisioning service and Image Builder.(Fourth option)

## Alternatives Considered

* Mentioned above

## Dependencies

* Integration with Inventory
