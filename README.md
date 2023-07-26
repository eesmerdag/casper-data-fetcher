# casper-info-fether

This is a project to fetch some data from Casper blockchain into local database. The idea behind of this is to have the data locally and somehow process it for stats and AI algorithms.
However, this only includes some data and project's skeleton. This is supposed to be improved over time. 

The repo includes 3 different commands to be executed:

### Backfill
this is when you need to run your after you created database. 

### Data Fetcher
A scheduled job to run the periodically and catches the chain by fetching data.

### API
This is to run app for the API endpoints which includes /blocks and /blocks/{height} as of now.