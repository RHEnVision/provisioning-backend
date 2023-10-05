<!-- insertion marker -->
<a name="Unreleased"></a>
## Unreleased ([compare](https://github.com/RHEnVision/provisioning-backend/compare/1.8.0...HEAD)) (2023-10-05)

<!-- insertion marker -->
<a name="1.8.0"></a>

## [1.8.0](https://github.com/RHEnVision/provisioning-backend/compare/1.7.0...1.8.0) (2023-10-04)

### Features

- Allow Azure resource group customization ([b369cf8](https://github.com/RHEnVision/provisioning-backend/commit/b369cf818fa9ea4ef160caec828cec17ae745504)), related to [HMS-1772](https://issues.redhat.com/browse/HMS-1772)
- **[HMS-2248](https://issues.redhat.com/browse/HMS-2248):** Add token support for list template ([e067a69](https://github.com/RHEnVision/provisioning-backend/commit/e067a69ac85ef8c8f3e247ccb186cd406fd02d09))

### Bug Fixes

- **[HMS-2703](https://issues.redhat.com/browse/HMS-2703):** allow to configure TLS without CA cert ([259c6f3](https://github.com/RHEnVision/provisioning-backend/commit/259c6f32122d1623d1e1591d1429fa5fb421c479))
- **[HMS-2674](https://issues.redhat.com/browse/HMS-2674):** Error when name pattern is invalid ([3a5cf3f](https://github.com/RHEnVision/provisioning-backend/commit/3a5cf3fa3f4094dbf50d6faa28ac7ce93121bf9b))
- **[HMS-2640](https://issues.redhat.com/browse/HMS-2640):** Support missing values when using launch template ([d7daafe](https://github.com/RHEnVision/provisioning-backend/commit/d7daafe873e902e64118515f0dae33b36ee77e73))
- **[HMS-2555](https://issues.redhat.com/browse/HMS-2555):** create reservation right before job ([168476c](https://github.com/RHEnVision/provisioning-backend/commit/168476c45aca794e59ea33eef851d8471b306445))
- **[HMS-2545](https://issues.redhat.com/browse/HMS-2545):** add rh labels to all providers ([643f5f7](https://github.com/RHEnVision/provisioning-backend/commit/643f5f7258b1155ad92e281d361f3d4c6da7f2f8))
- **[HMS-2519](https://issues.redhat.com/browse/HMS-2519):** error handling of jobs ([69d3e0f](https://github.com/RHEnVision/provisioning-backend/commit/69d3e0fdbc2df1163deaf74ae1a86921aa699e93))
- **[HMS-2258](https://issues.redhat.com/browse/HMS-2258):** Refactor error handling ([398d920](https://github.com/RHEnVision/provisioning-backend/commit/398d920f0bcb14e4940897d4d301eb12cc6cfbbd))

### Code Refactoring

- remove primitive string pointers ([c610aa7](https://github.com/RHEnVision/provisioning-backend/commit/c610aa7a0e36d7bfcc081279d32b92ddc16d0d63))
- remove fingerprint from error string ([88f4954](https://github.com/RHEnVision/provisioning-backend/commit/88f4954e047a090fce8624812ec442c7e5d93eb4))
- improve source messages ([325b594](https://github.com/RHEnVision/provisioning-backend/commit/325b59467e5dd0fbeff4de67e6fd7dd0b3961c87))
- drop unused errors ([ba3b4ee](https://github.com/RHEnVision/provisioning-backend/commit/ba3b4ee1c9716c98b8c0c2549495dcc4c295800f))
- use panic for non-runtime error ([7750339](https://github.com/RHEnVision/provisioning-backend/commit/77503393c17f3e18cbd94863d5fe67ca1bacdc6b))
- rename other errors ([41a53f0](https://github.com/RHEnVision/provisioning-backend/commit/41a53f0fd378c17dc3842128805d08e3e9acedd5))
- rename HTTP client errors ([f2bb8d3](https://github.com/RHEnVision/provisioning-backend/commit/f2bb8d3a331b66bc6d750db239374e5f6aa2ba4f))
- remove deprecated /account_identity ([512302e](https://github.com/RHEnVision/provisioning-backend/commit/512302e3da6da1183343578e557b6ca75d78627d))
- remove deprecated instance types endpoint ([2027eb9](https://github.com/RHEnVision/provisioning-backend/commit/2027eb91c6e0b2abe427024d2a7cd268cb2e002f))
- drop supported package ([53ff1e0](https://github.com/RHEnVision/provisioning-backend/commit/53ff1e02b67c1dad61ad5b11d7bb7b2fea6b9a11))
- Nest links under Metadata ([aae35e3](https://github.com/RHEnVision/provisioning-backend/commit/aae35e30f7b6f84f09a736919ca328b70ddf0662))

<a name="1.7.0"></a>

## [1.7.0](https://github.com/RHEnVision/provisioning-backend/compare/1.6.0...1.7.0) (2023-09-06)

### Features

- Allow request params to have description ([223649f](https://github.com/RHEnVision/provisioning-backend/commit/223649fb13154d85f8631a19ac6a70d0205595f6)), related to [HMS-2495](https://issues.redhat.com/browse/HMS-2495)
- **[HMS-2226](https://issues.redhat.com/browse/HMS-2226):** Add metadata info to List endpoints ([47e5cea](https://github.com/RHEnVision/provisioning-backend/commit/47e5cead6cf2ff25f1295334cd7beb932e2689b6))
- **[HMS-2225](https://issues.redhat.com/browse/HMS-2225):** Add limit and offset support ([f3ab041](https://github.com/RHEnVision/provisioning-backend/commit/f3ab0414d722300585bf0b7c22c90460550b15b8))

### Bug Fixes

- **[HMS-2462](https://issues.redhat.com/browse/HMS-2462):** Include unavailable details of src ([5e76bfa](https://github.com/RHEnVision/provisioning-backend/commit/5e76bfae6d342dd388b8095511dcdb55ca4a40a8))
- **[HMS-2350](https://issues.redhat.com/browse/HMS-2350):** Add status info for a source ([8428db8](https://github.com/RHEnVision/provisioning-backend/commit/8428db8bf7edcdd6adf55a4c488854cb1235babc))

<a name="1.6.0"></a>

## [1.6.0](https://github.com/RHEnVision/provisioning-backend/compare/1.5.0...1.6.0) (2023-08-23)

### Features

- **[HMS-2304](https://issues.redhat.com/browse/HMS-2304):** Nest list endpoints with 'data' ([8c45999](https://github.com/RHEnVision/provisioning-backend/commit/8c4599978459b257c066f18d531fd24e5d8a46ca))
- **[HMS-729](https://issues.redhat.com/browse/HMS-729):** Reservations cleanup ([77bd017](https://github.com/RHEnVision/provisioning-backend/commit/77bd017632576e1a2ac1590c95fbfda18e462ac0))

### Bug Fixes

- **[HMS-2033](https://issues.redhat.com/browse/HMS-2033):** cache for RBAC ([c45bc6f](https://github.com/RHEnVision/provisioning-backend/commit/c45bc6f4ce5c9ff82e2c77d6402b887e0ee25ff6))
- **[HMS-2327](https://issues.redhat.com/browse/HMS-2327):** fix empty resource perm check ([b32e40e](https://github.com/RHEnVision/provisioning-backend/commit/b32e40e5c2e209a5b789f4eeab88d6b74b07c7ca))
- **[HMS-2325](https://issues.redhat.com/browse/HMS-2325):** return 429 on rate error ([790d113](https://github.com/RHEnVision/provisioning-backend/commit/790d11327f91e1c7679ae739d659a926383efe8d))
- **[HMS-2328](https://issues.redhat.com/browse/HMS-2328):** fix context in failed launch notifications ([8a7286e](https://github.com/RHEnVision/provisioning-backend/commit/8a7286e607f732ff9ee04f6182874add84cc5c1c))

<a name="1.5.0"></a>

## [1.5.0](https://github.com/RHEnVision/provisioning-backend/compare/1.4.0...1.5.0) (2023-08-09)

### Bug Fixes

- **[HMS-1453](https://issues.redhat.com/browse/HMS-1453):** disable RBAC for sources ([dc8bb55](https://github.com/RHEnVision/provisioning-backend/commit/dc8bb555018696c422b46b4da7e30f485688596f))
- **[HMS-1453](https://issues.redhat.com/browse/HMS-1453):** list sources can be public ([c1f1cd9](https://github.com/RHEnVision/provisioning-backend/commit/c1f1cd991d1bf71540981dc10b12200ec3883a7c))
- **[HMS-1453](https://issues.redhat.com/browse/HMS-1453):** platform RBAC support enabled ([a812045](https://github.com/RHEnVision/provisioning-backend/commit/a8120450aadce624fe7401c6d6643904cbce9173))
- **[HMS-2295](https://issues.redhat.com/browse/HMS-2295):** add numeric rid to GCP ([0a66644](https://github.com/RHEnVision/provisioning-backend/commit/0a66644174d22772e29d5518b359f4212f65a570))
- **[HMS-2290](https://issues.redhat.com/browse/HMS-2290):** Make replicas configurable ([43be8f6](https://github.com/RHEnVision/provisioning-backend/commit/43be8f63f1705d3a6148961dc3b13f041a48133a))
- **[HMS-2274](https://issues.redhat.com/browse/HMS-2274):** Add setTags permission to GCP ([c6f22d8](https://github.com/RHEnVision/provisioning-backend/commit/c6f22d802f787a25d5daf3922173a957e44f9496))
- **[HMS-2233](https://issues.redhat.com/browse/HMS-2233):** reservation id tag for AWS ([3a74b17](https://github.com/RHEnVision/provisioning-backend/commit/3a74b17cf83bea73898a5000382d593ab00d8ebf))
- **[HMS-1453](https://issues.redhat.com/browse/HMS-1453):** platform RBAC support ([2a1951b](https://github.com/RHEnVision/provisioning-backend/commit/2a1951b7c6732e246939893549cd088bbe33a692))

### Code Refactoring

- Change launch template name to id in GCP ([2b0ccb9](https://github.com/RHEnVision/provisioning-backend/commit/2b0ccb936b106a7cde3c39a952990f4fac2009db))
- Add missing launch template to payload ([34e1c2f](https://github.com/RHEnVision/provisioning-backend/commit/34e1c2fcf01664a053aacaca80af5424c0ab9a7d))

<a name="1.4.0"></a>

## [1.4.0](https://github.com/RHEnVision/provisioning-backend/compare/1.3.0...1.4.0) (2023-07-26)

### Features

- **[HMS-1246](https://issues.redhat.com/browse/HMS-1246):** Use permission checks in statuser ([9a23611](https://github.com/RHEnVision/provisioning-backend/commit/9a2361155c6e23e35f0897b7b2a928c59e8bf765))

### Bug Fixes

- **[HMS-2178](https://issues.redhat.com/browse/HMS-2178):** followup context bug fix ([f9bcae5](https://github.com/RHEnVision/provisioning-backend/commit/f9bcae5c10a9a12db0c08a45dd85e8ca52a2ccc4))
- **[HMS-2178](https://issues.redhat.com/browse/HMS-2178):** followup context fix for statuser ([586a756](https://github.com/RHEnVision/provisioning-backend/commit/586a7567d2d6dbdc9a1fcb57ad521892b4ddc48b))
- **[HMS-2002](https://issues.redhat.com/browse/HMS-2002):** postpone db stats 10 minutes ([cfcaecb](https://github.com/RHEnVision/provisioning-backend/commit/cfcaecbfe1f4aecc974bdd4e9a75a3ca2712ce31))
- **[HMS-2181](https://issues.redhat.com/browse/HMS-2181):** statuser throttling configuration options ([dce8174](https://github.com/RHEnVision/provisioning-backend/commit/dce81749131c827410c69131c61e58395c72ea58))
- **[HMS-2178](https://issues.redhat.com/browse/HMS-2178):** fix context handling in Kafka consumer ([6c75031](https://github.com/RHEnVision/provisioning-backend/commit/6c75031130d22e5e3111ba2f6138f4b97d958bc0))
- **[HMS-2002](https://issues.redhat.com/browse/HMS-2002):** calculate db stats immediately ([662a2f4](https://github.com/RHEnVision/provisioning-backend/commit/662a2f4b254bbd326e29ccf964d136e818098ee6))
- **[HMS-2175](https://issues.redhat.com/browse/HMS-2175):** add AWS resource tag ([29bd8b6](https://github.com/RHEnVision/provisioning-backend/commit/29bd8b6a1d437ce7b8ef47ad3e699ffeaba491bf))
- **[HMS-1783](https://issues.redhat.com/browse/HMS-1783):** limit reservations per second ([0b0fab7](https://github.com/RHEnVision/provisioning-backend/commit/0b0fab755d3626148c3e6a09150ae954c1007864))

### Code Refactoring

- Add GCP job test ([a9e30b2](https://github.com/RHEnVision/provisioning-backend/commit/a9e30b22780185b62f97ba6a1cff1b2a83504815))

<a name="1.3.0"></a>

## [1.3.0](https://github.com/RHEnVision/provisioning-backend/compare/1.2.0...1.3.0) (2023-07-14)

### Features

- **[HMS-1955](https://issues.redhat.com/browse/HMS-1955):** Add fetch instance(s) description step for GCP ([3b3a3f4](https://github.com/RHEnVision/provisioning-backend/commit/3b3a3f45c29cd9b5c173d953939598a17aab69cc))
- **[HMSPROV-2088](https://issues.redhat.com/browse/HMSPROV-2088):** refacor notification's context message ([65a0960](https://github.com/RHEnVision/provisioning-backend/commit/65a096061168012f68ca4ef8efbe3e47cfc0f554))

### Bug Fixes

- **[HMS-2143](https://issues.redhat.com/browse/HMS-2143):** pending state stat ([6615db0](https://github.com/RHEnVision/provisioning-backend/commit/6615db0a2e5cf96cafff28859ac859b285923dff))
- **[HMS-2002](https://issues.redhat.com/browse/HMS-2002):** fix job queue duration registration ([175ed5e](https://github.com/RHEnVision/provisioning-backend/commit/175ed5e399aebd805fbe39a8edffb0b31767be70))
- **[HMS-2002](https://issues.redhat.com/browse/HMS-2002):** database stats are in secs not ms ([ee885ee](https://github.com/RHEnVision/provisioning-backend/commit/ee885ee81df38dc940aa190b93dbc0bf1ced0744))
- **[HMS-2002](https://issues.redhat.com/browse/HMS-2002):** cleanup stats logging ([f494205](https://github.com/RHEnVision/provisioning-backend/commit/f49420574340fe20578daa01f14d141ffa076f46))
- **[HMS-2002](https://issues.redhat.com/browse/HMS-2002):** fix memory request limits ([69024fb](https://github.com/RHEnVision/provisioning-backend/commit/69024fbfdd20ce132f9f7acc5113d6229e72d738))

<a name="1.2.0"></a>

## [1.2.0](https://github.com/RHEnVision/provisioning-backend/compare/1.1.0...1.2.0) (2023-06-28)

### Features

- **[HMS-1550](https://issues.redhat.com/browse/HMS-1550):** improve initialization code for kafka ([9f67a21](https://github.com/RHEnVision/provisioning-backend/commit/9f67a21b79312922b1219533b7a640964b5a41a8))
- **[HMS-1932](https://issues.redhat.com/browse/HMS-1932):** Add gcp source upload info ([9e41a3d](https://github.com/RHEnVision/provisioning-backend/commit/9e41a3dc3e4c4654b48b626a61b87667ea1b5c8f))
- **[HMS-1550](https://issues.redhat.com/browse/HMS-1550):** send notification after launch ([6e73aad](https://github.com/RHEnVision/provisioning-backend/commit/6e73aad859edbc11ee053781fc00d4c8f6020e14))
- **[HMS-1429](https://issues.redhat.com/browse/HMS-1429):** Add GCP to spec ([df68369](https://github.com/RHEnVision/provisioning-backend/commit/df683692215c477b003449f130ead51b96e84711))

### Bug Fixes

- **[HMS-2002](https://issues.redhat.com/browse/HMS-2002):** proper launch statistics ([7c2eb8d](https://github.com/RHEnVision/provisioning-backend/commit/7c2eb8d75277e8e60334e44fd85cc1cdba294302))

### Code Refactoring

- Add username to pubkey body ([ecc4591](https://github.com/RHEnVision/provisioning-backend/commit/ecc4591ae0b8b4cee4bd338e81677faff1d8d2b6))
- Add step titles info to gcp reservation ([084331f](https://github.com/RHEnVision/provisioning-backend/commit/084331fbde3675f826d71892290ffe54dd88a91b))
- Change Name to NamePattern in GCP ([313e7e5](https://github.com/RHEnVision/provisioning-backend/commit/313e7e53d8372d68472df7fb83516937dd682de9))
- Update permissions for provisioning role in GCP ([cc6a554](https://github.com/RHEnVision/provisioning-backend/commit/cc6a5546a6eeed241dcfb8ad6878e1322bebea8c))

<a name="1.1.0"></a>

## [1.1.0](https://github.com/RHEnVision/provisioning-backend/compare/1.0.0...1.1.0) (2023-06-13)

### Features

- **[HMS-1773](https://issues.redhat.com/browse/HMS-1773):** Add  IPv4 to instance ([9fc5788](https://github.com/RHEnVision/provisioning-backend/commit/9fc57889177871645f4bd3f02aa635735f210e3c))
- **[HMS-1828](https://issues.redhat.com/browse/HMS-1828):** Add listing templates launch from template ([b09349e](https://github.com/RHEnVision/provisioning-backend/commit/b09349e82b6a34f4aeec2bfcdda3fa82a868d557))

### Bug Fixes

- **[HMS-1616](https://issues.redhat.com/browse/HMS-1616):** cleanup reservation OpenAPI markdown ([bb54e11](https://github.com/RHEnVision/provisioning-backend/commit/bb54e115e16178de12e6da2895461b1d9b59f9e0))
- **[HMS-1877](https://issues.redhat.com/browse/HMS-1877):** filter authentications ([912e649](https://github.com/RHEnVision/provisioning-backend/commit/912e64921d49dbf48451cd8f33e19036a13a3470))
- **[HMS-1884](https://issues.redhat.com/browse/HMS-1884):** return 400 on invalid compose ID ([3bc1315](https://github.com/RHEnVision/provisioning-backend/commit/3bc1315a6016ba7cfcae3783d385cc9d9a149e12))

### Code Refactoring

- correct the function name to use label ([d98fe75](https://github.com/RHEnVision/provisioning-backend/commit/d98fe753305d43adaadedadbd3c6fb688b96de4d))
- add zone info when listing instance desc ([ee55eea](https://github.com/RHEnVision/provisioning-backend/commit/ee55eea25065f31143c6c4b75b5c1b83a033c53e))

<a name="1.0.0"></a>

## [1.0.0](https://github.com/RHEnVision/provisioning-backend/compare/0.21.0...1.0.0) (2023-06-01)

### Features

- **[HMS-1773](https://issues.redhat.com/browse/HMS-1773):** Add IPV4 to GCP instance desc ([684b9f3](https://github.com/RHEnVision/provisioning-backend/commit/684b9f3d92228a0eba23b7b1c1c8f1766960fd4b))
- **[HMS-1429](https://issues.redhat.com/browse/HMS-1429):** Add get GCP reservation by id ([32d5183](https://github.com/RHEnVision/provisioning-backend/commit/32d518332a4be1bf6770d1983858162bc362b67c))

### Bug Fixes

- **[HMS-1616](https://issues.redhat.com/browse/HMS-1616):** cleanup reservation OpenAPI docs ([43bd347](https://github.com/RHEnVision/provisioning-backend/commit/43bd347caab68d3c84d908e7ab93ae5b6c86e264))
- **[HMSPROV-451](https://issues.redhat.com/browse/HMSPROV-451):** Use sources filtering mechanism ([7bdc98d](https://github.com/RHEnVision/provisioning-backend/commit/7bdc98dcddf2c317e04669439168afd121e7b012))
- **[HMS-1830](https://issues.redhat.com/browse/HMS-1830):** improve concurrency of account creation ([fec282d](https://github.com/RHEnVision/provisioning-backend/commit/fec282df85f5f44751b058e6039f0a661124b20d))
- **[HMS-1800](https://issues.redhat.com/browse/HMS-1800):** fix a typo in upload info ([a867d02](https://github.com/RHEnVision/provisioning-backend/commit/a867d026b804465042052c7f7da8ffb219dba3b7))
- **[HMS-1785](https://issues.redhat.com/browse/HMS-1785):** check if public ip is nil during dereference ([07db863](https://github.com/RHEnVision/provisioning-backend/commit/07db863de6cad448a18f8b28c1290b3deb865ed3))
- **[HMS-1616](https://issues.redhat.com/browse/HMS-1616):** OpenAPI cleanup and examples ([1f198bd](https://github.com/RHEnVision/provisioning-backend/commit/1f198bdae3b4464c41b98c4f081feb3887127db3))
- **[HMS-1800](https://issues.redhat.com/browse/HMS-1800):** refactor and add more caching ([decfb33](https://github.com/RHEnVision/provisioning-backend/commit/decfb331a2e5642e904bed0dcd3ac41b319eb732))

### Code Refactoring

- split ctxval package ([b7cdac0](https://github.com/RHEnVision/provisioning-backend/commit/b7cdac06a0fa6a199abc338883af25805791bc27))
- replace ctxval.Logger ([bd92efe](https://github.com/RHEnVision/provisioning-backend/commit/bd92efee1d05b09b662dc64c22c77f14f3a1ad64))
- deprecate ctxval.Logger ([a9b2381](https://github.com/RHEnVision/provisioning-backend/commit/a9b23812df4e15083647cfca837836ee698ec71f))

<a name="0.21.0"></a>

## [0.21.0](https://github.com/RHEnVision/provisioning-backend/compare/summit23a...0.21.0) (2023-05-17)

### Features

- **[HMSPROV-449](https://issues.redhat.com/browse/HMSPROV-449):** Add correlation id to logger and context ([a7d6c6e](https://github.com/RHEnVision/provisioning-backend/commit/a7d6c6e2208ad9e4a393d0464de62424810d6b30))

### Bug Fixes

- **[HMS-1784](https://issues.redhat.com/browse/HMS-1784):** set GOMAXPROC for API workers ([dcceb5e](https://github.com/RHEnVision/provisioning-backend/commit/dcceb5ec6813d26e8fef19bfa96934921d0a0438))
- **[HMS-1782](https://issues.redhat.com/browse/HMS-1782):** cap job concurrency at 100 ([19655a6](https://github.com/RHEnVision/provisioning-backend/commit/19655a6b8faa763cf1fe16cc3f722633d8362233))
- **[HMS-1616](https://issues.redhat.com/browse/HMS-1616):** split payload and model for better OpenAPI ([539f62c](https://github.com/RHEnVision/provisioning-backend/commit/539f62c2257da91d1434fafc069e0788d6cbb5e6))

### Code Refactoring

- use Gob encoding for cache ([b31d961](https://github.com/RHEnVision/provisioning-backend/commit/b31d96170b977178f5dc14a6f416a63915feb06d))

<a name="summit23a"></a>

## [summit23a](https://github.com/RHEnVision/provisioning-backend/compare/0.20.0...summit23a) (2023-05-06)

### Features

- Fetch Azure image resource group from IB ([19d1704](https://github.com/RHEnVision/provisioning-backend/commit/19d1704828c8bc514f7f8c49e9fff09c0cc3de68)), related to [HMS-1691](https://issues.redhat.com/browse/HMS-1691)

### Bug Fixes

- Allow OPTIONS method for Azure template ([db88d24](https://github.com/RHEnVision/provisioning-backend/commit/db88d24fc36b68e7a5fbe6ab3e3c680ba809f424)), related to [HMS-1148](https://issues.redhat.com/browse/HMS-1148)
- Disable gzip on Azure Lighthouse template ([74ccf45](https://github.com/RHEnVision/provisioning-backend/commit/74ccf45e75e5e3183ff9244b726516084a38f4de)), related to [HMS-1148](https://issues.redhat.com/browse/HMS-1148)
- Allow access from Azure portal ([55b38dd](https://github.com/RHEnVision/provisioning-backend/commit/55b38dd495f7fb82ebf748004c57a8c695cca47f)), related to [HMS-1148](https://issues.redhat.com/browse/HMS-1148)

### Code Refactoring

- remove compression ([d479702](https://github.com/RHEnVision/provisioning-backend/commit/d479702dfc8f28810b2f14e065bfd12381b10ca1))

<a name="0.20.0"></a>

## [0.20.0](https://github.com/RHEnVision/provisioning-backend/compare/0.19.0...0.20.0) (2023-05-03)

### Bug Fixes

- **[HMS-1671](https://issues.redhat.com/browse/HMS-1671):** re-enable sonarcube ([25b0abe](https://github.com/RHEnVision/provisioning-backend/commit/25b0abe04b50e253e9f9ba3bf6cfe08f3f7f4a3e))
- **[HMS-719](https://issues.redhat.com/browse/HMS-719):** Azure image check ([5d5e207](https://github.com/RHEnVision/provisioning-backend/commit/5d5e207778a114bd51ae87a73272e82f557c6256))
- **[HMS-1567](https://issues.redhat.com/browse/HMS-1567):** insights tags support ([7c71a2e](https://github.com/RHEnVision/provisioning-backend/commit/7c71a2e41cd66ee37b4853a1573842895298bd8d))

### Code Refactoring

- Change metrics according to app-sre notes ([cf3616e](https://github.com/RHEnVision/provisioning-backend/commit/cf3616e0278e5693e4b24b599983185dfffc8d9c))
- Update permission check for not implemented sources ([4217524](https://github.com/RHEnVision/provisioning-backend/commit/4217524d57861ca006ca29a1a9c98146a7405c00))

<a name="0.19.0"></a>

## [0.19.0](https://github.com/RHEnVision/provisioning-backend/compare/0.18.0...0.19.0) (2023-04-14)

### Features

- store Azure instance's IP address ([df8a489](https://github.com/RHEnVision/provisioning-backend/commit/df8a489f4023d78c0a5fb510da5ae4c27418456c)), related to [HMS-1595](https://issues.redhat.com/browse/HMS-1595)
- allow Azure group principal in Lighthouse ([18db8e8](https://github.com/RHEnVision/provisioning-backend/commit/18db8e87c7ea137decf1f121574815d006016dff)), related to [HMS-1148](https://issues.redhat.com/browse/HMS-1148)
- details of Azure source ([7c681c4](https://github.com/RHEnVision/provisioning-backend/commit/7c681c49d6d4d8b23f207d6d07708bbe128cd3f2)), related to [HMS-1509](https://issues.redhat.com/browse/HMS-1509)
- Create Azure VMs in paralel ([e4d146c](https://github.com/RHEnVision/provisioning-backend/commit/e4d146c4e7b1f3d39c13a4cd960c4f4cd5afb281)), related to [HMS-1407](https://issues.redhat.com/browse/HMS-1407)

### Bug Fixes

- **[HMS-1110](https://issues.redhat.com/browse/HMS-1110):** use default region for perm check ([dec69bc](https://github.com/RHEnVision/provisioning-backend/commit/dec69bc2ee3a61dfcef555b32a56038708b50778))
- **[HMS-1105](https://issues.redhat.com/browse/HMS-1105):** measure jobs in seconds ([cf16b81](https://github.com/RHEnVision/provisioning-backend/commit/cf16b8191d4736b5769aad43642672cd953f0161))
- **[HMS-1105](https://issues.redhat.com/browse/HMS-1105):** check for make dashboard target ([c05c618](https://github.com/RHEnVision/provisioning-backend/commit/c05c6181e8d7788a27686456c85545a5a8d88653))
- **[HMS-1105](https://issues.redhat.com/browse/HMS-1105):** improve reservations dashboard ([92428fa](https://github.com/RHEnVision/provisioning-backend/commit/92428faaaca22e8db03858c7e9dd360aca0bae30))
- **[HMS-1105](https://issues.redhat.com/browse/HMS-1105):** add makefile dashboard target ([983b11c](https://github.com/RHEnVision/provisioning-backend/commit/983b11ce9e426491d98138344bbb41c9e6231425))
- **[HMS-1105](https://issues.redhat.com/browse/HMS-1105):** add reservations and jobs to dash ([3e71bf1](https://github.com/RHEnVision/provisioning-backend/commit/3e71bf16387954ad3e51ba1653fbc00e1e8698f7))

### Code Refactoring

- move cache hit metrics into the proper package ([12f1bf3](https://github.com/RHEnVision/provisioning-backend/commit/12f1bf3a5957efbde63ae2f2dac358b1c98ee3f5))
- remove unused singular CreateVM for Azure ([f26c527](https://github.com/RHEnVision/provisioning-backend/commit/f26c527c85d5ff364390d15e83cc4198e071c983))
- Create struct for create instances requests ([5e9a849](https://github.com/RHEnVision/provisioning-backend/commit/5e9a849dde52f962148e1045ff9e072a264c457f))

<a name="0.18.0"></a>

## [0.18.0](https://github.com/RHEnVision/provisioning-backend/compare/0.17.0...0.18.0) (2023-03-22)

### Features

- **[HMS-1001](https://issues.redhat.com/browse/HMS-1001):** Add created instances ids to GCP reservation ([9e8d936](https://github.com/RHEnVision/provisioning-backend/commit/9e8d9367da43c341e13dacdf188f9150d6ec46a3))
- Azure support for cloud init ([781753b](https://github.com/RHEnVision/provisioning-backend/commit/781753b02030dfc889a61faa8e890ec075a35260)), related to [HMS-1435](https://issues.redhat.com/browse/HMS-1435)
- shorten Azure polling intervals ([d097a9f](https://github.com/RHEnVision/provisioning-backend/commit/d097a9f336bb6b236b71ca2dafbacad2b69933fa)), related to [HMS-1404](https://issues.redhat.com/browse/HMS-1404)
- **[HMS-761](https://issues.redhat.com/browse/HMS-761):** add instance description to aws job ([d0291d2](https://github.com/RHEnVision/provisioning-backend/commit/d0291d288f8a75ab3cc2c81db8922252514a6fbd))
- Add sentry writer to zerolog ([3105091](https://github.com/RHEnVision/provisioning-backend/commit/31050918df0bad45d9834ccfb4f69296eebca5fa)), related to [HMS-851](https://issues.redhat.com/browse/HMS-851)
- **[HMS-1110](https://issues.redhat.com/browse/HMS-1110):** Add source permission validation check endpoint ([e84ec81](https://github.com/RHEnVision/provisioning-backend/commit/e84ec81c966b7d47bfd156bbf4981d92d82c220a))
- Azure lighthouse offering template ([2029433](https://github.com/RHEnVision/provisioning-backend/commit/20294330f9b82e703c799d7e87643e263210c816)), related to [HMS-1148](https://issues.redhat.com/browse/HMS-1148)

### Bug Fixes

- **[HMS-1105](https://issues.redhat.com/browse/HMS-1105):** add reservation counters ([2ac7e5b](https://github.com/RHEnVision/provisioning-backend/commit/2ac7e5bfdf8bea1714ce2bcb6c3bf015298b98e3))
- **[HMS-1396](https://issues.redhat.com/browse/HMS-1396):** integration job queue test ([11f1ec6](https://github.com/RHEnVision/provisioning-backend/commit/11f1ec6bd19d725bfcf9e685405796b48289e8fc))
- **[HMS-1403](https://issues.redhat.com/browse/HMS-1403):** timeout for job queue ([81582d2](https://github.com/RHEnVision/provisioning-backend/commit/81582d21702d79c294aa798fcbbf6f7fb9791942))
- cascade delete of pubkey to Azure details ([a8c5208](https://github.com/RHEnVision/provisioning-backend/commit/a8c52088adb86f6f0cad7ac846464885dba7d6b6)), related to [HMS-1402](https://issues.redhat.com/browse/HMS-1402)

### Code Refactoring

- change zone to match image builder ([b97730e](https://github.com/RHEnVision/provisioning-backend/commit/b97730e9b637dceebebc789942f2d0e0a1ee029f))
- Simplify OpenAPI generator ([4301088](https://github.com/RHEnVision/provisioning-backend/commit/43010881879b73b12f5f37ca651c9a0c7b18b666))
- Add test for valid region/zone/location ([9630bbb](https://github.com/RHEnVision/provisioning-backend/commit/9630bbb79f38a1fce639d9cd1a788e6cebf75e3a))

<a name="0.17.0"></a>

## [0.17.0](https://github.com/RHEnVision/provisioning-backend/compare/0.16.0...0.17.0) (2023-03-08)

### Features

- Azure reservation details endpoint ([db103ff](https://github.com/RHEnVision/provisioning-backend/commit/db103ff11bafd3152c2b45d0ef1d178954ee8286)), related to [HMS-1393](https://issues.redhat.com/browse/HMS-1393)
- consume Azure secret ([e27a875](https://github.com/RHEnVision/provisioning-backend/commit/e27a875e9a1837f1411ae922623a905b40702dfa)), related to [HMS-1374](https://issues.redhat.com/browse/HMS-1374)
- allow nullable fields iOpenAPI ([c2ae5a8](https://github.com/RHEnVision/provisioning-backend/commit/c2ae5a86ae7d0698d1c927c4e712d782da705794)), related to [HMS-1357](https://issues.redhat.com/browse/HMS-1357)
- **[HMS-894](https://issues.redhat.com/browse/HMS-894):** Add GCP reservation test ([64388f6](https://github.com/RHEnVision/provisioning-backend/commit/64388f635c800dca4716808700fbf87f9b20ca07))
- Add Azure reservation to OpenAPI ([288e3fb](https://github.com/RHEnVision/provisioning-backend/commit/288e3fb09850482f1c29aae63ca2bc00d18bd5ff)), related to [HMS-1182](https://issues.redhat.com/browse/HMS-1182)

### Bug Fixes

- **[HMS-879](https://issues.redhat.com/browse/HMS-879):** missing metric registration ([2055cf5](https://github.com/RHEnVision/provisioning-backend/commit/2055cf5d6dd0228e5fb6efa303de19e699e7f017))
- name public IP and nic to be VM specific ([cc74d2e](https://github.com/RHEnVision/provisioning-backend/commit/cc74d2eca2b4f5173d866a8b279d9a95e7e4a20b)), related to [HMS-1146](https://issues.redhat.com/browse/HMS-1146)
- allow dynamic naming for Azure disk ([f64c209](https://github.com/RHEnVision/provisioning-backend/commit/f64c209ef320249247ff9b9c1cabd1f0cc7dbe79)), related to [HMS-1146](https://issues.redhat.com/browse/HMS-1146)
- **[HMS-879](https://issues.redhat.com/browse/HMS-879):** workers metrics ([7aef3fc](https://github.com/RHEnVision/provisioning-backend/commit/7aef3fc752c86cc456820eafe00ef4b7f99db2cf))

### Code Refactoring

- logging initialization ([1d3dc01](https://github.com/RHEnVision/provisioning-backend/commit/1d3dc013e4f97a3ec9913232a8aa0c15e9750323))
- extract binary name getter ([eacc010](https://github.com/RHEnVision/provisioning-backend/commit/eacc010664ae721f12d5f83b8df8bbd9d44247ca))
- fix worker metrics registrations ([ca7893f](https://github.com/RHEnVision/provisioning-backend/commit/ca7893f65fbb062f6c6ed613825e456b5b4cfbba))

<a name="0.16.0"></a>

## [0.16.0](https://github.com/RHEnVision/provisioning-backend/compare/0.15.0...0.16.0) (2023-02-22)

### Features

- respect Amount in Azure deployments ([77942a2](https://github.com/RHEnVision/provisioning-backend/commit/77942a25a11d727d24bda3f86268ead6b2aa630c)), related to [HMS-1146](https://issues.redhat.com/browse/HMS-1146)
- Fetch image name from image builder ([077dd25](https://github.com/RHEnVision/provisioning-backend/commit/077dd25ae59e298bc1366943d2107319c8a71953)), related to [HMS-1219](https://issues.redhat.com/browse/HMS-1219)
- **[HMS-969](https://issues.redhat.com/browse/HMS-969):** List and filter sources by their hyperscaler ([74f26f5](https://github.com/RHEnVision/provisioning-backend/commit/74f26f5769829e4abd3a6e820dea6771077b2c6e))
- **[HMS-1110](https://issues.redhat.com/browse/HMS-1110):** ListAttachedPolicies feature ([befed63](https://github.com/RHEnVision/provisioning-backend/commit/befed63f39e036d8f7159fc2fff4b1ef40fa8711))
- add provider for Sources in OpenAPI spec ([877dae8](https://github.com/RHEnVision/provisioning-backend/commit/877dae89b414235f9c7ab6e81d06208ac384f20c)), related to [HMS-969](https://issues.redhat.com/browse/HMS-969)

### Bug Fixes

- **[HMS-951](https://issues.redhat.com/browse/HMS-951):** region refresh docs ([572622d](https://github.com/RHEnVision/provisioning-backend/commit/572622d879de1347cda87711314844a45163757e))
- **[HMS-1269](https://issues.redhat.com/browse/HMS-1269):** initialize clients in workers ([1da390d](https://github.com/RHEnVision/provisioning-backend/commit/1da390da29a18ee1d7f1e683fef590ff1742f98f))
- **[HMS-1259](https://issues.redhat.com/browse/HMS-1259):** update avail check buckets to ms ([bad43f8](https://github.com/RHEnVision/provisioning-backend/commit/bad43f81bf343b9f3898ab69ded82b6198804bd2))
- **[HMS-951](https://issues.redhat.com/browse/HMS-951):** region/location/zone validation ([36a8366](https://github.com/RHEnVision/provisioning-backend/commit/36a8366b04bdbe8d18a42f0e02a95a3cfba6e1a7))
- **[HMS-951](https://issues.redhat.com/browse/HMS-951):** refresh preloaded data ([a1ef8be](https://github.com/RHEnVision/provisioning-backend/commit/a1ef8bebf5eb78a980afca4e5bddb87df25a1357))
- **[HMS-1259](https://issues.redhat.com/browse/HMS-1259):** Add adjustable datasource and SLOs panels ([c382c91](https://github.com/RHEnVision/provisioning-backend/commit/c382c91b9b105ec76be0e82494ae262b6feb809a))
- **[HMS-951](https://issues.redhat.com/browse/HMS-951):** move preloaded types into separate package ([fd4d278](https://github.com/RHEnVision/provisioning-backend/commit/fd4d2783d8b9ca620b5d85853ed902dd58ee5499))
- **[HMS-1260](https://issues.redhat.com/browse/HMS-1260):** document pubkey and template behavior ([497c082](https://github.com/RHEnVision/provisioning-backend/commit/497c0823c256f8e21ffb47a7af192715de2e3d8c))
- **[HMS-860](https://issues.redhat.com/browse/HMS-860):** modify and update job queue metrics ([525d1d3](https://github.com/RHEnVision/provisioning-backend/commit/525d1d350f8cfeb0f1b310aeaf01691cef62c5bb))
- **[HMS-860](https://issues.redhat.com/browse/HMS-860):** fix typo in function name ([361da5a](https://github.com/RHEnVision/provisioning-backend/commit/361da5a485ba204c1212a7f3884b6f0ced66e11b))
- **[HMS-1242](https://issues.redhat.com/browse/HMS-1242):** atomically read statistics ([0b7a42b](https://github.com/RHEnVision/provisioning-backend/commit/0b7a42bc86394c7dd03aa1f53a498c8689b3c5fa))
- **[HMS-1240](https://issues.redhat.com/browse/HMS-1240):** add step titles back ([687e73a](https://github.com/RHEnVision/provisioning-backend/commit/687e73a28b6316605de1cde055b2f82481961582))
- allow setting proxy per client ([6e191e4](https://github.com/RHEnVision/provisioning-backend/commit/6e191e414ee46b7167be9739decd53122e1b3ff2)), related to [HMS-1227](https://issues.redhat.com/browse/HMS-1227)
- **[HMS-1209](https://issues.redhat.com/browse/HMS-1209):** launch templates for AWS ([01b4933](https://github.com/RHEnVision/provisioning-backend/commit/01b493303e12a27c3ea7865dbc0b346448101b6a))
- **[HMS-1106](https://issues.redhat.com/browse/HMS-1106):** rename ListInstanceTypes ([f224cb4](https://github.com/RHEnVision/provisioning-backend/commit/f224cb41ff56595d247042286ba9808dc629b1aa))
- **[HMSPROV-429](https://issues.redhat.com/browse/HMSPROV-429):** floorist exporter ([bf11efd](https://github.com/RHEnVision/provisioning-backend/commit/bf11efd62e9fe9cd6d4fbbd908abc95deba69b76))

### Code Refactoring

- Add logs to statuser and add invalid requests metric ([aacbf62](https://github.com/RHEnVision/provisioning-backend/commit/aacbf62f1cd731a6610ec52b5a1e932f63442e26))

<a name="0.15.0"></a>

## [0.15.0](https://github.com/RHEnVision/provisioning-backend/compare/0.14.0...0.15.0) (2023-02-06)

### Features

- **[HMS-953](https://issues.redhat.com/browse/HMS-953):** Put account id into the context for worker. ([ae8fc5a](https://github.com/RHEnVision/provisioning-backend/commit/ae8fc5a421139ef9cf860f1ed6f60791964a33a6))
- **[HMS-926](https://issues.redhat.com/browse/HMS-926):** User identity passed to jobs. ([5c4999b](https://github.com/RHEnVision/provisioning-backend/commit/5c4999bdb766bc0e7829ec81259ecf520f6ddd92))
- **[HMS-1122](https://issues.redhat.com/browse/HMS-1122):** Add total received availability check metric ([9ef3886](https://github.com/RHEnVision/provisioning-backend/commit/9ef3886fc9f23e3493b2ebd5b9fd1fb955d0ccd2))
- Azure deployment task ([02b55c3](https://github.com/RHEnVision/provisioning-backend/commit/02b55c39af439a06db029ec0fb7a02a1286406c4)), related to [HMS-1058](https://issues.redhat.com/browse/HMS-1058)
- Azure reservation service ([22520fb](https://github.com/RHEnVision/provisioning-backend/commit/22520fbb72044401e2461fed07c002b76f1122fb)), related to [HMS-1058](https://issues.redhat.com/browse/HMS-1058)
- minimal PoC Azure deployment ([00dc2d3](https://github.com/RHEnVision/provisioning-backend/commit/00dc2d33baf25336a75003cd85195e2d1cb08bc6)), related to [HMS-1058](https://issues.redhat.com/browse/HMS-1058)

### Bug Fixes

- **[HMS-1181](https://issues.redhat.com/browse/HMS-1181):** ignore pubkey resource deletion without SA ([e0648a9](https://github.com/RHEnVision/provisioning-backend/commit/e0648a90e90cb8357faa04b04b13fc1da6560f53))
- **[HMSPROV-1107](https://issues.redhat.com/browse/HMSPROV-1107):** update permissions to match sources ([379f26c](https://github.com/RHEnVision/provisioning-backend/commit/379f26c117738b79d4326ab01fe88c764bb99c8a))

### Code Refactoring

- Add step in Azure job ([8ae2020](https://github.com/RHEnVision/provisioning-backend/commit/8ae2020406d117eb10b6e898dfd76dadb4f76a28)), related to [HMS-1165](https://issues.redhat.com/browse/HMS-1165)
- Regenerate HTTP clients ([5933f8e](https://github.com/RHEnVision/provisioning-backend/commit/5933f8ef9169a3227e79e62b0c36b014a2525237))

<a name="0.14.0"></a>

## [0.14.0](https://github.com/RHEnVision/provisioning-backend/compare/0.13.0...0.14.0) (2023-01-25)

### Features

- **[HMSPROV-428](https://issues.redhat.com/browse/HMSPROV-428):** Add provisioning dashboard ([a93a9f5](https://github.com/RHEnVision/provisioning-backend/commit/a93a9f5589c3400f62fa86d1d93c66859eaf1f4e))

### Bug Fixes

- **[HMSPROV-390](https://issues.redhat.com/browse/HMSPROV-390):** unscoped update pubkey fix ([bd30ea8](https://github.com/RHEnVision/provisioning-backend/commit/bd30ea85cd57f19e77390f209c375e72d833eb33))
- **[HMSPROV-433](https://issues.redhat.com/browse/HMSPROV-433):** change resource type to application ([72d26d1](https://github.com/RHEnVision/provisioning-backend/commit/72d26d16f5324a0ac5b66d4558acc6a63f86c67c))
- **[HMSPROV-390](https://issues.redhat.com/browse/HMSPROV-390):** RSA fingerprint and migration ([e115286](https://github.com/RHEnVision/provisioning-backend/commit/e115286767281b9401c384c311feeda8820ca588))
- **[HMSPROV-425](https://issues.redhat.com/browse/HMSPROV-425):** recover panics in workers ([2a560d5](https://github.com/RHEnVision/provisioning-backend/commit/2a560d524040312ba3984197b8902dce3f1b1007))
- **[HMSPROV-425](https://issues.redhat.com/browse/HMSPROV-425):** incorporate dejq into the app ([f8a0b6f](https://github.com/RHEnVision/provisioning-backend/commit/f8a0b6f5fb5e767d2d0d57f1f11f892c7f014946))
- image builder clowder config ([29aa59d](https://github.com/RHEnVision/provisioning-backend/commit/29aa59d95cdb1b2fdc99731dcc44d78085932303)), related to [HMSPROV-421](https://issues.redhat.com/browse/HMSPROV-421)
- filtering Provisioning auth for Source ([cebbdb3](https://github.com/RHEnVision/provisioning-backend/commit/cebbdb3df4d21826fc7fdcc66c7fb33939b53e11)), related to [HMSPROV-426](https://issues.redhat.com/browse/HMSPROV-426)
- **[HMSPROV-387](https://issues.redhat.com/browse/HMSPROV-387):** filter out noisy kafka logs ([e0a7b21](https://github.com/RHEnVision/provisioning-backend/commit/e0a7b216744ae4232a6474a795c9ae2967eb99a6))
- **[HMSPROV-387](https://issues.redhat.com/browse/HMSPROV-387):** use time-based offset for statuser ([bbc59a9](https://github.com/RHEnVision/provisioning-backend/commit/bbc59a9cfc9960afe52db731507e761d6e4e2746))
- unique index on pubkey_resource ([2b68b0a](https://github.com/RHEnVision/provisioning-backend/commit/2b68b0a4240b5ec2e3d77ce64a5ce92292f65097)), related to [HMSPROV-415](https://issues.redhat.com/browse/HMSPROV-415)

### Code Refactoring

- Add numeric status code ([0c4591e](https://github.com/RHEnVision/provisioning-backend/commit/0c4591ecb299ea9e42c3036508bf74145966c427))

<a name="0.13.0"></a>

## [0.13.0](https://github.com/RHEnVision/provisioning-backend/compare/0.12.0...0.13.0) (2023-01-12)

### Features

- **[HMSPROV-177](https://issues.redhat.com/browse/HMSPROV-177):** Add availability check request duration metric ([5944953](https://github.com/RHEnVision/provisioning-backend/commit/594495321ddaa55f6802bdaeeefc7cb1b776c888))
- **[HMSPROV-177](https://issues.redhat.com/browse/HMSPROV-177):** Add total availability check request metric ([13d508d](https://github.com/RHEnVision/provisioning-backend/commit/13d508db6870a3ed93c1d8123f622cead6f22013))
- add account identity endpoint ([3df28fb](https://github.com/RHEnVision/provisioning-backend/commit/3df28fbbfc9ffc7ecfa4fa39deb00245d826b303)), related to [HMSPROV-357](https://issues.redhat.com/browse/HMSPROV-357)
- introduce availability status endpoint ([272f577](https://github.com/RHEnVision/provisioning-backend/commit/272f577b69852dbb4bd8a13a3e37fa81bf2c2e87)), related to [/HMSPROV-337](https://issues.redhat.com/browse/HMSPROV-337)
- **[HMSPROV-345](https://issues.redhat.com/browse/HMSPROV-345):** Add source check availability per each provider ([c20ef14](https://github.com/RHEnVision/provisioning-backend/commit/c20ef14924e8c5599eebed364092574f0e875f0c))

### Bug Fixes

- **[HMSPROV-407](https://issues.redhat.com/browse/HMSPROV-407):** fix cw config validation ([e661ac8](https://github.com/RHEnVision/provisioning-backend/commit/e661ac8f46fc314765ed383e0ddac853112c8961))
- **[HMSPROV-407](https://issues.redhat.com/browse/HMSPROV-407):** disable cw for migrations ([dbc1ea4](https://github.com/RHEnVision/provisioning-backend/commit/dbc1ea48c24d32da5c0590c11c0bb1a1aa763873))
- **[HMSPROV-389](https://issues.redhat.com/browse/HMSPROV-389):** drop memcache count ([bfc9a00](https://github.com/RHEnVision/provisioning-backend/commit/bfc9a00c1229fc3d64fdfe129e402edc4569e6de))
- **[HMSPROV-407](https://issues.redhat.com/browse/HMSPROV-407):** fix blank logic in cw initialization ([9de1b76](https://github.com/RHEnVision/provisioning-backend/commit/9de1b76abab17175fde3c7ae786376ba1e50d1e3))
- **[HMSPROV-387](https://issues.redhat.com/browse/HMSPROV-387):** set consumer group for statuser ([bd800fa](https://github.com/RHEnVision/provisioning-backend/commit/bd800fa8769b8abaed84cea944c42d1a07135803))
- **[HMSPROV-407](https://issues.redhat.com/browse/HMSPROV-407):** further improve logging of cw config ([681ccb0](https://github.com/RHEnVision/provisioning-backend/commit/681ccb0d76aa3781e600197510b3dc67ce419a09))
- **[HMSPROV-407](https://issues.redhat.com/browse/HMSPROV-407):** improve logging of cw config ([7598917](https://github.com/RHEnVision/provisioning-backend/commit/7598917111d67c696cccb1de51bcda29e03fe900))
- **[HMSPROV-414](https://issues.redhat.com/browse/HMSPROV-414):** start dequeue loop in api only for memory ([3b2eaae](https://github.com/RHEnVision/provisioning-backend/commit/3b2eaae6411a8b90342da7dc30d54924bedf2c3a))
- **[HMSPROV-407](https://issues.redhat.com/browse/HMSPROV-407):** enable cloudwatch in clowder ([cf5663b](https://github.com/RHEnVision/provisioning-backend/commit/cf5663bfe00677a84e3c4bff5ad05f6b520e5fae))
- **[HMSPROV-340](https://issues.redhat.com/browse/HMSPROV-340):** nice error on arch mismatch ([ca5d32e](https://github.com/RHEnVision/provisioning-backend/commit/ca5d32e61105a77ce86ba783a240afad4340e701))
- **[HMSPROV-399](https://issues.redhat.com/browse/HMSPROV-399):** add dejq job queue size metric ([6e49fcc](https://github.com/RHEnVision/provisioning-backend/commit/6e49fcc78c016c5e5b74e29430151773fed0bae6))
- **[HMSPROV-392](https://issues.redhat.com/browse/HMSPROV-392):** Check if image is an original or a shared one ([a03bde2](https://github.com/RHEnVision/provisioning-backend/commit/a03bde289659702e504335f5d837f50914a9a1f3))
- **[HMSPROV-352](https://issues.redhat.com/browse/HMSPROV-352):** improve error message ([0dd032f](https://github.com/RHEnVision/provisioning-backend/commit/0dd032fa183edba623f47a24f62aa11e67e59b63))
- **[HMSPROV-352](https://issues.redhat.com/browse/HMSPROV-352):** error out jobs early ([a85152a](https://github.com/RHEnVision/provisioning-backend/commit/a85152a747535acdd94d4b784fdd3f7600918d73))
- **[HMSPROV-390](https://issues.redhat.com/browse/HMSPROV-390):** calculate fingerprint for AWS ([df98843](https://github.com/RHEnVision/provisioning-backend/commit/df988435e24621496e170e5fe349b6af8b4096f6))
- **[HMSPROV-345](https://issues.redhat.com/browse/HMSPROV-345):** change to Source ([559bcfe](https://github.com/RHEnVision/provisioning-backend/commit/559bcfe2ad607e37875906db225229af64359821))
- Kafka headers are slice now ([5fdc3eb](https://github.com/RHEnVision/provisioning-backend/commit/5fdc3eb88ff827e8f28347829018fb6fe238bc7d)), related to [/HMSPROV-337](https://issues.redhat.com/browse/HMSPROV-337)
- **[HMSPROV-345](https://issues.redhat.com/browse/HMSPROV-345):** remove default tag ([5659d0d](https://github.com/RHEnVision/provisioning-backend/commit/5659d0d520402b721d5c6a3b2a20f95dbc285183))
- **[HMSPROV-345](https://issues.redhat.com/browse/HMSPROV-345):** Add event_type header and resource type to kafka msg ([28a69da](https://github.com/RHEnVision/provisioning-backend/commit/28a69da29e676801de84ff15d978ddb9103ba4fd))
- **[HMSPROV-170](https://issues.redhat.com/browse/HMSPROV-170):** change topic to platform.sources.status ([b63ff7a](https://github.com/RHEnVision/provisioning-backend/commit/b63ff7aaffd16b546e2f7b3c8fc7183977deab09))
- Use correct topic for sources availability check ([3e7f820](https://github.com/RHEnVision/provisioning-backend/commit/3e7f8208394010c2ecada463f44284245f8470a6)), related to [HMSPROV-170](https://issues.redhat.com/browse/HMSPROV-170)
- utilize clowder topic mapping ([8c2ef77](https://github.com/RHEnVision/provisioning-backend/commit/8c2ef7752b5331e5fc456781428bd2398328cf3d)), related to [HMSPROV-343](https://issues.redhat.com/browse/HMSPROV-343)
- **[HMSPROV-368](https://issues.redhat.com/browse/HMSPROV-368):** change version to BuildCommit ([f3991bb](https://github.com/RHEnVision/provisioning-backend/commit/f3991bb90a76034a9aa29d141f0f0ed340253a1e))
- ensure pubkey is always present on AWS ([2c310dd](https://github.com/RHEnVision/provisioning-backend/commit/2c310dd338b929619fee437a6da44f731f68e81e)), related to [HMSPROV-339](https://issues.redhat.com/browse/HMSPROV-339)

<a name="0.12.0"></a>

## [0.12.0](https://github.com/RHEnVision/provisioning-backend/compare/0.11.0...0.12.0) (2022-12-01)

### Features

- **[HMSPROV-368](https://issues.redhat.com/browse/HMSPROV-368):** add Version and BuildTime to ResponseError ([07e16ad](https://github.com/RHEnVision/provisioning-backend/commit/07e16adb1c4e17341cd4ff186e24bcd204531af0))
- better Kafka logging ([21d2f06](https://github.com/RHEnVision/provisioning-backend/commit/21d2f065976101ed20a0c69d5628f18cb35af959))
- increase default logging level to debug ([0e16d72](https://github.com/RHEnVision/provisioning-backend/commit/0e16d720f0102dfe88eb6455d610b9cb1bcdee31))
- statuser clowder deployment ([71c208a](https://github.com/RHEnVision/provisioning-backend/commit/71c208a3e7d1911e3049fc04371055a42e489841))

### Bug Fixes

- payload name not nullable ([996251d](https://github.com/RHEnVision/provisioning-backend/commit/996251d0c17dbf1b120d4c473a74f61854f77a61)), related to [HMSPROV-373](https://issues.redhat.com/browse/HMSPROV-373)
- intermittent failures on CI for ASM queue test ([656ee05](https://github.com/RHEnVision/provisioning-backend/commit/656ee05b8d9634d671aff0067ea7b1dc8336a48d))
- log topic alongside trace send message ([4f9a62a](https://github.com/RHEnVision/provisioning-backend/commit/4f9a62ac3c75612c4495615f641321ce9c7567ab))
- enable Kafka in Clowder ([8ea9023](https://github.com/RHEnVision/provisioning-backend/commit/8ea90236707a4f0cd080b48bd0f6aec6a2368deb))
- kafka port is a pointer ([4ca3076](https://github.com/RHEnVision/provisioning-backend/commit/4ca30768ed981a938f6d8d1de7d2142597ef29a9))
- create topics in kafka startup script ([f7b2fab](https://github.com/RHEnVision/provisioning-backend/commit/f7b2fabce583884a94df6aa9974254c5ee20b42d))
- scope existing pubkey search by source id ([af89244](https://github.com/RHEnVision/provisioning-backend/commit/af892449955be50118150015a6cf483c7d2ae97b)), related to [HMSPROV-366](https://issues.redhat.com/browse/HMSPROV-366)
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
- **reservation:** generic reservation by id ([5131c7b](https://github.com/RHEnVision/provisioning-backend/commit/5131c7b08c4164dcb11cb93ecb55916665132ccc)), related to [HMSPROV-349](https://issues.redhat.com/browse/HMSPROV-349)
- null for aws_reservation_id when pending ([eb5e353](https://github.com/RHEnVision/provisioning-backend/commit/eb5e353d2d17541331bc460d587b95c48c15a75d))
- print full errors in logs ([7cc2e10](https://github.com/RHEnVision/provisioning-backend/commit/7cc2e10181bd623549b6ff78d03480e82c47bff3))
- **config:** guard for non-exixtend kafka config ([a5b3d9c](https://github.com/RHEnVision/provisioning-backend/commit/a5b3d9c552953f3ddb7824c67712719b7a83bd27))
- **config:** unleash token as bearer header ([3bb424c](https://github.com/RHEnVision/provisioning-backend/commit/3bb424c889b8b84d8b982cae6e24ea9af1a927ba))
- **config:** correct Unleash URL prefix ([bd6ab5a](https://github.com/RHEnVision/provisioning-backend/commit/bd6ab5a02d317a51e9a5a7ca742bbd372b2807bf))
- **logging:** Disable middlewares for status routes ([811905d](https://github.com/RHEnVision/provisioning-backend/commit/811905dcf89335174a173f4362892c0f0931dce3)), related to [HMSPROV-333](https://issues.redhat.com/browse/HMSPROV-333)

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

