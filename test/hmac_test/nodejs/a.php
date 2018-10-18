<?php

$key = 'JZ5J39vFncv3j3453X2G45sCy6cOv5G3';
$message = 'v1.0%2CTiD3p6%2CHGTBv4hFj9%2C2018-08-29T15%3A06%3A56%2B0800%2C1eef7919-4af5-4a2c-bd6e-c2f60dde10db';
//$key = 'the shared secret key here';
//$message = 'the message to hash here';

// to lowercase hexits
hash_hmac('sha1', $message, $key);
//hash_hmac('sha256', $message, $key);

// to base64
base64_encode(hash_hmac('sha1', $message, $key, true));
