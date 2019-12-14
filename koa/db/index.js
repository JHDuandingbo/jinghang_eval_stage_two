const bluebird=require("bluebird")
const redis=require("redis")
bluebird.promisifyAll(redis);

const  mongoClient = require('mongodb').MongoClient;
const assert = require('assert');
//console.log(process.env.DB_USER)console.log(process.env.DB_PASSWORD)


module.exports = async function (app) {
	//const password = encodeURIComponent("fatcat@practice")
	const username = process.env.DB_USER
	const password = encodeURIComponent(process.env.DB_PASSWORD)
	const host = process.env.DB_HOST
	const MONGO_URI = `mongodb://${username}:${password}@${host}/admin`;
	console.log(MONGO_URI)
	//const dbName = "practice_question"
	const dbName = process.env.DB_NAME

	let conn = await mongoClient.connect(MONGO_URI);
	let myDB = conn.db(dbName); // 选择一个 db

	app.context.db = myDB
	app.context.table = {}
	app.context.table.userRequests = await myDB.collection("userRequests")
	console.log("Database conn established")

/*
	/////////////////////redis
	const rclient = redis.createClient()
	rclient.on("error", function (err) {
		 console.log("fatcat Error " + err);
	})
	app.context.rclient = rclient
*/

};
