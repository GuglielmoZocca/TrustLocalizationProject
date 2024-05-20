## GOAL

The target of my research is to give to the structure presented in the paper [A Trust Architecture for Blockchain in IoT](https://arxiv.org/pdf/1906.11461) some feature of scalability and privatization of data in the blockchain, with also the aim to show these new characteristics through an example of possible application building a prototype of the end to end process from device to the user in the case localization case. The report of my research can be read in `report.pdf` file.

### CODE

In this section it is explained the structure of the project, what contains every directory and file:
* Application_code: contains all code concern the Gateway application and how to manage the data come from devices.
  - Chipers: it contains the source code and the program for xor encrypting data from a device
  - DeviceFiles: it contains files that simulate are utilized to simulate some data calculated by a device. The device “11:a3” are data from the testbed, the others are example used for testing
  - Device_cript: it contains files that represent the encrypted and hash version of what can be find in DeviceFiles. This are the simulation of the data that can come from devices.
  - LogFiles: contains an example of log that come from an experiment in the testbed.
  - ParsingCode: it contains some python code for creating the appropriate device data from log “ParsingLogDeviceData.py” and for computing the average time to calculating the ranging “ParsingLogTimeRanging.py”, put ting this information in “AVGtemp[11:a3].txt” (in this case it is considered the device “11:a3”.
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
  - system-genesis-block: it will contain the genesis block of the blockchain.
  - organizations: Currently contains the configuration code that will be used in the initialization of the blockchain for the creation of the cryptographic material of the participants of the indicated organizations and the order nodes. The cryptographic material will put in this directory.

### TEST

To test the project you must follow following instruction:
1. Clone the project.
2. Install all [prerequisite](https://hyperledger-fabric.readthedocs.io/en/latest/prereqs.html).
3. Download the necessary images in docker:
   1. `curl -sSLO https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/install-fabric.sh && chmod +x install-fabric.sh`
   2. `./install-fabric.sh d`
4. Go to the network directory: cd Path/to/Network  .
5. Build a Hyperledger Fabric blockchain and deploy the smart contract PositionContract: `./initialize_for_test.sh`
6. For test Gateway and User application.
   1. `cd ../Application_code`
   2. `go test main.go Position_test.go -v`
7. For test Admin application:
   1. `cd ../Conf_code`
   2. `go test main.go Conf_test.go -v`
8. For test the performance of chaincode:
   1. `cd ../caliper-workspace`
   2. Install the prequisite: Node-version v12.22.10 and NPM version 6.14.16
   3. Install caliper:
      1. `npm init`
      2. `npm install --only=prod @hyperledger/caliper-cli@0.5.0`
      3. `npx caliper bind --caliper-bind-sut fabric:2.2`
   4. For test the reading of a device: `npx caliper launch manager --caliper-workspace ./ --caliper-networkconfig networks/networkConfig.yaml --caliper-benchconfig benchmarks/myDeviceBenchmark.yaml —caliper-flow-only-test`
   5. For test the reading of a target: `npx caliper launch manager --caliper-workspace ./ --caliper-networkconfig networks/networkConfig.yaml --caliper-benchconfig benchmarks/ReadTargetBenchmark.yaml—caliper-flow-only-test`
   6. For test a cycle of position target calculation: `npx caliper launch manager --caliper-workspace ./ --caliper-networkconfig networks/networkConfig.yaml --caliper-benchconfig benchmarks/totTest.yaml —caliper-flow-only-test`
   7. For test an update of a device by the admin: `npx caliper launch manager --caliper-workspace ./ --caliper-networkconfig networks/networkConfig.yaml --caliper-benchconfig benchmarks/UpdateDTest.yaml —caliper-flow-only-test`
   8. To see the result of one of this test open "report.html".
9. For do some experiments with testbed, that is a device ranging with a target:
   1. Register to [CLOVES](https://iottestbed.disi.unitn.it/cloves/getting-started/) testbed and install Cooja.
   2. `cd ../uwb-rng-radio-solution`
   3. Go to website of CLOVES and select [map](https://iottestbed.disi.unitn.it/cloves/infrastructure/) of devices you want to utilize, and download the file with devices information of that map.
   4. Substitute in rng-init.c “linkaddr_t init” with wanted device address (as you can see as example in the file) and “linkaddr_t resp_list[NUM_DEST]” with wanted target address (as you can see as example in the file).
   5. Compile the code for the ranging: `make TARGET=evb1000`
   6. Enter to the website of CLOVES.
   7. Go to "Create Job"
   8. In Timeslot info select an island (the map of devices), as start time indicate ASAP and as duration indicate how much time do you want during the experiment.
   9. In binary file 1 insert as hardware evb1000, as Bin file “rng-init.bin”, as targets the id of node device in respective island.
   10. In binary file 2 insert as hardware evb1000, as Bin file “rng-resp.bin”, as targets the id of node target in respective island.
   11. Start the experiment
   12. After the time indicated in Timeslot info, go to the "Download jobs". Download the experiment result and extract the log.
   13. After that you can do some operation on the job file "job.log" obtained: `cd ../Application_code`
       1. You can use the parsing code in “ParsingCode/ParsingLogDeviceData.py” passing the log file to extract the wanted data in device[idDevice].txt where the value of idDevice depend on device id: `ParsingCode/ParsingLogDeviceData.py /PATH/TO/job.log`
       2. You can use the parsing code in “ParsingCode/ParsingLogTimeRanging.py” passing the log file to calculate the average time for a ranging operation in AVGtemp[idDevice].txt where the value of idDevice depend on device id: `ParsingCode/ParsingLogTimeRanging.py /PATH/TO/job.log`
       3. You can take the data received from the device "/DeviceFiles/device[idDevice].txt" a and encrypt it with the key (a character) in "/Devices_cript/device[idDevice]_ript.txt": `chipers/xor_c "/DeviceFiles/device[idDevice].txt" "/Devices_cript/device[idDevice]_ript.txt" "in" "key"`
       4. You can use the encrypted for possible test of Gateway and Admin application
10. To prove Gateway application:
    1. `cd ../Application_code`
    2. Start the application until doesn’t print anything more: `go main.go`
    3. It can also be seen in TimeAppGateway.txt a list of periods of time calculated for every cycle of position target calculation
11. To prove Admin application:
    1. `cd ../Conf_code`
    2. Start the application and prove some operation suggested: `go main.go`
12. To prove User application:
    1. `cd ../User_code`
    2. Start the application and prove some operation suggested: `go main.go`

