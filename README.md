# Dude - Dynon user data explorer

> Each SkyView HDX display can act as a Primary Flight Display (PFD) with Synthetic Vision, an
Engine Monitoring System (EMS), and a Moving Map in a variety of customizable screen
layouts. Data is sourced from various connected modules and devices.
>
> SkyView HDX displays record and store flight information in several datalogs which can be
exported for analysis by the owner, and a high-resolution datalog which can be used by Dynon
for troubleshooting. 
>
> This tooling intend to use on an easy and efficient way the datalogs to:  
>  - Provide long term flight history,
>  - Provide capability to display flight map and related parameters for improving pilote usage,
>  - ...

## What are the dependencies
> - InfluxDB for timed long term storage
> - Gafana for data presentation
> - Nodejs for onboarding Dynon datalogs into InfluxDB.
>
> All of them are open source and free for personal usage.

## What is provided
> A local web interface to query flight data shuch as:
![Screenshot of web interface.](https://github.com/lla4u/Dude-Influx-Grafana/blob/main/Screenshots/Screenshot_web_Interface.png)


# Installation procedure

## 1 - Docker install (if not yet done)
Depending your operating system download and install docker software.
https://docs.docker.com/get-docker/


## 2 - Dude
```
1. go to your laptop install root directory (cd /home/lla/servers/ ) then execute:

2. git clone https://github.com/lla4u/Dude-Influx-Grafana.git
   or
   Download and unzip zip archive downloaded from github.

3. change directory to Dude-Influx-Grafana
   cd Dude-Influx-Grafana

4. Adapt the docker-compose.yml file to fit with your source datalogs directory ( where you plan  to store the Dynon files ).  

   IE: adapt the cli  
       volumes:  
         - /Users/lla/Documents/Laurent/Aviation/P300 Dude:/home/node/sessions  
       to 
         - <Your datalogs directory path>:/home/node/sessions

5. docker-compose up -d --build
   or 
   docker compose up -d --build (for recent docker version)

   after a while (depending your network speed) 3 containers will be created and available.
```

# Adding HDX datalogs to the solution
> Adding datalogs in the solution is made in two steps:
> - Collect datalog from the HDX
> - Import the datalogs into Influxdb using dude-cli container.

## Collecting datalog from the HDX
> Collecting datalog from the HDX is quite trivial and require usb key plugged into the Dynon:
> ( I use the same usb key that for the plates and map updates ...)
> 1. Fire up your HDX
> 2. Press button 7 & 8 Simultaneously for few seconds to startup the dynon setup screen
> 3. Navigate to ...

## Importing datalog into InfluxDB
> 1. copy the USER_DATA_LOG.csv file(s) into your datalogs directory
> 2. from teminal or cmd or powershell (windows) execute docker run -it dude-cli bash
> 3. Select the file(s) from the provided list (space bar) and press Enter to import ...
>
> Depending your laptop you might have issue selecting high count of files (Influxdb write error)
> I am on it for fixing. Any help are welcome.

# Todos
> - Make Grafana Dashboard(s) & Data source configuration automatically imported
> - Use Grafana variable to help finding the flights saved into the InfluxDB.

Have a safe flights.
Laurent



