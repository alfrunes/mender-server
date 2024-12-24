---
## [4.0.0-rc.5](https://github.com/alfrunes/mender-server/compare/v4.0.0-rc.3...v4.0.0-rc.5) (2024-12-24)


### Features

* add `version` command to all Go binaries ([ff439c9](https://github.com/alfrunes/mender-server/commit/ff439c93552ae7e32d3a0cb932339902f45271ec))
* Add `version` command to all Go binaries ([5ad7d86](https://github.com/alfrunes/mender-server/commit/5ad7d86168127ac3802b088671a89c1e371cf955))
* added feedback dialog ([8c0a3ba](https://github.com/alfrunes/mender-server/commit/8c0a3baa2fa4e4cf935d818235e651bd4c5ed85c))
* added service provider tenant UI and routes ([50125e0](https://github.com/alfrunes/mender-server/commit/50125e034c5f7cf3625dfa91f23105fe3331bd9a))
* added tenant audit logs type support ([77bf187](https://github.com/alfrunes/mender-server/commit/77bf187186c7ff8dd8ecbcabf5d845b5c281d689))
* added tenant creation form ([92664bc](https://github.com/alfrunes/mender-server/commit/92664bc2235f8f2edbab0d3ab2b9ea7cc191204f))
* added tenant edit form test ([a509e61](https://github.com/alfrunes/mender-server/commit/a509e611b8d5ac93aa0fb7e3ce673aae52a849f3))
* **deployments:** add filter field to deployment object ([a1153b9](https://github.com/alfrunes/mender-server/commit/a1153b97e97e52525a4b5abe70511285fa6b0e11))
* **deployments:** add filter field to deployment object ([fec5b91](https://github.com/alfrunes/mender-server/commit/fec5b91d59d07b1a0d85ccf077cd56aa5b192278))
* **deployments:** new endpoint for getting release by name ([db7c9c8](https://github.com/alfrunes/mender-server/commit/db7c9c8a220c64dc765b63397b33b6fc525a7713))
* **deployments:** new endpoint for getting release by name ([3a18e88](https://github.com/alfrunes/mender-server/commit/3a18e880ec5cddedc19ed08949777caedda4350d))
* **gui:** added possibility to trigger deployment & inventory data updates when troubleshooting ([11a9b7a](https://github.com/alfrunes/mender-server/commit/11a9b7a57a179c3d9605779b41f6d10b6dbc72fb))
* **gui:** added redux thunks rejection logger middleware ([ec773c3](https://github.com/alfrunes/mender-server/commit/ec773c38a01d93c0181a3e3d9a6cfba55b727416))
* **gui:** added redux thunks rejection logger middleware ([4cf82ee](https://github.com/alfrunes/mender-server/commit/4cf82ee53c1ab214e0f1810d34563d44d9117267))
* **gui:** added support for Personal Access Token auditlog entries ([9a9a6c3](https://github.com/alfrunes/mender-server/commit/9a9a6c3829611c35622e3812db7bbedd9bc9f9e5))
* **gui:** added the possibility to create service provider administering roles ([92d7e50](https://github.com/alfrunes/mender-server/commit/92d7e50e311d8c88f9847a83cec7b797ef9cebcc))
* **gui:** aligned role removal dialog with other parts of the UI ([8661704](https://github.com/alfrunes/mender-server/commit/866170425bef1f01f3a4a25f0d4e19fe5da94a6e))
* **gui:** aligned webhook listing with updated design ([80e55d1](https://github.com/alfrunes/mender-server/commit/80e55d15e361c21988e192bf715a219bb300f487))
* **gui:** enabled webhook scope selection ([cec277d](https://github.com/alfrunes/mender-server/commit/cec277d83adf930de47ca5bb03935aa350ea1af5))
* **gui:** ensured components are spelled consistently ([7dfe136](https://github.com/alfrunes/mender-server/commit/7dfe136f29afa680d32570c16b7fc09efd566104))
* **gui:** extended webhook event details ([0bfda40](https://github.com/alfrunes/mender-server/commit/0bfda409122ed6837e13cf7f5418b093bf3ef97b))
* **gui:** made deployment targets rely on filter information in the deployment to more reliably display target devices etc. ([47c92d4](https://github.com/alfrunes/mender-server/commit/47c92d4db494cfc77116258fc2ed7fdca8691400))
* **gui:** made troubleshooting quick actions slightly easier to access ([689a89e](https://github.com/alfrunes/mender-server/commit/689a89ea68cdc61896730f15cc0bad84029f87be))
* **inventory:** add support for "$in" operator in the device search API ([0fa739a](https://github.com/alfrunes/mender-server/commit/0fa739a65706b120adb5ec743a13682f19f2fa10))
* **inventory:** add support for "$in" operator in the device search API ([fd4eaf0](https://github.com/alfrunes/mender-server/commit/fd4eaf0ecc8b72ff7fa9cfe7b6f214bc4678a97f))
* new endpoint for listing deployments ([c1c9764](https://github.com/alfrunes/mender-server/commit/c1c9764beb33ba64024bd3b43dd1834274a2490a))
* new endpoint for listing deployments ([afb1566](https://github.com/alfrunes/mender-server/commit/afb15665474440751e0463582e5d08d07b626da8))
* reinstated source license check ([c502084](https://github.com/alfrunes/mender-server/commit/c502084b0329949333c5eac960ab29feb2667b2e))
* review of the supported cryptographic keys ([1ec473f](https://github.com/alfrunes/mender-server/commit/1ec473ff7ef05b92cd76696927018244a67525be))
* review of the supported cryptographic keys ([0cf10e0](https://github.com/alfrunes/mender-server/commit/0cf10e00cb5a0f9103977ed2510c4baebf385172))
* tenant details edit component ([38ae179](https://github.com/alfrunes/mender-server/commit/38ae179b50e82087a94b662f9c6ddb0edf064595))
* tenant list added ([522cc08](https://github.com/alfrunes/mender-server/commit/522cc0855d3ee8ef98166a13db2119c73d6415ce))
* tenant list added ([719dfeb](https://github.com/alfrunes/mender-server/commit/719dfebe015f79cc84235e7c8143e0d1b21716a8))
* updated upgrades and add-on page ([0bdfd2d](https://github.com/alfrunes/mender-server/commit/0bdfd2d84e89d098c4003178b7a22896bcec7929))


### Bug Fixes

* added missing api docs check after migration off of mendertesting ([25e5dac](https://github.com/alfrunes/mender-server/commit/25e5dacdb48d5dc93af3e9609a46c51b9bfbf13f))
* added missing api spec definitions for deployments v2 api ([25d0f6c](https://github.com/alfrunes/mender-server/commit/25d0f6c640fdc8c013a7d0332bc82a901843069e))
* additional compose alignment after compose setup changes ([fb25b25](https://github.com/alfrunes/mender-server/commit/fb25b2508b560b4548938c29f2de062b2b408a41))
* additional compose alignment after compose setup changes ([8e4a65e](https://github.com/alfrunes/mender-server/commit/8e4a65e0249e8f08e0cde7320648466a0a2d86b1))
* adjusted form to new requirements ([f6c3578](https://github.com/alfrunes/mender-server/commit/f6c357894bb7a349e1cc1ffe1d432a72a60ba8b4))
* **deployments:** deprecate v1 endpoint for listing deployments ([879b589](https://github.com/alfrunes/mender-server/commit/879b58986f7e34906cff649c687d47de9152455c))
* **deviceconfig:** enable multiplatform build ([183c4b8](https://github.com/alfrunes/mender-server/commit/183c4b85efc3d59d6c64410ddfeee6543a19695e))
* **deviceconfig:** enable multiplatform build ([fbbe646](https://github.com/alfrunes/mender-server/commit/fbbe6466981015f47f250ad673f35f00004d1589))
* **e2e-tests:** fixed an issue that could prevent the browser from closing release details by removing the focus of the release tag input input instead of triggering a close of the release details ([fee1686](https://github.com/alfrunes/mender-server/commit/fee16863a7d66ecf8b37349721557656e6825215))
* fixed an issue that prevented enterprise demo tenant creation ([f60ce2f](https://github.com/alfrunes/mender-server/commit/f60ce2fd094e8c2f493383f8cbb9a485f96496f4))
* fixed an issue that prevented enterprise demo tenant creation ([21374a1](https://github.com/alfrunes/mender-server/commit/21374a15d1630dc77325e0df5461b5c723150281))
* fixed an issue that prevented onboarding tips from showing ([ee4ee16](https://github.com/alfrunes/mender-server/commit/ee4ee16ecb400261f9a0c8246e4da340162e73d9))
* fixed an issue that prevented onboarding tips from showing ([c2ecfcf](https://github.com/alfrunes/mender-server/commit/c2ecfcffd0a21f17ea8b3485a5d87efa21ab233a))
* fixed an issue that prevented the UI from showing deeply nested software installations ([13496f3](https://github.com/alfrunes/mender-server/commit/13496f3468fd08dcc9656ba07463eba682cfaff4))
* fixed api specs linter errors ([55c3460](https://github.com/alfrunes/mender-server/commit/55c3460d0cec47b401e1c1fcc4ce1405a02c66c1))
* **gui:** added filter value preprocessing for added ltne filter ([dc18922](https://github.com/alfrunes/mender-server/commit/dc18922642608c165aa95ab321722412cb516bf1))
* **gui:** added missing link to rbac docs in the cooresponding section ([1d8c4ff](https://github.com/alfrunes/mender-server/commit/1d8c4ff3f71f5918ea98ff277f96c31a85ebffe5))
* **gui:** added readable name for ltne device filter ([a741011](https://github.com/alfrunes/mender-server/commit/a74101176c22df69455a9d0634494912e219daab))
* **gui:** addressed timestamp filter issue in audit logs ([60181c8](https://github.com/alfrunes/mender-server/commit/60181c838cbc3b058c6acea822af125d002e9c80))
* **gui:** addressed timestamp filter issue in audit logs ([0934bd6](https://github.com/alfrunes/mender-server/commit/0934bd6f34d0e6dbe3397b466fb05d2293c6dec0))
* **gui:** aligned auditlog spelling with rest in rbac settings ([2a3b47f](https://github.com/alfrunes/mender-server/commit/2a3b47f17f1c258728a0c2d1e51e0076b542ad32))
* **gui:** aligned loading behaviour for troubleshooting quick actions ([72dd44f](https://github.com/alfrunes/mender-server/commit/72dd44f51a81fff5475a07e34dc2c05f3a6daaf0))
* **gui:** aligned quick actions in release details with actually possibile actions ([365f564](https://github.com/alfrunes/mender-server/commit/365f5646f2c32956fa8c0cee22c20d8c3757948d))
* **gui:** aligned SP indicator chip with the rest of the UI ([15a3caf](https://github.com/alfrunes/mender-server/commit/15a3caf6b7e1ff45e75612a09c776dc247c9b1cc))
* **gui:** aligned tenant removal message with action ([9522803](https://github.com/alfrunes/mender-server/commit/9522803f373f5f8e509cdb5499282d989b7bfa8e))
* **gui:** aligned tenant removal message with action ([78c37c5](https://github.com/alfrunes/mender-server/commit/78c37c5485ba373c87d78354bb0d44526913696c))
* **gui:** allowed for async tenant management interactions ([da6455c](https://github.com/alfrunes/mender-server/commit/da6455c905bb6230ecb0594b6b50befd6d527ecb))
* **gui:** delayed potential SPs from retrieving release tags & device info ([b02ae24](https://github.com/alfrunes/mender-server/commit/b02ae2433ae6bf4c8813165a5c0b4aee573ac6af))
* **gui:** do not rely on user visible selector as the rhf + mui combo breaks the connection from select label to actual select ([e541e65](https://github.com/alfrunes/mender-server/commit/e541e655cc012cafbfb3679cdae5deb768978e29))
* **gui:** enable device configuration for non enterprise users ([2e2e966](https://github.com/alfrunes/mender-server/commit/2e2e966e3b8d5431439fbd464b82028e3d97f75f))
* **gui:** enable device configuration for non enterprise users ([67170c5](https://github.com/alfrunes/mender-server/commit/67170c5edb27a1061abf2826234fabab45e4dedf))
* **gui:** ensured effective user permissions are also shown for release tags ([b422e6f](https://github.com/alfrunes/mender-server/commit/b422e6f75c6cf08320393eb410a03504fc2f0649))
* **gui:** ensured pagination total is a number ([4606287](https://github.com/alfrunes/mender-server/commit/4606287b11c1616c1ac3e6266409cc582fedda55))
* **gui:** ensured tenant changing users get their location in the ui reset when they end up in an unsupported area in the changed to tenant ([b6a3bc3](https://github.com/alfrunes/mender-server/commit/b6a3bc3aa62cecde6d07bc6c7171371f430abfc0))
* **gui:** ensured tenant user creation settings are passed in full ([ee4edd1](https://github.com/alfrunes/mender-server/commit/ee4edd1a1294a17cb793cae84195ad173af6f026))
* **gui:** fixed an error that prevented upgrading account to paid plan ([03f50c0](https://github.com/alfrunes/mender-server/commit/03f50c0267ed14634be335855eeea846581a25b8))
* **gui:** fixed an error that prevented upgrading account to paid plan ([967a07a](https://github.com/alfrunes/mender-server/commit/967a07aba1d12069582a0457f584a1b3be3fc3f4))
* **gui:** fixed an issue that caused number comparisons in device filters to not work ([84e2398](https://github.com/alfrunes/mender-server/commit/84e2398fece6b10fddcf6f60e3ff744af903c707))
* **gui:** fixed an issue that caused version information to be parsed wrong ([2e2dbc6](https://github.com/alfrunes/mender-server/commit/2e2dbc616bc58388557b6ff01a11c24e505de35e))
* **gui:** fixed an issue that caused version information to be parsed wrong ([02dd86d](https://github.com/alfrunes/mender-server/commit/02dd86d09b8459395f69f8816bc28a118f92e7a3))
* **gui:** fixed an issue that could lead to unexpected locations in the UI when accessing unauthorized sections while authorized ([7938291](https://github.com/alfrunes/mender-server/commit/7938291f8ac37c7ee3366c0cf2773e2c0053438f))
* **gui:** fixed an issue that could prevent 2fa reconfiguration in the same session ([89b0430](https://github.com/alfrunes/mender-server/commit/89b0430720784b0732f988596fbdf7f2b66539b0))
* **gui:** fixed an issue that kept deployment selection on deployment abort ([f97ac98](https://github.com/alfrunes/mender-server/commit/f97ac98280aaea22e56754c7d32ffb4c7175033a))
* **gui:** fixed an issue that prevented deployment sizes from being shown ([eee8799](https://github.com/alfrunes/mender-server/commit/eee8799a0e64a22478dc7f6e47fb55416d681e08))
* **gui:** fixed an issue that prevented deployment sizes from being shown ([d2bbb8d](https://github.com/alfrunes/mender-server/commit/d2bbb8df54aea9288af6d77944a516a075816928))
* **gui:** fixed an issue that prevented device issue count retrieval after rtk migration ([f289c4b](https://github.com/alfrunes/mender-server/commit/f289c4b49970df3f21b67c568cbe265ca62950b5))
* **gui:** fixed an issue that prevented setting up a new webhook ([dcbb22f](https://github.com/alfrunes/mender-server/commit/dcbb22fc3a39dfd31e94f49ffacc2692ae3541a6))
* **gui:** fixed an issue that prevented setting up a new webhook ([d1c26c5](https://github.com/alfrunes/mender-server/commit/d1c26c5e82449a16a665f8b9cc2d1b5cd09fc5ee))
* **gui:** fixed an issue that prevented showing errors on preauth duplicate devices ([c3cffda](https://github.com/alfrunes/mender-server/commit/c3cffda732ecb31c883ab66c83632ae204a033aa))
* **gui:** fixed an issue that prevented showing role details in user list or user details for SP users ([bcfc045](https://github.com/alfrunes/mender-server/commit/bcfc0450b9b0dd477c204a75dc3f47d3b658c64b))
* **gui:** fixed an issue that prevented webhook creation as webhooks no longer can be edited, so their id was never adjusted ([9006b8b](https://github.com/alfrunes/mender-server/commit/9006b8b32d8a05c84771a4171b4a3ec370278da6))
* **gui:** fixed an issue that prevented webhook creation as webhooks no longer can be edited, so their id was never adjusted ([dde48bf](https://github.com/alfrunes/mender-server/commit/dde48bfd47b81143f5efeab404ee25b3eee03ae3))
* **gui:** fixed an issue that would sometimes prevent users from switching between tenants ([ce777fd](https://github.com/alfrunes/mender-server/commit/ce777fdc9ae558a21a630384152152872c94b7a5))
* **gui:** fixed an issue with the auditlog that prevented updating the window url after initializtation ([24a9600](https://github.com/alfrunes/mender-server/commit/24a9600aa736662d3712f41196b9709863fe6384))
* **gui:** fixed component nesting errors when using drawertitle ([1bd1d5b](https://github.com/alfrunes/mender-server/commit/1bd1d5b20528eda28e901642808caecbae189af2))
* **gui:** fixed layout issues on new tenant user creation ([7158daf](https://github.com/alfrunes/mender-server/commit/7158daf86bf2e0b81a1dae3e4838b89277ad5308))
* **gui:** fixed loggedin detection on dashboard ([2fa1be0](https://github.com/alfrunes/mender-server/commit/2fa1be022df7680edf8dd5e0d78fbfa5f5c13c4f))
* **gui:** fixed several selector stability issues ([286dc56](https://github.com/alfrunes/mender-server/commit/286dc5683e1e2726687604ff09c860e496bfe916))
* **gui:** fixed some key related render errors ([5201ef6](https://github.com/alfrunes/mender-server/commit/5201ef692a5be739ed9fc92ce219915b8a44cd11))
* **gui:** friendlier input interactions on new tenant user creation ([3d7c94c](https://github.com/alfrunes/mender-server/commit/3d7c94c872dcfc232031f4066769fb38099fbb8a))
* **gui:** gave delta generation helptip on SP creation some meaning ([bad99fe](https://github.com/alfrunes/mender-server/commit/bad99fed0e10d0d74421c10b08c08957d98a363e))
* **gui:** included current device limit when configuring a tenant as a SP ([91db07b](https://github.com/alfrunes/mender-server/commit/91db07b906aafe862fcb86f2e2de30544fc2cc83))
* **gui:** limited overview information in auditlog to a single line ([938dffa](https://github.com/alfrunes/mender-server/commit/938dffac402a1d88a97d6338abb5be93ed545ba2))
* **gui:** limited overview information in auditlog to a single line ([5330fb6](https://github.com/alfrunes/mender-server/commit/5330fb65c5f4bd0971fd0ae53e8a22ead4da59ff))
* **gui:** limited user retrieval when adding/ removing users ([d26133f](https://github.com/alfrunes/mender-server/commit/d26133ffbfaff8d37d6175ea7ed0e0cfcb2e1d10))
* **gui:** moved license disclaimer generation into regular build to allow including nt-gui packages ([2e467c6](https://github.com/alfrunes/mender-server/commit/2e467c67040cea6d322c1d96409dd1157d5fd30f))
* **gui:** moved license disclaimer generation into regular build to allow including nt-gui packages ([53cde48](https://github.com/alfrunes/mender-server/commit/53cde481deff45d0468932d2fb282d85f0cd5a32))
* **gui:** moved tenant creation to rhf to ensure password generation returns if no existing email is entered ([8eab32f](https://github.com/alfrunes/mender-server/commit/8eab32f962950f1c1fb95ecf1e31960fbd910092))
* **gui:** pinned docker version to the same as central ci ([a81e7d9](https://github.com/alfrunes/mender-server/commit/a81e7d97ea126010d8d998345bdb0a674dbc52a8))
* **gui:** pinned docker version to the same as central ci ([c49d5d9](https://github.com/alfrunes/mender-server/commit/c49d5d9ac7b55312cd2bcb6e676d68ebd07615fa))
* **gui:** prevent deployment config retrieval on OS setups ([9e43a49](https://github.com/alfrunes/mender-server/commit/9e43a49802a3aee67678ded8acd3827c348e120a))
* **gui:** prevented device count retrieval for sp tenants to limit error messages in the UI ([cf0b44d](https://github.com/alfrunes/mender-server/commit/cf0b44dba6fbe700b5ad3b61c7fef14da61e5eec))
* **gui:** prevented disabled form inputs from showing validation errors ([2e7215a](https://github.com/alfrunes/mender-server/commit/2e7215aa93a3d357cfad34ec24e852ca66faf7df))
* **gui:** rely on user visible selector for role name input ([33a358d](https://github.com/alfrunes/mender-server/commit/33a358d9faa34019c8d2c6ba28091f6e890ce9f2))
* **gui:** removed effect that caused unwanted state modifications ([c767c8a](https://github.com/alfrunes/mender-server/commit/c767c8a491ab52ad13af7ba0541cda40f3f49b5e))
* **gui:** removed effect that caused unwanted state modifications ([0fd8717](https://github.com/alfrunes/mender-server/commit/0fd871701111e303707dd16cd5b99f94f4c88308))
* **gui:** removed initial device limit for newly created tenants ([84c97d2](https://github.com/alfrunes/mender-server/commit/84c97d2047bd68e6dbdf6f21a7df3a6dc11bb80a))
* **gui:** removed overly broad selector to center header content & adjusted non-aligned component ([d2ea651](https://github.com/alfrunes/mender-server/commit/d2ea651e31f9e29061881c4034f4c8f94a4b0788))
* **gui:** removed reliance on role description id and rely on user visible selector instead ([63a9797](https://github.com/alfrunes/mender-server/commit/63a9797975a2f8b1ee18a47d3965879421d1b336))
* **gui:** removed reliance on role description id and rely on user visible selector instead ([86bc66c](https://github.com/alfrunes/mender-server/commit/86bc66cd90f7f8d16eeca5c279fc09e6ee8b5046))
* **gui:** staggered app initialization to reduce unauthorized API calls in SP configs ([d1e1f3a](https://github.com/alfrunes/mender-server/commit/d1e1f3a26865009a45c4cec62c5059c9792f1536))
* prevented backend license publishing on frontend pipelines ([6f39bfe](https://github.com/alfrunes/mender-server/commit/6f39bfe2db36202248156f2b9145d31d96228319))
* removed api spec validation restriction to management apis ([6818392](https://github.com/alfrunes/mender-server/commit/6818392334c2ace4e5aca630b5b25fb150ef8353))
* removed api spec validation restriction to management apis ([350a0f8](https://github.com/alfrunes/mender-server/commit/350a0f86a47c3e2df7c041f125412204a49c8e18))
* removed state assignment that caused silent thunk rejection ([523cebe](https://github.com/alfrunes/mender-server/commit/523cebe855b347dcad8a652ed4a4717083b47f55))
* stop user from having similar email and password ([a0881a1](https://github.com/alfrunes/mender-server/commit/a0881a1fa07b6095a405f9a5ef450d1b824e1352))
* stop user from having similar email and password ([3fa4a43](https://github.com/alfrunes/mender-server/commit/3fa4a432780a40fb9b8c37633c7feca6ba3445c5))
* switched to single release endpoint when selecting release ([f107234](https://github.com/alfrunes/mender-server/commit/f107234e2721d5453ae86c126e0d73189c372fac))


### Reverts

* "fix(gui): fixed an issue that caused number comparisons in device filters to not work" ([787237e](https://github.com/alfrunes/mender-server/commit/787237ec8689d96c73beefbc74bcea7b96b274ba))
* Revert "docs(deviceauth): migration to OpenAPI3" ([93ab08a](https://github.com/alfrunes/mender-server/commit/93ab08ab6051aec3508bb550a4455d30ba2a9b56))


### Miscellaneous Chores

* release 4.0.0-rc.1 ([52f2403](https://github.com/alfrunes/mender-server/commit/52f24030fda75246c1ae77dbcc70b64dd180ce04))
* release 4.0.0-rc.2 ([03e6748](https://github.com/alfrunes/mender-server/commit/03e674867f72d104b9496c24edffddae9de43bf1))
* **release:** Bump released version to v4.0.0-rc.5 ([73abefc](https://github.com/alfrunes/mender-server/commit/73abefcfb00e8a3ede05aa7bd4805baa23721941))


### Continuous Integration

* force version 4.0.0 ([97392ef](https://github.com/alfrunes/mender-server/commit/97392efe0bf50f6ae964721f0d709ac59862913b))

## v4.0.0-rc.3 - 2024-12-17


### Bug Fixes


- *(gui)* Fixed an issue that would sometimes prevent users from switching between tenants
([MEN-7774](https://northerntech.atlassian.net/browse/MEN-7774)) ([ce777fd](https://github.com/mendersoftware/mender-server/commit/ce777fdc9ae558a21a630384152152872c94b7a5))  by @mzedel


  can't rely on the user list data as it doesn't contain all the user details

- *(gui)* Fixed an issue that prevented deployment sizes from being shown
 ([d2bbb8d](https://github.com/mendersoftware/mender-server/commit/d2bbb8df54aea9288af6d77944a516a075816928))  by @mzedel

- *(gui)* Fixed an issue that caused number comparisons in device filters to not work
([MEN-7717](https://northerntech.atlassian.net/browse/MEN-7717)) ([84e2398](https://github.com/mendersoftware/mender-server/commit/84e2398fece6b10fddcf6f60e3ff744af903c707))  by @mzedel

- *(gui)* Added readable name for ltne device filter
([MEN-7717](https://northerntech.atlassian.net/browse/MEN-7717)) ([a741011](https://github.com/mendersoftware/mender-server/commit/a74101176c22df69455a9d0634494912e219daab))  by @mzedel

- *(gui)* Fixed an issue that could lead to unexpected locations in the UI when accessing unauthorized sections while authorized
([MEN-7842](https://northerntech.atlassian.net/browse/MEN-7842)) ([7938291](https://github.com/mendersoftware/mender-server/commit/7938291f8ac37c7ee3366c0cf2773e2c0053438f))  by @mzedel

- *(gui)* Enable device configuration for non enterprise users
 ([67170c5](https://github.com/mendersoftware/mender-server/commit/67170c5edb27a1061abf2826234fabab45e4dedf))  by @thall


  Currently it's not possible to see device configuration if you host
  Mender self and have environment variable `HAVE_DEVICECONFIG=true`.
  
  Changes the predicate to be the same as for `hasDeviceConnect`.

- Fixed an issue that prevented the UI from showing deeply nested software installations
([MEN-7640](https://northerntech.atlassian.net/browse/MEN-7640)) ([13496f3](https://github.com/mendersoftware/mender-server/commit/13496f3468fd08dcc9656ba07463eba682cfaff4))  by @mzedel




### Documentation


- *(README)* Add step to clone repository
 ([f9d3bbd](https://github.com/mendersoftware/mender-server/commit/f9d3bbde382bca4592f41e3d6be7e8292dcb221f))  by @alfrunes

- *(README)* Consistently add syntax highlighting to code blocks
 ([8583102](https://github.com/mendersoftware/mender-server/commit/8583102cbbf49882b9a9ab1b80257516ec13dc24))  by @alfrunes

- Update README.md
 ([f7a1b09](https://github.com/mendersoftware/mender-server/commit/f7a1b097726672dd40ed7df17551229c5cf6ce7f))  by @alfrunes


  Adjusted styling (note color, added 1st level indentation,  taxonomy i.e., Mender Server, Mender Enterprise) to make it easy to follow and read.
- Document how to bring up the Virtual Device for enterprise setup
 ([c674566](https://github.com/mendersoftware/mender-server/commit/c674566e6d834c64d6e64d321c6e09b5f2a36259))  by @lluiscampos




### Features


- *(deployments)* New endpoint for getting release by name
([MEN-7575](https://northerntech.atlassian.net/browse/MEN-7575)) ([3a18e88](https://github.com/mendersoftware/mender-server/commit/3a18e880ec5cddedc19ed08949777caedda4350d))  by @kjaskiewiczz

- *(gui)* Added the possibility to create service provider administering roles
([MEN-7570](https://northerntech.atlassian.net/browse/MEN-7570)) ([92d7e50](https://github.com/mendersoftware/mender-server/commit/92d7e50e311d8c88f9847a83cec7b797ef9cebcc))  by @mzedel

- *(gui)* Aligned role removal dialog with other parts of the UI
 ([8661704](https://github.com/mendersoftware/mender-server/commit/866170425bef1f01f3a4a25f0d4e19fe5da94a6e))  by @mzedel

- *(gui)* Added support for Personal Access Token auditlog entries
([MEN-7622](https://northerntech.atlassian.net/browse/MEN-7622)) ([9a9a6c3](https://github.com/mendersoftware/mender-server/commit/9a9a6c3829611c35622e3812db7bbedd9bc9f9e5))  by @mzedel

- *(gui)* Added possibility to trigger deployment & inventory data updates when troubleshooting
([MEN-7657](https://northerntech.atlassian.net/browse/MEN-7657)) ([11a9b7a](https://github.com/mendersoftware/mender-server/commit/11a9b7a57a179c3d9605779b41f6d10b6dbc72fb))  by @mzedel

- *(gui)* Made deployment targets rely on filter information in the deployment to more reliably display target devices etc.
([MEN-7647](https://northerntech.atlassian.net/browse/MEN-7647)) ([47c92d4](https://github.com/mendersoftware/mender-server/commit/47c92d4db494cfc77116258fc2ed7fdca8691400))  by @mzedel





### Revert


- "fix(gui): fixed an issue that caused number comparisons in device filters to not work"
 ([787237e](https://github.com/mendersoftware/mender-server/commit/787237ec8689d96c73beefbc74bcea7b96b274ba))  by @mzedel


  This reverts commit 84e2398fece6b10fddcf6f60e3ff744af903c707.
  Signed-off-by: Manuel Zedel <manuel.zedel@northern.tech>
- Revert "docs(deviceauth): migration to OpenAPI3"
 ([93ab08a](https://github.com/mendersoftware/mender-server/commit/93ab08ab6051aec3508bb550a4455d30ba2a9b56))  by @kjaskiewiczz


  This reverts commit f7a33e9a71339522ee33f3808e7d6a8598144d2a.





---
