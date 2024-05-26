'use strict';

//Test performance chaincode calculation target position

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

class MyWorkload extends WorkloadModuleBase {
    constructor() {
        super();
    }

    //Initialize leader
    async initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext) {
        await super.initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext);

        let device1Data = {deviceID: "12:22", coordinateX: 5, coordinateY:8, Key: "P", Neighborhood: ["12:21","12:23"], Reputation: 5};

        let tmapData = Buffer.from(JSON.stringify(device1Data));

        const assetID = `12:22`;
        console.log(`Worker ${this.workerIndex}: Creating device ${assetID}`);
        const request = {
            contractId: this.roundArguments.contractId,
            contractFunction: 'CreateDevice',
            invokerIdentity: 'Admin',
            contractArguments: ["DeviceAdmin1PrivateCollection"],
            targetOrganizations: ["Org1MSP"],
            transientMap: {Device_properties: tmapData},
            readOnly: false
        };

        await this.sutAdapter.sendRequests(request);

        let device2Data = {deviceID: "12:21", coordinateX: 3, coordinateY:2, Key: "3", Neighborhood: ["12:22","12:23"], Reputation: 5};

        let tmap2Data = Buffer.from(JSON.stringify(device2Data));

        const asset2ID = `12:21`;
        console.log(`Worker ${this.workerIndex}: Creating device ${asset2ID}`);
        const reques2t = {
            contractId: this.roundArguments.contractId,
            contractFunction: 'CreateDevice',
            invokerIdentity: 'Admin',
            contractArguments: ["DeviceAdmin1PrivateCollection"],
            targetOrganizations: ["Org1MSP"],
            transientMap: {Device_properties: tmap2Data},
            readOnly: false
        };

        await this.sutAdapter.sendRequests(reques2t);

        let device3Data = {deviceID: "12:23", coordinateX: 10, coordinateY:4, Key: "T", Neighborhood: ["12:21","12:22"], Reputation: 5};

        let tmap3Data = Buffer.from(JSON.stringify(device3Data));

        const asset3ID = `12:23`;
        console.log(`Worker ${this.workerIndex}: Creating device ${asset3ID}`);
        const reques3t = {
            contractId: this.roundArguments.contractId,
            contractFunction: 'CreateDevice',
            invokerIdentity: 'Admin',
            contractArguments: ["DeviceAdmin1PrivateCollection"],
            targetOrganizations: ["Org1MSP"],
            transientMap: {Device_properties: tmap3Data},
            readOnly: false
        };

        await this.sutAdapter.sendRequests(reques3t);

        let TargetData = { TargetID: "7"};

        let tarmapData = Buffer.from(JSON.stringify(TargetData));

        const targetID = `7`;
        console.log(`Worker ${this.workerIndex}: Creating target ${targetID}`);
        const reques4t = {
            contractId: this.roundArguments.contractId,
            contractFunction: 'CreateTarget',
            invokerIdentity: 'Admin',
            contractArguments: ["TargetOrg1PrivateCollection"],
            targetOrganizations: ["Org1MSP"],
            transientMap: {Target_properties: tarmapData},
            readOnly: false
        };

        await this.sutAdapter.sendRequests(reques4t);

        let deviceobsData = {deviceID: "12:22", observation: 3162, Confidence:1, MinConf: 0.4, MaxConf: 1};

        let tmapobsData = Buffer.from(JSON.stringify(deviceobsData));

        const assetIDu = `12:22`;
        console.log(`Worker ${this.workerIndex}: Add obs to ${assetIDu}`);
        const requestu = {
            contractId: this.roundArguments.contractId,
            contractFunction: 'UpdateDeviceObsConf',
            invokerIdentity: 'Admin',
            contractArguments: ["DeviceAdmin1PrivateCollection"],
            targetOrganizations: ["Org1MSP"],
            transientMap: {Device_properties: tmapobsData},
            readOnly: false
        };

        await this.sutAdapter.sendRequests(requestu);

        let deviceobs1Data = {deviceID: "12:21", observation: 4242, Confidence:1, MinConf: 0.4, MaxConf: 1};

        let tmapobs1Data = Buffer.from(JSON.stringify(deviceobs1Data));

        const asset1ID = `12:21`;
        console.log(`Worker ${this.workerIndex}: Add obs to ${asset1ID}`);
        const reques1t = {
            contractId: this.roundArguments.contractId,
            contractFunction: 'UpdateDeviceObsConf',
            invokerIdentity: 'Admin',
            contractArguments: ["DeviceAdmin1PrivateCollection"],
            targetOrganizations: ["Org1MSP"],
            transientMap: {Device_properties: tmapobs1Data},
            readOnly: false
        };

        await this.sutAdapter.sendRequests(reques1t);

        let deviceobs2Data = { deviceID: "12:23", observation: 4123, Confidence:1, MinConf: 0.4, MaxConf: 1};

        let tmapobs2Data = Buffer.from(JSON.stringify(deviceobs2Data));

        const asset2IDu = `12:23`;
        console.log(`Worker ${this.workerIndex}: Add obs to ${asset2IDu}`);
        const reques2tu = {
            contractId: this.roundArguments.contractId,
            contractFunction: 'UpdateDeviceObsConf',
            invokerIdentity: 'Admin',
            contractArguments: ["DeviceAdmin1PrivateCollection"],
            targetOrganizations: ["Org1MSP"],
            transientMap: {Device_properties: tmapobs2Data},
            readOnly: false
        };

        await this.sutAdapter.sendRequests(reques2tu);


        let deviceevData = {deviceID: "12:22", PRH: 2, PRL:1, threashConf: 0.7, threashEv: 0, maxRep: 20};

        let tmapevData = Buffer.from(JSON.stringify(deviceevData));
        const asset2IDu2 = `12:22`;
        console.log(`Worker ${this.workerIndex}: Adds ev to ${asset2IDu2}`);
        const reques3tu = {
            contractId: this.roundArguments.contractId,
            contractFunction: 'UpdateDeviceEvRep',
            invokerIdentity: 'Admin',
            contractArguments: ["DeviceAdmin1PrivateCollection"],
            targetOrganizations: ["Org1MSP"],
            transientMap: {Device_properties: tmapevData},
            readOnly: false
        };

        await this.sutAdapter.sendRequests(reques3tu);

        let deviceev1Data = {deviceID: "12:21", PRH: 2, PRL:1, threashConf: 0.7, threashEv: 0, maxRep: 20};

        let tmapev1Data = Buffer.from(JSON.stringify(deviceev1Data));
        console.log(`Worker ${this.workerIndex}: Adds ev to ${asset1ID}`);
        const reques4tu = {
            contractId: this.roundArguments.contractId,
            contractFunction: 'UpdateDeviceEvRep',
            invokerIdentity: 'Admin',
            contractArguments: ["DeviceAdmin1PrivateCollection"],
            targetOrganizations: ["Org1MSP"],
            transientMap: {Device_properties: tmapev1Data},
            readOnly: false
        };

        await this.sutAdapter.sendRequests(reques4tu);

        let deviceev2Data = {deviceID: "12:23", PRH: 2, PRL:1, threashConf: 0.7, threashEv: 0, maxRep: 20};

        let tmapev2Data = Buffer.from(JSON.stringify(deviceev2Data));

        console.log(`Worker ${this.workerIndex}: Adds ev to ${asset2ID}`);
        const reques5t = {
            contractId: this.roundArguments.contractId,
            contractFunction: 'UpdateDeviceEvRep',
            invokerIdentity: 'Admin',
            contractArguments: ["DeviceAdmin1PrivateCollection"],
            targetOrganizations: ["Org1MSP"],
            transientMap: {Device_properties: tmapev2Data},
            readOnly: false
        };

        await this.sutAdapter.sendRequests(reques5t);


    }

    //Calculate the target position
    async submitTransaction() {



        let PositionData = { TargetID: "7", ThreshErr: 0.01, DevicesUp: ["12:23","12:21","12:22"]};

        let tmapPotData = Buffer.from(JSON.stringify(PositionData));

        const targetID = `7`;
        console.log(`Worker ${this.workerIndex}: Calculate position of the target ${targetID}`);
        const reques6t = {
            contractId: this.roundArguments.contractId,
            contractFunction: 'PositionTarget',
            invokerIdentity: 'Admin',
            contractArguments: ["DeviceAdmin1PrivateCollection","TargetOrg1PrivateCollection","5/10/2022"],
            targetOrganizations: ["Org1MSP"],
            transientMap: {target_Posisiton: tmapPotData},
            readOnly: false
        };

        await this.sutAdapter.sendRequests(reques6t);




    }

    //Delete target and devices inserted
    async cleanupWorkloadModule() {

        let device1Data = {deviceID: "12:22"};

        let tmapData = Buffer.from(JSON.stringify(device1Data));

        const assetID = `12:22`;
        console.log(`Worker ${this.workerIndex}: Deleting device ${assetID}`);
        const reques1t = {
            contractId: this.roundArguments.contractId,
            contractFunction: 'DeleteDevice',
            invokerIdentity: 'Admin',
            targetOrganizations: ["Org1MSP"],
            contractArguments: ["DeviceAdmin1PrivateCollection"],
            transientMap: {device_delete: tmapData},
            readOnly: false
        };

        await this.sutAdapter.sendRequests(reques1t);

        let device2Data = {deviceID: "12:21"};

        let tmap2Data = Buffer.from(JSON.stringify(device2Data));

        const asset2ID = `12:21`;
        console.log(`Worker ${this.workerIndex}: Deleting device ${asset2ID}`);
        const reques2t = {
            contractId: this.roundArguments.contractId,
            contractFunction: 'DeleteDevice',
            invokerIdentity: 'Admin',
            targetOrganizations: ["Org1MSP"],
            contractArguments: ["DeviceAdmin1PrivateCollection"],
            transientMap: {device_delete: tmap2Data},
            readOnly: false
        };

        await this.sutAdapter.sendRequests(reques2t);

        let device3Data = {deviceID: "12:23"};

        let tmap3Data = Buffer.from(JSON.stringify(device3Data));

        const asset3ID = `12:23`;
        console.log(`Worker ${this.workerIndex}: Deleting device ${asset3ID}`);
        const reques3t = {
            contractId: this.roundArguments.contractId,
            contractFunction: 'DeleteDevice',
            invokerIdentity: 'Admin',
            targetOrganizations: ["Org1MSP"],
            contractArguments: ["DeviceAdmin1PrivateCollection"],
            transientMap: {device_delete: tmap3Data},
            readOnly: false
        };

        await this.sutAdapter.sendRequests(reques3t);

        let targetData = {targetID: "7"};

        let tmapTargetData = Buffer.from(JSON.stringify(targetData));

        const targetID = `7`;
        console.log(`Worker ${this.workerIndex}: Deleting target ${targetID}`);
        const reques4t = {
            contractId: this.roundArguments.contractId,
            contractFunction: 'DeleteTarget',
            invokerIdentity: 'Admin',
            targetOrganizations: ["Org1MSP"],
            contractArguments: ["TargetOrg1PrivateCollection"],
            transientMap: {target_delete: tmapTargetData},
            readOnly: false
        };

        await this.sutAdapter.sendRequests(reques4t);


    }
}

function createWorkloadModule() {
    return new MyWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;