# Dude - Dynon user data explorer

> Each SkyView HDX display can act as a Primary Flight Display (PFD) with Synthetic Vision, an
Engine Monitoring System (EMS), and a Moving Map in a variety of customizable screen
layouts. Data is sourced from various connected modules and devices.
>
> SkyView HDX displays record and store flight information in several datalogs which can be
exported for analysis by the owner, and a high-resolution datalog which can be used by Dynon
for troubleshooting. 
>
> This tool provide an easy and efficient way to enlighten the datalogs for:  
>  - Providing long term flight history,
>  - Providing capability to display flight map and related parameters for improving pilote usage and safety,
>  - Providing accurate informations like what is my average speed to considere preparing flight, what is my average landing speed,
>  - ...

## What are the dependencies
> This tool use:
> - Docker for containers and network management
> - Influx database for timed long term storage
> - Gafana for data presentation
> - Linux shell using a go program for onboarding Dynon datalogs into database.
>
> All of the selected are open source and free for personal usage.

## What is provided
> A local web interface to query flight data shuch as:
> ![Screenshot of web interface.](https://github.com/lla4u/Dude-Influx-Grafana/blob/main/Screenshots/Screenshot_web_Interface.png)

# Installation procedure

## 1 - Docker install (if not yet done)
Depending your operating system download and install docker software.
https://docs.docker.com/get-docker/


## 2 - Building Dude stack
```
1. Where to sit the stack:
   Open terminal, cmd and create your home install directory 
   cd /home/lla 
   mkdir dude 
   then move into: 
   cd dude

2. Clone de github
   git clone https://github.com/lla4u/Dude-Influx-Grafana.git
   or
   Download and unzip zip archive downloaded from github.

3. change directory to Dude-Influx-Grafana
   cd Dude-Influx-Grafana

4. Build the Docker stack using: 
   docker-compose --env-file config.env up --build -d 
   or 
   docker compose --env-file config.env up --build -d  (for recent docker version)

   after a while (mostly depending your network bandwith) 3 containers will be created and available.
5. Check:
   execute docker ps from the terminal

   Having 3 conainers running you are good to go further ...
```

# Adding HDX datalogs to the solution
> Adding datalogs in the tool is in two steps:
> - First, Collect datalog from the HDX
> - Second, Import the datalogs into Influxdb using dude-cli container.

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
> Dynon datalog storage is limited and rewrited over time. So collect datalogs around every 8 flight hours or accept to loose information.

## Importing datalog into dude stack
> 1. Copy the usb key csv file(s) (USER_DATA_LOG.csv) into the datalogs directory
> 2. From teminal or cmd or powershell (windows) execute: 
     docker exec -it dude-cli /bin/ash
> 
> 3. Execute:
>   ./dude-cli -file <path_to_file_you_want_to_load_into_the_stack>
>   Optional verbose mode:
>   ./dude-cli -file <path_to_file_you_want_to_load_into_the_stack> -verbose
> ![Screenshot of dude-cli.](https://github.com/lla4u/Dude-Influx-Grafana/blob/main/Screenshots/Screenshot_dude-cli.png)
>
> Onboarding datalogs required 26.54 seconds (mostly due to poor Influxdb synchrone writes)  
> Submited file is having 84347 csv rows  
> Import saved 19350 rows in database.  


# Ongoing
- [x] Fix InfluxDB write error having huge file count in import. (Fixed using synchrone write and go instead nodejs)
- [x] Make Grafana Dashboard(s) & Data source configuration automatically imported 
- [ ] Use Grafana variable to help finding the flights saved into the InfluxDB.
- [ ] Improve dude-cli performance and UI.
- [ ] Create document for helping users to use the Grafana UI and look at the Dynon datalogs.

Have a safe flights.  
Laurent



