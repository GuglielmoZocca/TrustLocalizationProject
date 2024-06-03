## GOAL

The aim of our research is to enhance a blockchain structure by incorporating data privatization features, drawing inspiration from the paper [A Trust Architecture for Blockchain in IoT](https://arxiv.org/pdf/1906.11461). We aim to demonstrate these novel characteristics through an example application—a prototype illustrating an end-to-end process from devices to users for calculating the position of a target. Our goal is to achieve data privatization within the architecture, ensuring that experiment application data stored in the blockchain are accessible only to those responsible for the experiment or authorized individuals. We strive to achieve this property without compromising the inherent advantages of blockchain technology, including auditability, integrity, and authenticity.

### CODE

In this section, the project's directory structure and the contents of each directory and file are explained.:
* Application_code: contains all the code related to the Gateway application, including how to manage the data received from devices.
  - Chipers: It includes the source code and program for the XOR encryption algorithm, which is applied to data from a device.
  - DeviceFiles: It contains files that simulate data calculated by a device. The data from the device "11:a3" are actual data from the testbed, while the others are examples used for testing purposes.
  - Device_cript: It contains files that represent the encrypted and hashed versions of the data found in DeviceFiles, simulating the data transmitted from devices.
  - LogFiles: It contains an example of logs generated from an experiment conducted in the testbed.
  - ParsingCode: It contains the Python code "ParsingLogDeviceData.py," which generates the appropriate device data from logs, and "ParsingLogTimeRanging.py," which computes the average time to calculate the ranging. The resulting information is stored in "AVGtemp[idDevice].txt," where idDevice is the ID of the considered device.
  - main.go: It is the source code of the Gateway application.
  - Position_test.go: It is the test code designed to verify the correctness, privacy, and security of various application functionalities.
  - TimeAppGateway.txt: This file contains a list of times representing the calculated durations of a target position computation cycle. These times account for the addition of observations as well as the updates to evidence, reputation, and trust.
  - pkg: They are the Golang packages used in the Gateway application.
* bin: These are the commands used for creating the blockchain.
* calliper-workspace: It contains the necessary code to evaluate the performance of blockchain code using Caliper.
  - benchmarks: It contains the configuration files for various chaincode tests.
    + myDeviceBenchmark.yaml: It is the configuration file for reading device test.
    + ReadTargetBenchmark.yaml: It is the configuration file for reading target test.
    + addObsTest.yaml: It is the configuration file for adding device observation an confidence test.
    + updateEvTrustRepTest.yaml: It is the configuration file for test the updating of evidence, reputation and trust.
    + positionTest: It is the configuration file for position computation test
    + UpdateDTest.yaml: It is the configuration file for updating device test.
  - networks: It contains the configuration file detailing the network, organizations, and accounts to be considered in the Caliper tests.
  - node_modules: it contains node.js modules utilized in Caliper tests.
  - workload: it contains the code used for Caliper tests.
    + readDevice.js: It is the test code for reading device.
    + readTarget.js: It is the test code for reading target.
    + addObs.js: It is the test code for adding the device observation and confidence.
    + updateEvTrustRep.js: It is the test code for updating the evidence, reputation and trust.
    + position.js: It is the test code for position computation.
    + UpdateDevice.js: It is the test code for updating device.
* Chaincode_dir: It houses the logic of the chaincode, constituting the project contract.
  - main.go: It is the source code utilized for installing the contract in the blockchain.
  - collections_config.json: It is the configuration file that describes the collections used in the blockchain.
  - chaincode: It contains the source codes of the contract.
    + Data_processing.go: It is the source code that contains all functionalities of the contract that changes the state of the blockchain.
    + Data_request.go: It is the source code encompassing all functionalities of the contract that request data from the blockchain.
* Conf_code: It is the directory of Admin application code.
  - main.go: It is the source code of the Admin application.
  - Conf_test.go: It is the test code for testing the correctness, privacy and security of application functionalities.
* config: It contains configuration files for the blockchain network.
  - core.yaml and configtx.yaml: They are the channel configuration files.
  - orderer.yaml: It is the order configuration file.
* User_code: It is the directory of User application code.
  - main.go: It is the source code of the User application.
* uwb-rng-radio-solution: It contains the source codes to do ranging experiment in testbed.
  - rng-init.c: It consists of the source code for each device for conducting ranging operations. The "linkaddr_t init" must be changed for each device.
  - rng-resp.c: It is the source code for the target for conducting ranging operation.
* Report.pdf: It is the document that describes the project.
* Presentation: They are the slides that describe the project.
* Network: It includes all the files related to the network aspect of the blockchain.
  - bft-config: It includes the configuration file for the blockchain channel specifically designed for the PBFT consensus algorithm.
  - compose: It contains the configuration files to build docker containers.
  - configtx: It contains the configuration file for the blockchain channel in the general case.
  - initialize_for_test.sh: It is the script for building the Hyperledger Fabric structure when utilizing the BFT consensus algorithm. Additionally, the script includes the "cryptogen" command for generating cryptographic material. Furthermore, it deploys the contract in the blockchain.
  - monitordocker.sh: It is the script to check the containers.
  - network.sh: It ia the general script to execute various operations with Hyperledger Fabric.
  - scripts: it contains some scripts to operate with the blockchain.
  - setOrgEnv.sh: It is the script for updating certain environment variables
  - system-genesis-block: it will contain the blockchain genesis block.
  - organizations: Currently, this directory contains the configuration files essential for initializing the blockchain. These files are utilized in creating cryptographic material for the members of organizations and the order nodes. The resulting cryptographic material will be stored within this directory.

### TEST

To test the project you must follow the following instructions:
1. Clone the project.
2. Install all [prerequisites](https://hyperledger-fabric.readthedocs.io/en/latest/prereqs.html).
3. Download the necessary docker images:
   1. `curl -sSLO https://raw.githubusercontent.com/hyperledger/fabric/main/scripts/install-fabric.sh && chmod +x install-fabric.sh`
   2. `./install-fabric.sh d`
4. Go to the network directory: cd Path/to/Network  .
5. Build a Hyperledger Fabric blockchain and deploy the smart contract "PositionContract": `./initialize_for_test.sh`
6. For testing Gateway and User applications:
   1. `cd ../Application_code`
   2. `go test main.go Position_test.go -v`
7. For testing Admin application:
   1. `cd ../Conf_code`
   2. `go test main.go Conf_test.go -v`
8. For testing the chaincode performance:
   1. `cd ../caliper-workspace`
   2. Install the prerequisites: Node-version v12.22.10 and NPM version 6.14.16
   3. Install caliper:
      1. `npm init`
      2. `npm install --only=prod @hyperledger/caliper-cli@0.5.0`
      3. `npx caliper bind --caliper-bind-sut fabric:2.2`
   4. For testing reading device data: `npx caliper launch manager --caliper-workspace ./ --caliper-networkconfig networks/networkConfig.yaml --caliper-benchconfig benchmarks/myDeviceBenchmark.yaml —caliper-flow-only-test`
   5. For testing reading target data: `npx caliper launch manager --caliper-workspace ./ --caliper-networkconfig networks/networkConfig.yaml --caliper-benchconfig benchmarks/ReadTargetBenchmark.yaml—caliper-flow-only-test`
   6. For testing adding a device observation and confidence: `npx caliper launch manager --caliper-workspace ./ --caliper-networkconfig networks/networkConfig.yaml --caliper-benchconfig benchmarks/addObsTest.yaml —caliper-flow-only-test`
   7. For testing updating evidence, reputation and trust: `npx caliper launch manager --caliper-workspace ./ --caliper-networkconfig networks/networkConfig.yaml --caliper-benchconfig benchmarks/updateEvTrustRepTest.yaml —caliper-flow-only-test`
   8. For testing position computation: `npx caliper launch manager --caliper-workspace ./ --caliper-networkconfig networks/networkConfig.yaml --caliper-benchconfig benchmarks/positionTest.yaml —caliper-flow-only-test`
   9. For testing updating device by an admin: `npx caliper launch manager --caliper-workspace ./ --caliper-networkconfig networks/networkConfig.yaml --caliper-benchconfig benchmarks/UpdateDTest.yaml —caliper-flow-only-test`
   10. To see the result of one of these tests, open "report.html".
9. For conducting experiments with the testbed, which involves ranging operations between a device and a target:
   1. Register to [CLOVES](https://iottestbed.disi.unitn.it/cloves/getting-started/) testbed and install Cooja.
   2. `cd ../uwb-rng-radio-solution`
   3. Go to the CLOVES website and navigate to the [device map](https://iottestbed.disi.unitn.it/cloves/infrastructure/) section. Choose the specific map of devices you intend to utilize and download the corresponding file containing device information.
   4. Replace "linkaddr_t init" in rng-init.c with the desired device address, as shown in the example in the C file. Also, replace "linkaddr_t resp_list[NUM_DEST]" with the desired target address, following the example provided in the C file.
   5. Compile the source codes for the ranging: `make TARGET=evb1000`
   6. Enter to the CLOVES website.
   7. Go to "Create Job".
   8. Choose an island (the map of devices) in Timeslot info. Set the start time to ASAP and specify the duration for the experiment.
   9. In binary file 1 insert as hardware "evb1000", as Bin file “rng-init.bin”, as "targets" the id of node device in respective island.
   10. In binary file 2 insert as hardware "evb1000", as Bin file “rng-resp.bin”, as "targets" the id of node target in respective island.
   11. Start the experiment.
   12. After the time indicated in Timeslot info, go to the "Download jobs". Download the experiment results and extract the log file.
   13. Then, you can do some operation on the job file "job.log" obtained: `cd ../Application_code`
       1. You can utilize the parsing code located at "ParsingCode/ParsingLogDeviceData.py". By running this code and providing the log file as input, you can extract the desired device data. The extracted data will be stored in "/DeviceFiles/device[idDevice].txt", where the value of idDevice corresponds to the device ID: `ParsingCode/ParsingLogDeviceData.py /PATH/TO/job.log`
       2. You can utilize the parsing code "ParsingCode/ParsingLogTimeRanging.py" by providing the log file as input to calculate the average time for a ranging operation. The results will be stored in "AVGtemp[idDevice].txt", where the value of idDevice corresponds to the device ID.: `ParsingCode/ParsingLogTimeRanging.py /PATH/TO/job.log`
       3. You can take the data received from the device "/DeviceFiles/device[idDevice].txt",encrypt it with the key (a character) and store it in "/Devices_cript/device[idDevice]_ript.txt": `chipers/xor_c "/DeviceFiles/device[idDevice].txt" "/Devices_cript/device[idDevice]_ript.txt" "in" "key"`
       4. You can use the encrypted file for possible tests of Gateway and Admin applications.
10. To simulate Gateway application:
    1. `cd ../Application_code`
    2. Initiate the application and continue until it ceases printing any further output: `go main.go`
    3. You can also refer to the "TimeAppGateway.txt" file for a comprehensive list of time periods calculated for each cycle of target position computation, taking into account the addition of observations and the updating of evidence, reputation, and trust.
11. To simulate Admin application:
    1. `cd ../Conf_code`
    2. Start the application and try some operation suggested: `go main.go`
12. To prove User application:
    1. `cd ../User_code`
    2. Start the application and try some operation suggested: `go main.go`
13. To stop the test network:
    1. `cd ../Network`
    2. `./network.sh down`

