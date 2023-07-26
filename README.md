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


Database Tables: 
```
-- Table: public.blocks

-- DROP TABLE IF EXISTS public.blocks;

CREATE TABLE IF NOT EXISTS public.blocks
(
"timestamp" timestamp without time zone NOT NULL,
era_id integer NOT NULL,
height integer NOT NULL,
hash text COLLATE pg_catalog."default" NOT NULL,
CONSTRAINT blocks_pkey PRIMARY KEY (height)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.blocks
OWNER to postgres;
-- Index: idx_hash_blocks

-- DROP INDEX IF EXISTS public.idx_hash_blocks;

CREATE UNIQUE INDEX IF NOT EXISTS idx_hash_blocks
ON public.blocks USING btree
(hash COLLATE pg_catalog."default" ASC NULLS LAST)
TABLESPACE pg_default;


-- Table: public.transfers

-- DROP TABLE IF EXISTS public.transfers;

CREATE TABLE IF NOT EXISTS public.transfers
(
    block_hash text COLLATE pg_catalog."default" NOT NULL,
    block_height integer NOT NULL,
    from_account character varying COLLATE pg_catalog."default" NOT NULL,
    to_account character varying COLLATE pg_catalog."default" NOT NULL,
    amount character varying COLLATE pg_catalog."default" NOT NULL,
    gas character varying COLLATE pg_catalog."default" NOT NULL,
    CONSTRAINT transfers_pkey PRIMARY KEY (block_height)
)

TABLESPACE pg_default;

ALTER TABLE IF EXISTS public.transfers
    OWNER to postgres;
-- Index: idx_hash_transfer

-- DROP INDEX IF EXISTS public.idx_hash_transfer;

CREATE UNIQUE INDEX IF NOT EXISTS idx_hash_transfer
    ON public.transfers USING btree
    (block_hash COLLATE pg_catalog."default" ASC NULLS LAST)
    TABLESPACE pg_default;

```