'use strict';

//Test performance chaincode update Device

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

class MyWorkload extends WorkloadModuleBase {
    constructor() {
        super();
    }

    //Initialize ledger
    async initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext) {
        await super.initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext);


        let device1Data = { deviceID: "12:22", coordinateX: 4.3, coordinateY:2.5, Key: "P", Neighborhood: ["3","4"], Reputation: 5};

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

    }

    //Update device
    async submitTransaction() {


        let deviceData = {deviceID:"12:22",Key:"T",Neighborhood:["1","3"]};

        let tmapData = Buffer.from(JSON.stringify(deviceData));


        const myArgs = {
            contractId: this.roundArguments.contractId,
            contractFunction: 'UpdateDevice',
            invokerIdentity: 'Admin',
            contractArguments: ["DeviceAdmin1PrivateCollection"],
            targetOrganizations: ["Org1MSP"],
            transientMap: {Device_properties: tmapData},
            readOnly: false
        };

        await this.sutAdapter.sendRequests(myArgs);
    }

    //Delete device inserted
    async cleanupWorkloadModule() {

        let device1Data = {deviceID: "12:22"};

        let tmapData = Buffer.from(JSON.stringify(device1Data));

        const assetID = `12:22`;
        console.log(`Worker ${this.workerIndex}: Deleting asset 12:22`);
        const request = {
            contractId: this.roundArguments.contractId,
            contractFunction: 'DeleteDevice',
            invokerIdentity: 'Admin',
            targetOrganizations: ["Org1MSP"],
            contractArguments: ["DeviceAdmin1PrivateCollection"],
            transientMap: {device_delete: tmapData},
            readOnly: false
        };

        await this.sutAdapter.sendRequests(request);
    }
}

function createWorkloadModule() {
    return new MyWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;