<a name="0.14.0"></a>
## [0.14.0](https://github.com/RHEnVision/provisioning-backend/compare/0.13.0...0.14.0) (2023-01-25)

### Bug Fixes
- **HMSPROV-387:** filter out noisy kafka logs ([e0a7b21](https://github.com/RHEnVision/provisioning-backend/commit/e0a7b216744ae4232a6474a795c9ae2967eb99a6))
- **HMSPROV-387:** use time-based offset for statuser ([bbc59a9](https://github.com/RHEnVision/provisioning-backend/commit/bbc59a9cfc9960afe52db731507e761d6e4e2746))
- **HMSPROV-390:** RSA fingerprint and migration ([e115286](https://github.com/RHEnVision/provisioning-backend/commit/e115286767281b9401c384c311feeda8820ca588))
- **HMSPROV-390:** RSA fingerprint and migration ([6ddf275](https://github.com/RHEnVision/provisioning-backend/commit/6ddf275f49f3c947a4ea7357b0a66808b0ac131e))
- **HMSPROV-390:** unscoped update pubkey fix ([bd30ea8](https://github.com/RHEnVision/provisioning-backend/commit/bd30ea85cd57f19e77390f209c375e72d833eb33))
- **HMSPROV-425:** incorporate dejq into the app ([f8a0b6f](https://github.com/RHEnVision/provisioning-backend/commit/f8a0b6f5fb5e767d2d0d57f1f11f892c7f014946))
- **HMSPROV-425:** incorporate dejq into the app ([bf11c14](https://github.com/RHEnVision/provisioning-backend/commit/bf11c14ed178bd98727eb348647ea63732f38f14))
- **HMSPROV-425:** recover panics in workers ([2a560d5](https://github.com/RHEnVision/provisioning-backend/commit/2a560d524040312ba3984197b8902dce3f1b1007))
- **HMSPROV-433:** change resource type to application ([72d26d1](https://github.com/RHEnVision/provisioning-backend/commit/72d26d16f5324a0ac5b66d4558acc6a63f86c67c))
- filtering Provisioning auth for Source ([cebbdb3](https://github.com/RHEnVision/provisioning-backend/commit/cebbdb3df4d21826fc7fdcc66c7fb33939b53e11)), related to [HMSPROV-426](https://issues.redhat.com/browse/HMSPROV-426)
- image builder clowder config ([29aa59d](https://github.com/RHEnVision/provisioning-backend/commit/29aa59d95cdb1b2fdc99731dcc44d78085932303)), related to [HMSPROV-421](https://issues.redhat.com/browse/HMSPROV-421)
- unique index on pubkey_resource ([2b68b0a](https://github.com/RHEnVision/provisioning-backend/commit/2b68b0a4240b5ec2e3d77ce64a5ce92292f65097)), related to [HMSPROV-415](https://issues.redhat.com/browse/HMSPROV-415)

### Build
- Change unleash prefix to only app name ([2932b02](https://github.com/RHEnVision/provisioning-backend/commit/2932b020634011a7f305df8139e71dfe7b4148c8))

### Chore
- changelog for 0.13.0 ([4a5b30e](https://github.com/RHEnVision/provisioning-backend/commit/4a5b30ef6c0a1a18deecaa28d06c9e4cedff2c57))
- document and split ctxval package ([f6e8db3](https://github.com/RHEnVision/provisioning-backend/commit/f6e8db38907220c4c911a666ff67bf440ec643ff))
- improve logger middleware to better visibility ([24919ed](https://github.com/RHEnVision/provisioning-backend/commit/24919edf2fbbc3d79ce7dc82c2fb63c00233ff0b))
- move EnforceIdentity so it logs through our logger ([603d4f6](https://github.com/RHEnVision/provisioning-backend/commit/603d4f675bb9ba7e06e7e384749f2b7bf5148f1c))
- optimize ctxval getters ([f73e3a4](https://github.com/RHEnVision/provisioning-backend/commit/f73e3a4ccc78233c92a797957366df0c9944f7b8))
- prevent kafka SASL creds log leak ([793746c](https://github.com/RHEnVision/provisioning-backend/commit/793746c4ce4c3af36350e8a411ce11a93958ec46))
- run migrations for statuser ([f5b2e19](https://github.com/RHEnVision/provisioning-backend/commit/f5b2e1992b1987bdb7c95685f7ae2e607a420990))
- simplify error payload logging ([542f74e](https://github.com/RHEnVision/provisioning-backend/commit/542f74e7965a70c6e2b560e3ad21afb342d5e9e8))
- use reservation_id instead reservation ([65e5ac3](https://github.com/RHEnVision/provisioning-backend/commit/65e5ac3dfee8b25c561a4fbe9fe90ff2fa495739))

### CI
- Add testing.yaml for IQE CJI ([a978bda](https://github.com/RHEnVision/provisioning-backend/commit/a978bda0702b238329cfb9bf94ab59e82753272a))

### Code Refactoring
- Add numeric status code ([0c4591e](https://github.com/RHEnVision/provisioning-backend/commit/0c4591ecb299ea9e42c3036508bf74145966c427))

### Docs
- Safe zero downtime upgrades ([75caf69](https://github.com/RHEnVision/provisioning-backend/commit/75caf695b7dbffb3a490916f83b35c399fa12346)), related to [HMSPROV-386](https://issues.redhat.com/browse/HMSPROV-386)

### Features
- **HMSPROV-428:** Add provisioning dashboard ([a93a9f5](https://github.com/RHEnVision/provisioning-backend/commit/a93a9f5589c3400f62fa86d1d93c66859eaf1f4e))


<a name="0.13.0"></a>
## [0.13.0](https://github.com/RHEnVision/provisioning-backend/compare/0.12.0...0.13.0) (2023-01-12)

### Bug Fixes
- **HMSPROV-170:** change topic to platform.sources.status ([b63ff7a](https://github.com/RHEnVision/provisioning-backend/commit/b63ff7aaffd16b546e2f7b3c8fc7183977deab09))
- **HMSPROV-340:** nice error on arch mismatch ([ca5d32e](https://github.com/RHEnVision/provisioning-backend/commit/ca5d32e61105a77ce86ba783a240afad4340e701))
- **HMSPROV-345:** Add event_type header and resource type to kafka msg ([28a69da](https://github.com/RHEnVision/provisioning-backend/commit/28a69da29e676801de84ff15d978ddb9103ba4fd))
- **HMSPROV-345:** change to Source ([559bcfe](https://github.com/RHEnVision/provisioning-backend/commit/559bcfe2ad607e37875906db225229af64359821))
- **HMSPROV-345:** remove default tag ([5659d0d](https://github.com/RHEnVision/provisioning-backend/commit/5659d0d520402b721d5c6a3b2a20f95dbc285183))
- **HMSPROV-352:** error out jobs early ([a85152a](https://github.com/RHEnVision/provisioning-backend/commit/a85152a747535acdd94d4b784fdd3f7600918d73))
- **HMSPROV-352:** improve error message ([0dd032f](https://github.com/RHEnVision/provisioning-backend/commit/0dd032fa183edba623f47a24f62aa11e67e59b63))
- **HMSPROV-368:** change version to BuildCommit ([f3991bb](https://github.com/RHEnVision/provisioning-backend/commit/f3991bb90a76034a9aa29d141f0f0ed340253a1e))
- **HMSPROV-387:** set consumer group for statuser ([bd800fa](https://github.com/RHEnVision/provisioning-backend/commit/bd800fa8769b8abaed84cea944c42d1a07135803))
- **HMSPROV-389:** drop memcache count ([bfc9a00](https://github.com/RHEnVision/provisioning-backend/commit/bfc9a00c1229fc3d64fdfe129e402edc4569e6de))
- **HMSPROV-390:** calculate fingerprint for AWS ([df98843](https://github.com/RHEnVision/provisioning-backend/commit/df988435e24621496e170e5fe349b6af8b4096f6))
- **HMSPROV-392:** Check if image is an original or a shared one ([a03bde2](https://github.com/RHEnVision/provisioning-backend/commit/a03bde289659702e504335f5d837f50914a9a1f3))
- **HMSPROV-399:** add dejq job queue size metric ([6e49fcc](https://github.com/RHEnVision/provisioning-backend/commit/6e49fcc78c016c5e5b74e29430151773fed0bae6))
- **HMSPROV-407:** disable cw for migrations ([dbc1ea4](https://github.com/RHEnVision/provisioning-backend/commit/dbc1ea48c24d32da5c0590c11c0bb1a1aa763873))
- **HMSPROV-407:** enable cloudwatch in clowder ([cf5663b](https://github.com/RHEnVision/provisioning-backend/commit/cf5663bfe00677a84e3c4bff5ad05f6b520e5fae))
- **HMSPROV-407:** enable cloudwatch in clowder ([e4a0083](https://github.com/RHEnVision/provisioning-backend/commit/e4a0083fd6a3a9d7a32f644af5230b9e6d9cdbc7))
- **HMSPROV-407:** fix blank logic in cw initialization ([9de1b76](https://github.com/RHEnVision/provisioning-backend/commit/9de1b76abab17175fde3c7ae786376ba1e50d1e3))
- **HMSPROV-407:** fix cw config validation ([e661ac8](https://github.com/RHEnVision/provisioning-backend/commit/e661ac8f46fc314765ed383e0ddac853112c8961))
- **HMSPROV-407:** further improve logging of cw config ([681ccb0](https://github.com/RHEnVision/provisioning-backend/commit/681ccb0d76aa3781e600197510b3dc67ce419a09))
- **HMSPROV-407:** improve logging of cw config ([7598917](https://github.com/RHEnVision/provisioning-backend/commit/7598917111d67c696cccb1de51bcda29e03fe900))
- **HMSPROV-414:** start dequeue loop in api only for memory ([3b2eaae](https://github.com/RHEnVision/provisioning-backend/commit/3b2eaae6411a8b90342da7dc30d54924bedf2c3a))
- ensure pubkey is always present on AWS ([2c310dd](https://github.com/RHEnVision/provisioning-backend/commit/2c310dd338b929619fee437a6da44f731f68e81e)), related to [HMSPROV-339](https://issues.redhat.com/browse/HMSPROV-339)
- Kafka headers are slice now ([5fdc3eb](https://github.com/RHEnVision/provisioning-backend/commit/5fdc3eb88ff827e8f28347829018fb6fe238bc7d)), related to [/HMSPROV-337](https://issues.redhat.com/browse/HMSPROV-337)
- Use correct topic for sources availability check ([3e7f820](https://github.com/RHEnVision/provisioning-backend/commit/3e7f8208394010c2ecada463f44284245f8470a6)), related to [HMSPROV-170](https://issues.redhat.com/browse/HMSPROV-170)
- utilize clowder topic mapping ([8c2ef77](https://github.com/RHEnVision/provisioning-backend/commit/8c2ef7752b5331e5fc456781428bd2398328cf3d)), related to [HMSPROV-343](https://issues.redhat.com/browse/HMSPROV-343)

### Build
- add sources kafka topics in clowdapp ([d279a89](https://github.com/RHEnVision/provisioning-backend/commit/d279a893196aeb583d6e768fae7f065abfd48f33))
- drop OpenAPI filepath env var ([7dce9ac](https://github.com/RHEnVision/provisioning-backend/commit/7dce9ac83e2073407d1fde29f03d06f7c7cd699c))
- Use external changelog module ([0cb3b3e](https://github.com/RHEnVision/provisioning-backend/commit/0cb3b3ee37a79a479536ee4f250ee529488480ae))

### Chore
- bump dejq version and use 2 workers ([b58ef29](https://github.com/RHEnVision/provisioning-backend/commit/b58ef292b8087e629be2924de6308adb7276d484))
- regenerate Azure types ([ea671bb](https://github.com/RHEnVision/provisioning-backend/commit/ea671bb43865fd1bc46fdd8bfe8e916e7ff1ba08))
- regenerate EC2 types ([3eb68a9](https://github.com/RHEnVision/provisioning-backend/commit/3eb68a9452c953df4807616e1a399a20e8e35f24))
- simplify replica config and set all back to 1 ([6621e0a](https://github.com/RHEnVision/provisioning-backend/commit/6621e0a3fa5c192a1480b7428286f559a446536d))
- update all dependencies ([2eadb7c](https://github.com/RHEnVision/provisioning-backend/commit/2eadb7c857af335b1b3f81a779db79ca6a9c3cef))
- update dependencies ([0de9ebb](https://github.com/RHEnVision/provisioning-backend/commit/0de9ebb3e761b640da4e3917f3ed187991da84d2))
- use one API worker ([c01814f](https://github.com/RHEnVision/provisioning-backend/commit/c01814f435449e7c44497015b853fcd5de465d0b))
- version 0.12.0 ([01c1cf7](https://github.com/RHEnVision/provisioning-backend/commit/01c1cf76f80dba38ea2a37645bd0595d28d92ea2))

### CI
- make fmt does not check commits ([3e22126](https://github.com/RHEnVision/provisioning-backend/commit/3e22126799e73657520c10dc202e206b48eee863))
- Update iqe markers for backend ([4e871fa](https://github.com/RHEnVision/provisioning-backend/commit/4e871fa48d17d85ec086daeab30afbf0d4a90d10))

### Docs
- adding info about trust relationship ([70fe8ed](https://github.com/RHEnVision/provisioning-backend/commit/70fe8eddc44624430d68774ece10215809a2b1e1))
- kafka and worker info ([c9ece9a](https://github.com/RHEnVision/provisioning-backend/commit/c9ece9a6b822286aa729def339e2593e71f1fd67))

### Features
- **HMSPROV-177:** Add availability check request duration metric ([5944953](https://github.com/RHEnVision/provisioning-backend/commit/594495321ddaa55f6802bdaeeefc7cb1b776c888))
- **HMSPROV-177:** Add total availability check request metric ([13d508d](https://github.com/RHEnVision/provisioning-backend/commit/13d508db6870a3ed93c1d8123f622cead6f22013))
- **HMSPROV-345:** Add source check availability per each provider ([c20ef14](https://github.com/RHEnVision/provisioning-backend/commit/c20ef14924e8c5599eebed364092574f0e875f0c))
- add account identity endpoint ([3df28fb](https://github.com/RHEnVision/provisioning-backend/commit/3df28fbbfc9ffc7ecfa4fa39deb00245d826b303)), related to [HMSPROV-357](https://issues.redhat.com/browse/HMSPROV-357)
- introduce availability status endpoint ([272f577](https://github.com/RHEnVision/provisioning-backend/commit/272f577b69852dbb4bd8a13a3e37fa81bf2c2e87)), related to [/HMSPROV-337](https://issues.redhat.com/browse/HMSPROV-337)

### Tests
- Add EC2 pubkey stub ([447d686](https://github.com/RHEnVision/provisioning-backend/commit/447d686d687559c0dd466ef6e4bd3c27e66b3374))
- use path to the main config dir ([1f4907d](https://github.com/RHEnVision/provisioning-backend/commit/1f4907d5200c01d81ae0dda3ce8ec87d765c9582))


<a name="0.12.0"></a>
## [0.12.0](https://github.com/RHEnVision/provisioning-backend/compare/0.11.0...0.12.0) (2022-12-01)

### Bug Fixes
- **kafka:** CA config from string ([34f4c59](https://github.com/RHEnVision/provisioning-backend/commit/34f4c59af28b52087e52dc2f12d8273bc0966f6a))
- **sources:** handle source without application type correctly ([75bc847](https://github.com/RHEnVision/provisioning-backend/commit/75bc8471cff0ffbb30908ad0bfec125fba3216e6))
- break availability queue sender loop on context cancel ([0ec4201](https://github.com/RHEnVision/provisioning-backend/commit/0ec420179b2372e875d829910123efcf5203fdd9))
- break consume look on context cancel ([4795c58](https://github.com/RHEnVision/provisioning-backend/commit/4795c58aac69e0e8dc262f53382428715ad06a9f))
- create topics in kafka startup script ([f7b2fab](https://github.com/RHEnVision/provisioning-backend/commit/f7b2fabce583884a94df6aa9974254c5ee20b42d))
- enable Kafka in Clowder ([8ea9023](https://github.com/RHEnVision/provisioning-backend/commit/8ea90236707a4f0cd080b48bd0f6aec6a2368deb))
- intermittent failures on CI for ASM queue test ([656ee05](https://github.com/RHEnVision/provisioning-backend/commit/656ee05b8d9634d671aff0067ea7b1dc8336a48d))
- kafka port is a pointer ([4ca3076](https://github.com/RHEnVision/provisioning-backend/commit/4ca30768ed981a938f6d8d1de7d2142597ef29a9))
- log topic alongside trace send message ([4f9a62a](https://github.com/RHEnVision/provisioning-backend/commit/4f9a62ac3c75612c4495615f641321ce9c7567ab))
- payload name not nullable ([996251d](https://github.com/RHEnVision/provisioning-backend/commit/996251d0c17dbf1b120d4c473a74f61854f77a61)), related to [HMSPROV-373](https://issues.redhat.com/browse/HMSPROV-373)
- scope existing pubkey search by source id ([af89244](https://github.com/RHEnVision/provisioning-backend/commit/af892449955be50118150015a6cf483c7d2ae97b)), related to [HMSPROV-366](https://issues.redhat.com/browse/HMSPROV-366)
- share HTTP transport across platform clients ([dcf7c38](https://github.com/RHEnVision/provisioning-backend/commit/dcf7c3894bfaac7d391437c175af0f75b9f31ead))

### Build
- **changelog:** Add support for Jira issues ([1235b4a](https://github.com/RHEnVision/provisioning-backend/commit/1235b4a29d7fe41ce66068eb53364ff7cb0bab2e))
- **changelog:** allow empty scope in changelog for feat and fix ([99c4e69](https://github.com/RHEnVision/provisioning-backend/commit/99c4e69175ad7e0583a48df4a7bfadf55325c638))
- **container:** Switch back to oficial go-tool build ([9930373](https://github.com/RHEnVision/provisioning-backend/commit/99303730cc564f2d8d702d02723833ccde6c3d33)), related to [HMSPROV-365](https://issues.redhat.com/browse/HMSPROV-365)

### Chore
- version 0.11.0 ([274ec81](https://github.com/RHEnVision/provisioning-backend/commit/274ec81ae0a48bccdcbfde8c6069f031aa0073f8))

### Code Refactoring
- Error payload with messages ([fa9ae9b](https://github.com/RHEnVision/provisioning-backend/commit/fa9ae9b553fff5e2543cf497108ff1008ec2e792))

### Features
- **HMSPROV-368:** add Version and BuildTime to ResponseError ([07e16ad](https://github.com/RHEnVision/provisioning-backend/commit/07e16adb1c4e17341cd4ff186e24bcd204531af0))
- better Kafka logging ([21d2f06](https://github.com/RHEnVision/provisioning-backend/commit/21d2f065976101ed20a0c69d5628f18cb35af959))
- increase default logging level to debug ([0e16d72](https://github.com/RHEnVision/provisioning-backend/commit/0e16d720f0102dfe88eb6455d610b9cb1bcdee31))
- statuser clowder deployment ([71c208a](https://github.com/RHEnVision/provisioning-backend/commit/71c208a3e7d1911e3049fc04371055a42e489841))

### Tests
- add unit test for Authentication failure ([99e4819](https://github.com/RHEnVision/provisioning-backend/commit/99e4819945dade484eeba7bb1443bebefb4b40e0)), related to [HMSPROV-347](https://issues.redhat.com/browse/HMSPROV-347)


<a name="0.11.0"></a>
## [0.11.0](https://github.com/RHEnVision/provisioning-backend/compare/0.10.0...0.11.0) (2022-11-21)

### Bug Fixes
- **config:** correct Unleash URL prefix ([bd6ab5a](https://github.com/RHEnVision/provisioning-backend/commit/bd6ab5a02d317a51e9a5a7ca742bbd372b2807bf))
- **config:** guard for non-exixtend kafka config ([a5b3d9c](https://github.com/RHEnVision/provisioning-backend/commit/a5b3d9c552953f3ddb7824c67712719b7a83bd27))
- **config:** unleash token as bearer header ([3bb424c](https://github.com/RHEnVision/provisioning-backend/commit/3bb424c889b8b84d8b982cae6e24ea9af1a927ba))
- **logging:** Disable middlewares for status routes ([811905d](https://github.com/RHEnVision/provisioning-backend/commit/811905dcf89335174a173f4362892c0f0931dce3)), related to [HMSPROV-333](https://issues.redhat.com/browse/HMSPROV-333)
- **reservation:** generic reservation by id ([5131c7b](https://github.com/RHEnVision/provisioning-backend/commit/5131c7b08c4164dcb11cb93ecb55916665132ccc)), related to [HMSPROV-349](https://issues.redhat.com/browse/HMSPROV-349)
- missing cache type variable for api ([0ed8c1b](https://github.com/RHEnVision/provisioning-backend/commit/0ed8c1bc8ade1d3b8346240855adcb76d5ab5a3f))
- null for aws_reservation_id when pending ([eb5e353](https://github.com/RHEnVision/provisioning-backend/commit/eb5e353d2d17541331bc460d587b95c48c15a75d))
- print full errors in logs ([7cc2e10](https://github.com/RHEnVision/provisioning-backend/commit/7cc2e10181bd623549b6ff78d03480e82c47bff3))

### Build
- **clowder:** Add image builder as clowder dependency ([01b1c06](https://github.com/RHEnVision/provisioning-backend/commit/01b1c06a9d59ff1334035699f46ae0915e4ac430)), related to [HMSPROV-194](https://issues.redhat.com/browse/HMSPROV-194)

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


