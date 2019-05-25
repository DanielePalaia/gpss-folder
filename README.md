## Introduction
This software is part of the set of GPSS connectors developed by the Pivotal PDE EMEA team to test Pivotal Greenplum on several streaming scenarios. It continues the prototypes started here (for reference and better understanding): </br></br>
https://github.com/DanielePalaia/gpss-rabbit-greenplum-connector</br>
https://github.com/DanielePalaia/gpss-pipe</br>
https://github.com/DanielePalaia/gpss-splunk</br>
https://github.com/DanielePalaia/splunkExternalTables</br>
</br> </br>
In this case the connector will listen to a specific directory for new files in .json and .csv format and it will then ingest Greenplum using gpss and gpscli (which is using Kafka as default broker)</br></br>
https://gpdb.docs.pivotal.io/5160/greenplum-stream/loading-gpss.html</br>
https://gpdb.docs.pivotal.io/5160/greenplum-stream/ref/gpsscli.html</br></br>

## Summary
The connector is highly multithreaded, using all the potentiality of Go routines </br></br>
https://golangbot.com/goroutines/ </br>
This translate to be able to load several files at the same time with very minimal overhead </br></br>
Once started on a given directory the software will create inside it three subdirectories: </br></br>
**staging** </br>
**processing** </br>
**completed** </br></br>
Files are expected to be provided inside staging (even several in parallel) </br>
Once landed the software will copy the files inside processing </br>
it will then read line by line and send the result to a kafka topic (based on the extension .json or .csv) </br>
After completation it will move the files from processing to completed folder </br>
The following checks are provided:
1. If a file is submitted and it's already in completed the two files are compared. If different a new ingestion process will start otherwise if equals it will not happen again </br>
2. If a file is been processed and at the same time a file with the same name is provided in staging then the file is skipped (you have to resubmit it when the original one is completed and moved to completed folder) </br>
3. If the connector is stopped or something bad happens and files are in the middle to be processed, files remain in processing folder, so you need to do an investigation on what was ingested and what is not. At the moment the software is not able to recontinue with the ingestion if interrupted </br>
4. Checking for .json and .csv correctess is delegated to gpsscli which will skip entries eventually not valid and log them in gpsscli logs.

## Prerequisites:
### 1. Start a kafka broker and create 2 new topics
This software is using gpsscli to ingest Greenplum tables which is using Kafka as default broker. Please start a kafka broker and create two topics, one will contains .json entries and the other one .csv entries

### 2. Create a Greenplum database, and create two sample tables to be ingested
Create 2 tables to be ingested based on what you have in your input data. </br>
In this sample we will use these two ones:</br></br>
**create table personJson(id int, name text, surname text, email text, address text);** </br>
**create table personCsv(id int, name text, surname text, email text, address text);** </br></br>
Take also in mind that gpsscli will allow you to do some filtering as we will see afterwards </br>

### 3. Enable GPSS and create two gpsscli job
Enable GPSS for the database which contains the two tables, also start a gpss instance and two gpsscli jobs with the configuration specified here: (Please modified it based on your own parameters) </br>

**For .json** </br>
DATABASE: dashboard</br>
USER: gpadmin</br>
HOST: localhost</br>
PORT: 5432</br>
KAFKA:</br>
   INPUT:</br>
     SOURCE:</br>
        BROKERS: 172.16.125.1:9092</br>
        TOPIC: foldernewjson2</br>
     COLUMNS:</br>
        - NAME: jdata</br>
          TYPE: json</br>
     FORMAT: json</br>
     ERROR_LIMIT: 10</br>
   OUTPUT:</br>
     TABLE: person</br>
     MAPPING:</br>
        - NAME: id</br>
          EXPRESSION: (jdata->>'id')::int</br>
        - NAME: name</br>
          EXPRESSION: (jdata->>'name')::text</br>
        - NAME: surname</br>
          EXPRESSION: (jdata->>'surname')::text</br>
        - NAME: email</br>
          EXPRESSION: (jdata->>'email')::text</br>
        - NAME: address</br>
          EXPRESSION: (jdata->>'address')::text</br>

   COMMIT:</br>
     MAX_ROW: 1000</br></br>
     
     
**For .csv** </br>
DATABASE: dashboard</br>
USER: gpadmin</br>
HOST: localhost</br>
PORT: 5432</br>
KAFKA:</br>
   INPUT:</br>
     SOURCE:</br>
        BROKERS: 172.16.125.1:9092</br>
        TOPIC: foldernewcsv5</br>
     COLUMNS:</br>
        - NAME: id</br>
          TYPE: int</br>
        - NAME: name</br>
          TYPE: text</br>
        - NAME: surname</br>
          TYPE: text</br>
        - NAME: email</br>
          TYPE: text</br>
        - NAME: address</br>
          TYPE: text</br>
     FORMAT: csv</br>
     ERROR_LIMIT: 10</br>
   OUTPUT:</br>
     TABLE: personcsv</br>

   COMMIT:</br>
     MAX_ROW: 1000</br></br>
    
## Running the software:

### 1. Configuration file: </br>  
The software is written in GO, binaries are already provided in:</br>  
./bin/macosx/gpss-folder</br>  
./bin/linux/gpss-folder</br>  
There is a configuration file to be put at the same folder of the binary: properties.ini with these configurations to specify</br></br>  
**folder=/Users/dpalaia/foldertest**</br>  
**kafkaIp=localhost:9092**</br>  
**topicJson=foldernewjson2**</br>  
**topicCsv=foldernewcsv5**</br>  </br> 
Please specify folder, the ip where the kafka broder has been started and the two topics created in the previous step </br>

### 2. Simply run the binary: </br> 
./gpss-folder</br> 

### 3. the program is on hold waiting for files, have look to the folder specified and see that staging, processing and completed folders are created

### 4. Provide files in staging
Provide the sample files provided in ./test in the staging directories (even in parallel)

### 5. Some logs entries will be diplayed in the consolle (NB to remove them for performance in the future)

### 6. Stop the 2 gpsscli jobs when you want to finalize the writing (software will continue to send to kafka broker).

### 6. Do some experiment on the files, maybe change tables or gpsscli mapping.


