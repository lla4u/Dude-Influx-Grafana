/** InfluxDB v2 URL */
const url = process.env['INFLUX_URL'] || 'http://dude-influxdb:8086'

/** InfluxDB authorization token */
const token = process.env['INFLUX_TOKEN'] || 'my-super-secret-auth-token'

/** Organization within InfluxDB  */
const org = process.env['INFLUX_ORG'] || 'dude'

/**InfluxDB bucket used in examples  */
const bucket = 'dude'

/**Dynon csv file directory */
// const dynonDir = process.env['DYNON_DIR'] || '/Users/lla/Documents/Laurent/Aviation/P300 Dude/'
const dynonDir = process.env['DYNON_DIR'] || '/home/node/sessions/'


export {url, token, org, bucket, dynonDir}