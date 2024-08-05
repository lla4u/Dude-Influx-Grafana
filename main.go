package main

import (
	"context"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gocarina/gocsv"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/joho/godotenv"
)

type Datalog struct {
	GpsFix                  string `csv:"GPS Fix Quality"`
	NumSatellites           string `csv:"Number of Satellites"`
	GpsDateTime             string `csv:"GPS Date & Time"`
	Lat                     string `csv:"Latitude (deg)"`
	Lon                     string `csv:"Longitude (deg)"`
	Alt                     string `csv:"GPS Altitude (feet)"`
	GroundSpeed_Knots       string `csv:"Ground Speed (knots)"`
	Pitch_Deg               string `csv:"Pitch (deg)"`
	Roll_Deg                string `csv:"Roll (deg)"`
	MagneticHeading_Deg     string `csv:"Magnetic Heading (deg)"`
	IndicatedAirspeed_Knots string `csv:"Indicated Airspeed (knots)"`
	LateralAccel_G          string `csv:"Lateral Accel (g)"`
	VerticalAccel_G         string `csv:"Vertical Accel(g)"`
	VerticalSpeed_ft_min    string `csv:"Vertical Speed (ft/min)"`
	OAT_Deg_C               string `csv:"OAT (deg C)"`
	TrueAirspeed_Knots      string `csv:"True Airspeed (knots)"`
	WindDirection_Deg       string `csv:"Wind Direction (deg)"`
	WindSpeed_Knots         string `csv:"Wind Speed (knots)"`
	Oil_Pressure_PSI        string `csv:"Oil Pressure (PSI)"`
	OilTemp_Deg_C           string `csv:"Oil Temp (deg C)"`
	RPM                     string `csv:"RPM L"`
	ManifoldPressure_inHg   string `csv:"Manifold Pressure (inHg)"`
	FuelFlow1_Gal_hr        string `csv:"Fuel Flow 1 (gal/hr)"`
	FuelPressure_PSI        string `csv:"Fuel Pressure (PSI)"`
	FuelRemaining_Gal       string `csv:"Fuel Remaining (gal)"`
	Volts                   string `csv:"Volts 1"`
	Amps                    string `csv:"Amps"`
	EGT1_Deg_C              string `csv:"EGT 1 (deg C)"`
	EGT2_Deg_C              string `csv:"EGT 2 (deg C)"`
	CHTL_Deg_C              string `csv:"CHTL TEMPERATURE (deg C)"`
	CHTR_Deg_C              string `csv:"CHTR TEMPERATURE (deg C)"`
}

var (
	file        string
	verboseFlag bool
)

