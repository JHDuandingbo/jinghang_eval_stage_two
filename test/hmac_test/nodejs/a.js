var crypto    = require('crypto');
var algorithm = 'sha1';   //consider using sha256
var hash, hmac;

var secret = 'JZ5J39vFncv3j3453X2G45sCy6cOv5G3';
var appId = "TiD3p6";
var text = 'v1.0%2CTiD3p6%2CHGTBv4hFj9%2C2018-08-29T15%3A06%3A56%2B0800%2C1eef7919-4af5-4a2c-bd6e-c2f60dde10db';
/*
// Method 1 - Writing to a stream
hmac = crypto.createHmac(algorithm, secret);    
hmac.write(text); // write in to the stream
hmac.end();       // can't read from the stream until you call end()
hash = hmac.read().toString('hex');    // read out hmac digest
console.log("Method 1: ", hash);
*/

// Method 2 - Using update and digest:
hexHmac = crypto.createHmac(algorithm, secret);
hexHmac.update(text);
hexHash = hexHmac.digest('hex');

hmac = crypto.createHmac(algorithm, secret);
hmac.update(text);
hash = hmac.digest('base64');
console.log("base64Hash: ", hash);
console.log(`base64Hash:<${hash}> len:${hash.length}`);
console.log(`hex:<${hexHash}>, bytelen:${hexHash.length/2}`);
