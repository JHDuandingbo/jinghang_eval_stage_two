'use strict'
const router = require('koa-router')()
const logger = require("../tools/logger")
const lodash = require("lodash")

router.prefix('/stats')
router.post('/', async function (ctx, next) {
		let rsp = {}
		rsp.code= 0
		rsp.msg= "ok"
		rsp.data= {}

		let body = ctx.request.body
		try{
				const userRequests = ctx.table.userRequests
				body.date=new Date()
				let returned = await userRequests.insert(body)
				rsp.data = returned
		}catch(err){
				logger.debug("err:", err)
				rsp.code = -1;
				rsp.msg = err.toString()
		}

		ctx.body = rsp;
})
router.get('/', async function (ctx, next) {
		let rsp = {}
		rsp.code= 0
		rsp.msg= "ok"
		rsp.data= {}

		//let query = ctx.request.query
		try{
				const userRequests = ctx.table.userRequests
				let returned = await userRequests.aggregate([
						{
								$group:{
										_id:"$requestKey",
										count:{$sum:1},
								}
						}

				]).toArray()
				rsp.data = returned
		}catch(err){
				logger.debug("err:", err)
				rsp.code = -1;
				rsp.msg = err.toString()
		}

		ctx.body = rsp;
})

module.exports=router;
