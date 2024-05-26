## GOAL

The target of my research is to give to the structure presented in the paper [A Trust Architecture for Blockchain in IoT](https://arxiv.org/pdf/1906.11461) some feature of privatization of data in the blockchain, with also the aim to show these new characteristics through an example of possible application building a prototype of the end to end process from device to the user in the case localization case. The report of my research can be read in `Report.pdf` file.

### CODE

In this section it is explained the directory structure of the project, what contains every directory and file:
* Application_code: contains all code concern the Gateway application and how to manage the data that come from devices.
  - Chipers: it contains the source code and the program for xor encryption algorithm to apply on data from a device.
  - DeviceFiles: it contains files that simulate some data calculated by a device. The device “11:a3” are data from the testbed, the others are example used for testing.
  - Device_cript: it contains files that represent the encrypted and hash version of what can be find in DeviceFiles. These are the simulation of the data that come from devices.
  - LogFiles: contains an example of logs that come from an experiment in the testbed.
  - ParsingCode: it contains the python code “ParsingLogDeviceData.py” for creating the appropriate device data from log and the python code “ParsingLogTimeRanging.py” for computing the average time to calculate the ranging, putting this information in “AVGtemp[idDevice].txt” (idDevice is the id of considered device).
  - main.go: source code of the Gateway application.
  - Position_test.go: test code for testing the correctness, privacy and security of some application functionality.
  - TimeAppGateway.txt: List of times that indicates the calculated periods of a cycle of target position computation.
  - pkg: packages used in the Gateway application.
* bin: command used in the creation of the blockchain.
* calliper-workspace: it contains the necessary code to evaluate the performance of blockchain code with caliper.
  - benchmarks: it contains the configuration files for the different test on chaincode.
    + myDeviceBenchmark.yaml: configuration file for reading device test.
    + ReadTargetBenchmark.yaml: configuration file for reading target test.
    + addObsTest.yaml: configuration file for adding device observation test.
    + updateEvTrustRepTest.yaml: configuration file for test the updating of evidence, reputation and trust.
    + positionTest: configuration file for position computation test
    + UpdateDTest.yaml: configuration file for updating device test.
  - networks: it contains the configuration file for the network, organizations and accounts to consider in caliper test.
  - node_modules: it contains node.js modules utilized in caliper test.
  - workload: it contains the code used for test.
    + readDevice.js: test code for reading device.
    + readTarget.js: test code for reading target.
    + addObs.js: test code for adding the device observation.
    + updateEvTrustRep.js: test code for updating the evidence, reputation and trust.
    + position.js: test code for position computation.
    + UpdateDevice.js: test code for updating device.
* Chaincode_dir: it contains the logic of chaincode, that is the project contract.
  - main.go: source code used for installation of the contract in the blockchain.
  - collections_config.json: configuration file that describes the collections used in the blockchain.
  - chaincode: it contains the code of the contract.
    + Data_processing.go: It contains all functionalities of the contract that changes the state of the blockchain.
    + Data_request.go: It contains all functionalities of the contract that ask data from the blockchain.
* Conf_code: it is the directory of Admin application code.
  - main.go: source code of the Admin application.
  - Conf_test.go: test code for testing the correctness, privacy and security of some application functionality.
* config: it contains configuration files for the blockchain network.
  - core.yaml and configtx.yaml: channel configuration files.
  - orderer.yaml: order configuration file.
* User_code: it is the directory of User application code.
  - main.go: source code of the User application.
* uwb-rng-radio-solution: it contains the source codes to do ranging experiment in testbed.
  - rng-init.c: source code for every device for doing ranging operation. Must change “linkaddr_t init” for every device.
  - rng-resp.c: source code for the target for doing ranging operation.
* Report: document that describes the project.
* Presentation: slides that describe the project.
* Network: it contains all the files for the network aspect of the blockchain.
  - bft-config: it contains the configuration file for the blockchain channel in case it is used the PBFT consensus algorithm.
  - compose: it contains the configuration files to build docker containers.
  - configtx: it contains the configuration file for the blockchain channel in general case.
  - initialize_for_test.sh: script for building the Hyperledger fabric structure in case it is used the bft consensus algorithm and the “cryptogen” command for the creation of the cryptographic material. Moreover, the script deploys the contract in the blockchain.
  - monitordocker.sh: script to check the containers
  - network.sh: general script to execute various operations with Hyperledger Fabric deployed.
  - scripts: it contains some scripts to operate with the built blockchain.
  - setOrgEnv.sh: script for update some environment variable.
  - system-genesis-block: it will contain the blockchain genesis block.
  - organizations: Currently contains the configuration files that will be used in the initialization of the blockchain for the cryptographic material creation of the organizations members and the order nodes. The cryptographic material will be put in this directory.

### TEST

To test the project you must follow the following instructions:
1. Clone the project.
2. Install all [prerequisites](https://hyperledger-fabric.readthedocs.io/en/latest/prereqs.html).
3. Download the necessary docker images:
   1. `curl -sSLO https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/install-fabric.sh && chmod +x install-fabric.sh`
   2. `./install-fabric.sh d`
4. Go to the network directory: cd Path/to/Network  .
5. Build a Hyperledger Fabric blockchain and deploy the smart contract "PositionContract": `./initialize_for_test.sh`
6. For test Gateway and User application:
   1. `cd ../Application_code`
   2. `go test main.go Position_test.go -v`
7. For test Admin application:
   1. `cd ../Conf_code`
   2. `go test main.go Conf_test.go -v`
8. For test the chaincode performance:
   1. `cd ../caliper-workspace`
   2. Install the prerequisites: Node-version v12.22.10 and NPM version 6.14.16
   3. Install caliper:
      1. `npm init`
      2. `npm install --only=prod @hyperledger/caliper-cli@0.5.0`
      3. `npx caliper bind --caliper-bind-sut fabric:2.2`
   4. For test the device reading: `npx caliper launch manager --caliper-workspace ./ --caliper-networkconfig networks/networkConfig.yaml --caliper-benchconfig benchmarks/myDeviceBenchmark.yaml —caliper-flow-only-test`
   5. For test the target reading: `npx caliper launch manager --caliper-workspace ./ --caliper-networkconfig networks/networkConfig.yaml --caliper-benchconfig benchmarks/ReadTargetBenchmark.yaml—caliper-flow-only-test`
   6. For test an add of a device observation: `npx caliper launch manager --caliper-workspace ./ --caliper-networkconfig networks/networkConfig.yaml --caliper-benchconfig benchmarks/addObsTest.yaml —caliper-flow-only-test`
   7. For test an update of evidence, reputation and trust: `npx caliper launch manager --caliper-workspace ./ --caliper-networkconfig networks/networkConfig.yaml --caliper-benchconfig benchmarks/updateEvTrustRepTest.yaml —caliper-flow-only-test`
   8. For test a position computation: `npx caliper launch manager --caliper-workspace ./ --caliper-networkconfig networks/networkConfig.yaml --caliper-benchconfig benchmarks/positionTest.yaml —caliper-flow-only-test`
   9. For test an update of a device by an admin: `npx caliper launch manager --caliper-workspace ./ --caliper-networkconfig networks/networkConfig.yaml --caliper-benchconfig benchmarks/UpdateDTest.yaml —caliper-flow-only-test`
   10. To see the result of one of this test open "report.html".
9. For do some experiments with testbed, that is a ranging operation between a device and a target:
   1. Register to [CLOVES](https://iottestbed.disi.unitn.it/cloves/getting-started/) testbed and install Cooja.
   2. `cd ../uwb-rng-radio-solution`
   3. Go to website of CLOVES and select [map](https://iottestbed.disi.unitn.it/cloves/infrastructure/) of devices you want to utilize, and download the file with devices information of that map.
   4. Substitute in rng-init.c “linkaddr_t init” with wanted device address (as you can see as an example in the c file) and “linkaddr_t resp_list[NUM_DEST]” with wanted target address (as you can see as sn example in the c file).
   5. Compile the source codes for the ranging: `make TARGET=evb1000`
   6. Enter to the CLOVES website.
   7. Go to "Create Job"
   8. In Timeslot info select an island (the map of devices), as start time indicate ASAP and as duration indicate how much time do you want the experiment last.
   9. In binary file 1 insert as hardware "evb1000", as Bin file “rng-init.bin”, as "targets" the id of node device in respective island.
   10. In binary file 2 insert as hardware "evb1000", as Bin file “rng-resp.bin”, as "targets" the id of node target in respective island.
   11. Start the experiment
   12. After the time indicated in Timeslot info, go to the "Download jobs". Download the experiment result and extract the log file.
   13. Then, you can do some operation on the job file "job.log" obtained: `cd ../Application_code`
       1. You can use the parsing code “ParsingCode/ParsingLogDeviceData.py”, passing the log file, to extract the wanted data in "/DeviceFiles/device[idDevice].txt" where the value of idDevice depend on device id: `ParsingCode/ParsingLogDeviceData.py /PATH/TO/job.log`
       2. You can use the parsing code “ParsingCode/ParsingLogTimeRanging.py”, passing the log file, to calculate the average time for a ranging operation in AVGtemp[idDevice].txt where the value of idDevice depend on device id: `ParsingCode/ParsingLogTimeRanging.py /PATH/TO/job.log`
       3. You can take the data received from the device "/DeviceFiles/device[idDevice].txt" and encrypt it with the key (a character) in "/Devices_cript/device[idDevice]_ript.txt": `chipers/xor_c "/DeviceFiles/device[idDevice].txt" "/Devices_cript/device[idDevice]_ript.txt" "in" "key"`
       4. You can use the encrypted file for possible test of Gateway and Admin application.
10. To prove Gateway application:
    1. `cd ../Application_code`
    2. Start the application until it doesn’t print anything more: `go main.go`
    3. It can also be seen in TimeAppGateway.txt a list of time periods calculated for every cycle of position target calculation.
11. To prove Admin application:
    1. `cd ../Conf_code`
    2. Start the application and prove some operation suggested: `go main.go`
12. To prove User application:
    1. `cd ../User_code`
    2. Start the application and prove some operation suggested: `go main.go`
13. To stop the test:
    1. `cd ../Network`
    2. `./network.sh down`

