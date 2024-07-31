import chalk from "chalk" // Chain styles
import clear from "clear" // Clear screen
import figlet from "figlet" // Create ASCII Art from text 
import fs from 'fs'
import _ from 'lodash'
import inquirer from 'inquirer'
import { parse } from 'csv-parse'

import { InfluxDB, Point, HttpError, DEFAULT_WriteOptions } from '@influxdata/influxdb-client'
import { url, token, org, bucket, dynonDir } from './env.mjs'


const flushBatchSize = DEFAULT_WriteOptions.batchSize


// advanced write options
const writeOptions = {
    /* the maximum points/lines to send in a single batch to InfluxDB server */
    batchSize: flushBatchSize + 1, // don't let automatically flush data

    /* maximum time in millis to keep points in an unflushed batch, 0 means don't periodically flush */
    flushInterval: 0,

    /* maximum size of the retry buffer - it contains items that could not be sent for the first time */
    maxBufferLines: 30_000,

    /* the count of internally-scheduled retries upon write failure, the delays between write attempts 
    follow an exponential backoff strategy if there is no Retry-After HTTP header */
    maxRetries: 0, // do not retry writes
    // ... there are more write options that can be customized, see
    // https://influxdata.github.io/influxdb-client-js/influxdb-client.writeoptions.html and
    // https://influxdata.github.io/influxdb-client-js/influxdb-client.writeretryoptions.html
}


// Node.js HTTP client OOTB does not reuse established TCP connections, a custom node HTTP agent
// can be used to reuse them and thus reduce the count of newly established networking sockets
import { Agent } from 'http'
const keepAliveAgent = new Agent({
    keepAlive: true, // reuse existing connections
    keepAliveMsecs: 60 * 1000, // 60 seconds keep alive
})
process.on('exit', () => keepAliveAgent.destroy())

// create InfluxDB with a custom HTTP agent
const influxDB = new InfluxDB({
    url,
    token,
    transportOptions: {
        agent: keepAliveAgent,
    },
})


