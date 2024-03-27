-- all credit to https://gist.github.com/jarnaldich/d5952a134d89dfac48d034ed141e86c5?permalink_comment_id=4401600#gistcomment-4401600
--
-- Usage:
--
-- SELECT * FROM v_introspect_schemas;
-- SELECT * FROM v_introspect_tables;
-- SELECT * FROM v_introspect_columns;
-- SELECT * FROM v_introspect_table_oids;
-- SELECT * FROM v_introspect_tables WHERE schema = 'public' AND relkind IN ('r', 'v');
--
-- v_introspect_oids
-- here we get all tables names then all columns then group by table then schema
DROP VIEW IF EXISTS v_introspect_schemas;

DROP VIEW IF EXISTS v_introspect_tables;

DROP VIEW IF EXISTS v_introspect_columns;

DROP VIEW IF EXISTS v_introspect_table_oids;

CREATE VIEW
    v_introspect_table_oids AS (
        WITH
            table_oids AS (
                SELECT
                    c.relname,
                    c.oid,
                    nspname,
                    reltuples,
                    relkind,
                    relam,
                    nspacl,
                    relacl,
                    reltype,
                    relowner,
                    relhasindex
                FROM
                    pg_catalog.pg_class c
                    LEFT JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
                WHERE
                    pg_catalog.PG_TABLE_IS_VISIBLE (c.oid)
            )
        SELECT
            *
        FROM
            table_oids
        ORDER BY
            nspname,
            relname
    );

COMMENT ON VIEW v_introspect_table_oids IS 'This table contains all objects in the pg catalog and visible namespace/schemas (probably public)';

-- v_introspect_columns
CREATE VIEW
    v_introspect_columns AS (
        WITH
            -- we'll use this CTE to track foreign keys between tables
            foreign_keys AS (
                SELECT DISTINCT
                    tc.table_schema AS schema_name,
                    kcu.table_name,
                    kcu.column_name,
                    ccu.table_name AS foreign_table_name,
                    ccu.column_name AS foreign_column_name
                FROM
                    information_schema.table_constraints AS tc
                    JOIN information_schema.key_column_usage AS kcu ON tc.constraint_name = kcu.constraint_name
                    AND tc.table_schema = kcu.table_schema
                    JOIN information_schema.constraint_column_usage AS ccu ON ccu.constraint_name = tc.constraint_name
                    AND ccu.table_schema = tc.table_schema
                WHERE
                    tc.constraint_type = 'FOREIGN KEY'
            ),
            -- we'll use this CTE to track primary keys
            primary_keys AS (
                SELECT DISTINCT
                    tc.table_schema AS schema_name,
                    kcu.table_name,
                    kcu.column_name
                FROM
                    information_schema.table_constraints AS tc
                    JOIN information_schema.key_column_usage AS kcu ON tc.constraint_name = kcu.constraint_name
                    AND tc.table_schema = kcu.table_schema
                WHERE
                    tc.constraint_type = 'PRIMARY KEY'
            )
        SELECT
            toids.relname AS "tablename",
            a.attname AS "column",
            pg_catalog.FORMAT_TYPE (a.atttypid, a.atttypmod) AS "datatype",
            nspname AS "schema",
            oid,
            ROW_NUMBER() OVER (
                ORDER BY
                    nspname DESC,
                    "relname",
                    "attnum"
            ) AS column_id,
            -- toids.nspacl AS "nspacl",
            -- toids.relacl AS "tableacl",
            attnum AS "pos",
            atttypid AS "typeid",
            attlen AS "typelen",
            atttypmod AS "typemod",
            attndims AS "ndims",
            attnotnull AS "notnull",
            atthasdef AS "hasdefault",
            atthasmissing AS "hasmissing",
            primary_keys.column_name IS NOT NULL
            AND foreign_keys.foreign_column_name IS NULL AS "ispkey",
            foreign_keys.foreign_table_name AS "ftable",
            foreign_keys.foreign_column_name AS "fcolumn"
        FROM
            pg_catalog.pg_attribute a,
            v_introspect_table_oids toids
            -- mix in knowledge we have about primary keys
            LEFT JOIN LATERAL (
                SELECT
                    *
                FROM
                    primary_keys
                WHERE
                    primary_keys.schema_name = nspname
                    AND primary_keys.table_name = toids.relname
                    AND primary_keys.column_name = a.attname
            ) AS primary_keys ON true
            -- optionally mix in any knowledge we have about foreign keys
            LEFT JOIN LATERAL (
                SELECT
                    *
                FROM
                    foreign_keys
                WHERE
                    foreign_keys.schema_name = nspname
                    AND foreign_keys.table_name = toids.relname
                    AND foreign_keys.column_name = a.attname
            ) AS foreign_keys ON true
        WHERE
            a.attnum > 0
            AND NOT a.attisdropped
            AND a.attrelid = toids.oid
        ORDER BY
            schema DESC,
            tablename,
            pos
    );

COMMENT ON VIEW v_introspect_columns IS 'This view contains a row for every single column in pg_catalog and public (or query pool schema). I have added some nice names but left all pg data for nonlimiting aplication of end API. todo remove useless fields and simplify';

-- v_introspect_tables
CREATE VIEW
    v_introspect_tables AS (
        SELECT
            tablename, -- here we could add other fields to our column definition
            col.oid,
            schema, -- add some table info back in
            reltuples,
            relkind,
            relam,
            -- nspacl,
            relacl,
            reltype,
            relowner,
            relhasindex,
            JSON_AGG(
                JSON_BUILD_OBJECT( -- condense columns into object
                    'column',
                    "column",
                    'datatype',
                    datatype,
                    'table',
                    tablename,
                    'pos',
                    pos,
                    'typeid',
                    typeid,
                    'typelen',
                    typelen,
                    'typemod',
                    "typemod",
                    'notnull',
                    "notnull",
                    'hasdefault',
                    hasdefault,
                    'hasmissing',
                    hasmissing,
                    'ispkey',
                    ispkey,
                    'ftable',
                    ftable,
                    'fcolumn',
                    fcolumn,
                    'parent_id',
                    col.oid
                )
            ) AS "columns"
        FROM
            v_introspect_columns col,
            v_introspect_table_oids ts
        WHERE
            ts.oid = col.oid
        GROUP BY
            tablename,
            col.oid,
            "schema",
            reltuples,
            relkind,
            relam,
            nspacl,
            relacl,
            reltype,
            relowner,
            relhasindex
        ORDER BY
            "schema" DESC,
            tablename,
            relkind
    );

COMMENT ON VIEW v_introspect_tables IS 'Ok now we have all the columns and can group them by the table. json aggregates the columns ';

-- v_introspect_schemas
CREATE VIEW
    v_introspect_schemas AS ( -- lets bundle it all up at the top level, we could ad more info here maybe
        SELECT
            "schema",
            ROW_NUMBER() OVER (
                ORDER BY
                    "schema"
            ) AS id,
            JSON_AGG(
                JSON_BUILD_OBJECT(
                    'id',
                    oid,
                    'tablename',
                    tablename,
                    'columns',
                    "columns",
                    'schema',
                    schema
                )
            ) AS "tables"
        FROM
            v_introspect_tables
        GROUP BY
            "schema"
        ORDER BY
            "schema"
    );

COMMENT ON VIEW v_introspect_schemas IS 'just another aggregation to group by schema, if for example, we dont care about public ';
