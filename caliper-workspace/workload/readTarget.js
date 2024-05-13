'use strict';

//Test performance chaincode read Target

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

class MyWorkload extends WorkloadModuleBase {
    constructor() {
        super();
    }

    //Initialize ledger
    async initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext) {
        await super.initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext);

        let TargetData = {TargetID: "7"};

        let tarmapData = Buffer.from(JSON.stringify(TargetData));

        const targetID = `7`;
        console.log(`Worker ${this.workerIndex}: Creating target ${targetID}`);
        const request = {
            contractId: this.roundArguments.contractId,
            contractFunction: 'CreateTarget',
            invokerIdentity: 'Admin',
            contractArguments: ["TargetOrg1PrivateCollection"],
            targetOrganizations: ["Org1MSP"],
            transientMap: {Target_properties: tarmapData},
            readOnly: false
        };

        await this.sutAdapter.sendRequests(request);

    }

    //Read device
    async submitTransaction() {

        let targetData = {TargetID: "7"};

        let tmapTarData = Buffer.from(JSON.stringify(targetData));


        const myArgs = {
            contractId: this.roundArguments.contractId,
            contractFunction: 'ReadTarget',
            invokerIdentity: 'Admin',
            contractArguments: ["TargetOrg1PrivateCollection"],
            transientMap: {target_data: tmapTarData},
            readOnly: true
        };

        await this.sutAdapter.sendRequests(myArgs);
    }

    //Delete target inserted
    async cleanupWorkloadModule() {

        let targetData = {targetID: "7"};

        let tmapTargetData = Buffer.from(JSON.stringify(targetData));

        const targetID = `7`;
        console.log(`Worker ${this.workerIndex}: Deleting target ${targetID}`);
        const request = {
            contractId: this.roundArguments.contractId,
            contractFunction: 'DeleteTarget',
            invokerIdentity: 'Admin',
            targetOrganizations: ["Org1MSP"],
            contractArguments: ["TargetOrg1PrivateCollection"],
            transientMap: {target_delete: tmapTargetData},
            readOnly: false
        };

        await this.sutAdapter.sendRequests(request);
    }
}

function createWorkloadModule() {
    return new MyWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;