// 5.57ms -> 600 records (read)
func main() {

	flag.StringVar(&file, "file", "test2.csv", "csv file name to import")
	flag.BoolVar(&verboseFlag, "verbose", false, "enable verbose")
	flag.Parse()

	fmt.Println("Starting")
	fmt.Println("Got file: ", file, "verbose", verboseFlag)

	if _, err := os.Stat(file); err != nil {
		log.Fatal("File does not exist")
	}

	// Just for timing import
	now := time.Now()

	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatal("Error loading config.env file")
	}

	// Create a new client using an InfluxDB server base URL and an authentication token
	client, err := ConnectToInfluxDB()
	if err != nil {
		log.Fatal("Error connecting influx db")
	}

	// Use blocking write client for writes to desired bucket
	writeAPI := client.WriteAPIBlocking("dude", "dude")

	readChannel := make(chan Datalog, 1)

	readFilePath := file

	// Open the CSV readFile
	readFile, err := os.OpenFile(readFilePath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer readFile.Close()

	csvCount := 0
	readFromCSV(readFile, readChannel)

	var gpsDateTime string
	influxCount := 0

	// Deal with the parsed csv records
	for r := range readChannel {

		if verboseFlag {

			fmt.Println("========================================")
			// 6
			fmt.Println("Fix:", r.GpsFix, "Sats:", r.NumSatellites, "Date:", r.GpsDateTime, "Lat:", r.Lat, "Lon:", r.Lon, "Alt:", r.Alt)
			// 6
			fmt.Println("GS:", r.GroundSpeed_Knots, "IAS:", r.IndicatedAirspeed_Knots, "TAS:", r.TrueAirspeed_Knots, "VSpeed:", r.VerticalSpeed_ft_min, "Wdir:", r.WindDirection_Deg, "WSpeed:", r.WindSpeed_Knots)
			// 6
			fmt.Println("Volts:", r.Volts, "Amps:", r.Amps, "CHTR:", r.CHTR_Deg_C, "CHTL:", r.CHTL_Deg_C, "EGT1:", r.EGT1_Deg_C, "EGT2:", r.EGT2_Deg_C)
			// 5
			fmt.Println("Pitch:", r.Pitch_Deg, "Roll:", r.Roll_Deg, "Mag:", r.MagneticHeading_Deg, "VertAccel:", r.VerticalAccel_G, "LatAccel:", r.LateralAccel_G)
			// 5
			fmt.Println("OAT:", r.OAT_Deg_C, "OilTemp:", r.OilTemp_Deg_C, "OilPress:", r.Oil_Pressure_PSI, "RPM:", r.RPM, "MAP:", r.ManifoldPressure_inHg)
			// 3
			fmt.Println("FuelPress:", r.FuelPressure_PSI, "FuelFlow:", r.FuelFlow1_Gal_hr, "FuelRemain:", r.FuelRemaining_Gal)
			// 		fmt.Println("========================================")
			fmt.Println()
		}

		//
		// Import Influxdb
		//

		// A valid gps Fix and Number of satellites up to 6 are required
		if (StringToInt(r.GpsFix) >= 1) && (StringToInt(r.NumSatellites) >= 6) {

			// Save the first record
			if r.GpsDateTime != gpsDateTime {

				p := influxdb2.NewPointWithMeasurement("datalog").
					AddField("lat", StringToFloat(r.Lat)).
					AddField("lon", StringToFloat(r.Lon)).
					AddField("alt", StringToInt(r.Alt)).
					//
					AddField("GS", StringToFloat(r.GroundSpeed_Knots)).
					AddField("IAS", StringToFloat(r.IndicatedAirspeed_Knots)).
					AddField("TAS", StringToFloat(r.TrueAirspeed_Knots)).
					//
					AddField("Volts", StringToFloat(r.Volts)).
					AddField("Amps", StringToFloat(r.Amps)).
					AddField("CHTR", StringToFloat(r.CHTR_Deg_C)).
					AddField("CHTL", StringToFloat(r.CHTL_Deg_C)).
					AddField("EGT1", StringToInt(r.EGT1_Deg_C)).
					AddField("EGT2", StringToInt(r.EGT2_Deg_C)).
					//
					AddField("Pitch", StringToFloat(r.Pitch_Deg)).
					AddField("Roll", StringToFloat(r.Roll_Deg)).
					AddField("Mag", StringToFloat(r.MagneticHeading_Deg)).
					//
					AddField("VertAccel", StringToFloat(r.VerticalAccel_G)).
					AddField("LatAccel", StringToFloat(r.LateralAccel_G)).
					//
					AddField("OAT", StringToInt(r.OAT_Deg_C)).
					//
					AddField("OilTemp", StringToInt(r.OilTemp_Deg_C)).
					AddField("OilPress", StringToInt(r.Oil_Pressure_PSI)).
					AddField("RPM", StringToInt(r.RPM)).
					AddField("MAP", StringToFloat(r.ManifoldPressure_inHg)).
					//
					AddField("FuelPress", StringToFloat(r.FuelPressure_PSI)).
					AddField("FuelFlow", StringToFloat(r.FuelFlow1_Gal_hr)).
					AddField("FuelRemaining", StringToFloat(r.FuelRemaining_Gal)).
					//
					SetTime(dateStringToUnix(r.GpsDateTime))

				err := writeAPI.WritePoint(context.Background(), p)

				if err != nil {
					log.Println(err)
				}
				gpsDateTime = r.GpsDateTime
				influxCount++
			}

		}

		csvCount++
	}

	// Ensures background processes finishes
	client.Close()

	fmt.Println(time.Since(now), csvCount, influxCount)
}

func readFromCSV(file *os.File, c chan Datalog) {
	gocsv.SetCSVReader(func(r io.Reader) gocsv.CSVReader {
		reader := csv.NewReader(r)
		reader.LazyQuotes = true
		reader.FieldsPerRecord = -1
		return reader
	})

	// Read the CSV file into a slice of Record structs
	go func() {
		err := gocsv.UnmarshalToChan(file, c)
		if err != nil {
			panic(err)
		}
	}()
}

func ConnectToInfluxDB() (influxdb2.Client, error) {

	dbToken := os.Getenv("INFLUXDB_TOKEN")
	if dbToken == "" {
		return nil, errors.New("INFLUXDB_TOKEN must be set")
	}

	dbURL := os.Getenv("INFLUXDB_URL")
	if dbURL == "" {
		return nil, errors.New("INFLUXDB_URL must be set")
	}

	// client := influxdb2.NewClient(dbURL, dbToken)
	client := influxdb2.NewClientWithOptions(dbURL, dbToken,
		influxdb2.DefaultOptions().SetBatchSize(20))

	// validate client connection health
	_, err := client.Health(context.Background())

	return client, err
}

func StringToFloat(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Println(err)
	}
	return f
}

func StringToInt(s string) int {
	f, err := strconv.Atoi(s)
	if err != nil {
		log.Println(err)
	}
	return f
}

func dateStringToUnix(s string) time.Time {
	layout := "2006-01-02 15:04:05"
	date, _ := time.Parse(layout, s)
	return date
}