async function analyzeFile(files) {

    // For file in files
    for (const file of files) {

        var csvFile = dynonDir + file

        fs.readFile(csvFile, (err, data) => {

            if (err) throw err

            // Empty array resulting the csv parsing
            const records = []

            // Initialize the parser params
            // The Dynon csv data is not clean especially in shutdown and columns heading change over version so:
            const parser = parse({
                // We want le columns heading for being consistant over Dynon versions
                columns: true,

                // Dynon csv column count might have corruption (Dynon shutdown)
                skip_records_with_error: true,

                // Dynon csv headers are not compatible with db column naming convention
                columns: header => header.map(column => column.replace(/[ ()&/%]/g, '')),

                // Don't generate records for lines having empty values
                skip_records_with_empty_values: true,
            })

            // Use the readable stream api to consume records
            parser.on('readable', function () {
                let record
                while ((record = parser.read()) !== null) {
                    records.push(record)
                }
            })

            // Catch any parsing error and report it
            parser.on('error', function (err) {
                console.error(err.message)
            })

            // Finalyze the parsed records
            // As a reminder the records hold datas having a valid GPS fix and Speed up or aqual to 10 kts 
            parser.on('end', function () {
                try {
                    let influxCnt = 0

                    const writeApi = influxDB.getWriteApi(org, bucket, 's', writeOptions)

                    var recordTime = 0
                    var lastRecordTime = 0

                    console.log(file, "- Found:", records.length, "records")

                    records.forEach(async function (record, index) {

                        // Record logic require only one key to be analyzed 
                        var gpsSpeed = Number(record['GroundSpeedknots'])
                        
                        recordTime = record['GPSDateTime']

                        // if GPS speed up to 10 we save in db
                        if (gpsSpeed >= 10 && recordTime !== lastRecordTime) {

                            // We might have multiple record at the same time. So save the first one only ...
                            lastRecordTime = record['GPSDateTime']


                            // Count the records to save
                            influxCnt++

                            // console.log(lastRecordTime)

                            // Dispatch record into 2 buckets : avionic & engine
                            // Avionic
                            const pointAvionic = new Point('avionic')
                                .floatField('lat', record['Latitudedeg'])
                                .floatField('lon', record['Longitudedeg'])
                                .intField('alt', record['GPSAltitudefeet'])
                                .intField('GroundSpeedknots', record['GroundSpeedknots'])
                                .floatField('Pitchdeg', record['Pitchdeg'])
                                .floatField('Rolldeg', record['Rolldeg'])
                                .floatField('MagneticHeadingdeg', record['MagneticHeadingdeg'])
                                .intField('IndicatedAirspeedknots', record['IndicatedAirspeedknots'])
                                .intField('TrueAirspeedknots', record['TrueAirspeedknots'])
                                .intField('IndicatedAirspeedknots', record['IndicatedAirspeedknots'])
                                .floatField('LateralAccelg', record['LateralAccelg'])
                                .floatField('VerticalAccelg', record['VerticalAccelg'])
                                .floatField('VerticalSpeedftmin', record['VerticalSpeedftmin'])
                                .floatField('OATdegC', record['OATdegC'])
                                .timestamp(new Date(record['GPSDateTime']))

                            // engine
                            const pointEngine = new Point('engine')
                                .floatField('OilPressurePSI:', record['OILPRESSUREPSI'])
                                .intField('OilTempdegC', record['OilTempdegC'])
                                .intField('RPM', record['RPML'])
                                .floatField('ManifoldPressureinHg', record['ManifoldPressureinHg'])
                                .floatField('FuelFlow1galhr', record['FuelFlow1galhr'])
                                .floatField('FuelPressurePSI', record['FuelPressurePSI'])
                                .floatField('Volts1', record['Volts1'])
                                .floatField('Amps', record['Amps'])
                                .intField('EGT1degC', record['EGT1degC'])
                                .intField('EGT2degC', record['EGT2degC'])
                                .intField('ChtLeftTemperaturedegC', record['CHTLTEMPERATUREdegC'])
                                .intField('ChtRightTemperaturedegC', record['CHTRTEMPERATUREdegC'])
                                .timestamp(new Date(record['GPSDateTime']))

                            // Write points into Influx buffer
                            writeApi.writePoint(pointAvionic)
                            writeApi.writePoint(pointEngine)

                            // control the way of how data are flushed
                            if ((influxCnt + 1) % flushBatchSize === 0) {
                                console.log(`flush writeApi: chunk #${(influxCnt + 1) / flushBatchSize}`)
                                try {
                                    // write the data to InfluxDB server, wait for it
                                    await writeApi.flush()
                                } catch (e) {
                                    console.error()
                                }
                            }

                        }
                        // Force buffer flush so db receive the unflushed points
                        writeApi.close()

                    })

                    console.log(file, "- Found:", records.length, "records", "Inserted:", influxCnt)

                } catch (err) {
                    console.log(err)
                }
            })

            // Write csv data into the parser and close stream
            parser.write(data.toString())
            parser.end()

        })

    }
}


async function askFiles(filelist) {
    const questions = [
        {
            type: 'checkbox',
            name: 'select',
            message: 'Select the files you wish to import influxDB:',
            choices: filelist,
            default: []
        }
    ]
    return inquirer.prompt(questions)
}


try {

    // let’s clear the screen and then display a banner:
    clear()
    console.log(
        chalk.yellow.bold(
            figlet.textSync('Dude Loader', { horizontalLayout: 'full' })
        )
    )

    // console.log(flushBatchSize)
    
    // Find files in dynon directory
    const filelist = _.without(fs.readdirSync(dynonDir), '.DS_Store')
    // console.log(filelist)

    // Cleanup for only supported files
    var reg1 = new RegExp("ALERT_DATA|DIAGNOSTIC", "g")
    const cleanfilelist = _.remove(filelist, function (v) { return !v.match(reg1) })
    // console.log(cleanfilelist)

    // Allow Dynon csv file(s) selection
    const selectedFiles = await askFiles(cleanfilelist)
    // console.log(selectedFiles)

    // Execute the parsing and later db import
    await analyzeFile(selectedFiles.select)

} catch (e) {
    console.error(e)
    console.log('FINISHED ERROR')
}