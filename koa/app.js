const log4js = require('log4js')
const logger = log4js.getLogger();
const Koa = require('koa')
const app = new Koa()

require("./db")(app)
const compress = require('koa-compress')

const views = require('koa-views')
const json = require('koa-json')
const onerror = require('koa-onerror')
const bodyparser = require('koa-bodyparser')
const koaLogger = require('koa-logger')

const index = require('./routes/index')
const users = require('./routes/users')
const stats = require('./routes/stats')

// error handler
onerror(app)

// middlewares
app.use(bodyparser({
  enableTypes:['json', 'form', 'text']
}))
app.use(json())
app.use(koaLogger())
app.use(require('koa-static')(__dirname + '/public'))

app.use(views(__dirname + '/views', {
  extension: 'pug'
}))
//response compress
app.use(compress({
  filter: function (content_type) {
	logger.debug(`content_type:${content_type}`)
	return true
  	//return /text/i.test(content_type)
  },
  threshold: 2048,
  flush: require('zlib').Z_SYNC_FLUSH
}))


// logger
app.use(async (ctx, next) => {
/*
  if(ctx.method === "POST" ){
	  logger.debug(`ctx.request.body:${JSON.stringify(ctx.request.body, null,2)}`)
  }
*/
  const start = new Date()
  await next()
  const ms = new Date() - start
  logger.debug(`${ctx.method} ${ctx.url} - ${ms}ms`)
})

// routes
//app.use(index.routes(), index.allowedMethods())
//app.use(users.routes(), users.allowedMethods())
app.use(stats.routes(), stats.allowedMethods())

// error-handling
app.on('error', (err, ctx) => {
  console.error('server error', err, ctx)
});



module.exports = app
