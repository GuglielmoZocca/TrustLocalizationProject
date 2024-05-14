### CODE

In this section it is explained the structure of the project, what contains every directory and file:
* Application_code: contains all code concern the Gateway application and how to manage the data come from devices.
  - Chipers: it contains the source code and the program for xor encrypting data from a device
  - DeviceFiles: it contains files that simulate are utilized to simulate some data calculated by a device. The device “11/a3” are data from the testbed, the others are example used for testing
  - Device_cript: it contains files that represent the encrypted and hash version of what can be find in DeviceFiles. This are the simulation of the data that can come from devices.
  - LogFiles: contains an example of log that come from an experiment in the testbed.
  - ParsingCode: it contains some python code for creating the appropriate device data from log “ParsingLogDeviceData.py” and for computing the average time to calculating the ranging “ParsingLogTimeRanging.py”, put ting this information in “AVGtemp[11/a3].txt” (in this case it is considered the device “11/a3”.
  - main.go: source code of the gateway application.
  - Position_test.go: test code for the correctness of some application functionality.
  - TimeAppGateway.txt: List of times that indicates the calculated period of a cycle of target position computation.
  - pkg: packages used in the gateway application.
* bin: command used in the creation of the blockchain
* calliper-workspace: it contains the necessary code to evaluate the performance of blockchain code with caliper.
  - benchmarks: it contains the configuration files for the different test on chaincode.
    + myDeviceBenchmark.yaml: configuration file for reading device test.
    + ReadTargetBenchmark.yaml: configuration file for reading target test.
    + totTest.yaml: configuration file for calculation position target test.
    + UpdateDTest.yaml: configuration file for updating device test.
  - networks: it contains the configuration file for the network, organizations and accounts to consider in caliper test.
  - node_modules: it contains node.js modules utilized in caliper test.
  - workload: it contains the code used for test
    + readDevice.js: test code for reading device.
    + ReadTargetBenchmark.yaml: test code for reading target.
     + tot_test.js: test code for calculation position target.
     + UpdateDevice.js: test code for updating device.
* Chaincode_dir: it contains the logic of chaincode of the blockchain.
  - main.go: code used for installation of the contract in the blockchain
  - collections_config.json: configuration file to describe the collection used in the blockchain.
  - chaincode: it contains the code of chaincode
    + Data_processing.go: It contains all functionalities of the contract that change the state of the blockchain.
    + Data_request.go: It contains all functionalities of the contract that ask data from the blockchain.
* Conf_code: it is the directory of Admin application code.
  - main.go: source code of the Admin application.
  - Conf_test.go: test code for the correctness of some application functionality.
* config: it contains configuration files for the blockchain network.
  - core.yaml and configtx.yaml: channel configuration files.
  - orderer.yaml: order configuration file.
* User_code: it is the directory of User application code.
  - main.go: source code of the User application.
* uwb-rng-radio-solution: it contains code to do some ranging experiment in testbed.
  - rng-init.c: source code for every device to do ranging operation. Must change “linkaddr_t init” for every device.
  - rng-resp.c: source code for the target.
* Report: document that describes the project.
* Presentation: slides that describe the project.
* Network: it contains all the codes for the structure of the blockchain.
  - bft-config: it contains the configuration of the blockchain channel in case it used the bot consensus algorithm.
  - compose: it contains the configurations file to build docker containers.
  - configtx: it contains the configuration of the blockchain channel in general case.
  - initialize_for_test.sh: script for building the Hyperledger fabric structure in case of bft consensus algorithm, use “cryptogen” command for the creation of  the cryptographic material  and deploy the contract in the blockchain.
  - monitordocker.sh: script to process the processing of the containers
  - network.sh: general script to execute various operations with Hyperledger Fabric deployed.
  - scripts: it contains some scripts to operate with the built blockchain.
  - setOrgEnv.sh: script for update some environment variable.
  - system-genesis-block: it will contains the genesis block of the blockchain.
  - organizations: Currently contains the configuration code that will be used in the initialization of the blockchain for the creation of the cryptographic material of the participants of the indicated organizations and the order nodes. The cryptographic material will put in this directory.