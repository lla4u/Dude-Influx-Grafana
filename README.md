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
>  - Provide capability to display flight map and related parameters for improving pilote usage and safety,
>  - Provide accurate informations like what is my average speed to considere preparing flight, what is my average landing speed,
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
> 3. Navigate to SYSTEM SOFTWARE -> EXPORT USER DATA LOGS 
> 4. (Otional) Define label
> 5. Export pressing button 8
> 
> Video: https://www.youtube.com/watch?v=fS6H_8gNd90&ab_channel=RobertHamilton

> [!IMPORTANT]
> Dynon datalog storage is limited and file is rewrited. So collect around every 8 hours flight or so not to loose information.

## Importing datalog into InfluxDB
> 1. copy the usb key csv file(s) (USER_DATA_LOG.csv) into your datalogs directory
> 2. from teminal or cmd or powershell (windows) execute: docker exec -it dude-cli ./dudeloader.mjs 
> ![Screenshot of cli select datalogs.](https://github.com/lla4u/Dude-Influx-Grafana/blob/main/Screenshots/Screenshot_cli_datalogs_select.png)
> 3. Select the file(s) from the provided list (space bar) and press Enter to import selected file(s).
> ![Screenshot of cli import datalog(s) result.](https://github.com/lla4u/Dude-Influx-Grafana/blob/main/Screenshots/Screenshot_cli_import_result.png)
  
> [!CAUTION]
> Depending your laptop you might have influxDB write issue selecting high count of files.  


# Todos
- [ ] Fix InfluxDB write error having huge file count in import. 
- [x] Make Grafana Dashboard(s) & Data source configuration automatically imported 
- [ ] Use Grafana variable to help finding the flights saved into the InfluxDB.

Have a safe flights.  
Laurent



