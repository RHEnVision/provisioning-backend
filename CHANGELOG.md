<!-- insertion marker -->
<a name="0.20.0"></a>

## [0.20.0](https://github.com/RHEnVision/provisioning-backend/compare/0.19.0...0.20.0) (2023-05-03)

### Bug Fixes

- **HMS-1671:** re-enable sonarcube ([25b0abe](https://github.com/RHEnVision/provisioning-backend/commit/25b0abe04b50e253e9f9ba3bf6cfe08f3f7f4a3e))
- **HMS-719:** Azure image check ([5d5e207](https://github.com/RHEnVision/provisioning-backend/commit/5d5e207778a114bd51ae87a73272e82f557c6256))
- **HMS-1567:** insights tags support ([7c71a2e](https://github.com/RHEnVision/provisioning-backend/commit/7c71a2e41cd66ee37b4853a1573842895298bd8d))

### Code Refactoring

- Change metrics according to app-sre notes ([cf3616e](https://github.com/RHEnVision/provisioning-backend/commit/cf3616e0278e5693e4b24b599983185dfffc8d9c))
- Update permission check for not implemented sources ([4217524](https://github.com/RHEnVision/provisioning-backend/commit/4217524d57861ca006ca29a1a9c98146a7405c00))

<a name="0.19.0"></a>

## [0.19.0](https://github.com/RHEnVision/provisioning-backend/compare/0.18.0...0.19.0) (2023-04-14)

### Features

- store Azure instance's IP address ([df8a489](https://github.com/RHEnVision/provisioning-backend/commit/df8a489f4023d78c0a5fb510da5ae4c27418456c))
- allow Azure group principal in Lighthouse ([18db8e8](https://github.com/RHEnVision/provisioning-backend/commit/18db8e87c7ea137decf1f121574815d006016dff))
- details of Azure source ([7c681c4](https://github.com/RHEnVision/provisioning-backend/commit/7c681c49d6d4d8b23f207d6d07708bbe128cd3f2))
- Create Azure VMs in paralel ([e4d146c](https://github.com/RHEnVision/provisioning-backend/commit/e4d146c4e7b1f3d39c13a4cd960c4f4cd5afb281))

### Bug Fixes

- **HMS-1110:** use default region for perm check ([dec69bc](https://github.com/RHEnVision/provisioning-backend/commit/dec69bc2ee3a61dfcef555b32a56038708b50778))
- **HMS-1105:** measure jobs in seconds ([cf16b81](https://github.com/RHEnVision/provisioning-backend/commit/cf16b8191d4736b5769aad43642672cd953f0161))
- **HMS-1105:** check for make dashboard target ([c05c618](https://github.com/RHEnVision/provisioning-backend/commit/c05c6181e8d7788a27686456c85545a5a8d88653))
- **HMS-1105:** improve reservations dashboard ([92428fa](https://github.com/RHEnVision/provisioning-backend/commit/92428faaaca22e8db03858c7e9dd360aca0bae30))
- **HMS-1105:** add makefile dashboard target ([983b11c](https://github.com/RHEnVision/provisioning-backend/commit/983b11ce9e426491d98138344bbb41c9e6231425))
- **HMS-1105:** add reservations and jobs to dash ([3e71bf1](https://github.com/RHEnVision/provisioning-backend/commit/3e71bf16387954ad3e51ba1653fbc00e1e8698f7))

### Code Refactoring

- move cache hit metrics into the proper package ([12f1bf3](https://github.com/RHEnVision/provisioning-backend/commit/12f1bf3a5957efbde63ae2f2dac358b1c98ee3f5))
- remove unused singular CreateVM for Azure ([f26c527](https://github.com/RHEnVision/provisioning-backend/commit/f26c527c85d5ff364390d15e83cc4198e071c983))
- Create struct for create instances requests ([5e9a849](https://github.com/RHEnVision/provisioning-backend/commit/5e9a849dde52f962148e1045ff9e072a264c457f))

<a name="0.18.0"></a>

## [0.18.0](https://github.com/RHEnVision/provisioning-backend/compare/0.17.0...0.18.0) (2023-03-22)

### Features

- **HMS-1001:** Add created instances ids to GCP reservation ([9e8d936](https://github.com/RHEnVision/provisioning-backend/commit/9e8d9367da43c341e13dacdf188f9150d6ec46a3))
- Azure support for cloud init ([781753b](https://github.com/RHEnVision/provisioning-backend/commit/781753b02030dfc889a61faa8e890ec075a35260))
- shorten Azure polling intervals ([d097a9f](https://github.com/RHEnVision/provisioning-backend/commit/d097a9f336bb6b236b71ca2dafbacad2b69933fa))
- **HMS-761:** add instance description to aws job ([d0291d2](https://github.com/RHEnVision/provisioning-backend/commit/d0291d288f8a75ab3cc2c81db8922252514a6fbd))
- Add sentry writer to zerolog ([3105091](https://github.com/RHEnVision/provisioning-backend/commit/31050918df0bad45d9834ccfb4f69296eebca5fa))
- **HMS-1110:** Add source permission validation check endpoint ([e84ec81](https://github.com/RHEnVision/provisioning-backend/commit/e84ec81c966b7d47bfd156bbf4981d92d82c220a))
- Azure lighthouse offering template ([2029433](https://github.com/RHEnVision/provisioning-backend/commit/20294330f9b82e703c799d7e87643e263210c816))

### Bug Fixes

- **HMS-1105:** add reservation counters ([2ac7e5b](https://github.com/RHEnVision/provisioning-backend/commit/2ac7e5bfdf8bea1714ce2bcb6c3bf015298b98e3))
- **HMS-1396:** integration job queue test ([11f1ec6](https://github.com/RHEnVision/provisioning-backend/commit/11f1ec6bd19d725bfcf9e685405796b48289e8fc))
- **HMS-1403:** timeout for job queue ([81582d2](https://github.com/RHEnVision/provisioning-backend/commit/81582d21702d79c294aa798fcbbf6f7fb9791942))
- cascade delete of pubkey to Azure details ([a8c5208](https://github.com/RHEnVision/provisioning-backend/commit/a8c52088adb86f6f0cad7ac846464885dba7d6b6))

### Code Refactoring

- change zone to match image builder ([b97730e](https://github.com/RHEnVision/provisioning-backend/commit/b97730e9b637dceebebc789942f2d0e0a1ee029f))
- Simplify OpenAPI generator ([4301088](https://github.com/RHEnVision/provisioning-backend/commit/43010881879b73b12f5f37ca651c9a0c7b18b666))
- Add test for valid region/zone/location ([9630bbb](https://github.com/RHEnVision/provisioning-backend/commit/9630bbb79f38a1fce639d9cd1a788e6cebf75e3a))

<a name="0.17.0"></a>

## [0.17.0](https://github.com/RHEnVision/provisioning-backend/compare/0.16.0...0.17.0) (2023-03-08)

### Features

- Azure reservation details endpoint ([db103ff](https://github.com/RHEnVision/provisioning-backend/commit/db103ff11bafd3152c2b45d0ef1d178954ee8286))
- consume Azure secret ([e27a875](https://github.com/RHEnVision/provisioning-backend/commit/e27a875e9a1837f1411ae922623a905b40702dfa))
- allow nullable fields iOpenAPI ([c2ae5a8](https://github.com/RHEnVision/provisioning-backend/commit/c2ae5a86ae7d0698d1c927c4e712d782da705794))
- **HMS-894:** Add GCP reservation test ([64388f6](https://github.com/RHEnVision/provisioning-backend/commit/64388f635c800dca4716808700fbf87f9b20ca07))
- Add Azure reservation to OpenAPI ([288e3fb](https://github.com/RHEnVision/provisioning-backend/commit/288e3fb09850482f1c29aae63ca2bc00d18bd5ff))

### Bug Fixes

- **HMS-879:** missing metric registration ([2055cf5](https://github.com/RHEnVision/provisioning-backend/commit/2055cf5d6dd0228e5fb6efa303de19e699e7f017))
- name public IP and nic to be VM specific ([cc74d2e](https://github.com/RHEnVision/provisioning-backend/commit/cc74d2eca2b4f5173d866a8b279d9a95e7e4a20b))
- allow dynamic naming for Azure disk ([f64c209](https://github.com/RHEnVision/provisioning-backend/commit/f64c209ef320249247ff9b9c1cabd1f0cc7dbe79))
- **HMS-879:** workers metrics ([7aef3fc](https://github.com/RHEnVision/provisioning-backend/commit/7aef3fc752c86cc456820eafe00ef4b7f99db2cf))

### Code Refactoring

- logging initialization ([1d3dc01](https://github.com/RHEnVision/provisioning-backend/commit/1d3dc013e4f97a3ec9913232a8aa0c15e9750323))
- extract binary name getter ([eacc010](https://github.com/RHEnVision/provisioning-backend/commit/eacc010664ae721f12d5f83b8df8bbd9d44247ca))
- fix worker metrics registrations ([ca7893f](https://github.com/RHEnVision/provisioning-backend/commit/ca7893f65fbb062f6c6ed613825e456b5b4cfbba))

<a name="0.16.0"></a>

## [0.16.0](https://github.com/RHEnVision/provisioning-backend/compare/0.15.0...0.16.0) (2023-02-22)

### Features

- respect Amount in Azure deployments ([77942a2](https://github.com/RHEnVision/provisioning-backend/commit/77942a25a11d727d24bda3f86268ead6b2aa630c))
- Fetch image name from image builder ([077dd25](https://github.com/RHEnVision/provisioning-backend/commit/077dd25ae59e298bc1366943d2107319c8a71953))
- **HMS-969:** List and filter sources by their hyperscaler ([74f26f5](https://github.com/RHEnVision/provisioning-backend/commit/74f26f5769829e4abd3a6e820dea6771077b2c6e))
- **HMS-1110:** ListAttachedPolicies feature ([befed63](https://github.com/RHEnVision/provisioning-backend/commit/befed63f39e036d8f7159fc2fff4b1ef40fa8711))
- add provider for Sources in OpenAPI spec ([877dae8](https://github.com/RHEnVision/provisioning-backend/commit/877dae89b414235f9c7ab6e81d06208ac384f20c))

### Bug Fixes

- **HMS-951:** region refresh docs ([572622d](https://github.com/RHEnVision/provisioning-backend/commit/572622d879de1347cda87711314844a45163757e))
- **HMS-1269:** initialize clients in workers ([1da390d](https://github.com/RHEnVision/provisioning-backend/commit/1da390da29a18ee1d7f1e683fef590ff1742f98f))
- **HMS-1259:** update avail check buckets to ms ([bad43f8](https://github.com/RHEnVision/provisioning-backend/commit/bad43f81bf343b9f3898ab69ded82b6198804bd2))
- **HMS-951:** region/location/zone validation ([36a8366](https://github.com/RHEnVision/provisioning-backend/commit/36a8366b04bdbe8d18a42f0e02a95a3cfba6e1a7))
- **HMS-951:** refresh preloaded data ([a1ef8be](https://github.com/RHEnVision/provisioning-backend/commit/a1ef8bebf5eb78a980afca4e5bddb87df25a1357))
- **HMS-1259:** Add adjustable datasource and SLOs panels ([c382c91](https://github.com/RHEnVision/provisioning-backend/commit/c382c91b9b105ec76be0e82494ae262b6feb809a))
- **HMS-951:** move preloaded types into separate package ([fd4d278](https://github.com/RHEnVision/provisioning-backend/commit/fd4d2783d8b9ca620b5d85853ed902dd58ee5499))
- **HMS-1260:** document pubkey and template behavior ([497c082](https://github.com/RHEnVision/provisioning-backend/commit/497c0823c256f8e21ffb47a7af192715de2e3d8c))
- **HMS-860:** modify and update job queue metrics ([525d1d3](https://github.com/RHEnVision/provisioning-backend/commit/525d1d350f8cfeb0f1b310aeaf01691cef62c5bb))
- **HMS-860:** fix typo in function name ([361da5a](https://github.com/RHEnVision/provisioning-backend/commit/361da5a485ba204c1212a7f3884b6f0ced66e11b))
- **HMS-1242:** atomically read statistics ([0b7a42b](https://github.com/RHEnVision/provisioning-backend/commit/0b7a42bc86394c7dd03aa1f53a498c8689b3c5fa))
- **HMS-1240:** add step titles back ([687e73a](https://github.com/RHEnVision/provisioning-backend/commit/687e73a28b6316605de1cde055b2f82481961582))
- allow setting proxy per client ([6e191e4](https://github.com/RHEnVision/provisioning-backend/commit/6e191e414ee46b7167be9739decd53122e1b3ff2))
- **HMS-1209:** launch templates for AWS ([01b4933](https://github.com/RHEnVision/provisioning-backend/commit/01b493303e12a27c3ea7865dbc0b346448101b6a))
- **HMS-1106:** rename ListInstanceTypes ([f224cb4](https://github.com/RHEnVision/provisioning-backend/commit/f224cb41ff56595d247042286ba9808dc629b1aa))
- **HMSPROV-429:** floorist exporter ([bf11efd](https://github.com/RHEnVision/provisioning-backend/commit/bf11efd62e9fe9cd6d4fbbd908abc95deba69b76))

### Code Refactoring

- Add logs to statuser and add invalid requests metric ([aacbf62](https://github.com/RHEnVision/provisioning-backend/commit/aacbf62f1cd731a6610ec52b5a1e932f63442e26))

<a name="0.15.0"></a>

## [0.15.0](https://github.com/RHEnVision/provisioning-backend/compare/0.14.0...0.15.0) (2023-02-06)

### Features

- **HMS-953:** Put account id into the context for worker. ([ae8fc5a](https://github.com/RHEnVision/provisioning-backend/commit/ae8fc5a421139ef9cf860f1ed6f60791964a33a6))
- **HMS-926:** User identity passed to jobs. ([5c4999b](https://github.com/RHEnVision/provisioning-backend/commit/5c4999bdb766bc0e7829ec81259ecf520f6ddd92))
- **HMS-1122:** Add total received availability check metric ([9ef3886](https://github.com/RHEnVision/provisioning-backend/commit/9ef3886fc9f23e3493b2ebd5b9fd1fb955d0ccd2))
- Azure deployment task ([02b55c3](https://github.com/RHEnVision/provisioning-backend/commit/02b55c39af439a06db029ec0fb7a02a1286406c4))
- Azure reservation service ([22520fb](https://github.com/RHEnVision/provisioning-backend/commit/22520fbb72044401e2461fed07c002b76f1122fb))
- minimal PoC Azure deployment ([00dc2d3](https://github.com/RHEnVision/provisioning-backend/commit/00dc2d33baf25336a75003cd85195e2d1cb08bc6))

### Bug Fixes

- **HMS-1181:** ignore pubkey resource deletion without SA ([e0648a9](https://github.com/RHEnVision/provisioning-backend/commit/e0648a90e90cb8357faa04b04b13fc1da6560f53))
- **HMSPROV-1107:** update permissions to match sources ([379f26c](https://github.com/RHEnVision/provisioning-backend/commit/379f26c117738b79d4326ab01fe88c764bb99c8a))

### Code Refactoring

- Add step in Azure job ([8ae2020](https://github.com/RHEnVision/provisioning-backend/commit/8ae2020406d117eb10b6e898dfd76dadb4f76a28))
- Regenerate HTTP clients ([5933f8e](https://github.com/RHEnVision/provisioning-backend/commit/5933f8ef9169a3227e79e62b0c36b014a2525237))

<a name="0.14.0"></a>

## [0.14.0](https://github.com/RHEnVision/provisioning-backend/compare/0.13.0...0.14.0) (2023-01-25)

### Features

- **HMSPROV-428:** Add provisioning dashboard ([a93a9f5](https://github.com/RHEnVision/provisioning-backend/commit/a93a9f5589c3400f62fa86d1d93c66859eaf1f4e))

### Bug Fixes

- **HMSPROV-390:** unscoped update pubkey fix ([bd30ea8](https://github.com/RHEnVision/provisioning-backend/commit/bd30ea85cd57f19e77390f209c375e72d833eb33))
- **HMSPROV-433:** change resource type to application ([72d26d1](https://github.com/RHEnVision/provisioning-backend/commit/72d26d16f5324a0ac5b66d4558acc6a63f86c67c))
- **HMSPROV-390:** RSA fingerprint and migration ([e115286](https://github.com/RHEnVision/provisioning-backend/commit/e115286767281b9401c384c311feeda8820ca588))
- **HMSPROV-425:** recover panics in workers ([2a560d5](https://github.com/RHEnVision/provisioning-backend/commit/2a560d524040312ba3984197b8902dce3f1b1007))
- **HMSPROV-425:** incorporate dejq into the app ([f8a0b6f](https://github.com/RHEnVision/provisioning-backend/commit/f8a0b6f5fb5e767d2d0d57f1f11f892c7f014946))
- image builder clowder config ([29aa59d](https://github.com/RHEnVision/provisioning-backend/commit/29aa59d95cdb1b2fdc99731dcc44d78085932303))
- filtering Provisioning auth for Source ([cebbdb3](https://github.com/RHEnVision/provisioning-backend/commit/cebbdb3df4d21826fc7fdcc66c7fb33939b53e11))
- **HMSPROV-387:** filter out noisy kafka logs ([e0a7b21](https://github.com/RHEnVision/provisioning-backend/commit/e0a7b216744ae4232a6474a795c9ae2967eb99a6))
- **HMSPROV-387:** use time-based offset for statuser ([bbc59a9](https://github.com/RHEnVision/provisioning-backend/commit/bbc59a9cfc9960afe52db731507e761d6e4e2746))
- unique index on pubkey_resource ([2b68b0a](https://github.com/RHEnVision/provisioning-backend/commit/2b68b0a4240b5ec2e3d77ce64a5ce92292f65097))

### Code Refactoring

- Add numeric status code ([0c4591e](https://github.com/RHEnVision/provisioning-backend/commit/0c4591ecb299ea9e42c3036508bf74145966c427))

<a name="0.13.0"></a>

## [0.13.0](https://github.com/RHEnVision/provisioning-backend/compare/0.12.0...0.13.0) (2023-01-12)

### Features

- **HMSPROV-177:** Add availability check request duration metric ([5944953](https://github.com/RHEnVision/provisioning-backend/commit/594495321ddaa55f6802bdaeeefc7cb1b776c888))
- **HMSPROV-177:** Add total availability check request metric ([13d508d](https://github.com/RHEnVision/provisioning-backend/commit/13d508db6870a3ed93c1d8123f622cead6f22013))
- add account identity endpoint ([3df28fb](https://github.com/RHEnVision/provisioning-backend/commit/3df28fbbfc9ffc7ecfa4fa39deb00245d826b303))
- introduce availability status endpoint ([272f577](https://github.com/RHEnVision/provisioning-backend/commit/272f577b69852dbb4bd8a13a3e37fa81bf2c2e87))
- **HMSPROV-345:** Add source check availability per each provider ([c20ef14](https://github.com/RHEnVision/provisioning-backend/commit/c20ef14924e8c5599eebed364092574f0e875f0c))

### Bug Fixes

- **HMSPROV-407:** fix cw config validation ([e661ac8](https://github.com/RHEnVision/provisioning-backend/commit/e661ac8f46fc314765ed383e0ddac853112c8961))
- **HMSPROV-407:** disable cw for migrations ([dbc1ea4](https://github.com/RHEnVision/provisioning-backend/commit/dbc1ea48c24d32da5c0590c11c0bb1a1aa763873))
- **HMSPROV-389:** drop memcache count ([bfc9a00](https://github.com/RHEnVision/provisioning-backend/commit/bfc9a00c1229fc3d64fdfe129e402edc4569e6de))
- **HMSPROV-407:** fix blank logic in cw initialization ([9de1b76](https://github.com/RHEnVision/provisioning-backend/commit/9de1b76abab17175fde3c7ae786376ba1e50d1e3))
- **HMSPROV-387:** set consumer group for statuser ([bd800fa](https://github.com/RHEnVision/provisioning-backend/commit/bd800fa8769b8abaed84cea944c42d1a07135803))
- **HMSPROV-407:** further improve logging of cw config ([681ccb0](https://github.com/RHEnVision/provisioning-backend/commit/681ccb0d76aa3781e600197510b3dc67ce419a09))
- **HMSPROV-407:** improve logging of cw config ([7598917](https://github.com/RHEnVision/provisioning-backend/commit/7598917111d67c696cccb1de51bcda29e03fe900))
- **HMSPROV-414:** start dequeue loop in api only for memory ([3b2eaae](https://github.com/RHEnVision/provisioning-backend/commit/3b2eaae6411a8b90342da7dc30d54924bedf2c3a))
- **HMSPROV-407:** enable cloudwatch in clowder ([cf5663b](https://github.com/RHEnVision/provisioning-backend/commit/cf5663bfe00677a84e3c4bff5ad05f6b520e5fae))
- **HMSPROV-340:** nice error on arch mismatch ([ca5d32e](https://github.com/RHEnVision/provisioning-backend/commit/ca5d32e61105a77ce86ba783a240afad4340e701))
- **HMSPROV-399:** add dejq job queue size metric ([6e49fcc](https://github.com/RHEnVision/provisioning-backend/commit/6e49fcc78c016c5e5b74e29430151773fed0bae6))
- **HMSPROV-392:** Check if image is an original or a shared one ([a03bde2](https://github.com/RHEnVision/provisioning-backend/commit/a03bde289659702e504335f5d837f50914a9a1f3))
- **HMSPROV-352:** improve error message ([0dd032f](https://github.com/RHEnVision/provisioning-backend/commit/0dd032fa183edba623f47a24f62aa11e67e59b63))
- **HMSPROV-352:** error out jobs early ([a85152a](https://github.com/RHEnVision/provisioning-backend/commit/a85152a747535acdd94d4b784fdd3f7600918d73))
- **HMSPROV-390:** calculate fingerprint for AWS ([df98843](https://github.com/RHEnVision/provisioning-backend/commit/df988435e24621496e170e5fe349b6af8b4096f6))
- **HMSPROV-345:** change to Source ([559bcfe](https://github.com/RHEnVision/provisioning-backend/commit/559bcfe2ad607e37875906db225229af64359821))
- Kafka headers are slice now ([5fdc3eb](https://github.com/RHEnVision/provisioning-backend/commit/5fdc3eb88ff827e8f28347829018fb6fe238bc7d))
- **HMSPROV-345:** remove default tag ([5659d0d](https://github.com/RHEnVision/provisioning-backend/commit/5659d0d520402b721d5c6a3b2a20f95dbc285183))
- **HMSPROV-345:** Add event_type header and resource type to kafka msg ([28a69da](https://github.com/RHEnVision/provisioning-backend/commit/28a69da29e676801de84ff15d978ddb9103ba4fd))
- **HMSPROV-170:** change topic to platform.sources.status ([b63ff7a](https://github.com/RHEnVision/provisioning-backend/commit/b63ff7aaffd16b546e2f7b3c8fc7183977deab09))
- Use correct topic for sources availability check ([3e7f820](https://github.com/RHEnVision/provisioning-backend/commit/3e7f8208394010c2ecada463f44284245f8470a6))
- utilize clowder topic mapping ([8c2ef77](https://github.com/RHEnVision/provisioning-backend/commit/8c2ef7752b5331e5fc456781428bd2398328cf3d))
- **HMSPROV-368:** change version to BuildCommit ([f3991bb](https://github.com/RHEnVision/provisioning-backend/commit/f3991bb90a76034a9aa29d141f0f0ed340253a1e))
- ensure pubkey is always present on AWS ([2c310dd](https://github.com/RHEnVision/provisioning-backend/commit/2c310dd338b929619fee437a6da44f731f68e81e))

<a name="0.12.0"></a>

## [0.12.0](https://github.com/RHEnVision/provisioning-backend/compare/0.11.0...0.12.0) (2022-12-01)

### Features

- **HMSPROV-368:** add Version and BuildTime to ResponseError ([07e16ad](https://github.com/RHEnVision/provisioning-backend/commit/07e16adb1c4e17341cd4ff186e24bcd204531af0))
- better Kafka logging ([21d2f06](https://github.com/RHEnVision/provisioning-backend/commit/21d2f065976101ed20a0c69d5628f18cb35af959))
- increase default logging level to debug ([0e16d72](https://github.com/RHEnVision/provisioning-backend/commit/0e16d720f0102dfe88eb6455d610b9cb1bcdee31))
- statuser clowder deployment ([71c208a](https://github.com/RHEnVision/provisioning-backend/commit/71c208a3e7d1911e3049fc04371055a42e489841))

### Bug Fixes

- payload name not nullable ([996251d](https://github.com/RHEnVision/provisioning-backend/commit/996251d0c17dbf1b120d4c473a74f61854f77a61))
- intermittent failures on CI for ASM queue test ([656ee05](https://github.com/RHEnVision/provisioning-backend/commit/656ee05b8d9634d671aff0067ea7b1dc8336a48d))
- log topic alongside trace send message ([4f9a62a](https://github.com/RHEnVision/provisioning-backend/commit/4f9a62ac3c75612c4495615f641321ce9c7567ab))
- enable Kafka in Clowder ([8ea9023](https://github.com/RHEnVision/provisioning-backend/commit/8ea90236707a4f0cd080b48bd0f6aec6a2368deb))
- kafka port is a pointer ([4ca3076](https://github.com/RHEnVision/provisioning-backend/commit/4ca30768ed981a938f6d8d1de7d2142597ef29a9))
- create topics in kafka startup script ([f7b2fab](https://github.com/RHEnVision/provisioning-backend/commit/f7b2fabce583884a94df6aa9974254c5ee20b42d))
- scope existing pubkey search by source id ([af89244](https://github.com/RHEnVision/provisioning-backend/commit/af892449955be50118150015a6cf483c7d2ae97b))
- **kafka:** CA config from string ([34f4c59](https://github.com/RHEnVision/provisioning-backend/commit/34f4c59af28b52087e52dc2f12d8273bc0966f6a))
- **sources:** handle source without application type correctly ([75bc847](https://github.com/RHEnVision/provisioning-backend/commit/75bc8471cff0ffbb30908ad0bfec125fba3216e6))
- break availability queue sender loop on context cancel ([0ec4201](https://github.com/RHEnVision/provisioning-backend/commit/0ec420179b2372e875d829910123efcf5203fdd9))
- share HTTP transport across platform clients ([dcf7c38](https://github.com/RHEnVision/provisioning-backend/commit/dcf7c3894bfaac7d391437c175af0f75b9f31ead))
- break consume look on context cancel ([4795c58](https://github.com/RHEnVision/provisioning-backend/commit/4795c58aac69e0e8dc262f53382428715ad06a9f))

### Code Refactoring

- Error payload with messages ([fa9ae9b](https://github.com/RHEnVision/provisioning-backend/commit/fa9ae9b553fff5e2543cf497108ff1008ec2e792))

<a name="0.11.0"></a>

## [0.11.0](https://github.com/RHEnVision/provisioning-backend/compare/0.10.0...0.11.0) (2022-11-21)

### Features

- **refactor:** Add required true post aws in apispec ([73f3956](https://github.com/RHEnVision/provisioning-backend/commit/73f3956b2a08914d3adb2628b4addb6836a85941))
- availability check kafka topic ([cea9ea0](https://github.com/RHEnVision/provisioning-backend/commit/cea9ea0e6aaac9a71f9d2171d38fb73ab1c27eb6))
- reservation detail returns instance ids ([30a1f8e](https://github.com/RHEnVision/provisioning-backend/commit/30a1f8ec766a2da552c27da70c39c2b7b8863111))
- **kafka:** setup, configuration, availability check ([971c64d](https://github.com/RHEnVision/provisioning-backend/commit/971c64d37d62787778e4ed476cbcc3255ba8f6bd))

### Bug Fixes

- missing cache type variable for api ([0ed8c1b](https://github.com/RHEnVision/provisioning-backend/commit/0ed8c1bc8ade1d3b8346240855adcb76d5ab5a3f))
- **reservation:** generic reservation by id ([5131c7b](https://github.com/RHEnVision/provisioning-backend/commit/5131c7b08c4164dcb11cb93ecb55916665132ccc))
- null for aws_reservation_id when pending ([eb5e353](https://github.com/RHEnVision/provisioning-backend/commit/eb5e353d2d17541331bc460d587b95c48c15a75d))
- print full errors in logs ([7cc2e10](https://github.com/RHEnVision/provisioning-backend/commit/7cc2e10181bd623549b6ff78d03480e82c47bff3))
- **config:** guard for non-exixtend kafka config ([a5b3d9c](https://github.com/RHEnVision/provisioning-backend/commit/a5b3d9c552953f3ddb7824c67712719b7a83bd27))
- **config:** unleash token as bearer header ([3bb424c](https://github.com/RHEnVision/provisioning-backend/commit/3bb424c889b8b84d8b982cae6e24ea9af1a927ba))
- **config:** correct Unleash URL prefix ([bd6ab5a](https://github.com/RHEnVision/provisioning-backend/commit/bd6ab5a02d317a51e9a5a7ca742bbd372b2807bf))
- **logging:** Disable middlewares for status routes ([811905d](https://github.com/RHEnVision/provisioning-backend/commit/811905dcf89335174a173f4362892c0f0931dce3))

<a name="0.10.0"></a>

## [0.10.0](https://github.com/RHEnVision/provisioning-backend/compare/0.9.0...0.10.0) (2022-11-02)

### Features

- **clients:** Preload machine types for GCP ([741b820](https://github.com/RHEnVision/provisioning-backend/commit/741b8207c6ee22d95dc33c84ef001ac413f5e68a))
- **clients:** Generated machine types and types per zone ([a654baa](https://github.com/RHEnVision/provisioning-backend/commit/a654baa0132724f6c1a5c34e8d3efdad9c826344))
- **flags:** feature flags endpoint ([727a72f](https://github.com/RHEnVision/provisioning-backend/commit/727a72f86e5f0e89a9166ebb5dc74b5545fdfb27))
- **gcp:** add gcp request payload ([0c718cc](https://github.com/RHEnVision/provisioning-backend/commit/0c718cc6321fd965d055288514808f9b66ce110e))
- **cache:** add Redis cache and queue ([1d60e0f](https://github.com/RHEnVision/provisioning-backend/commit/1d60e0f59a8e2f01fda1451546f36dd3044effa3))
- welcome HTML page ([a2fe12c](https://github.com/RHEnVision/provisioning-backend/commit/a2fe12cdde2d6c3525cb1e706dfc27e6ce2c4d36))

### Bug Fixes

- **queue:** recognize unknown worker config values ([a4114e8](https://github.com/RHEnVision/provisioning-backend/commit/a4114e8c304ebd9922d8da5e5c01b0962bf0c72c))
- **queue:** improve errors and logging for dejq init ([dff6b11](https://github.com/RHEnVision/provisioning-backend/commit/dff6b11c7a17cc5f044d71064b83b9dbb00542c3))
- remove DAOInit Error ([5800a20](https://github.com/RHEnVision/provisioning-backend/commit/5800a20aa4f6a768c6b533eea752fd010c88bb25))
- Add error messages ([e937d94](https://github.com/RHEnVision/provisioning-backend/commit/e937d948ad8faaddc5a4632cabbe93926360c412))
- **ec2:** typo in etag prefix ([ea86eb2](https://github.com/RHEnVision/provisioning-backend/commit/ea86eb23c48c78fbd894100207aa4bfb59726030))
- throw an error when pubkey is duplicated (HMSPROV-309) ([5a5078d](https://github.com/RHEnVision/provisioning-backend/commit/5a5078dd355b39d17bf2fc7a11ffc20f204e4ee8))
- **changelog:** correct upstream links ([d34b161](https://github.com/RHEnVision/provisioning-backend/commit/d34b16184094fd51b11d1ad084b8a1201b910291))
- **services:** refactor errors ([c1b05f4](https://github.com/RHEnVision/provisioning-backend/commit/c1b05f44a60580e3618446a19598c6a161401337))
- **logging:** no colors for clowder ([cff318a](https://github.com/RHEnVision/provisioning-backend/commit/cff318a0ece63abe253915239e72fc9d51faf9ff))

### Code Refactoring

- **azure:** rename client type ([c7ec156](https://github.com/RHEnVision/provisioning-backend/commit/c7ec156c79d26ae1d1baa51c1871589ff839a7cb))
- **clients:** remove Customer prefix from EC2 client ([2ec34db](https://github.com/RHEnVision/provisioning-backend/commit/2ec34dbc9867b72f90fb20a44f53a23f29350f0a))
- **azure:** split service and cust ifaces ([924216b](https://github.com/RHEnVision/provisioning-backend/commit/924216b3e3eb5b8c1680116346889771d23152a6))
- **clients:** change arn to authentication for generic purpose ([5b19411](https://github.com/RHEnVision/provisioning-backend/commit/5b19411b4fbef8030867fbfc4e49eecc0c65a56b))
- **changelog:** different tool for changelog ([febcb38](https://github.com/RHEnVision/provisioning-backend/commit/febcb38efffad06ab590d848b52cce79059100f3))

<a name="0.9.0"></a>

## [0.9.0](https://github.com/RHEnVision/provisioning-backend/compare/0.8.0...0.9.0) (2022-10-17)

### Features

- delete pubkey resource ([a148c89](https://github.com/RHEnVision/provisioning-backend/commit/a148c89124987b9c336d4340b93565098a768cbd))

### Bug Fixes

- **changelog:** messages without subtypes ([c798bb0](https://github.com/RHEnVision/provisioning-backend/commit/c798bb05455fa1721806ba07fe819e1aaadd6952))

### Code Refactoring

- **clients:** HTTP transport improved ([70f2c07](https://github.com/RHEnVision/provisioning-backend/commit/70f2c07b94afa001a3aceef5a33a8aea0f77d1d0))

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

