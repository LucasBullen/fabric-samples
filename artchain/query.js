'use strict';
/*
* Copyright IBM Corp All Rights Reserved
*
* SPDX-License-Identifier: Apache-2.0
*/
/*
 * Hyperledger Fabric Sample Query Program
 */

var query=function(functionString, callback){//, args){
//queryArtById
//queryAllArt
var fcn_param = functionString;//'queryArtById';
var args_param = []//args;//['6b86b273ff34fce19d6b804eff5a3f574'];

var hfc = require('fabric-client');
var path = require('path');

var options = {
    wallet_path: path.join(__dirname, './creds'),
    user_id: 'PeerAdmin',
    channel_id: 'mychannel',
    chaincode_id: 'artchain',
    network_url: 'grpc://localhost:7051',
};

var channel = {};
var client = null;

Promise.resolve().then(() => {
//    console.log("Create a client and set the wallet location");
    client = new hfc();
    return hfc.newDefaultKeyValueStore({ path: options.wallet_path });
}).then((wallet) => {
  //  console.log("Set wallet path, and associate user ", options.user_id, " with application");
    client.setStateStore(wallet);
    return client.getUserContext(options.user_id, true);
}).then((user) => {
  //  console.log("Check user is enrolled, and set a query URL in the network");
    if (user === undefined || user.isEnrolled() === false) {
        console.error("User not defined, or not enrolled - error");
    }
    channel = client.newChannel(options.channel_id);
    channel.addPeer(client.newPeer(options.network_url));
    return;
}).then(() => {
    //console.log("Make query:", fcn_param);
    var transaction_id = client.newTransactionID();
    //console.log("Assigning transaction_id: ", transaction_id._transaction_id);
    const request = {
        chaincodeId: options.chaincode_id,
        txId: transaction_id,
        fcn: fcn_param,
        args: args_param
    };
    return channel.queryByChaincode(request);
}).then((query_responses) => {
    //console.log("returned from query");
    if (!query_responses.length) {
        console.log("false");
    } else {
        //console.log("Query result count = ", query_responses.length)
    }
    if (query_responses[0] instanceof Error) {
        console.error("error from query = ", query_responses[0]);
    }
    console.log("before responce");    
    callback(JSON.parse(query_responses[0].toString()));
    return query_responses[0].toString();
    //console.log( query_responses[0].toString());
}).catch((err) => {
    console.error("Caught Error", err);
});
};

module.exports.query=query;
