<a name="0.11.0"></a>
## [0.11.0](https://github.com/RHEnVision/provisioning-backend/compare/0.10.0...0.11.0) (2022-11-21)

### Bug Fixes
- **config:** correct Unleash URL prefix ([bd6ab5a](https://github.com/RHEnVision/provisioning-backend/commit/bd6ab5a02d317a51e9a5a7ca742bbd372b2807bf))
- **config:** guard for non-exixtend kafka config ([a5b3d9c](https://github.com/RHEnVision/provisioning-backend/commit/a5b3d9c552953f3ddb7824c67712719b7a83bd27))
- **config:** unleash token as bearer header ([3bb424c](https://github.com/RHEnVision/provisioning-backend/commit/3bb424c889b8b84d8b982cae6e24ea9af1a927ba))
- **logging:** Disable middlewares for status routes ([811905d](https://github.com/RHEnVision/provisioning-backend/commit/811905dcf89335174a173f4362892c0f0931dce3))
- **reservation:** generic reservation by id ([5131c7b](https://github.com/RHEnVision/provisioning-backend/commit/5131c7b08c4164dcb11cb93ecb55916665132ccc))
- missing cache type variable for api ([0ed8c1b](https://github.com/RHEnVision/provisioning-backend/commit/0ed8c1bc8ade1d3b8346240855adcb76d5ab5a3f))
- null for aws_reservation_id when pending ([eb5e353](https://github.com/RHEnVision/provisioning-backend/commit/eb5e353d2d17541331bc460d587b95c48c15a75d))
- print full errors in logs ([7cc2e10](https://github.com/RHEnVision/provisioning-backend/commit/7cc2e10181bd623549b6ff78d03480e82c47bff3))

### Build
- **clowder:** Add image builder as clowder dependency ([01b1c06](https://github.com/RHEnVision/provisioning-backend/commit/01b1c06a9d59ff1334035699f46ae0915e4ac430))

### Chore
- add direnv.net to gitignore ([9ea4d84](https://github.com/RHEnVision/provisioning-backend/commit/9ea4d84764eb44e56716a07d6a2d478956d7b94d))
- update all deps ([e3f5b54](https://github.com/RHEnVision/provisioning-backend/commit/e3f5b54d505217e12b1a49077651215d80a21fc5))
- version 0.10.0 and changelog ([b7314fd](https://github.com/RHEnVision/provisioning-backend/commit/b7314fdc007cc7c4fc94dfde4871bb5c868a59d5))

### Features
- **kafka:** setup, configuration, availability check ([971c64d](https://github.com/RHEnVision/provisioning-backend/commit/971c64d37d62787778e4ed476cbcc3255ba8f6bd))
- **refactor:** Add required true post aws in apispec ([73f3956](https://github.com/RHEnVision/provisioning-backend/commit/73f3956b2a08914d3adb2628b4addb6836a85941))
- availability check kafka topic ([cea9ea0](https://github.com/RHEnVision/provisioning-backend/commit/cea9ea0e6aaac9a71f9d2171d38fb73ab1c27eb6))
- reservation detail returns instance ids ([30a1f8e](https://github.com/RHEnVision/provisioning-backend/commit/30a1f8ec766a2da552c27da70c39c2b7b8863111))


<a name="0.10.0"></a>
## [0.10.0](https://github.com/RHEnVision/provisioning-backend/compare/0.9.0...0.10.0) (2022-11-02)

### Bug Fixes
- **changelog:** correct upstream links ([d34b161](https://github.com/RHEnVision/provisioning-backend/commit/d34b16184094fd51b11d1ad084b8a1201b910291))
- **ec2:** typo in etag prefix ([ea86eb2](https://github.com/RHEnVision/provisioning-backend/commit/ea86eb23c48c78fbd894100207aa4bfb59726030))
- **logging:** no colors for clowder ([cff318a](https://github.com/RHEnVision/provisioning-backend/commit/cff318a0ece63abe253915239e72fc9d51faf9ff))
- **queue:** improve errors and logging for dejq init ([dff6b11](https://github.com/RHEnVision/provisioning-backend/commit/dff6b11c7a17cc5f044d71064b83b9dbb00542c3))
- **queue:** recognize unknown worker config values ([a4114e8](https://github.com/RHEnVision/provisioning-backend/commit/a4114e8c304ebd9922d8da5e5c01b0962bf0c72c))
- **services:** refactor errors ([c1b05f4](https://github.com/RHEnVision/provisioning-backend/commit/c1b05f44a60580e3618446a19598c6a161401337))
- Add error messages ([e937d94](https://github.com/RHEnVision/provisioning-backend/commit/e937d948ad8faaddc5a4632cabbe93926360c412))
- remove DAOInit Error ([5800a20](https://github.com/RHEnVision/provisioning-backend/commit/5800a20aa4f6a768c6b533eea752fd010c88bb25))
- throw an error when pubkey is duplicated (HMSPROV-309) ([5a5078d](https://github.com/RHEnVision/provisioning-backend/commit/5a5078dd355b39d17bf2fc7a11ffc20f204e4ee8))

### Build
- **clowder:** Add full image builder URL ([ba4c795](https://github.com/RHEnVision/provisioning-backend/commit/ba4c795b25a1009c56f438e0f050e41e675ff843))
- **deploy:** Allow passing image builder url ([e5d85bc](https://github.com/RHEnVision/provisioning-backend/commit/e5d85bcf9564e43930b93d702b93baeda794803e))

### Chore
- **azure:** refresh types ([5020551](https://github.com/RHEnVision/provisioning-backend/commit/5020551e17ea64c9985c9b17b2e2398ade83423c))
- **errors:** change payload structure ([19d31bf](https://github.com/RHEnVision/provisioning-backend/commit/19d31bf974a9442d12f1119171ab93dc288e6392))
- add JetBrains Fleet to gitignore ([304f719](https://github.com/RHEnVision/provisioning-backend/commit/304f71943b638c69672ccd8917851f30995b8f9a))
- allow local seed scripts ([e59d597](https://github.com/RHEnVision/provisioning-backend/commit/e59d59702e03b1e6fa2cd6207e7730bd0dc3518b))
- change gitignore for sources ([699b2eb](https://github.com/RHEnVision/provisioning-backend/commit/699b2ebe7fdf2ba3bc0bc69dde32b1b1a604d07a))

### Code Refactoring
- **azure:** rename client type ([c7ec156](https://github.com/RHEnVision/provisioning-backend/commit/c7ec156c79d26ae1d1baa51c1871589ff839a7cb))
- **azure:** split service and cust ifaces ([924216b](https://github.com/RHEnVision/provisioning-backend/commit/924216b3e3eb5b8c1680116346889771d23152a6))
- **changelog:** different tool for changelog ([febcb38](https://github.com/RHEnVision/provisioning-backend/commit/febcb38efffad06ab590d848b52cce79059100f3))
- **clients:** change arn to authentication for generic purpose ([5b19411](https://github.com/RHEnVision/provisioning-backend/commit/5b19411b4fbef8030867fbfc4e49eecc0c65a56b))
- **clients:** remove Customer prefix from EC2 client ([2ec34db](https://github.com/RHEnVision/provisioning-backend/commit/2ec34dbc9867b72f90fb20a44f53a23f29350f0a))

### Docs
- **readme:** Add info for creating an image for gcp ([02c171a](https://github.com/RHEnVision/provisioning-backend/commit/02c171adabec92ab1fa998fc20f02bf1d6256912))

### Features
- **cache:** add Redis cache and queue ([1d60e0f](https://github.com/RHEnVision/provisioning-backend/commit/1d60e0f59a8e2f01fda1451546f36dd3044effa3))
- **clients:** Generated machine types and types per zone ([a654baa](https://github.com/RHEnVision/provisioning-backend/commit/a654baa0132724f6c1a5c34e8d3efdad9c826344))
- **clients:** Preload machine types for GCP ([741b820](https://github.com/RHEnVision/provisioning-backend/commit/741b8207c6ee22d95dc33c84ef001ac413f5e68a))
- **flags:** feature flags endpoint ([727a72f](https://github.com/RHEnVision/provisioning-backend/commit/727a72f86e5f0e89a9166ebb5dc74b5545fdfb27))
- **gcp:** add gcp request payload ([0c718cc](https://github.com/RHEnVision/provisioning-backend/commit/0c718cc6321fd965d055288514808f9b66ce110e))
- welcome HTML page ([a2fe12c](https://github.com/RHEnVision/provisioning-backend/commit/a2fe12cdde2d6c3525cb1e706dfc27e6ce2c4d36))


<a name="0.9.0"></a>
## [0.9.0](https://github.com/RHEnVision/provisioning-backend/compare/0.8.0...0.9.0) (2022-10-17)

### Bug Fixes
- **changelog:** messages without subtypes ([c798bb0](https://github.com/RHEnVision/provisioning-backend/commit/c798bb05455fa1721806ba07fe819e1aaadd6952))

### Chore
- **changelog:** GHA action for commits ([293b33b](https://github.com/RHEnVision/provisioning-backend/commit/293b33b24ce5c2b85ca53f7b147d6f30c0ad95f2))
- **changelog:** introduced changelog generator ([0f0596f](https://github.com/RHEnVision/provisioning-backend/commit/0f0596fec334070a96b593b120e1d1926ca10f8a))
- **changelog:** update ([01b8f5f](https://github.com/RHEnVision/provisioning-backend/commit/01b8f5f8daa15137a290e6304eeeebeb3268e76f))
- **oapi:** update clients, missing make doc ([96f00bf](https://github.com/RHEnVision/provisioning-backend/commit/96f00bf32c4eac694d85a7c3eb3f9560410a2814))
- **readme:** missing pbworker info ([c237de8](https://github.com/RHEnVision/provisioning-backend/commit/c237de871de7fea9c9b7260a20a886d731d8a067))
- **scripts:** update sources seed script ([735ea16](https://github.com/RHEnVision/provisioning-backend/commit/735ea16f2daff678183aa925288b2cc3d460698d))

### Code Refactoring
- **clients:** HTTP transport improved ([70f2c07](https://github.com/RHEnVision/provisioning-backend/commit/70f2c07b94afa001a3aceef5a33a8aea0f77d1d0))

### Docs
- **changelog:** developer instructions ([fed4d45](https://github.com/RHEnVision/provisioning-backend/commit/fed4d457d1f4db12b76927c7be8e6e7dd0ba6038))

### Features
- delete pubkey resource ([a148c89](https://github.com/RHEnVision/provisioning-backend/commit/a148c89124987b9c336d4340b93565098a768cbd))


<a name="0.8.0"></a>
## [0.8.0](https://github.com/RHEnVision/provisioning-backend/compare/0.7.0...0.8.0) (2022-10-11)


<a name="0.7.0"></a>
## [0.7.0](https://github.com/RHEnVision/provisioning-backend/compare/0.6.0...0.7.0) (2022-09-15)


<a name="0.6.0"></a>
## [0.6.0](https://github.com/RHEnVision/provisioning-backend/compare/0.5.0...0.6.0) (2022-09-08)


<a name="0.5.0"></a>
## [0.5.0](https://github.com/RHEnVision/provisioning-backend/compare/0.4.0...0.5.0) (2022-08-31)


<a name="0.4.0"></a>
## [0.4.0](https://github.com/RHEnVision/provisioning-backend/compare/0.3.0...0.4.0) (2022-08-22)


<a name="0.3.0"></a>
## [0.3.0](https://github.com/RHEnVision/provisioning-backend/compare/0.2.0...0.3.0) (2022-08-15)


<a name="0.2.0"></a>
## [0.2.0](https://github.com/RHEnVision/provisioning-backend/compare/0.1.0...0.2.0) (2022-07-27)


<a name="0.1.0"></a>
## [0.1.0](https://github.com/RHEnVision/provisioning-backend/compare/9d638e99279166b7f27e14feb9468e9b7c98a390...0.1.0) (2022-06-15)


