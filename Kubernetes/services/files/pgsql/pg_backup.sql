--
-- PostgreSQL database cluster dump
--

SET default_transaction_read_only = off;

SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;

--
-- Roles
--

CREATE ROLE aman;
ALTER ROLE aman WITH NOSUPERUSER INHERIT NOCREATEROLE CREATEDB LOGIN NOREPLICATION NOBYPASSRLS PASSWORD 'SCRAM-SHA-256$4096:Sw5wp2EkykO/Hv5WH3YuIg==$zeQp0iY71LwI+qmbBCPNAEE8mN8q/Lxrnp1OJw+av64=:BHR8Z1+YcHeDEb8gs2bIt400/Uh7MmETMGo+5WfF3GQ=';
CREATE ROLE postgres;
ALTER ROLE postgres WITH SUPERUSER INHERIT CREATEROLE CREATEDB LOGIN REPLICATION BYPASSRLS PASSWORD 'SCRAM-SHA-256$4096:BEMJDV9ZyF6RMEV36k4+yA==$ezAHL5JsBJVovnHTGi5PgtjjDtfEp0jkBSVnQ/GMfVc=:RNfeW7CZFEiFMIeQ2fsKXwUjlkksnd8fViu89c2gGrA=';
CREATE ROLE repl_user;
ALTER ROLE repl_user WITH NOSUPERUSER INHERIT NOCREATEROLE NOCREATEDB LOGIN REPLICATION NOBYPASSRLS PASSWORD 'SCRAM-SHA-256$4096:myUC6E9qJT6+LNfAP3odTg==$RpcrS3G1TxJtBDonXB8XqjsuNjrp7ISZr8ruZ4GvTwA=:ENExUAnWkBb3+li/66P2eurqs9ldKuRtLewFotn8ahE=';

--
-- User Configurations
--








--
-- Databases
--

--
-- Database "template1" dump
--

\connect template1

--
-- PostgreSQL database dump
--

-- Dumped from database version 16.3
-- Dumped by pg_dump version 16.3

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- PostgreSQL database dump complete
--

--
-- Database "compute" dump
--

--
-- PostgreSQL database dump
--

-- Dumped from database version 16.3
-- Dumped by pg_dump version 16.3

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: compute; Type: DATABASE; Schema: -; Owner: aman
--

CREATE DATABASE compute WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE_PROVIDER = libc LOCALE = 'en_US.UTF-8';


ALTER DATABASE compute OWNER TO aman;

\connect compute

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: public; Type: SCHEMA; Schema: -; Owner: aman
--

-- *not* creating schema, since initdb creates it


ALTER SCHEMA public OWNER TO aman;

--
-- PostgreSQL database dump complete
--

--
-- Database "metabase" dump
--

--
-- PostgreSQL database dump
--

-- Dumped from database version 16.3
-- Dumped by pg_dump version 16.3

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: metabase; Type: DATABASE; Schema: -; Owner: postgres
--

CREATE DATABASE metabase WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE_PROVIDER = libc LOCALE = 'en_US.UTF-8';


ALTER DATABASE metabase OWNER TO postgres;

\connect metabase

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: citext; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS citext WITH SCHEMA public;


--
-- Name: EXTENSION citext; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION citext IS 'data type for case-insensitive character strings';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: action; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.action (
    id integer NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    type text NOT NULL,
    model_id integer NOT NULL,
    name character varying(254) NOT NULL,
    description text,
    parameters text,
    parameter_mappings text,
    visualization_settings text,
    public_uuid character(36),
    made_public_by_id integer,
    creator_id integer,
    archived boolean DEFAULT false NOT NULL,
    entity_id character(21)
);


ALTER TABLE public.action OWNER TO aman;

--
-- Name: TABLE action; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON TABLE public.action IS 'An action is something you can do, such as run a readwrite query';


--
-- Name: COLUMN action.created_at; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.action.created_at IS 'The timestamp of when the action was created';


--
-- Name: COLUMN action.updated_at; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.action.updated_at IS 'The timestamp of when the action was updated';


--
-- Name: COLUMN action.type; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.action.type IS 'Type of action';


--
-- Name: COLUMN action.model_id; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.action.model_id IS 'The associated model';


--
-- Name: COLUMN action.name; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.action.name IS 'The name of the action';


--
-- Name: COLUMN action.description; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.action.description IS 'The description of the action';


--
-- Name: COLUMN action.parameters; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.action.parameters IS 'The saved parameters for this action';


--
-- Name: COLUMN action.parameter_mappings; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.action.parameter_mappings IS 'The saved parameter mappings for this action';


--
-- Name: COLUMN action.visualization_settings; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.action.visualization_settings IS 'The UI visualization_settings for this action';


--
-- Name: COLUMN action.public_uuid; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.action.public_uuid IS 'Unique UUID used to in publically-accessible links to this Action.';


--
-- Name: COLUMN action.made_public_by_id; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.action.made_public_by_id IS 'The ID of the User who first publically shared this Action.';


--
-- Name: COLUMN action.creator_id; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.action.creator_id IS 'The user who created the action';


--
-- Name: COLUMN action.archived; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.action.archived IS 'Whether or not the action has been archived';


--
-- Name: COLUMN action.entity_id; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.action.entity_id IS 'Random NanoID tag for unique identity.';


--
-- Name: action_id_seq; Type: SEQUENCE; Schema: public; Owner: aman
--

ALTER TABLE public.action ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.action_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: api_key; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.api_key (
    id integer NOT NULL,
    user_id integer,
    key character varying(254) NOT NULL,
    key_prefix character varying(7) NOT NULL,
    creator_id integer NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    name character varying(254) NOT NULL,
    updated_by_id integer NOT NULL,
    scope character varying(64)
);


ALTER TABLE public.api_key OWNER TO aman;

--
-- Name: TABLE api_key; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON TABLE public.api_key IS 'An API Key';


--
-- Name: COLUMN api_key.id; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.api_key.id IS 'The ID of the API Key itself';


--
-- Name: COLUMN api_key.user_id; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.api_key.user_id IS 'The ID of the user who this API Key acts as';


--
-- Name: COLUMN api_key.key; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.api_key.key IS 'The hashed API key';


--
-- Name: COLUMN api_key.key_prefix; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.api_key.key_prefix IS 'The first 7 characters of the unhashed key';


--
-- Name: COLUMN api_key.creator_id; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.api_key.creator_id IS 'The ID of the user that created this API key';


--
-- Name: COLUMN api_key.created_at; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.api_key.created_at IS 'The timestamp when the key was created';


--
-- Name: COLUMN api_key.updated_at; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.api_key.updated_at IS 'The timestamp when the key was last updated';


--
-- Name: COLUMN api_key.name; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.api_key.name IS 'The user-defined name of the API key.';


--
-- Name: COLUMN api_key.updated_by_id; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.api_key.updated_by_id IS 'The ID of the user that last updated this API key';


--
-- Name: COLUMN api_key.scope; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.api_key.scope IS 'The scope of the API key, if applicable';


--
-- Name: api_key_id_seq; Type: SEQUENCE; Schema: public; Owner: aman
--

ALTER TABLE public.api_key ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.api_key_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: application_permissions_revision; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.application_permissions_revision (
    id integer NOT NULL,
    before text NOT NULL,
    after text NOT NULL,
    user_id integer NOT NULL,
    created_at timestamp with time zone NOT NULL,
    remark text
);


ALTER TABLE public.application_permissions_revision OWNER TO aman;

--
-- Name: application_permissions_revision_id_seq; Type: SEQUENCE; Schema: public; Owner: aman
--

ALTER TABLE public.application_permissions_revision ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.application_permissions_revision_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: audit_log; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.audit_log (
    id integer NOT NULL,
    topic character varying(32) NOT NULL,
    "timestamp" timestamp with time zone NOT NULL,
    end_timestamp timestamp with time zone,
    user_id integer,
    model character varying(32),
    model_id integer,
    details text NOT NULL
);


ALTER TABLE public.audit_log OWNER TO aman;

--
-- Name: TABLE audit_log; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON TABLE public.audit_log IS 'Used to store application events for auditing use cases';


--
-- Name: COLUMN audit_log.topic; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.audit_log.topic IS 'The topic of a given audit event';


--
-- Name: COLUMN audit_log."timestamp"; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.audit_log."timestamp" IS 'The time an event was recorded';


--
-- Name: COLUMN audit_log.end_timestamp; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.audit_log.end_timestamp IS 'The time an event ended, if applicable';


--
-- Name: COLUMN audit_log.user_id; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.audit_log.user_id IS 'The user who performed an action or triggered an event';


--
-- Name: COLUMN audit_log.model; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.audit_log.model IS 'The name of the model this event applies to (e.g. Card, Dashboard), if applicable';


--
-- Name: COLUMN audit_log.model_id; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.audit_log.model_id IS 'The ID of the model this event applies to, if applicable';


--
-- Name: COLUMN audit_log.details; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.audit_log.details IS 'A JSON map with metadata about the event';


--
-- Name: audit_log_id_seq; Type: SEQUENCE; Schema: public; Owner: aman
--

ALTER TABLE public.audit_log ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.audit_log_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: bookmark_ordering; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.bookmark_ordering (
    id integer NOT NULL,
    user_id integer NOT NULL,
    type character varying(255) NOT NULL,
    item_id integer NOT NULL,
    ordering integer NOT NULL
);


ALTER TABLE public.bookmark_ordering OWNER TO aman;

--
-- Name: bookmark_ordering_id_seq; Type: SEQUENCE; Schema: public; Owner: aman
--

ALTER TABLE public.bookmark_ordering ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.bookmark_ordering_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: cache_config; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.cache_config (
    id integer NOT NULL,
    model character varying(32) NOT NULL,
    model_id integer NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    strategy text NOT NULL,
    config text NOT NULL,
    state text,
    invalidated_at timestamp with time zone,
    next_run_at timestamp with time zone
);


ALTER TABLE public.cache_config OWNER TO aman;

--
-- Name: TABLE cache_config; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON TABLE public.cache_config IS 'Cache Configuration';


--
-- Name: COLUMN cache_config.id; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.cache_config.id IS 'Unique ID';


--
-- Name: COLUMN cache_config.model; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.cache_config.model IS 'Name of an entity model';


--
-- Name: COLUMN cache_config.model_id; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.cache_config.model_id IS 'ID of the said entity';


--
-- Name: COLUMN cache_config.created_at; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.cache_config.created_at IS 'Timestamp when the config was inserted';


--
-- Name: COLUMN cache_config.updated_at; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.cache_config.updated_at IS 'Timestamp when the config was updated';


--
-- Name: COLUMN cache_config.strategy; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.cache_config.strategy IS 'caching strategy name';


--
-- Name: COLUMN cache_config.config; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.cache_config.config IS 'caching strategy configuration';


--
-- Name: COLUMN cache_config.state; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.cache_config.state IS 'state for strategies needing to keep some data between runs';


--
-- Name: COLUMN cache_config.invalidated_at; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.cache_config.invalidated_at IS 'indicates when a cache was invalidated last time for schedule-based strategies';


--
-- Name: COLUMN cache_config.next_run_at; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.cache_config.next_run_at IS 'keeps next time to run for schedule-based strategies';


--
-- Name: cache_config_id_seq; Type: SEQUENCE; Schema: public; Owner: aman
--

ALTER TABLE public.cache_config ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.cache_config_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: card_bookmark; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.card_bookmark (
    id integer NOT NULL,
    user_id integer NOT NULL,
    card_id integer NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.card_bookmark OWNER TO aman;

--
-- Name: card_bookmark_id_seq; Type: SEQUENCE; Schema: public; Owner: aman
--

ALTER TABLE public.card_bookmark ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.card_bookmark_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: card_label; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.card_label (
    id integer NOT NULL,
    card_id integer NOT NULL,
    label_id integer NOT NULL
);


ALTER TABLE public.card_label OWNER TO aman;

--
-- Name: card_label_id_seq; Type: SEQUENCE; Schema: public; Owner: aman
--

ALTER TABLE public.card_label ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.card_label_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: cloud_migration; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.cloud_migration (
    id integer NOT NULL,
    external_id text NOT NULL,
    upload_url text NOT NULL,
    state character varying(32) DEFAULT 'init'::character varying NOT NULL,
    progress integer DEFAULT 0 NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL
);


ALTER TABLE public.cloud_migration OWNER TO aman;

--
-- Name: TABLE cloud_migration; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON TABLE public.cloud_migration IS 'Migrate to cloud directly from Metabase';


--
-- Name: COLUMN cloud_migration.id; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.cloud_migration.id IS 'Unique ID';


--
-- Name: COLUMN cloud_migration.external_id; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.cloud_migration.external_id IS 'Matching ID in Cloud for this migration';


--
-- Name: COLUMN cloud_migration.upload_url; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.cloud_migration.upload_url IS 'URL where the backup will be uploaded to';


--
-- Name: COLUMN cloud_migration.state; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.cloud_migration.state IS 'Current state of the migration: init, setup, dump, upload, done, error, cancelled';


--
-- Name: COLUMN cloud_migration.progress; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.cloud_migration.progress IS 'Number between 0 to 100 representing progress as a percentage';


--
-- Name: COLUMN cloud_migration.created_at; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.cloud_migration.created_at IS 'Timestamp when the config was inserted';


--
-- Name: COLUMN cloud_migration.updated_at; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.cloud_migration.updated_at IS 'Timestamp when the config was updated';


--
-- Name: cloud_migration_id_seq; Type: SEQUENCE; Schema: public; Owner: aman
--

ALTER TABLE public.cloud_migration ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.cloud_migration_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: collection; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.collection (
    id integer NOT NULL,
    name text NOT NULL,
    description text,
    archived boolean DEFAULT false NOT NULL,
    location character varying(254) DEFAULT '/'::character varying NOT NULL,
    personal_owner_id integer,
    slug character varying(510) NOT NULL,
    namespace character varying(254),
    authority_level character varying(255),
    entity_id character(21),
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    type character varying(256),
    is_sample boolean DEFAULT false NOT NULL
);


ALTER TABLE public.collection OWNER TO aman;

--
-- Name: COLUMN collection.created_at; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.collection.created_at IS 'Timestamp of when this Collection was created.';


--
-- Name: COLUMN collection.type; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.collection.type IS 'This is used to differentiate instance-analytics collections from all other collections.';


--
-- Name: COLUMN collection.is_sample; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.collection.is_sample IS 'Is the collection part of the sample content?';


--
-- Name: collection_bookmark; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.collection_bookmark (
    id integer NOT NULL,
    user_id integer NOT NULL,
    collection_id integer NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.collection_bookmark OWNER TO aman;

--
-- Name: collection_bookmark_id_seq; Type: SEQUENCE; Schema: public; Owner: aman
--

ALTER TABLE public.collection_bookmark ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.collection_bookmark_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: collection_id_seq; Type: SEQUENCE; Schema: public; Owner: aman
--

ALTER TABLE public.collection ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.collection_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: collection_permission_graph_revision; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.collection_permission_graph_revision (
    id integer NOT NULL,
    before text NOT NULL,
    after text NOT NULL,
    user_id integer NOT NULL,
    created_at timestamp with time zone NOT NULL,
    remark text
);


ALTER TABLE public.collection_permission_graph_revision OWNER TO aman;

--
-- Name: collection_permission_graph_revision_id_seq; Type: SEQUENCE; Schema: public; Owner: aman
--

ALTER TABLE public.collection_permission_graph_revision ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.collection_permission_graph_revision_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: connection_impersonations; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.connection_impersonations (
    id integer NOT NULL,
    db_id integer NOT NULL,
    group_id integer NOT NULL,
    attribute text
);


ALTER TABLE public.connection_impersonations OWNER TO aman;

--
-- Name: TABLE connection_impersonations; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON TABLE public.connection_impersonations IS 'Table for holding connection impersonation policies';


--
-- Name: COLUMN connection_impersonations.db_id; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.connection_impersonations.db_id IS 'ID of the database this connection impersonation policy affects';


--
-- Name: COLUMN connection_impersonations.group_id; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.connection_impersonations.group_id IS 'ID of the permissions group this connection impersonation policy affects';


--
-- Name: COLUMN connection_impersonations.attribute; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.connection_impersonations.attribute IS 'User attribute associated with the database role to use for this connection impersonation policy';


--
-- Name: connection_impersonations_id_seq; Type: SEQUENCE; Schema: public; Owner: aman
--

ALTER TABLE public.connection_impersonations ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.connection_impersonations_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: core_session; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.core_session (
    id character varying(254) NOT NULL,
    user_id integer NOT NULL,
    created_at timestamp with time zone NOT NULL,
    anti_csrf_token text
);


ALTER TABLE public.core_session OWNER TO aman;

--
-- Name: core_user; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.core_user (
    id integer NOT NULL,
    email public.citext NOT NULL,
    first_name character varying(254),
    last_name character varying(254),
    password character varying(254),
    password_salt character varying(254) DEFAULT 'default'::character varying,
    date_joined timestamp with time zone NOT NULL,
    last_login timestamp with time zone,
    is_superuser boolean DEFAULT false NOT NULL,
    is_active boolean DEFAULT true NOT NULL,
    reset_token character varying(254),
    reset_triggered bigint,
    is_qbnewb boolean DEFAULT true NOT NULL,
    login_attributes text,
    updated_at timestamp with time zone,
    sso_source character varying(254),
    locale character varying(5),
    is_datasetnewb boolean DEFAULT true NOT NULL,
    settings text,
    type character varying(64) DEFAULT 'personal'::character varying NOT NULL,
    entity_id character(21)
);


ALTER TABLE public.core_user OWNER TO aman;

--
-- Name: COLUMN core_user.type; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.core_user.type IS 'The type of user';


--
-- Name: COLUMN core_user.entity_id; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.core_user.entity_id IS 'NanoID tag for each user';


--
-- Name: core_user_id_seq; Type: SEQUENCE; Schema: public; Owner: aman
--

ALTER TABLE public.core_user ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.core_user_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: dashboard_bookmark; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.dashboard_bookmark (
    id integer NOT NULL,
    user_id integer NOT NULL,
    dashboard_id integer NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.dashboard_bookmark OWNER TO aman;

--
-- Name: dashboard_bookmark_id_seq; Type: SEQUENCE; Schema: public; Owner: aman
--

ALTER TABLE public.dashboard_bookmark ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.dashboard_bookmark_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: dashboard_favorite; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.dashboard_favorite (
    id integer NOT NULL,
    user_id integer NOT NULL,
    dashboard_id integer NOT NULL
);


ALTER TABLE public.dashboard_favorite OWNER TO aman;

--
-- Name: dashboard_favorite_id_seq; Type: SEQUENCE; Schema: public; Owner: aman
--

ALTER TABLE public.dashboard_favorite ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.dashboard_favorite_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: dashboard_tab; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.dashboard_tab (
    id integer NOT NULL,
    dashboard_id integer NOT NULL,
    name text NOT NULL,
    "position" integer NOT NULL,
    entity_id character(21),
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.dashboard_tab OWNER TO aman;

--
-- Name: TABLE dashboard_tab; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON TABLE public.dashboard_tab IS 'Join table connecting dashboard to dashboardcards';


--
-- Name: COLUMN dashboard_tab.dashboard_id; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.dashboard_tab.dashboard_id IS 'The dashboard that a tab is on';


--
-- Name: COLUMN dashboard_tab.name; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.dashboard_tab.name IS 'Displayed name of the tab';


--
-- Name: COLUMN dashboard_tab."position"; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.dashboard_tab."position" IS 'Position of the tab with respect to others tabs in dashboard';


--
-- Name: COLUMN dashboard_tab.entity_id; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.dashboard_tab.entity_id IS 'Random NanoID tag for unique identity.';


--
-- Name: COLUMN dashboard_tab.created_at; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.dashboard_tab.created_at IS 'The timestamp at which the tab was created';


--
-- Name: COLUMN dashboard_tab.updated_at; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.dashboard_tab.updated_at IS 'The timestamp at which the tab was last updated';


--
-- Name: dashboard_tab_id_seq; Type: SEQUENCE; Schema: public; Owner: aman
--

ALTER TABLE public.dashboard_tab ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.dashboard_tab_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: dashboardcard_series; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.dashboardcard_series (
    id integer NOT NULL,
    dashboardcard_id integer NOT NULL,
    card_id integer NOT NULL,
    "position" integer NOT NULL
);


ALTER TABLE public.dashboardcard_series OWNER TO aman;

--
-- Name: dashboardcard_series_id_seq; Type: SEQUENCE; Schema: public; Owner: aman
--

ALTER TABLE public.dashboardcard_series ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.dashboardcard_series_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: data_permissions; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.data_permissions (
    id integer NOT NULL,
    group_id integer NOT NULL,
    perm_type character varying(64) NOT NULL,
    db_id integer NOT NULL,
    schema_name character varying(254),
    table_id integer,
    perm_value character varying(64) NOT NULL
);


ALTER TABLE public.data_permissions OWNER TO aman;

--
-- Name: TABLE data_permissions; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON TABLE public.data_permissions IS 'A table to store database and table permissions';


--
-- Name: COLUMN data_permissions.id; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.data_permissions.id IS 'The ID of the permission';


--
-- Name: COLUMN data_permissions.group_id; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.data_permissions.group_id IS 'The ID of the associated permission group';


--
-- Name: COLUMN data_permissions.perm_type; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.data_permissions.perm_type IS 'The type of the permission (e.g. "data", "collection", "download"...)';


--
-- Name: COLUMN data_permissions.db_id; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.data_permissions.db_id IS 'A database ID, for DB and table-level permissions';


--
-- Name: COLUMN data_permissions.schema_name; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.data_permissions.schema_name IS 'A schema name, for table-level permissions';


--
-- Name: COLUMN data_permissions.table_id; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.data_permissions.table_id IS 'A table ID';


--
-- Name: COLUMN data_permissions.perm_value; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.data_permissions.perm_value IS 'The value this permission is set to.';


--
-- Name: data_permissions_id_seq; Type: SEQUENCE; Schema: public; Owner: aman
--

ALTER TABLE public.data_permissions ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.data_permissions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: databasechangelog; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.databasechangelog (
    id character varying(255) NOT NULL,
    author character varying(255) NOT NULL,
    filename character varying(255) NOT NULL,
    dateexecuted timestamp without time zone NOT NULL,
    orderexecuted integer NOT NULL,
    exectype character varying(10) NOT NULL,
    md5sum character varying(35),
    description character varying(255),
    comments character varying(255),
    tag character varying(255),
    liquibase character varying(20),
    contexts character varying(255),
    labels character varying(255),
    deployment_id character varying(10)
);


ALTER TABLE public.databasechangelog OWNER TO aman;

--
-- Name: databasechangeloglock; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.databasechangeloglock (
    id integer NOT NULL,
    locked boolean NOT NULL,
    lockgranted timestamp without time zone,
    lockedby character varying(255)
);


ALTER TABLE public.databasechangeloglock OWNER TO aman;

--
-- Name: dependency; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.dependency (
    id integer NOT NULL,
    model character varying(32) NOT NULL,
    model_id integer NOT NULL,
    dependent_on_model character varying(32) NOT NULL,
    dependent_on_id integer NOT NULL,
    created_at timestamp with time zone NOT NULL
);


ALTER TABLE public.dependency OWNER TO aman;

--
-- Name: dependency_id_seq; Type: SEQUENCE; Schema: public; Owner: aman
--

ALTER TABLE public.dependency ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.dependency_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: dimension; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.dimension (
    id integer NOT NULL,
    field_id integer NOT NULL,
    name character varying(254) NOT NULL,
    type character varying(254) NOT NULL,
    human_readable_field_id integer,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    entity_id character(21)
);


ALTER TABLE public.dimension OWNER TO aman;

--
-- Name: dimension_id_seq; Type: SEQUENCE; Schema: public; Owner: aman
--

ALTER TABLE public.dimension ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.dimension_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: field_usage; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.field_usage (
    id integer NOT NULL,
    field_id integer NOT NULL,
    query_execution_id integer NOT NULL,
    used_in character varying(25) NOT NULL,
    filter_op character varying(25),
    aggregation_function character varying(25),
    breakout_temporal_unit character varying(25),
    breakout_binning_strategy character varying(25),
    breakout_binning_num_bins integer,
    breakout_binning_bin_width integer,
    created_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.field_usage OWNER TO aman;

--
-- Name: TABLE field_usage; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON TABLE public.field_usage IS 'Used to store field usage during query execution';


--
-- Name: COLUMN field_usage.id; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.field_usage.id IS 'Unique ID';


--
-- Name: COLUMN field_usage.field_id; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.field_usage.field_id IS 'ID of the field';


--
-- Name: COLUMN field_usage.query_execution_id; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.field_usage.query_execution_id IS 'referenced query execution';


--
-- Name: COLUMN field_usage.used_in; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.field_usage.used_in IS 'which part of the query the field was used in';


--
-- Name: COLUMN field_usage.filter_op; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.field_usage.filter_op IS 'filter''s operator that applied to the field';


--
-- Name: COLUMN field_usage.aggregation_function; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.field_usage.aggregation_function IS 'the aggregation function that field applied to';


--
-- Name: COLUMN field_usage.breakout_temporal_unit; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.field_usage.breakout_temporal_unit IS 'temporal unit options of the breakout';


--
-- Name: COLUMN field_usage.breakout_binning_strategy; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.field_usage.breakout_binning_strategy IS 'the strategy of breakout';


--
-- Name: COLUMN field_usage.breakout_binning_num_bins; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.field_usage.breakout_binning_num_bins IS 'The numbin option of breakout';


--
-- Name: COLUMN field_usage.breakout_binning_bin_width; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.field_usage.breakout_binning_bin_width IS 'The numbin option of breakout';


--
-- Name: COLUMN field_usage.created_at; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.field_usage.created_at IS 'The time a field usage was recorded';


--
-- Name: field_usage_id_seq; Type: SEQUENCE; Schema: public; Owner: aman
--

ALTER TABLE public.field_usage ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.field_usage_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: sandboxes; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.sandboxes (
    id integer NOT NULL,
    group_id integer NOT NULL,
    table_id integer NOT NULL,
    card_id integer,
    attribute_remappings text,
    permission_id integer
);


ALTER TABLE public.sandboxes OWNER TO aman;

--
-- Name: COLUMN sandboxes.permission_id; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.sandboxes.permission_id IS 'The ID of the corresponding permissions path for this sandbox';


--
-- Name: group_table_access_policy_id_seq; Type: SEQUENCE; Schema: public; Owner: aman
--

ALTER TABLE public.sandboxes ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.group_table_access_policy_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: http_action; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.http_action (
    action_id integer NOT NULL,
    template text NOT NULL,
    response_handle text,
    error_handle text
);


ALTER TABLE public.http_action OWNER TO aman;

--
-- Name: TABLE http_action; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON TABLE public.http_action IS 'An http api call type of action';


--
-- Name: COLUMN http_action.action_id; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.http_action.action_id IS 'The related action';


--
-- Name: COLUMN http_action.template; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.http_action.template IS 'A template that defines method,url,body,headers required to make an api call';


--
-- Name: COLUMN http_action.response_handle; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.http_action.response_handle IS 'A program to take an api response and transform to an appropriate response for emitters';


--
-- Name: COLUMN http_action.error_handle; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.http_action.error_handle IS 'A program to take an api response to determine if an error occurred';


--
-- Name: implicit_action; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.implicit_action (
    action_id integer NOT NULL,
    kind text NOT NULL
);


ALTER TABLE public.implicit_action OWNER TO aman;

--
-- Name: TABLE implicit_action; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON TABLE public.implicit_action IS 'An action with dynamic parameters based on the underlying model';


--
-- Name: COLUMN implicit_action.action_id; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.implicit_action.action_id IS 'The associated action';


--
-- Name: COLUMN implicit_action.kind; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.implicit_action.kind IS 'The kind of implicit action create/update/delete';


--
-- Name: label; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.label (
    id integer NOT NULL,
    name character varying(254) NOT NULL,
    slug character varying(254) NOT NULL,
    icon character varying(128)
);


ALTER TABLE public.label OWNER TO aman;

--
-- Name: label_id_seq; Type: SEQUENCE; Schema: public; Owner: aman
--

ALTER TABLE public.label ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.label_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: login_history; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.login_history (
    id integer NOT NULL,
    "timestamp" timestamp with time zone DEFAULT now() NOT NULL,
    user_id integer NOT NULL,
    session_id character varying(254),
    device_id character(36) NOT NULL,
    device_description text NOT NULL,
    ip_address text NOT NULL
);


ALTER TABLE public.login_history OWNER TO aman;

--
-- Name: login_history_id_seq; Type: SEQUENCE; Schema: public; Owner: aman
--

ALTER TABLE public.login_history ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.login_history_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: metabase_database; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.metabase_database (
    id integer NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    name character varying(254) NOT NULL,
    description text,
    details text NOT NULL,
    engine character varying(254) NOT NULL,
    is_sample boolean DEFAULT false NOT NULL,
    is_full_sync boolean DEFAULT true NOT NULL,
    points_of_interest text,
    caveats text,
    metadata_sync_schedule character varying(254) DEFAULT '0 50 * * * ? *'::character varying NOT NULL,
    cache_field_values_schedule character varying(254) DEFAULT NULL::character varying,
    timezone character varying(254),
    is_on_demand boolean DEFAULT false NOT NULL,
    auto_run_queries boolean DEFAULT true NOT NULL,
    refingerprint boolean,
    cache_ttl integer,
    initial_sync_status character varying(32) DEFAULT 'complete'::character varying NOT NULL,
    creator_id integer,
    settings text,
    dbms_version text,
    is_audit boolean DEFAULT false NOT NULL,
    uploads_enabled boolean DEFAULT false NOT NULL,
    uploads_schema_name text,
    uploads_table_prefix text
);


ALTER TABLE public.metabase_database OWNER TO aman;

--
-- Name: COLUMN metabase_database.dbms_version; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.metabase_database.dbms_version IS 'A JSON object describing the flavor and version of the DBMS.';


--
-- Name: COLUMN metabase_database.is_audit; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.metabase_database.is_audit IS 'Only the app db, visible to admins via auditing should have this set true.';


--
-- Name: COLUMN metabase_database.uploads_enabled; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.metabase_database.uploads_enabled IS 'Whether uploads are enabled for this database';


--
-- Name: COLUMN metabase_database.uploads_schema_name; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.metabase_database.uploads_schema_name IS 'The schema name for uploads';


--
-- Name: COLUMN metabase_database.uploads_table_prefix; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.metabase_database.uploads_table_prefix IS 'The prefix for upload table names';


--
-- Name: metabase_database_id_seq; Type: SEQUENCE; Schema: public; Owner: aman
--

ALTER TABLE public.metabase_database ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.metabase_database_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: metabase_field; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.metabase_field (
    id integer NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    name character varying(254) NOT NULL,
    base_type character varying(255) NOT NULL,
    semantic_type character varying(255),
    active boolean DEFAULT true NOT NULL,
    description text,
    preview_display boolean DEFAULT true NOT NULL,
    "position" integer DEFAULT 0 NOT NULL,
    table_id integer NOT NULL,
    parent_id integer,
    display_name character varying(254),
    visibility_type character varying(32) DEFAULT 'normal'::character varying NOT NULL,
    fk_target_field_id integer,
    last_analyzed timestamp with time zone,
    points_of_interest text,
    caveats text,
    fingerprint text,
    fingerprint_version integer DEFAULT 0 NOT NULL,
    database_type text NOT NULL,
    has_field_values text,
    settings text,
    database_position integer DEFAULT 0 NOT NULL,
    custom_position integer DEFAULT 0 NOT NULL,
    effective_type character varying(255),
    coercion_strategy character varying(255),
    nfc_path character varying(254),
    database_required boolean DEFAULT false NOT NULL,
    json_unfolding boolean DEFAULT false NOT NULL,
    database_is_auto_increment boolean DEFAULT false NOT NULL,
    database_indexed boolean,
    database_partitioned boolean
);


ALTER TABLE public.metabase_field OWNER TO aman;

--
-- Name: COLUMN metabase_field.json_unfolding; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.metabase_field.json_unfolding IS 'Enable/disable JSON unfolding for a field';


--
-- Name: COLUMN metabase_field.database_is_auto_increment; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.metabase_field.database_is_auto_increment IS 'Indicates this field is auto incremented';


--
-- Name: COLUMN metabase_field.database_indexed; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.metabase_field.database_indexed IS 'If the database supports indexing, this column indicate whether or not a field is indexed, or is the 1st column in a composite index';


--
-- Name: COLUMN metabase_field.database_partitioned; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.metabase_field.database_partitioned IS 'Whether the table is partitioned by this field';


--
-- Name: metabase_field_id_seq; Type: SEQUENCE; Schema: public; Owner: aman
--

ALTER TABLE public.metabase_field ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.metabase_field_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: metabase_fieldvalues; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.metabase_fieldvalues (
    id integer NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    "values" text,
    human_readable_values text,
    field_id integer NOT NULL,
    has_more_values boolean DEFAULT false,
    type character varying(32) DEFAULT 'full'::character varying NOT NULL,
    hash_key text,
    last_used_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.metabase_fieldvalues OWNER TO aman;

--
-- Name: COLUMN metabase_fieldvalues.last_used_at; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.metabase_fieldvalues.last_used_at IS 'Timestamp of when these FieldValues were last used.';


--
-- Name: metabase_fieldvalues_id_seq; Type: SEQUENCE; Schema: public; Owner: aman
--

ALTER TABLE public.metabase_fieldvalues ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.metabase_fieldvalues_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: metabase_table; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.metabase_table (
    id integer NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    name character varying(256) NOT NULL,
    description text,
    entity_type character varying(254),
    active boolean NOT NULL,
    db_id integer NOT NULL,
    display_name character varying(256),
    visibility_type character varying(254),
    schema character varying(254),
    points_of_interest text,
    caveats text,
    show_in_getting_started boolean DEFAULT false NOT NULL,
    field_order character varying(254) DEFAULT 'database'::character varying NOT NULL,
    initial_sync_status character varying(32) DEFAULT 'complete'::character varying NOT NULL,
    is_upload boolean DEFAULT false NOT NULL,
    database_require_filter boolean,
    estimated_row_count bigint,
    view_count integer DEFAULT 0 NOT NULL
);


ALTER TABLE public.metabase_table OWNER TO aman;

--
-- Name: COLUMN metabase_table.is_upload; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.metabase_table.is_upload IS 'Was the table created from user-uploaded (i.e., from a CSV) data?';


--
-- Name: COLUMN metabase_table.database_require_filter; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.metabase_table.database_require_filter IS 'If true, the table requires a filter to be able to query it';


--
-- Name: COLUMN metabase_table.estimated_row_count; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.metabase_table.estimated_row_count IS 'The estimated row count';


--
-- Name: COLUMN metabase_table.view_count; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.metabase_table.view_count IS 'Keeps a running count of card views';


--
-- Name: metabase_table_id_seq; Type: SEQUENCE; Schema: public; Owner: aman
--

ALTER TABLE public.metabase_table ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.metabase_table_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: metric; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.metric (
    id integer NOT NULL,
    table_id integer NOT NULL,
    creator_id integer NOT NULL,
    name character varying(254) NOT NULL,
    description text,
    archived boolean DEFAULT false NOT NULL,
    definition text NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    points_of_interest text,
    caveats text,
    how_is_this_calculated text,
    show_in_getting_started boolean DEFAULT false NOT NULL,
    entity_id character(21)
);


ALTER TABLE public.metric OWNER TO aman;

--
-- Name: metric_id_seq; Type: SEQUENCE; Schema: public; Owner: aman
--

ALTER TABLE public.metric ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.metric_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: metric_important_field; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.metric_important_field (
    id integer NOT NULL,
    metric_id integer NOT NULL,
    field_id integer NOT NULL
);


ALTER TABLE public.metric_important_field OWNER TO aman;

--
-- Name: metric_important_field_id_seq; Type: SEQUENCE; Schema: public; Owner: aman
--

ALTER TABLE public.metric_important_field ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.metric_important_field_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: model_index; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.model_index (
    id integer NOT NULL,
    model_id integer,
    pk_ref text NOT NULL,
    value_ref text NOT NULL,
    schedule text NOT NULL,
    state text NOT NULL,
    indexed_at timestamp with time zone,
    error text,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    creator_id integer NOT NULL
);


ALTER TABLE public.model_index OWNER TO aman;

--
-- Name: TABLE model_index; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON TABLE public.model_index IS 'Used to keep track of which models have indexed columns.';


--
-- Name: COLUMN model_index.model_id; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.model_index.model_id IS 'The ID of the indexed model.';


--
-- Name: COLUMN model_index.pk_ref; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.model_index.pk_ref IS 'Serialized JSON of the primary key field ref.';


--
-- Name: COLUMN model_index.value_ref; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.model_index.value_ref IS 'Serialized JSON of the label field ref.';


--
-- Name: COLUMN model_index.schedule; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.model_index.schedule IS 'The cron schedule for when value syncing should happen.';


--
-- Name: COLUMN model_index.state; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.model_index.state IS 'The status of the index: initializing, indexed, error, overflow.';


--
-- Name: COLUMN model_index.indexed_at; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.model_index.indexed_at IS 'When the status changed';


--
-- Name: COLUMN model_index.error; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.model_index.error IS 'The error message if the status is error.';


--
-- Name: COLUMN model_index.created_at; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.model_index.created_at IS 'The timestamp of when these changes were made.';


--
-- Name: COLUMN model_index.creator_id; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.model_index.creator_id IS 'ID of the user who created the event';


--
-- Name: model_index_id_seq; Type: SEQUENCE; Schema: public; Owner: aman
--

ALTER TABLE public.model_index ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.model_index_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: model_index_value; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.model_index_value (
    model_index_id integer,
    model_pk integer NOT NULL,
    name text NOT NULL
);


ALTER TABLE public.model_index_value OWNER TO aman;

--
-- Name: TABLE model_index_value; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON TABLE public.model_index_value IS 'Used to keep track of the values indexed in a model';


--
-- Name: COLUMN model_index_value.model_index_id; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.model_index_value.model_index_id IS 'The ID of the indexed model.';


--
-- Name: COLUMN model_index_value.model_pk; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.model_index_value.model_pk IS 'The primary key of the indexed value';


--
-- Name: COLUMN model_index_value.name; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.model_index_value.name IS 'The label to display identifying the indexed value.';


--
-- Name: moderation_review; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.moderation_review (
    id integer NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    status character varying(255),
    text text,
    moderated_item_id integer NOT NULL,
    moderated_item_type character varying(255) NOT NULL,
    moderator_id integer NOT NULL,
    most_recent boolean NOT NULL
);


ALTER TABLE public.moderation_review OWNER TO aman;

--
-- Name: moderation_review_id_seq; Type: SEQUENCE; Schema: public; Owner: aman
--

ALTER TABLE public.moderation_review ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.moderation_review_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: native_query_snippet; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.native_query_snippet (
    id integer NOT NULL,
    name character varying(254) NOT NULL,
    description text,
    content text NOT NULL,
    creator_id integer NOT NULL,
    archived boolean DEFAULT false NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    collection_id integer,
    entity_id character(21)
);


ALTER TABLE public.native_query_snippet OWNER TO aman;

--
-- Name: native_query_snippet_id_seq; Type: SEQUENCE; Schema: public; Owner: aman
--

ALTER TABLE public.native_query_snippet ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.native_query_snippet_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: parameter_card; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.parameter_card (
    id integer NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    card_id integer NOT NULL,
    parameterized_object_type character varying(32) NOT NULL,
    parameterized_object_id integer NOT NULL,
    parameter_id character varying(36) NOT NULL
);


ALTER TABLE public.parameter_card OWNER TO aman;

--
-- Name: TABLE parameter_card; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON TABLE public.parameter_card IS 'Join table connecting cards to entities (dashboards, other cards, etc.) that use the values generated by the card for filter values';


--
-- Name: COLUMN parameter_card.updated_at; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.parameter_card.updated_at IS 'most recent modification time';


--
-- Name: COLUMN parameter_card.created_at; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.parameter_card.created_at IS 'creation time';


--
-- Name: COLUMN parameter_card.card_id; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.parameter_card.card_id IS 'ID of the card generating the values';


--
-- Name: COLUMN parameter_card.parameterized_object_type; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.parameter_card.parameterized_object_type IS 'Type of the entity consuming the values (dashboard, card, etc.)';


--
-- Name: COLUMN parameter_card.parameterized_object_id; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.parameter_card.parameterized_object_id IS 'ID of the entity consuming the values';


--
-- Name: COLUMN parameter_card.parameter_id; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.parameter_card.parameter_id IS 'The parameter ID';


--
-- Name: parameter_card_id_seq; Type: SEQUENCE; Schema: public; Owner: aman
--

ALTER TABLE public.parameter_card ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.parameter_card_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: permissions; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.permissions (
    id integer NOT NULL,
    object character varying(254) NOT NULL,
    group_id integer NOT NULL
);


ALTER TABLE public.permissions OWNER TO aman;

--
-- Name: permissions_group; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.permissions_group (
    id integer NOT NULL,
    name character varying(255) NOT NULL,
    entity_id character(21)
);


ALTER TABLE public.permissions_group OWNER TO aman;

--
-- Name: COLUMN permissions_group.entity_id; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.permissions_group.entity_id IS 'NanoID tag for each user';


--
-- Name: permissions_group_id_seq; Type: SEQUENCE; Schema: public; Owner: aman
--

ALTER TABLE public.permissions_group ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.permissions_group_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: permissions_group_membership; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.permissions_group_membership (
    id integer NOT NULL,
    user_id integer NOT NULL,
    group_id integer NOT NULL,
    is_group_manager boolean DEFAULT false NOT NULL
);


ALTER TABLE public.permissions_group_membership OWNER TO aman;

--
-- Name: permissions_group_membership_id_seq; Type: SEQUENCE; Schema: public; Owner: aman
--

ALTER TABLE public.permissions_group_membership ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.permissions_group_membership_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: permissions_id_seq; Type: SEQUENCE; Schema: public; Owner: aman
--

ALTER TABLE public.permissions ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.permissions_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: permissions_revision; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.permissions_revision (
    id integer NOT NULL,
    before text NOT NULL,
    after text NOT NULL,
    user_id integer NOT NULL,
    created_at timestamp with time zone NOT NULL,
    remark text
);


ALTER TABLE public.permissions_revision OWNER TO aman;

--
-- Name: permissions_revision_id_seq; Type: SEQUENCE; Schema: public; Owner: aman
--

ALTER TABLE public.permissions_revision ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.permissions_revision_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: persisted_info; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.persisted_info (
    id integer NOT NULL,
    database_id integer NOT NULL,
    card_id integer NOT NULL,
    question_slug text NOT NULL,
    table_name text NOT NULL,
    definition text,
    query_hash text,
    active boolean DEFAULT false NOT NULL,
    state text NOT NULL,
    refresh_begin timestamp with time zone NOT NULL,
    refresh_end timestamp with time zone,
    state_change_at timestamp with time zone,
    error text,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    creator_id integer
);


ALTER TABLE public.persisted_info OWNER TO aman;

--
-- Name: persisted_info_id_seq; Type: SEQUENCE; Schema: public; Owner: aman
--

ALTER TABLE public.persisted_info ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.persisted_info_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: pulse; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.pulse (
    id integer NOT NULL,
    creator_id integer NOT NULL,
    name character varying(254),
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    skip_if_empty boolean DEFAULT false NOT NULL,
    alert_condition character varying(254),
    alert_first_only boolean,
    alert_above_goal boolean,
    collection_id integer,
    collection_position smallint,
    archived boolean DEFAULT false,
    dashboard_id integer,
    parameters text NOT NULL,
    entity_id character(21)
);


ALTER TABLE public.pulse OWNER TO aman;

--
-- Name: pulse_card; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.pulse_card (
    id integer NOT NULL,
    pulse_id integer NOT NULL,
    card_id integer NOT NULL,
    "position" integer NOT NULL,
    include_csv boolean DEFAULT false NOT NULL,
    include_xls boolean DEFAULT false NOT NULL,
    dashboard_card_id integer,
    entity_id character(21),
    format_rows boolean DEFAULT true
);


ALTER TABLE public.pulse_card OWNER TO aman;

--
-- Name: COLUMN pulse_card.format_rows; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.pulse_card.format_rows IS 'Whether or not to apply formatting to the rows of the export';


--
-- Name: pulse_card_id_seq; Type: SEQUENCE; Schema: public; Owner: aman
--

ALTER TABLE public.pulse_card ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.pulse_card_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: pulse_channel; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.pulse_channel (
    id integer NOT NULL,
    pulse_id integer NOT NULL,
    channel_type character varying(32) NOT NULL,
    details text NOT NULL,
    schedule_type character varying(32) NOT NULL,
    schedule_hour integer,
    schedule_day character varying(64),
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    schedule_frame character varying(32),
    enabled boolean DEFAULT true NOT NULL,
    entity_id character(21)
);


ALTER TABLE public.pulse_channel OWNER TO aman;

--
-- Name: pulse_channel_id_seq; Type: SEQUENCE; Schema: public; Owner: aman
--

ALTER TABLE public.pulse_channel ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.pulse_channel_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: pulse_channel_recipient; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.pulse_channel_recipient (
    id integer NOT NULL,
    pulse_channel_id integer NOT NULL,
    user_id integer NOT NULL
);


ALTER TABLE public.pulse_channel_recipient OWNER TO aman;

--
-- Name: pulse_channel_recipient_id_seq; Type: SEQUENCE; Schema: public; Owner: aman
--

ALTER TABLE public.pulse_channel_recipient ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.pulse_channel_recipient_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: pulse_id_seq; Type: SEQUENCE; Schema: public; Owner: aman
--

ALTER TABLE public.pulse ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.pulse_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: qrtz_blob_triggers; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.qrtz_blob_triggers (
    sched_name character varying(120) NOT NULL,
    trigger_name character varying(200) NOT NULL,
    trigger_group character varying(200) NOT NULL,
    blob_data bytea
);


ALTER TABLE public.qrtz_blob_triggers OWNER TO aman;

--
-- Name: qrtz_calendars; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.qrtz_calendars (
    sched_name character varying(120) NOT NULL,
    calendar_name character varying(200) NOT NULL,
    calendar bytea NOT NULL
);


ALTER TABLE public.qrtz_calendars OWNER TO aman;

--
-- Name: qrtz_cron_triggers; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.qrtz_cron_triggers (
    sched_name character varying(120) NOT NULL,
    trigger_name character varying(200) NOT NULL,
    trigger_group character varying(200) NOT NULL,
    cron_expression character varying(120) NOT NULL,
    time_zone_id character varying(80)
);


ALTER TABLE public.qrtz_cron_triggers OWNER TO aman;

--
-- Name: qrtz_fired_triggers; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.qrtz_fired_triggers (
    sched_name character varying(120) NOT NULL,
    entry_id character varying(95) NOT NULL,
    trigger_name character varying(200) NOT NULL,
    trigger_group character varying(200) NOT NULL,
    instance_name character varying(200) NOT NULL,
    fired_time bigint NOT NULL,
    sched_time bigint,
    priority integer NOT NULL,
    state character varying(16) NOT NULL,
    job_name character varying(200),
    job_group character varying(200),
    is_nonconcurrent boolean,
    requests_recovery boolean
);


ALTER TABLE public.qrtz_fired_triggers OWNER TO aman;

--
-- Name: qrtz_job_details; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.qrtz_job_details (
    sched_name character varying(120) NOT NULL,
    job_name character varying(200) NOT NULL,
    job_group character varying(200) NOT NULL,
    description character varying(250),
    job_class_name character varying(250) NOT NULL,
    is_durable boolean NOT NULL,
    is_nonconcurrent boolean NOT NULL,
    is_update_data boolean NOT NULL,
    requests_recovery boolean NOT NULL,
    job_data bytea
);


ALTER TABLE public.qrtz_job_details OWNER TO aman;

--
-- Name: qrtz_locks; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.qrtz_locks (
    sched_name character varying(120) NOT NULL,
    lock_name character varying(40) NOT NULL
);


ALTER TABLE public.qrtz_locks OWNER TO aman;

--
-- Name: qrtz_paused_trigger_grps; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.qrtz_paused_trigger_grps (
    sched_name character varying(120) NOT NULL,
    trigger_group character varying(200) NOT NULL
);


ALTER TABLE public.qrtz_paused_trigger_grps OWNER TO aman;

--
-- Name: qrtz_scheduler_state; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.qrtz_scheduler_state (
    sched_name character varying(120) NOT NULL,
    instance_name character varying(200) NOT NULL,
    last_checkin_time bigint NOT NULL,
    checkin_interval bigint NOT NULL
);


ALTER TABLE public.qrtz_scheduler_state OWNER TO aman;

--
-- Name: qrtz_simple_triggers; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.qrtz_simple_triggers (
    sched_name character varying(120) NOT NULL,
    trigger_name character varying(200) NOT NULL,
    trigger_group character varying(200) NOT NULL,
    repeat_count bigint NOT NULL,
    repeat_interval bigint NOT NULL,
    times_triggered bigint NOT NULL
);


ALTER TABLE public.qrtz_simple_triggers OWNER TO aman;

--
-- Name: qrtz_simprop_triggers; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.qrtz_simprop_triggers (
    sched_name character varying(120) NOT NULL,
    trigger_name character varying(200) NOT NULL,
    trigger_group character varying(200) NOT NULL,
    str_prop_1 character varying(512),
    str_prop_2 character varying(512),
    str_prop_3 character varying(512),
    int_prop_1 integer,
    int_prop_2 integer,
    long_prop_1 bigint,
    long_prop_2 bigint,
    dec_prop_1 numeric(13,4),
    dec_prop_2 numeric(13,4),
    bool_prop_1 boolean,
    bool_prop_2 boolean
);


ALTER TABLE public.qrtz_simprop_triggers OWNER TO aman;

--
-- Name: qrtz_triggers; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.qrtz_triggers (
    sched_name character varying(120) NOT NULL,
    trigger_name character varying(200) NOT NULL,
    trigger_group character varying(200) NOT NULL,
    job_name character varying(200) NOT NULL,
    job_group character varying(200) NOT NULL,
    description character varying(250),
    next_fire_time bigint,
    prev_fire_time bigint,
    priority integer,
    trigger_state character varying(16) NOT NULL,
    trigger_type character varying(8) NOT NULL,
    start_time bigint NOT NULL,
    end_time bigint,
    calendar_name character varying(200),
    misfire_instr smallint,
    job_data bytea
);


ALTER TABLE public.qrtz_triggers OWNER TO aman;

--
-- Name: query; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.query (
    query_hash bytea NOT NULL,
    average_execution_time integer NOT NULL,
    query text
);


ALTER TABLE public.query OWNER TO aman;

--
-- Name: query_action; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.query_action (
    action_id integer NOT NULL,
    database_id integer NOT NULL,
    dataset_query text NOT NULL
);


ALTER TABLE public.query_action OWNER TO aman;

--
-- Name: TABLE query_action; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON TABLE public.query_action IS 'A readwrite query type of action';


--
-- Name: COLUMN query_action.action_id; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.query_action.action_id IS 'The related action';


--
-- Name: COLUMN query_action.database_id; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.query_action.database_id IS 'The associated database';


--
-- Name: COLUMN query_action.dataset_query; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.query_action.dataset_query IS 'The MBQL writeback query';


--
-- Name: query_cache; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.query_cache (
    query_hash bytea NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    results bytea NOT NULL
);


ALTER TABLE public.query_cache OWNER TO aman;

--
-- Name: query_execution; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.query_execution (
    id integer NOT NULL,
    hash bytea NOT NULL,
    started_at timestamp with time zone NOT NULL,
    running_time integer NOT NULL,
    result_rows integer NOT NULL,
    native boolean NOT NULL,
    context character varying(32),
    error text,
    executor_id integer,
    card_id integer,
    dashboard_id integer,
    pulse_id integer,
    database_id integer,
    cache_hit boolean,
    action_id integer,
    is_sandboxed boolean,
    cache_hash bytea
);


ALTER TABLE public.query_execution OWNER TO aman;

--
-- Name: COLUMN query_execution.action_id; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.query_execution.action_id IS 'The ID of the action associated with this query execution, if any.';


--
-- Name: COLUMN query_execution.is_sandboxed; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.query_execution.is_sandboxed IS 'Is query from a sandboxed user';


--
-- Name: COLUMN query_execution.cache_hash; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.query_execution.cache_hash IS 'Hash of normalized query, calculated in middleware.cache';


--
-- Name: query_execution_id_seq; Type: SEQUENCE; Schema: public; Owner: aman
--

ALTER TABLE public.query_execution ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.query_execution_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: query_field; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.query_field (
    id integer NOT NULL,
    card_id integer NOT NULL,
    field_id integer NOT NULL,
    direct_reference boolean DEFAULT true NOT NULL
);


ALTER TABLE public.query_field OWNER TO aman;

--
-- Name: TABLE query_field; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON TABLE public.query_field IS 'Fields used by a card''s query';


--
-- Name: COLUMN query_field.id; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.query_field.id IS 'PK';


--
-- Name: COLUMN query_field.card_id; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.query_field.card_id IS 'referenced card';


--
-- Name: COLUMN query_field.field_id; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.query_field.field_id IS 'referenced field';


--
-- Name: COLUMN query_field.direct_reference; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.query_field.direct_reference IS 'Is the Field referenced directly or via a wildcard';


--
-- Name: query_field_id_seq; Type: SEQUENCE; Schema: public; Owner: aman
--

ALTER TABLE public.query_field ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.query_field_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: recent_views; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.recent_views (
    id integer NOT NULL,
    user_id integer NOT NULL,
    model character varying(16) NOT NULL,
    model_id integer NOT NULL,
    "timestamp" timestamp with time zone NOT NULL
);


ALTER TABLE public.recent_views OWNER TO aman;

--
-- Name: TABLE recent_views; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON TABLE public.recent_views IS 'Used to store recently viewed objects for each user';


--
-- Name: COLUMN recent_views.user_id; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.recent_views.user_id IS 'The user associated with this view';


--
-- Name: COLUMN recent_views.model; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.recent_views.model IS 'The name of the model that was viewed';


--
-- Name: COLUMN recent_views.model_id; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.recent_views.model_id IS 'The ID of the model that was viewed';


--
-- Name: COLUMN recent_views."timestamp"; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.recent_views."timestamp" IS 'The time a view was recorded';


--
-- Name: recent_views_id_seq; Type: SEQUENCE; Schema: public; Owner: aman
--

ALTER TABLE public.recent_views ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.recent_views_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: report_card; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.report_card (
    id integer NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    name character varying(254) NOT NULL,
    description text,
    display character varying(254) NOT NULL,
    dataset_query text NOT NULL,
    visualization_settings text NOT NULL,
    creator_id integer NOT NULL,
    database_id integer NOT NULL,
    table_id integer,
    query_type character varying(16),
    archived boolean DEFAULT false NOT NULL,
    collection_id integer,
    public_uuid character(36),
    made_public_by_id integer,
    enable_embedding boolean DEFAULT false NOT NULL,
    embedding_params text,
    cache_ttl integer,
    result_metadata text,
    collection_position smallint,
    entity_id character(21),
    parameters text,
    parameter_mappings text,
    collection_preview boolean DEFAULT true NOT NULL,
    metabase_version character varying(100),
    type character varying(16) DEFAULT 'question'::character varying NOT NULL,
    initially_published_at timestamp with time zone,
    cache_invalidated_at timestamp with time zone,
    last_used_at timestamp with time zone,
    view_count integer DEFAULT 0 NOT NULL
);


ALTER TABLE public.report_card OWNER TO aman;

--
-- Name: COLUMN report_card.metabase_version; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.report_card.metabase_version IS 'Metabase version used to create the card.';


--
-- Name: COLUMN report_card.type; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.report_card.type IS 'The type of card, could be ''question'', ''model'', ''metric''';


--
-- Name: COLUMN report_card.initially_published_at; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.report_card.initially_published_at IS 'The timestamp when the card was first published in a static embed';


--
-- Name: COLUMN report_card.cache_invalidated_at; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.report_card.cache_invalidated_at IS 'An invalidation time that can supersede cache_config.invalidated_at';


--
-- Name: COLUMN report_card.last_used_at; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.report_card.last_used_at IS 'The timestamp of when the card is last used';


--
-- Name: COLUMN report_card.view_count; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.report_card.view_count IS 'Keeps a running count of card views';


--
-- Name: report_card_id_seq; Type: SEQUENCE; Schema: public; Owner: aman
--

ALTER TABLE public.report_card ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.report_card_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: report_cardfavorite; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.report_cardfavorite (
    id integer NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    card_id integer NOT NULL,
    owner_id integer NOT NULL
);


ALTER TABLE public.report_cardfavorite OWNER TO aman;

--
-- Name: report_cardfavorite_id_seq; Type: SEQUENCE; Schema: public; Owner: aman
--

ALTER TABLE public.report_cardfavorite ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.report_cardfavorite_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: report_dashboard; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.report_dashboard (
    id integer NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    name character varying(254) NOT NULL,
    description text,
    creator_id integer NOT NULL,
    parameters text NOT NULL,
    points_of_interest text,
    caveats text,
    show_in_getting_started boolean DEFAULT false NOT NULL,
    public_uuid character(36),
    made_public_by_id integer,
    enable_embedding boolean DEFAULT false NOT NULL,
    embedding_params text,
    archived boolean DEFAULT false NOT NULL,
    "position" integer,
    collection_id integer,
    collection_position smallint,
    cache_ttl integer,
    entity_id character(21),
    auto_apply_filters boolean DEFAULT true NOT NULL,
    width character varying(16) DEFAULT 'fixed'::character varying NOT NULL,
    initially_published_at timestamp with time zone,
    view_count integer DEFAULT 0 NOT NULL
);


ALTER TABLE public.report_dashboard OWNER TO aman;

--
-- Name: COLUMN report_dashboard.auto_apply_filters; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.report_dashboard.auto_apply_filters IS 'Whether or not to auto-apply filters on a dashboard';


--
-- Name: COLUMN report_dashboard.width; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.report_dashboard.width IS 'The value of the dashboard''s width setting can be fixed or full. New dashboards will be set to fixed';


--
-- Name: COLUMN report_dashboard.initially_published_at; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.report_dashboard.initially_published_at IS 'The timestamp when the dashboard was first published in a static embed';


--
-- Name: COLUMN report_dashboard.view_count; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.report_dashboard.view_count IS 'Keeps a running count of dashboard views';


--
-- Name: report_dashboard_id_seq; Type: SEQUENCE; Schema: public; Owner: aman
--

ALTER TABLE public.report_dashboard ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.report_dashboard_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: report_dashboardcard; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.report_dashboardcard (
    id integer NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    size_x integer NOT NULL,
    size_y integer NOT NULL,
    "row" integer NOT NULL,
    col integer NOT NULL,
    card_id integer,
    dashboard_id integer NOT NULL,
    parameter_mappings text NOT NULL,
    visualization_settings text NOT NULL,
    entity_id character(21),
    action_id integer,
    dashboard_tab_id integer
);


ALTER TABLE public.report_dashboardcard OWNER TO aman;

--
-- Name: COLUMN report_dashboardcard.action_id; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.report_dashboardcard.action_id IS 'The related action';


--
-- Name: COLUMN report_dashboardcard.dashboard_tab_id; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.report_dashboardcard.dashboard_tab_id IS 'The referenced tab id that dashcard is on, it''s nullable for dashboard with no tab';


--
-- Name: report_dashboardcard_id_seq; Type: SEQUENCE; Schema: public; Owner: aman
--

ALTER TABLE public.report_dashboardcard ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.report_dashboardcard_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: revision; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.revision (
    id integer NOT NULL,
    model character varying(16) NOT NULL,
    model_id integer NOT NULL,
    user_id integer NOT NULL,
    "timestamp" timestamp with time zone NOT NULL,
    object text NOT NULL,
    is_reversion boolean DEFAULT false NOT NULL,
    is_creation boolean DEFAULT false NOT NULL,
    message text,
    most_recent boolean DEFAULT false NOT NULL,
    metabase_version character varying(100)
);


ALTER TABLE public.revision OWNER TO aman;

--
-- Name: COLUMN revision.most_recent; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.revision.most_recent IS 'Whether a revision is the most recent one';


--
-- Name: COLUMN revision.metabase_version; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.revision.metabase_version IS 'Metabase version used to create the revision.';


--
-- Name: revision_id_seq; Type: SEQUENCE; Schema: public; Owner: aman
--

ALTER TABLE public.revision ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.revision_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: secret; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.secret (
    id integer NOT NULL,
    version integer DEFAULT 1 NOT NULL,
    creator_id integer,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone,
    name character varying(254) NOT NULL,
    kind character varying(254) NOT NULL,
    source character varying(254),
    value bytea NOT NULL
);


ALTER TABLE public.secret OWNER TO aman;

--
-- Name: secret_id_seq; Type: SEQUENCE; Schema: public; Owner: aman
--

ALTER TABLE public.secret ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.secret_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: segment; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.segment (
    id integer NOT NULL,
    table_id integer NOT NULL,
    creator_id integer NOT NULL,
    name character varying(254) NOT NULL,
    description text,
    archived boolean DEFAULT false NOT NULL,
    definition text NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    points_of_interest text,
    caveats text,
    show_in_getting_started boolean DEFAULT false NOT NULL,
    entity_id character(21)
);


ALTER TABLE public.segment OWNER TO aman;

--
-- Name: segment_id_seq; Type: SEQUENCE; Schema: public; Owner: aman
--

ALTER TABLE public.segment ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.segment_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: setting; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.setting (
    key character varying(254) NOT NULL,
    value text NOT NULL
);


ALTER TABLE public.setting OWNER TO aman;

--
-- Name: table_privileges; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.table_privileges (
    table_id integer NOT NULL,
    role character varying(255),
    "select" boolean DEFAULT false NOT NULL,
    update boolean DEFAULT false NOT NULL,
    insert boolean DEFAULT false NOT NULL,
    delete boolean DEFAULT false NOT NULL
);


ALTER TABLE public.table_privileges OWNER TO aman;

--
-- Name: TABLE table_privileges; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON TABLE public.table_privileges IS 'Table for user and role privileges by table';


--
-- Name: COLUMN table_privileges.table_id; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.table_privileges.table_id IS 'Table ID';


--
-- Name: COLUMN table_privileges.role; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.table_privileges.role IS 'Role name. NULL indicates the privileges are the current user''s';


--
-- Name: COLUMN table_privileges."select"; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.table_privileges."select" IS 'Privilege to select from the table';


--
-- Name: COLUMN table_privileges.update; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.table_privileges.update IS 'Privilege to update records in the table';


--
-- Name: COLUMN table_privileges.insert; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.table_privileges.insert IS 'Privilege to insert records into the table';


--
-- Name: COLUMN table_privileges.delete; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.table_privileges.delete IS 'Privilege to delete records from the table';


--
-- Name: task_history; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.task_history (
    id integer NOT NULL,
    task character varying(254) NOT NULL,
    db_id integer,
    started_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    ended_at timestamp with time zone,
    duration integer,
    task_details text,
    status character varying(21) DEFAULT 'started'::character varying NOT NULL
);


ALTER TABLE public.task_history OWNER TO aman;

--
-- Name: COLUMN task_history.status; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.task_history.status IS 'the status of task history, could be started, failed, success, unknown';


--
-- Name: task_history_id_seq; Type: SEQUENCE; Schema: public; Owner: aman
--

ALTER TABLE public.task_history ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.task_history_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: timeline; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.timeline (
    id integer NOT NULL,
    name character varying(255) NOT NULL,
    description character varying(255),
    icon character varying(128) NOT NULL,
    collection_id integer,
    archived boolean DEFAULT false NOT NULL,
    creator_id integer NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL,
    "default" boolean DEFAULT false NOT NULL,
    entity_id character(21)
);


ALTER TABLE public.timeline OWNER TO aman;

--
-- Name: timeline_event; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.timeline_event (
    id integer NOT NULL,
    timeline_id integer NOT NULL,
    name character varying(255) NOT NULL,
    description character varying(255),
    "timestamp" timestamp with time zone NOT NULL,
    time_matters boolean NOT NULL,
    timezone character varying(255) NOT NULL,
    icon character varying(128) NOT NULL,
    archived boolean DEFAULT false NOT NULL,
    creator_id integer NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT now() NOT NULL
);


ALTER TABLE public.timeline_event OWNER TO aman;

--
-- Name: timeline_event_id_seq; Type: SEQUENCE; Schema: public; Owner: aman
--

ALTER TABLE public.timeline_event ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.timeline_event_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: timeline_id_seq; Type: SEQUENCE; Schema: public; Owner: aman
--

ALTER TABLE public.timeline ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.timeline_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: user_parameter_value; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.user_parameter_value (
    id integer NOT NULL,
    user_id integer NOT NULL,
    parameter_id character varying(36) NOT NULL,
    value text
);


ALTER TABLE public.user_parameter_value OWNER TO aman;

--
-- Name: TABLE user_parameter_value; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON TABLE public.user_parameter_value IS 'Table holding last set value of a parameter per user';


--
-- Name: COLUMN user_parameter_value.user_id; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.user_parameter_value.user_id IS 'ID of the User who has set the parameter value';


--
-- Name: COLUMN user_parameter_value.parameter_id; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.user_parameter_value.parameter_id IS 'The parameter ID';


--
-- Name: COLUMN user_parameter_value.value; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.user_parameter_value.value IS 'Value of the parameter';


--
-- Name: user_parameter_value_id_seq; Type: SEQUENCE; Schema: public; Owner: aman
--

ALTER TABLE public.user_parameter_value ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.user_parameter_value_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Name: v_alerts; Type: VIEW; Schema: public; Owner: aman
--

CREATE VIEW public.v_alerts AS
 WITH agg_recipients AS (
         SELECT pulse_channel_recipient.pulse_channel_id,
            string_agg((core_user.email)::text, ','::text) AS recipients
           FROM (public.pulse_channel_recipient
             LEFT JOIN public.core_user ON ((pulse_channel_recipient.user_id = core_user.id)))
          GROUP BY pulse_channel_recipient.pulse_channel_id
        )
 SELECT pulse.id AS entity_id,
    ('pulse_'::text || pulse.id) AS entity_qualified_id,
    pulse.created_at,
    pulse.updated_at,
    pulse.creator_id,
    pulse_card.card_id,
    ('card_'::text || pulse_card.card_id) AS card_qualified_id,
    pulse.alert_condition,
    pulse_channel.schedule_type,
    pulse_channel.schedule_day,
    pulse_channel.schedule_hour,
    pulse.archived,
    pulse_channel.channel_type AS recipient_type,
    agg_recipients.recipients,
    pulse_channel.details AS recipient_external
   FROM (((public.pulse
     LEFT JOIN public.pulse_card ON ((pulse.id = pulse_card.pulse_id)))
     LEFT JOIN public.pulse_channel ON ((pulse.id = pulse_channel.pulse_id)))
     LEFT JOIN agg_recipients ON ((pulse_channel.id = agg_recipients.pulse_channel_id)))
  WHERE (pulse.alert_condition IS NOT NULL);


ALTER VIEW public.v_alerts OWNER TO aman;

--
-- Name: v_audit_log; Type: VIEW; Schema: public; Owner: aman
--

CREATE VIEW public.v_audit_log AS
 SELECT id,
        CASE
            WHEN ((topic)::text = 'card-create'::text) THEN 'card-create'::character varying
            WHEN ((topic)::text = 'card-delete'::text) THEN 'card-delete'::character varying
            WHEN ((topic)::text = 'card-update'::text) THEN 'card-update'::character varying
            WHEN ((topic)::text = 'pulse-create'::text) THEN 'subscription-create'::character varying
            WHEN ((topic)::text = 'pulse-delete'::text) THEN 'subscription-delete'::character varying
            ELSE topic
        END AS topic,
    "timestamp",
    NULL::text AS end_timestamp,
    COALESCE(user_id, 0) AS user_id,
    lower((model)::text) AS entity_type,
    model_id AS entity_id,
        CASE
            WHEN ((model)::text = 'Dataset'::text) THEN ('card_'::text || model_id)
            WHEN (model_id IS NULL) THEN NULL::text
            ELSE ((lower((model)::text) || '_'::text) || model_id)
        END AS entity_qualified_id,
    details
   FROM public.audit_log
  WHERE ((topic)::text <> ALL ((ARRAY['card-read'::character varying, 'card-query'::character varying, 'dashboard-read'::character varying, 'dashboard-query'::character varying, 'table-read'::character varying])::text[]));


ALTER VIEW public.v_audit_log OWNER TO aman;

--
-- Name: v_content; Type: VIEW; Schema: public; Owner: aman
--

CREATE VIEW public.v_content AS
 SELECT action.id AS entity_id,
    ('action_'::text || action.id) AS entity_qualified_id,
    'action'::text AS entity_type,
    action.created_at,
    action.updated_at,
    action.creator_id,
    action.name,
    action.description,
    NULL::integer AS collection_id,
    action.made_public_by_id AS made_public_by_user,
    NULL::boolean AS is_embedding_enabled,
    action.archived,
    action.type AS action_type,
    action.model_id AS action_model_id,
    NULL::boolean AS collection_is_official,
    NULL::boolean AS collection_is_personal,
    NULL::text AS question_viz_type,
    NULL::text AS question_database_id,
    NULL::boolean AS question_is_native,
    NULL::timestamp without time zone AS event_timestamp
   FROM public.action
UNION
 SELECT collection.id AS entity_id,
    ('collection_'::text || collection.id) AS entity_qualified_id,
    'collection'::text AS entity_type,
    collection.created_at,
    NULL::timestamp with time zone AS updated_at,
    NULL::integer AS creator_id,
    collection.name,
    collection.description,
    NULL::integer AS collection_id,
    NULL::integer AS made_public_by_user,
    NULL::boolean AS is_embedding_enabled,
    collection.archived,
    NULL::text AS action_type,
    NULL::integer AS action_model_id,
        CASE
            WHEN ((collection.authority_level)::text = 'official'::text) THEN true
            ELSE false
        END AS collection_is_official,
        CASE
            WHEN (collection.personal_owner_id IS NOT NULL) THEN true
            ELSE false
        END AS collection_is_personal,
    NULL::text AS question_viz_type,
    NULL::text AS question_database_id,
    NULL::boolean AS question_is_native,
    NULL::timestamp without time zone AS event_timestamp
   FROM public.collection
UNION
 SELECT report_card.id AS entity_id,
    ('card_'::text || report_card.id) AS entity_qualified_id,
    report_card.type AS entity_type,
    report_card.created_at,
    report_card.updated_at,
    report_card.creator_id,
    report_card.name,
    report_card.description,
    report_card.collection_id,
    report_card.made_public_by_id AS made_public_by_user,
    report_card.enable_embedding AS is_embedding_enabled,
    report_card.archived,
    NULL::text AS action_type,
    NULL::integer AS action_model_id,
    NULL::boolean AS collection_is_official,
    NULL::boolean AS collection_is_personal,
    report_card.display AS question_viz_type,
    ('database_'::text || report_card.database_id) AS question_database_id,
        CASE
            WHEN ((report_card.query_type)::text = 'native'::text) THEN true
            ELSE false
        END AS question_is_native,
    NULL::timestamp without time zone AS event_timestamp
   FROM public.report_card
UNION
 SELECT report_dashboard.id AS entity_id,
    ('dashboard_'::text || report_dashboard.id) AS entity_qualified_id,
    'dashboard'::text AS entity_type,
    report_dashboard.created_at,
    report_dashboard.updated_at,
    report_dashboard.creator_id,
    report_dashboard.name,
    report_dashboard.description,
    report_dashboard.collection_id,
    report_dashboard.made_public_by_id AS made_public_by_user,
    report_dashboard.enable_embedding AS is_embedding_enabled,
    report_dashboard.archived,
    NULL::text AS action_type,
    NULL::integer AS action_model_id,
    NULL::boolean AS collection_is_official,
    NULL::boolean AS collection_is_personal,
    NULL::text AS question_viz_type,
    NULL::text AS question_database_id,
    NULL::boolean AS question_is_native,
    NULL::timestamp without time zone AS event_timestamp
   FROM public.report_dashboard
UNION
 SELECT event.id AS entity_id,
    ('event_'::text || event.id) AS entity_qualified_id,
    'event'::text AS entity_type,
    event.created_at,
    event.updated_at,
    event.creator_id,
    event.name,
    event.description,
    timeline.collection_id,
    NULL::integer AS made_public_by_user,
    NULL::boolean AS is_embedding_enabled,
    event.archived,
    NULL::text AS action_type,
    NULL::integer AS action_model_id,
    NULL::boolean AS collection_is_official,
    NULL::boolean AS collection_is_personal,
    NULL::text AS question_viz_type,
    NULL::text AS question_database_id,
    NULL::boolean AS question_is_native,
    event."timestamp" AS event_timestamp
   FROM (public.timeline_event event
     LEFT JOIN public.timeline ON ((event.timeline_id = timeline.id)));


ALTER VIEW public.v_content OWNER TO aman;

--
-- Name: v_dashboardcard; Type: VIEW; Schema: public; Owner: aman
--

CREATE VIEW public.v_dashboardcard AS
 SELECT id AS entity_id,
    concat('dashboardcard_', id) AS entity_qualified_id,
    concat('dashboard_', dashboard_id) AS dashboard_qualified_id,
    concat('dashboardtab_', dashboard_tab_id) AS dashboardtab_id,
    concat('card_', card_id) AS card_qualified_id,
    created_at,
    updated_at,
    size_x,
    size_y,
    visualization_settings,
    parameter_mappings
   FROM public.report_dashboardcard;


ALTER VIEW public.v_dashboardcard OWNER TO aman;

--
-- Name: v_databases; Type: VIEW; Schema: public; Owner: aman
--

CREATE VIEW public.v_databases AS
 SELECT id AS entity_id,
    concat('database_', id) AS entity_qualified_id,
    created_at,
    updated_at,
    name,
    description,
    engine AS database_type,
    metadata_sync_schedule,
    cache_field_values_schedule,
    timezone,
    is_on_demand,
    auto_run_queries,
    cache_ttl,
    creator_id,
    dbms_version AS db_version
   FROM public.metabase_database
  WHERE (id <> 13371337);


ALTER VIEW public.v_databases OWNER TO aman;

--
-- Name: v_fields; Type: VIEW; Schema: public; Owner: aman
--

CREATE VIEW public.v_fields AS
 SELECT id AS entity_id,
    ('field_'::text || id) AS entity_qualified_id,
    created_at,
    updated_at,
    name,
    display_name,
    description,
    base_type,
    visibility_type,
    fk_target_field_id,
    has_field_values,
    active,
    table_id
   FROM public.metabase_field;


ALTER VIEW public.v_fields OWNER TO aman;

--
-- Name: v_group_members; Type: VIEW; Schema: public; Owner: aman
--

CREATE VIEW public.v_group_members AS
 SELECT permissions_group_membership.user_id,
    permissions_group.id AS group_id,
    permissions_group.name AS group_name
   FROM (public.permissions_group_membership
     LEFT JOIN public.permissions_group ON ((permissions_group_membership.group_id = permissions_group.id)))
UNION
 SELECT 0 AS user_id,
    0 AS group_id,
    'Anonymous users'::character varying AS group_name;


ALTER VIEW public.v_group_members OWNER TO aman;

--
-- Name: v_query_log; Type: VIEW; Schema: public; Owner: aman
--

CREATE VIEW public.v_query_log AS
 SELECT query_execution.id AS entity_id,
    query_execution.started_at,
    ((query_execution.running_time)::double precision / (1000)::double precision) AS running_time_seconds,
    query_execution.result_rows,
    query_execution.native AS is_native,
    query_execution.context AS query_source,
    query_execution.error,
    COALESCE(query_execution.executor_id, 0) AS user_id,
    query_execution.card_id,
    ('card_'::text || query_execution.card_id) AS card_qualified_id,
    query_execution.dashboard_id,
    ('dashboard_'::text || query_execution.dashboard_id) AS dashboard_qualified_id,
    query_execution.pulse_id,
    query_execution.database_id,
    ('database_'::text || query_execution.database_id) AS database_qualified_id,
    query_execution.cache_hit,
    query_execution.action_id,
    ('action_'::text || query_execution.action_id) AS action_qualified_id,
    query.query
   FROM (public.query_execution
     LEFT JOIN public.query ON ((query_execution.hash = query.query_hash)));


ALTER VIEW public.v_query_log OWNER TO aman;

--
-- Name: v_subscriptions; Type: VIEW; Schema: public; Owner: aman
--

CREATE VIEW public.v_subscriptions AS
 WITH agg_recipients AS (
         SELECT pulse_channel_recipient.pulse_channel_id,
            string_agg((core_user.email)::text, ','::text) AS recipients
           FROM (public.pulse_channel_recipient
             LEFT JOIN public.core_user ON ((pulse_channel_recipient.user_id = core_user.id)))
          GROUP BY pulse_channel_recipient.pulse_channel_id
        )
 SELECT pulse.id AS entity_id,
    ('pulse_'::text || pulse.id) AS entity_qualified_id,
    pulse.created_at,
    pulse.updated_at,
    pulse.creator_id,
    pulse.archived,
    ('dashboard_'::text || pulse.dashboard_id) AS dashboard_qualified_id,
    pulse_channel.schedule_type,
    pulse_channel.schedule_day,
    pulse_channel.schedule_hour,
    pulse_channel.channel_type AS recipient_type,
    agg_recipients.recipients,
    pulse_channel.details AS recipient_external,
    pulse.parameters
   FROM ((public.pulse
     LEFT JOIN public.pulse_channel ON ((pulse.id = pulse_channel.pulse_id)))
     LEFT JOIN agg_recipients ON ((pulse_channel.id = agg_recipients.pulse_channel_id)))
  WHERE (pulse.alert_condition IS NULL);


ALTER VIEW public.v_subscriptions OWNER TO aman;

--
-- Name: v_tables; Type: VIEW; Schema: public; Owner: aman
--

CREATE VIEW public.v_tables AS
 SELECT id AS entity_id,
    ('table_'::text || id) AS entity_qualified_id,
    created_at,
    updated_at,
    name,
    display_name,
    description,
    active,
    db_id AS database_id,
    schema,
    is_upload
   FROM public.metabase_table;


ALTER VIEW public.v_tables OWNER TO aman;

--
-- Name: v_tasks; Type: VIEW; Schema: public; Owner: aman
--

CREATE VIEW public.v_tasks AS
 SELECT id,
    task,
    ('database_'::text || db_id) AS database_qualified_id,
    started_at,
    ended_at,
    ((duration)::double precision / (1000)::double precision) AS duration_seconds,
    task_details AS details
   FROM public.task_history;


ALTER VIEW public.v_tasks OWNER TO aman;

--
-- Name: v_users; Type: VIEW; Schema: public; Owner: aman
--

CREATE VIEW public.v_users AS
 SELECT core_user.id AS user_id,
    ('user_'::text || core_user.id) AS entity_qualified_id,
    core_user.type,
        CASE
            WHEN ((core_user.type)::text = 'api-key'::text) THEN NULL::public.citext
            ELSE core_user.email
        END AS email,
    core_user.first_name,
    core_user.last_name,
    (((core_user.first_name)::text || ' '::text) || (core_user.last_name)::text) AS full_name,
    core_user.date_joined,
    core_user.last_login,
    core_user.updated_at,
    core_user.is_superuser AS is_admin,
    core_user.is_active,
    core_user.sso_source,
    core_user.locale
   FROM public.core_user
UNION
 SELECT 0 AS user_id,
    'user_0'::text AS entity_qualified_id,
    'anonymous'::character varying AS type,
    NULL::public.citext AS email,
    'External'::character varying AS first_name,
    'User'::character varying AS last_name,
    'External User'::text AS full_name,
    NULL::timestamp with time zone AS date_joined,
    NULL::timestamp with time zone AS last_login,
    NULL::timestamp with time zone AS updated_at,
    false AS is_admin,
    NULL::boolean AS is_active,
    NULL::character varying AS sso_source,
    NULL::character varying AS locale;


ALTER VIEW public.v_users OWNER TO aman;

--
-- Name: view_log; Type: TABLE; Schema: public; Owner: aman
--

CREATE TABLE public.view_log (
    id integer NOT NULL,
    user_id integer,
    model character varying(16) NOT NULL,
    model_id integer NOT NULL,
    "timestamp" timestamp with time zone NOT NULL,
    metadata text,
    has_access boolean,
    context character varying(32)
);


ALTER TABLE public.view_log OWNER TO aman;

--
-- Name: COLUMN view_log.has_access; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.view_log.has_access IS 'Whether the user who initiated the view had read access to the item being viewed.';


--
-- Name: COLUMN view_log.context; Type: COMMENT; Schema: public; Owner: aman
--

COMMENT ON COLUMN public.view_log.context IS 'The context of the view, can be collection, question, or dashboard. Only for cards.';


--
-- Name: v_view_log; Type: VIEW; Schema: public; Owner: aman
--

CREATE VIEW public.v_view_log AS
 SELECT id,
    "timestamp",
    COALESCE(user_id, 0) AS user_id,
    model AS entity_type,
    model_id AS entity_id,
    (((model)::text || '_'::text) || model_id) AS entity_qualified_id
   FROM public.view_log;


ALTER VIEW public.v_view_log OWNER TO aman;

--
-- Name: view_log_id_seq; Type: SEQUENCE; Schema: public; Owner: aman
--

ALTER TABLE public.view_log ALTER COLUMN id ADD GENERATED BY DEFAULT AS IDENTITY (
    SEQUENCE NAME public.view_log_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1
);


--
-- Data for Name: action; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.action (id, created_at, updated_at, type, model_id, name, description, parameters, parameter_mappings, visualization_settings, public_uuid, made_public_by_id, creator_id, archived, entity_id) FROM stdin;
\.


--
-- Data for Name: api_key; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.api_key (id, user_id, key, key_prefix, creator_id, created_at, updated_at, name, updated_by_id, scope) FROM stdin;
\.


--
-- Data for Name: application_permissions_revision; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.application_permissions_revision (id, before, after, user_id, created_at, remark) FROM stdin;
\.


--
-- Data for Name: audit_log; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.audit_log (id, topic, "timestamp", end_timestamp, user_id, model, model_id, details) FROM stdin;
\.


--
-- Data for Name: bookmark_ordering; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.bookmark_ordering (id, user_id, type, item_id, ordering) FROM stdin;
\.


--
-- Data for Name: cache_config; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.cache_config (id, model, model_id, created_at, updated_at, strategy, config, state, invalidated_at, next_run_at) FROM stdin;
\.


--
-- Data for Name: card_bookmark; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.card_bookmark (id, user_id, card_id, created_at) FROM stdin;
\.


--
-- Data for Name: card_label; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.card_label (id, card_id, label_id) FROM stdin;
\.


--
-- Data for Name: cloud_migration; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.cloud_migration (id, external_id, upload_url, state, progress, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: collection; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.collection (id, name, description, archived, location, personal_owner_id, slug, namespace, authority_level, entity_id, created_at, type, is_sample) FROM stdin;
1	Examples	\N	f	/	\N	examples	\N	\N	53YGAg4EE6MC76nxx-f5f	2024-06-22 08:19:17.839482+00	\N	t
2	Amanpreet Singh's Personal Collection	\N	f	/	1	amanpreet_singh_s_personal_collection	\N	\N	pIx3LpuPxJeV0_U93fwEP	2024-06-22 08:25:41.424253+00	\N	f
\.


--
-- Data for Name: collection_bookmark; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.collection_bookmark (id, user_id, collection_id, created_at) FROM stdin;
\.


--
-- Data for Name: collection_permission_graph_revision; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.collection_permission_graph_revision (id, before, after, user_id, created_at, remark) FROM stdin;
\.


--
-- Data for Name: connection_impersonations; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.connection_impersonations (id, db_id, group_id, attribute) FROM stdin;
\.


--
-- Data for Name: core_session; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.core_session (id, user_id, created_at, anti_csrf_token) FROM stdin;
aff0e644-0b3a-4483-8134-897f53e96131	1	2024-06-22 08:25:50.578832+00	\N
\.


--
-- Data for Name: core_user; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.core_user (id, email, first_name, last_name, password, password_salt, date_joined, last_login, is_superuser, is_active, reset_token, reset_triggered, is_qbnewb, login_attributes, updated_at, sso_source, locale, is_datasetnewb, settings, type, entity_id) FROM stdin;
13371338	internal@metabase.com	Metabase	Internal	$2a$10$fp5n6.DSzeuV8xMGYQWaPO6D6N1Dm0M4jD6oePE2xq2wnNskl.YEO	de2c4024-e5fa-45f2-84af-7d44aed9e1d6	2024-06-22 08:19:16.518337+00	\N	f	f	\N	\N	t	\N	\N	\N	\N	t	\N	internal	\N
1	aman@punjab.com	Amanpreet	Singh	$2a$10$1r0mAnpyhZ0Z6bsE2tjplO2ID.Gwy2ny5VGQgGIHCR5X9RS.BEmby	7eadc72d-d34e-4373-8e5c-a17ab4fa5590	2024-06-22 08:25:03.366674+00	2024-06-22 08:25:50.587559+00	t	t	\N	\N	t	\N	2024-06-22 08:25:50.587559+00	\N	\N	t	{"last-acknowledged-version":"v0.50.6"}	personal	Tq8fPkfolnfdHphriI_Ct
\.


--
-- Data for Name: dashboard_bookmark; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.dashboard_bookmark (id, user_id, dashboard_id, created_at) FROM stdin;
\.


--
-- Data for Name: dashboard_favorite; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.dashboard_favorite (id, user_id, dashboard_id) FROM stdin;
\.


--
-- Data for Name: dashboard_tab; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.dashboard_tab (id, dashboard_id, name, "position", entity_id, created_at, updated_at) FROM stdin;
1	1	Overview	0	PS7GW5IfD3ov0haogAPLu	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00
2	1	Portfolio performance	1	z5IymDhXzRY2kF79LxQIN	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00
3	1	Demographics	2	TKyJ0onLUPOuZgVfwZfU1	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00
\.


--
-- Data for Name: dashboardcard_series; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.dashboardcard_series (id, dashboardcard_id, card_id, "position") FROM stdin;
\.


--
-- Data for Name: data_permissions; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.data_permissions (id, group_id, perm_type, db_id, schema_name, table_id, perm_value) FROM stdin;
1	1	perms/view-data	1	\N	\N	unrestricted
2	1	perms/create-queries	1	\N	\N	query-builder-and-native
3	1	perms/download-results	1	\N	\N	one-million-rows
4	1	perms/manage-table-metadata	1	\N	\N	no
5	1	perms/manage-database	1	\N	\N	no
7	1	perms/create-queries	2	\N	\N	query-builder-and-native
8	1	perms/view-data	2	\N	\N	unrestricted
9	1	perms/download-results	2	\N	\N	one-million-rows
10	1	perms/manage-table-metadata	2	\N	\N	no
11	1	perms/manage-database	2	\N	\N	no
\.


--
-- Data for Name: databasechangelog; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.databasechangelog (id, author, filename, dateexecuted, orderexecuted, exectype, md5sum, description, comments, tag, liquibase, contexts, labels, deployment_id) FROM stdin;
v49.00-007	johnswanson	migrations/001_update_migrations.yaml	2024-06-22 08:19:13.693617	216	EXECUTED	9:638a7394870315abd52742b739ee49ff	sql	Set the `type` of the internal user	\N	4.26.0	\N	\N	9044340460
v00.00-000	qnkhuat	migrations/001_update_migrations.yaml	2024-06-22 08:19:02.605362	1	EXECUTED	9:e346841fa1cd9d142e1237b37fdc8a20	sqlFile path=initialization/metabase_postgres.sql; sqlFile path=initialization/metabase_mysql.sql; sqlFile path=initialization/metabase_h2.sql	Initialze metabase	\N	4.26.0	\N	\N	9044340460
v45.00-001	snoe	migrations/001_update_migrations.yaml	2024-06-22 08:19:02.625104	2	EXECUTED	9:15c13a8aa3fdc72ef0c54f4cccfc39e1	createTable tableName=action	Added 0.44.0 - writeback	\N	4.26.0	\N	\N	9044340460
v45.00-002	snoe	migrations/001_update_migrations.yaml	2024-06-22 08:19:02.63515	3	EXECUTED	9:d2c9f50f5a29947a07e4808957d63ab6	createTable tableName=query_action	Added 0.44.0 - writeback	\N	4.26.0	\N	\N	9044340460
v45.00-003	snoe	migrations/001_update_migrations.yaml	2024-06-22 08:19:02.645029	4	EXECUTED	9:dafaaf7c9f0efbf92670ea93c001f7a1	addPrimaryKey constraintName=pk_query_action, tableName=query_action	Added 0.44.0 - writeback	\N	4.26.0	\N	\N	9044340460
v45.00-011	snoe	migrations/001_update_migrations.yaml	2024-06-22 08:19:02.652952	5	EXECUTED	9:c539f152aa1c2287c5b602c7a395f9e8	addColumn tableName=report_card	Added 0.44.0 - writeback	\N	4.26.0	\N	\N	9044340460
v45.00-012	snoe	migrations/001_update_migrations.yaml	2024-06-22 08:19:02.667927	6	EXECUTED	9:b872219f47d12ec80db8f1731be4ce94	createTable tableName=http_action	Added 0.44.0 - writeback	\N	4.26.0	\N	\N	9044340460
v45.00-013	snoe	migrations/001_update_migrations.yaml	2024-06-22 08:19:02.678211	7	EXECUTED	9:54c5d6a9659b7ae62e8c42f60f9620d2	addPrimaryKey constraintName=pk_http_action, tableName=http_action	Added 0.44.0 - writeback	\N	4.26.0	\N	\N	9044340460
v45.00-022	snoe	migrations/001_update_migrations.yaml	2024-06-22 08:19:02.697825	8	EXECUTED	9:06e8a3ba8c4c5e5cf8e6aa3f35081cac	createTable tableName=app	Added 0.45.0 - add app container	\N	4.26.0	\N	\N	9044340460
v45.00-023	snoe	migrations/001_update_migrations.yaml	2024-06-22 08:19:02.753885	9	EXECUTED	9:082d296233ee6dbbf2871b9d93c3a6a4	addForeignKeyConstraint baseTableName=app, constraintName=fk_app_ref_dashboard_id, referencedTableName=report_dashboard	Added 0.45.0 - add app container	\N	4.26.0	\N	\N	9044340460
v45.00-025	metamben	migrations/001_update_migrations.yaml	2024-06-22 08:19:02.762126	10	EXECUTED	9:2214b0d71acc8a8cf90781a2aca98664	addColumn tableName=report_dashboard	Added 0.45.0 - mark app pages	\N	4.26.0	\N	\N	9044340460
v45.00-026	snoe	migrations/001_update_migrations.yaml	2024-06-22 08:19:02.770376	11	EXECUTED	9:9f2ce2d2d79d0dce365ddf3464d1f648	addColumn tableName=report_dashboardcard	Added 0.45.0 - apps add action_id to report_dashboardcard	\N	4.26.0	\N	\N	9044340460
v45.00-027	snoe	migrations/001_update_migrations.yaml	2024-06-22 08:19:02.778598	12	EXECUTED	9:73718f7b7c3fb4ef30f71dc6e6170528	addForeignKeyConstraint baseTableName=report_dashboardcard, constraintName=fk_report_dashboardcard_ref_action_id, referencedTableName=action	Added 0.45.0 - apps add fk for action_id to report_dashboardcard	\N	4.26.0	\N	\N	9044340460
v45.00-028	camsaul	migrations/001_update_migrations.yaml	2024-06-22 08:19:02.786745	13	EXECUTED	9:883315d70f3fc10b07858aa0e48ed9da	renameColumn newColumnName=size_x, oldColumnName=sizeX, tableName=report_dashboardcard	Added 0.45.0 -- rename DashboardCard sizeX to size_x. See https://github.com/metabase/metabase/issues/16344	\N	4.26.0	\N	\N	9044340460
v45.00-029	camsaul	migrations/001_update_migrations.yaml	2024-06-22 08:19:02.856456	14	EXECUTED	9:fe3b8aca811ef5b541f922a83c7ded8c	renameColumn newColumnName=size_y, oldColumnName=sizeY, tableName=report_dashboardcard	Added 0.45.0 -- rename DashboardCard size_y to size_y. See https://github.com/metabase/metabase/issues/16344	\N	4.26.0	\N	\N	9044340460
v45.00-030	camsaul	migrations/001_update_migrations.yaml	2024-06-22 08:19:02.86398	15	EXECUTED	9:63da8f2f82baf396ad30f3fd451c501d	addDefaultValue columnName=size_x, tableName=report_dashboardcard	Added 0.45.0 -- add default value to DashboardCard size_x -- this was previously done by Toucan	\N	4.26.0	\N	\N	9044340460
v45.00-031	camsaul	migrations/001_update_migrations.yaml	2024-06-22 08:19:02.871972	16	EXECUTED	9:3628c1c692bec0f7258ea983b18340b5	addDefaultValue columnName=size_y, tableName=report_dashboardcard	Added 0.45.0 -- add default value to DashboardCard size_y -- this was previously done by Toucan	\N	4.26.0	\N	\N	9044340460
v45.00-032	camsaul	migrations/001_update_migrations.yaml	2024-06-22 08:19:02.881273	17	EXECUTED	9:7a6f0320210b82c5eafe836ea98f477d	addDefaultValue columnName=created_at, tableName=report_dashboardcard	Added 0.45.0 -- add default value for DashboardCard created_at (Postgres/H2)	\N	4.26.0	\N	\N	9044340460
v45.00-033	camsaul	migrations/001_update_migrations.yaml	2024-06-22 08:19:02.886104	18	MARK_RAN	9:b93dab321e4a1fcaeb74d12b52867fd4	sql	Added 0.45.0 -- add default value for DashboardCard created_at (MySQL/MariaDB)	\N	4.26.0	\N	\N	9044340460
v45.00-034	camsaul	migrations/001_update_migrations.yaml	2024-06-22 08:19:02.892796	19	EXECUTED	9:39c965f9dc521d2a1e196f645d11de9e	addDefaultValue columnName=updated_at, tableName=report_dashboardcard	Added 0.45.0 -- add default value for DashboardCard updated_at (Postgres/H2)	\N	4.26.0	\N	\N	9044340460
v45.00-035	camsaul	migrations/001_update_migrations.yaml	2024-06-22 08:19:02.896619	20	MARK_RAN	9:207984e21c44681a6135592cab3d0f3f	sql	Added 0.45.0 -- add default value for DashboardCard updated_at (MySQL/MariaDB)	\N	4.26.0	\N	\N	9044340460
v45.00-036	snoe	migrations/001_update_migrations.yaml	2024-06-22 08:19:02.965649	21	EXECUTED	9:8632f7a046e1094399c517b05f0feea3	createTable tableName=model_action	Added 0.45.0 - add model action table	\N	4.26.0	\N	\N	9044340460
v45.00-037	snoe	migrations/001_update_migrations.yaml	2024-06-22 08:19:02.976009	22	EXECUTED	9:6600702fb0cc9dcd8628fca8df9c0b39	addUniqueConstraint constraintName=unique_model_action_card_id_slug, tableName=model_action	Added 0.45.0 - model action	\N	4.26.0	\N	\N	9044340460
v45.00-038	camsaul	migrations/001_update_migrations.yaml	2024-06-22 08:19:02.982901	23	EXECUTED	9:5ae52c12861e5eec6a2a9a8c5a442826	addDefaultValue columnName=created_at, tableName=metabase_database	Added 0.45.0 -- add default value for Database created_at (Postgres/H2)	\N	4.26.0	\N	\N	9044340460
v45.00-039	camsaul	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.043837	24	MARK_RAN	9:7b331f47a0260218275c58fa21fdcc60	sql	Added 0.45.0 -- add default value for Database created_at (MySQL/MariaDB)	\N	4.26.0	\N	\N	9044340460
v45.00-040	camsaul	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.055686	25	EXECUTED	9:6d6fb2c7cd62868b9878a951c8596cec	addDefaultValue columnName=updated_at, tableName=metabase_database	Added 0.45.0 -- add default value for Database updated_at (Postgres/H2)	\N	4.26.0	\N	\N	9044340460
v45.00-041	camsaul	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.059409	26	MARK_RAN	9:84d479472b81809b8067f72b28934954	sql	Added 0.45.0 -- add default value for Database updated_at (MySQL/MariaDB)	\N	4.26.0	\N	\N	9044340460
v45.00-042	camsaul	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.064493	27	EXECUTED	9:f600f2e052bf44938d081165c8d87364	sql	Added 0.45.0 -- add default value for Database with NULL details	\N	4.26.0	\N	\N	9044340460
v45.00-043	camsaul	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.071249	28	EXECUTED	9:f97ef9506c96084958075bb7c8d67b37	addNotNullConstraint columnName=details, tableName=metabase_database	Added 0.45.0 -- make Database details NOT NULL	\N	4.26.0	\N	\N	9044340460
v45.00-044	metamben	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.084952	29	EXECUTED	9:3a2bdd2e615d4577394828afda7fe61b	createTable tableName=app_permission_graph_revision	Added 0.45.0 -- create app permission graph revision table	\N	4.26.0	\N	\N	9044340460
v45.00-048	camsaul	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.153383	30	EXECUTED	9:20fea59b307a381f506485c07b3434f0	addColumn tableName=collection	Added 0.45.0 -- add created_at to Collection	\N	4.26.0	\N	\N	9044340460
v45.00-049	camsaul	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.158658	31	EXECUTED	9:525cdee9190fec1c22f431bbc4c33165	sql; sql; sql	Added 0.45.0 -- set Collection.created_at to User.date_joined for Personal Collections	\N	4.26.0	\N	\N	9044340460
v45.00-050	camsaul	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.163792	32	EXECUTED	9:6a28bbf7cf34c8730e1d3151fcc5b090	sql; sql; sql	Added 0.45.0 -- seed Collection.created_at with value of oldest item for non-Personal Collections	\N	4.26.0	\N	\N	9044340460
v45.00-051	qnkhuat	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.166218	33	MARK_RAN	9:21a53bba289feab23d8a322e28c1e281	modifyDataType columnName=after, tableName=collection_permission_graph_revision	Added 0.45.0 - modify type of collection_permission_graph_revision.after from text to text on mysql,mariadb	\N	4.26.0	\N	\N	9044340460
v45.00-052	qnkhuat	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.169916	34	MARK_RAN	9:6cbe88994cc644d01d336851484f7d04	modifyDataType columnName=before, tableName=collection_permission_graph_revision	Added 0.45.0 - modify type of collection_permission_graph_revision.before from text to text on mysql,mariadb	\N	4.26.0	\N	\N	9044340460
v45.00-053	qnkhuat	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.173949	35	MARK_RAN	9:ea7fc60a2091ce7a78266b6b2a0f91ef	modifyDataType columnName=remark, tableName=collection_permission_graph_revision	Added 0.45.0 - modify type of collection_permission_graph_revision.remark from text to text on mysql,mariadb	\N	4.26.0	\N	\N	9044340460
v45.00-054	qnkhuat	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.24539	36	MARK_RAN	9:b150b811b511b9a2f9a063f1e72e144b	modifyDataType columnName=after, tableName=permissions_revision	Added 0.45.0 - modify type of permissions_revision.after from text to text on mysql,mariadb	\N	4.26.0	\N	\N	9044340460
v45.00-055	qnkhuat	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.252418	37	MARK_RAN	9:b0e45cae16c5f3ba5b3becf3b767dced	modifyDataType columnName=before, tableName=permissions_revision	Added 0.45.0 - modify type of permissions_revision.before from text to text on mysql,mariadb	\N	4.26.0	\N	\N	9044340460
v45.00-056	qnkhuat	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.256095	38	MARK_RAN	9:2b7e01dd0c5f8d720cbc346526f712ae	modifyDataType columnName=remark, tableName=permissions_revision	Added 0.45.0 - modify type of permissions_revision.remark from text to text on mysql,mariadb	\N	4.26.0	\N	\N	9044340460
v45.00-057	qnkhuat	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.259703	39	MARK_RAN	9:d1125daee3e40f7a316c59c6b7a0fb1b	modifyDataType columnName=value, tableName=secret	Added 0.45.0 - modify type of secret.value from blob to longblob on mysql,mariadb	\N	4.26.0	\N	\N	9044340460
v46.00-000	snoe	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.269902	40	EXECUTED	9:ebf6161a0fd64634e0032cf8c44e2c64	createTable tableName=implicit_action	Added 0.46.0 - Unify action representation	\N	4.26.0	\N	\N	9044340460
v46.00-001	snoe	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.280286	41	EXECUTED	9:9a3a543cd836c34d8131b6c929061425	addColumn tableName=action	Added 0.46.0 - Unify action representation	\N	4.26.0	\N	\N	9044340460
v46.00-002	snoe	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.346503	42	EXECUTED	9:a352e46d75236605148308fc7e95cfe6	addColumn tableName=action	Added 0.46.0 - Unify action representation	\N	4.26.0	\N	\N	9044340460
v46.00-003	snoe	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.353899	43	EXECUTED	9:c9ae88de84e869dbda6681c543d9701c	addColumn tableName=action	Added 0.46.0 - Unify action representation	\N	4.26.0	\N	\N	9044340460
v46.00-004	snoe	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.364396	44	EXECUTED	9:f7b94fa036afd26110610156dc0054ea	addColumn tableName=action	Added 0.46.0 - Unify action representation	\N	4.26.0	\N	\N	9044340460
v46.00-005	snoe	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.377622	45	EXECUTED	9:60be619c1c562da167a15ea8cc7421b3	addColumn tableName=action	Added 0.46.0 - Unify action representation	\N	4.26.0	\N	\N	9044340460
v46.00-006	snoe	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.391296	46	EXECUTED	9:3c24b56f156b5891db5ee8a3b0a68195	addColumn tableName=action	Added 0.46.0 - Unify action representation	\N	4.26.0	\N	\N	9044340460
v46.00-007	snoe	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.452937	47	EXECUTED	9:604733020cb8d94b0b21728393b9227a	addForeignKeyConstraint baseTableName=action, constraintName=fk_action_model_id, referencedTableName=report_card	Added 0.46.0 - Unify action representation	\N	4.26.0	\N	\N	9044340460
v46.00-008	snoe	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.460665	48	EXECUTED	9:24adc5d36405f95e479157afa6ac7090	addColumn tableName=query_action	Added 0.46.0 - Unify action representation	\N	4.26.0	\N	\N	9044340460
v46.00-009	snoe	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.47054	49	EXECUTED	9:a54fda21eeecdf97eab5c25e7ad616be	addColumn tableName=query_action	Added 0.46.0 - Unify action representation	\N	4.26.0	\N	\N	9044340460
v46.00-010	snoe	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.478119	50	EXECUTED	9:17661d9a15cea0e046aa22bbe05f8ca9	addForeignKeyConstraint baseTableName=query_action, constraintName=fk_query_action_database_id, referencedTableName=metabase_database	Added 0.46.0 - Unify action representation	\N	4.26.0	\N	\N	9044340460
v46.00-011	snoe	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.488343	51	EXECUTED	9:704a13bc4f66dba37ce7fbf130c2d207	sql; sql; sql	Added 0.46.0 - Unify action representation	\N	4.26.0	\N	\N	9044340460
v46.00-012	snoe	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.55242	52	EXECUTED	9:f5e43052660cb6fd61dca7385a41690f	dropNotNullConstraint columnName=card_id, tableName=query_action	Added 0.46.0 - Unify action representation	\N	4.26.0	\N	\N	9044340460
v46.00-013	snoe	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.563669	53	EXECUTED	9:12bf37b0e6732f2f09122f76cb8b59b5	sql	Added 0.46.0 - Unify action representation	\N	4.26.0	\N	\N	9044340460
v46.00-014	snoe	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.57015	54	EXECUTED	9:f5856fab23b69dbc669c354d5cb14d36	dropForeignKeyConstraint baseTableName=query_action, constraintName=fk_query_action_ref_card_id	Added 0.46.0 - Unify action representation	\N	4.26.0	\N	\N	9044340460
v46.00-015	snoe	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.577631	55	EXECUTED	9:1ef231990b57b1f8496070f5f63f4579	dropColumn columnName=card_id, tableName=query_action	Added 0.46.0 - Unify action representation	\N	4.26.0	\N	\N	9044340460
v46.00-016	snoe	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.58225	56	EXECUTED	9:cf5ae825070bb05fd046b0a22f299403	sql	Added 0.46.0 - Unify action representation	\N	4.26.0	\N	\N	9044340460
v46.00-017	snoe	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.652995	57	EXECUTED	9:7428779ccdce61097d0fe5e761e05b18	dropColumn columnName=name, tableName=http_action	Added 0.46.0 - Unify action representation	\N	4.26.0	\N	\N	9044340460
v46.00-018	snoe	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.660391	58	EXECUTED	9:1ae9ab13080fc2dd4df6ccdeb28baf72	dropColumn columnName=description, tableName=http_action	Added 0.46.0 - Unify action representation	\N	4.26.0	\N	\N	9044340460
v46.00-019	snoe	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.668356	59	EXECUTED	9:ffa18e1fb9a06d0a673a1d6aab1edfbc	dropColumn columnName=is_write, tableName=report_card	Added 0.46.0 - Unify action representation	\N	4.26.0	\N	\N	9044340460
v46.00-020	snoe	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.676015	60	EXECUTED	9:f96413cb1260923de7be41fd0a665543	addNotNullConstraint columnName=database_id, tableName=query_action	Added 0.46.0 - Unify action representation	\N	4.26.0	\N	\N	9044340460
v46.00-021	snoe	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.683806	61	EXECUTED	9:92176e007f82994418edd21c13ea648d	addNotNullConstraint columnName=dataset_query, tableName=query_action	Added 0.46.0 - Unify action representation	\N	4.26.0	\N	\N	9044340460
v46.00-022	snoe	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.691685	62	EXECUTED	9:92176e007f82994418edd21c13ea648d	addNotNullConstraint columnName=dataset_query, tableName=query_action	Added 0.46.0 - Unify action representation	\N	4.26.0	\N	\N	9044340460
v46.00-023	snoe	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.699585	63	EXECUTED	9:072dea52b8b04cbb2829741c8523d768	addNotNullConstraint columnName=model_id, tableName=action	Added 0.46.0 - Unify action representation	\N	4.26.0	\N	\N	9044340460
v46.00-024	snoe	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.755892	64	EXECUTED	9:92798789cb9756596896f74e7055988d	addNotNullConstraint columnName=name, tableName=action	Added 0.46.0 - Unify action representation	\N	4.26.0	\N	\N	9044340460
v46.00-025	snoe	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.763787	65	EXECUTED	9:cedeeebf30271f749b77f79609b197ba	dropTable tableName=model_action	Added 0.46.0 - Unify action representation	\N	4.26.0	\N	\N	9044340460
v46.00-026	metamben	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.771377	66	EXECUTED	9:2466f8ea4308426e2dfaa1a0f89d2807	addColumn tableName=metabase_database	Added 0.46.0 -- add field for tracking DBMS versions	\N	4.26.0	\N	\N	9044340460
v46.00-027	snoe	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.778916	67	EXECUTED	9:1450a0a51d88a368393d2202a8e194fd	addColumn tableName=metabase_fieldvalues	Added 0.46.0 -- add last_used_at to FieldValues	\N	4.26.0	\N	\N	9044340460
v46.00-028	tsmacdonald	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.790784	68	EXECUTED	9:cc1b9daf6ce06234decfadd3a2629a80	createTable tableName=parameter_card	Added 0.46.0 -- Join table connecting cards to dashboards/cards's parameters that need custom filter values from the card	\N	4.26.0	\N	\N	9044340460
v46.00-029	camsaul	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.798319	69	EXECUTED	9:c95d07d656ec9c2d06e9d83ab14e5170	dropUniqueConstraint constraintName=unique_dimension_field_id_name, tableName=dimension	Make Dimension <=> Field a 1t1 relationship. Drop unique constraint on field_id + name. (1/3)	\N	4.26.0	\N	\N	9044340460
v46.00-030	camsaul	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.851121	70	EXECUTED	9:3342b8ad5e8160325705f4df619c1416	sql	Make Dimension <=> Field a 1t1 relationship. Delete duplicate entries. (2/3)	\N	4.26.0	\N	\N	9044340460
v46.00-031	camsaul	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.859767	71	EXECUTED	9:5822fde920ae9c048d75e951674c6570	addUniqueConstraint constraintName=unique_dimension_field_id, tableName=dimension	Make Dimension <=> Field a 1t1 relationship. Add unique constraint on field_id. (3/3)	\N	4.26.0	\N	\N	9044340460
v46.00-032	tsmacdonald	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.870201	72	EXECUTED	9:2e882eb92ce78edf23d417dcbc7c4f03	addUniqueConstraint constraintName=unique_parameterized_object_card_parameter, tableName=parameter_card	Added 0.46.0 -- Unique parameter_card	\N	4.26.0	\N	\N	9044340460
v46.00-033	tsmacdonald	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.880565	73	EXECUTED	9:2ca24580b18e8a0736cc1c0b463d43a6	createIndex indexName=idx_parameter_card_parameterized_object_id, tableName=parameter_card	Added 0.46.0 -- parameter_card index on connected object	\N	4.26.0	\N	\N	9044340460
v46.00-034	tsmacdonald	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.891078	74	EXECUTED	9:654f10ff8db9163f27f87fc7da686912	createIndex indexName=idx_parameter_card_card_id, tableName=parameter_card	Added 0.46.0 -- parameter_card index on connected card	\N	4.26.0	\N	\N	9044340460
v46.00-035	tsmacdonald	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.899215	75	EXECUTED	9:54d48bc5d1b60faa10ba000c3409c973	addForeignKeyConstraint baseTableName=parameter_card, constraintName=fk_parameter_card_ref_card_id, referencedTableName=report_card	Added 0.46.0 - parameter_card.card_id foreign key	\N	4.26.0	\N	\N	9044340460
v46.00-036	metamben	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.907941	76	EXECUTED	9:5e155b0c1edc895bddf0e5de96def85e	dropTable tableName=app_permission_graph_revision	App containers are removed in 0.46.0	\N	4.26.0	\N	\N	9044340460
v46.00-037	metamben	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.916614	77	EXECUTED	9:ba61742e04377180572fb011d5ad3263	dropColumn columnName=is_app_page, tableName=report_dashboard	App pages are removed in 0.46.0	\N	4.26.0	\N	\N	9044340460
v46.00-038	metamben	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.926727	78	EXECUTED	9:a45d7cfa18714cd1d536bd17b927e3ac	dropTable tableName=app	App containers are removed in 0.46.0	\N	4.26.0	\N	\N	9044340460
v46.00-039	qnkhuat	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.93812	79	EXECUTED	9:a850f5255ad1f2abb52d4a7c93fa68a0	addColumn tableName=parameter_card	Added 0.46.0 - add entity_id to parameter_card	\N	4.26.0	\N	\N	9044340460
v46.00-089	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:04.758728	107	EXECUTED	9:7a7761ef80cae24914abb65143fb0503	sql	Added 0.46.5 -- remove orphaned entries in `sandboxes`	\N	4.26.0	\N	\N	9044340460
v46.00-040	tsmacdonald	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.94559	80	EXECUTED	9:400fe00b6e522509e4949a71cc12e196	addDefaultValue columnName=size_x, tableName=report_dashboardcard	Added 0.46.0 -- Bump default dashcard size to 4x4	\N	4.26.0	\N	\N	9044340460
v46.00-041	tsmacdonald	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.957017	81	EXECUTED	9:0f7635baf2acb131464865f2a94a17ed	addDefaultValue columnName=size_y, tableName=report_dashboardcard	Added 0.46.0 -- Bump default dashcard size to 4x4	\N	4.26.0	\N	\N	9044340460
v46.00-042	tsmacdonald	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.965712	82	EXECUTED	9:6b99f953adbc8c3f536c7d729ef780a2	createIndex indexName=idx_query_execution_executor_id, tableName=query_execution	Added 0.46.0 -- index query_execution.executor_id	\N	4.26.0	\N	\N	9044340460
v46.00-043	tsmacdonald	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.975412	83	EXECUTED	9:6fd04070a9c19ac34c15a5cbe27def47	createIndex indexName=idx_query_execution_context, tableName=query_execution	Added 0.46.0 -- index query_execution.context	\N	4.26.0	\N	\N	9044340460
v46.00-045	calherries	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.985134	84	EXECUTED	9:1f61e64d7ce15c315ff530ef8a2dc7f3	addColumn tableName=action	Added 0.46.0 -- add public_uuid to action.	\N	4.26.0	\N	\N	9044340460
v46.00-051	calherries	migrations/001_update_migrations.yaml	2024-06-22 08:19:03.993054	85	EXECUTED	9:7f54a26b4edf1aa4bc86b95f9c316dea	dropDefaultValue columnName=row, tableName=report_dashboardcard	Added 0.46.0 -- drop defaults for dashcard's position and size	\N	4.26.0	\N	\N	9044340460
v46.00-052	calherries	migrations/001_update_migrations.yaml	2024-06-22 08:19:04.001197	86	EXECUTED	9:4fc7e45d19ab28e91bd27325d23afd83	dropDefaultValue columnName=col, tableName=report_dashboardcard	Added 0.46.0 -- drop defaults for dashcard's position and size	\N	4.26.0	\N	\N	9044340460
v46.00-053	calherries	migrations/001_update_migrations.yaml	2024-06-22 08:19:04.008962	87	EXECUTED	9:2d2836a1c76d1c5db400219364aa070e	dropDefaultValue columnName=size_x, tableName=report_dashboardcard	Added 0.46.0 -- drop defaults for dashcard's position and size	\N	4.26.0	\N	\N	9044340460
v46.00-054	calherries	migrations/001_update_migrations.yaml	2024-06-22 08:19:04.016791	88	EXECUTED	9:60b7a066e6ca1826efcb8b562b33f331	dropDefaultValue columnName=size_y, tableName=report_dashboardcard	Added 0.46.0 -- drop defaults for dashcard's position and size	\N	4.26.0	\N	\N	9044340460
v46.00-055	calherries	migrations/001_update_migrations.yaml	2024-06-22 08:19:04.024805	89	EXECUTED	9:5b583b1b5eac915914f40a27163263b3	addColumn tableName=action	Added 0.46.0 -- add made_public_by_id	\N	4.26.0	\N	\N	9044340460
v46.00-056	calherries	migrations/001_update_migrations.yaml	2024-06-22 08:19:04.037987	90	EXECUTED	9:9e37acc40f5685ec3992dff91bbfa618	createIndex indexName=idx_action_public_uuid, tableName=action	Added 0.46.0 -- add public_uuid and made_public_by_id to action. public_uuid is indexed	\N	4.26.0	\N	\N	9044340460
v46.00-057	dpsutton	migrations/001_update_migrations.yaml	2024-06-22 08:19:04.048322	91	EXECUTED	9:86b2d39951f748b94abd2e1c7b68144a	modifyDataType columnName=parameter_id, tableName=parameter_card	Added 0.46.0 -- parameter_card.parameter_id long enough to hold a uuid	\N	4.26.0	\N	\N	9044340460
v46.00-058	calherries	migrations/001_update_migrations.yaml	2024-06-22 08:19:04.057681	92	EXECUTED	9:7458ddbf194a385219df05900d78185d	addForeignKeyConstraint baseTableName=action, constraintName=fk_action_made_public_by_id, referencedTableName=core_user	Added 0.46.0 -- add FK constraint for action.made_public_by_id with core_user.id	\N	4.26.0	\N	\N	9044340460
v46.00-059	tsmacdonald	migrations/001_update_migrations.yaml	2024-06-22 08:19:04.06852	93	EXECUTED	9:55a077fb646639308046b80c3f64267e	addColumn tableName=action	Added 0.46.0 -- add actions.creator_id	\N	4.26.0	\N	\N	9044340460
v46.00-060	tsmacdonald	migrations/001_update_migrations.yaml	2024-06-22 08:19:04.086904	94	EXECUTED	9:a0726ea2ef7c354af8b60a8ab37d24bd	createIndex indexName=idx_action_creator_id, tableName=action	Added 0.46.0 -- action.creator_id index	\N	4.26.0	\N	\N	9044340460
v46.00-061	tsmacdonald	migrations/001_update_migrations.yaml	2024-06-22 08:19:04.156261	95	EXECUTED	9:7497cc10ea1fdad211179b36d53bde6a	addForeignKeyConstraint baseTableName=action, constraintName=fk_action_creator_id, referencedTableName=core_user	Added 0.46.0 -- action.creator_id index	\N	4.26.0	\N	\N	9044340460
v46.00-062	tsmacdonald	migrations/001_update_migrations.yaml	2024-06-22 08:19:04.164447	96	EXECUTED	9:35cc67731bc19abd498bcdbb0aeb688e	addColumn tableName=action	Added 0.46.0 -- add actions.archived	\N	4.26.0	\N	\N	9044340460
v46.00-064	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:04.172538	97	EXECUTED	9:a73357bb088af23043336f048172a1f3	renameTable newTableName=sandboxes, oldTableName=group_table_access_policy	Added 0.46.0 -- rename `group_table_access_policy` to `sandboxes`	\N	4.26.0	\N	\N	9044340460
v46.00-065	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:04.180486	98	EXECUTED	9:b5b5d61c3e6f7daa528acce1a52b6f75	addColumn tableName=sandboxes	Added 0.46.0 -- add `permission_id` to `sandboxes`	\N	4.26.0	\N	\N	9044340460
v46.00-070	calherries	migrations/001_update_migrations.yaml	2024-06-22 08:19:04.256607	99	EXECUTED	9:e52174e49aa14e61eb922ba200cfc002	addColumn tableName=action	Added 0.46.0 - add entity_id column to action	\N	4.26.0	\N	\N	9044340460
v46.00-074	metamben	migrations/001_update_migrations.yaml	2024-06-22 08:19:04.264859	100	EXECUTED	9:c2c5004951b49617a624ba2cf79fb617	modifyDataType columnName=updated_at, tableName=report_card	Added 0.46.0 -- increase precision of updated_at of report_card	\N	4.26.0	\N	\N	9044340460
v46.00-079	john-metabase	migrations/001_update_migrations.yaml	2024-06-22 08:19:04.269972	101	EXECUTED	9:e5ce98dd7ac26fb102db98625e10dfab	sql	Added 0.46.0 -- migrates Databases using deprecated and removed presto driver to presto-jdbc	\N	4.26.0	\N	\N	9044340460
v46.00-080	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:04.450001	102	EXECUTED	9:742aaa27012538f89770095551808dbc	customChange	Migrate data permission paths from v1 to v2 (splitting them into separate data and query permissions)	\N	4.26.0	\N	\N	9044340460
v46.00-084	qnkhuat	migrations/001_update_migrations.yaml	2024-06-22 08:19:04.456395	103	EXECUTED	9:ec26e53892d1293ef822005b0e2d5d0d	dropForeignKeyConstraint baseTableName=action, constraintName=fk_action_model_id	Added 0.46.0 - CASCADE delete for action.model_id	\N	4.26.0	\N	\N	9044340460
v46.00-085	qnkhuat	migrations/001_update_migrations.yaml	2024-06-22 08:19:04.46427	104	EXECUTED	9:013d55806a3c819a2e94ff2d5cb71df2	addForeignKeyConstraint baseTableName=action, constraintName=fk_action_model_id, referencedTableName=report_card	Added 0.46.0 - CASCADE delete for action.model_id	\N	4.26.0	\N	\N	9044340460
v46.00-086	calherries	migrations/001_update_migrations.yaml	2024-06-22 08:19:04.751434	105	EXECUTED	9:3029f7b7b204834ce65fc573378b3425	customChange	Added 0.46.0 - Delete the abandonment email task	\N	4.26.0	\N	\N	9044340460
v46.00-088	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:04.75565	106	EXECUTED	9:c38b799746664b436dc8c63f2c6214c6	sql	Added 0.46.5 -- backfill `permission_id` values in `sandboxes`. This is a fixed verison of v46.00-066 which has been removed, since it had a bug that blocked a customer from upgrading.	\N	4.26.0	\N	\N	9044340460
v46.00-090	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:10.571851	108	EXECUTED	9:d17eaed2e43c5589332682c53e4a6458	addForeignKeyConstraint baseTableName=sandboxes, constraintName=fk_sandboxes_ref_permissions, referencedTableName=permissions	Add foreign key constraint on sandboxes.permission_id	\N	4.26.0	\N	\N	9044340460
v47.00-001	calherries	migrations/001_update_migrations.yaml	2024-06-22 08:19:10.577218	109	EXECUTED	9:b655eb44da8863f3021bd71dc14e69b3	sql	Added 0.47.0 -- set base-type to type/JSON for JSON database-types for postgres and mysql	\N	4.26.0	\N	\N	9044340460
v47.00-002	calherries	migrations/001_update_migrations.yaml	2024-06-22 08:19:10.585229	110	EXECUTED	9:1983e3e2e005932513c2146770aa7f37	addColumn tableName=metabase_field	Added 0.47.0 - Add json_unfolding column to metabase_field	\N	4.26.0	\N	\N	9044340460
v47.00-003	calherries	migrations/001_update_migrations.yaml	2024-06-22 08:19:10.591128	111	EXECUTED	9:a3cdc062588a6c5c4459d48dc1b76519	sql	Added 0.47.0 - Populate metabase_field.json_unfolding based on base_type	\N	4.26.0	\N	\N	9044340460
v47.00-004	qnkhuat	migrations/001_update_migrations.yaml	2024-06-22 08:19:10.597807	112	EXECUTED	9:9457f62a9c6533da79c4881ca21d727f	addColumn tableName=metabase_field	Added 0.47.0 - Add auto_incremented to metabase_field	\N	4.26.0	\N	\N	9044340460
v47.00-005	winlost	migrations/001_update_migrations.yaml	2024-06-22 08:19:10.605656	113	EXECUTED	9:afc9309217305117e3b9e88018c5437e	addColumn tableName=report_dashboard	Added 0.47.0 - Add auto_apply_filters to dashboard	\N	4.26.0	\N	\N	9044340460
v47.00-006	qnkhuat	migrations/001_update_migrations.yaml	2024-06-22 08:19:10.661341	114	EXECUTED	9:98647900eba8f78d7bef0678415ebe2a	createTable tableName=dashboard_tab	Added 0.47.0 - Add dashboard_tab table	\N	4.26.0	\N	\N	9044340460
v47.00-007	qnkhuat	migrations/001_update_migrations.yaml	2024-06-22 08:19:10.669352	115	EXECUTED	9:7e6c0250df2957d89875cd0d14d271ee	addColumn tableName=report_dashboardcard	Added 0.47.0 -- add report_dashboardcard.dashboard_tab_id	\N	4.26.0	\N	\N	9044340460
v47.00-008	qnkhuat	migrations/001_update_migrations.yaml	2024-06-22 08:19:10.677059	116	EXECUTED	9:606ece18373d1e534db0714add1c41e8	addForeignKeyConstraint baseTableName=report_dashboardcard, constraintName=fk_report_dashboardcard_ref_dashboard_tab_id, referencedTableName=dashboard_tab	Added 0.47.0 -- add report_dashboardcard.dashboard_tab_id fk constraint	\N	4.26.0	\N	\N	9044340460
v47.00-009	qwef	migrations/001_update_migrations.yaml	2024-06-22 08:19:10.683418	117	EXECUTED	9:738a13b99f62269ff366f57c09652ebc	sql	Added 0.47.0 - Replace user google_auth and ldap_auth columns with sso_source values	\N	4.26.0	\N	\N	9044340460
v47.00-010	tsmacdonald	migrations/001_update_migrations.yaml	2024-06-22 08:19:10.750016	118	EXECUTED	9:5d9a509a79dedadff743042b6f82ddbe	modifyDataType columnName=name, tableName=metabase_table	Added 0.47.0 - Make metabase_table.name long enough for H2 names	\N	4.26.0	\N	\N	9044340460
v47.00-011	tsmacdonald	migrations/001_update_migrations.yaml	2024-06-22 08:19:10.756475	119	EXECUTED	9:7a6a06886b1428bf54d5105d1d4fcf0a	modifyDataType columnName=display_name, tableName=metabase_table	Added 0.47.0 - Make metabase_table.display_name long enough for H2 names	\N	4.26.0	\N	\N	9044340460
v47.00-012	qwef	migrations/001_update_migrations.yaml	2024-06-22 08:19:10.764382	120	EXECUTED	9:42ca291fec3fb0a7bdd74dd17d03339a	dropColumn columnName=google_auth, tableName=core_user	Added 0.47.0 - Replace user google_auth and ldap_auth columns with sso_source values	\N	4.26.0	\N	\N	9044340460
v47.00-013	qwef	migrations/001_update_migrations.yaml	2024-06-22 08:19:10.769106	121	EXECUTED	9:fcc22e7f3fd2f6e52739a8b9778f8e50	sql	Added 0.47.0 - Replace user google_auth and ldap_auth columns with sso_source values	\N	4.26.0	\N	\N	9044340460
v47.00-014	qwef	migrations/001_update_migrations.yaml	2024-06-22 08:19:10.777825	122	EXECUTED	9:e8309c65ac1f122f944e64d702409377	dropColumn columnName=ldap_auth, tableName=core_user	Added 0.47.0 - Replace user google_auth and ldap_auth columns with sso_source values	\N	4.26.0	\N	\N	9044340460
v47.00-015	escherize	migrations/001_update_migrations.yaml	2024-06-22 08:19:10.786148	123	EXECUTED	9:a05fe66edd2512f17a9fc7f9ff122669	addColumn tableName=metabase_database	added 0.47.0 - Add is_audit to metabase_database	\N	4.26.0	\N	\N	9044340460
v47.00-016	calherres	migrations/001_update_migrations.yaml	2024-06-22 08:19:10.802214	124	EXECUTED	9:dbd88ce575f7976114e5b8c7e0382a5a	customChange	Added 0.47.0 - Migrate the report_card.visualization_settings.column_settings field refs from legacy format	\N	4.26.0	\N	\N	9044340460
v47.00-018	dpsutton	migrations/001_update_migrations.yaml	2024-06-22 08:19:10.81493	125	EXECUTED	9:ec0895ba935e438f7cd104534b3b13f4	createTable tableName=model_index	Indexed Entities information table	\N	4.26.0	\N	\N	9044340460
v47.00-019	dpsutton	migrations/001_update_migrations.yaml	2024-06-22 08:19:10.855875	126	EXECUTED	9:d3d383af3f3c901efdd891a405c2639b	createTable tableName=model_index_value	Indexed Entities values table	\N	4.26.0	\N	\N	9044340460
v47.00-020	dpsutton	migrations/001_update_migrations.yaml	2024-06-22 08:19:10.865949	127	EXECUTED	9:96991ffa2c754f2f70d61928c7574e31	addUniqueConstraint constraintName=unique_model_index_value_model_index_id_model_pk, tableName=model_index_value	Add unique constraint on index_id and pk	\N	4.26.0	\N	\N	9044340460
v47.00-023	dpsutton	migrations/001_update_migrations.yaml	2024-06-22 08:19:10.875272	128	EXECUTED	9:ab9a7a030405fc844b36f10d77fef039	createIndex indexName=idx_model_index_model_id, tableName=model_index	Added 0.47.0 -- model_index index	\N	4.26.0	\N	\N	9044340460
v47.00-024	dpsutton	migrations/001_update_migrations.yaml	2024-06-22 08:19:10.883043	129	EXECUTED	9:ab4cd02348503fcaa7cedd49d11d217d	addForeignKeyConstraint baseTableName=model_index, constraintName=fk_model_index_model_id, referencedTableName=report_card	Added 0.47.0 -- model_index foriegn key to report_card	\N	4.26.0	\N	\N	9044340460
v47.00-025	dpsutton	migrations/001_update_migrations.yaml	2024-06-22 08:19:10.891433	130	EXECUTED	9:e002aab419bb327f68a09b604eb24fa5	addForeignKeyConstraint baseTableName=model_index_value, constraintName=fk_model_index_value_model_id, referencedTableName=model_index	Added 0.47.0 -- model_index_value foriegn key to model_index	\N	4.26.0	\N	\N	9044340460
v47.00-026	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:10.952647	131	EXECUTED	9:cd23ef8fb587021b6a7190811db35b90	createTable tableName=connection_impersonations	Added 0.47.0 - New table for connection impersonation policies	\N	4.26.0	\N	\N	9044340460
v47.00-027	calherries	migrations/001_update_migrations.yaml	2024-06-22 08:19:10.962027	132	EXECUTED	9:a40730a20a1f9ae345d15f9c2bfa443e	customChange	Added 0.47.0 - Migrate field_ref in report_card.result_metadata from legacy format	\N	4.26.0	\N	\N	9044340460
v47.00-028	calherries	migrations/001_update_migrations.yaml	2024-06-22 08:19:10.967549	133	EXECUTED	9:66ec973062109ac00d89c8bb1a867d1e	customChange	Added 0.47.0 - Add join-alias to the report_card.visualization_settings.column_settings field refs	\N	4.26.0	\N	\N	9044340460
v47.00-029	qnkhuat	migrations/001_update_migrations.yaml	2024-06-22 08:19:10.971922	134	EXECUTED	9:4e63a84adeeacdbe8ecf771e3f6cf65e	customChange	Added 0.47.0 - Stack cards vertically for dashboard with tabs on downgrade	\N	4.26.0	\N	\N	9044340460
v48.00-005	qnkhuat	migrations/001_update_migrations.yaml	2024-06-22 08:19:11.346945	161	EXECUTED	9:46dbd09ead2298f0e37d28f011b8986e	addColumn tableName=query_execution	Added 0.48.0 - Add query_execution.action_id	\N	4.26.0	\N	\N	9044340460
v47.00-030	escherize	migrations/001_update_migrations.yaml	2024-06-22 08:19:10.991015	135	EXECUTED	9:e6d3c56859e8c41a5c5f6b7ebcc489d7	addColumn tableName=collection	Added 0.47.0 - Type column for collections for instance-analytics	\N	4.26.0	\N	\N	9044340460
v47.00-031	qnkhuat	migrations/001_update_migrations.yaml	2024-06-22 08:19:11.051626	136	EXECUTED	9:6c0047955d2fc52cd25e4e4aabdc7143	sql; sql	Added 0.47.0 - migrate dashboard grid size from 18 to 24	\N	4.26.0	\N	\N	9044340460
v47.00-032	qnkhuat	migrations/001_update_migrations.yaml	2024-06-22 08:19:11.141541	137	EXECUTED	9:e9621439a79c5255a2cfe540c1e8df73	customChange	Added 0.47.0 - migrate dashboard grid size from 18 to 24 for revisions	\N	4.26.0	\N	\N	9044340460
v47.00-033	calherries	migrations/001_update_migrations.yaml	2024-06-22 08:19:11.152205	138	EXECUTED	9:ef48cabec16353dc09a6297d06e27a9f	customChange	Added 0.47.0 - Migrate field refs in visualization_settings.column_settings keys from legacy format	\N	4.26.0	\N	\N	9044340460
v47.00-034	calherries	migrations/001_update_migrations.yaml	2024-06-22 08:19:11.158792	139	EXECUTED	9:9032d7052ab38a476e7f83d9486e8609	customChange	Added 0.47.0 - Add join-alias to the visualization_settings.column_settings field refs in card revisions	\N	4.26.0	\N	\N	9044340460
v47.00-035	calherries	migrations/001_update_migrations.yaml	2024-06-22 08:19:11.16574	140	EXECUTED	9:57f53ab4caba81bea788f23487b6888a	dropForeignKeyConstraint baseTableName=implicit_action, constraintName=fk_implicit_action_action_id	Added 0.47.0 - Drop foreign key constraint on implicit_action.action_id	\N	4.26.0	\N	\N	9044340460
v47.00-036	calherries	migrations/001_update_migrations.yaml	2024-06-22 08:19:11.174306	141	EXECUTED	9:c45b0246a2068f7bed037662a25e52c0	addPrimaryKey constraintName=pk_implicit_action, tableName=implicit_action	Added 0.47.0 - Set primary key to action_id for implicit_action table	\N	4.26.0	\N	\N	9044340460
v47.00-037	calherries	migrations/001_update_migrations.yaml	2024-06-22 08:19:11.181489	142	EXECUTED	9:8a7a6c51a3f52acb8b69bf10782e3a6e	addForeignKeyConstraint baseTableName=implicit_action, constraintName=fk_implicit_action_action_id, referencedTableName=action	Added 0.47.0 - Add foreign key constraint on implicit_action.action_id	\N	4.26.0	\N	\N	9044340460
v47.00-043	calherres	migrations/001_update_migrations.yaml	2024-06-22 08:19:11.190951	143	EXECUTED	9:d7a6e633d99b064539cbc7563c7f1905	customChange	Added 0.47.0 - Migrate report_dashboardcard.visualization_settings.column_settings field refs from legacy format	\N	4.26.0	\N	\N	9044340460
v47.00-044	calherries	migrations/001_update_migrations.yaml	2024-06-22 08:19:11.198253	144	EXECUTED	9:905e1385202b3626c9e67e4645b44375	customChange	Added 0.47.0 - Add join-alias to the report_dashboardcard.visualization_settings.column_settings field refs	\N	4.26.0	\N	\N	9044340460
v47.00-045	calherres	migrations/001_update_migrations.yaml	2024-06-22 08:19:11.203866	145	EXECUTED	9:ad1660852fc3eacd2c35e164d3c97609	customChange	Added 0.47.0 - Migrate dashboard revision dashboard cards' visualization_settings.column_settings field refs from legacy format	\N	4.26.0	\N	\N	9044340460
v47.00-046	calherries	migrations/001_update_migrations.yaml	2024-06-22 08:19:11.209162	146	EXECUTED	9:e0cd8b50dc855ce3362b562dc80afbb0	customChange	Added 0.47.0 - Add join-alias to dashboard revision dashboard cards' visualization_settings.column_settings field refs	\N	4.26.0	\N	\N	9044340460
v47.00-050	tsmacdonald	migrations/001_update_migrations.yaml	2024-06-22 08:19:11.214763	147	EXECUTED	9:8e9aad6950c3b6f4799f51f0a7277457	addColumn tableName=metabase_table	Added 0.47.0 - table.is_upload	\N	4.26.0	\N	\N	9044340460
v47.00-051	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:11.221528	148	EXECUTED	9:ec914ed2458e64fdd021bc6ca19e85c1	dropForeignKeyConstraint baseTableName=connection_impersonations, constraintName=fk_conn_impersonation_db_id	Added 0.47.0 - Drop foreign key constraint on connection_impersonations.db_id	\N	4.26.0	\N	\N	9044340460
v47.00-052	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:11.228361	149	EXECUTED	9:9dc7bed69c542749be182f5d76233488	dropForeignKeyConstraint baseTableName=connection_impersonations, constraintName=fk_conn_impersonation_group_id	Added 0.47.0 - Drop foreign key constraint on connection_impersonations.group_id	\N	4.26.0	\N	\N	9044340460
v47.00-053	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:11.23665	150	EXECUTED	9:91a175d49e28b1653ac8a7bad2a8204f	createIndex indexName=idx_conn_impersonations_db_id, tableName=connection_impersonations	Added 0.47.0 -- connection_impersonations index for db_id column	\N	4.26.0	\N	\N	9044340460
v47.00-054	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:11.245156	151	EXECUTED	9:5245b1f06e5d1ae3a4c929f6cfcde9d4	createIndex indexName=idx_conn_impersonations_group_id, tableName=connection_impersonations	Added 0.47.0 -- connection_impersonations index for group_id column	\N	4.26.0	\N	\N	9044340460
v47.00-055	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:11.253419	152	EXECUTED	9:aef0d48122a09e653a9683d3b7074165	addUniqueConstraint constraintName=conn_impersonation_unique_group_id_db_id, tableName=connection_impersonations	Added 0.47.0 - unique constraint for connection impersonations	\N	4.26.0	\N	\N	9044340460
v47.00-056	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:11.260296	153	EXECUTED	9:bac4590e83ce8522e9fad8703aa88295	addForeignKeyConstraint baseTableName=connection_impersonations, constraintName=fk_conn_impersonation_db_id, referencedTableName=metabase_database	Added 0.47.0 - re-add foreign key constraint on connection_impersonations.db_id	\N	4.26.0	\N	\N	9044340460
v47.00-057	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:11.267077	154	EXECUTED	9:21e3b2ee3c1ae91f35c7c1fbd0f82dac	addForeignKeyConstraint baseTableName=connection_impersonations, constraintName=fk_conn_impersonation_group_id, referencedTableName=permissions_group	Added 0.47.0 - re-add foreign key constraint on connection_impersonations.group_id	\N	4.26.0	\N	\N	9044340460
v47.00-058	qnkhuat	migrations/001_update_migrations.yaml	2024-06-22 08:19:11.274896	155	EXECUTED	9:d268796eef28773e767c286904f50cee	dropColumn columnName=entity_id, tableName=parameter_card	Drop parameter_card.entity_id	\N	4.26.0	\N	\N	9044340460
v47.00-059	piranha	migrations/001_update_migrations.yaml	2024-06-22 08:19:11.281378	156	EXECUTED	9:087e632c348027aeb01bc49bb428f67a	dropNotNullConstraint columnName=entity_id, tableName=dashboard_tab	Drops not null from dashboard_tab.entity_id since it breaks drop-entity-ids command	\N	4.26.0	\N	\N	9044340460
v48.00-001	qnkhuat	migrations/001_update_migrations.yaml	2024-06-22 08:19:11.292033	157	EXECUTED	9:2dd9aa417cd83b00a59b5eb1deb9ae55	customChange	Added 0.47.0 - Migrate database.options to database.settings	\N	4.26.0	\N	\N	9044340460
v48.00-002	qnkhuat	migrations/001_update_migrations.yaml	2024-06-22 08:19:11.297805	158	EXECUTED	9:7a9312346a8d041d62726cd99d27e02b	dropColumn columnName=options, tableName=metabase_database	Added 0.47.0 - drop metabase_database.options	\N	4.26.0	\N	\N	9044340460
v48.00-003	qnkhuat	migrations/001_update_migrations.yaml	2024-06-22 08:19:11.315356	159	EXECUTED	9:b623ee408e368c53ef1865af15a5edac	dropTable tableName=computation_job_result	Added 0.48.0 - drop computation_job_result table	\N	4.26.0	\N	\N	9044340460
v48.00-004	qnkhuat	migrations/001_update_migrations.yaml	2024-06-22 08:19:11.332377	160	EXECUTED	9:e9bdc7a7ca09fe29d5b3479d9925bbd7	dropTable tableName=computation_job	Added 0.48.0 - drop computation_job table	\N	4.26.0	\N	\N	9044340460
v48.00-006	qnkhuat	migrations/001_update_migrations.yaml	2024-06-22 08:19:11.375535	162	EXECUTED	9:87cb5019893ddcb08ce06856825f1e6a	createIndex indexName=idx_query_execution_action_id, tableName=query_execution	Added 0.48.0 - Index query_execution.action_id	\N	4.26.0	\N	\N	9044340460
v48.00-007	qnkhuat	migrations/001_update_migrations.yaml	2024-06-22 08:19:11.400872	163	EXECUTED	9:efa5dcca04d3887bca0eafc95754b9ee	addColumn tableName=revision	Added 0.48.0 - Add revision.most_recent	\N	4.26.0	\N	\N	9044340460
v48.00-008	qnkhuat	migrations/001_update_migrations.yaml	2024-06-22 08:19:11.442897	164	EXECUTED	9:ad3f69be24960003ce545c16ac7922eb	sql; sql	Set revision.most_recent = true for latest revisions	\N	4.26.0	\N	\N	9044340460
v48.00-009	calherries	migrations/001_update_migrations.yaml	2024-06-22 08:19:11.462559	165	EXECUTED	9:4243d81311d2dc265d687e4b665215dc	createTable tableName=table_privileges	Added 0.48.0 - Create table_privileges table	\N	4.26.0	\N	\N	9044340460
v48.00-010	qnkhuat	migrations/001_update_migrations.yaml	2024-06-22 08:19:11.475285	166	MARK_RAN	9:79908afdc739c3e9dbdddd1e29a20e5d	sql	Remove ON UPDATE for revision.timestamp on mysql, mariadb	\N	4.26.0	\N	\N	9044340460
v48.00-011	qnkhuat	migrations/001_update_migrations.yaml	2024-06-22 08:19:11.507286	167	EXECUTED	9:2da4950057f3a912cd044bd95a189087	createIndex indexName=idx_revision_most_recent, tableName=revision	Index revision.most_recent	\N	4.26.0	\N	\N	9044340460
v48.00-013	qnkhuat	migrations/001_update_migrations.yaml	2024-06-22 08:19:11.774346	168	EXECUTED	9:9478b5f507bc34df51aaa4d719fc1c37	sql	Index unindexed FKs for postgres	\N	4.26.0	\N	\N	9044340460
v48.00-014	calherries	migrations/001_update_migrations.yaml	2024-06-22 08:19:11.807378	169	EXECUTED	9:f22e8ba6b47eaa27344835a6f5069c7f	createIndex indexName=idx_table_privileges_table_id, tableName=table_privileges	Added 0.48.0 - Create table_privileges.table_id index	\N	4.26.0	\N	\N	9044340460
v48.00-015	calherries	migrations/001_update_migrations.yaml	2024-06-22 08:19:11.843016	170	EXECUTED	9:7ce7965ef94006d568d2298c4e5e007c	createIndex indexName=idx_table_privileges_role, tableName=table_privileges	Added 0.48.0 - Create table_privileges.role index	\N	4.26.0	\N	\N	9044340460
v48.00-016	calherries	migrations/001_update_migrations.yaml	2024-06-22 08:19:11.876097	171	EXECUTED	9:99bac12bbd1ff64b35f0c48f5f9bf824	modifyDataType columnName=slug, tableName=collection	Added 0.48.0 - Change the type of collection.slug to varchar(510)	\N	4.26.0	\N	\N	9044340460
v48.00-018	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:11.915932	172	EXECUTED	9:7e190e044ea31c175d382aa6eb2b80c4	createTable tableName=recent_views	Add new recent_views table	\N	4.26.0	\N	\N	9044340460
v48.00-019	nemanjaglumac	migrations/001_update_migrations.yaml	2024-06-22 08:19:11.943925	173	EXECUTED	9:3af8327e2dab1b5344292966c4e1d8e0	dropColumn columnName=color, tableName=collection	Collection color is removed in 0.48.0	\N	4.26.0	\N	\N	9044340460
v48.00-020	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:11.97975	174	EXECUTED	9:53aa607fc4b962053bffbca4e4283391	createIndex indexName=idx_recent_views_user_id, tableName=recent_views	Added 0.48.0 - Create recent_views.user_id index	\N	4.26.0	\N	\N	9044340460
v48.00-021	piranha	migrations/001_update_migrations.yaml	2024-06-22 08:19:12.01483	175	EXECUTED	9:8fd8b97bfbae458ea7eb91e117cabc33	addColumn tableName=report_card	Cards store Metabase version used to create them	\N	4.26.0	\N	\N	9044340460
v48.00-022	johnswanson	migrations/001_update_migrations.yaml	2024-06-22 08:19:12.038787	176	EXECUTED	9:0f826a4d0f9f512ca63fb667409ad999	customChange	Migrate migrate-click-through to a custom migration	\N	4.26.0	\N	\N	9044340460
v48.00-023	piranha	migrations/001_update_migrations.yaml	2024-06-22 08:19:12.072235	177	EXECUTED	9:6ee76f86f68f7a2120cf8160cb34fd99	customChange	Data migration migrate-remove-admin-from-group-mapping-if-needed	\N	4.26.0	\N	\N	9044340460
v48.00-024	piranha	migrations/001_update_migrations.yaml	2024-06-22 08:19:12.093421	178	EXECUTED	9:87ed320552869766323c1d12f7969b17	dropTable tableName=data_migrations	All data migrations were transferred to custom_migrations!	\N	4.26.0	\N	\N	9044340460
v48.00-025	piranha	migrations/001_update_migrations.yaml	2024-06-22 08:19:12.158783	179	EXECUTED	9:0e54628ce18964128e827ad05ab2a448	addColumn tableName=revision	Revisions store Metabase version used to create them	\N	4.26.0	\N	\N	9044340460
v48.00-026	lbrdnk	migrations/001_update_migrations.yaml	2024-06-22 08:19:12.174769	180	EXECUTED	9:d1ce706fe25d39767a77089171859856	update tableName=metabase_field	Set semantic_type with value type/Number to null (#18754)	\N	4.26.0	\N	\N	9044340460
v48.00-027	qnkhuat	migrations/001_update_migrations.yaml	2024-06-22 08:19:12.190184	181	EXECUTED	9:0a7ed49c904abc7110aa09324b49f106	sql	No op migration	\N	4.26.0	\N	\N	9044340460
v48.00-028	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:12.239044	182	EXECUTED	9:9340f3026bae444e9e98bc714a17a2d0	createTable tableName=audit_log	Add new audit_log table	\N	4.26.0	\N	\N	9044340460
v48.00-029	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:12.270461	183	EXECUTED	9:0ea006a4b23969abbb4574f1e7fc35c0	sqlFile path=instance_analytics_views/audit_log/v1/postgres-audit_log.sql; sqlFile path=instance_analytics_views/audit_log/v1/mysql-audit_log.sql; sqlFile path=instance_analytics_views/audit_log/v1/h2-audit_log.sql	Added 0.48.0 - new view v_audit_log	\N	4.26.0	\N	\N	9044340460
v48.00-030	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:12.302461	184	EXECUTED	9:d59b20e1ab52d390d4bd9577b9a54439	sqlFile path=instance_analytics_views/content/v1/postgres-content.sql; sqlFile path=instance_analytics_views/content/v1/mysql-content.sql; sqlFile path=instance_analytics_views/content/v1/h2-content.sql	Added 0.48.0 - new view v_content	\N	4.26.0	\N	\N	9044340460
v48.00-031	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:12.358828	185	EXECUTED	9:4493c1a05ff71a8cf629c2e469e63b29	sqlFile path=instance_analytics_views/dashboardcard/v1/dashboardcard.sql	Added 0.48.0 - new view v_dashboardcard	\N	4.26.0	\N	\N	9044340460
v48.00-032	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:12.391553	186	EXECUTED	9:c73d37997ac7c32b730b761bd1a7ecc0	sqlFile path=instance_analytics_views/group_members/v1/group_members.sql	Added 0.48.0 - new view v_group_members	\N	4.26.0	\N	\N	9044340460
v48.00-033	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:12.424927	187	EXECUTED	9:5b0b30fbf86a695d93e04de3131294c4	sqlFile path=instance_analytics_views/subscriptions/v1/postgres-subscriptions.sql; sqlFile path=instance_analytics_views/subscriptions/v1/mysql-subscriptions.sql; sqlFile path=instance_analytics_views/subscriptions/v1/h2-subscriptions.sql	Added 0.48.0 - new view v_subscriptions for postgres	\N	4.26.0	\N	\N	9044340460
v48.00-034	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:12.459163	188	EXECUTED	9:0e64a610c2a763eb8366d8619b60b6fa	sqlFile path=instance_analytics_views/users/v1/postgres-users.sql; sqlFile path=instance_analytics_views/users/v1/mysql-users.sql; sqlFile path=instance_analytics_views/users/v1/h2-users.sql	Added 0.48.0 - new view v_users	\N	4.26.0	\N	\N	9044340460
v49.00-006	johnswanson	migrations/001_update_migrations.yaml	2024-06-22 08:19:13.677216	215	EXECUTED	9:46cfb5a9bab696272d0211d7f34870ef	createIndex indexName=idx_api_key_user_id, tableName=api_key	Add an index on `api_key.user_id`	\N	4.26.0	\N	\N	9044340460
v48.00-035	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:12.49856	189	EXECUTED	9:54569b1677c151e2fd052cef8c9d69cc	sqlFile path=instance_analytics_views/alerts/v1/postgres-alerts.sql; sqlFile path=instance_analytics_views/alerts/v1/mysql-alerts.sql; sqlFile path=instance_analytics_views/alerts/v1/h2-alerts.sql	Added 0.48.0 - new view v_alerts	\N	4.26.0	\N	\N	9044340460
v48.00-036	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:12.560712	190	EXECUTED	9:4a4b8b9d1d537618546664e041121858	sqlFile path=instance_analytics_views/databases/v1/databases.sql	Added 0.48.0 - new view v_databases	\N	4.26.0	\N	\N	9044340460
v48.00-037	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:12.59745	191	EXECUTED	9:a304624882c0f8bd2fd4241786376989	sqlFile path=instance_analytics_views/fields/v1/postgres-fields.sql; sqlFile path=instance_analytics_views/fields/v1/mysql-fields.sql; sqlFile path=instance_analytics_views/fields/v1/h2-fields.sql	Added 0.48.0 - new view v_fields	\N	4.26.0	\N	\N	9044340460
v48.00-038	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:12.674962	192	EXECUTED	9:a305f2ec4549f7d5afb8cd2318e87e4b	sqlFile path=instance_analytics_views/query_log/v1/postgres-query_log.sql; sqlFile path=instance_analytics_views/query_log/v1/mysql-query_log.sql; sqlFile path=instance_analytics_views/query_log/v1/h2-query_log.sql	Added 0.48.0 - new view v_query_log	\N	4.26.0	\N	\N	9044340460
v48.00-039	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:12.759852	193	EXECUTED	9:71d9d867d7798a728baf858802d6c255	sqlFile path=instance_analytics_views/tables/v1/postgres-tables.sql; sqlFile path=instance_analytics_views/tables/v1/mysql-tables.sql; sqlFile path=instance_analytics_views/tables/v1/h2-tables.sql	Added 0.48.0 - new view v_tables	\N	4.26.0	\N	\N	9044340460
v48.00-040	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:12.800419	194	EXECUTED	9:8e4eba04e03bb4edb326f7d4389e9dbb	sqlFile path=instance_analytics_views/view_log/v1/postgres-view_log.sql; sqlFile path=instance_analytics_views/view_log/v1/mysql-view_log.sql; sqlFile path=instance_analytics_views/view_log/v1/h2-view_log.sql	Added 0.48.0 - new view v_view_log	\N	4.26.0	\N	\N	9044340460
v48.00-045	qwef	migrations/001_update_migrations.yaml	2024-06-22 08:19:12.876416	195	EXECUTED	9:3a8a23ac9c57ddb43351c8a0d57edbb7	addColumn tableName=query_execution	Added 0.48.0 - add is_sandboxed to query_execution	\N	4.26.0	\N	\N	9044340460
v48.00-046	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:12.989013	196	EXECUTED	9:6015a45ac465e5430f7482280427b182	sqlFile path=instance_analytics_views/indexes/v1/postgres-indexes.sql; sqlFile path=instance_analytics_views/indexes/v1/mysql-indexes.sql; sqlFile path=instance_analytics_views/indexes/v1/mariadb-indexes.sql; sqlFile path=instance_analytics_views/...	Added 0.48.0 - new indexes to optimize audit v2 queries	\N	4.26.0	\N	\N	9044340460
v48.00-047	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:13.059928	197	EXECUTED	9:83099e323661cb585639e5ee8d48d7bf	dropForeignKeyConstraint baseTableName=recent_views, constraintName=fk_recent_views_ref_user_id	Drop foreign key on recent_views so that it can be recreated with onDelete policy	\N	4.26.0	\N	\N	9044340460
v48.00-048	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:13.094385	198	EXECUTED	9:1ccb1544d732d0c6944ee1c4ac998801	addForeignKeyConstraint baseTableName=recent_views, constraintName=fk_recent_views_ref_user_id, referencedTableName=core_user	Add foreign key on recent_views with onDelete CASCADE	\N	4.26.0	\N	\N	9044340460
v48.00-049	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:13.150141	199	EXECUTED	9:0064304265cd53bb1965e6956bbe410e	sql; sql; sql	Migrate data from activity to audit_log	\N	4.26.0	\N	\N	9044340460
v48.00-050	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:13.159414	200	EXECUTED	9:d41d8cd98f00b204e9800998ecf8427e	empty	Added 0.48.0 - no-op migration to remove audit DB and collection on downgrade	\N	4.26.0	\N	\N	9044340460
v48.00-051	calherries	migrations/001_update_migrations.yaml	2024-06-22 08:19:13.181473	201	EXECUTED	9:ed306ef18d2546ccbe84f621fed21f1a	sql; sql	Migrate metabase_field when the fk target field is inactive	\N	4.26.0	\N	\N	9044340460
v48.00-053	johnswanson	migrations/001_update_migrations.yaml	2024-06-22 08:19:13.217931	202	EXECUTED	9:7997563ebc2f57b53a1f19679d89b3b7	modifyDataType columnName=model, tableName=activity	Increase length of `activity.model` to fit longer model names	\N	4.26.0	\N	\N	9044340460
v48.00-054	escherize	migrations/001_update_migrations.yaml	2024-06-22 08:19:13.233158	203	EXECUTED	9:d41d8cd98f00b204e9800998ecf8427e	empty	Added 0.48.0 - no-op migration to remove Internal Metabase User on downgrade	\N	4.26.0	\N	\N	9044340460
v48.00-055	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:13.265789	204	EXECUTED	9:6c1a8ab293774009e006e0c629513fe5	sqlFile path=instance_analytics_views/tasks/v1/postgres-tasks.sql; sqlFile path=instance_analytics_views/tasks/v1/mysql-tasks.sql; sqlFile path=instance_analytics_views/tasks/v1/h2-tasks.sql	Added 0.48.0 - new view v_tasks	\N	4.26.0	\N	\N	9044340460
v48.00-056	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:13.297615	205	EXECUTED	9:cd66e7282ad1d92ae91d2694906f8762	addColumn tableName=view_log	Adjust view_log schema for Audit Log v2	\N	4.26.0	\N	\N	9044340460
v48.00-057	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:13.329792	206	EXECUTED	9:24967224f6153a884f8f588772dfaf7d	addColumn tableName=view_log	Adjust view_log schema for Audit Log v2	\N	4.26.0	\N	\N	9044340460
v48.00-059	qwef	migrations/001_update_migrations.yaml	2024-06-22 08:19:13.346872	207	EXECUTED	9:aac112e35358d50dd099737f32e491bc	sql	Update the namespace of any audit collections that are already loaded	\N	4.26.0	\N	\N	9044340460
v48.00-060	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:13.38805	208	EXECUTED	9:97dc864df6ae5bca6e14ff4b1a455030	createIndex indexName=idx_task_history_started_at, tableName=task_history	Added 0.48.0 - task_history.started_at	\N	4.26.0	\N	\N	9044340460
v48.00-061	piranha	migrations/001_update_migrations.yaml	2024-06-22 08:19:13.420899	209	EXECUTED	9:5afa172dba1d6093b94356463912628b	addColumn tableName=query_execution	Adds query_execution.cache_hash -> query_cache.query_hash	\N	4.26.0	\N	\N	9044340460
v48.00-067	qnkhuat	migrations/001_update_migrations.yaml	2024-06-22 08:19:13.462117	210	EXECUTED	9:77853099d1da534596765591803a5d4c	addUniqueConstraint constraintName=idx_databasechangelog_id_author_filename, tableName=databasechangelog	Add unique constraint idx_databasechangelog_id_author_filename	\N	4.26.0	\N	\N	9044340460
v49.00-000	qnkhuat	migrations/001_update_migrations.yaml	2024-06-22 08:19:13.479573	211	EXECUTED	9:a6eaedb204bd70b999a7b7ed7524b904	sql	Remove leagcy pulses	\N	4.26.0	\N	\N	9044340460
v49.00-003	johnswanson	migrations/001_update_migrations.yaml	2024-06-22 08:19:13.529058	212	EXECUTED	9:6c49190e041265b22255a447302236b8	addColumn tableName=core_user	Add a `type` to users	\N	4.26.0	\N	\N	9044340460
v49.00-004	johnswanson	migrations/001_update_migrations.yaml	2024-06-22 08:19:13.591711	213	EXECUTED	9:d3d62eb2ff5b27040c5e4d1f5a263b4c	createTable tableName=api_key	Add a table for API keys	\N	4.26.0	\N	\N	9044340460
v49.00-005	johnswanson	migrations/001_update_migrations.yaml	2024-06-22 08:19:13.636491	214	EXECUTED	9:790e75675265f53ca11f32994b9d1d12	createIndex indexName=idx_api_key_created_by, tableName=api_key	Add an index on `api_key.created_by`	\N	4.26.0	\N	\N	9044340460
v49.00-008	qnkhuat	migrations/001_update_migrations.yaml	2024-06-22 08:19:13.720948	217	EXECUTED	9:381669dee234455e2d6e975afb9e95b0	addColumn tableName=metabase_field	Add metabase_field.database_indexed	\N	4.26.0	\N	\N	9044340460
v49.00-009	adam-james	migrations/001_update_migrations.yaml	2024-06-22 08:19:13.739978	218	EXECUTED	9:20b6791390f852c195f25bfcfb9a77a2	sql; sql	Migrate pulse_card.include_csv to 'true' when the joined card.display is 'table'	\N	4.26.0	\N	\N	9044340460
v49.00-010	johnswanson	migrations/001_update_migrations.yaml	2024-06-22 08:19:13.780566	219	EXECUTED	9:c618fafea8561e66c97ac5d9a56106f3	addColumn tableName=api_key	Add a name to API Keys	\N	4.26.0	\N	\N	9044340460
v49.00-011	qnkhuat	migrations/001_update_migrations.yaml	2024-06-22 08:19:13.813655	220	EXECUTED	9:9dae0a3606d63cc6c47ec94f181d203c	addColumn tableName=metabase_table	Add metabase_table.database_require_filter	\N	4.26.0	\N	\N	9044340460
v49.00-012	qnkhuat	migrations/001_update_migrations.yaml	2024-06-22 08:19:13.845981	221	EXECUTED	9:129b215b73fec358adbd3cf9ba478515	addColumn tableName=metabase_field	Add metabase_field.database_partitioned	\N	4.26.0	\N	\N	9044340460
v49.00-013	johnswanson	migrations/001_update_migrations.yaml	2024-06-22 08:19:13.877148	222	EXECUTED	9:fd99fac43bba1d45dfa8f17cd635e55b	addColumn tableName=api_key	Add `api_key.updated_by_id`	\N	4.26.0	\N	\N	9044340460
v49.00-014	johnswanson	migrations/001_update_migrations.yaml	2024-06-22 08:19:13.918389	223	EXECUTED	9:ab6ff0fa8e1088222e01466c38295012	createIndex indexName=idx_api_key_updated_by_id, tableName=api_key	Add an index on `api_key.updated_by_id`	\N	4.26.0	\N	\N	9044340460
v49.00-015	johnswanson	migrations/001_update_migrations.yaml	2024-06-22 08:19:13.946542	224	EXECUTED	9:85e3063c4b008044fe6729789af4d2d0	renameColumn newColumnName=creator_id, oldColumnName=created_by, tableName=api_key	Rename `created_by` to `creator_id`	\N	4.26.0	\N	\N	9044340460
v49.00-017	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:13.963393	225	MARK_RAN	9:add3e47d69b3c7c7a407a186bba06564	addNotNullConstraint columnName=archived, tableName=action	Add NOT NULL constraint to action.archived	\N	4.26.0	\N	\N	9044340460
v49.00-018	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:13.978779	226	MARK_RAN	9:dfa7799dd05f439f01c646f35810d37a	addDefaultValue columnName=archived, tableName=action	Add default value to action.archived	\N	4.26.0	\N	\N	9044340460
v49.00-020	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:13.995001	227	MARK_RAN	9:3ed2c593b8e803a416caa3f66c791c18	addNotNullConstraint columnName=json_unfolding, tableName=metabase_field	Add NOT NULL constraint to metabase_field.json_unfolding	\N	4.26.0	\N	\N	9044340460
v49.00-021	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:14.008264	228	MARK_RAN	9:eda3e041b4def2f0c9188b131330a743	addDefaultValue columnName=json_unfolding, tableName=metabase_field	Add default value to metabase_field.json_unfolding	\N	4.26.0	\N	\N	9044340460
v49.00-023	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:14.025002	229	MARK_RAN	9:8a8187fb01a8e5e2e23be4d64a068f54	addNotNullConstraint columnName=database_is_auto_increment, tableName=metabase_field	Add NOT NULL constraint to metabase_field.database_is_auto_increment	\N	4.26.0	\N	\N	9044340460
v49.00-024	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:14.042479	230	MARK_RAN	9:438a77b956692c7e3703d96b913e5566	addDefaultValue columnName=database_is_auto_increment, tableName=metabase_field	Add default value to metabase_field.database_is_auto_increment	\N	4.26.0	\N	\N	9044340460
v49.00-026	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:14.057127	231	MARK_RAN	9:7432b8bccc366fdeb6611b7899140ee7	addNotNullConstraint columnName=auto_apply_filters, tableName=report_dashboard	Add NOT NULL constraint to report_dashboard.auto_apply_filters	\N	4.26.0	\N	\N	9044340460
v49.00-027	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:14.070736	232	MARK_RAN	9:e8e27cec1e1cb5801ddfca828e3404a2	addDefaultValue columnName=auto_apply_filters, tableName=report_dashboard	Add default value to report_dashboard.auto_apply_filters	\N	4.26.0	\N	\N	9044340460
v49.00-029	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:14.087847	233	MARK_RAN	9:69f4b29b93422535b22df7ceb4c9dfb5	addNotNullConstraint columnName=is_audit, tableName=metabase_database	Add NOT NULL constraint to metabase_database.is_audit	\N	4.26.0	\N	\N	9044340460
v49.00-030	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:14.105761	234	MARK_RAN	9:93140751149cec3fdf7a186e6bac564a	addDefaultValue columnName=is_audit, tableName=metabase_database	Add default value to metabase_database.is_audit	\N	4.26.0	\N	\N	9044340460
v49.00-032	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:14.116286	235	MARK_RAN	9:bd9098a5361f7ab1aa631615e39e7eea	addNotNullConstraint columnName=is_upload, tableName=metabase_table	Add NOT NULL constraint to metabase_table.is_upload	\N	4.26.0	\N	\N	9044340460
v49.00-033	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:14.128777	236	MARK_RAN	9:75f535cee2ac99c5eabfd8d23007cec5	addDefaultValue columnName=is_upload, tableName=metabase_table	Add default value to metabase_table.is_upload	\N	4.26.0	\N	\N	9044340460
v49.00-036	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:14.144085	237	MARK_RAN	9:e2f3bc795300e8cba1ce4f594624eb98	addNotNullConstraint columnName=most_recent, tableName=revision	Add NOT NULL constraint to revision.most_recent	\N	4.26.0	\N	\N	9044340460
v49.00-037	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:14.16296	238	MARK_RAN	9:ddaf7704e565f42c3599823018f7f0cd	addDefaultValue columnName=most_recent, tableName=revision	Add default value to revision.most_recent	\N	4.26.0	\N	\N	9044340460
v49.00-039	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:14.175968	239	MARK_RAN	9:e085f675a04674caa7eefae9555faa98	addNotNullConstraint columnName=select, tableName=table_privileges	Add NOT NULL constraint to table_privileges.select	\N	4.26.0	\N	\N	9044340460
v49.00-040	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:14.191169	240	MARK_RAN	9:a429c76f53ca1f6e40cb97e10f5bbb13	addDefaultValue columnName=select, tableName=table_privileges	Add default value to table_privileges.select	\N	4.26.0	\N	\N	9044340460
v49.00-042	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:14.206319	241	MARK_RAN	9:a2a53332a9e16f2510537f27694f177f	addNotNullConstraint columnName=update, tableName=table_privileges	Add NOT NULL constraint to table_privileges.update	\N	4.26.0	\N	\N	9044340460
v49.00-043	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:14.224997	242	MARK_RAN	9:1a18955a2d01bb8fb1b7edfc74f5d976	addDefaultValue columnName=update, tableName=table_privileges	Add default value to table_privileges.update	\N	4.26.0	\N	\N	9044340460
v49.00-045	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:14.240896	243	MARK_RAN	9:9ecb48a04626476919f49615614a1166	addNotNullConstraint columnName=insert, tableName=table_privileges	Add NOT NULL constraint to table_privileges.insert	\N	4.26.0	\N	\N	9044340460
v49.00-046	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:14.254305	244	MARK_RAN	9:32c3a206904e45e4d6b3d639a2477d4e	addDefaultValue columnName=insert, tableName=table_privileges	Add default value to table_privileges.insert	\N	4.26.0	\N	\N	9044340460
v49.00-048	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:14.272123	245	MARK_RAN	9:8f313a8ba8b16a19466d75d93df683f2	addNotNullConstraint columnName=delete, tableName=table_privileges	Add NOT NULL constraint to table_privileges.delete	\N	4.26.0	\N	\N	9044340460
v49.00-049	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:14.303054	246	MARK_RAN	9:fa6e7cafdbbf7880ddd8eef9b1cd33c9	addDefaultValue columnName=delete, tableName=table_privileges	Add default value to table_privileges.delete	\N	4.26.0	\N	\N	9044340460
v49.00-059	qnkhuat	migrations/001_update_migrations.yaml	2024-06-22 08:19:14.353643	247	EXECUTED	9:4046506923db48fd7a7c11021b58a4b5	customChange	Added 0.49.0 - unify type of time columns	\N	4.26.0	\N	\N	9044340460
v50.2024-01-04T13:52:51	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:15.207556	273	EXECUTED	9:19f8b6614c4fe95ff71b42830785df04	createTable tableName=data_permissions	Data permissions table	\N	4.26.0	\N	\N	9044340460
v49.2023-01-24T12:00:00	piranha	migrations/001_update_migrations.yaml	2024-06-22 08:19:14.399321	248	EXECUTED	9:6f381d427bef9022f3ac00618af2112b	createIndex indexName=idx_field_name_lower, tableName=metabase_field	This significantly speeds up api.database/autocomplete-fields query and is an improvement for issue 30588. H2 does not support this: https://github.com/h2database/h2database/issues/3535 Mariadb does not support this. Mysql says it does, but report...	\N	4.26.0	\N	\N	9044340460
v49.2024-01-22T11:50:00	qnkhuat	migrations/001_update_migrations.yaml	2024-06-22 08:19:14.432937	249	EXECUTED	9:643a028e650ad7fcd6789a702249a179	addColumn tableName=report_card	Add `report_card.type`	\N	4.26.0	\N	\N	9044340460
v49.2024-01-22T11:51:00	qnkhuat	migrations/001_update_migrations.yaml	2024-06-22 08:19:14.4491	250	EXECUTED	9:6e53593c40c63e38d6c1ffd0331753c5	sql	Backfill `report_card.type`	\N	4.26.0	\N	\N	9044340460
v49.2024-01-22T11:52:00	qnkhuat	migrations/001_update_migrations.yaml	2024-06-22 08:19:14.473134	251	EXECUTED	9:7d260ed302016469ea5950bba09cb471	customChange	Backfill `report_card.type`	\N	4.26.0	\N	\N	9044340460
v49.2024-01-29T19:26:40	adam-james	migrations/001_update_migrations.yaml	2024-06-22 08:19:14.498409	252	EXECUTED	9:04671d8b09ff919ab5154604377ba247	addColumn tableName=report_dashboard	Add width setting to Dashboards	\N	4.26.0	\N	\N	9044340460
v49.2024-01-29T19:30:00	adam-james	migrations/001_update_migrations.yaml	2024-06-22 08:19:14.516869	253	EXECUTED	9:856b864c837cd25b3de595691f2b5712	update tableName=report_dashboard	Update existing report_dashboard width values to full	\N	4.26.0	\N	\N	9044340460
v49.2024-01-29T19:56:40	adam-james	migrations/001_update_migrations.yaml	2024-06-22 08:19:14.558427	254	EXECUTED	9:c092ab859ee6994fcd3b4719872abb97	addNotNullConstraint columnName=width, tableName=report_dashboard	Dashboard width setting must have a value	\N	4.26.0	\N	\N	9044340460
v49.2024-01-29T19:59:12	adam-james	migrations/001_update_migrations.yaml	2024-06-22 08:19:14.576685	255	MARK_RAN	9:663770b688f1ec9aee90d067d789ddbe	addDefaultValue columnName=width, tableName=report_dashboard	Add default value to report_dashboard.width for mysql and mariadb	\N	4.26.0	\N	\N	9044340460
v49.2024-02-02T11:27:49	oisincoveney	migrations/001_update_migrations.yaml	2024-06-22 08:19:14.612271	256	EXECUTED	9:b98d0471c408e0c13d3e6760275a7252	addColumn tableName=report_card	Add report_card.initially_published_at	\N	4.26.0	\N	\N	9044340460
v49.2024-02-02T11:28:36	oisincoveney	migrations/001_update_migrations.yaml	2024-06-22 08:19:14.646473	257	EXECUTED	9:749dba47625d9fdb89d366e3c43dc510	addColumn tableName=report_dashboard	Add report_dashboard.initially_published_at	\N	4.26.0	\N	\N	9044340460
v49.2024-02-07T21:52:01	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:14.680949	258	EXECUTED	9:147145d4f57327ceeddf6c0e64605488	sqlFile path=instance_analytics_views/view_log/v2/postgres-view_log.sql; sqlFile path=instance_analytics_views/view_log/v2/mysql-view_log.sql; sqlFile path=instance_analytics_views/view_log/v2/h2-view_log.sql	Added 0.49.0 - updated view v_view_log	\N	4.26.0	\N	\N	9044340460
v49.2024-02-07T21:52:02	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:14.712499	259	EXECUTED	9:ad858cd4bea316258d9a3894f2ef27d3	sqlFile path=instance_analytics_views/audit_log/v2/postgres-audit_log.sql; sqlFile path=instance_analytics_views/audit_log/v2/mysql-audit_log.sql; sqlFile path=instance_analytics_views/audit_log/v2/h2-audit_log.sql	Added 0.49.0 - updated view v_audit_log	\N	4.26.0	\N	\N	9044340460
v49.2024-02-07T21:52:03	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:14.746509	260	EXECUTED	9:9072a2d1352e44778ba6ecda18ba2177	sqlFile path=instance_analytics_views/group_members/v2/group_members.sql	Added 0.49.0 - updated view v_group_members	\N	4.26.0	\N	\N	9044340460
v49.2024-02-07T21:52:04	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:14.781282	261	EXECUTED	9:f89937596b5e1dc112bdd6f142dbce78	sqlFile path=instance_analytics_views/query_log/v2/postgres-query_log.sql; sqlFile path=instance_analytics_views/query_log/v2/mysql-query_log.sql; sqlFile path=instance_analytics_views/query_log/v2/h2-query_log.sql	Added 0.49.0 - updated view v_query_log	\N	4.26.0	\N	\N	9044340460
v49.2024-02-07T21:52:05	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:14.812134	262	EXECUTED	9:55199a2ff1141c96910ba3e171da5d34	sqlFile path=instance_analytics_views/users/v2/postgres-users.sql; sqlFile path=instance_analytics_views/users/v2/mysql-users.sql; sqlFile path=instance_analytics_views/users/v2/h2-users.sql	Added 0.49.0 - updated view v_users	\N	4.26.0	\N	\N	9044340460
v49.2024-02-09T13:55:26	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:14.850026	263	MARK_RAN	9:d6f29528c18573873a35c84be228238a	sql; sql; sql	Set default value for enable-public-sharing to `false` for existing instances with users, if not already set	\N	4.26.0	\N	\N	9044340460
v49.2024-03-26T10:23:12	adam-james	migrations/001_update_migrations.yaml	2024-06-22 08:19:14.873133	264	EXECUTED	9:32218014a4cff780ad4c8e042fc98f23	addColumn tableName=pulse_card	Add pulse_card.format_rows	\N	4.26.0	\N	\N	9044340460
v49.2024-03-26T20:27:58	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:14.993751	265	EXECUTED	9:469ae1b8af4545f2f48d9505fd34059d	customChange	Added 0.46.0 - Delete the truncate audit log task (renamed to truncate audit tables)	\N	4.26.0	\N	\N	9044340460
v49.2024-04-09T10:00:00	qnkhuat	migrations/001_update_migrations.yaml	2024-06-22 08:19:15.023589	266	EXECUTED	9:0584ae323d98a60f80b433bedab2a0a2	dropNotNullConstraint columnName=cache_field_values_schedule, tableName=metabase_database	Drop not null constraint on metabase_database.cache_field_values_schedule	\N	4.26.0	\N	\N	9044340460
v49.2024-04-09T10:00:01	qnkhuat	migrations/001_update_migrations.yaml	2024-06-22 08:19:15.058043	267	EXECUTED	9:68341c448ecdca425eae019974b64c09	dropDefaultValue columnName=cache_field_values_schedule, tableName=metabase_database	Drop default value on metabase_database.cache_field_values_schedule	\N	4.26.0	\N	\N	9044340460
v49.2024-04-09T10:00:02	qnkhuat	migrations/001_update_migrations.yaml	2024-06-22 08:19:15.096685	268	EXECUTED	9:adb15f2d80ea2b037b44d1eb2dbcb3c5	addDefaultValue columnName=cache_field_values_schedule, tableName=metabase_database	Add null as default value for metabase_database.cache_field_values_schedule	\N	4.26.0	\N	\N	9044340460
v50.2024-03-22T00:38:28	qnkhuat	migrations/001_update_migrations.yaml	2024-06-22 08:19:16.400473	296	EXECUTED	9:24cd167582544658226746847e145181	createTable tableName=field_usage	Add field_usage table	\N	4.26.0	\N	\N	9044340460
v49.2024-04-09T10:00:03	qnkhuat	migrations/001_update_migrations.yaml	2024-06-22 08:19:15.124475	269	EXECUTED	9:fb84fa8ea82ea520edae45164aa167d9	customChange	This clears the schedule for caching field values for databases with period scanning disabled.	\N	4.26.0	\N	\N	9044340460
v49.2024-05-07T10:00:00	qnkhuat	migrations/001_update_migrations.yaml	2024-06-22 08:19:15.131675	270	MARK_RAN	9:d41d8cd98f00b204e9800998ecf8427e	sql	Set revision.most_recent = true to latest revision and false to others. A redo of v48.00-008 for mysql	\N	4.26.0	\N	\N	9044340460
v49.2024-05-20T20:37:55	johnswanson	migrations/001_update_migrations.yaml	2024-06-22 08:19:15.148407	271	MARK_RAN	9:99f4a155e5f4eb4debf6af01cdf68b8e	addNotNullConstraint columnName=collection_preview, tableName=report_card	Add NOT NULL constraint to report_card.collection_preview	\N	4.26.0	\N	\N	9044340460
v49.2024-05-20T20:38:34	johnswanson	migrations/001_update_migrations.yaml	2024-06-22 08:19:15.163752	272	MARK_RAN	9:4915c63aed5f55151a7fe470a040862f	addDefaultValue columnName=collection_preview, tableName=report_card	Add default value to report_card.collection_preview	\N	4.26.0	\N	\N	9044340460
v50.2024-01-09T13:52:21	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:15.245659	274	EXECUTED	9:3c438c9a400361ed3dedc2818c7ed4b8	createIndex indexName=idx_data_permissions_table_id, tableName=data_permissions	Index on data_permissions.table_id	\N	4.26.0	\N	\N	9044340460
v50.2024-01-09T13:53:50	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:15.303436	275	EXECUTED	9:21ad67aa43df73593d4e86492afdfb4e	createIndex indexName=idx_data_permissions_db_id, tableName=data_permissions	Index on data_permissions.db_id	\N	4.26.0	\N	\N	9044340460
v50.2024-01-09T13:53:54	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:15.373158	276	EXECUTED	9:34a2e61c8298c84b18f7a7c7bf5d2119	createIndex indexName=idx_data_permissions_group_id, tableName=data_permissions	Index on data_permissions.group_id	\N	4.26.0	\N	\N	9044340460
v50.2024-01-10T03:27:29	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:15.435886	277	EXECUTED	9:a8c8b1c823344ea148826e9772c2b9b7	dropForeignKeyConstraint baseTableName=sandboxes, constraintName=fk_sandboxes_ref_permissions	Drop foreign key constraint on sandboxes.permissions_id	\N	4.26.0	\N	\N	9044340460
v50.2024-01-10T03:27:30	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:15.475449	278	EXECUTED	9:9fea61b6c7e6a2bb9c023280801eb921	sqlFile path=permissions/data_access.sql; sqlFile path=permissions/mysql_data_access.sql	Migrate data-access permissions from `permissions` to `data_permissions`	\N	4.26.0	\N	\N	9044340460
v50.2024-01-10T03:27:31	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:15.497314	279	EXECUTED	9:4ea632c5ef082ebf005d923bd5527be3	sqlFile path=permissions/native_query_editing.sql	Migrate native-query-editing permissions from `permissions` to `data_permissions`	\N	4.26.0	\N	\N	9044340460
v50.2024-01-10T03:27:32	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:15.552807	280	EXECUTED	9:dad70b4fb60d65b115e48eb2c0f66642	sqlFile path=permissions/download_results.sql; sqlFile path=permissions/mysql_download_results.sql	Migrate download-results permissions from `permissions` to `data_permissions`	\N	4.26.0	\N	\N	9044340460
v50.2024-01-10T03:27:33	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:15.581131	281	EXECUTED	9:62a4a5c7fad5e2bbd48469b748576c48	sqlFile path=permissions/manage_table_metadata.sql; sqlFile path=permissions/mysql_manage_table_metadata.sql	Migrate manage-data-metadata permissions from `permissions` to `data_permissions`	\N	4.26.0	\N	\N	9044340460
v50.2024-01-10T03:27:34	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:15.650867	282	EXECUTED	9:9c37edf4ba12133ce298fb806af1a9cb	sqlFile path=permissions/manage_database.sql	Migrate manage-database permissions from `permissions` to `data_permissions`	\N	4.26.0	\N	\N	9044340460
v50.2024-02-19T21:32:04	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:15.6662	283	EXECUTED	9:2a02ffb8f60532c88e645d2cd9053d95	sql	Clear data permission paths	\N	4.26.0	\N	\N	9044340460
v50.2024-02-20T19:21:04	camsaul	migrations/001_update_migrations.yaml	2024-06-22 08:19:15.696087	284	EXECUTED	9:32e06eedee74b10b4f03aa0c64e5f19c	sql	Drop v1 version of v_content view since it references report_card.dataset which we are dropping in next migration	\N	4.26.0	\N	\N	9044340460
v50.2024-02-20T19:25:40	camsaul	migrations/001_update_migrations.yaml	2024-06-22 08:19:15.758425	285	EXECUTED	9:1a0d14160b0a4e346ffe38d0e1009b7e	dropColumn columnName=dataset, tableName=report_card	Remove report_card.dataset (indicated whether Card was a Model; migrated to report_card.type in 49)	\N	4.26.0	\N	\N	9044340460
v50.2024-02-20T19:26:38	camsaul	migrations/001_update_migrations.yaml	2024-06-22 08:19:15.791428	286	EXECUTED	9:ba2590be4f62444dd99d688ffff286ba	sqlFile path=instance_analytics_views/content/v2/postgres-content.sql; sqlFile path=instance_analytics_views/content/v2/mysql-content.sql; sqlFile path=instance_analytics_views/content/v2/h2-content.sql	Add new v2 version of v_content view which uses report_card.type instead of report_card.dataset (removed in previous migration)	\N	4.26.0	\N	\N	9044340460
v50.2024-02-26T22:15:54	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:15.856118	287	EXECUTED	9:68d8f2fdf00395a9e614b0db8a2ec7bb	sqlFile path=permissions/view_data.sql; sqlFile path=permissions/mysql_view_data.sql	New `view-data` permission	\N	4.26.0	\N	\N	9044340460
v50.2024-02-26T22:15:55	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:15.869736	288	EXECUTED	9:7c2bc517b6def4e3278394b0399c3d4b	sqlFile path=permissions/create_queries.sql	New `create_queries` permission	\N	4.26.0	\N	\N	9044340460
v50.2024-02-29T15:06:43	tsmacdonald	migrations/001_update_migrations.yaml	2024-06-22 08:19:15.962557	289	EXECUTED	9:bf61f630cb8d21a16ef6566d842a8b46	createTable tableName=query_field	Add the query_field join table	\N	4.26.0	\N	\N	9044340460
v50.2024-02-29T15:07:43	tsmacdonald	migrations/001_update_migrations.yaml	2024-06-22 08:19:16.006745	290	EXECUTED	9:0abc329e07ca6746b3e3da79a1bb4fdb	createIndex indexName=idx_query_field_card_id, tableName=query_field	Index query_field.card_id	\N	4.26.0	\N	\N	9044340460
v50.2024-02-29T15:08:43	tsmacdonald	migrations/001_update_migrations.yaml	2024-06-22 08:19:16.085393	291	EXECUTED	9:27ce36c055a2b252a6c6e52d5f5f1d76	createIndex indexName=idx_query_field_field_id, tableName=query_field	Index query_field.field_id	\N	4.26.0	\N	\N	9044340460
v50.2024-03-12T17:16:38	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:16.158107	292	EXECUTED	9:ee42f814c124ba84636b9b5fd1c34bdf	dropTable tableName=activity	Drops the `activity` table which is now unused	\N	4.26.0	\N	\N	9044340460
v50.2024-03-18T16:00:00	piranha	migrations/001_update_migrations.yaml	2024-06-22 08:19:16.255208	293	EXECUTED	9:f0cd082514b5a1ffdc432034df0d837e	createTable tableName=cache_config	Effective caching #36847	\N	4.26.0	\N	\N	9044340460
v50.2024-03-18T16:00:01	piranha	migrations/001_update_migrations.yaml	2024-06-22 08:19:16.296843	294	EXECUTED	9:5007337d4de322587ce6056e0b9eac9b	addUniqueConstraint constraintName=idx_cache_config_unique_model, tableName=cache_config	Effective caching #36847	\N	4.26.0	\N	\N	9044340460
v50.2024-03-21T17:41:00	qnkhuat	migrations/001_update_migrations.yaml	2024-06-22 08:19:16.357539	295	EXECUTED	9:ecace972ee38bae5dd02e23ff91e2794	addColumn tableName=metabase_table	Add metabase_table.estimated_row_count	\N	4.26.0	\N	\N	9044340460
v50.2024-03-22T00:39:28	qnkhuat	migrations/001_update_migrations.yaml	2024-06-22 08:19:16.438844	297	EXECUTED	9:c94534b092f89579385156f18faee569	createIndex indexName=idx_field_usage_field_id, tableName=field_usage	Index field_usage.field_id	\N	4.26.0	\N	\N	9044340460
v50.2024-03-24T19:34:11	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:16.454795	298	EXECUTED	9:ceac2712d4c61fbbafa21b905ee7cf35	sql	Clean up deprecated view-data and native-query-editing permissions	\N	4.26.0	\N	\N	9044340460
v50.2024-03-25T14:53:00	tsmacdonald	migrations/001_update_migrations.yaml	2024-06-22 08:19:16.48777	299	EXECUTED	9:86d60ff7f9ab572c53887c13e96e4315	addColumn tableName=query_field	Add query_field.direct_reference	\N	4.26.0	\N	\N	9044340460
v50.2024-03-28T16:30:35	calherries	migrations/001_update_migrations.yaml	2024-06-22 08:19:16.692663	300	EXECUTED	9:db68375fac0a9fe40f0caf475dd1fd92	customChange	Create internal user	\N	4.26.0	\N	\N	9044340460
v50.2024-03-29T10:00:00	piranha	migrations/001_update_migrations.yaml	2024-06-22 08:19:16.757802	301	EXECUTED	9:5b64f09469c0d093e147cc93dbaa94e5	addColumn tableName=report_card	Granular cache invalidation	\N	4.26.0	\N	\N	9044340460
v50.2024-04-09T15:55:19	calherries	migrations/001_update_migrations.yaml	2024-06-22 08:19:16.792705	302	EXECUTED	9:3af4d0d7b06784b822eed6b5b7c525f3	addColumn tableName=collection	Add collection.is_sample column	\N	4.26.0	\N	\N	9044340460
v50.2024-04-15T16:30:35	qnkhuat	migrations/001_update_migrations.yaml	2024-06-22 08:19:16.827312	303	EXECUTED	9:7f083cce96c96c2b5d9b21cc6ad91ef8	addColumn tableName=report_card	Add report_card.last_used_at	\N	4.26.0	\N	\N	9044340460
v50.2024-04-19T17:04:04	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:16.855446	304	EXECUTED	9:ceac2712d4c61fbbafa21b905ee7cf35	sql	Clean up deprecated view-data and native-query-editing permissions (again)	\N	4.26.0	\N	\N	9044340460
v50.2024-04-25T01:04:05	qnkhuat	migrations/001_update_migrations.yaml	2024-06-22 08:19:16.924129	305	EXECUTED	9:61c4fc54c2c8ab0ea6b072ec1409fa48	customChange	Delete the old SendPulses job and trigger	\N	4.26.0	\N	\N	9044340460
v50.2024-04-25T01:04:06	qnkhuat	migrations/001_update_migrations.yaml	2024-06-22 08:19:16.945596	306	EXECUTED	9:ca85ee4382798cec047044ab394061c0	customChange	Delete SendPulse Job on downgrade	\N	4.26.0	\N	\N	9044340460
v50.2024-04-25T01:04:07	qnkhuat	migrations/001_update_migrations.yaml	2024-06-22 08:19:16.959517	307	EXECUTED	9:8380c81746dd65909faceeef3461133e	customChange	Delete InitSendPulseTriggers Job on downgrade	\N	4.26.0	\N	\N	9044340460
v50.2024-04-25T03:15:01	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:16.995605	308	EXECUTED	9:669daa77b4a7c048fbb8bfb4fee54df0	addColumn tableName=core_user	Add entity_id to core_user	\N	4.26.0	\N	\N	9044340460
v50.2024-04-25T03:15:02	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:17.036273	309	EXECUTED	9:477fd6e94e5b5bcd818b811928af19ff	addColumn tableName=permissions_group	Add entity_id to permissions_group	\N	4.26.0	\N	\N	9044340460
v50.2024-04-25T16:29:31	calherries	migrations/001_update_migrations.yaml	2024-06-22 08:19:17.06297	310	EXECUTED	9:8241c6c3caca2f4eeccfdd157063da58	addColumn tableName=report_card	Add report_card.view_count	\N	4.26.0	\N	\N	9044340460
v50.2024-04-25T16:29:32	calherries	migrations/001_update_migrations.yaml	2024-06-22 08:19:17.080814	311	EXECUTED	9:807983e5cff6386bb21e3cd2197c076d	sql; sql	Populate report_card.view_count	\N	4.26.0	\N	\N	9044340460
v50.2024-04-25T16:29:33	calherries	migrations/001_update_migrations.yaml	2024-06-22 08:19:17.149182	312	EXECUTED	9:c443c22fa32cc1e0dbee3669ef0398ee	addColumn tableName=report_dashboard	Add report_dashboard.view_count	\N	4.26.0	\N	\N	9044340460
v50.2024-04-25T16:29:34	calherries	migrations/001_update_migrations.yaml	2024-06-22 08:19:17.165319	313	EXECUTED	9:a0eef6561d9759a96d72b22c16b29f1e	sql; sql	Populate report_dashboard.view_count	\N	4.26.0	\N	\N	9044340460
v50.2024-04-25T16:29:35	calherries	migrations/001_update_migrations.yaml	2024-06-22 08:19:17.192373	314	EXECUTED	9:5873cdd28eb07ce7ea60ab44e9646338	addColumn tableName=metabase_table	Add metabase_table.view_count	\N	4.26.0	\N	\N	9044340460
v50.2024-04-25T16:29:36	calherries	migrations/001_update_migrations.yaml	2024-06-22 08:19:17.206963	315	EXECUTED	9:b55ed9c3e59124df25d78d191c165508	sql; sql	Populate metabase_table.view_count	\N	4.26.0	\N	\N	9044340460
v50.2024-04-26T09:19:00	adam-james	migrations/001_update_migrations.yaml	2024-06-22 08:19:17.255709	316	EXECUTED	9:713beca44f427ff4832de85b1e88320b	createTable tableName=user_parameter_value	Added 0.50.0 - Per-user Dashboard Parameter values table	\N	4.26.0	\N	\N	9044340460
v50.2024-04-26T09:25:00	adam-james	migrations/001_update_migrations.yaml	2024-06-22 08:19:17.297054	317	EXECUTED	9:2af6f33e022984b2abd41fb0a6820b92	createIndex indexName=idx_user_parameter_value_user_id, tableName=user_parameter_value	Index user_parameter_value.user_id	\N	4.26.0	\N	\N	9044340460
v50.2024-04-30T23:57:23	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:17.329297	318	EXECUTED	9:19d659847dc567d1abd78734b8672ea9	addColumn tableName=api_key	Add `scope` column to api_key to support SCIM authentication	\N	4.26.0	\N	\N	9044340460
v50.2024-04-30T23:58:24	noahmoss	migrations/001_update_migrations.yaml	2024-06-22 08:19:17.358189	319	EXECUTED	9:a6914c25365739b16b573067fae3a383	dropNotNullConstraint columnName=user_id, tableName=api_key	Drop NOT NULL constraint on api_key.user_id to support SCIM-scoped API keys	\N	4.26.0	\N	\N	9044340460
v50.2024-05-08T09:00:00	qnkhuat	migrations/001_update_migrations.yaml	2024-06-22 08:19:17.395686	320	EXECUTED	9:e28d26db534e046b411c8fe176c1c0e0	addColumn tableName=task_history	Add task_history.status	\N	4.26.0	\N	\N	9044340460
v50.2024-05-08T09:00:01	qnkhuat	migrations/001_update_migrations.yaml	2024-06-22 08:19:17.426432	321	EXECUTED	9:c8df0ec666d8c736f83bdb85e446fed1	dropDefaultValue columnName=status, tableName=task_history	Drop default value task_history.status	\N	4.26.0	\N	\N	9044340460
v50.2024-05-08T09:00:02	qnkhuat	migrations/001_update_migrations.yaml	2024-06-22 08:19:17.456424	322	EXECUTED	9:c9cd34125a3c9445f4d7076527b95589	addDefaultValue columnName=status, tableName=task_history	Add "started" as default value for task_history.status, now that backfill is done.	\N	4.26.0	\N	\N	9044340460
v50.2024-05-08T09:00:03	qnkhuat	migrations/001_update_migrations.yaml	2024-06-22 08:19:17.505161	323	EXECUTED	9:9a753f9d1a43011f20be4eff7634cd33	dropNotNullConstraint columnName=ended_at, tableName=task_history	Drop not null constraint for task_history.ended_at	\N	4.26.0	\N	\N	9044340460
v50.2024-05-08T09:00:04	qnkhuat	migrations/001_update_migrations.yaml	2024-06-22 08:19:17.530343	324	EXECUTED	9:20c322e8ee6d62c43ae49d8891f81f89	dropNotNullConstraint columnName=duration, tableName=task_history	Drop not null constraint for task_history.duration	\N	4.26.0	\N	\N	9044340460
v50.2024-05-08T09:00:05	qnkhuat	migrations/001_update_migrations.yaml	2024-06-22 08:19:17.57328	325	EXECUTED	9:b02406ab39b905df65acccf4bdc82ee1	dropDefaultValue columnName=ended_at, tableName=task_history	Drop default value task_history.ended_at	\N	4.26.0	\N	\N	9044340460
v50.2024-05-08T09:00:06	qnkhuat	migrations/001_update_migrations.yaml	2024-06-22 08:19:17.600908	326	EXECUTED	9:7ab667ce198a98ccb03ec77b2f1b3655	addDefaultValue columnName=ended_at, tableName=task_history	Add null as default value for task_history.ended_at	\N	4.26.0	\N	\N	9044340460
v50.2024-05-13T16:00:00	filipesilva	migrations/001_update_migrations.yaml	2024-06-22 08:19:17.655145	327	EXECUTED	9:5a355ae0134f6b643ff3b5d4814a15f2	createTable tableName=cloud_migration	Create cloud migration	\N	4.26.0	\N	\N	9044340460
v50.2024-05-15T13:13:13	adam-james	migrations/001_update_migrations.yaml	2024-06-22 08:19:17.71887	328	EXECUTED	9:a7f15cd934779ad0d52c634c3fdca04f	customChange	Fix visualization settings for stacked area/bar/combo displays	\N	4.26.0	\N	\N	9044340460
v50.2024-05-17T19:54:23	calherries	migrations/001_update_migrations.yaml	2024-06-22 08:19:17.743133	329	EXECUTED	9:a3525795c74bd10591b66a2b9138aec3	addColumn tableName=metabase_database	Add metabase_database.uploads_enabled column	\N	4.26.0	\N	\N	9044340460
v50.2024-05-17T19:54:24	calherries	migrations/001_update_migrations.yaml	2024-06-22 08:19:17.774317	330	EXECUTED	9:76516d4fd0f8a99a1232d999e5f834a8	addColumn tableName=metabase_database	Add metabase_database.uploads_schema_name column	\N	4.26.0	\N	\N	9044340460
v50.2024-05-17T19:54:25	calherries	migrations/001_update_migrations.yaml	2024-06-22 08:19:17.808973	331	EXECUTED	9:9df4d358f294a0541342ceeb9a43c12a	addColumn tableName=metabase_database	Add metabase_database.uploads_table_prefix column	\N	4.26.0	\N	\N	9044340460
v50.2024-05-17T19:54:26	calherries	migrations/001_update_migrations.yaml	2024-06-22 08:19:17.832527	332	EXECUTED	9:2ba2a860ad40c70381cddf179a766d08	customChange	Update metabase_database.uploads_enabled value	\N	4.26.0	\N	\N	9044340460
v50.2024-05-27T15:55:22	calherries	migrations/001_update_migrations.yaml	2024-06-22 08:19:19.332207	333	EXECUTED	9:47cab1070db15a6d225dcf467747120a	customChange	Create sample content	\N	4.26.0	\N	\N	9044340460
v50.2024-06-12T12:33:07	piranha	migrations/001_update_migrations.yaml	2024-06-22 08:19:20.787559	334	EXECUTED	9:6d76b0bddae5c21628ed0e2c3600e5ce	customChange	Decrypt some settings so the next migration runs well	\N	4.26.0	\N	\N	9044340460
v50.2024-06-12T12:33:08	piranha	migrations/001_update_migrations.yaml	2024-06-22 08:19:21.068536	335	EXECUTED	9:b9c980de6526154b21d4780b57c19f05	sqlFile path=custom_sql/fill_cache_config.pg.sql; sqlFile path=custom_sql/fill_cache_config.my.sql; sqlFile path=custom_sql/fill_cache_config.h2.sql	Copy old cache configurations to cache_config table	\N	4.26.0	\N	\N	9044340460
\.


--
-- Data for Name: databasechangeloglock; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.databasechangeloglock (id, locked, lockgranted, lockedby) FROM stdin;
1	f	\N	\N
\.


--
-- Data for Name: dependency; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.dependency (id, model, model_id, dependent_on_model, dependent_on_id, created_at) FROM stdin;
\.


--
-- Data for Name: dimension; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.dimension (id, field_id, name, type, human_readable_field_id, created_at, updated_at, entity_id) FROM stdin;
\.


--
-- Data for Name: field_usage; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.field_usage (id, field_id, query_execution_id, used_in, filter_op, aggregation_function, breakout_temporal_unit, breakout_binning_strategy, breakout_binning_num_bins, breakout_binning_bin_width, created_at) FROM stdin;
\.


--
-- Data for Name: http_action; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.http_action (action_id, template, response_handle, error_handle) FROM stdin;
\.


--
-- Data for Name: implicit_action; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.implicit_action (action_id, kind) FROM stdin;
\.


--
-- Data for Name: label; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.label (id, name, slug, icon) FROM stdin;
\.


--
-- Data for Name: login_history; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.login_history (id, "timestamp", user_id, session_id, device_id, device_description, ip_address) FROM stdin;
1	2024-06-22 08:25:50.592174+00	1	aff0e644-0b3a-4483-8134-897f53e96131	6c36effc-c064-4105-9e01-65cb68b6abc2	Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36	127.0.0.1
\.


--
-- Data for Name: metabase_database; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.metabase_database (id, created_at, updated_at, name, description, details, engine, is_sample, is_full_sync, points_of_interest, caveats, metadata_sync_schedule, cache_field_values_schedule, timezone, is_on_demand, auto_run_queries, refingerprint, cache_ttl, initial_sync_status, creator_id, settings, dbms_version, is_audit, uploads_enabled, uploads_schema_name, uploads_table_prefix) FROM stdin;
1	2024-06-22 08:19:17.839482+00	2024-06-22 08:23:40.924338+00	Sample Database	Some example data for you to play around with.	{"db":"file:/plugins/sample-database.db;USER=GUEST;PASSWORD=guest"}	h2	t	t	You can find all sorts of different joinable tables ranging from products to people to reviews here.	You probably don't want to use this for your business-critical analyses, but hey, it's your world, we're just living in it.	0 26 * * * ? *	0 0 17 * * ? *	GMT	f	t	\N	\N	complete	\N	\N	{"flavor":"H2","version":"2.1.214 (2022-06-13)","semantic-version":[2,1]}	f	f	\N	\N
2	2024-06-22 08:25:35.534063+00	2024-06-22 08:25:37.013467+00	pgsql	\N	{"ssl":false,"password":"aman","port":5432,"advanced-options":false,"schema-filters-type":"all","dbname":"compute","host":"pg-primary","tunnel-enabled":false,"user":"aman"}	postgres	f	t	\N	\N	0 48 * * * ? *	0 0 2 * * ? *	GMT	f	t	\N	\N	complete	1	\N	{"flavor":"PostgreSQL","version":"16.3","semantic-version":[16,3]}	f	f	\N	\N
\.


--
-- Data for Name: metabase_field; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.metabase_field (id, created_at, updated_at, name, base_type, semantic_type, active, description, preview_display, "position", table_id, parent_id, display_name, visibility_type, fk_target_field_id, last_analyzed, points_of_interest, caveats, fingerprint, fingerprint_version, database_type, has_field_values, settings, database_position, custom_position, effective_type, coercion_strategy, nfc_path, database_required, json_unfolding, database_is_auto_increment, database_indexed, database_partitioned) FROM stdin;
1	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	SEATS	type/Integer	\N	t	\N	t	6	6	\N	Seats	normal	\N	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":102,"nil%":0.0},"type":{"type/Number":{"min":1.0,"q1":2.4309856865966593,"q3":10.553778422458695,"max":1325.0,"sd":51.198301031505444,"avg":16.21763527054108}}}	5	INTEGER	auto-list	\N	6	0	type/Integer	\N	\N	f	f	f	f	\N
2	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	TRIAL_ENDS_AT	type/DateTime	\N	t	\N	t	8	6	\N	Trial Ends At	normal	\N	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":1712,"nil%":0.001202404809619238},"type":{"type/DateTime":{"earliest":"2020-09-30T12:00:00Z","latest":"2031-10-25T12:00:00Z"}}}	5	TIMESTAMP	\N	\N	8	0	type/DateTime	\N	\N	f	f	f	f	\N
3	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	CANCELED_AT	type/DateTime	type/CancelationTimestamp	t	\N	t	9	6	\N	Canceled At	normal	\N	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":2021,"nil%":0.1859719438877756},"type":{"type/DateTime":{"earliest":"2020-10-01T15:43:40Z","latest":"2032-06-03T14:01:15Z"}}}	5	TIMESTAMP	\N	\N	9	0	type/DateTime	\N	\N	f	f	f	f	\N
4	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	PLAN	type/Text	type/Category	t	\N	t	4	6	\N	Plan	normal	\N	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":3,"nil%":0.0},"type":{"type/Text":{"percent-json":0.0,"percent-url":0.0,"percent-email":0.0,"percent-state":0.0,"average-length":5.1062124248497}}}	5	CHARACTER VARYING	auto-list	\N	4	0	type/Text	\N	\N	f	f	f	f	\N
5	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	ACTIVE_SUBSCRIPTION	type/Boolean	type/Category	t	\N	t	11	6	\N	Active Subscription	normal	\N	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":2,"nil%":0.0}}	5	BOOLEAN	auto-list	\N	11	0	type/Boolean	\N	\N	f	f	f	f	\N
6	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	LATITUDE	type/Float	type/Latitude	t	\N	t	13	6	\N	Latitude	normal	\N	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":2472,"nil%":4.008016032064128E-4},"type":{"type/Number":{"min":-48.75,"q1":19.430679334308675,"q3":47.24585743676113,"max":69.23111,"sd":23.492041679980137,"avg":31.35760681046913}}}	5	DOUBLE PRECISION	\N	\N	13	0	type/Float	\N	\N	f	f	f	f	\N
7	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	ID	type/BigInteger	type/PK	t	\N	t	0	6	\N	ID	normal	\N	\N	\N	\N	\N	0	BIGINT	\N	\N	0	0	type/BigInteger	\N	\N	t	f	f	t	\N
8	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	LEGACY_PLAN	type/Boolean	type/Category	t	\N	t	12	6	\N	Legacy Plan	normal	\N	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":2,"nil%":0.0}}	5	BOOLEAN	auto-list	\N	12	0	type/Boolean	\N	\N	f	f	f	f	\N
9	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	FIRST_NAME	type/Text	type/Name	t	\N	t	2	6	\N	First Name	normal	\N	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":1687,"nil%":0.0},"type":{"type/Text":{"percent-json":0.0,"percent-url":0.0,"percent-email":0.0,"percent-state":0.001603206412825651,"average-length":5.997595190380761}}}	5	CHARACTER VARYING	\N	\N	2	0	type/Text	\N	\N	f	f	f	f	\N
10	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	LAST_NAME	type/Text	type/Name	t	\N	t	3	6	\N	Last Name	normal	\N	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":473,"nil%":0.0},"type":{"type/Text":{"percent-json":0.0,"percent-url":0.0,"percent-email":0.0,"percent-state":0.0,"average-length":6.536673346693386}}}	5	CHARACTER VARYING	auto-list	\N	3	0	type/Text	\N	\N	f	f	f	f	\N
11	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	SOURCE	type/Text	type/Source	t	\N	t	5	6	\N	Source	normal	\N	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":5,"nil%":0.3346693386773547},"type":{"type/Text":{"percent-json":0.0,"percent-url":0.0,"percent-email":0.0,"percent-state":0.0,"average-length":4.4705410821643286}}}	5	CHARACTER VARYING	auto-list	\N	5	0	type/Text	\N	\N	f	f	f	f	\N
12	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	EMAIL	type/Text	type/Email	t	\N	t	1	6	\N	Email	normal	\N	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":2494,"nil%":0.0},"type":{"type/Text":{"percent-json":0.0,"percent-url":0.0,"percent-email":1.0,"percent-state":0.0,"average-length":28.185971943887775}}}	5	CHARACTER VARYING	\N	\N	1	0	type/Text	\N	\N	f	f	f	f	\N
13	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	CREATED_AT	type/DateTime	type/CreationTimestamp	t	\N	t	7	6	\N	Created At	normal	\N	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":2495,"nil%":0.0},"type":{"type/DateTime":{"earliest":"2020-09-15T16:11:50Z","latest":"2031-10-10T19:14:48Z"}}}	5	TIMESTAMP	\N	\N	7	0	type/DateTime	\N	\N	f	f	f	f	\N
14	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	TRIAL_CONVERTED	type/Boolean	type/Category	t	\N	t	10	6	\N	Trial Converted	normal	\N	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":2,"nil%":0.0}}	5	BOOLEAN	auto-list	\N	10	0	type/Boolean	\N	\N	f	f	f	f	\N
15	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	COUNTRY	type/Text	type/Country	t	\N	t	15	6	\N	Country	normal	\N	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":134,"nil%":8.016032064128256E-4},"type":{"type/Text":{"percent-json":0.0,"percent-url":0.0,"percent-email":0.0,"percent-state":0.1130260521042084,"average-length":1.9983967935871743}}}	5	CHARACTER	auto-list	\N	15	0	type/Text	\N	\N	f	f	f	f	\N
16	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	LONGITUDE	type/Float	type/Longitude	t	\N	t	14	6	\N	Longitude	normal	\N	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":2484,"nil%":4.008016032064128E-4},"type":{"type/Number":{"min":-175.06667,"q1":-55.495929410727236,"q3":28.627359769389155,"max":176.21667,"sd":68.51011002740533,"avg":2.6042336031796345}}}	5	DOUBLE PRECISION	\N	\N	14	0	type/Float	\N	\N	f	f	f	f	\N
17	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	EVENT	type/Text	type/Category	t	\N	t	2	1	\N	Event	normal	\N	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":2,"nil%":0.0},"type":{"type/Text":{"percent-json":0.0,"percent-url":0.0,"percent-email":0.0,"percent-state":0.0,"average-length":11.3906}}}	5	CHARACTER VARYING	auto-list	\N	2	0	type/Text	\N	\N	f	f	f	f	\N
18	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	PAGE_URL	type/Text	type/URL	t	\N	t	4	1	\N	Page URL	normal	\N	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":6,"nil%":0.1302},"type":{"type/Text":{"percent-json":0.0,"percent-url":0.0,"percent-email":0.0,"percent-state":0.0,"average-length":22.2674}}}	5	CHARACTER VARYING	auto-list	\N	4	0	type/Text	\N	\N	f	f	f	f	\N
19	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	BUTTON_LABEL	type/Text	type/Category	t	\N	t	5	1	\N	Button Label	normal	\N	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":6,"nil%":0.8698},"type":{"type/Text":{"percent-json":0.0,"percent-url":0.0,"percent-email":0.0,"percent-state":0.0,"average-length":1.0552}}}	5	CHARACTER VARYING	auto-list	\N	5	0	type/Text	\N	\N	f	f	f	f	\N
20	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	ID	type/BigInteger	type/PK	t	\N	t	0	1	\N	ID	normal	\N	\N	\N	\N	\N	0	BIGINT	\N	\N	0	0	type/BigInteger	\N	\N	t	f	f	t	\N
21	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	ACCOUNT_ID	type/BigInteger	type/FK	t	\N	t	1	1	\N	Account ID	normal	7	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":589,"nil%":0.0}}	5	BIGINT	\N	\N	1	0	type/BigInteger	\N	\N	f	f	f	t	\N
22	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	TIMESTAMP	type/DateTime	\N	t	\N	t	3	1	\N	Timestamp	normal	\N	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":8576,"nil%":0.0},"type":{"type/DateTime":{"earliest":"2022-03-15T00:18:25Z","latest":"2022-04-11T20:24:02Z"}}}	5	TIMESTAMP	\N	\N	3	0	type/DateTime	\N	\N	f	f	f	f	\N
23	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	RATING	type/Integer	type/Score	t	\N	t	4	2	\N	Rating	normal	\N	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":5,"nil%":0.0},"type":{"type/Number":{"min":1.0,"q1":2.7545289729206877,"q3":4.004191340512663,"max":5.0,"sd":0.8137255616667736,"avg":3.3629283489096573}}}	5	SMALLINT	auto-list	\N	4	0	type/Integer	\N	\N	f	f	f	f	\N
24	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	ID	type/BigInteger	type/PK	t	\N	t	0	2	\N	ID	normal	\N	\N	\N	\N	\N	0	BIGINT	\N	\N	0	0	type/BigInteger	\N	\N	t	f	f	t	\N
25	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	ACCOUNT_ID	type/BigInteger	type/FK	t	\N	t	1	2	\N	Account ID	normal	7	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":642,"nil%":0.0}}	5	BIGINT	\N	\N	1	0	type/BigInteger	\N	\N	f	f	f	t	\N
26	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	EMAIL	type/Text	type/Email	t	\N	t	2	2	\N	Email	normal	\N	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":642,"nil%":0.0},"type":{"type/Text":{"percent-json":0.0,"percent-url":0.0,"percent-email":1.0,"percent-state":0.0,"average-length":28.327102803738317}}}	5	CHARACTER VARYING	auto-list	\N	2	0	type/Text	\N	\N	f	f	f	f	\N
27	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	BODY	type/Text	\N	t	\N	f	6	2	\N	Body	normal	\N	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":642,"nil%":0.0},"type":{"type/Text":{"percent-json":0.0,"percent-url":0.0,"percent-email":0.0,"percent-state":0.0,"average-length":438.15264797507785}}}	5	CHARACTER LARGE OBJECT	\N	\N	6	0	type/Text	\N	\N	f	f	f	f	\N
28	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	RATING_MAPPED	type/Text	type/Category	t	\N	t	5	2	\N	Rating Mapped	normal	\N	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":5,"nil%":0.0},"type":{"type/Text":{"percent-json":0.0,"percent-url":0.0,"percent-email":0.0,"percent-state":0.0,"average-length":6.453271028037383}}}	5	CHARACTER VARYING	auto-list	\N	5	0	type/Text	\N	\N	f	f	f	f	\N
29	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	DATE_RECEIVED	type/DateTime	\N	t	\N	t	3	2	\N	Date Received	normal	\N	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":576,"nil%":0.0},"type":{"type/DateTime":{"earliest":"2020-11-20T00:00:00Z","latest":"2031-12-01T00:00:00Z"}}}	5	TIMESTAMP	\N	\N	3	0	type/DateTime	\N	\N	f	f	f	f	\N
30	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	PLAN	type/Text	type/Category	t	\N	t	4	7	\N	Plan	normal	\N	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":3,"nil%":0.0},"type":{"type/Text":{"percent-json":0.0,"percent-url":0.0,"percent-email":0.0,"percent-state":0.0,"average-length":5.2931}}}	5	CHARACTER VARYING	auto-list	\N	4	0	type/Text	\N	\N	f	f	f	f	\N
31	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	ID	type/BigInteger	type/PK	t	\N	t	0	7	\N	ID	normal	\N	\N	\N	\N	\N	0	BIGINT	\N	\N	0	0	type/BigInteger	\N	\N	t	f	f	t	\N
32	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	ACCOUNT_ID	type/BigInteger	type/FK	t	\N	t	1	7	\N	Account ID	normal	7	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":1449,"nil%":0.0}}	5	BIGINT	\N	\N	1	0	type/BigInteger	\N	\N	f	f	f	t	\N
33	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	DATE_RECEIVED	type/DateTime	\N	t	\N	t	5	7	\N	Date Received	normal	\N	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":714,"nil%":0.0},"type":{"type/DateTime":{"earliest":"2020-09-30T00:00:00Z","latest":"2027-05-02T00:00:00Z"}}}	5	TIMESTAMP	\N	\N	5	0	type/DateTime	\N	\N	f	f	f	f	\N
34	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	PAYMENT	type/Float	\N	t	\N	t	2	7	\N	Payment	normal	\N	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":707,"nil%":0.0},"type":{"type/Number":{"min":13.7,"q1":233.1870107122195,"q3":400.5965814842149,"max":33714.6,"sd":763.7961603932441,"avg":519.4153400000004}}}	5	DOUBLE PRECISION	\N	\N	2	0	type/Float	\N	\N	f	f	f	f	\N
35	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	EXPECTED_INVOICE	type/Boolean	type/Category	t	\N	t	3	7	\N	Expected Invoice	normal	\N	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":2,"nil%":0.0}}	5	BOOLEAN	auto-list	\N	3	0	type/Boolean	\N	\N	f	f	f	f	\N
38	2024-06-22 08:19:17.839482+00	2024-06-22 08:23:39.619455+00	TAX	type/Float	\N	t	This is the amount of local and federal taxes that are collected on the purchase. Note that other governmental fees on some products are not included here, but instead are accounted for in the subtotal.	t	4	5	\N	Tax	normal	\N	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":797,"nil%":0.0},"type":{"type/Number":{"min":0.0,"q1":2.273340386603857,"q3":5.337275338216307,"max":11.12,"sd":2.3206651358900316,"avg":3.8722100000000004}}}	5	DOUBLE PRECISION	\N	\N	4	0	type/Float	\N	\N	f	f	f	f	\N
42	2024-06-22 08:19:17.839482+00	2024-06-22 08:23:39.633046+00	TOTAL	type/Float	\N	t	The total billed amount.	t	5	5	\N	Total	normal	\N	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":4426,"nil%":0.0},"type":{"type/Number":{"min":8.93914247937167,"q1":51.34535490743823,"q3":110.29428389265787,"max":159.34900526552292,"sd":34.26469575709948,"avg":80.35871658771228}}}	5	DOUBLE PRECISION	\N	\N	5	0	type/Float	\N	\N	f	f	f	f	\N
37	2024-06-22 08:19:17.839482+00	2024-06-22 08:23:39.714136+00	ID	type/BigInteger	type/PK	t	This is a unique ID for the product. It is also called the “Invoice number” or “Confirmation number” in customer facing emails and screens.	t	0	5	\N	ID	normal	\N	\N	\N	\N	\N	0	BIGINT	\N	\N	0	0	type/BigInteger	\N	\N	f	f	t	t	\N
43	2024-06-22 08:19:17.839482+00	2024-06-22 08:23:39.725661+00	USER_ID	type/Integer	type/FK	t	The id of the user who made this order. Note that in some cases where an order was created on behalf of a customer who phoned the order in, this might be the employee who handled the request.	t	1	5	\N	User ID	normal	46	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":929,"nil%":0.0}}	5	INTEGER	\N	\N	1	0	type/Integer	\N	\N	f	f	f	t	\N
39	2024-06-22 08:19:17.839482+00	2024-06-22 08:23:39.736624+00	QUANTITY	type/Integer	type/Quantity	t	Number of products bought.	t	8	5	\N	Quantity	normal	\N	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":62,"nil%":0.0},"type":{"type/Number":{"min":0.0,"q1":1.755882607764982,"q3":4.882654507928044,"max":100.0,"sd":4.214258386403798,"avg":3.7015}}}	5	INTEGER	auto-list	\N	8	0	type/Integer	\N	\N	f	f	f	f	\N
36	2024-06-22 08:19:17.839482+00	2024-06-22 08:23:39.814528+00	DISCOUNT	type/Float	type/Discount	t	Discount amount.	t	6	5	\N	Discount	normal	\N	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":701,"nil%":0.898},"type":{"type/Number":{"min":0.17088996672584322,"q1":2.9786226681458743,"q3":7.338187788658235,"max":61.69684269960571,"sd":3.053663125001991,"avg":5.161255547580326}}}	5	DOUBLE PRECISION	\N	\N	6	0	type/Float	\N	\N	f	f	f	f	\N
41	2024-06-22 08:19:17.839482+00	2024-06-22 08:23:39.828142+00	CREATED_AT	type/DateTime	type/CreationTimestamp	t	The date and time an order was submitted.	t	7	5	\N	Created At	normal	\N	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":10001,"nil%":0.0},"type":{"type/DateTime":{"earliest":"2022-04-30T18:56:13.352Z","latest":"2026-04-19T14:07:15.657Z"}}}	5	TIMESTAMP	\N	\N	7	0	type/DateTime	\N	\N	f	f	f	f	\N
47	2024-06-22 08:19:17.839482+00	2024-06-22 08:23:40.0143+00	NAME	type/Text	type/Name	t	The name of the user who owns an account	t	4	3	\N	Name	normal	\N	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":2499,"nil%":0.0},"type":{"type/Text":{"percent-json":0.0,"percent-url":0.0,"percent-email":0.0,"percent-state":0.0,"average-length":13.532}}}	5	CHARACTER VARYING	\N	\N	4	0	type/Text	\N	\N	f	f	f	f	\N
45	2024-06-22 08:19:17.839482+00	2024-06-22 08:23:40.112886+00	SOURCE	type/Text	type/Source	t	The channel through which we acquired this user. Valid values include: Affiliate, Facebook, Google, Organic and Twitter	t	8	3	\N	Source	normal	\N	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":5,"nil%":0.0},"type":{"type/Text":{"percent-json":0.0,"percent-url":0.0,"percent-email":0.0,"percent-state":0.0,"average-length":7.4084}}}	5	CHARACTER VARYING	auto-list	\N	8	0	type/Text	\N	\N	f	f	f	f	\N
46	2024-06-22 08:19:17.839482+00	2024-06-22 08:23:40.143313+00	ID	type/BigInteger	type/PK	t	A unique identifier given to each user.	t	0	3	\N	ID	normal	\N	\N	\N	\N	\N	0	BIGINT	\N	\N	0	0	type/BigInteger	\N	\N	f	f	t	t	\N
40	2024-06-22 08:19:17.839482+00	2024-06-22 08:23:39.382657+00	PRODUCT_ID	type/Integer	type/FK	t	The product ID. This is an internal identifier for the product, NOT the SKU.	t	2	5	\N	Product ID	normal	62	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":200,"nil%":0.0}}	5	INTEGER	\N	\N	2	0	type/Integer	\N	\N	f	f	f	t	\N
51	2024-06-22 08:19:17.839482+00	2024-06-22 08:23:39.928945+00	EMAIL	type/Text	type/Email	t	The contact email for the account.	t	2	3	\N	Email	normal	\N	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":2500,"nil%":0.0},"type":{"type/Text":{"percent-json":0.0,"percent-url":0.0,"percent-email":1.0,"percent-state":0.0,"average-length":24.1824}}}	5	CHARACTER VARYING	\N	\N	2	0	type/Text	\N	\N	f	f	f	f	\N
53	2024-06-22 08:19:17.839482+00	2024-06-22 08:23:39.940753+00	PASSWORD	type/Text	\N	t	This is the salted password of the user. It should not be visible	t	3	3	\N	Password	normal	\N	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":2500,"nil%":0.0},"type":{"type/Text":{"percent-json":0.0,"percent-url":0.0,"percent-email":0.0,"percent-state":0.0,"average-length":36.0}}}	5	CHARACTER VARYING	\N	\N	3	0	type/Text	\N	\N	f	f	f	f	\N
55	2024-06-22 08:19:17.839482+00	2024-06-22 08:23:40.02719+00	CITY	type/Text	type/City	t	The city of the account’s billing address	t	5	3	\N	City	normal	\N	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":1966,"nil%":0.0},"type":{"type/Text":{"percent-json":0.0,"percent-url":0.0,"percent-email":0.0,"percent-state":0.002,"average-length":8.284}}}	5	CHARACTER VARYING	\N	\N	5	0	type/Text	\N	\N	f	f	f	f	\N
48	2024-06-22 08:19:17.839482+00	2024-06-22 08:23:40.055719+00	STATE	type/Text	type/State	t	The state or province of the account’s billing address	t	7	3	\N	State	normal	\N	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":49,"nil%":0.0},"type":{"type/Text":{"percent-json":0.0,"percent-url":0.0,"percent-email":0.0,"percent-state":1.0,"average-length":2.0}}}	5	CHARACTER	auto-list	\N	7	0	type/Text	\N	\N	f	f	f	f	\N
49	2024-06-22 08:19:17.839482+00	2024-06-22 08:23:40.124439+00	BIRTH_DATE	type/Date	\N	t	The date of birth of the user	t	9	3	\N	Birth Date	normal	\N	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":2308,"nil%":0.0},"type":{"type/DateTime":{"earliest":"1958-04-26","latest":"2000-04-03"}}}	5	DATE	\N	\N	9	0	type/Date	\N	\N	f	f	f	f	\N
52	2024-06-22 08:19:17.839482+00	2024-06-22 08:23:39.916778+00	ADDRESS	type/Text	\N	t	The street address of the account’s billing address	t	1	3	\N	Address	normal	\N	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":2490,"nil%":0.0},"type":{"type/Text":{"percent-json":0.0,"percent-url":0.0,"percent-email":0.0,"percent-state":0.0,"average-length":20.85}}}	5	CHARACTER VARYING	\N	\N	1	0	type/Text	\N	\N	f	f	f	f	\N
56	2024-06-22 08:19:17.839482+00	2024-06-22 08:23:40.225419+00	CREATED_AT	type/DateTime	type/CreationTimestamp	t	The date the user record was created. Also referred to as the user’s "join date"	t	12	3	\N	Created At	normal	\N	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":2500,"nil%":0.0},"type":{"type/DateTime":{"earliest":"2022-04-19T21:35:18.752Z","latest":"2025-04-19T14:06:27.3Z"}}}	5	TIMESTAMP	\N	\N	12	0	type/DateTime	\N	\N	f	f	f	f	\N
58	2024-06-22 08:19:17.839482+00	2024-06-22 08:23:40.325379+00	CATEGORY	type/Text	type/Category	t	The type of product, valid values include: Doohicky, Gadget, Gizmo and Widget	t	3	8	\N	Category	normal	\N	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":4,"nil%":0.0},"type":{"type/Text":{"percent-json":0.0,"percent-url":0.0,"percent-email":0.0,"percent-state":0.0,"average-length":6.375}}}	5	CHARACTER VARYING	auto-list	\N	3	0	type/Text	\N	\N	f	f	f	f	\N
64	2024-06-22 08:19:17.839482+00	2024-06-22 08:23:40.33684+00	CREATED_AT	type/DateTime	type/CreationTimestamp	t	The date the product was added to our catalog.	t	7	8	\N	Created At	normal	\N	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":200,"nil%":0.0},"type":{"type/DateTime":{"earliest":"2022-04-26T19:29:55.147Z","latest":"2025-04-15T13:34:19.931Z"}}}	5	TIMESTAMP	\N	\N	7	0	type/DateTime	\N	\N	f	f	f	f	\N
63	2024-06-22 08:19:17.839482+00	2024-06-22 08:23:40.415271+00	EAN	type/Text	\N	t	The international article number. A 13 digit number uniquely identifying the product.	t	1	8	\N	Ean	normal	\N	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":200,"nil%":0.0},"type":{"type/Text":{"percent-json":0.0,"percent-url":0.0,"percent-email":0.0,"percent-state":0.0,"average-length":13.0}}}	5	CHARACTER	auto-list	\N	1	0	type/Text	\N	\N	f	f	f	f	\N
62	2024-06-22 08:19:17.839482+00	2024-06-22 08:23:40.427393+00	ID	type/BigInteger	type/PK	t	The numerical product number. Only used internally. All external communication should use the title or EAN.	t	0	8	\N	ID	normal	\N	\N	\N	\N	\N	0	BIGINT	\N	\N	0	0	type/BigInteger	\N	\N	f	f	t	t	\N
59	2024-06-22 08:19:17.839482+00	2024-06-22 08:23:40.440654+00	PRICE	type/Float	\N	t	The list price of the product. Note that this is not always the price the product sold for due to discounts, promotions, etc.	t	5	8	\N	Price	normal	\N	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":170,"nil%":0.0},"type":{"type/Number":{"min":15.691943673970439,"q1":37.25154462926434,"q3":75.45898071609447,"max":98.81933684368194,"sd":21.711481557852057,"avg":55.74639966792074}}}	5	DOUBLE PRECISION	\N	\N	5	0	type/Float	\N	\N	f	f	f	f	\N
61	2024-06-22 08:19:17.839482+00	2024-06-22 08:23:40.523759+00	RATING	type/Float	type/Score	t	The average rating users have given the product. This ranges from 1 - 5	t	6	8	\N	Rating	normal	\N	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":23,"nil%":0.0},"type":{"type/Number":{"min":0.0,"q1":3.5120465053408525,"q3":4.216124969497314,"max":5.0,"sd":1.3605488657451452,"avg":3.4715}}}	5	DOUBLE PRECISION	\N	\N	6	0	type/Float	\N	\N	f	f	f	f	\N
65	2024-06-22 08:19:17.839482+00	2024-06-22 08:23:40.534969+00	TITLE	type/Text	type/Title	t	The name of the product as it should be displayed to customers.	t	2	8	\N	Title	normal	\N	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":199,"nil%":0.0},"type":{"type/Text":{"percent-json":0.0,"percent-url":0.0,"percent-email":0.0,"percent-state":0.0,"average-length":21.495}}}	5	CHARACTER VARYING	auto-list	\N	2	0	type/Text	\N	\N	f	f	f	f	\N
60	2024-06-22 08:19:17.839482+00	2024-06-22 08:23:40.546641+00	VENDOR	type/Text	type/Company	t	The source of the product.	t	4	8	\N	Vendor	normal	\N	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":200,"nil%":0.0},"type":{"type/Text":{"percent-json":0.0,"percent-url":0.0,"percent-email":0.0,"percent-state":0.0,"average-length":20.6}}}	5	CHARACTER VARYING	auto-list	\N	4	0	type/Text	\N	\N	f	f	f	f	\N
66	2024-06-22 08:19:17.839482+00	2024-06-22 08:23:40.713286+00	RATING	type/Integer	type/Score	t	The rating (on a scale of 1-5) the user left.	t	3	4	\N	Rating	normal	\N	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":5,"nil%":0.0},"type":{"type/Number":{"min":1.0,"q1":3.54744353181696,"q3":4.764807071650455,"max":5.0,"sd":1.0443899855660577,"avg":3.987410071942446}}}	5	SMALLINT	auto-list	\N	3	0	type/Integer	\N	\N	f	f	f	f	\N
54	2024-06-22 08:19:17.839482+00	2024-06-22 08:23:40.748766+00	ZIP	type/Text	type/ZipCode	t	The postal code of the account’s billing address	t	10	3	\N	Zip	normal	\N	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":2234,"nil%":0.0},"type":{"type/Text":{"percent-json":0.0,"percent-url":0.0,"percent-email":0.0,"percent-state":0.0,"average-length":5.0}}}	5	CHARACTER	\N	\N	10	0	type/Text	\N	\N	f	f	f	f	\N
57	2024-06-22 08:19:17.839482+00	2024-06-22 08:23:40.214408+00	LATITUDE	type/Float	type/Latitude	t	This is the latitude of the user on sign-up. It might be updated in the future to the last seen location.	t	11	3	\N	Latitude	normal	\N	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":2491,"nil%":0.0},"type":{"type/Number":{"min":25.775827,"q1":35.302705923023126,"q3":43.773802584662,"max":70.6355001,"sd":6.390832341883712,"avg":39.87934670484002}}}	5	DOUBLE PRECISION	\N	\N	11	0	type/Float	\N	\N	f	f	f	f	\N
67	2024-06-22 08:19:17.839482+00	2024-06-22 08:23:40.812915+00	BODY	type/Text	type/Description	t	The review the user left. Limited to 2000 characters.	t	4	4	\N	Body	normal	\N	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":1112,"nil%":0.0},"type":{"type/Text":{"percent-json":0.0,"percent-url":0.0,"percent-email":0.0,"percent-state":0.0,"average-length":177.41996402877697}}}	5	CHARACTER VARYING	\N	\N	4	0	type/Text	\N	\N	f	f	f	f	\N
44	2024-06-22 08:19:17.839482+00	2024-06-22 08:23:39.510747+00	SUBTOTAL	type/Float	\N	t	The raw, pre-tax cost of the order. Note that this might be different in the future from the product price due to promotions, credits, etc.	t	3	5	\N	Subtotal	normal	\N	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":340,"nil%":0.0},"type":{"type/Number":{"min":15.691943673970439,"q1":49.74894519060184,"q3":105.42965746993103,"max":148.22900526552291,"sd":32.53705013056317,"avg":77.01295465356547}}}	5	DOUBLE PRECISION	\N	\N	3	0	type/Float	\N	\N	f	f	f	f	\N
50	2024-06-22 08:19:17.839482+00	2024-06-22 08:23:40.040342+00	LONGITUDE	type/Float	type/Longitude	t	This is the longitude of the user on sign-up. It might be updated in the future to the last seen location.	t	6	3	\N	Longitude	normal	\N	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":2491,"nil%":0.0},"type":{"type/Number":{"min":-166.5425726,"q1":-101.58350792373135,"q3":-84.65289348288829,"max":-67.96735199999999,"sd":15.399698968175663,"avg":-95.18741780363999}}}	5	DOUBLE PRECISION	\N	\N	6	0	type/Float	\N	\N	f	f	f	f	\N
69	2024-06-22 08:19:17.839482+00	2024-06-22 08:23:40.637221+00	CREATED_AT	type/DateTime	type/CreationTimestamp	t	The day and time a review was written by a user.	t	5	4	\N	Created At	normal	\N	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":1112,"nil%":0.0},"type":{"type/DateTime":{"earliest":"2022-06-03T00:37:05.818Z","latest":"2026-04-19T14:15:25.677Z"}}}	5	TIMESTAMP	\N	\N	5	0	type/DateTime	\N	\N	f	f	f	f	\N
68	2024-06-22 08:19:17.839482+00	2024-06-22 08:23:40.647553+00	ID	type/BigInteger	type/PK	t	A unique internal identifier for the review. Should not be used externally.	t	0	4	\N	ID	normal	\N	\N	\N	\N	\N	0	BIGINT	\N	\N	0	0	type/BigInteger	\N	\N	f	f	t	t	\N
71	2024-06-22 08:19:17.839482+00	2024-06-22 08:23:40.678511+00	PRODUCT_ID	type/Integer	type/FK	t	The product the review was for	t	1	4	\N	Product ID	normal	62	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":176,"nil%":0.0}}	5	INTEGER	\N	\N	1	0	type/Integer	\N	\N	f	f	f	t	\N
70	2024-06-22 08:19:17.839482+00	2024-06-22 08:23:40.727832+00	REVIEWER	type/Text	\N	t	The user who left the review	t	2	4	\N	Reviewer	normal	\N	2024-06-22 08:19:17.839482+00	\N	\N	{"global":{"distinct-count":1076,"nil%":0.0},"type":{"type/Text":{"percent-json":0.0,"percent-url":0.0,"percent-email":0.0,"percent-state":0.001798561151079137,"average-length":9.972122302158274}}}	5	CHARACTER VARYING	\N	\N	2	0	type/Text	\N	\N	f	f	f	f	\N
\.


--
-- Data for Name: metabase_fieldvalues; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.metabase_fieldvalues (id, created_at, updated_at, "values", human_readable_values, field_id, has_more_values, type, hash_key, last_used_at) FROM stdin;
1	2024-06-22 08:20:32.254468+00	2024-06-22 08:20:32.254468+00	[1,3,5,8,10,12,14,16,19,21,23,25,27,30,32,34,36,38,41,43,45,47,49,52,54,56,58,60,63,65,67,69,71,74,76,78,80,82,85,87,89,91,93,96,98,102,104,107,109,111,113,115,120,122,124,129,131,133,135,137,140,142,148,155,162,168,173,175,177,179,186,188,190,197,201,203,208,210,219,223,236,252,254,261,267,269,274,289,331,335,351,390,393,401,423,426,450,503,522,639,668,1325]	\N	1	f	full	\N	2024-06-22 08:20:32.254468+00
2	2024-06-22 08:20:32.364801+00	2024-06-22 08:20:32.364801+00	["Basic","Business","Premium"]	\N	4	f	full	\N	2024-06-22 08:20:32.364801+00
3	2024-06-22 08:20:32.550293+00	2024-06-22 08:20:32.550293+00	[false,true]	\N	5	f	full	\N	2024-06-22 08:20:32.550293+00
4	2024-06-22 08:20:32.654551+00	2024-06-22 08:20:32.654551+00	[false,true]	\N	8	f	full	\N	2024-06-22 08:20:32.654551+00
5	2024-06-22 08:20:32.767079+00	2024-06-22 08:20:32.767079+00	["Abbott","Abernathy","Abshire","Adams","Altenwerth","Anderson","Ankunding","Armstrong","Auer","Aufderhar","Bahringer","Bailey","Balistreri","Barrows","Bartell","Bartoletti","Barton","Batz","Bauch","Baumbach","Bayer","Beahan","Beatty","Bechtelar","Becker","Bednar","Beer","Beier","Berge","Bergnaum","Bergstrom","Bernhard","Bernier","Bins","Blanda","Blick","Block","Bode","Boehm","Bogan","Bogisich","Borer","Bosco","Botsford","Boyer","Boyle","Bradtke","Brakus","Braun","Breitenberg","Brekke","Brown","Bruen","Buckridge","Carroll","Carter","Cartwright","Casper","Cassin","Champlin","Christiansen","Cole","Collier","Collins","Conn","Connelly","Conroy","Considine","Corkery","Cormier","Corwin","Cremin","Crist","Crona","Cronin","Crooks","Cruickshank","Cummerata","Cummings","D'Amore","Dach","Daniel","Dare","Daugherty","Davis","Deckow","Denesik","Dibbert","Dickens","Dicki","Dickinson","Dietrich","Donnelly","Dooley","Douglas","Doyle","DuBuque","Durgan","Ebert","Effertz","Eichmann","Emard","Emmerich","Erdman","Ernser","Fadel","Fahey","Farrell","Fay","Feeney","Feest","Feil","Ferry","Fisher","Flatley","Frami","Franecki","Friesen","Fritsch","Funk","Gaylord","Gerhold","Gerlach","Gibson","Gislason","Gleason","Gleichner","Glover","Goldner","Goodwin","Gorczany","Gottlieb","Goyette","Grady","Graham","Grant","Green","Greenfelder","Greenholt","Grimes","Gulgowski","Gusikowski","Gutkowski","Gutmann","Haag","Hackett","Hagenes","Hahn","Haley","Halvorson","Hamill","Hammes","Hand","Hane","Hansen","Harber","Harris","Hartmann","Harvey","Hauck","Hayes","Heaney","Heathcote","Hegmann","Heidenreich","Heller","Herman","Hermann","Hermiston","Herzog","Hessel","Hettinger","Hickle","Hilll","Hills","Hilpert","Hintz","Hirthe","Hodkiewicz","Hoeger","Homenick","Hoppe","Howe","Howell","Hudson","Huel","Huels","Hyatt","Jacobi","Jacobs","Jacobson","Jakubowski","Jaskolski","Jast","Jenkins","Jerde","Jewess","Johns","Johnson","Johnston","Jones","Kassulke","Kautzer","Keebler","Keeling","Kemmer","Kerluke","Kertzmann","Kessler","Kiehn","Kihn","Kilback","King","Kirlin","Klein","Kling","Klocko","Koch","Koelpin","Koepp","Kohler","Konopelski","Koss","Kovacek","Kozey","Krajcik","Kreiger","Kris","Kshlerin","Kub","Kuhic","Kuhlman","Kuhn","Kulas","Kunde","Kunze","Kuphal","Kutch","Kuvalis","Labadie","Lakin","Lang","Langosh","Langworth","Larkin","Larson","Leannon","Lebsack","Ledner","Leffler","Legros","Lehner","Lemke","Lesch","Leuschke","Lind","Lindgren","Littel","Little","Lockman","Lowe","Lubowitz","Lueilwitz","Luettgen","Lynch","Macejkovic","Maggio","Mann","Mante","Marks","Marquardt","Marvin","Mayer","Mayert","McClure","McCullough","McDermott","McGlynn","McKenzie","McLaughlin","Medhurst","Mertz","Metz","Miller","Mills","Mitchell","Moen","Mohr","Monahan","Moore","Morar","Morissette","Mosciski","Mraz","Mueller","Muller","Murazik","Murphy","Murray","Nader","Nicolas","Nienow","Nikolaus","Nitzsche","Nolan","O'Connell","O'Conner","O'Hara","O'Keefe","O'Kon","O'Reilly","Oberbrunner","Okuneva","Olson","Ondricka","Orn","Ortiz","Osinski","Pacocha","Padberg","Pagac","Parisian","Parker","Paucek","Pfannerstill","Pfeffer","Pollich","Pouros","Powlowski","Predovic","Price","Prohaska","Prosacco","Purdy","Quigley","Quitzon","Rath","Ratke","Rau","Raynor","Reichel","Reichert","Reilly","Reinger","Rempel","Renner","Reynolds","Rice","Rippin","Ritchie","Robel","Roberts","Rodriguez","Rogahn","Rohan","Rolfson","Romaguera","Roob","Rosenbaum","Rowe","Ruecker","Runolfsdottir","Runolfsson","Runte","Russel","Rutherford","Ryan","Sanford","Satterfield","Sauer","Sawayn","Schaden","Schaefer","Schamberger","Schiller","Schimmel","Schinner","Schmeler","Schmidt","Schmitt","Schneider","Schoen","Schowalter","Schroeder","Schulist","Schultz","Schumm","Schuppe","Schuster","Senger","Shanahan","Shields","Simonis","Sipes","Skiles","Smith","Smitham","Spencer","Spinka","Sporer","Stamm","Stanton","Stark","Stehr","Steuber","Stiedemann","Stokes","Stoltenberg","Stracke","Streich","Stroman","Strosin","Swaniawski","Swift","Terry","Thiel","Thompson","Tillman","Torp","Torphy","Towne","Toy","Trantow","Tremblay","Treutel","Tromp","Turcotte","Turner","Ullrich","Upton","Vandervort","Veum","Volkman","Von","VonRueden","Waelchi","Walker","Walsh","Walter","Ward","Waters","Watsica","Weber","Wehner","Weimann","Weissnat","Welch","West","White","Wiegand","Wilderman","Wilkinson","Will","Williamson","Willms","Windler","Wintheiser","Wisoky","Wisozk","Witting","Wiza","Wolf","Wolff","Wuckert","Wunsch","Wyman","Yost","Yundt","Zboncak","Zemlak","Ziemann","Zieme","Zulauf"]	\N	10	f	full	\N	2024-06-22 08:20:32.767079+00
6	2024-06-22 08:20:32.958518+00	2024-06-22 08:20:32.958518+00	[null,"Facebook","Google","Invite","Twitter"]	\N	11	f	full	\N	2024-06-22 08:20:32.958518+00
7	2024-06-22 08:20:33.061378+00	2024-06-22 08:20:33.061378+00	[false,true]	\N	14	f	full	\N	2024-06-22 08:20:33.061378+00
8	2024-06-22 08:20:33.241377+00	2024-06-22 08:20:33.241377+00	[null,"AE","AF","AG","AL","AM","AR","AT","AU","BA","BD","BE","BF","BG","BN","BO","BR","BT","BW","BY","CA","CD","CH","CI","CL","CM","CN","CO","CR","CU","CV","CY","CZ","DE","DK","DO","DZ","EE","EG","ES","ET","FI","FR","GB","GE","GM","GN","GR","GT","HN","HR","HT","HU","ID","IE","IL","IN","IQ","IR","IT","JM","JO","JP","KE","KH","KI","KM","KR","KZ","LA","LC","LK","LR","LS","LT","LU","LV","LY","MA","MD","MG","MK","ML","MM","MT","MU","MW","MX","MY","NE","NG","NI","NL","NO","NZ","PA","PE","PH","PK","PL","PT","PW","PY","RO","RS","RU","RW","SA","SE","SI","SK","SL","SM","SN","SO","SV","SY","SZ","TH","TJ","TN","TO","TR","TZ","UA","UG","US","UZ","VE","VN","YE","ZA","ZM","ZW"]	\N	15	f	full	\N	2024-06-22 08:20:33.241377+00
9	2024-06-22 08:20:33.362192+00	2024-06-22 08:20:33.362192+00	["Button Clicked","Page Viewed"]	\N	17	f	full	\N	2024-06-22 08:20:33.362192+00
10	2024-06-22 08:20:33.552487+00	2024-06-22 08:20:33.552487+00	[null,"www.piespace.example/help","www.piespace.example/home","www.piespace.example/invite","www.piespace.example/login","www.piespace.example/pies"]	\N	18	f	full	\N	2024-06-22 08:20:33.552487+00
11	2024-06-22 08:20:33.660108+00	2024-06-22 08:20:33.660108+00	[null,"Checkout","Create Item","Invite","Signup","Subscribe"]	\N	19	f	full	\N	2024-06-22 08:20:33.660108+00
12	2024-06-22 08:23:44.822131+00	2024-06-22 08:23:44.822131+00	[1,2,3,4,5]	\N	23	f	full	\N	2024-06-22 08:23:44.822131+00
13	2024-06-22 08:23:45.012191+00	2024-06-22 08:23:45.012191+00	["abbott-berneice@hotmail.example","abdullah-kerluke@gmail.example","adan-weissnat@yahoo.example","aida.schneider@gmail.example","aidan-hagenes@hotmail.example","aidan.rodriguez@gmail.example","alaina-howell@gmail.example","alayna.halvorson@hotmail.example","alberto.gulgowski@gmail.example","alicia.schimmel@gmail.example","alisa-morissette@hotmail.example","alisa-schmitt@hotmail.example","altenwerth.onie@yahoo.example","alva.conroy@hotmail.example","alvena-legros@hotmail.example","alverta-rogahn@hotmail.example","alvina.mertz@gmail.example","alvis.emmerich@yahoo.example","alycia.collins@yahoo.example","alysson-cartwright@yahoo.example","anastacio.jaskolski@hotmail.example","anderson-eliza@hotmail.example","anderson.schinner@yahoo.example","andy-skiles@yahoo.example","angela-botsford@hotmail.example","anjali-parker@hotmail.example","ankunding-rudolph@hotmail.example","aracely.jenkins@gmail.example","arch-ryan@hotmail.example","archibald-lowe@hotmail.example","archibald-turner@hotmail.example","arne-o-hara@gmail.example","art-graham@yahoo.example","arvel-lakin@gmail.example","ashton-herman@hotmail.example","aubree-dibbert@hotmail.example","aubree-walter@hotmail.example","aufderhar-mya@hotmail.example","aufderhar.john@yahoo.example","aurore-yundt@yahoo.example","bahringer-laura@yahoo.example","bailey.kenna@yahoo.example","balistreri-oral@yahoo.example","balistreri-unique@gmail.example","bauch-wilford@gmail.example","bayer-mattie@hotmail.example","bayer.mark@hotmail.example","beatty-emmie@gmail.example","beatty.julio@gmail.example","beatty.mohammed@gmail.example","bechtelar.antone@gmail.example","beer.humberto@hotmail.example","berge-halie@hotmail.example","bergstrom-chelsie@yahoo.example","bernhard.kathleen@yahoo.example","bins-evans@hotmail.example","blair.heaney@gmail.example","blake-leffler@hotmail.example","blaze-daugherty@hotmail.example","blick-candelario@yahoo.example","block.emiliano@hotmail.example","bode-sydnie@gmail.example","bode.richmond@gmail.example","boehm-amanda@hotmail.example","bogan.rodger@gmail.example","bosco-zachariah@gmail.example","bosco.haylie@hotmail.example","botsford.okey@hotmail.example","boyer-bernhard@yahoo.example","boyle-christiana@gmail.example","brakus-kimberly@gmail.example","brakus.marlene@hotmail.example","brant.klein@yahoo.example","braun.madisyn@gmail.example","breanna.strosin@yahoo.example","breitenberg-louie@gmail.example","brekke.kirsten@yahoo.example","brennon-gerlach@hotmail.example","bret-quigley@gmail.example","brianne-jacobson@yahoo.example","brown-deontae@gmail.example","buddy-hills@gmail.example","caleigh-hodkiewicz@yahoo.example","camron-homenick@gmail.example","camryn-schmeler@hotmail.example","candida-turcotte@yahoo.example","carolanne-upton@gmail.example","carroll.chanel@yahoo.example","carroll.kohler@hotmail.example","carter-fern@hotmail.example","casey.robel@yahoo.example","casper-alfonzo@yahoo.example","cassin-cleta@hotmail.example","cassin.mario@hotmail.example","cayla.vonrueden@hotmail.example","cecilia.stark@hotmail.example","cedrick-kessler@gmail.example","champlin.jensen@yahoo.example","chanel.rippin@yahoo.example","charlene-bayer@hotmail.example","chet-blick@yahoo.example","christophe.wilderman@hotmail.example","ciara-larson@hotmail.example","ciara.green@yahoo.example","clark-luettgen@gmail.example","claudie-dare@yahoo.example","claudine.mccullough@yahoo.example","clay-pfannerstill@gmail.example","clemens.hansen@gmail.example","cole.christophe@yahoo.example","conn-gideon@yahoo.example","connell-o-henriette@yahoo.example","connell.lisette.o@yahoo.example","connelly-alice@gmail.example","connelly.bessie@hotmail.example","conner.windler@gmail.example","conroy-orlando@yahoo.example","conroy-yadira@gmail.example","corbin.mertz@hotmail.example","corbin.wiegand@yahoo.example","corkery.theresa@yahoo.example","cornelius-bogisich@hotmail.example","crawford.rath@gmail.example","cremin-jerome@hotmail.example","cremin.tyler@gmail.example","cronin-marley@yahoo.example","curtis.morar@hotmail.example","d-amore-geoffrey@yahoo.example","dagmar-sawayn@gmail.example","dana-orn@hotmail.example","dana.kozey@yahoo.example","darwin-abshire@yahoo.example","dawson-kuvalis@gmail.example","dax-bartell@hotmail.example","dayne.strosin@hotmail.example","deckow.alisha@hotmail.example","dell-schimmel@hotmail.example","demario-hand@yahoo.example","demetris.hauck@hotmail.example","dena-schiller@yahoo.example","denesik-delphia@hotmail.example","deron-cremin@gmail.example","destiny-murazik@hotmail.example","deven.brekke@gmail.example","domenico.bailey@yahoo.example","dominic.jacobi@yahoo.example","donavon.lowe@gmail.example","dooley-karen@gmail.example","douglas-prosacco@hotmail.example","douglas.anais@hotmail.example","durgan-emiliano@hotmail.example","earnestine-lockman@hotmail.example","easton-koch@gmail.example","effertz-elnora@yahoo.example","eileen-mayert@gmail.example","eldon.herman@yahoo.example","elisa-grady@yahoo.example","ellie-oberbrunner@yahoo.example","ellsworth.west@hotmail.example","elmo.schimmel@yahoo.example","elsa.klocko@gmail.example","elvera.lowe@yahoo.example","elwin.okuneva@gmail.example","emanuel-corwin@gmail.example","emard-janiya@gmail.example","emerson-o-keefe@gmail.example","emery.gerlach@hotmail.example","emmie-mertz@yahoo.example","enola.bayer@yahoo.example","erich.kris@gmail.example","ernestina-gerhold@gmail.example","ernser-ardella@gmail.example","esther-douglas@yahoo.example","estrella.goyette@hotmail.example","ethan.rutherford@gmail.example","eudora-renner@gmail.example","eugenia-stroman@yahoo.example","eula-connell-o@hotmail.example","eve.mante@yahoo.example","fadel-philip@hotmail.example","feest-angus@gmail.example","feil.sterling@gmail.example","felicity-greenfelder@hotmail.example","felipe-johnston@yahoo.example","ferry.enrico@hotmail.example","fisher-antwan@hotmail.example","florence.donnelly@gmail.example","foster-gusikowski@yahoo.example","foster-marks@yahoo.example","francisco-robel@hotmail.example","freddie.wisoky@gmail.example","fredrick-gulgowski@yahoo.example","fritz.dickens@hotmail.example","funk.nichole@yahoo.example","gabrielle-considine@yahoo.example","gabrielle-frami@gmail.example","gaetano-rogahn@hotmail.example","gaylord-granville@yahoo.example","gene-lueilwitz@gmail.example","gerhold.lempi@yahoo.example","germaine-brakus@yahoo.example","gibson.eveline@hotmail.example","gilberto-mueller@gmail.example","gino.johnston@yahoo.example","giovani-thompson@hotmail.example","giovani.lesch@hotmail.example","gislason-kaelyn@hotmail.example","giuseppe.morar@hotmail.example","gleichner-joshuah@gmail.example","glover-eryn@gmail.example","glover.kelsie@yahoo.example","goldner.ruthe@hotmail.example","gorczany-eulah@yahoo.example","gottlieb-ola@gmail.example","gottlieb-ruthe@yahoo.example","grady.raynor@hotmail.example","graham-liam@yahoo.example","greenfelder-hulda@gmail.example","greenfelder.wilbert@gmail.example","greg-purdy@gmail.example","gretchen.muller@hotmail.example","greyson.boyle@gmail.example","grimes-terrence@gmail.example","grimes.melisa@gmail.example","guido-mckenzie@yahoo.example","gulgowski.ubaldo@gmail.example","gutkowski-pattie@hotmail.example","gutmann-lura@hotmail.example","hagenes-rosie@hotmail.example","hahn.hugh@hotmail.example","halvorson.dale@gmail.example","hane.audie@gmail.example","hane.carter@gmail.example","hansen-karl@hotmail.example","hansen.alta@yahoo.example","hansen.anibal@yahoo.example","hansen.magnolia@gmail.example","harris-myrtice@gmail.example","harris.constantin@hotmail.example","harris.richard@yahoo.example","heath-dare@gmail.example","heathcote.jamar@gmail.example","heidenreich-patience@hotmail.example","heidenreich.pearlie@hotmail.example","heidi-glover@gmail.example","henry-rowe@yahoo.example","hermann-madelyn@hotmail.example","hermiston.gerald@yahoo.example","hertha.price@gmail.example","herzog-ophelia@hotmail.example","hessel.arnoldo@yahoo.example","hettinger-brendon@yahoo.example","hettinger.david@yahoo.example","hettinger.orval@hotmail.example","hettinger.shyanne@yahoo.example","hills-violet@hotmail.example","hilpert.gunnar@yahoo.example","hollis-hettinger@gmail.example","homenick-omari@hotmail.example","hoppe.kathryne@gmail.example","hoppe.lewis@gmail.example","howell-reba@gmail.example","howell.jacinthe@yahoo.example","hudson-audra@hotmail.example","hudson.larkin@hotmail.example","huels-earnest@hotmail.example","huels-gunnar@yahoo.example","hyatt.rowan@hotmail.example","jace-kihn@yahoo.example","jacey.schoen@hotmail.example","jacobs-oliver@yahoo.example","jacobs-ronny@hotmail.example","jacobson.stan@gmail.example","jairo-simonis@yahoo.example","jakob.hansen@gmail.example","jakubowski.nyasia@yahoo.example","jaleel.collins@gmail.example","jamel.stanton@yahoo.example","jannie-balistreri@yahoo.example","jasen.stanton@yahoo.example","jast.leann@yahoo.example","jayden.kris@hotmail.example","jazmin.brekke@gmail.example","jedediah-huels@hotmail.example","jeffry-schowalter@hotmail.example","jenkins-sandy@yahoo.example","jennifer-klocko@hotmail.example","jerrod-king@yahoo.example","jessika.funk@yahoo.example","jo-gusikowski@hotmail.example","joe.becker@yahoo.example","joelle-ullrich@gmail.example","johns-myrtle@yahoo.example","johnston.benny@gmail.example","jorge.bins@gmail.example","josh-schimmel@gmail.example","judd-hickle@gmail.example","kade-kub@hotmail.example","kaela-kunze@gmail.example","kariane.hintz@gmail.example","karine.mante@gmail.example","katharina-heathcote@yahoo.example","kavon-dach@gmail.example","kaya.schoen@yahoo.example","kayley.powlowski@hotmail.example","keefe-o-jonas@yahoo.example","keenan.ferry@hotmail.example","kellie.price@hotmail.example","kelsi.douglas@gmail.example","kemmer-gene@yahoo.example","kemmer-matt@yahoo.example","kemmer.bonnie@hotmail.example","kennedy-kunde@gmail.example","kerluke.jakob@gmail.example","kertzmann-coty@hotmail.example","keshaun-carroll@hotmail.example","khalid-pouros@yahoo.example","khalid.blanda@yahoo.example","kihn.alfred@yahoo.example","kilback-alisha@hotmail.example","kilback-carmelo@yahoo.example","kitty.hilll@gmail.example","koelpin-karelle@gmail.example","koepp-melyna@gmail.example","kohler.jermain@hotmail.example","konopelski.beaulah@hotmail.example","koss-ella@hotmail.example","koss.letha@hotmail.example","kovacek-dawson@hotmail.example","kristoffer.blanda@yahoo.example","krystel.boyle@yahoo.example","kshlerin-bernardo@gmail.example","kshlerin-stella@yahoo.example","kulas-armani@yahoo.example","kunze.eleanora@hotmail.example","kuphal.colton@gmail.example","kurtis.parker@gmail.example","kuvalis-cierra@yahoo.example","kuvalis-willis@yahoo.example","kyler-altenwerth@yahoo.example","kyler.abshire@yahoo.example","kyra-lynch@hotmail.example","lacey.dickinson@hotmail.example","langosh.cathrine@gmail.example","langworth-savion@gmail.example","larkin-lilliana@hotmail.example","larkin.cedrick@gmail.example","larson-adrianna@gmail.example","laurel.pfannerstill@gmail.example","laurie-sanford@hotmail.example","lavern.botsford@hotmail.example","lavern.boyle@hotmail.example","leannon-clay@gmail.example","lebsack-tristin@yahoo.example","ledner-nichole@yahoo.example","leila-considine@gmail.example","leta-heidenreich@yahoo.example","leuschke-estefania@yahoo.example","liam-schoen@hotmail.example","lillie.wilderman@hotmail.example","lind.annamae@gmail.example","linnea.dickens@gmail.example","littel.otto@hotmail.example","little-anika@hotmail.example","little.john@yahoo.example","lockman-janiya@hotmail.example","logan-weber@yahoo.example","lon-friesen@yahoo.example","lorna.greenholt@yahoo.example","lou.runte@gmail.example","lowell-daniel@gmail.example","loyal.wintheiser@yahoo.example","loyce-lemke@yahoo.example","lucas-beer@gmail.example","lucile-bednar@gmail.example","lueilwitz.osbaldo@yahoo.example","luna-nienow@gmail.example","lynch.tyson@gmail.example","lysanne-brekke@yahoo.example","mabel-grimes@hotmail.example","macejkovic-cyrus@gmail.example","macejkovic.andrew@hotmail.example","mackenzie-ullrich@yahoo.example","madge-friesen@gmail.example","madie.bayer@hotmail.example","maeve.hilpert@hotmail.example","maiya-beier@hotmail.example","malika-kuphal@hotmail.example","mante-dakota@yahoo.example","marcelle-rippin@gmail.example","marcelo-ferry@gmail.example","margarete.tillman@gmail.example","marley.gorczany@hotmail.example","marvin.kris@gmail.example","mathilde.quigley@gmail.example","maurine-considine@hotmail.example","maximillia.ebert@hotmail.example","maximillian-zboncak@hotmail.example","mayer.arne@hotmail.example","mayert.jessyca@gmail.example","mckenzie.eduardo@yahoo.example","mclaughlin.ezekiel@yahoo.example","meagan.cremin@yahoo.example","melba-witting@hotmail.example","melisa.hilpert@hotmail.example","melissa.cormier@hotmail.example","merle.blick@yahoo.example","merle.moen@yahoo.example","mertz.antoinette@hotmail.example","mertz.melissa@gmail.example","micaela.kerluke@gmail.example","milan-ritchie@gmail.example","miller-geovanni@hotmail.example","miller.morgan@gmail.example","mills-andy@hotmail.example","milton.schiller@hotmail.example","mina.reynolds@yahoo.example","miracle.erdman@gmail.example","misty-botsford@hotmail.example","mitchell-lacey@gmail.example","moen-evalyn@gmail.example","mohr.johnson@gmail.example","mollie.bogan@hotmail.example","monahan.loma@yahoo.example","monserrate-doyle@hotmail.example","morar-maddison@hotmail.example","morissette.jailyn@gmail.example","mraz-tomas@gmail.example","mraz.caitlyn@yahoo.example","muller-russell@yahoo.example","murazik-donny@hotmail.example","murray-idell@gmail.example","murray-zemlak@hotmail.example","murray.gleason@gmail.example","mya-gleason@yahoo.example","myles.deckow@gmail.example","myrtle.bahringer@hotmail.example","nader-arnaldo@hotmail.example","nader-ryley@hotmail.example","nayeli.becker@yahoo.example","nicolas.dameon@yahoo.example","nicolas.karen@hotmail.example","nikko.bartoletti@gmail.example","nikolas-hilpert@gmail.example","nikolaus-willie@yahoo.example","nils.gaylord@hotmail.example","nolan-amy@hotmail.example","nolan.samantha@hotmail.example","o-issac-kon@hotmail.example","o.janelle.hara@gmail.example","oceane.runte@yahoo.example","odell.stehr@hotmail.example","olaf.sipes@gmail.example","ollie.corkery@gmail.example","ondricka-lamont@hotmail.example","ondricka-madge@yahoo.example","ondricka.rollin@gmail.example","orie-sipes@yahoo.example","ortiz.harrison@gmail.example","orville-effertz@hotmail.example","oscar-olson@hotmail.example","osinski-joanne@gmail.example","pacocha-khalil@hotmail.example","padberg-albert@gmail.example","pagac-yessenia@yahoo.example","parker-lilliana@yahoo.example","pattie.senger@yahoo.example","paucek-larry@hotmail.example","petra.durgan@hotmail.example","peyton-barton@gmail.example","powlowski-mohammed@yahoo.example","price-rosalyn@yahoo.example","rau.arnaldo@hotmail.example","raynor.chasity@gmail.example","rebekah-dickinson@hotmail.example","rebekah.ledner@gmail.example","reichel-antwon@hotmail.example","reichel.gracie@gmail.example","reichert-evangeline@yahoo.example","reid-reilly@hotmail.example","reilly.o.franco@yahoo.example","rempel.brooke@gmail.example","rene.muller@gmail.example","reuben-koelpin@yahoo.example","rey-schumm@hotmail.example","reyna-greenholt@yahoo.example","reynolds-melisa@hotmail.example","richmond-adams@gmail.example","roberts-lilian@yahoo.example","rogahn-meta@hotmail.example","rolfson-ford@yahoo.example","rolfson.natalie@yahoo.example","romaguera-angeline@yahoo.example","roob-lila@gmail.example","rosalinda-stamm@gmail.example","rosella-bergstrom@gmail.example","rowe-celestine@gmail.example","ruecker-kathlyn@yahoo.example","ruecker-tad@yahoo.example","runolfsdottir-tyreek@gmail.example","runolfsdottir.augustine@gmail.example","runolfsson.davonte@gmail.example","rutherford-beau@yahoo.example","rylee-upton@yahoo.example","ryleigh-padberg@hotmail.example","sabrina-schmidt@hotmail.example","sabryna-schumm@gmail.example","sallie.wehner@yahoo.example","samir.hayes@yahoo.example","sanford.leilani@yahoo.example","satterfield-abbey@hotmail.example","satterfield.chris@hotmail.example","satterfield.creola@hotmail.example","sauer-franco@yahoo.example","schaden.johathan@hotmail.example","schamberger-zora@gmail.example","schiller.loyal@gmail.example","schinner.verna@yahoo.example","schmeler-lucinda@hotmail.example","schmeler-rita@yahoo.example","schmeler.annabelle@gmail.example","schoen.viola@hotmail.example","schuster.geovanny@yahoo.example","scottie.schmidt@hotmail.example","senger-kamron@hotmail.example","shanie.spinka@gmail.example","shanny-kuvalis@gmail.example","shirley-okuneva@yahoo.example","sidney-kling@gmail.example","skiles-devan@gmail.example","skye.heidenreich@gmail.example","smith-meaghan@yahoo.example","smith-price@yahoo.example","spencer.efrain@gmail.example","spinka-donato@yahoo.example","spinka-jessy@gmail.example","sporer-nyasia@hotmail.example","sporer.toby@yahoo.example","stamm.davin@yahoo.example","stanley.kuphal@gmail.example","stark.yasmin@gmail.example","stehr-freeman@hotmail.example","stehr.tyson@gmail.example","steuber-vernice@yahoo.example","steuber.dedrick@gmail.example","stewart.sawayn@gmail.example","stiedemann-gage@gmail.example","stiedemann.coby@gmail.example","stokes.cordelia@yahoo.example","stoltenberg-miguel@gmail.example","swaniawski-kaleb@gmail.example","swaniawski-luisa@hotmail.example","swift.pietro@hotmail.example","tavares.metz@yahoo.example","terry-gregorio@gmail.example","terry-william@yahoo.example","terry.darlene@gmail.example","terry.joe@gmail.example","theodore.mcglynn@hotmail.example","theresa.grant@yahoo.example","theresia-russel@gmail.example","thompson-fay@yahoo.example","thompson-hosea@hotmail.example","thurman-pouros@hotmail.example","toby.yundt@gmail.example","torp.magdalen@yahoo.example","trantow-daphnee@hotmail.example","tressie-smitham@yahoo.example","treutel-jessika@yahoo.example","treutel.jaquan@hotmail.example","tromp-demario@yahoo.example","tromp.emelia@hotmail.example","trudie-muller@gmail.example","trudie.koch@yahoo.example","turner.dahlia@yahoo.example","turner.kelley@hotmail.example","tyrel-beatty@yahoo.example","ullrich-ladarius@gmail.example","unique.jerde@yahoo.example","vanessa-jaskolski@gmail.example","vergie.borer@hotmail.example","veronica.weissnat@gmail.example","virginia-prohaska@yahoo.example","vonrueden-sheridan@hotmail.example","waelchi-filomena@yahoo.example","waelchi-jaqueline@hotmail.example","waelchi.iva@yahoo.example","walker-derrick@yahoo.example","walker-nicole@yahoo.example","walker-trace@gmail.example","walker.carter@yahoo.example","walter-melisa@gmail.example","walter.chris@yahoo.example","ward.anais@gmail.example","ward.mabel@gmail.example","warren.gulgowski@hotmail.example","watsica-olen@gmail.example","watsica.stanley@hotmail.example","weber-breana@hotmail.example","weber-gino@yahoo.example","weber-mohammad@yahoo.example","weimann.keyshawn@hotmail.example","weimann.maryam@gmail.example","weissnat-victoria@gmail.example","weissnat.elmore@gmail.example","weissnat.mathilde@hotmail.example","welch-lucinda@yahoo.example","wendell-becker@yahoo.example","west.laisha@yahoo.example","white.berneice@gmail.example","wiegand.guy@gmail.example","wilderman-nellie@hotmail.example","wilhelmine.erdman@gmail.example","wilkinson-penelope@gmail.example","wilkinson.edmund@hotmail.example","will.garrison@yahoo.example","willms-ardella@yahoo.example","willms-tressie@hotmail.example","willms.seth@gmail.example","willms.wilhelm@gmail.example","wilton-senger@yahoo.example","winfield.donnelly@hotmail.example","winona.cassin@yahoo.example","wintheiser-broderick@yahoo.example","wintheiser-celestino@gmail.example","wintheiser-murray@yahoo.example","wisoky-rebeka@hotmail.example","witting-cindy@yahoo.example","witting-maud@gmail.example","witting-raegan@gmail.example","wiza-andreanne@gmail.example","wiza.lisette@yahoo.example","wolf-jewell@yahoo.example","wuckert.iva@hotmail.example","wyman.hilma@gmail.example","yundt-haven@hotmail.example","yundt.merl@yahoo.example","zackery.bailey@gmail.example","zane.paucek@yahoo.example","zetta.nitzsche@gmail.example","ziemann-serena@gmail.example","zula.boehm@hotmail.example"]	\N	26	f	full	\N	2024-06-22 08:23:45.012191+00
14	2024-06-22 08:23:45.125602+00	2024-06-22 08:23:45.125602+00	["Average","Below Average","Good","Great","Poor"]	\N	28	f	full	\N	2024-06-22 08:23:45.125602+00
15	2024-06-22 08:23:45.229527+00	2024-06-22 08:23:45.229527+00	["Basic","Business","Premium"]	\N	30	f	full	\N	2024-06-22 08:23:45.229527+00
16	2024-06-22 08:23:45.328885+00	2024-06-22 08:23:45.328885+00	[false,true]	\N	35	f	full	\N	2024-06-22 08:23:45.328885+00
17	2024-06-22 08:23:45.514248+00	2024-06-22 08:23:45.514248+00	[0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26,27,28,29,30,31,32,33,34,35,36,37,38,39,40,41,42,43,44,45,46,47,48,49,50,51,52,53,54,55,56,57,58,59,60,61,63,65,67,68,69,70,71,72,73,75,78,82,83,88,100]	\N	39	f	full	\N	2024-06-22 08:23:45.514248+00
18	2024-06-22 08:23:45.623908+00	2024-06-22 08:23:45.623908+00	["Affiliate","Facebook","Google","Organic","Twitter"]	\N	45	f	full	\N	2024-06-22 08:23:45.623908+00
19	2024-06-22 08:23:45.727552+00	2024-06-22 08:23:45.727552+00	["AK","AL","AR","AZ","CA","CO","CT","DE","FL","GA","IA","ID","IL","IN","KS","KY","LA","MA","MD","ME","MI","MN","MO","MS","MT","NC","ND","NE","NH","NJ","NM","NV","NY","OH","OK","OR","PA","RI","SC","SD","TN","TX","UT","VA","VT","WA","WI","WV","WY"]	\N	48	f	full	\N	2024-06-22 08:23:45.727552+00
20	2024-06-22 08:23:45.827777+00	2024-06-22 08:23:45.827777+00	["Doohickey","Gadget","Gizmo","Widget"]	\N	58	f	full	\N	2024-06-22 08:23:45.827777+00
21	2024-06-22 08:23:45.930738+00	2024-06-22 08:23:45.930738+00	["0001664425970","0006590063715","0010465925138","0038948983000","0095774502751","0096051986448","0157967025871","0212722801067","0225815844582","0236197465609","0255973714120","0272643267465","0335243754848","0399569209871","0498395047364","0698616313838","0743731223606","0832349515187","0848056924761","0899183128263","1018947080336","1078766578568","1087115303928","1144906750559","1157194463322","1272575087123","1404411876176","1408483808240","1464781960745","1468999794635","1484994799123","1538211018396","1559730016366","1576499102253","1613730311804","1613963249998","1726123595351","1770178011663","1790740189682","1807963902339","1838229841499","1878073010375","1909194306167","1943523619306","1960588072419","2084705637233","2091630691049","2117622168280","2125923238175","2293343551454","2315609605258","2339358820724","2434890445616","2448500145612","2484897511500","2516506541834","2529776156387","2543248750439","2562717359713","2646001599860","2703547723491","2820850288674","2890379323668","2952766751666","3084140869281","3301617687934","3307124431763","3576267834421","3621077291879","3642408008706","3661250556340","3685697688891","3691313722887","3769015137275","3772022926992","3806751355669","3828680930458","3987140172453","3988126680641","4009053735033","4093428378987","4134502155718","4168050315812","4198118078267","4201407654834","4284914664558","4307721071729","4312472827051","4347934129886","4406572671024","4504719641739","4516685534489","4561421124790","4665903801947","4686859196154","4709231420798","4733532233696","4734618834332","4760375596107","4785470010730","4819782507258","4863291591550","4886504321812","4893655420066","4945934419923","4963935336179","4966277046676","5010710584900","5050273180195","5065846711133","5099742600901","5176352942567","5272733645116","5291392809646","5408760500061","5433448189252","5499736705597","5522456328132","5528517133622","5592486096660","5626486088179","5738533322232","5778452195678","5856636800041","5881647583898","5935916054838","5955704607626","6009279470754","6154584840805","6190070243323","6201199361567","6248889948356","6316992933962","6372029072158","6403813628678","6409491343148","6424174443243","6575325360237","6588791601926","6704641545275","6858015278648","6875096496570","6906120611895","6966709160725","7059492880556","7067375149041","7080123588503","7153630876392","7167715379463","7177157744491","7217466997444","7317365230007","7345418848909","7384311074268","7485639601133","7494558044822","7532074237028","7570673549500","7595223735110","7663515285824","7667946672475","7668932199532","7760442733661","7813908779724","7854842811538","8002754191821","8163753213485","8207931408888","8222420544052","8245402607613","8271165200181","8296484749050","8368305700967","8469939413398","8590367775021","8687358946192","8703661046340","8725228831589","8769809778856","8825217022124","8833419218504","8844419430964","8909358907493","8933669659420","9031323475252","9095019841233","9131148018211","9182640035008","9216642429807","9347854191845","9458076657016","9482467478850","9522454376759","9633135585459","9644009305424","9687547218818","9753065345920","9786855487647","9802920493181","9978391194435"]	\N	63	f	full	\N	2024-06-22 08:23:45.930738+00
22	2024-06-22 08:23:46.025745+00	2024-06-22 08:23:46.025745+00	["Aerodynamic Bronze Hat","Aerodynamic Concrete Bench","Aerodynamic Concrete Lamp","Aerodynamic Copper Knife","Aerodynamic Cotton Bottle","Aerodynamic Cotton Lamp","Aerodynamic Granite Bench","Aerodynamic Granite Bottle","Aerodynamic Leather Computer","Aerodynamic Leather Toucan","Aerodynamic Linen Coat","Aerodynamic Paper Coat","Aerodynamic Paper Computer","Aerodynamic Rubber Bench","Awesome Aluminum Keyboard","Awesome Aluminum Table","Awesome Bronze Plate","Awesome Concrete Shoes","Awesome Cotton Shoes","Awesome Granite Car","Awesome Iron Hat","Awesome Plastic Watch","Awesome Rubber Wallet","Awesome Silk Car","Awesome Steel Toucan","Awesome Wool Bench","Durable Aluminum Bag","Durable Copper Clock","Durable Cotton Bench","Durable Cotton Shirt","Durable Iron Knife","Durable Leather Wallet","Durable Marble Watch","Durable Rubber Computer","Durable Steel Toucan","Durable Wool Toucan","Enormous Aluminum Clock","Enormous Aluminum Shirt","Enormous Copper Shirt","Enormous Cotton Pants","Enormous Granite Bottle","Enormous Granite Wallet","Enormous Leather Wallet","Enormous Marble Gloves","Enormous Marble Shoes","Enormous Marble Wallet","Enormous Plastic Coat","Enormous Steel Watch","Enormous Wool Car","Ergonomic Aluminum Plate","Ergonomic Concrete Lamp","Ergonomic Cotton Bag","Ergonomic Granite Bottle","Ergonomic Iron Watch","Ergonomic Leather Pants","Ergonomic Linen Toucan","Ergonomic Marble Computer","Ergonomic Marble Hat","Ergonomic Paper Wallet","Ergonomic Plastic Bench","Ergonomic Rubber Bench","Ergonomic Silk Coat","Ergonomic Silk Keyboard","Ergonomic Silk Table","Ergonomic Wool Bag","Fantastic Aluminum Bottle","Fantastic Copper Hat","Fantastic Leather Watch","Fantastic Rubber Knife","Fantastic Silk Bottle","Fantastic Steel Knife","Fantastic Wool Shirt","Gorgeous Aluminum Plate","Gorgeous Bronze Hat","Gorgeous Concrete Chair","Gorgeous Concrete Shoes","Gorgeous Copper Knife","Gorgeous Linen Bottle","Gorgeous Linen Keyboard","Gorgeous Marble Computer","Gorgeous Marble Plate","Gorgeous Paper Bag","Gorgeous Wooden Car","Heavy-Duty Copper Gloves","Heavy-Duty Copper Toucan","Heavy-Duty Copper Watch","Heavy-Duty Cotton Bottle","Heavy-Duty Linen Gloves","Heavy-Duty Linen Toucan","Heavy-Duty Rubber Bottle","Heavy-Duty Rubber Gloves","Heavy-Duty Silk Car","Heavy-Duty Silk Chair","Heavy-Duty Steel Watch","Heavy-Duty Wooden Clock","Incredible Aluminum Knife","Incredible Bronze Pants","Incredible Bronze Wallet","Incredible Concrete Keyboard","Incredible Concrete Watch","Incredible Granite Toucan","Incredible Linen Knife","Incredible Plastic Chair","Incredible Plastic Watch","Incredible Silk Shoes","Intelligent Bronze Knife","Intelligent Granite Hat","Intelligent Iron Shirt","Intelligent Paper Car","Intelligent Paper Hat","Intelligent Steel Car","Intelligent Wooden Gloves","Lightweight Bronze Table","Lightweight Copper Wallet","Lightweight Granite Hat","Lightweight Leather Bench","Lightweight Leather Gloves","Lightweight Linen Bottle","Lightweight Linen Coat","Lightweight Linen Hat","Lightweight Marble Bag","Lightweight Paper Bottle","Lightweight Steel Knife","Lightweight Steel Watch","Lightweight Wool Bag","Lightweight Wool Computer","Lightweight Wool Plate","Mediocre Aluminum Lamp","Mediocre Aluminum Shirt","Mediocre Cotton Coat","Mediocre Cotton Toucan","Mediocre Leather Coat","Mediocre Leather Computer","Mediocre Marble Lamp","Mediocre Paper Car","Mediocre Plastic Clock","Mediocre Rubber Shoes","Mediocre Silk Bottle","Mediocre Wooden Bench","Mediocre Wooden Table","Mediocre Wool Toucan","Practical Aluminum Coat","Practical Aluminum Table","Practical Bronze Computer","Practical Bronze Watch","Practical Copper Car","Practical Granite Plate","Practical Paper Bag","Practical Plastic Keyboard","Practical Silk Bottle","Practical Silk Computer","Practical Steel Table","Practical Wool Hat","Rustic Concrete Bottle","Rustic Copper Hat","Rustic Copper Knife","Rustic Iron Bench","Rustic Iron Keyboard","Rustic Linen Keyboard","Rustic Marble Bottle","Rustic Paper Bench","Rustic Paper Car","Rustic Paper Wallet","Rustic Rubber Clock","Rustic Rubber Knife","Rustic Silk Knife","Rustic Silk Pants","Sleek Aluminum Clock","Sleek Aluminum Watch","Sleek Bronze Lamp","Sleek Copper Watch","Sleek Granite Pants","Sleek Leather Table","Sleek Leather Toucan","Sleek Marble Clock","Sleek Marble Table","Sleek Paper Toucan","Sleek Plastic Shoes","Sleek Steel Table","Sleek Wool Wallet","Sleek Wool Watch","Small Concrete Knife","Small Copper Clock","Small Copper Plate","Small Cotton Chair","Small Granite Gloves","Small Marble Hat","Small Marble Knife","Small Marble Shoes","Small Plastic Computer","Small Rubber Clock","Small Wool Wallet","Synergistic Copper Computer","Synergistic Granite Chair","Synergistic Leather Coat","Synergistic Marble Keyboard","Synergistic Rubber Shoes","Synergistic Steel Chair","Synergistic Wool Coat"]	\N	65	f	full	\N	2024-06-22 08:23:46.025745+00
23	2024-06-22 08:23:46.125755+00	2024-06-22 08:23:46.125755+00	["Alfreda Konopelski II Group","Alfredo Kuhlman Group","Americo Sipes and Sons","Annetta Wyman and Sons","Aufderhar-Boehm","Balistreri-Ankunding","Balistreri-Muller","Barrows-Johns","Batz-Schroeder","Baumbach-Hilpert","Bednar, Berge and Boyle","Berge, Mraz and Sawayn","Bernhard-Grady","Blake Greenfelder Group","Bosco-Breitenberg","Bradtke, Wilkinson and Reilly","Braeden Gislason and Sons","Brittany Mueller Inc","Cale Thompson V and Sons","Carmela Douglas Inc","Carol Marvin LLC","Casper-Schimmel","Cassin-Collins","Claude Thompson Group","Connelly-Mitchell","Connelly-Ritchie","Considine, Bogisich and Bauch","Considine, Lehner and Maggio","Considine, Schamberger and Schiller","Cremin-Williamson","Crona, Block and Homenick","Cruickshank-Abernathy","Daugherty-Dach","Delphia Bauch Inc","Demarcus Brakus Inc","Denesik-Ortiz","Devonte Gleichner Inc","Dominic Mann Group","Donnelly, Renner and Barton","Dooley-Cummings","Dora Fay and Sons","Dorothea Balistreri Inc","Emmerich-Nienow","Erika Volkman Group","Eugenia Kunze LLC","Fisher-Kemmer","Fisher-Purdy","Flatley-Kunde","Ford Runolfsson Group","Francis Wolff Group","Friesen-Anderson","Friesen-Langworth","Gail Bergstrom Inc","Gaylord-Lesch","Gibson, Turner and Douglas","Goyette-Smitham","Grady, Greenfelder and Welch","Gulgowski, Grimes and Mayer","Gutmann-Breitenberg","Hackett-Reynolds","Halle Kulas I LLC","Halvorson, Lockman and Ruecker","Hane, Hamill and Jerde","Hartmann, Mohr and Stiedemann","Hauck, Ernser and Barton","Heaney-Windler","Heathcote-Kirlin","Herman Flatley Group","Herman, Gleason and Renner","Hermiston, O'Hara and Wunsch","Herta Skiles and Sons","Hills, Fahey and Jones","Hilpert, Jacobs and Hauck","Hodkiewicz-Brekke","Howe, Kiehn and Price","Israel Spinka and Sons","Izabella Dach I and Sons","Jacobson-Daniel","Janick Harvey LLC","Jefferey Volkman LLC","Jerrell Gulgowski Inc","Jerrod McLaughlin LLC","Jones, Hayes and Kshlerin","Jordi Effertz LLC","Keely Stehr Group","Keshaun Mueller Group","Kiehn-Pacocha","Kiel Kassulke Group","Kirlin, Hermann and Stokes","Koch-Ruecker","Koepp, Ondricka and Larkin","Kuhlman-Kuphal","Kuhlman-McKenzie","Kuhn-O'Reilly","Kuphal, Brown and Koss","Kuphal, Friesen and Rowe","Kuphal, Schowalter and Bogan","Lakin-Stroman","Larson, Pfeffer and Klocko","Ledner-Satterfield","Ledner-Watsica","Legros, Lynch and Howell","Little-Pagac","Lon Wiegand DVM and Sons","Lorenza Mayer Inc","Maegan Casper Group","Marge Effertz Jr. Inc","Marquardt, Crooks and Abshire","Marvin, Turcotte and Wisozk","Mason Bashirian and Sons","Maxime Haley and Sons","Mayer, Kiehn and Turcotte","McClure-Lockman","McClure-Murphy","McDermott, Kiehn and Becker","McGlynn, Fay and Kertzmann","Medhurst-Reichert","Miles Ryan Group","Miss Annamae Kutch Group","Morar-Schamberger","Morissette, Bartoletti and Cummings","Morissette, Dare and Schimmel","Mr. Colton Mayer Group","Mr. Johanna Koepp and Sons","Mr. Tanya Stracke and Sons","Mrs. Eugenia Koelpin and Sons","Mueller, Mayert and Johnston","Mueller-Dare","Murray, Watsica and Wunsch","Myriam Macejkovic Inc","Nikolaus-Hudson","Noah Anderson and Sons","Nolan-Heller","Nolan-Wolff","Odessa Emmerich Inc","Okuneva, Kutch and Monahan","Ora Monahan and Sons","Oran D'Amore Inc","Orn, Hilpert and Pfannerstill","Pacocha-Volkman","Padberg, Senger and Williamson","Parker, O'Connell and Beahan","Pouros, Nitzsche and Mayer","Powlowski, Keebler and Quigley","Price Kuhic Inc","Price, Schultz and Daniel","Prohaska-Quigley","Quigley, Von and Will","Regan Bradtke and Sons","Reichert, Johnson and Roob","Reid Pfannerstill and Sons","Reynolds, Gleason and Brekke","Ritchie, Haley and Pacocha","Robyn Padberg Inc","Rodriguez-Kuhlman","Rosanna Murazik Inc","Roscoe Oberbrunner Group","Rowan Kautzer LLC","Ruecker, Carter and Ortiz","Ruecker-Jakubowski","Schamberger-Maggio","Schamberger-Wehner","Schiller, Bogisich and Lockman","Schinner, Schmitt and Crona","Schumm, Brown and Wehner","Schuster-Wyman","Senger, Mertz and Murray","Senger-Doyle","Senger-Stamm","Smitham, Dach and Bode","Spinka-Stokes","Stamm, Crist and Labadie","Stanton-Fritsch","Stark-Bayer","Stroman-Carroll","Swaniawski, Casper and Hilll","Theodora Terry and Sons","Theodore Hansen Inc","Thompson-Wolf","Tia Goyette Group","Toy, Deckow and Nitzsche","Trantow-Bartell","Turner, Kiehn and Schmitt","Una Fadel Group","Upton, Kovacek and Halvorson","Upton, Schoen and Streich","Ursula Collins LLC","Volkman, Greenfelder and Kiehn","Von-Gulgowski","Weimann-Cummings","West, Prohaska and Wunsch","Wilkinson, Donnelly and Gulgowski","Wilkinson-Gottlieb","Wisoky, Pagac and Heaney","Wiza, Abbott and Deckow","Wolf, Beahan and Thiel","Wolff, Ebert and Hansen","Wuckert, Murazik and Ernser","Zemlak, Botsford and Corkery","Zemlak-Wiegand"]	\N	60	f	full	\N	2024-06-22 08:23:46.125755+00
24	2024-06-22 08:23:46.225483+00	2024-06-22 08:23:46.225483+00	[1,2,3,4,5]	\N	66	f	full	\N	2024-06-22 08:23:46.225483+00
\.


--
-- Data for Name: metabase_table; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.metabase_table (id, created_at, updated_at, name, description, entity_type, active, db_id, display_name, visibility_type, schema, points_of_interest, caveats, show_in_getting_started, field_order, initial_sync_status, is_upload, database_require_filter, estimated_row_count, view_count) FROM stdin;
6	2024-06-22 08:19:17.839482+00	2024-06-22 08:23:41.325648+00	ACCOUNTS	Information on customer accounts registered with Piespace. Each account represents a new organization signing up for on-demand pies.	entity/UserTable	t	1	Accounts	\N	PUBLIC	Is it? We’ll let you be the judge of that.	Piespace’s business operates with a two week trial period. If you see that “Canceled At” is null then that account is still happily paying for their pies.	f	database	complete	f	\N	\N	0
1	2024-06-22 08:19:17.839482+00	2024-06-22 08:23:41.355868+00	ANALYTIC_EVENTS	Piespace does some anonymous analytics tracking on how users interact with their platform. They’ve only had time to implement a few events, but you know how it is. Pies come first.	entity/EventTable	t	1	Analytic Events	\N	PUBLIC	Is it? We’ll let you be the judge of that.	Piespace has cracked time travel, so keep in mind that some events may have already happened in the future.	f	database	complete	f	\N	\N	0
2	2024-06-22 08:19:17.839482+00	2024-06-22 08:23:41.413263+00	FEEDBACK	With each order of pies sent out, Piespace includes a place for customers to submit feedback and review their order.	entity/GenericTable	t	1	Feedback	\N	PUBLIC	Is it? We’ll let you be the judge of that.	Not every account feels inclined to submit feedback. That’s cool. There’s still quite a few responses here.	f	database	complete	f	\N	\N	0
8	2024-06-22 08:19:17.839482+00	2024-06-22 08:23:41.112968+00	PRODUCTS	Includes a catalog of all the products ever sold by the famed Sample Company.	entity/ProductTable	t	1	Products	\N	PUBLIC	Is it? You tell us!	The rating column is an integer from 1-5 where 1 is dreadful and 5 is the best thing ever.	f	database	complete	f	\N	\N	0
5	2024-06-22 08:19:17.839482+00	2024-06-22 08:23:41.125478+00	ORDERS	Confirmed Sample Company orders for a product, from a user.	entity/TransactionTable	t	1	Orders	\N	PUBLIC	Is it? You tell us!	You can join this on the Products and Orders table using the ID fields. Discount is left null if not applicable.	f	database	complete	f	\N	\N	0
3	2024-06-22 08:19:17.839482+00	2024-06-22 08:23:41.138348+00	PEOPLE	Information on the user accounts registered with Sample Company.	entity/UserTable	t	1	People	\N	PUBLIC	Is it? You tell us!	Note that employees and customer support staff will have accounts.	f	database	complete	f	\N	\N	0
4	2024-06-22 08:19:17.839482+00	2024-06-22 08:23:41.179011+00	REVIEWS	Reviews that Sample Company customers have left on our products.	entity/GenericTable	t	1	Reviews	\N	PUBLIC	Is it? You tell us!	These reviews aren't tied to orders so it is possible people have reviewed products they did not purchase from us.	f	database	complete	f	\N	\N	0
7	2024-06-22 08:19:17.839482+00	2024-06-22 08:23:41.31358+00	INVOICES	Confirmed payments from Piespace’s customers. Most accounts pay for their pie subscription on a monthly basis.	entity/GenericTable	t	1	Invoices	\N	PUBLIC	Is it? We’ll let you be the judge of that.	You can group by “Account ID” to see all the payments from an account and unveil information like total amount paid to date.	f	database	complete	f	\N	\N	0
\.


--
-- Data for Name: metric; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.metric (id, table_id, creator_id, name, description, archived, definition, created_at, updated_at, points_of_interest, caveats, how_is_this_calculated, show_in_getting_started, entity_id) FROM stdin;
\.


--
-- Data for Name: metric_important_field; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.metric_important_field (id, metric_id, field_id) FROM stdin;
\.


--
-- Data for Name: model_index; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.model_index (id, model_id, pk_ref, value_ref, schedule, state, indexed_at, error, created_at, creator_id) FROM stdin;
\.


--
-- Data for Name: model_index_value; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.model_index_value (model_index_id, model_pk, name) FROM stdin;
\.


--
-- Data for Name: moderation_review; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.moderation_review (id, updated_at, created_at, status, text, moderated_item_id, moderated_item_type, moderator_id, most_recent) FROM stdin;
\.


--
-- Data for Name: native_query_snippet; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.native_query_snippet (id, name, description, content, creator_id, archived, created_at, updated_at, collection_id, entity_id) FROM stdin;
\.


--
-- Data for Name: parameter_card; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.parameter_card (id, updated_at, created_at, card_id, parameterized_object_type, parameterized_object_id, parameter_id) FROM stdin;
\.


--
-- Data for Name: permissions; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.permissions (id, object, group_id) FROM stdin;
1	/	2
2	/collection/root/	1
3	/application/subscription/	1
4	/collection/namespace/snippets/root/	1
5	/collection/1/	1
\.


--
-- Data for Name: permissions_group; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.permissions_group (id, name, entity_id) FROM stdin;
1	All Users	\N
2	Administrators	\N
\.


--
-- Data for Name: permissions_group_membership; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.permissions_group_membership (id, user_id, group_id, is_group_manager) FROM stdin;
1	13371338	1	f
2	1	1	f
3	1	2	f
\.


--
-- Data for Name: permissions_revision; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.permissions_revision (id, before, after, user_id, created_at, remark) FROM stdin;
\.


--
-- Data for Name: persisted_info; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.persisted_info (id, database_id, card_id, question_slug, table_name, definition, query_hash, active, state, refresh_begin, refresh_end, state_change_at, error, created_at, creator_id) FROM stdin;
\.


--
-- Data for Name: pulse; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.pulse (id, creator_id, name, created_at, updated_at, skip_if_empty, alert_condition, alert_first_only, alert_above_goal, collection_id, collection_position, archived, dashboard_id, parameters, entity_id) FROM stdin;
\.


--
-- Data for Name: pulse_card; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.pulse_card (id, pulse_id, card_id, "position", include_csv, include_xls, dashboard_card_id, entity_id, format_rows) FROM stdin;
\.


--
-- Data for Name: pulse_channel; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.pulse_channel (id, pulse_id, channel_type, details, schedule_type, schedule_hour, schedule_day, created_at, updated_at, schedule_frame, enabled, entity_id) FROM stdin;
\.


--
-- Data for Name: pulse_channel_recipient; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.pulse_channel_recipient (id, pulse_channel_id, user_id) FROM stdin;
\.


--
-- Data for Name: qrtz_blob_triggers; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.qrtz_blob_triggers (sched_name, trigger_name, trigger_group, blob_data) FROM stdin;
\.


--
-- Data for Name: qrtz_calendars; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.qrtz_calendars (sched_name, calendar_name, calendar) FROM stdin;
\.


--
-- Data for Name: qrtz_cron_triggers; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.qrtz_cron_triggers (sched_name, trigger_name, trigger_group, cron_expression, time_zone_id) FROM stdin;
MetabaseScheduler	metabase.task.upgrade-checks.trigger	DEFAULT	0 15 6,18 * * ? *	GMT
MetabaseScheduler	metabase.task.anonymous-stats.trigger	DEFAULT	0 7 0 * * ? *	GMT
MetabaseScheduler	metabase.task.refresh-channel-cache.trigger	DEFAULT	0 37 0/4 1/1 * ? *	GMT
MetabaseScheduler	metabase.task.truncate-audit-tables.trigger	DEFAULT	0 0 */12 * * ? *	GMT
MetabaseScheduler	metabase.task.follow-up-emails.trigger	DEFAULT	0 0 12 * * ? *	GMT
MetabaseScheduler	metabase.task.creator-sentiment-emails.trigger	DEFAULT	0 0 2 ? * 7	GMT
MetabaseScheduler	metabase.task.task-history-cleanup.trigger	DEFAULT	0 0 0 * * ? *	GMT
MetabaseScheduler	metabase.task.sync-and-analyze.trigger.2	DEFAULT	0 48 * * * ? *	GMT
MetabaseScheduler	metabase.task.update-field-values.trigger.2	DEFAULT	0 0 2 * * ? *	GMT
\.


--
-- Data for Name: qrtz_fired_triggers; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.qrtz_fired_triggers (sched_name, entry_id, trigger_name, trigger_group, instance_name, fired_time, sched_time, priority, state, job_name, job_group, is_nonconcurrent, requests_recovery) FROM stdin;
\.


--
-- Data for Name: qrtz_job_details; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.qrtz_job_details (sched_name, job_name, job_group, description, job_class_name, is_durable, is_nonconcurrent, is_update_data, requests_recovery, job_data) FROM stdin;
MetabaseScheduler	metabase.task.sync-and-analyze.job	DEFAULT	sync-and-analyze for all databases	metabase.task.sync_databases.SyncAndAnalyzeDatabase	t	t	f	f	\\xaced0005737200156f72672e71756172747a2e4a6f62446174614d61709fb083e8bfa9b0cb020000787200266f72672e71756172747a2e7574696c732e537472696e674b65794469727479466c61674d61708208e8c3fbc55d280200015a0013616c6c6f77735472616e7369656e74446174617872001d6f72672e71756172747a2e7574696c732e4469727479466c61674d617013e62ead28760ace0200025a000564697274794c00036d617074000f4c6a6176612f7574696c2f4d61703b787000737200116a6176612e7574696c2e486173684d61700507dac1c31660d103000246000a6c6f6164466163746f724900097468726573686f6c6478703f40000000000010770800000010000000007800
MetabaseScheduler	metabase.task.update-field-values.job	DEFAULT	update-field-values for all databases	metabase.task.sync_databases.UpdateFieldValues	t	t	f	f	\\xaced0005737200156f72672e71756172747a2e4a6f62446174614d61709fb083e8bfa9b0cb020000787200266f72672e71756172747a2e7574696c732e537472696e674b65794469727479466c61674d61708208e8c3fbc55d280200015a0013616c6c6f77735472616e7369656e74446174617872001d6f72672e71756172747a2e7574696c732e4469727479466c61674d617013e62ead28760ace0200025a000564697274794c00036d617074000f4c6a6176612f7574696c2f4d61703b787000737200116a6176612e7574696c2e486173684d61700507dac1c31660d103000246000a6c6f6164466163746f724900097468726573686f6c6478703f40000000000010770800000010000000007800
MetabaseScheduler	metabase.task.PersistenceRefresh.job	DEFAULT	Persisted Model refresh task	metabase.task.persist_refresh.PersistenceRefresh	t	t	f	f	\\xaced0005737200156f72672e71756172747a2e4a6f62446174614d61709fb083e8bfa9b0cb020000787200266f72672e71756172747a2e7574696c732e537472696e674b65794469727479466c61674d61708208e8c3fbc55d280200015a0013616c6c6f77735472616e7369656e74446174617872001d6f72672e71756172747a2e7574696c732e4469727479466c61674d617013e62ead28760ace0200025a000564697274794c00036d617074000f4c6a6176612f7574696c2f4d61703b787000737200116a6176612e7574696c2e486173684d61700507dac1c31660d103000246000a6c6f6164466163746f724900097468726573686f6c6478703f40000000000010770800000010000000007800
MetabaseScheduler	metabase.task.upgrade-checks.job	DEFAULT	\N	metabase.task.upgrade_checks.CheckForNewVersions	f	f	f	f	\\xaced0005737200156f72672e71756172747a2e4a6f62446174614d61709fb083e8bfa9b0cb020000787200266f72672e71756172747a2e7574696c732e537472696e674b65794469727479466c61674d61708208e8c3fbc55d280200015a0013616c6c6f77735472616e7369656e74446174617872001d6f72672e71756172747a2e7574696c732e4469727479466c61674d617013e62ead28760ace0200025a000564697274794c00036d617074000f4c6a6176612f7574696c2f4d61703b787000737200116a6176612e7574696c2e486173684d61700507dac1c31660d103000246000a6c6f6164466163746f724900097468726573686f6c6478703f40000000000010770800000010000000007800
MetabaseScheduler	metabase.task.PersistencePrune.job	DEFAULT	Persisted Model prune task	metabase.task.persist_refresh.PersistencePrune	t	t	f	f	\\xaced0005737200156f72672e71756172747a2e4a6f62446174614d61709fb083e8bfa9b0cb020000787200266f72672e71756172747a2e7574696c732e537472696e674b65794469727479466c61674d61708208e8c3fbc55d280200015a0013616c6c6f77735472616e7369656e74446174617872001d6f72672e71756172747a2e7574696c732e4469727479466c61674d617013e62ead28760ace0200025a000564697274794c00036d617074000f4c6a6176612f7574696c2f4d61703b787000737200116a6176612e7574696c2e486173684d61700507dac1c31660d103000246000a6c6f6164466163746f724900097468726573686f6c6478703f40000000000010770800000010000000007800
MetabaseScheduler	metabase.task.anonymous-stats.job	DEFAULT	\N	metabase.task.send_anonymous_stats.SendAnonymousUsageStats	f	f	f	f	\\xaced0005737200156f72672e71756172747a2e4a6f62446174614d61709fb083e8bfa9b0cb020000787200266f72672e71756172747a2e7574696c732e537472696e674b65794469727479466c61674d61708208e8c3fbc55d280200015a0013616c6c6f77735472616e7369656e74446174617872001d6f72672e71756172747a2e7574696c732e4469727479466c61674d617013e62ead28760ace0200025a000564697274794c00036d617074000f4c6a6176612f7574696c2f4d61703b787000737200116a6176612e7574696c2e486173684d61700507dac1c31660d103000246000a6c6f6164466163746f724900097468726573686f6c6478703f40000000000010770800000010000000007800
MetabaseScheduler	metabase.task.IndexValues.job	DEFAULT	Indexed Value Refresh task	metabase.task.index_values.ModelIndexRefresh	t	t	f	f	\\xaced0005737200156f72672e71756172747a2e4a6f62446174614d61709fb083e8bfa9b0cb020000787200266f72672e71756172747a2e7574696c732e537472696e674b65794469727479466c61674d61708208e8c3fbc55d280200015a0013616c6c6f77735472616e7369656e74446174617872001d6f72672e71756172747a2e7574696c732e4469727479466c61674d617013e62ead28760ace0200025a000564697274794c00036d617074000f4c6a6176612f7574696c2f4d61703b787000737200116a6176612e7574696c2e486173684d61700507dac1c31660d103000246000a6c6f6164466163746f724900097468726573686f6c6478703f40000000000010770800000010000000007800
MetabaseScheduler	metabase.task.refresh-channel-cache.job	DEFAULT	\N	metabase.task.refresh_slack_channel_user_cache.RefreshCache	f	f	f	f	\\xaced0005737200156f72672e71756172747a2e4a6f62446174614d61709fb083e8bfa9b0cb020000787200266f72672e71756172747a2e7574696c732e537472696e674b65794469727479466c61674d61708208e8c3fbc55d280200015a0013616c6c6f77735472616e7369656e74446174617872001d6f72672e71756172747a2e7574696c732e4469727479466c61674d617013e62ead28760ace0200025a000564697274794c00036d617074000f4c6a6176612f7574696c2f4d61703b787000737200116a6176612e7574696c2e486173684d61700507dac1c31660d103000246000a6c6f6164466163746f724900097468726573686f6c6478703f40000000000010770800000010000000007800
MetabaseScheduler	metabase.task.truncate-audit-tables.job	DEFAULT	\N	metabase.task.truncate_audit_tables.TruncateAuditTables	f	f	f	f	\\xaced0005737200156f72672e71756172747a2e4a6f62446174614d61709fb083e8bfa9b0cb020000787200266f72672e71756172747a2e7574696c732e537472696e674b65794469727479466c61674d61708208e8c3fbc55d280200015a0013616c6c6f77735472616e7369656e74446174617872001d6f72672e71756172747a2e7574696c732e4469727479466c61674d617013e62ead28760ace0200025a000564697274794c00036d617074000f4c6a6176612f7574696c2f4d61703b787000737200116a6176612e7574696c2e486173684d61700507dac1c31660d103000246000a6c6f6164466163746f724900097468726573686f6c6478703f40000000000010770800000010000000007800
MetabaseScheduler	metabase.task.send-pulses.send-pulse.job	DEFAULT	Send Pulse	metabase.task.send_pulses.SendPulse	t	f	f	f	\\xaced0005737200156f72672e71756172747a2e4a6f62446174614d61709fb083e8bfa9b0cb020000787200266f72672e71756172747a2e7574696c732e537472696e674b65794469727479466c61674d61708208e8c3fbc55d280200015a0013616c6c6f77735472616e7369656e74446174617872001d6f72672e71756172747a2e7574696c732e4469727479466c61674d617013e62ead28760ace0200025a000564697274794c00036d617074000f4c6a6176612f7574696c2f4d61703b787000737200116a6176612e7574696c2e486173684d61700507dac1c31660d103000246000a6c6f6164466163746f724900097468726573686f6c6478703f40000000000010770800000010000000007800
MetabaseScheduler	metabase.task.send-pulses.init-send-pulse-triggers.job	DEFAULT	\N	metabase.task.send_pulses.InitSendPulseTriggers	t	f	f	f	\\xaced0005737200156f72672e71756172747a2e4a6f62446174614d61709fb083e8bfa9b0cb020000787200266f72672e71756172747a2e7574696c732e537472696e674b65794469727479466c61674d61708208e8c3fbc55d280200015a0013616c6c6f77735472616e7369656e74446174617872001d6f72672e71756172747a2e7574696c732e4469727479466c61674d617013e62ead28760ace0200025a000564697274794c00036d617074000f4c6a6176612f7574696c2f4d61703b787000737200116a6176612e7574696c2e486173684d61700507dac1c31660d103000246000a6c6f6164466163746f724900097468726573686f6c6478703f40000000000010770800000010000000007800
MetabaseScheduler	metabase.task.follow-up-emails.job	DEFAULT	\N	metabase.task.follow_up_emails.FollowUpEmail	f	f	f	f	\\xaced0005737200156f72672e71756172747a2e4a6f62446174614d61709fb083e8bfa9b0cb020000787200266f72672e71756172747a2e7574696c732e537472696e674b65794469727479466c61674d61708208e8c3fbc55d280200015a0013616c6c6f77735472616e7369656e74446174617872001d6f72672e71756172747a2e7574696c732e4469727479466c61674d617013e62ead28760ace0200025a000564697274794c00036d617074000f4c6a6176612f7574696c2f4d61703b787000737200116a6176612e7574696c2e486173684d61700507dac1c31660d103000246000a6c6f6164466163746f724900097468726573686f6c6478703f40000000000010770800000010000000007800
MetabaseScheduler	metabase.task.creator-sentiment-emails.job	DEFAULT	\N	metabase.task.creator_sentiment_emails.CreatorSentimentEmail	f	f	f	f	\\xaced0005737200156f72672e71756172747a2e4a6f62446174614d61709fb083e8bfa9b0cb020000787200266f72672e71756172747a2e7574696c732e537472696e674b65794469727479466c61674d61708208e8c3fbc55d280200015a0013616c6c6f77735472616e7369656e74446174617872001d6f72672e71756172747a2e7574696c732e4469727479466c61674d617013e62ead28760ace0200025a000564697274794c00036d617074000f4c6a6176612f7574696c2f4d61703b787000737200116a6176612e7574696c2e486173684d61700507dac1c31660d103000246000a6c6f6164466163746f724900097468726573686f6c6478703f40000000000010770800000010000000007800
MetabaseScheduler	metabase.task.legacy-no-self-service-emails.job	DEFAULT	\N	metabase.task.legacy_no_self_service_emails.LegacyNoSelfServiceEmail	t	f	f	f	\\xaced0005737200156f72672e71756172747a2e4a6f62446174614d61709fb083e8bfa9b0cb020000787200266f72672e71756172747a2e7574696c732e537472696e674b65794469727479466c61674d61708208e8c3fbc55d280200015a0013616c6c6f77735472616e7369656e74446174617872001d6f72672e71756172747a2e7574696c732e4469727479466c61674d617013e62ead28760ace0200025a000564697274794c00036d617074000f4c6a6176612f7574696c2f4d61703b787000737200116a6176612e7574696c2e486173684d61700507dac1c31660d103000246000a6c6f6164466163746f724900097468726573686f6c6478703f40000000000010770800000010000000007800
MetabaseScheduler	metabase.task.task-history-cleanup.job	DEFAULT	\N	metabase.task.task_history_cleanup.TaskHistoryCleanup	f	f	f	f	\\xaced0005737200156f72672e71756172747a2e4a6f62446174614d61709fb083e8bfa9b0cb020000787200266f72672e71756172747a2e7574696c732e537472696e674b65794469727479466c61674d61708208e8c3fbc55d280200015a0013616c6c6f77735472616e7369656e74446174617872001d6f72672e71756172747a2e7574696c732e4469727479466c61674d617013e62ead28760ace0200025a000564697274794c00036d617074000f4c6a6176612f7574696c2f4d61703b787000737200116a6176612e7574696c2e486173684d61700507dac1c31660d103000246000a6c6f6164466163746f724900097468726573686f6c6478703f40000000000010770800000010000000007800
MetabaseScheduler	metabase.task.email-remove-legacy-pulse.job	DEFAULT	\N	metabase.task.email_remove_legacy_pulse.EmailRemoveLegacyPulse	t	f	f	f	\\xaced0005737200156f72672e71756172747a2e4a6f62446174614d61709fb083e8bfa9b0cb020000787200266f72672e71756172747a2e7574696c732e537472696e674b65794469727479466c61674d61708208e8c3fbc55d280200015a0013616c6c6f77735472616e7369656e74446174617872001d6f72672e71756172747a2e7574696c732e4469727479466c61674d617013e62ead28760ace0200025a000564697274794c00036d617074000f4c6a6176612f7574696c2f4d61703b787000737200116a6176612e7574696c2e486173684d61700507dac1c31660d103000246000a6c6f6164466163746f724900097468726573686f6c6478703f40000000000010770800000010000000007800
\.


--
-- Data for Name: qrtz_locks; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.qrtz_locks (sched_name, lock_name) FROM stdin;
MetabaseScheduler	STATE_ACCESS
MetabaseScheduler	TRIGGER_ACCESS
\.


--
-- Data for Name: qrtz_paused_trigger_grps; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.qrtz_paused_trigger_grps (sched_name, trigger_group) FROM stdin;
\.


--
-- Data for Name: qrtz_scheduler_state; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.qrtz_scheduler_state (sched_name, instance_name, last_checkin_time, checkin_interval) FROM stdin;
MetabaseScheduler	metabase-d7fc79fd8-jl55m1719044626413	1719044754590	7500
\.


--
-- Data for Name: qrtz_simple_triggers; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.qrtz_simple_triggers (sched_name, trigger_name, trigger_group, repeat_count, repeat_interval, times_triggered) FROM stdin;
\.


--
-- Data for Name: qrtz_simprop_triggers; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.qrtz_simprop_triggers (sched_name, trigger_name, trigger_group, str_prop_1, str_prop_2, str_prop_3, int_prop_1, int_prop_2, long_prop_1, long_prop_2, dec_prop_1, dec_prop_2, bool_prop_1, bool_prop_2) FROM stdin;
\.


--
-- Data for Name: qrtz_triggers; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.qrtz_triggers (sched_name, trigger_name, trigger_group, job_name, job_group, description, next_fire_time, prev_fire_time, priority, trigger_state, trigger_type, start_time, end_time, calendar_name, misfire_instr, job_data) FROM stdin;
MetabaseScheduler	metabase.task.upgrade-checks.trigger	DEFAULT	metabase.task.upgrade-checks.job	DEFAULT	\N	1719080100000	-1	5	WAITING	CRON	1719044626000	0	\N	0	\\x
MetabaseScheduler	metabase.task.anonymous-stats.trigger	DEFAULT	metabase.task.anonymous-stats.job	DEFAULT	\N	1719101220000	-1	5	WAITING	CRON	1719044626000	0	\N	0	\\x
MetabaseScheduler	metabase.task.refresh-channel-cache.trigger	DEFAULT	metabase.task.refresh-channel-cache.job	DEFAULT	\N	1719045420000	-1	5	WAITING	CRON	1719044626000	0	\N	2	\\x
MetabaseScheduler	metabase.task.truncate-audit-tables.trigger	DEFAULT	metabase.task.truncate-audit-tables.job	DEFAULT	\N	1719057600000	-1	5	WAITING	CRON	1719044626000	0	\N	2	\\x
MetabaseScheduler	metabase.task.follow-up-emails.trigger	DEFAULT	metabase.task.follow-up-emails.job	DEFAULT	\N	1719057600000	-1	5	WAITING	CRON	1719044626000	0	\N	0	\\x
MetabaseScheduler	metabase.task.creator-sentiment-emails.trigger	DEFAULT	metabase.task.creator-sentiment-emails.job	DEFAULT	\N	1719626400000	-1	5	WAITING	CRON	1719044627000	0	\N	0	\\x
MetabaseScheduler	metabase.task.task-history-cleanup.trigger	DEFAULT	metabase.task.task-history-cleanup.job	DEFAULT	\N	1719100800000	-1	5	WAITING	CRON	1719044627000	0	\N	0	\\x
MetabaseScheduler	metabase.task.sync-and-analyze.trigger.2	DEFAULT	metabase.task.sync-and-analyze.job	DEFAULT	sync-and-analyze Database 2	1719046080000	-1	5	WAITING	CRON	1719044736000	0	\N	2	\\xaced0005737200156f72672e71756172747a2e4a6f62446174614d61709fb083e8bfa9b0cb020000787200266f72672e71756172747a2e7574696c732e537472696e674b65794469727479466c61674d61708208e8c3fbc55d280200015a0013616c6c6f77735472616e7369656e74446174617872001d6f72672e71756172747a2e7574696c732e4469727479466c61674d617013e62ead28760ace0200025a000564697274794c00036d617074000f4c6a6176612f7574696c2f4d61703b787000737200116a6176612e7574696c2e486173684d61700507dac1c31660d103000246000a6c6f6164466163746f724900097468726573686f6c6478703f4000000000000c7708000000100000000174000564622d6964737200116a6176612e6c616e672e496e746567657212e2a0a4f781873802000149000576616c7565787200106a6176612e6c616e672e4e756d62657286ac951d0b94e08b0200007870000000027800
MetabaseScheduler	metabase.task.update-field-values.trigger.2	DEFAULT	metabase.task.update-field-values.job	DEFAULT	update-field-values Database 2	1719108000000	-1	5	WAITING	CRON	1719044736000	0	\N	2	\\xaced0005737200156f72672e71756172747a2e4a6f62446174614d61709fb083e8bfa9b0cb020000787200266f72672e71756172747a2e7574696c732e537472696e674b65794469727479466c61674d61708208e8c3fbc55d280200015a0013616c6c6f77735472616e7369656e74446174617872001d6f72672e71756172747a2e7574696c732e4469727479466c61674d617013e62ead28760ace0200025a000564697274794c00036d617074000f4c6a6176612f7574696c2f4d61703b787000737200116a6176612e7574696c2e486173684d61700507dac1c31660d103000246000a6c6f6164466163746f724900097468726573686f6c6478703f4000000000000c7708000000100000000174000564622d6964737200116a6176612e6c616e672e496e746567657212e2a0a4f781873802000149000576616c7565787200106a6176612e6c616e672e4e756d62657286ac951d0b94e08b0200007870000000027800
\.


--
-- Data for Name: query; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.query (query_hash, average_execution_time, query) FROM stdin;
\.


--
-- Data for Name: query_action; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.query_action (action_id, database_id, dataset_query) FROM stdin;
\.


--
-- Data for Name: query_cache; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.query_cache (query_hash, updated_at, results) FROM stdin;
\.


--
-- Data for Name: query_execution; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.query_execution (id, hash, started_at, running_time, result_rows, native, context, error, executor_id, card_id, dashboard_id, pulse_id, database_id, cache_hit, action_id, is_sandboxed, cache_hash) FROM stdin;
\.


--
-- Data for Name: query_field; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.query_field (id, card_id, field_id, direct_reference) FROM stdin;
\.


--
-- Data for Name: recent_views; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.recent_views (id, user_id, model, model_id, "timestamp") FROM stdin;
\.


--
-- Data for Name: report_card; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.report_card (id, created_at, updated_at, name, description, display, dataset_query, visualization_settings, creator_id, database_id, table_id, query_type, archived, collection_id, public_uuid, made_public_by_id, enable_embedding, embedding_params, cache_ttl, result_metadata, collection_position, entity_id, parameters, parameter_mappings, collection_preview, metabase_version, type, initially_published_at, cache_invalidated_at, last_used_at, view_count) FROM stdin;
1	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	Aerodynamic Copper Knife trend	Compares the total number of orders placed for this product this month with the previous period	smartscalar	{"database":1,"type":"query","query":{"aggregation":[["sum",["field",39,{"base-type":"type/Integer"}]]],"breakout":[["field",41,{"base-type":"type/DateTime","temporal-unit":"month"}]],"source-table":5,"filter":["and",["between",["field",41,{"base-type":"type/DateTime"}],"2021-01-01","2024-01-23"],["=",["field",62,{"base-type":"type/BigInteger","source-field":40}],165]]}}	{"column_settings":null}	13371338	1	5	query	f	1	\N	\N	f	\N	\N	[{"description":"The date and time an order was submitted.","semantic_type":"type/CreationTimestamp","coercion_strategy":null,"unit":"month","name":"CREATED_AT","settings":null,"fk_target_field_id":null,"field_ref":["field",41,{"base-type":"type/DateTime","temporal-unit":"month"}],"effective_type":"type/DateTime","id":41,"visibility_type":"normal","display_name":"Created At","fingerprint":{"global":{"distinct-count":10001,"nil%":0.0},"type":{"type/DateTime":{"earliest":"2022-04-30T18:56:13.352Z","latest":"2026-04-19T14:07:15.657Z"}}},"base_type":"type/DateTime"},{"display_name":"Sum of Quantity","semantic_type":"type/Quantity","settings":null,"field_ref":["aggregation",0],"name":"sum","base_type":"type/BigInteger","effective_type":"type/BigInteger","fingerprint":{"global":{"distinct-count":7,"nil%":0.0},"type":{"type/Number":{"min":1.0,"q1":3.25,"q3":6.0,"max":74.0,"sd":22.086949389477336,"avg":11.5}}}}]	\N	oShRjudlgUvJRezvGTQTD	[]	[]	t	\N	question	\N	\N	2024-06-22 08:19:17.839482+00	0
2	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	Orders with People	\N	table	{"database":1,"type":"query","query":{"expressions":{"Age":["datetime-diff",["field",49,{"base-type":"type/Date","join-alias":"People - User"}],["now"],"year"]},"joins":[{"alias":"People - User","strategy":"left-join","condition":["=",["field",43,{"base-type":"type/Integer"}],["field",46,{"base-type":"type/BigInteger","join-alias":"People - User"}]],"source-table":3}],"source-table":5}}	{"table.column_widths":[null,null,null,null,null,null,null,null,null,164],"graph.show_values":true,"table.cell_column":"SUBTOTAL","graph.series_order_dimension":null,"graph.metrics":[],"graph.series_order":null,"table.pivot_column":"Age","graph.dimensions":["Age"],"stackable.stack_type":"stacked"}	13371338	1	5	query	f	1	\N	\N	f	\N	\N	[{"description":"This is a unique ID for the product. It is also called the “Invoice number” or “Confirmation number” in customer facing emails and screens.","semantic_type":"type/PK","coercion_strategy":null,"name":"ID","settings":null,"fk_target_field_id":null,"field_ref":["field",37,null],"effective_type":"type/BigInteger","id":37,"visibility_type":"normal","display_name":"ID","fingerprint":null,"base_type":"type/BigInteger"},{"description":"The id of the user who made this order. Note that in some cases where an order was created on behalf of a customer who phoned the order in, this might be the employee who handled the request.","semantic_type":"type/FK","coercion_strategy":null,"name":"USER_ID","settings":null,"fk_target_field_id":46,"field_ref":["field",43,null],"effective_type":"type/Integer","id":43,"visibility_type":"normal","display_name":"User ID","fingerprint":{"global":{"distinct-count":929,"nil%":0}},"base_type":"type/Integer"},{"description":"The product ID. This is an internal identifier for the product, NOT the SKU.","semantic_type":"type/FK","coercion_strategy":null,"name":"PRODUCT_ID","settings":null,"fk_target_field_id":62,"field_ref":["field",40,null],"effective_type":"type/Integer","id":40,"visibility_type":"normal","display_name":"Product ID","fingerprint":{"global":{"distinct-count":200,"nil%":0}},"base_type":"type/Integer"},{"description":"The raw, pre-tax cost of the order. Note that this might be different in the future from the product price due to promotions, credits, etc.","semantic_type":null,"coercion_strategy":null,"name":"SUBTOTAL","settings":null,"fk_target_field_id":null,"field_ref":["field",44,null],"effective_type":"type/Float","id":44,"visibility_type":"normal","display_name":"Subtotal","fingerprint":{"global":{"distinct-count":340,"nil%":0},"type":{"type/Number":{"min":15.691943673970439,"q1":49.74894519060184,"q3":105.42965746993103,"max":148.22900526552291,"sd":32.53705013056317,"avg":77.01295465356547}}},"base_type":"type/Float"},{"description":"This is the amount of local and federal taxes that are collected on the purchase. Note that other governmental fees on some products are not included here, but instead are accounted for in the subtotal.","semantic_type":null,"coercion_strategy":null,"name":"TAX","settings":null,"fk_target_field_id":null,"field_ref":["field",38,null],"effective_type":"type/Float","id":38,"visibility_type":"normal","display_name":"Tax","fingerprint":{"global":{"distinct-count":797,"nil%":0},"type":{"type/Number":{"min":0,"q1":2.273340386603857,"q3":5.337275338216307,"max":11.12,"sd":2.3206651358900316,"avg":3.8722100000000004}}},"base_type":"type/Float"},{"description":"The total billed amount.","semantic_type":null,"coercion_strategy":null,"name":"TOTAL","settings":null,"fk_target_field_id":null,"field_ref":["field",42,null],"effective_type":"type/Float","id":42,"visibility_type":"normal","display_name":"Total","fingerprint":{"global":{"distinct-count":4426,"nil%":0},"type":{"type/Number":{"min":8.93914247937167,"q1":51.34535490743823,"q3":110.29428389265787,"max":159.34900526552292,"sd":34.26469575709948,"avg":80.35871658771228}}},"base_type":"type/Float"},{"description":"Discount amount.","semantic_type":"type/Discount","coercion_strategy":null,"name":"DISCOUNT","settings":null,"fk_target_field_id":null,"field_ref":["field",36,null],"effective_type":"type/Float","id":36,"visibility_type":"normal","display_name":"Discount","fingerprint":{"global":{"distinct-count":701,"nil%":0.898},"type":{"type/Number":{"min":0.17088996672584322,"q1":2.9786226681458743,"q3":7.338187788658235,"max":61.69684269960571,"sd":3.053663125001991,"avg":5.161255547580326}}},"base_type":"type/Float"},{"description":"The date and time an order was submitted.","semantic_type":"type/CreationTimestamp","coercion_strategy":null,"unit":"default","name":"CREATED_AT","settings":null,"fk_target_field_id":null,"field_ref":["field",41,{"temporal-unit":"default"}],"effective_type":"type/DateTime","id":41,"visibility_type":"normal","display_name":"Created At","fingerprint":{"global":{"distinct-count":10001,"nil%":0},"type":{"type/DateTime":{"earliest":"2022-04-30T18:56:13.352Z","latest":"2026-04-19T14:07:15.657Z"}}},"base_type":"type/DateTime"},{"description":"Number of products bought.","semantic_type":"type/Quantity","coercion_strategy":null,"name":"QUANTITY","settings":null,"fk_target_field_id":null,"field_ref":["field",39,null],"effective_type":"type/Integer","id":39,"visibility_type":"normal","display_name":"Quantity","fingerprint":{"global":{"distinct-count":62,"nil%":0},"type":{"type/Number":{"min":0,"q1":1.755882607764982,"q3":4.882654507928044,"max":100,"sd":4.214258386403798,"avg":3.7015}}},"base_type":"type/Integer"},{"display_name":"Age","field_ref":["field","Age",{"base-type":"type/BigInteger"}],"name":"Age","base_type":"type/BigInteger","effective_type":"type/BigInteger","fingerprint":{"global":{"distinct-count":42,"nil%":0},"type":{"type/Number":{"min":24,"q1":33.36752836803635,"q3":55.20362176071121,"max":65,"sd":12.063315373018085,"avg":44.572}}},"semantic_type":null}]	1	hjNoawcRfsDrC32g7LSOE	[]	[]	t	\N	model	\N	\N	2024-06-22 08:19:17.839482+00	0
3	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	Revenue per age group	Shows the revenue distributed by age group	bar	{"database":1,"type":"query","query":{"aggregation":[["sum",["field","TOTAL",{"base-type":"type/Float"}]]],"breakout":[["field","Age",{"base-type":"type/BigInteger","binning":{"strategy":"num-bins","num-bins":10}}]],"source-table":"card__2"}}	{"column_settings":null,"graph.dimensions":["Age"],"graph.metrics":["sum"],"graph.series_order":null,"graph.series_order_dimension":null,"graph.show_values":true,"stackable.stack_type":"stacked","table.cell_column":"SUBTOTAL"}	13371338	1	5	query	f	1	\N	\N	f	\N	\N	[{"field_ref":["field","Age",{"base-type":"type/BigInteger","binning":{"strategy":"num-bins","num-bins":10,"min-value":20.0,"max-value":65.0,"bin-width":5.0}}],"base_type":"type/BigInteger","name":"Age","effective_type":"type/BigInteger","display_name":"Age","fingerprint":{"global":{"distinct-count":42,"nil%":0},"type":{"type/Number":{"min":24,"q1":33.36752836803635,"q3":55.20362176071121,"max":65,"sd":12.063315373018085,"avg":44.572}}},"binning_info":{"num_bins":10,"min_value":20.0,"max_value":65.0,"bin_width":5.0,"binning_strategy":"num-bins"},"source":"breakout"},{"base_type":"type/Float","name":"sum","display_name":"Sum of Total","source":"aggregation","field_ref":["aggregation",0],"aggregation_index":0}]	\N	rqxxsDv8zhSBXgJmVjyWa	[]	[]	t	\N	question	\N	\N	2024-06-22 08:19:17.839482+00	0
4	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	Buyers by age group	Shows a distribution of our customers in age groups	bar	{"database":1,"type":"query","query":{"aggregation":[["count"]],"breakout":[["field","Age",{"base-type":"type/BigInteger","binning":{"strategy":"num-bins","num-bins":10}}]],"source-table":"card__26"}}	{"graph.dimensions":["Age"],"series_settings":{"count":{"color":"#A989C5"}},"graph.series_order_dimension":null,"graph.series_order":null,"graph.metrics":["count"],"graph.y_axis.title_text":"Buyers"}	13371338	1	3	query	f	1	\N	\N	f	\N	\N	[{"display_name":"Age","field_ref":["field","Age",{"base-type":"type/BigInteger","binning":{"strategy":"num-bins","num-bins":10,"min-value":20,"max-value":65,"bin-width":5}}],"name":"Age","base_type":"type/Decimal","effective_type":"type/Decimal","fingerprint":{"global":{"distinct-count":42,"nil%":0},"type":{"type/Number":{"min":24,"q1":33.340572873934306,"q3":55.17166516756599,"max":65,"sd":12.263883782175668,"avg":44.434}}},"semantic_type":null},{"display_name":"Count","semantic_type":"type/Quantity","field_ref":["aggregation",0],"name":"count","base_type":"type/BigInteger","effective_type":"type/BigInteger","fingerprint":{"global":{"distinct-count":9,"nil%":0},"type":{"type/Number":{"min":61,"q1":270,"q3":304,"max":334,"sd":99.91663191547909,"avg":250}}}}]	\N	qj0jT7SXwEUezz1wSjtaZ	[]	[]	t	\N	question	\N	\N	2024-06-22 08:19:17.839482+00	0
5	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	Heavy-Duty Silk Chair trend	Compares the total number of orders placed for this product this month with the previous period	smartscalar	{"database":1,"type":"query","query":{"aggregation":[["sum",["field",39,{"base-type":"type/Integer"}]]],"breakout":[["field",41,{"base-type":"type/DateTime","temporal-unit":"month"}]],"source-table":5,"filter":["and",["between",["field",41,{"base-type":"type/DateTime"}],"2021-01-01","2024-01-23"],["=",["field",62,{"base-type":"type/BigInteger","source-field":40}],86]]}}	{"column_settings":null}	13371338	1	5	query	f	1	\N	\N	f	\N	\N	[{"description":"The date and time an order was submitted.","semantic_type":"type/CreationTimestamp","coercion_strategy":null,"unit":"month","name":"CREATED_AT","settings":null,"fk_target_field_id":null,"field_ref":["field",41,{"base-type":"type/DateTime","temporal-unit":"month"}],"effective_type":"type/DateTime","id":41,"visibility_type":"normal","display_name":"Created At","fingerprint":{"global":{"distinct-count":10001,"nil%":0.0},"type":{"type/DateTime":{"earliest":"2022-04-30T18:56:13.352Z","latest":"2026-04-19T14:07:15.657Z"}}},"base_type":"type/DateTime"},{"display_name":"Sum of Quantity","semantic_type":"type/Quantity","settings":null,"field_ref":["aggregation",0],"name":"sum","base_type":"type/BigInteger","effective_type":"type/BigInteger","fingerprint":{"global":{"distinct-count":7,"nil%":0.0},"type":{"type/Number":{"min":1.0,"q1":2.5,"q3":9.414213562373096,"max":10.0,"sd":3.693623849670827,"avg":5.75}}}}]	\N	nhFZBIIZA_KnK6RuLp36z	[]	[]	t	\N	question	\N	\N	2024-06-22 08:19:17.839482+00	0
6	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	Total orders by category	Breaks down the overall performance of each of the product categories	pie	{"database":1,"type":"query","query":{"aggregation":[["count"]],"breakout":[["field",58,{"base-type":"type/Text","source-field":40}]],"source-table":5}}	{"column_settings":null,"pie.colors":{"Doohickey":"#7172AD","Gadget":"#A989C5","Gizmo":"#C7EAEA","Widget":"#227FD2"}}	13371338	1	5	query	f	1	\N	\N	f	\N	\N	[{"description":"The type of product, valid values include: Doohicky, Gadget, Gizmo and Widget","semantic_type":"type/Category","coercion_strategy":null,"name":"CATEGORY","settings":null,"fk_target_field_id":null,"field_ref":["field",58,{"base-type":"type/Text","source-field":40}],"effective_type":"type/Text","id":58,"visibility_type":"normal","display_name":"Product → Category","fingerprint":{"global":{"distinct-count":4,"nil%":0.0},"type":{"type/Text":{"percent-json":0.0,"percent-url":0.0,"percent-email":0.0,"percent-state":0.0,"average-length":6.375}}},"base_type":"type/Text"},{"display_name":"Count","semantic_type":"type/Quantity","field_ref":["aggregation",0],"name":"count","base_type":"type/BigInteger","effective_type":"type/BigInteger","fingerprint":{"global":{"distinct-count":4,"nil%":0.0},"type":{"type/Number":{"min":3976.0,"q1":4380.0,"q3":5000.0,"max":5061.0,"sd":489.3103990992493,"avg":4690.0}}}}]	\N	1aVGgD6-onGIPUeWCLiMA	[]	[]	t	\N	question	\N	\N	2024-06-22 08:19:17.839482+00	0
7	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	Value of orders over time (before taxes)	\N	bar	{"database":1,"type":"query","query":{"aggregation":[["sum",["field",44,{"base-type":"type/Float"}]]],"breakout":[["field",41,{"base-type":"type/DateTime","temporal-unit":"month"}],["field",58,{"base-type":"type/Text","source-field":40}]],"source-table":5,"filter":["time-interval",["field",41,{"base-type":"type/DateTime"}],-24,"month"]}}	{"column_settings":null,"graph.dimensions":["CREATED_AT","CATEGORY"],"graph.metrics":["sum"],"graph.series_order":null,"graph.series_order_dimension":null,"stackable.stack_type":"stacked"}	13371338	1	5	query	f	1	\N	\N	f	\N	\N	[{"description":"The date and time an order was submitted.","semantic_type":"type/CreationTimestamp","coercion_strategy":null,"unit":"month","name":"CREATED_AT","settings":null,"fk_target_field_id":null,"field_ref":["field",41,{"base-type":"type/DateTime","temporal-unit":"month"}],"effective_type":"type/DateTime","id":41,"visibility_type":"normal","display_name":"Created At","fingerprint":{"global":{"distinct-count":10001,"nil%":0.0},"type":{"type/DateTime":{"earliest":"2022-04-30T18:56:13.352Z","latest":"2026-04-19T14:07:15.657Z"}}},"base_type":"type/DateTime"},{"description":"The type of product, valid values include: Doohicky, Gadget, Gizmo and Widget","semantic_type":"type/Category","coercion_strategy":null,"name":"CATEGORY","settings":null,"fk_target_field_id":null,"field_ref":["field",58,{"base-type":"type/Text","source-field":40}],"effective_type":"type/Text","id":58,"visibility_type":"normal","display_name":"Product → Category","fingerprint":{"global":{"distinct-count":4,"nil%":0.0},"type":{"type/Text":{"percent-json":0.0,"percent-url":0.0,"percent-email":0.0,"percent-state":0.0,"average-length":6.375}}},"base_type":"type/Text"},{"display_name":"Sum of Subtotal","semantic_type":null,"settings":null,"field_ref":["aggregation",0],"name":"sum","base_type":"type/Float","effective_type":"type/Float","fingerprint":{"global":{"distinct-count":96,"nil%":0.0},"type":{"type/Number":{"min":194.15721878715678,"q1":1682.037094502916,"q3":5250.149916853326,"max":10113.40941589646,"sd":2757.8851826916916,"avg":3956.9519658839577}}}}]	\N	f9RqS8HspKvcsGEH24_Af	[]	[]	t	\N	question	\N	\N	2024-06-22 08:19:17.839482+00	0
8	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	Enormous Wool Car trend	Compares the total number of orders placed for this product this month with the previous period	smartscalar	{"database":1,"type":"query","query":{"aggregation":[["sum",["field",39,{"base-type":"type/Integer"}]]],"breakout":[["field",41,{"base-type":"type/DateTime","temporal-unit":"month"}]],"source-table":5,"filter":["and",["between",["field",41,{"base-type":"type/DateTime"}],"2021-01-01","2024-01-23"],["=",["field",62,{"base-type":"type/BigInteger","source-field":40}],17]]}}	{"column_settings":null}	13371338	1	5	query	f	1	\N	\N	f	\N	\N	[{"description":"The date and time an order was submitted.","semantic_type":"type/CreationTimestamp","coercion_strategy":null,"unit":"month","name":"CREATED_AT","settings":null,"fk_target_field_id":null,"field_ref":["field",41,{"base-type":"type/DateTime","temporal-unit":"month"}],"effective_type":"type/DateTime","id":41,"visibility_type":"normal","display_name":"Created At","fingerprint":{"global":{"distinct-count":10001,"nil%":0.0},"type":{"type/DateTime":{"earliest":"2022-04-30T18:56:13.352Z","latest":"2026-04-19T14:07:15.657Z"}}},"base_type":"type/DateTime"},{"display_name":"Sum of Quantity","semantic_type":"type/Quantity","settings":null,"field_ref":["aggregation",0],"name":"sum","base_type":"type/BigInteger","effective_type":"type/BigInteger","fingerprint":{"global":{"distinct-count":11,"nil%":0.0},"type":{"type/Number":{"min":3.0,"q1":4.76393202250021,"q3":17.0,"max":21.0,"sd":6.199964551476116,"avg":10.142857142857142}}}}]	\N	BP4o72CSvPOTLe5Nk6twf	[]	[]	t	\N	question	\N	\N	2024-06-22 08:19:17.839482+00	0
9	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	Revenue by state	\N	map	{"database":1,"type":"query","query":{"aggregation":[["sum",["field",42,{"base-type":"type/Float"}]]],"breakout":[["field",48,{"base-type":"type/Text","source-field":43}]],"source-table":5}}	{"column_settings":null,"map.colors":["#e5e5f1","#c0c0da","#9b9cc3","#7677ab","hsl(239, 29.5%, 39.3%)"]}	13371338	1	5	query	f	1	\N	\N	f	\N	\N	[{"description":"The state or province of the account’s billing address","semantic_type":"type/State","coercion_strategy":null,"name":"STATE","settings":null,"fk_target_field_id":null,"field_ref":["field",48,{"base-type":"type/Text","source-field":43}],"effective_type":"type/Text","id":48,"visibility_type":"normal","display_name":"User → State","fingerprint":{"global":{"distinct-count":49,"nil%":0.0},"type":{"type/Text":{"percent-json":0.0,"percent-url":0.0,"percent-email":0.0,"percent-state":1.0,"average-length":2.0}}},"base_type":"type/Text"},{"display_name":"Sum of Total","semantic_type":null,"settings":null,"field_ref":["aggregation",0],"name":"sum","base_type":"type/Float","effective_type":"type/Float","fingerprint":{"global":{"distinct-count":48,"nil%":0.0},"type":{"type/Number":{"min":1358.6880585850759,"q1":14606.102496036605,"q3":45423.71594332504,"max":108466.5974651383,"sd":21267.69039176716,"avg":31471.285063554824}}}}]	\N	hBnes5i3LGSYZeoGIYwvm	[]	[]	t	\N	question	\N	\N	2024-06-22 08:19:17.839482+00	0
10	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	Orders by source per age group	Shows a distribution of orders broken down by source across our customers' age groups	bar	{"database":1,"type":"query","query":{"aggregation":[["count"]],"breakout":[["field","Age",{"base-type":"type/BigInteger","binning":{"strategy":"num-bins","num-bins":10}}],["field",45,{"base-type":"type/Text","join-alias":"People - User"}]],"joins":[{"alias":"People - User","fields":"all","strategy":"left-join","condition":["=",["field","USER_ID",{"base-type":"type/Integer"}],["field",46,{"base-type":"type/BigInteger","join-alias":"People - User"}]],"source-table":3}],"source-table":"card__2"}}	{"column_settings":null,"graph.dimensions":["Age","SOURCE"],"graph.metrics":["count"],"graph.series_order":null,"graph.series_order_dimension":null,"graph.show_values":true,"graph.x_axis.axis_enabled":true,"graph.x_axis.labels_enabled":true,"graph.x_axis.scale":"ordinal","graph.y_axis.auto_range":true,"stackable.stack_type":"stacked","table.cell_column":"SUBTOTAL"}	13371338	1	5	query	f	1	\N	\N	f	\N	\N	[{"field_ref":["field","Age",{"base-type":"type/BigInteger","binning":{"strategy":"num-bins","num-bins":10,"min-value":20.0,"max-value":65.0,"bin-width":5.0}}],"base_type":"type/BigInteger","name":"Age","effective_type":"type/BigInteger","display_name":"Age","fingerprint":{"global":{"distinct-count":42,"nil%":0},"type":{"type/Number":{"min":24,"q1":33.36752836803635,"q3":55.20362176071121,"max":65,"sd":12.063315373018085,"avg":44.572}}},"binning_info":{"num_bins":10,"min_value":20.0,"max_value":65.0,"bin_width":5.0,"binning_strategy":"num-bins"},"source":"breakout"},{"description":"The channel through which we acquired this user. Valid values include: Affiliate, Facebook, Google, Organic and Twitter","semantic_type":"type/Source","table_id":3,"coercion_strategy":null,"name":"SOURCE","settings":null,"source":"breakout","fk_target_field_id":null,"field_ref":["field",45,{"base-type":"type/Text","join-alias":"People - User"}],"effective_type":"type/Text","nfc_path":null,"parent_id":null,"id":45,"position":8,"visibility_type":"normal","display_name":"People - User → Source","fingerprint":{"global":{"distinct-count":5,"nil%":0.0},"type":{"type/Text":{"percent-json":0.0,"percent-url":0.0,"percent-email":0.0,"percent-state":0.0,"average-length":7.4084}}},"base_type":"type/Text","source_alias":"People - User"},{"base_type":"type/Integer","name":"count","display_name":"Count","semantic_type":"type/Quantity","source":"aggregation","field_ref":["aggregation",0],"aggregation_index":0}]	\N	t1rQijQdCIZoQI6bhwBEG	[]	[]	t	\N	question	\N	\N	2024-06-22 08:19:17.839482+00	0
11	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	Product category orders per age group	Shows a distribution of orders broken down by product categories across our customers' age groups	bar	{"database":1,"type":"query","query":{"aggregation":[["count"]],"breakout":[["field","Age",{"base-type":"type/BigInteger","binning":{"strategy":"num-bins","num-bins":10}}],["field",58,{"base-type":"type/Text","join-alias":"Products"}]],"joins":[{"alias":"Products","fields":"all","strategy":"left-join","condition":["=",["field","PRODUCT_ID",{"base-type":"type/Integer"}],["field",62,{"base-type":"type/BigInteger","join-alias":"Products"}]],"source-table":8}],"source-table":"card__2"}}	{"column_settings":null,"graph.dimensions":["Age","CATEGORY"],"graph.metrics":["count"],"graph.series_order":null,"graph.series_order_dimension":null,"graph.show_values":true,"graph.x_axis.axis_enabled":true,"graph.x_axis.labels_enabled":true,"graph.x_axis.scale":"ordinal","graph.y_axis.auto_range":true,"stackable.stack_type":"stacked","table.cell_column":"SUBTOTAL"}	13371338	1	5	query	f	1	\N	\N	f	\N	\N	[{"field_ref":["field","Age",{"base-type":"type/BigInteger","binning":{"strategy":"num-bins","num-bins":10,"min-value":20.0,"max-value":65.0,"bin-width":5.0}}],"base_type":"type/BigInteger","name":"Age","effective_type":"type/BigInteger","display_name":"Age","fingerprint":{"global":{"distinct-count":42,"nil%":0},"type":{"type/Number":{"min":24,"q1":33.36752836803635,"q3":55.20362176071121,"max":65,"sd":12.063315373018085,"avg":44.572}}},"binning_info":{"num_bins":10,"min_value":20.0,"max_value":65.0,"bin_width":5.0,"binning_strategy":"num-bins"},"source":"breakout"},{"description":"The type of product, valid values include: Doohicky, Gadget, Gizmo and Widget","semantic_type":"type/Category","table_id":8,"coercion_strategy":null,"name":"CATEGORY","settings":null,"source":"breakout","fk_target_field_id":null,"field_ref":["field",58,{"base-type":"type/Text","join-alias":"Products"}],"effective_type":"type/Text","nfc_path":null,"parent_id":null,"id":58,"position":3,"visibility_type":"normal","display_name":"Products → Category","fingerprint":{"global":{"distinct-count":4,"nil%":0.0},"type":{"type/Text":{"percent-json":0.0,"percent-url":0.0,"percent-email":0.0,"percent-state":0.0,"average-length":6.375}}},"base_type":"type/Text","source_alias":"Products"},{"base_type":"type/Integer","name":"count","display_name":"Count","semantic_type":"type/Quantity","source":"aggregation","field_ref":["aggregation",0],"aggregation_index":0}]	\N	avC199ZiYvp1aIWdejfFf	[]	[]	t	\N	question	\N	\N	2024-06-22 08:19:17.839482+00	0
12	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	Revenue and orders over time	Matches the cumulative revenue month over month with the number of orders placed each month	combo	{"database":1,"type":"query","query":{"aggregation":[["sum",["field",42,{"base-type":"type/Float"}]],["sum",["field",39,{"base-type":"type/Integer"}]]],"breakout":[["field",41,{"base-type":"type/DateTime","temporal-unit":"month"}]],"source-table":5,"filter":["time-interval",["field",41,{"base-type":"type/DateTime"}],-24,"month"]}}	{"column_settings":null,"graph.dimensions":["CREATED_AT"],"graph.metrics":["sum","sum_2"],"graph.show_trendline":false,"graph.x_axis.title_text":"Orders date","graph.y_axis.title_text":"Revenue","series_settings":{"sum":{"display":"line","line.interpolate":"linear","line.marker_enabled":false,"show_series_values":true,"title":"Revenue"},"sum_2":{"color":"#51528D","title":"Number of orders"}}}	13371338	1	5	query	f	1	\N	\N	f	\N	\N	[{"description":"The date and time an order was submitted.","semantic_type":"type/CreationTimestamp","coercion_strategy":null,"unit":"month","name":"CREATED_AT","settings":null,"fk_target_field_id":null,"field_ref":["field",41,{"base-type":"type/DateTime","temporal-unit":"month"}],"effective_type":"type/DateTime","id":41,"visibility_type":"normal","display_name":"Created At","fingerprint":{"global":{"distinct-count":10001,"nil%":0.0},"type":{"type/DateTime":{"earliest":"2022-04-30T18:56:13.352Z","latest":"2026-04-19T14:07:15.657Z"}}},"base_type":"type/DateTime"},{"display_name":"Sum of Total","semantic_type":null,"settings":null,"field_ref":["aggregation",0],"name":"sum","base_type":"type/Float","effective_type":"type/Float","fingerprint":{"global":{"distinct-count":24,"nil%":0.0},"type":{"type/Number":{"min":1265.7162964063327,"q1":7814.799491244121,"q3":21566.481581896165,"max":38569.69761756364,"sd":11384.837065416124,"avg":16481.762292966578}}}},{"display_name":"Sum of Quantity","semantic_type":"type/Quantity","settings":null,"field_ref":["aggregation",1],"name":"sum_2","base_type":"type/BigInteger","effective_type":"type/BigInteger","fingerprint":{"global":{"distinct-count":24,"nil%":0.0},"type":{"type/Number":{"min":75.0,"q1":369.0,"q3":1286.0,"max":2260.0,"sd":630.7594628699596,"avg":877.25}}}}]	\N	evd6f3MoDUtTNsoXHwu0T	[]	[]	t	\N	question	\N	\N	2024-06-22 08:19:17.839482+00	0
13	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	Revenue goal for this quarter	Compares the current total in revenue with the set goal for the quarter	progress	{"database":1,"type":"query","query":{"aggregation":[["cum-sum",["field",42,{"base-type":"type/Float"}]]],"source-table":5,"filter":["time-interval",["field",41,{"base-type":"type/DateTime"}],"current","quarter"]}}	{"column_settings":null,"progress.goal":120000}	13371338	1	5	query	f	1	\N	\N	f	\N	\N	[{"display_name":"Cumulative sum of Total","semantic_type":null,"settings":null,"field_ref":["aggregation",0],"name":"sum","base_type":"type/Float","effective_type":"type/Float","fingerprint":{"global":{"distinct-count":1,"nil%":0.0},"type":{"type/Number":{"min":123889.60365320997,"q1":123889.60365320997,"q3":123889.60365320997,"max":123889.60365320997,"sd":null,"avg":123889.60365320997}}}}]	\N	Bkb4GEr5dH2_LbDOKhObR	[]	[]	t	\N	question	\N	\N	2024-06-22 08:19:17.839482+00	0
14	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	Numbers of orders by category per quarter	Compares the orders of each category quarter over quarter	line	{"database":1,"type":"query","query":{"aggregation":[["count"]],"breakout":[["field",58,{"base-type":"type/Text","source-field":40}],["field",41,{"base-type":"type/DateTime","temporal-unit":"quarter"}]],"source-table":5,"filter":["time-interval",["field",41,{"base-type":"type/DateTime"}],-24,"month"]}}	{"graph.dimensions":["CREATED_AT","CATEGORY"],"graph.metrics":["count"],"table.pivot_column":"CATEGORY","table.cell_column":"count"}	13371338	1	5	query	f	1	\N	\N	f	\N	\N	[{"description":"The type of product, valid values include: Doohicky, Gadget, Gizmo and Widget","semantic_type":"type/Category","coercion_strategy":null,"name":"CATEGORY","settings":null,"fk_target_field_id":null,"field_ref":["field",58,{"base-type":"type/Text","source-field":40}],"effective_type":"type/Text","id":58,"visibility_type":"normal","display_name":"Product → Category","fingerprint":{"global":{"distinct-count":4,"nil%":0.0},"type":{"type/Text":{"percent-json":0.0,"percent-url":0.0,"percent-email":0.0,"percent-state":0.0,"average-length":6.375}}},"base_type":"type/Text"},{"description":"The date and time an order was submitted.","semantic_type":"type/CreationTimestamp","coercion_strategy":null,"unit":"quarter","name":"CREATED_AT","settings":null,"fk_target_field_id":null,"field_ref":["field",41,{"base-type":"type/DateTime","temporal-unit":"quarter"}],"effective_type":"type/DateTime","id":41,"visibility_type":"normal","display_name":"Created At","fingerprint":{"global":{"distinct-count":10001,"nil%":0.0},"type":{"type/DateTime":{"earliest":"2022-04-30T18:56:13.352Z","latest":"2026-04-19T14:07:15.657Z"}}},"base_type":"type/DateTime"},{"display_name":"Count","semantic_type":"type/Quantity","field_ref":["aggregation",0],"name":"count","base_type":"type/BigInteger","effective_type":"type/BigInteger","fingerprint":{"global":{"distinct-count":34,"nil%":0.0},"type":{"type/Number":{"min":13.0,"q1":81.0,"q3":259.5,"max":338.0,"sd":102.85448188797824,"avg":169.38888888888889}}}}]	\N	NdqehUq98ebQ4H_-1d9v3	[]	[]	t	\N	question	\N	\N	2024-06-22 08:19:17.839482+00	0
15	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	Revenue per quarter	Compares the total revenue this quarter with the previous period	smartscalar	{"database":1,"type":"query","query":{"aggregation":[["sum",["field",42,{"base-type":"type/Float"}]]],"breakout":[["field",41,{"base-type":"type/DateTime","temporal-unit":"quarter"}]],"source-table":5,"filter":["time-interval",["field",41,{"base-type":"type/DateTime"}],-2,"quarter"]}}	{"column_settings":null}	13371338	1	5	query	f	1	\N	\N	f	\N	\N	[{"description":"The date and time an order was submitted.","semantic_type":"type/CreationTimestamp","coercion_strategy":null,"unit":"quarter","name":"CREATED_AT","settings":null,"fk_target_field_id":null,"field_ref":["field",41,{"base-type":"type/DateTime","temporal-unit":"quarter"}],"effective_type":"type/DateTime","id":41,"visibility_type":"normal","display_name":"Created At","fingerprint":{"global":{"distinct-count":10001,"nil%":0.0},"type":{"type/DateTime":{"earliest":"2022-04-30T18:56:13.352Z","latest":"2026-04-19T14:07:15.657Z"}}},"base_type":"type/DateTime"},{"display_name":"Sum of Total","semantic_type":null,"settings":null,"field_ref":["aggregation",0],"name":"sum","base_type":"type/Float","effective_type":"type/Float","fingerprint":{"global":{"distinct-count":2,"nil%":0.0},"type":{"type/Number":{"min":67863.39650548704,"q1":67863.39650548704,"q3":111121.25780141751,"max":111121.25780141751,"sd":30587.92706197953,"avg":89492.32715345226}}}}]	\N	d3KPtMfbe0IIK0MDrHTNn	[]	[]	t	\N	question	\N	\N	2024-06-22 08:19:17.839482+00	0
16	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	Orders according to sources per quarter	Orders placed per quarter broken down by source and formatted to highlight best and worst quarters	pivot	{"database":1,"type":"query","query":{"aggregation":[["sum",["field",44,{"base-type":"type/Float"}]]],"breakout":[["field",45,{"base-type":"type/Text","source-field":43}],["field",41,{"base-type":"type/DateTime","temporal-unit":"quarter"}]],"source-table":5,"filter":["time-interval",["field",41,{"base-type":"type/DateTime"}],-24,"month"]}}	{"graph.dimensions":["CREATED_AT","SOURCE"],"graph.series_labels":[null],"pivot_table.column_split":{"columns":[["field",45,{"base-type":"type/Text","source-field":43}]],"rows":[["field",41,{"base-type":"type/DateTime","temporal-unit":"quarter"}]],"values":[["aggregation",0]]},"pivot_table.column_widths":{"leftHeaderWidths":[141],"totalLeftHeaderWidths":141,"valueHeaderWidths":{}},"stackable.stack_type":"stacked","table.column_formatting":[{"color":"#A7D07C","columns":["sum"],"highlight_row":false,"operator":">","type":"single","value":20000}],"graph.metrics":["sum"]}	13371338	1	5	query	f	1	\N	\N	f	\N	\N	[{"description":"The channel through which we acquired this user. Valid values include: Affiliate, Facebook, Google, Organic and Twitter","semantic_type":null,"coercion_strategy":null,"name":"pivot-grouping","settings":null,"fk_target_field_id":null,"field_ref":["expression","pivot-grouping"],"effective_type":"type/Text","id":45,"visibility_type":"normal","display_name":"pivot-grouping","fingerprint":{"global":{"distinct-count":1,"nil%":0.0},"type":{"type/Number":{"min":3.0,"q1":3.0,"q3":3.0,"max":3.0,"sd":null,"avg":3.0}}},"base_type":"type/Integer"},{"description":"The date and time an order was submitted.","semantic_type":null,"coercion_strategy":null,"name":"sum","settings":null,"fk_target_field_id":null,"field_ref":["aggregation",0],"effective_type":"type/DateTime","id":41,"visibility_type":"normal","display_name":"Sum of Subtotal","fingerprint":{"global":{"distinct-count":1,"nil%":0.0},"type":{"type/Number":{"min":379867.38872485986,"q1":379867.38872485986,"q3":379867.38872485986,"max":379867.38872485986,"sd":null,"avg":379867.38872485986}}},"base_type":"type/Float"}]	\N	-mfPjg0QOoIKbqJS0SYne	[]	[]	t	\N	question	\N	\N	2024-06-22 08:19:17.839482+00	0
17	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	Product category orders per individual age	Shows a distribution of orders broken down by product category across our customers' individual age values	bar	{"database":1,"type":"query","query":{"aggregation":[["count"]],"breakout":[["expression","Age",{"base-type":"type/Integer"}],["field",58,{"base-type":"type/Text","source-field":40}]],"expressions":{"Age":["datetime-diff",["field",49,{"base-type":"type/Date","join-alias":"People - User"}],["now"],"year"]},"joins":[{"alias":"People - User","strategy":"left-join","condition":["=",["field",43,{"base-type":"type/Integer"}],["field",46,{"base-type":"type/BigInteger","join-alias":"People - User"}]],"source-table":3}],"source-table":5}}	{"column_settings":null,"graph.dimensions":["Age","CATEGORY"],"graph.metrics":["count"],"graph.series_order":null,"graph.series_order_dimension":null,"stackable.stack_type":"stacked"}	13371338	1	5	query	f	1	\N	\N	f	\N	\N	[{"display_name":"Age","field_ref":["expression","Age"],"name":"Age","base_type":"type/BigInteger","effective_type":"type/BigInteger","semantic_type":null,"fingerprint":{"global":{"distinct-count":43,"nil%":0.0},"type":{"type/Number":{"min":24.0,"q1":34.25,"q3":55.75,"max":66.0,"sd":12.445906346880554,"avg":45.0}}}},{"description":"The type of product, valid values include: Doohicky, Gadget, Gizmo and Widget","semantic_type":"type/Category","coercion_strategy":null,"name":"CATEGORY","settings":null,"fk_target_field_id":null,"field_ref":["field",58,{"base-type":"type/Text","source-field":40}],"effective_type":"type/Text","id":58,"visibility_type":"normal","display_name":"Product → Category","fingerprint":{"global":{"distinct-count":4,"nil%":0.0},"type":{"type/Text":{"percent-json":0.0,"percent-url":0.0,"percent-email":0.0,"percent-state":0.0,"average-length":6.375}}},"base_type":"type/Text"},{"display_name":"Count","semantic_type":"type/Quantity","field_ref":["aggregation",0],"name":"count","base_type":"type/BigInteger","effective_type":"type/BigInteger","fingerprint":{"global":{"distinct-count":84,"nil%":0.0},"type":{"type/Number":{"min":11.0,"q1":95.5,"q3":124.48331477354789,"max":204.0,"sd":28.015985634286857,"avg":109.06976744186046}}}}]	\N	P-shWkp5lozD2nuo5xusW	[]	[]	t	\N	question	\N	\N	2024-06-22 08:19:17.839482+00	0
18	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	Discount given per quarter	Compares the total discount given this quarter with the previous period	smartscalar	{"database":1,"type":"query","query":{"aggregation":[["sum",["field",36,{"base-type":"type/Float"}]]],"breakout":[["field",41,{"base-type":"type/DateTime","temporal-unit":"quarter"}]],"source-table":5,"filter":["time-interval",["field",41,{"base-type":"type/DateTime"}],-2,"quarter"]}}	{"column_settings":null}	13371338	1	5	query	f	1	\N	\N	f	\N	\N	[{"description":"The date and time an order was submitted.","semantic_type":"type/CreationTimestamp","coercion_strategy":null,"unit":"quarter","name":"CREATED_AT","settings":null,"fk_target_field_id":null,"field_ref":["field",41,{"base-type":"type/DateTime","temporal-unit":"quarter"}],"effective_type":"type/DateTime","id":41,"visibility_type":"normal","display_name":"Created At","fingerprint":{"global":{"distinct-count":10001,"nil%":0.0},"type":{"type/DateTime":{"earliest":"2022-04-30T18:56:13.352Z","latest":"2026-04-19T14:07:15.657Z"}}},"base_type":"type/DateTime"},{"display_name":"Sum of Discount","semantic_type":"type/Discount","settings":null,"field_ref":["aggregation",0],"name":"sum","base_type":"type/Float","effective_type":"type/Float","fingerprint":{"global":{"distinct-count":2,"nil%":0.0},"type":{"type/Number":{"min":738.5575011690233,"q1":738.5575011690233,"q3":795.9833992557524,"max":795.9833992557524,"sd":40.606241952853686,"avg":767.2704502123879}}}}]	\N	wd1PEi_Z217Gr-ceeZSNj	[]	[]	t	\N	question	\N	\N	2024-06-22 08:19:17.839482+00	0
19	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	Customer satisfaction per category	Shows the distribution of the product categories along the scale of customer ratings	bar	{"database":1,"type":"query","query":{"aggregation":[["count"]],"breakout":[["field",61,{"base-type":"type/Float","source-field":40}],["field",58,{"base-type":"type/Text","source-field":40}]],"order-by":[["desc",["field",61,{"base-type":"type/Float","source-field":40}]]],"source-table":5,"filter":["!=",["field",65,{"base-type":"type/Text","source-field":40}],"Incredible Aluminum Knife"]}}	{"column_settings":null,"graph.dimensions":["RATING","CATEGORY"],"graph.metrics":["count"],"graph.series_order":null,"graph.series_order_dimension":null,"graph.x_axis.title_text":"Products","graph.y_axis.auto_range":true,"graph.y_axis.title_text":"total orders","series_settings":{"count":{"color":"#999AC4"}},"stackable.stack_type":"stacked","table.pivot":false}	13371338	1	5	query	f	1	\N	\N	f	\N	\N	[{"description":"The average rating users have given the product. This ranges from 1 - 5","semantic_type":"type/Score","coercion_strategy":null,"name":"RATING","settings":null,"fk_target_field_id":null,"field_ref":["field",61,{"base-type":"type/Float","source-field":40}],"effective_type":"type/Float","id":61,"visibility_type":"normal","display_name":"Product → Rating","fingerprint":{"global":{"distinct-count":23,"nil%":0.0},"type":{"type/Number":{"min":0.0,"q1":3.5120465053408525,"q3":4.216124969497314,"max":5.0,"sd":1.3605488657451452,"avg":3.4715}}},"base_type":"type/Float"},{"description":"The type of product, valid values include: Doohicky, Gadget, Gizmo and Widget","semantic_type":"type/Category","coercion_strategy":null,"name":"CATEGORY","settings":null,"fk_target_field_id":null,"field_ref":["field",58,{"base-type":"type/Text","source-field":40}],"effective_type":"type/Text","id":58,"visibility_type":"normal","display_name":"Product → Category","fingerprint":{"global":{"distinct-count":4,"nil%":0.0},"type":{"type/Text":{"percent-json":0.0,"percent-url":0.0,"percent-email":0.0,"percent-state":0.0,"average-length":6.375}}},"base_type":"type/Text"},{"display_name":"Count","semantic_type":"type/Quantity","field_ref":["aggregation",0],"name":"count","base_type":"type/BigInteger","effective_type":"type/BigInteger","fingerprint":{"global":{"distinct-count":57,"nil%":0.0},"type":{"type/Number":{"min":70.0,"q1":105.0,"q3":370.0,"max":1093.0,"sd":229.88812521865958,"avg":281.6060606060606}}}}]	\N	cO53DQ6k0DueW5Aglwp7Y	[]	[]	t	\N	question	\N	\N	2024-06-22 08:19:17.839482+00	0
20	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	Number of orders per month	Compares the total number of orders placed this month with the previous period	smartscalar	{"database":1,"type":"query","query":{"aggregation":[["sum",["field",39,{"base-type":"type/Integer"}]]],"breakout":[["field",41,{"base-type":"type/DateTime","temporal-unit":"month"}]],"source-table":5,"filter":["time-interval",["field",41,{"base-type":"type/DateTime"}],-2,"month"]}}	{"column_settings":null}	13371338	1	5	query	f	1	\N	\N	f	\N	\N	[{"description":"The date and time an order was submitted.","semantic_type":"type/CreationTimestamp","coercion_strategy":null,"unit":"month","name":"CREATED_AT","settings":null,"fk_target_field_id":null,"field_ref":["field",41,{"base-type":"type/DateTime","temporal-unit":"month"}],"effective_type":"type/DateTime","id":41,"visibility_type":"normal","display_name":"Created At","fingerprint":{"global":{"distinct-count":10001,"nil%":0.0},"type":{"type/DateTime":{"earliest":"2022-04-30T18:56:13.352Z","latest":"2026-04-19T14:07:15.657Z"}}},"base_type":"type/DateTime"},{"display_name":"Sum of Quantity","semantic_type":"type/Quantity","settings":null,"field_ref":["aggregation",0],"name":"sum","base_type":"type/BigInteger","effective_type":"type/BigInteger","fingerprint":{"global":{"distinct-count":2,"nil%":0.0},"type":{"type/Number":{"min":487.0,"q1":487.0,"q3":1330.0,"max":1330.0,"sd":596.0910165402596,"avg":908.5}}}}]	\N	OouzOUPGIjf6zXJsSDe0f	[]	[]	t	\N	question	\N	\N	2024-06-22 08:19:17.839482+00	0
21	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	Orders by source per individual age	Shows a distribution of orders broken down by source across our customers' individual age values	bar	{"database":1,"type":"query","query":{"aggregation":[["count"]],"breakout":[["expression","Age",{"base-type":"type/Integer"}],["field",45,{"base-type":"type/Text","join-alias":"People - User"}]],"expressions":{"Age":["datetime-diff",["field",49,{"base-type":"type/Date","join-alias":"People - User"}],["now"],"year"]},"joins":[{"alias":"People - User","strategy":"left-join","condition":["=",["field",43,{"base-type":"type/Integer"}],["field",46,{"base-type":"type/BigInteger","join-alias":"People - User"}]],"source-table":3}],"source-table":5}}	{"column_settings":null,"graph.dimensions":["Age","SOURCE"],"graph.metrics":["count"],"graph.series_order":null,"graph.series_order_dimension":null,"stackable.stack_type":"stacked"}	13371338	1	5	query	f	1	\N	\N	f	\N	\N	[{"display_name":"Age","field_ref":["expression","Age"],"name":"Age","base_type":"type/BigInteger","effective_type":"type/BigInteger","semantic_type":null,"fingerprint":{"global":{"distinct-count":43,"nil%":0.0},"type":{"type/Number":{"min":24.0,"q1":34.15,"q3":55.45,"max":66.0,"sd":12.328008896346248,"avg":44.80281690140845}}}},{"description":"The channel through which we acquired this user. Valid values include: Affiliate, Facebook, Google, Organic and Twitter","semantic_type":"type/Source","coercion_strategy":null,"name":"SOURCE","settings":null,"fk_target_field_id":null,"field_ref":["field",45,{"base-type":"type/Text","join-alias":"People - User"}],"effective_type":"type/Text","id":45,"visibility_type":"normal","display_name":"People - User → Source","fingerprint":{"global":{"distinct-count":5,"nil%":0.0},"type":{"type/Text":{"percent-json":0.0,"percent-url":0.0,"percent-email":0.0,"percent-state":0.0,"average-length":7.4084}}},"base_type":"type/Text"},{"display_name":"Count","semantic_type":"type/Quantity","field_ref":["aggregation",0],"name":"count","base_type":"type/BigInteger","effective_type":"type/BigInteger","fingerprint":{"global":{"distinct-count":111,"nil%":0.0},"type":{"type/Number":{"min":14.0,"q1":61.30452018978108,"q3":110.7900825561565,"max":214.0,"sd":36.95993697001328,"avg":88.07511737089202}}}}]	\N	DDT6zuq6O7bsnw1UEWtQQ	[]	[]	t	\N	question	\N	\N	2024-06-22 08:19:17.839482+00	0
22	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	Average product rating overall	Indicates the average customer review of our products	gauge	{"database":1,"type":"query","query":{"aggregation":[["avg",["field",61,{"base-type":"type/Float"}]]],"source-table":8}}	{"column_settings":null,"gauge.segments":[{"color":"#EF8C8C","label":"awful","max":1,"min":0},{"color":"#F2A86F","label":"bad","max":2,"min":1},{"color":"#F9D45C","label":"alright","max":3,"min":2},{"color":"#A7D07C","label":"good","max":4,"min":3},{"color":"#689636","label":"great","max":5,"min":4}],"table.cell_column":"avg"}	13371338	1	8	query	f	1	\N	\N	f	\N	\N	[{"display_name":"Average of Rating","semantic_type":"type/Score","settings":null,"field_ref":["aggregation",0],"name":"avg","base_type":"type/Float","effective_type":"type/Float","fingerprint":{"global":{"distinct-count":1,"nil%":0.0},"type":{"type/Number":{"min":3.4715,"q1":3.4715,"q3":3.4715,"max":3.4715,"sd":null,"avg":3.4715}}}}]	\N	-BcEE22ptXyDOgUNRqyct	[]	[]	t	\N	question	\N	\N	2024-06-22 08:19:17.839482+00	0
23	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	Best selling products	An ordered list of our most successful products	row	{"database":1,"type":"query","query":{"filter":[">",["field","count",{"base-type":"type/Integer"}],108],"source-query":{"aggregation":[["count"]],"breakout":[["field",65,{"base-type":"type/Text","source-field":40}],["field",58,{"base-type":"type/Text","source-field":40}]],"order-by":[["desc",["aggregation",0]]],"source-table":5,"filter":["!=",["field",65,{"base-type":"type/Text","source-field":40}],"Incredible Aluminum Knife"]}}}	{"column_settings":null,"graph.dimensions":["TITLE"],"graph.metrics":["count"],"graph.x_axis.title_text":"Products","graph.y_axis.title_text":"total orders","series_settings":{"count":{"color":"#999AC4"}},"table.pivot":false}	13371338	1	5	query	f	1	\N	\N	f	\N	\N	[{"description":"The name of the product as it should be displayed to customers.","semantic_type":"type/Title","coercion_strategy":null,"name":"TITLE","settings":null,"fk_target_field_id":null,"field_ref":["field",65,{"base-type":"type/Text","source-field":40}],"effective_type":"type/Text","id":65,"visibility_type":"normal","display_name":"Product → Title","fingerprint":{"global":{"distinct-count":199,"nil%":0.0},"type":{"type/Text":{"percent-json":0.0,"percent-url":0.0,"percent-email":0.0,"percent-state":0.0,"average-length":21.495}}},"base_type":"type/Text"},{"description":"The type of product, valid values include: Doohicky, Gadget, Gizmo and Widget","semantic_type":"type/Category","coercion_strategy":null,"name":"CATEGORY","settings":null,"fk_target_field_id":null,"field_ref":["field",58,{"base-type":"type/Text","source-field":40}],"effective_type":"type/Text","id":58,"visibility_type":"normal","display_name":"Product → Category","fingerprint":{"global":{"distinct-count":4,"nil%":0.0},"type":{"type/Text":{"percent-json":0.0,"percent-url":0.0,"percent-email":0.0,"percent-state":0.0,"average-length":6.375}}},"base_type":"type/Text"},{"display_name":"Count","semantic_type":"type/Quantity","field_ref":["field","count",{"base-type":"type/Integer"}],"name":"count","base_type":"type/BigInteger","effective_type":"type/BigInteger","fingerprint":{"global":{"distinct-count":8,"nil%":0.0},"type":{"type/Number":{"min":109.0,"q1":109.36150801752578,"q3":117.41886116991581,"max":120.0,"sd":4.254710813035841,"avg":113.53846153846153}}}}]	\N	dsZMohGmK4P7eQRBdLFmb	[]	[]	t	\N	question	\N	\N	2024-06-22 08:19:17.839482+00	0
24	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	Unique customers per month	Compares the number of unique customers placing orders this month with the previous period	smartscalar	{"database":1,"type":"query","query":{"aggregation":[["distinct",["field",47,{"base-type":"type/Text"}]]],"breakout":[["field",56,{"base-type":"type/DateTime","temporal-unit":"month"}]],"source-table":3,"filter":["time-interval",["field",56,{"base-type":"type/DateTime"}],-2,"month"]}}	{"column_settings":null}	13371338	1	3	query	f	1	\N	\N	f	\N	\N	[{"description":"The date the user record was created. Also referred to as the user’s \\"join date\\"","semantic_type":"type/CreationTimestamp","coercion_strategy":null,"unit":"month","name":"CREATED_AT","settings":null,"fk_target_field_id":null,"field_ref":["field",56,{"base-type":"type/DateTime","temporal-unit":"month"}],"effective_type":"type/DateTime","id":56,"visibility_type":"normal","display_name":"Created At","fingerprint":{"global":{"distinct-count":2500,"nil%":0.0},"type":{"type/DateTime":{"earliest":"2022-04-19T21:35:18.752Z","latest":"2025-04-19T14:06:27.3Z"}}},"base_type":"type/DateTime"},{"display_name":"Distinct values of Name","semantic_type":"type/Quantity","settings":null,"field_ref":["aggregation",0],"name":"count","base_type":"type/BigInteger","effective_type":"type/BigInteger","fingerprint":{"global":{"distinct-count":2,"nil%":0.0},"type":{"type/Number":{"min":70.0,"q1":70.0,"q3":76.0,"max":76.0,"sd":4.242640687119285,"avg":73.0}}}}]	\N	fPYpb3KcV2TJZLRw-L841	[]	[]	t	\N	question	\N	\N	2024-06-22 08:19:17.839482+00	0
25	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	Revenue per individual age	Shows a distribution of revenue per individual age values	bar	{"database":1,"type":"query","query":{"aggregation":[["sum",["field",42,{"base-type":"type/Float"}]]],"breakout":[["expression","Age",{"base-type":"type/Integer"}]],"expressions":{"Age":["datetime-diff",["field",49,{"base-type":"type/Date","join-alias":"People - User"}],["now"],"year"]},"joins":[{"alias":"People - User","strategy":"left-join","condition":["=",["field",43,{"base-type":"type/Integer"}],["field",46,{"base-type":"type/BigInteger","join-alias":"People - User"}]],"source-table":3}],"source-table":5}}	{"column_settings":null,"graph.dimensions":["Age"],"graph.metrics":["sum"],"graph.series_order":null,"graph.series_order_dimension":null,"graph.show_values":true,"stackable.stack_type":"stacked"}	13371338	1	5	query	f	1	\N	\N	f	\N	\N	[{"display_name":"Age","field_ref":["expression","Age"],"name":"Age","base_type":"type/BigInteger","effective_type":"type/BigInteger","semantic_type":null,"fingerprint":{"global":{"distinct-count":43,"nil%":0.0},"type":{"type/Number":{"min":24.0,"q1":34.25,"q3":55.75,"max":66.0,"sd":12.556538801224908,"avg":45.0}}}},{"display_name":"Sum of Total","semantic_type":null,"settings":null,"field_ref":["aggregation",0],"name":"sum","base_type":"type/Float","effective_type":"type/Float","fingerprint":{"global":{"distinct-count":43,"nil%":0.0},"type":{"type/Number":{"min":4158.786206744129,"q1":31397.769615345165,"q3":38764.15212640109,"max":58676.42899556268,"sd":8106.619825736762,"avg":35130.73681513096}}}}]	\N	KMq74J_Gh78vIimczx0LH	[]	[]	t	\N	question	\N	\N	2024-06-22 08:19:17.839482+00	0
26	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	People with age	\N	table	{"database":1,"type":"query","query":{"source-table":3,"expressions":{"Age":["datetime-diff",["field",49,{"base-type":"type/Date"}],["now"],"year"]}}}	{"graph.dimensions":["age"],"graph.metrics":["count"],"table.pivot_column":"SOURCE","table.cell_column":"LONGITUDE"}	13371338	1	3	query	f	1	\N	\N	f	\N	\N	[{"description":"A unique identifier given to each user.","semantic_type":"type/PK","coercion_strategy":null,"name":"ID","settings":null,"fk_target_field_id":null,"field_ref":["field",46,null],"effective_type":"type/BigInteger","id":46,"visibility_type":"normal","display_name":"ID","fingerprint":null,"base_type":"type/BigInteger"},{"description":"The street address of the account’s billing address","semantic_type":null,"coercion_strategy":null,"name":"ADDRESS","settings":null,"fk_target_field_id":null,"field_ref":["field",52,null],"effective_type":"type/Text","id":52,"visibility_type":"normal","display_name":"Address","fingerprint":{"global":{"distinct-count":2490,"nil%":0.0},"type":{"type/Text":{"percent-json":0.0,"percent-url":0.0,"percent-email":0.0,"percent-state":0.0,"average-length":20.85}}},"base_type":"type/Text"},{"description":"The contact email for the account.","semantic_type":"type/Email","coercion_strategy":null,"name":"EMAIL","settings":null,"fk_target_field_id":null,"field_ref":["field",51,null],"effective_type":"type/Text","id":51,"visibility_type":"normal","display_name":"Email","fingerprint":{"global":{"distinct-count":2500,"nil%":0.0},"type":{"type/Text":{"percent-json":0.0,"percent-url":0.0,"percent-email":1.0,"percent-state":0.0,"average-length":24.1824}}},"base_type":"type/Text"},{"description":"This is the salted password of the user. It should not be visible","semantic_type":null,"coercion_strategy":null,"name":"PASSWORD","settings":null,"fk_target_field_id":null,"field_ref":["field",53,null],"effective_type":"type/Text","id":53,"visibility_type":"normal","display_name":"Password","fingerprint":{"global":{"distinct-count":2500,"nil%":0.0},"type":{"type/Text":{"percent-json":0.0,"percent-url":0.0,"percent-email":0.0,"percent-state":0.0,"average-length":36.0}}},"base_type":"type/Text"},{"description":"The name of the user who owns an account","semantic_type":"type/Name","coercion_strategy":null,"name":"NAME","settings":null,"fk_target_field_id":null,"field_ref":["field",47,null],"effective_type":"type/Text","id":47,"visibility_type":"normal","display_name":"Name","fingerprint":{"global":{"distinct-count":2499,"nil%":0.0},"type":{"type/Text":{"percent-json":0.0,"percent-url":0.0,"percent-email":0.0,"percent-state":0.0,"average-length":13.532}}},"base_type":"type/Text"},{"description":"The city of the account’s billing address","semantic_type":"type/City","coercion_strategy":null,"name":"CITY","settings":null,"fk_target_field_id":null,"field_ref":["field",55,null],"effective_type":"type/Text","id":55,"visibility_type":"normal","display_name":"City","fingerprint":{"global":{"distinct-count":1966,"nil%":0.0},"type":{"type/Text":{"percent-json":0.0,"percent-url":0.0,"percent-email":0.0,"percent-state":0.002,"average-length":8.284}}},"base_type":"type/Text"},{"description":"This is the longitude of the user on sign-up. It might be updated in the future to the last seen location.","semantic_type":"type/Longitude","coercion_strategy":null,"name":"LONGITUDE","settings":null,"fk_target_field_id":null,"field_ref":["field",50,null],"effective_type":"type/Float","id":50,"visibility_type":"normal","display_name":"Longitude","fingerprint":{"global":{"distinct-count":2491,"nil%":0.0},"type":{"type/Number":{"min":-166.5425726,"q1":-101.58350792373135,"q3":-84.65289348288829,"max":-67.96735199999999,"sd":15.399698968175663,"avg":-95.18741780363999}}},"base_type":"type/Float"},{"description":"The state or province of the account’s billing address","semantic_type":"type/State","coercion_strategy":null,"name":"STATE","settings":null,"fk_target_field_id":null,"field_ref":["field",48,null],"effective_type":"type/Text","id":48,"visibility_type":"normal","display_name":"State","fingerprint":{"global":{"distinct-count":49,"nil%":0.0},"type":{"type/Text":{"percent-json":0.0,"percent-url":0.0,"percent-email":0.0,"percent-state":1.0,"average-length":2.0}}},"base_type":"type/Text"},{"description":"The channel through which we acquired this user. Valid values include: Affiliate, Facebook, Google, Organic and Twitter","semantic_type":"type/Source","coercion_strategy":null,"name":"SOURCE","settings":null,"fk_target_field_id":null,"field_ref":["field",45,null],"effective_type":"type/Text","id":45,"visibility_type":"normal","display_name":"Source","fingerprint":{"global":{"distinct-count":5,"nil%":0.0},"type":{"type/Text":{"percent-json":0.0,"percent-url":0.0,"percent-email":0.0,"percent-state":0.0,"average-length":7.4084}}},"base_type":"type/Text"},{"description":"The date of birth of the user","semantic_type":null,"coercion_strategy":null,"unit":"default","name":"BIRTH_DATE","settings":null,"fk_target_field_id":null,"field_ref":["field",49,{"temporal-unit":"default"}],"effective_type":"type/Date","id":49,"visibility_type":"normal","display_name":"Birth Date","fingerprint":{"global":{"distinct-count":2308,"nil%":0.0},"type":{"type/DateTime":{"earliest":"1958-04-26","latest":"2000-04-03"}}},"base_type":"type/Date"},{"description":"The postal code of the account’s billing address","semantic_type":"type/ZipCode","coercion_strategy":null,"name":"ZIP","settings":null,"fk_target_field_id":null,"field_ref":["field",54,null],"effective_type":"type/Text","id":54,"visibility_type":"normal","display_name":"Zip","fingerprint":{"global":{"distinct-count":2234,"nil%":0.0},"type":{"type/Text":{"percent-json":0.0,"percent-url":0.0,"percent-email":0.0,"percent-state":0.0,"average-length":5.0}}},"base_type":"type/Text"},{"description":"This is the latitude of the user on sign-up. It might be updated in the future to the last seen location.","semantic_type":"type/Latitude","coercion_strategy":null,"name":"LATITUDE","settings":null,"fk_target_field_id":null,"field_ref":["field",57,null],"effective_type":"type/Float","id":57,"visibility_type":"normal","display_name":"Latitude","fingerprint":{"global":{"distinct-count":2491,"nil%":0.0},"type":{"type/Number":{"min":25.775827,"q1":35.302705923023126,"q3":43.773802584662,"max":70.6355001,"sd":6.390832341883712,"avg":39.87934670484002}}},"base_type":"type/Float"},{"description":"The date the user record was created. Also referred to as the user’s \\"join date\\"","semantic_type":"type/CreationTimestamp","coercion_strategy":null,"unit":"default","name":"CREATED_AT","settings":null,"fk_target_field_id":null,"field_ref":["field",56,{"temporal-unit":"default"}],"effective_type":"type/DateTime","id":56,"visibility_type":"normal","display_name":"Created At","fingerprint":{"global":{"distinct-count":2500,"nil%":0.0},"type":{"type/DateTime":{"earliest":"2022-04-19T21:35:18.752Z","latest":"2025-04-19T14:06:27.3Z"}}},"base_type":"type/DateTime"},{"display_name":"Age","field_ref":["expression","Age"],"name":"Age","base_type":"type/BigInteger","effective_type":"type/BigInteger","semantic_type":null,"fingerprint":{"global":{"distinct-count":42,"nil%":0.0},"type":{"type/Number":{"min":24.0,"q1":33.340572873934306,"q3":55.17166516756599,"max":65.0,"sd":12.263883782175668,"avg":44.434}}}}]	\N	lY4hbjNofxepxQGRhJa0s	[]	[]	t	\N	question	\N	\N	2024-06-22 08:19:17.839482+00	0
\.


--
-- Data for Name: report_cardfavorite; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.report_cardfavorite (id, created_at, updated_at, card_id, owner_id) FROM stdin;
\.


--
-- Data for Name: report_dashboard; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.report_dashboard (id, created_at, updated_at, name, description, creator_id, parameters, points_of_interest, caveats, show_in_getting_started, public_uuid, made_public_by_id, enable_embedding, embedding_params, archived, "position", collection_id, collection_position, cache_ttl, entity_id, auto_apply_filters, width, initially_published_at, view_count) FROM stdin;
1	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	E-commerce insights	Quickly take an overview of an e-commerce reseller business and dive into separate tabs that focus on top selling products and demographic insights. Each vendor can log in as a tenant and see their own data sandboxed from all the others.	13371338	[{"id":"fc2cd1be","isMultiSelect":false,"name":"Vendor","sectionId":"string","slug":"vendor","type":"string/="},{"id":"afa56954","name":"Date Range","sectionId":"date","slug":"date_range","type":"date/range"},{"id":"5eeec658","name":"Category","sectionId":"string","slug":"category","type":"string/=","values_query_type":"list"},{"id":"512c560a","name":"Location","sectionId":"location","slug":"location","type":"string/=","values_query_type":"search"}]	\N	\N	f	\N	\N	f	\N	f	\N	1	2	\N	DlK2jXoIHPXyVkEuo6Uy6	t	full	\N	0
\.


--
-- Data for Name: report_dashboardcard; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.report_dashboardcard (id, created_at, updated_at, size_x, size_y, "row", col, card_id, dashboard_id, parameter_mappings, visualization_settings, entity_id, action_id, dashboard_tab_id) FROM stdin;
1	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	24	7	15	0	7	1	[{"parameter_id":"fc2cd1be","card_id":7,"target":["dimension",["field",60,{"base-type":"type/Text","source-field":40}]]},{"parameter_id":"afa56954","card_id":7,"target":["dimension",["field",41,{"base-type":"type/DateTime"}]]},{"parameter_id":"5eeec658","card_id":7,"target":["dimension",["field",58,{"base-type":"type/Text","source-field":40}]]},{"parameter_id":"512c560a","card_id":7,"target":["dimension",["field",48,{"base-type":"type/Text","source-field":43}]]}]	{"card.title":"Value of orders over time (before taxes)","column_settings":null,"graph.dimensions":["CREATED_AT","CATEGORY"],"graph.metrics":["sum"],"graph.series_order":null,"graph.series_order_dimension":null,"graph.x_axis.labels_enabled":false,"series_settings":{"Doohickey":{"color":"#7172AD"},"Gadget":{"color":"#A989C5"},"Gizmo":{"color":"#C7EAEA"},"Widget":{"color":"#227FD2"}},"stackable.stack_type":"stacked"}	T3G6Is3XniVtrqzKbZ2vx	\N	1
2	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	24	6	9	0	16	1	[{"parameter_id":"fc2cd1be","card_id":16,"target":["dimension",["field",60,{"base-type":"type/Text","source-field":40}]]},{"parameter_id":"afa56954","card_id":16,"target":["dimension",["field",41,{"base-type":"type/DateTime"}]]},{"parameter_id":"5eeec658","card_id":16,"target":["dimension",["field",58,{"base-type":"type/Text","source-field":40}]]},{"parameter_id":"512c560a","card_id":16,"target":["dimension",["field",48,{"base-type":"type/Text","source-field":43}]]}]	{"card.description":"Orders placed per quarter broken down by source and formatted to highlight best and worst quarters","column_settings":null,"graph.dimensions":["CREATED_AT","SOURCE"],"graph.metrics":["sum"],"graph.series_labels":[null],"pivot.show_row_totals":true,"pivot_table.column_split":{"columns":[["field",45,{"base-type":"type/Text","source-field":43}]],"rows":[["field",41,{"base-type":"type/DateTime","temporal-unit":"quarter"}]],"values":[["aggregation",0]]},"pivot_table.column_widths":{"leftHeaderWidths":[141],"totalLeftHeaderWidths":141,"valueHeaderWidths":{"0":234,"1":263,"2":254,"3":236.3636474609375,"4":227,"5":170}},"stackable.stack_type":"stacked","table.column_formatting":[{"color":"#FBE499","colors":["#ED6E6E","#FFFFFF","#84BB4C"],"columns":["sum"],"highlight_row":false,"max_type":null,"max_value":100,"min_type":null,"min_value":0,"operator":"<","type":"range","value":1000}]}	gXkQgKfsYghl-XjDpSefW	\N	1
3	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	24	1	0	0	\N	1	[]	{"column_settings":null,"dashcard.background":false,"text":"Overall business health","virtual_card":{"archived":false,"dataset_query":{},"display":"heading","name":null,"visualization_settings":{}}}	Cb92ioon7vWaER-4OcdDS	\N	1
4	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	15	6	1	9	4	1	[{"parameter_id":"512c560a","card_id":4,"target":["dimension",["field",48,{"base-type":"type/Text"}]]}]	{"graph.y_axis.title_text":"Buyers","graph.series_order_dimension":null,"graph.x_axis.title_text":"Age","card.title":"Buyers by age group","graph.metrics":["count_2"],"graph.series_order":null,"card.description":"Shows a distribution of our customers in age groups","series_settings":{"count":{"color":"#A989C5"},"count_2":{"color":"#999AC4"}},"graph.x_axis.scale":"ordinal","graph.dimensions":["count"]}	vJZR2mN2xRrXiqqwAc0eb	\N	3
5	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	24	1	0	0	\N	1	[]	{"column_settings":null,"dashcard.background":false,"text":"Top performing products","virtual_card":{"archived":false,"dataset_query":{},"display":"heading","name":null,"visualization_settings":{}}}	s2ZdKLKOql1B0gCbOh3XO	\N	2
6	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	7	5	13	17	6	1	[{"parameter_id":"fc2cd1be","card_id":6,"target":["dimension",["field",60,{"base-type":"type/Text","source-field":40}]]},{"parameter_id":"afa56954","card_id":6,"target":["dimension",["field",41,{"base-type":"type/DateTime"}]]},{"parameter_id":"5eeec658","card_id":6,"target":["dimension",["field",58,{"base-type":"type/Text","source-field":40}]]},{"parameter_id":"512c560a","card_id":6,"target":["dimension",["field",48,{"base-type":"type/Text","source-field":43}]]}]	{"card.description":"Breaks down the overall performance of each of the product categories ","card.title":"Total orders by category","column_settings":null,"pie.colors":{"Doohickey":"#7172AD","Gadget":"#A989C5","Gizmo":"#C7EAEA","Widget":"#227FD2"}}	C3ByFOdgRYsNwwXVTGfY-	\N	2
7	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	17	5	13	0	14	1	[{"parameter_id":"fc2cd1be","card_id":14,"target":["dimension",["field",60,{"base-type":"type/Text","source-field":40}]]},{"parameter_id":"afa56954","card_id":14,"target":["dimension",["field",41,{"base-type":"type/DateTime"}]]},{"parameter_id":"5eeec658","card_id":14,"target":["dimension",["field",58,{"base-type":"type/Text","source-field":40}]]},{"parameter_id":"512c560a","card_id":14,"target":["dimension",["field",48,{"base-type":"type/Text","source-field":43}]]}]	{"card.description":"Compares the orders of each category quarter over quarter","card.title":"Numbers of orders by category per quarter","column_settings":null,"graph.dimensions":["CREATED_AT","CATEGORY"],"graph.metrics":["count"],"graph.x_axis.labels_enabled":false,"graph.x_axis.title_text":"Created At","graph.y_axis.title_text":"Number of orders","series_settings":{"Doohickey":{"color":"#7172AD"},"Gadget":{"color":"#A989C5"},"Gizmo":{"color":"#C7EAEA"},"Widget":{"color":"#227FD2"}}}	ndNZEBY96LMGGWefgnG3o	\N	2
8	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	7	3	3	17	1	1	[{"parameter_id":"fc2cd1be","card_id":1,"target":["dimension",["field",60,{"base-type":"type/Text","source-field":40}]]},{"parameter_id":"afa56954","card_id":1,"target":["dimension",["field",41,{"base-type":"type/DateTime"}]]},{"parameter_id":"5eeec658","card_id":1,"target":["dimension",["field",58,{"base-type":"type/Text","source-field":40}]]},{"parameter_id":"512c560a","card_id":1,"target":["dimension",["field",48,{"base-type":"type/Text","source-field":43}]]}]	{"card.description":"Compares the total number of orders placed for this product this month with the previous period","column_settings":null}	g6Z9N7QjvtdQbEHaAgxhY	\N	2
9	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	7	3	6	17	5	1	[{"parameter_id":"fc2cd1be","card_id":5,"target":["dimension",["field",60,{"base-type":"type/Text","source-field":40}]]},{"parameter_id":"afa56954","card_id":5,"target":["dimension",["field",41,{"base-type":"type/DateTime"}]]},{"parameter_id":"5eeec658","card_id":5,"target":["dimension",["field",58,{"base-type":"type/Text","source-field":40}]]},{"parameter_id":"512c560a","card_id":5,"target":["dimension",["field",48,{"base-type":"type/Text","source-field":43}]]}]	{"card.description":"Compares the total number of orders placed for this product this month with the previous period","column_settings":null}	Q4sD1F6XrD8dVPD2JHV1Z	\N	2
10	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	7	3	9	17	8	1	[{"parameter_id":"fc2cd1be","card_id":8,"target":["dimension",["field",60,{"base-type":"type/Text","source-field":40}]]},{"parameter_id":"afa56954","card_id":8,"target":["dimension",["field",41,{"base-type":"type/DateTime"}]]},{"parameter_id":"5eeec658","card_id":8,"target":["dimension",["field",58,{"base-type":"type/Text","source-field":40}]]},{"parameter_id":"512c560a","card_id":8,"target":["dimension",["field",48,{"base-type":"type/Text","source-field":43}]]}]	{"card.description":"Compares the total number of orders placed for this product this month with the previous period","column_settings":null}	3d1iG7J48euO467UN-hCL	\N	2
11	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	7	2	1	17	\N	1	[]	{"text":"### Top three all-time products \\nMoM performance","virtual_card":{"archived":false,"dataset_query":{},"display":"text","name":null,"visualization_settings":{}}}	Jtd-7AGoT0MNS4CxhVjei	\N	2
12	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	24	1	0	0	\N	1	[]	{"column_settings":null,"dashcard.background":false,"text":"Who  and where are our customers","virtual_card":{"archived":false,"dataset_query":{},"display":"heading","name":null,"visualization_settings":{}}}	xmP43pM4LHbrbRP-rjT-j	\N	3
13	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	15	6	14	9	17	1	[{"parameter_id":"fc2cd1be","card_id":17,"target":["dimension",["field",60,{"base-type":"type/Text","source-field":40}]]},{"parameter_id":"afa56954","card_id":17,"target":["dimension",["field",41,{"base-type":"type/DateTime"}]]},{"parameter_id":"5eeec658","card_id":17,"target":["dimension",["field",58,{"base-type":"type/Text","source-field":40}]]},{"parameter_id":"512c560a","card_id":17,"target":["dimension",["field",48,{"base-type":"type/Text","join-alias":"People - User"}]]}]	{"card.description":"Shows a distribution of orders broken down by product category across our customers' age groups","card.title":"Product category orders per individual age","column_settings":null,"graph.dimensions":["Age","CATEGORY"],"graph.metrics":["count"],"graph.series_order":null,"graph.series_order_dimension":null,"graph.x_axis.scale":"ordinal","graph.y_axis.title_text":"Orders by category","series_settings":{"Doohickey":{"color":"#7172AD"},"Gadget":{"color":"#A989C5"},"Gizmo":{"color":"#C7EAEA"},"Widget":{"color":"#227FD2"}},"stackable.stack_type":"normalized"}	f-hh5qVJwpgrj3_JWDtpX	\N	3
14	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	15	6	8	9	25	1	[{"parameter_id":"fc2cd1be","card_id":25,"target":["dimension",["field",60,{"base-type":"type/Text","source-field":40}]]},{"parameter_id":"afa56954","card_id":25,"target":["dimension",["field",41,{"base-type":"type/DateTime"}]]},{"parameter_id":"5eeec658","card_id":25,"target":["dimension",["field",58,{"base-type":"type/Text","source-field":40}]]},{"parameter_id":"512c560a","card_id":25,"target":["dimension",["field",48,{"base-type":"type/Text","join-alias":"People - User"}]]}]	{"card.description":"Shows a distribution of revenue in age groups","card.title":"Revenue per individual age","column_settings":null,"graph.dimensions":["Age"],"graph.label_value_formatting":"compact","graph.metrics":["sum"],"graph.series_order":null,"graph.series_order_dimension":null,"graph.show_values":true,"graph.x_axis.scale":"ordinal","graph.y_axis.title_text":"Total revenue","stackable.stack_type":"stacked"}	VjAahFlVkOA0hp1LVBrH5	\N	3
15	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	17	11	1	0	23	1	[{"parameter_id":"afa56954","card_id":23,"target":["dimension",["field","PRODUCTS__via__PRODUCT_ID__TITLE",{"base-type":"type/Text"}]]}]	{"card.description":"An ordered list of our most successful products","card.title":"Best selling products","column_settings":null,"graph.dimensions":["TITLE"],"graph.metrics":["count"],"graph.x_axis.labels_enabled":false,"graph.x_axis.title_text":"Products","graph.y_axis.title_text":"number of orders","series_settings":{"count":{"color":"#999AC4"}},"table.pivot":false}	zfSwuSU-8e6O3eqyt9r16	\N	2
16	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	24	1	12	0	\N	1	[]	{"column_settings":null,"dashcard.background":false,"text":"Product category insights","virtual_card":{"archived":false,"dataset_query":{},"display":"heading","name":null,"visualization_settings":{}}}	tNvUvy6ezg0rUhV9baI0O	\N	2
17	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	17	6	19	0	19	1	[{"parameter_id":"fc2cd1be","card_id":19,"target":["dimension",["field",60,{"base-type":"type/Text","source-field":40}]]},{"parameter_id":"afa56954","card_id":19,"target":["dimension",["field",41,{"base-type":"type/DateTime"}]]},{"parameter_id":"5eeec658","card_id":19,"target":["dimension",["field",58,{"base-type":"type/Text","source-field":40}]]},{"parameter_id":"512c560a","card_id":19,"target":["dimension",["field",48,{"base-type":"type/Text","source-field":43}]]}]	{"card.description":"Shows the distribution of the product categories along the scale of customer ratings","card.title":"Customer Satisfaction per category","column_settings":null,"graph.dimensions":["RATING","CATEGORY"],"graph.metrics":["count"],"graph.series_order":null,"graph.series_order_dimension":null,"graph.x_axis.title_text":"Customer Ratings","graph.y_axis.auto_range":true,"graph.y_axis.title_text":"Orders by category","series_settings":{"Doohickey":{"color":"#7172AD"},"Gadget":{"color":"#A989C5"},"Gizmo":{"color":"#C7EAEA"},"Widget":{"color":"#227FD2"},"count":{"color":"#999AC4"}},"stackable.stack_type":"stacked","table.pivot":false}	OMCqqRP9CQTT3pSp6sfo5	\N	2
18	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	7	6	19	17	22	1	[{"parameter_id":"fc2cd1be","card_id":22,"target":["dimension",["field",60,{"base-type":"type/Text"}]]},{"parameter_id":"5eeec658","card_id":22,"target":["dimension",["field",58,{"base-type":"type/Text"}]]}]	{"card.description":"Indicates the average customer review of our products","card.title":"Average product rating overall","column_settings":null,"gauge.segments":[{"color":"#EF8C8C","label":"awful","max":1,"min":0},{"color":"#F2A86F","label":"bad","max":2,"min":1},{"color":"#F9D45C","label":"alright","max":3,"min":2},{"color":"#A7D07C","label":"good","max":4,"min":3},{"color":"#689636","label":"great","max":5,"min":4}],"table.cell_column":"avg"}	iSUCvbFniuVFlN5NDV3Y7	\N	2
19	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	15	6	20	9	21	1	[{"parameter_id":"fc2cd1be","card_id":21,"target":["dimension",["field",60,{"base-type":"type/Text","source-field":40}]]},{"parameter_id":"afa56954","card_id":21,"target":["dimension",["field",41,{"base-type":"type/DateTime"}]]},{"parameter_id":"5eeec658","card_id":21,"target":["dimension",["field",58,{"base-type":"type/Text","source-field":40}]]},{"parameter_id":"512c560a","card_id":21,"target":["dimension",["field",48,{"base-type":"type/Text","join-alias":"People - User"}]]}]	{"card.description":"Shows a distribution of orders broken down by source across our customers' age groups","card.title":"Orders by source per individual age","column_settings":null,"graph.dimensions":["Age","SOURCE"],"graph.metrics":["count"],"graph.series_order":null,"graph.series_order_dimension":null,"graph.x_axis.axis_enabled":true,"graph.x_axis.scale":"ordinal","graph.y_axis.title_text":"Orders by Sources","stackable.stack_type":"normalized"}	s_DDqee44Wp_gHh16bAuR	\N	3
20	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	24	1	7	0	\N	1	[]	{"column_settings":null,"dashcard.background":false,"text":"Age insights breakdowns","virtual_card":{"archived":false,"dataset_query":{},"display":"heading","name":null,"visualization_settings":{}}}	P37ixkaPafIr47ep0oypD	\N	3
21	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	18	5	4	6	12	1	[{"parameter_id":"fc2cd1be","card_id":12,"target":["dimension",["field",60,{"base-type":"type/Text","source-field":40}]]},{"parameter_id":"afa56954","card_id":12,"target":["dimension",["field",41,{"base-type":"type/DateTime"}]]},{"parameter_id":"5eeec658","card_id":12,"target":["dimension",["field",58,{"base-type":"type/Text","source-field":40}]]},{"parameter_id":"512c560a","card_id":12,"target":["dimension",["field",48,{"base-type":"type/Text","source-field":43}]]}]	{"card.description":"Matches the cumulative revenue  month over month with the number of orders placed each month","column_settings":null,"graph.dimensions":["CREATED_AT"],"graph.metrics":["sum","sum_2"],"graph.show_trendline":false,"graph.x_axis.labels_enabled":false,"graph.x_axis.title_text":"Orders date","graph.y_axis.labels_enabled":false,"graph.y_axis.title_text":"Revenue","series_settings":{"sum":{"display":"line","line.interpolate":"linear","line.marker_enabled":false,"show_series_values":true,"title":"Revenue"},"sum_2":{"color":"#51528D","title":"Number of orders"}}}	Cc0jFy1OZZbmC1CEA4eu2	\N	1
22	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	24	1	18	0	\N	1	[]	{"column_settings":null,"dashcard.background":false,"text":"Customer satisfaction insights","virtual_card":{"archived":false,"dataset_query":{},"display":"heading","name":null,"visualization_settings":{}}}	wRqWXKjFJ_d3t34HXl7Vg	\N	2
23	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	6	3	1	12	20	1	[{"parameter_id":"fc2cd1be","card_id":20,"target":["dimension",["field",60,{"base-type":"type/Text","source-field":40}]]},{"parameter_id":"afa56954","card_id":20,"target":["dimension",["field",41,{"base-type":"type/DateTime"}]]},{"parameter_id":"5eeec658","card_id":20,"target":["dimension",["field",58,{"base-type":"type/Text","source-field":40}]]},{"parameter_id":"512c560a","card_id":20,"target":["dimension",["field",48,{"base-type":"type/Text","source-field":43}]]}]	{"card.description":"Compares the total number of orders placed this month with the previous period","column_settings":null}	xZFChb3ro-2jzg7I767yY	\N	1
24	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	6	3	1	6	18	1	[{"parameter_id":"fc2cd1be","card_id":18,"target":["dimension",["field",60,{"base-type":"type/Text","source-field":40}]]},{"parameter_id":"afa56954","card_id":18,"target":["dimension",["field",41,{"base-type":"type/DateTime"}]]},{"parameter_id":"5eeec658","card_id":18,"target":["dimension",["field",58,{"base-type":"type/Text","source-field":40}]]},{"parameter_id":"512c560a","card_id":18,"target":["dimension",["field",48,{"base-type":"type/Text","source-field":43}]]}]	{"card.description":"Compares the total discount given this quarter with the previous period","card.title":"Discount given per quarter","column_settings":null}	NGar3mVnCROw5p7djSBa3	\N	1
25	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	6	3	1	0	15	1	[{"parameter_id":"fc2cd1be","card_id":15,"target":["dimension",["field",60,{"base-type":"type/Text","source-field":40}]]},{"parameter_id":"afa56954","card_id":15,"target":["dimension",["field",41,{"base-type":"type/DateTime"}]]},{"parameter_id":"5eeec658","card_id":15,"target":["dimension",["field",58,{"base-type":"type/Text","source-field":40}]]},{"parameter_id":"512c560a","card_id":15,"target":["dimension",["field",48,{"base-type":"type/Text","source-field":43}]]}]	{"card.description":"Compares the total revenue this quarter with the previous period","column_settings":null}	dX_lc-n9uNhVcRmbhZfdQ	\N	1
26	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	6	3	1	18	24	1	[{"parameter_id":"512c560a","card_id":24,"target":["dimension",["field",48,{"base-type":"type/Text"}]]}]	{"card.description":"Compares the number of unique customers placing orders this month with the previous period","card.title":"Unique customers per month","column_settings":null}	1mSLJrWPvehXkcIzoK1Yq	\N	1
27	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	6	5	4	0	13	1	[{"parameter_id":"fc2cd1be","card_id":13,"target":["dimension",["field",60,{"base-type":"type/Text","source-field":40}]]},{"parameter_id":"afa56954","card_id":13,"target":["dimension",["field",41,{"base-type":"type/DateTime"}]]},{"parameter_id":"5eeec658","card_id":13,"target":["dimension",["field",58,{"base-type":"type/Text","source-field":40}]]},{"parameter_id":"512c560a","card_id":13,"target":["dimension",["field",48,{"base-type":"type/Text","source-field":43}]]}]	{"card.description":"Compares the current total in revenue with the set goal for the quarter","progress.goal":250000,"column_settings":{"[\\"name\\",\\"sum\\"]":{"number_style":"currency"}}}	e-8h-dZx-B9RYybRoZzci	\N	1
28	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	9	6	1	0	9	1	[{"parameter_id":"fc2cd1be","card_id":9,"target":["dimension",["field",60,{"base-type":"type/Text","source-field":40}]]},{"parameter_id":"afa56954","card_id":9,"target":["dimension",["field",41,{"base-type":"type/DateTime"}]]},{"parameter_id":"5eeec658","card_id":9,"target":["dimension",["field",58,{"base-type":"type/Text","source-field":40}]]},{"parameter_id":"512c560a","card_id":9,"target":["dimension",["field",48,{"base-type":"type/Text","source-field":43}]]}]	{"column_settings":null}	M89Lh2AQnAunejYWXNUeT	\N	3
29	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	9	6	14	0	11	1	[{"parameter_id":"fc2cd1be","card_id":11,"target":["dimension",["field",60,{"base-type":"type/Text","join-alias":"Products"}]]},{"parameter_id":"afa56954","card_id":11,"target":["dimension",["field","CREATED_AT",{"base-type":"type/DateTime"}]]},{"parameter_id":"5eeec658","card_id":11,"target":["dimension",["field",58,{"base-type":"type/Text","join-alias":"Products"}]]},{"parameter_id":"512c560a","card_id":11,"target":["dimension",["field",48,{"base-type":"type/Text","source-field":43}]]}]	{"graph.show_values":true,"graph.x_axis.labels_enabled":true,"table.cell_column":"SUBTOTAL","graph.series_order_dimension":null,"graph.x_axis.axis_enabled":true,"graph.y_axis.labels_enabled":false,"card.title":"Product category orders per age group","graph.metrics":["count"],"graph.series_order":null,"series_settings":{"Doohickey":{"color":"#7172AD"},"Gadget":{"color":"#A989C5"},"Gizmo":{"color":"#C7EAEA"},"Widget":{"color":"#227FD2"}},"graph.y_axis.auto_range":true,"graph.x_axis.scale":"ordinal","graph.dimensions":["Age","CATEGORY"],"stackable.stack_type":"stacked"}	3f4Uo5mlklic7PBxY32DZ	\N	3
30	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	9	6	8	0	3	1	[{"card_id":3,"parameter_id":"afa56954","target":["dimension",["field",41,{"base-type":"type/DateTime"}]]},{"card_id":3,"parameter_id":"512c560a","target":["dimension",["field",48,{"base-type":"type/Text","source-field":43}]]},{"card_id":3,"parameter_id":"5eeec658","target":["dimension",["field",58,{"base-type":"type/Text","source-field":40}]]},{"card_id":3,"parameter_id":"fc2cd1be","target":["dimension",["field",60,{"base-type":"type/Text","source-field":40}]]}]	{"column_settings":null,"graph.dimensions":["Age"],"graph.metrics":["sum"],"graph.series_order":null,"graph.series_order_dimension":null,"graph.show_values":true,"graph.x_axis.scale":"ordinal","stackable.stack_type":"stacked","table.cell_column":"SUBTOTAL"}	h_f6hPOFg0FN3m6VA5B03	\N	3
31	2024-06-22 08:19:17.839482+00	2024-06-22 08:19:17.839482+00	9	6	20	0	10	1	[{"parameter_id":"fc2cd1be","card_id":10,"target":["dimension",["field",60,{"base-type":"type/Text","source-field":40}]]},{"parameter_id":"afa56954","card_id":10,"target":["dimension",["field","CREATED_AT",{"base-type":"type/DateTime"}]]},{"parameter_id":"5eeec658","card_id":10,"target":["dimension",["field",58,{"base-type":"type/Text","source-field":40}]]},{"parameter_id":"512c560a","card_id":10,"target":["dimension",["field",48,{"base-type":"type/Text","join-alias":"People - User"}]]}]	{"column_settings":null,"graph.dimensions":["Age","SOURCE"],"graph.label_value_formatting":"compact","graph.metrics":["count"],"graph.series_order":null,"graph.series_order_dimension":null,"graph.show_values":true,"graph.x_axis.axis_enabled":true,"graph.x_axis.labels_enabled":true,"graph.x_axis.scale":"ordinal","graph.y_axis.auto_range":true,"graph.y_axis.labels_enabled":false,"graph.y_axis.title_text":"Orders","stackable.stack_type":"stacked","table.cell_column":"SUBTOTAL"}	JDDgMLIDoRDC24ufF2yeI	\N	3
\.


--
-- Data for Name: revision; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.revision (id, model, model_id, user_id, "timestamp", object, is_reversion, is_creation, message, most_recent, metabase_version) FROM stdin;
\.


--
-- Data for Name: sandboxes; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.sandboxes (id, group_id, table_id, card_id, attribute_remappings, permission_id) FROM stdin;
\.


--
-- Data for Name: secret; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.secret (id, version, creator_id, created_at, updated_at, name, kind, source, value) FROM stdin;
\.


--
-- Data for Name: segment; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.segment (id, table_id, creator_id, name, description, archived, definition, created_at, updated_at, points_of_interest, caveats, show_in_getting_started, entity_id) FROM stdin;
\.


--
-- Data for Name: setting; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.setting (key, value) FROM stdin;
example-dashboard-id	1
setup-token	49a99579-1bf2-4ba1-a4c7-1ff5b02e8215
startup-time-millis	52985
redirect-all-requests-to-https	false
site-url	http://metabase.docker
analytics-uuid	0d7804dd-92e8-405d-ac43-6529cc382835
instance-creation	2024-06-22T08:19:16.518337Z
site-name	Autonomy
admin-email	aman@punjab.com
site-locale	en
anon-tracking-enabled	true
site-uuid	aec085ee-a8cb-49a8-b344-823b98a83fee
setup-license-active-at-setup	false
settings-last-updated	2024-06-22 08:25:38.0763+00
\.


--
-- Data for Name: table_privileges; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.table_privileges (table_id, role, "select", update, insert, delete) FROM stdin;
\.


--
-- Data for Name: task_history; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.task_history (id, task, db_id, started_at, ended_at, duration, task_details, status) FROM stdin;
2	sync-dbms-version	1	2024-06-22 08:20:23.490725+00	2024-06-22 08:20:23.619041+00	128	{"flavor":"H2","version":"2.1.214 (2022-06-13)","semantic-version":[2,1]}	success
3	sync-timezone	1	2024-06-22 08:20:23.697369+00	2024-06-22 08:20:24.114802+00	417	{"timezone-id":"GMT"}	success
4	sync-tables	1	2024-06-22 08:20:24.205398+00	2024-06-22 08:20:24.40296+00	197	{"updated-tables":0,"total-tables":8}	success
5	sync-fields	1	2024-06-22 08:20:24.447377+00	2024-06-22 08:20:25.766162+00	1318	{"total-fields":71,"updated-fields":0}	success
6	sync-fks	1	2024-06-22 08:20:25.89782+00	2024-06-22 08:20:26.681006+00	783	{"total-fks":6,"updated-fks":0,"total-failed":0}	success
7	sync-indexes	1	2024-06-22 08:20:26.756721+00	2024-06-22 08:20:27.855348+00	1098	{"total-indexes":0,"added-indexes":0,"removed-indexes":0}	success
8	sync-metabase-metadata	1	2024-06-22 08:20:27.864119+00	2024-06-22 08:20:30.086226+00	2222	{}	success
9	sync-table-privileges	1	2024-06-22 08:20:30.158947+00	2024-06-22 08:20:30.174571+00	15	\N	success
1	sync	1	2024-06-22 08:20:23.11242+00	2024-06-22 08:20:30.187882+00	7098	\N	success
11	fingerprint-fields	1	2024-06-22 08:20:30.279263+00	2024-06-22 08:20:30.452401+00	173	{"no-data-fingerprints":0,"failed-fingerprints":0,"updated-fingerprints":0,"fingerprints-attempted":0}	success
12	classify-fields	1	2024-06-22 08:20:30.464675+00	2024-06-22 08:20:30.490497+00	25	{"fields-classified":0,"fields-failed":0}	success
13	classify-tables	1	2024-06-22 08:20:30.499564+00	2024-06-22 08:20:30.569661+00	70	{"total-tables":8,"tables-classified":0}	success
10	analyze	1	2024-06-22 08:20:30.270304+00	2024-06-22 08:20:30.579304+00	309	\N	success
14	field values scanning	1	2024-06-22 08:20:30.598137+00	\N	\N	\N	started
15	delete-expired-advanced-field-values	1	2024-06-22 08:20:30.609143+00	2024-06-22 08:20:31.455157+00	846	{"deleted":0}	success
16	update-field-values	1	2024-06-22 08:20:31.464435+00	\N	\N	\N	started
18	sync-dbms-version	1	2024-06-22 08:23:36.98025+00	2024-06-22 08:23:37.013151+00	32	{"flavor":"H2","version":"2.1.214 (2022-06-13)","semantic-version":[2,1]}	success
19	sync-timezone	1	2024-06-22 08:23:37.029325+00	2024-06-22 08:23:37.12179+00	92	{"timezone-id":"GMT"}	success
20	sync-tables	1	2024-06-22 08:23:37.132871+00	2024-06-22 08:23:37.280072+00	147	{"updated-tables":0,"total-tables":8}	success
21	sync-fields	1	2024-06-22 08:23:37.322115+00	2024-06-22 08:23:38.219061+00	896	{"total-fields":71,"updated-fields":0}	success
22	sync-fks	1	2024-06-22 08:23:38.232374+00	2024-06-22 08:23:38.422041+00	189	{"total-fks":6,"updated-fks":0,"total-failed":0}	success
23	sync-indexes	1	2024-06-22 08:23:38.478259+00	2024-06-22 08:23:39.221868+00	743	{"total-indexes":0,"added-indexes":0,"removed-indexes":0}	success
24	sync-metabase-metadata	1	2024-06-22 08:23:39.286989+00	2024-06-22 08:23:41.420929+00	2133	{}	success
25	sync-table-privileges	1	2024-06-22 08:23:41.430376+00	2024-06-22 08:23:41.441404+00	11	\N	success
17	sync	1	2024-06-22 08:23:36.822614+00	2024-06-22 08:23:41.486919+00	4664	\N	success
27	fingerprint-fields	1	2024-06-22 08:23:41.53579+00	2024-06-22 08:23:41.713517+00	177	{"no-data-fingerprints":0,"failed-fingerprints":0,"updated-fingerprints":0,"fingerprints-attempted":0}	success
28	classify-fields	1	2024-06-22 08:23:41.725176+00	2024-06-22 08:23:41.814968+00	89	{"fields-classified":0,"fields-failed":0}	success
29	classify-tables	1	2024-06-22 08:23:41.825383+00	2024-06-22 08:23:41.83993+00	14	{"total-tables":8,"tables-classified":0}	success
26	analyze	1	2024-06-22 08:23:41.526321+00	2024-06-22 08:23:41.852421+00	326	\N	success
31	delete-expired-advanced-field-values	1	2024-06-22 08:23:41.939511+00	2024-06-22 08:23:42.83032+00	890	{"deleted":0}	success
32	update-field-values	1	2024-06-22 08:23:42.885359+00	2024-06-22 08:23:46.280981+00	3395	{"errors":0,"created":13,"updated":0,"deleted":0}	success
30	field values scanning	1	2024-06-22 08:23:41.930125+00	2024-06-22 08:23:46.321681+00	4391	\N	success
34	sync-dbms-version	2	2024-06-22 08:25:36.534036+00	2024-06-22 08:25:36.634247+00	100	{"flavor":"PostgreSQL","version":"16.3","semantic-version":[16,3]}	success
35	sync-timezone	2	2024-06-22 08:25:36.639417+00	2024-06-22 08:25:36.729058+00	89	{"timezone-id":"GMT"}	success
36	sync-tables	2	2024-06-22 08:25:36.733305+00	2024-06-22 08:25:36.741251+00	7	{"updated-tables":0,"total-tables":0}	success
37	sync-fields	2	2024-06-22 08:25:36.815573+00	2024-06-22 08:25:36.821049+00	5	{"total-fields":0,"updated-fields":0}	success
38	sync-fks	2	2024-06-22 08:25:36.824981+00	2024-06-22 08:25:36.830761+00	5	{"total-fks":0,"updated-fks":0,"total-failed":0}	success
39	sync-indexes	2	2024-06-22 08:25:36.839805+00	2024-06-22 08:25:36.849368+00	9	{"total-indexes":0,"added-indexes":0,"removed-indexes":0}	success
40	sync-metabase-metadata	2	2024-06-22 08:25:36.915504+00	2024-06-22 08:25:36.920104+00	4	{}	success
41	sync-table-privileges	2	2024-06-22 08:25:36.92749+00	2024-06-22 08:25:36.944337+00	16	{"total-table-privileges":0}	success
33	sync	2	2024-06-22 08:25:36.527011+00	2024-06-22 08:25:36.948213+00	421	\N	success
43	fingerprint-fields	2	2024-06-22 08:25:37.028802+00	2024-06-22 08:25:37.034489+00	5	{"no-data-fingerprints":0,"failed-fingerprints":0,"updated-fingerprints":0,"fingerprints-attempted":0}	success
44	classify-fields	2	2024-06-22 08:25:37.038553+00	2024-06-22 08:25:37.042276+00	3	{"fields-classified":0,"fields-failed":0}	success
45	classify-tables	2	2024-06-22 08:25:37.11694+00	2024-06-22 08:25:37.120488+00	3	{"total-tables":0,"tables-classified":0}	success
42	analyze	2	2024-06-22 08:25:37.024963+00	2024-06-22 08:25:37.124389+00	99	\N	success
47	delete-expired-advanced-field-values	2	2024-06-22 08:25:37.133955+00	2024-06-22 08:25:37.139986+00	6	{"deleted":0}	success
48	update-field-values	2	2024-06-22 08:25:37.145657+00	2024-06-22 08:25:37.150888+00	5	\N	success
46	field values scanning	2	2024-06-22 08:25:37.129909+00	2024-06-22 08:25:37.218814+00	88	\N	success
\.


--
-- Data for Name: timeline; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.timeline (id, name, description, icon, collection_id, archived, creator_id, created_at, updated_at, "default", entity_id) FROM stdin;
\.


--
-- Data for Name: timeline_event; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.timeline_event (id, timeline_id, name, description, "timestamp", time_matters, timezone, icon, archived, creator_id, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: user_parameter_value; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.user_parameter_value (id, user_id, parameter_id, value) FROM stdin;
\.


--
-- Data for Name: view_log; Type: TABLE DATA; Schema: public; Owner: aman
--

COPY public.view_log (id, user_id, model, model_id, "timestamp", metadata, has_access, context) FROM stdin;
\.


--
-- Name: action_id_seq; Type: SEQUENCE SET; Schema: public; Owner: aman
--

SELECT pg_catalog.setval('public.action_id_seq', 1, false);


--
-- Name: api_key_id_seq; Type: SEQUENCE SET; Schema: public; Owner: aman
--

SELECT pg_catalog.setval('public.api_key_id_seq', 1, false);


--
-- Name: application_permissions_revision_id_seq; Type: SEQUENCE SET; Schema: public; Owner: aman
--

SELECT pg_catalog.setval('public.application_permissions_revision_id_seq', 1, false);


--
-- Name: audit_log_id_seq; Type: SEQUENCE SET; Schema: public; Owner: aman
--

SELECT pg_catalog.setval('public.audit_log_id_seq', 1, false);


--
-- Name: bookmark_ordering_id_seq; Type: SEQUENCE SET; Schema: public; Owner: aman
--

SELECT pg_catalog.setval('public.bookmark_ordering_id_seq', 1, false);


--
-- Name: cache_config_id_seq; Type: SEQUENCE SET; Schema: public; Owner: aman
--

SELECT pg_catalog.setval('public.cache_config_id_seq', 1, false);


--
-- Name: card_bookmark_id_seq; Type: SEQUENCE SET; Schema: public; Owner: aman
--

SELECT pg_catalog.setval('public.card_bookmark_id_seq', 1, false);


--
-- Name: card_label_id_seq; Type: SEQUENCE SET; Schema: public; Owner: aman
--

SELECT pg_catalog.setval('public.card_label_id_seq', 1, false);


--
-- Name: cloud_migration_id_seq; Type: SEQUENCE SET; Schema: public; Owner: aman
--

SELECT pg_catalog.setval('public.cloud_migration_id_seq', 1, false);


--
-- Name: collection_bookmark_id_seq; Type: SEQUENCE SET; Schema: public; Owner: aman
--

SELECT pg_catalog.setval('public.collection_bookmark_id_seq', 1, false);


--
-- Name: collection_id_seq; Type: SEQUENCE SET; Schema: public; Owner: aman
--

SELECT pg_catalog.setval('public.collection_id_seq', 2, true);


--
-- Name: collection_permission_graph_revision_id_seq; Type: SEQUENCE SET; Schema: public; Owner: aman
--

SELECT pg_catalog.setval('public.collection_permission_graph_revision_id_seq', 1, false);


--
-- Name: connection_impersonations_id_seq; Type: SEQUENCE SET; Schema: public; Owner: aman
--

SELECT pg_catalog.setval('public.connection_impersonations_id_seq', 1, false);


--
-- Name: core_user_id_seq; Type: SEQUENCE SET; Schema: public; Owner: aman
--

SELECT pg_catalog.setval('public.core_user_id_seq', 1, true);


--
-- Name: dashboard_bookmark_id_seq; Type: SEQUENCE SET; Schema: public; Owner: aman
--

SELECT pg_catalog.setval('public.dashboard_bookmark_id_seq', 1, false);


--
-- Name: dashboard_favorite_id_seq; Type: SEQUENCE SET; Schema: public; Owner: aman
--

SELECT pg_catalog.setval('public.dashboard_favorite_id_seq', 1, false);


--
-- Name: dashboard_tab_id_seq; Type: SEQUENCE SET; Schema: public; Owner: aman
--

SELECT pg_catalog.setval('public.dashboard_tab_id_seq', 3, true);


--
-- Name: dashboardcard_series_id_seq; Type: SEQUENCE SET; Schema: public; Owner: aman
--

SELECT pg_catalog.setval('public.dashboardcard_series_id_seq', 1, false);


--
-- Name: data_permissions_id_seq; Type: SEQUENCE SET; Schema: public; Owner: aman
--

SELECT pg_catalog.setval('public.data_permissions_id_seq', 11, true);


--
-- Name: dependency_id_seq; Type: SEQUENCE SET; Schema: public; Owner: aman
--

SELECT pg_catalog.setval('public.dependency_id_seq', 1, false);


--
-- Name: dimension_id_seq; Type: SEQUENCE SET; Schema: public; Owner: aman
--

SELECT pg_catalog.setval('public.dimension_id_seq', 1, false);


--
-- Name: field_usage_id_seq; Type: SEQUENCE SET; Schema: public; Owner: aman
--

SELECT pg_catalog.setval('public.field_usage_id_seq', 1, false);


--
-- Name: group_table_access_policy_id_seq; Type: SEQUENCE SET; Schema: public; Owner: aman
--

SELECT pg_catalog.setval('public.group_table_access_policy_id_seq', 1, false);


--
-- Name: label_id_seq; Type: SEQUENCE SET; Schema: public; Owner: aman
--

SELECT pg_catalog.setval('public.label_id_seq', 1, false);


--
-- Name: login_history_id_seq; Type: SEQUENCE SET; Schema: public; Owner: aman
--

SELECT pg_catalog.setval('public.login_history_id_seq', 1, true);


--
-- Name: metabase_database_id_seq; Type: SEQUENCE SET; Schema: public; Owner: aman
--

SELECT pg_catalog.setval('public.metabase_database_id_seq', 2, true);


--
-- Name: metabase_field_id_seq; Type: SEQUENCE SET; Schema: public; Owner: aman
--

SELECT pg_catalog.setval('public.metabase_field_id_seq', 71, true);


--
-- Name: metabase_fieldvalues_id_seq; Type: SEQUENCE SET; Schema: public; Owner: aman
--

SELECT pg_catalog.setval('public.metabase_fieldvalues_id_seq', 24, true);


--
-- Name: metabase_table_id_seq; Type: SEQUENCE SET; Schema: public; Owner: aman
--

SELECT pg_catalog.setval('public.metabase_table_id_seq', 8, true);


--
-- Name: metric_id_seq; Type: SEQUENCE SET; Schema: public; Owner: aman
--

SELECT pg_catalog.setval('public.metric_id_seq', 1, false);


--
-- Name: metric_important_field_id_seq; Type: SEQUENCE SET; Schema: public; Owner: aman
--

SELECT pg_catalog.setval('public.metric_important_field_id_seq', 1, false);


--
-- Name: model_index_id_seq; Type: SEQUENCE SET; Schema: public; Owner: aman
--

SELECT pg_catalog.setval('public.model_index_id_seq', 1, false);


--
-- Name: moderation_review_id_seq; Type: SEQUENCE SET; Schema: public; Owner: aman
--

SELECT pg_catalog.setval('public.moderation_review_id_seq', 1, false);


--
-- Name: native_query_snippet_id_seq; Type: SEQUENCE SET; Schema: public; Owner: aman
--

SELECT pg_catalog.setval('public.native_query_snippet_id_seq', 1, false);


--
-- Name: parameter_card_id_seq; Type: SEQUENCE SET; Schema: public; Owner: aman
--

SELECT pg_catalog.setval('public.parameter_card_id_seq', 1, false);


--
-- Name: permissions_group_id_seq; Type: SEQUENCE SET; Schema: public; Owner: aman
--

SELECT pg_catalog.setval('public.permissions_group_id_seq', 2, true);


--
-- Name: permissions_group_membership_id_seq; Type: SEQUENCE SET; Schema: public; Owner: aman
--

SELECT pg_catalog.setval('public.permissions_group_membership_id_seq', 3, true);


--
-- Name: permissions_id_seq; Type: SEQUENCE SET; Schema: public; Owner: aman
--

SELECT pg_catalog.setval('public.permissions_id_seq', 5, true);


--
-- Name: permissions_revision_id_seq; Type: SEQUENCE SET; Schema: public; Owner: aman
--

SELECT pg_catalog.setval('public.permissions_revision_id_seq', 1, false);


--
-- Name: persisted_info_id_seq; Type: SEQUENCE SET; Schema: public; Owner: aman
--

SELECT pg_catalog.setval('public.persisted_info_id_seq', 1, false);


--
-- Name: pulse_card_id_seq; Type: SEQUENCE SET; Schema: public; Owner: aman
--

SELECT pg_catalog.setval('public.pulse_card_id_seq', 1, false);


--
-- Name: pulse_channel_id_seq; Type: SEQUENCE SET; Schema: public; Owner: aman
--

SELECT pg_catalog.setval('public.pulse_channel_id_seq', 1, false);


--
-- Name: pulse_channel_recipient_id_seq; Type: SEQUENCE SET; Schema: public; Owner: aman
--

SELECT pg_catalog.setval('public.pulse_channel_recipient_id_seq', 1, false);


--
-- Name: pulse_id_seq; Type: SEQUENCE SET; Schema: public; Owner: aman
--

SELECT pg_catalog.setval('public.pulse_id_seq', 1, false);


--
-- Name: query_execution_id_seq; Type: SEQUENCE SET; Schema: public; Owner: aman
--

SELECT pg_catalog.setval('public.query_execution_id_seq', 1, false);


--
-- Name: query_field_id_seq; Type: SEQUENCE SET; Schema: public; Owner: aman
--

SELECT pg_catalog.setval('public.query_field_id_seq', 1, false);


--
-- Name: recent_views_id_seq; Type: SEQUENCE SET; Schema: public; Owner: aman
--

SELECT pg_catalog.setval('public.recent_views_id_seq', 1, false);


--
-- Name: report_card_id_seq; Type: SEQUENCE SET; Schema: public; Owner: aman
--

SELECT pg_catalog.setval('public.report_card_id_seq', 26, true);


--
-- Name: report_cardfavorite_id_seq; Type: SEQUENCE SET; Schema: public; Owner: aman
--

SELECT pg_catalog.setval('public.report_cardfavorite_id_seq', 1, false);


--
-- Name: report_dashboard_id_seq; Type: SEQUENCE SET; Schema: public; Owner: aman
--

SELECT pg_catalog.setval('public.report_dashboard_id_seq', 1, true);


--
-- Name: report_dashboardcard_id_seq; Type: SEQUENCE SET; Schema: public; Owner: aman
--

SELECT pg_catalog.setval('public.report_dashboardcard_id_seq', 31, true);


--
-- Name: revision_id_seq; Type: SEQUENCE SET; Schema: public; Owner: aman
--

SELECT pg_catalog.setval('public.revision_id_seq', 1, false);


--
-- Name: secret_id_seq; Type: SEQUENCE SET; Schema: public; Owner: aman
--

SELECT pg_catalog.setval('public.secret_id_seq', 1, false);


--
-- Name: segment_id_seq; Type: SEQUENCE SET; Schema: public; Owner: aman
--

SELECT pg_catalog.setval('public.segment_id_seq', 1, false);


--
-- Name: task_history_id_seq; Type: SEQUENCE SET; Schema: public; Owner: aman
--

SELECT pg_catalog.setval('public.task_history_id_seq', 48, true);


--
-- Name: timeline_event_id_seq; Type: SEQUENCE SET; Schema: public; Owner: aman
--

SELECT pg_catalog.setval('public.timeline_event_id_seq', 1, false);


--
-- Name: timeline_id_seq; Type: SEQUENCE SET; Schema: public; Owner: aman
--

SELECT pg_catalog.setval('public.timeline_id_seq', 1, false);


--
-- Name: user_parameter_value_id_seq; Type: SEQUENCE SET; Schema: public; Owner: aman
--

SELECT pg_catalog.setval('public.user_parameter_value_id_seq', 1, false);


--
-- Name: view_log_id_seq; Type: SEQUENCE SET; Schema: public; Owner: aman
--

SELECT pg_catalog.setval('public.view_log_id_seq', 1, false);


--
-- Name: action action_entity_id_key; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.action
    ADD CONSTRAINT action_entity_id_key UNIQUE (entity_id);


--
-- Name: action action_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.action
    ADD CONSTRAINT action_pkey PRIMARY KEY (id);


--
-- Name: action action_public_uuid_key; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.action
    ADD CONSTRAINT action_public_uuid_key UNIQUE (public_uuid);


--
-- Name: api_key api_key_key_prefix_key; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.api_key
    ADD CONSTRAINT api_key_key_prefix_key UNIQUE (key_prefix);


--
-- Name: api_key api_key_name_key; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.api_key
    ADD CONSTRAINT api_key_name_key UNIQUE (name);


--
-- Name: api_key api_key_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.api_key
    ADD CONSTRAINT api_key_pkey PRIMARY KEY (id);


--
-- Name: audit_log audit_log_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.audit_log
    ADD CONSTRAINT audit_log_pkey PRIMARY KEY (id);


--
-- Name: bookmark_ordering bookmark_ordering_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.bookmark_ordering
    ADD CONSTRAINT bookmark_ordering_pkey PRIMARY KEY (id);


--
-- Name: cache_config cache_config_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.cache_config
    ADD CONSTRAINT cache_config_pkey PRIMARY KEY (id);


--
-- Name: card_bookmark card_bookmark_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.card_bookmark
    ADD CONSTRAINT card_bookmark_pkey PRIMARY KEY (id);


--
-- Name: card_label card_label_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.card_label
    ADD CONSTRAINT card_label_pkey PRIMARY KEY (id);


--
-- Name: cloud_migration cloud_migration_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.cloud_migration
    ADD CONSTRAINT cloud_migration_pkey PRIMARY KEY (id);


--
-- Name: collection_bookmark collection_bookmark_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.collection_bookmark
    ADD CONSTRAINT collection_bookmark_pkey PRIMARY KEY (id);


--
-- Name: collection collection_entity_id_key; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.collection
    ADD CONSTRAINT collection_entity_id_key UNIQUE (entity_id);


--
-- Name: collection collection_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.collection
    ADD CONSTRAINT collection_pkey PRIMARY KEY (id);


--
-- Name: collection_permission_graph_revision collection_revision_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.collection_permission_graph_revision
    ADD CONSTRAINT collection_revision_pkey PRIMARY KEY (id);


--
-- Name: connection_impersonations conn_impersonation_unique_group_id_db_id; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.connection_impersonations
    ADD CONSTRAINT conn_impersonation_unique_group_id_db_id UNIQUE (group_id, db_id);


--
-- Name: connection_impersonations connection_impersonations_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.connection_impersonations
    ADD CONSTRAINT connection_impersonations_pkey PRIMARY KEY (id);


--
-- Name: core_session core_session_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.core_session
    ADD CONSTRAINT core_session_pkey PRIMARY KEY (id);


--
-- Name: core_user core_user_email_key; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.core_user
    ADD CONSTRAINT core_user_email_key UNIQUE (email);


--
-- Name: core_user core_user_entity_id_key; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.core_user
    ADD CONSTRAINT core_user_entity_id_key UNIQUE (entity_id);


--
-- Name: core_user core_user_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.core_user
    ADD CONSTRAINT core_user_pkey PRIMARY KEY (id);


--
-- Name: dashboard_bookmark dashboard_bookmark_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.dashboard_bookmark
    ADD CONSTRAINT dashboard_bookmark_pkey PRIMARY KEY (id);


--
-- Name: dashboard_favorite dashboard_favorite_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.dashboard_favorite
    ADD CONSTRAINT dashboard_favorite_pkey PRIMARY KEY (id);


--
-- Name: dashboard_tab dashboard_tab_entity_id_key; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.dashboard_tab
    ADD CONSTRAINT dashboard_tab_entity_id_key UNIQUE (entity_id);


--
-- Name: dashboard_tab dashboard_tab_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.dashboard_tab
    ADD CONSTRAINT dashboard_tab_pkey PRIMARY KEY (id);


--
-- Name: dashboardcard_series dashboardcard_series_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.dashboardcard_series
    ADD CONSTRAINT dashboardcard_series_pkey PRIMARY KEY (id);


--
-- Name: data_permissions data_permissions_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.data_permissions
    ADD CONSTRAINT data_permissions_pkey PRIMARY KEY (id);


--
-- Name: databasechangeloglock databasechangeloglock_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.databasechangeloglock
    ADD CONSTRAINT databasechangeloglock_pkey PRIMARY KEY (id);


--
-- Name: dependency dependency_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.dependency
    ADD CONSTRAINT dependency_pkey PRIMARY KEY (id);


--
-- Name: dimension dimension_entity_id_key; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.dimension
    ADD CONSTRAINT dimension_entity_id_key UNIQUE (entity_id);


--
-- Name: dimension dimension_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.dimension
    ADD CONSTRAINT dimension_pkey PRIMARY KEY (id);


--
-- Name: field_usage field_usage_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.field_usage
    ADD CONSTRAINT field_usage_pkey PRIMARY KEY (id);


--
-- Name: application_permissions_revision general_permissions_revision_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.application_permissions_revision
    ADD CONSTRAINT general_permissions_revision_pkey PRIMARY KEY (id);


--
-- Name: sandboxes group_table_access_policy_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.sandboxes
    ADD CONSTRAINT group_table_access_policy_pkey PRIMARY KEY (id);


--
-- Name: cache_config idx_cache_config_unique_model; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.cache_config
    ADD CONSTRAINT idx_cache_config_unique_model UNIQUE (model, model_id);


--
-- Name: databasechangelog idx_databasechangelog_id_author_filename; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.databasechangelog
    ADD CONSTRAINT idx_databasechangelog_id_author_filename UNIQUE (id, author, filename);


--
-- Name: metabase_field idx_uniq_field_table_id_parent_id_name; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.metabase_field
    ADD CONSTRAINT idx_uniq_field_table_id_parent_id_name UNIQUE (table_id, parent_id, name);


--
-- Name: metabase_table idx_uniq_table_db_id_schema_name; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.metabase_table
    ADD CONSTRAINT idx_uniq_table_db_id_schema_name UNIQUE (db_id, schema, name);


--
-- Name: report_cardfavorite idx_unique_cardfavorite_card_id_owner_id; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.report_cardfavorite
    ADD CONSTRAINT idx_unique_cardfavorite_card_id_owner_id UNIQUE (card_id, owner_id);


--
-- Name: label label_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.label
    ADD CONSTRAINT label_pkey PRIMARY KEY (id);


--
-- Name: label label_slug_key; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.label
    ADD CONSTRAINT label_slug_key UNIQUE (slug);


--
-- Name: login_history login_history_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.login_history
    ADD CONSTRAINT login_history_pkey PRIMARY KEY (id);


--
-- Name: metabase_database metabase_database_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.metabase_database
    ADD CONSTRAINT metabase_database_pkey PRIMARY KEY (id);


--
-- Name: metabase_field metabase_field_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.metabase_field
    ADD CONSTRAINT metabase_field_pkey PRIMARY KEY (id);


--
-- Name: metabase_fieldvalues metabase_fieldvalues_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.metabase_fieldvalues
    ADD CONSTRAINT metabase_fieldvalues_pkey PRIMARY KEY (id);


--
-- Name: metabase_table metabase_table_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.metabase_table
    ADD CONSTRAINT metabase_table_pkey PRIMARY KEY (id);


--
-- Name: metric metric_entity_id_key; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.metric
    ADD CONSTRAINT metric_entity_id_key UNIQUE (entity_id);


--
-- Name: metric_important_field metric_important_field_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.metric_important_field
    ADD CONSTRAINT metric_important_field_pkey PRIMARY KEY (id);


--
-- Name: metric metric_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.metric
    ADD CONSTRAINT metric_pkey PRIMARY KEY (id);


--
-- Name: model_index model_index_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.model_index
    ADD CONSTRAINT model_index_pkey PRIMARY KEY (id);


--
-- Name: moderation_review moderation_review_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.moderation_review
    ADD CONSTRAINT moderation_review_pkey PRIMARY KEY (id);


--
-- Name: native_query_snippet native_query_snippet_entity_id_key; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.native_query_snippet
    ADD CONSTRAINT native_query_snippet_entity_id_key UNIQUE (entity_id);


--
-- Name: native_query_snippet native_query_snippet_name_key; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.native_query_snippet
    ADD CONSTRAINT native_query_snippet_name_key UNIQUE (name);


--
-- Name: native_query_snippet native_query_snippet_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.native_query_snippet
    ADD CONSTRAINT native_query_snippet_pkey PRIMARY KEY (id);


--
-- Name: parameter_card parameter_card_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.parameter_card
    ADD CONSTRAINT parameter_card_pkey PRIMARY KEY (id);


--
-- Name: permissions_group permissions_group_entity_id_key; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.permissions_group
    ADD CONSTRAINT permissions_group_entity_id_key UNIQUE (entity_id);


--
-- Name: permissions permissions_group_id_object_key; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.permissions
    ADD CONSTRAINT permissions_group_id_object_key UNIQUE (group_id, object);


--
-- Name: permissions_group_membership permissions_group_membership_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.permissions_group_membership
    ADD CONSTRAINT permissions_group_membership_pkey PRIMARY KEY (id);


--
-- Name: permissions_group permissions_group_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.permissions_group
    ADD CONSTRAINT permissions_group_pkey PRIMARY KEY (id);


--
-- Name: permissions permissions_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.permissions
    ADD CONSTRAINT permissions_pkey PRIMARY KEY (id);


--
-- Name: permissions_revision permissions_revision_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.permissions_revision
    ADD CONSTRAINT permissions_revision_pkey PRIMARY KEY (id);


--
-- Name: persisted_info persisted_info_card_id_key; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.persisted_info
    ADD CONSTRAINT persisted_info_card_id_key UNIQUE (card_id);


--
-- Name: persisted_info persisted_info_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.persisted_info
    ADD CONSTRAINT persisted_info_pkey PRIMARY KEY (id);


--
-- Name: http_action pk_http_action; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.http_action
    ADD CONSTRAINT pk_http_action PRIMARY KEY (action_id);


--
-- Name: implicit_action pk_implicit_action; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.implicit_action
    ADD CONSTRAINT pk_implicit_action PRIMARY KEY (action_id);


--
-- Name: qrtz_blob_triggers pk_qrtz_blob_triggers; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.qrtz_blob_triggers
    ADD CONSTRAINT pk_qrtz_blob_triggers PRIMARY KEY (sched_name, trigger_name, trigger_group);


--
-- Name: qrtz_calendars pk_qrtz_calendars; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.qrtz_calendars
    ADD CONSTRAINT pk_qrtz_calendars PRIMARY KEY (sched_name, calendar_name);


--
-- Name: qrtz_cron_triggers pk_qrtz_cron_triggers; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.qrtz_cron_triggers
    ADD CONSTRAINT pk_qrtz_cron_triggers PRIMARY KEY (sched_name, trigger_name, trigger_group);


--
-- Name: qrtz_fired_triggers pk_qrtz_fired_triggers; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.qrtz_fired_triggers
    ADD CONSTRAINT pk_qrtz_fired_triggers PRIMARY KEY (sched_name, entry_id);


--
-- Name: qrtz_job_details pk_qrtz_job_details; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.qrtz_job_details
    ADD CONSTRAINT pk_qrtz_job_details PRIMARY KEY (sched_name, job_name, job_group);


--
-- Name: qrtz_locks pk_qrtz_locks; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.qrtz_locks
    ADD CONSTRAINT pk_qrtz_locks PRIMARY KEY (sched_name, lock_name);


--
-- Name: qrtz_scheduler_state pk_qrtz_scheduler_state; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.qrtz_scheduler_state
    ADD CONSTRAINT pk_qrtz_scheduler_state PRIMARY KEY (sched_name, instance_name);


--
-- Name: qrtz_simple_triggers pk_qrtz_simple_triggers; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.qrtz_simple_triggers
    ADD CONSTRAINT pk_qrtz_simple_triggers PRIMARY KEY (sched_name, trigger_name, trigger_group);


--
-- Name: qrtz_simprop_triggers pk_qrtz_simprop_triggers; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.qrtz_simprop_triggers
    ADD CONSTRAINT pk_qrtz_simprop_triggers PRIMARY KEY (sched_name, trigger_name, trigger_group);


--
-- Name: qrtz_triggers pk_qrtz_triggers; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.qrtz_triggers
    ADD CONSTRAINT pk_qrtz_triggers PRIMARY KEY (sched_name, trigger_name, trigger_group);


--
-- Name: query_action pk_query_action; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.query_action
    ADD CONSTRAINT pk_query_action PRIMARY KEY (action_id);


--
-- Name: qrtz_paused_trigger_grps pk_sched_name; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.qrtz_paused_trigger_grps
    ADD CONSTRAINT pk_sched_name PRIMARY KEY (sched_name, trigger_group);


--
-- Name: pulse_card pulse_card_entity_id_key; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.pulse_card
    ADD CONSTRAINT pulse_card_entity_id_key UNIQUE (entity_id);


--
-- Name: pulse_card pulse_card_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.pulse_card
    ADD CONSTRAINT pulse_card_pkey PRIMARY KEY (id);


--
-- Name: pulse_channel pulse_channel_entity_id_key; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.pulse_channel
    ADD CONSTRAINT pulse_channel_entity_id_key UNIQUE (entity_id);


--
-- Name: pulse_channel pulse_channel_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.pulse_channel
    ADD CONSTRAINT pulse_channel_pkey PRIMARY KEY (id);


--
-- Name: pulse_channel_recipient pulse_channel_recipient_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.pulse_channel_recipient
    ADD CONSTRAINT pulse_channel_recipient_pkey PRIMARY KEY (id);


--
-- Name: pulse pulse_entity_id_key; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.pulse
    ADD CONSTRAINT pulse_entity_id_key UNIQUE (entity_id);


--
-- Name: pulse pulse_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.pulse
    ADD CONSTRAINT pulse_pkey PRIMARY KEY (id);


--
-- Name: query_cache query_cache_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.query_cache
    ADD CONSTRAINT query_cache_pkey PRIMARY KEY (query_hash);


--
-- Name: query_execution query_execution_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.query_execution
    ADD CONSTRAINT query_execution_pkey PRIMARY KEY (id);


--
-- Name: query_field query_field_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.query_field
    ADD CONSTRAINT query_field_pkey PRIMARY KEY (id);


--
-- Name: query query_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.query
    ADD CONSTRAINT query_pkey PRIMARY KEY (query_hash);


--
-- Name: recent_views recent_views_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.recent_views
    ADD CONSTRAINT recent_views_pkey PRIMARY KEY (id);


--
-- Name: report_card report_card_entity_id_key; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.report_card
    ADD CONSTRAINT report_card_entity_id_key UNIQUE (entity_id);


--
-- Name: report_card report_card_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.report_card
    ADD CONSTRAINT report_card_pkey PRIMARY KEY (id);


--
-- Name: report_card report_card_public_uuid_key; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.report_card
    ADD CONSTRAINT report_card_public_uuid_key UNIQUE (public_uuid);


--
-- Name: report_cardfavorite report_cardfavorite_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.report_cardfavorite
    ADD CONSTRAINT report_cardfavorite_pkey PRIMARY KEY (id);


--
-- Name: report_dashboard report_dashboard_entity_id_key; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.report_dashboard
    ADD CONSTRAINT report_dashboard_entity_id_key UNIQUE (entity_id);


--
-- Name: report_dashboard report_dashboard_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.report_dashboard
    ADD CONSTRAINT report_dashboard_pkey PRIMARY KEY (id);


--
-- Name: report_dashboard report_dashboard_public_uuid_key; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.report_dashboard
    ADD CONSTRAINT report_dashboard_public_uuid_key UNIQUE (public_uuid);


--
-- Name: report_dashboardcard report_dashboardcard_entity_id_key; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.report_dashboardcard
    ADD CONSTRAINT report_dashboardcard_entity_id_key UNIQUE (entity_id);


--
-- Name: report_dashboardcard report_dashboardcard_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.report_dashboardcard
    ADD CONSTRAINT report_dashboardcard_pkey PRIMARY KEY (id);


--
-- Name: revision revision_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.revision
    ADD CONSTRAINT revision_pkey PRIMARY KEY (id);


--
-- Name: secret secret_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.secret
    ADD CONSTRAINT secret_pkey PRIMARY KEY (id, version);


--
-- Name: segment segment_entity_id_key; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.segment
    ADD CONSTRAINT segment_entity_id_key UNIQUE (entity_id);


--
-- Name: segment segment_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.segment
    ADD CONSTRAINT segment_pkey PRIMARY KEY (id);


--
-- Name: setting setting_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.setting
    ADD CONSTRAINT setting_pkey PRIMARY KEY (key);


--
-- Name: task_history task_history_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.task_history
    ADD CONSTRAINT task_history_pkey PRIMARY KEY (id);


--
-- Name: timeline timeline_entity_id_key; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.timeline
    ADD CONSTRAINT timeline_entity_id_key UNIQUE (entity_id);


--
-- Name: timeline_event timeline_event_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.timeline_event
    ADD CONSTRAINT timeline_event_pkey PRIMARY KEY (id);


--
-- Name: timeline timeline_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.timeline
    ADD CONSTRAINT timeline_pkey PRIMARY KEY (id);


--
-- Name: bookmark_ordering unique_bookmark_user_id_ordering; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.bookmark_ordering
    ADD CONSTRAINT unique_bookmark_user_id_ordering UNIQUE (user_id, ordering);


--
-- Name: bookmark_ordering unique_bookmark_user_id_type_item_id; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.bookmark_ordering
    ADD CONSTRAINT unique_bookmark_user_id_type_item_id UNIQUE (user_id, type, item_id);


--
-- Name: card_bookmark unique_card_bookmark_user_id_card_id; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.card_bookmark
    ADD CONSTRAINT unique_card_bookmark_user_id_card_id UNIQUE (user_id, card_id);


--
-- Name: card_label unique_card_label_card_id_label_id; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.card_label
    ADD CONSTRAINT unique_card_label_card_id_label_id UNIQUE (card_id, label_id);


--
-- Name: collection_bookmark unique_collection_bookmark_user_id_collection_id; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.collection_bookmark
    ADD CONSTRAINT unique_collection_bookmark_user_id_collection_id UNIQUE (user_id, collection_id);


--
-- Name: collection unique_collection_personal_owner_id; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.collection
    ADD CONSTRAINT unique_collection_personal_owner_id UNIQUE (personal_owner_id);


--
-- Name: dashboard_bookmark unique_dashboard_bookmark_user_id_dashboard_id; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.dashboard_bookmark
    ADD CONSTRAINT unique_dashboard_bookmark_user_id_dashboard_id UNIQUE (user_id, dashboard_id);


--
-- Name: dashboard_favorite unique_dashboard_favorite_user_id_dashboard_id; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.dashboard_favorite
    ADD CONSTRAINT unique_dashboard_favorite_user_id_dashboard_id UNIQUE (user_id, dashboard_id);


--
-- Name: dimension unique_dimension_field_id; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.dimension
    ADD CONSTRAINT unique_dimension_field_id UNIQUE (field_id);


--
-- Name: sandboxes unique_gtap_table_id_group_id; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.sandboxes
    ADD CONSTRAINT unique_gtap_table_id_group_id UNIQUE (table_id, group_id);


--
-- Name: metric_important_field unique_metric_important_field_metric_id_field_id; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.metric_important_field
    ADD CONSTRAINT unique_metric_important_field_metric_id_field_id UNIQUE (metric_id, field_id);


--
-- Name: model_index_value unique_model_index_value_model_index_id_model_pk; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.model_index_value
    ADD CONSTRAINT unique_model_index_value_model_index_id_model_pk UNIQUE (model_index_id, model_pk);


--
-- Name: parameter_card unique_parameterized_object_card_parameter; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.parameter_card
    ADD CONSTRAINT unique_parameterized_object_card_parameter UNIQUE (parameterized_object_id, parameterized_object_type, parameter_id);


--
-- Name: permissions_group_membership unique_permissions_group_membership_user_id_group_id; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.permissions_group_membership
    ADD CONSTRAINT unique_permissions_group_membership_user_id_group_id UNIQUE (user_id, group_id);


--
-- Name: permissions_group unique_permissions_group_name; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.permissions_group
    ADD CONSTRAINT unique_permissions_group_name UNIQUE (name);


--
-- Name: user_parameter_value user_parameter_value_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.user_parameter_value
    ADD CONSTRAINT user_parameter_value_pkey PRIMARY KEY (id);


--
-- Name: view_log view_log_pkey; Type: CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.view_log
    ADD CONSTRAINT view_log_pkey PRIMARY KEY (id);


--
-- Name: idx_action_creator_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_action_creator_id ON public.action USING btree (creator_id);


--
-- Name: idx_action_made_public_by_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_action_made_public_by_id ON public.action USING btree (made_public_by_id);


--
-- Name: idx_action_model_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_action_model_id ON public.action USING btree (model_id);


--
-- Name: idx_action_public_uuid; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_action_public_uuid ON public.action USING btree (public_uuid);


--
-- Name: idx_api_key_created_by; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_api_key_created_by ON public.api_key USING btree (creator_id);


--
-- Name: idx_api_key_updated_by_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_api_key_updated_by_id ON public.api_key USING btree (updated_by_id);


--
-- Name: idx_api_key_user_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_api_key_user_id ON public.api_key USING btree (user_id);


--
-- Name: idx_application_permissions_revision_user_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_application_permissions_revision_user_id ON public.application_permissions_revision USING btree (user_id);


--
-- Name: idx_audit_log_entity_qualified_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_audit_log_entity_qualified_id ON public.audit_log USING btree ((
CASE
    WHEN ((model)::text = 'Dataset'::text) THEN ('card_'::text || model_id)
    WHEN (model_id IS NULL) THEN NULL::text
    ELSE ((lower((model)::text) || '_'::text) || model_id)
END));


--
-- Name: idx_bookmark_ordering_user_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_bookmark_ordering_user_id ON public.bookmark_ordering USING btree (user_id);


--
-- Name: idx_card_bookmark_card_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_card_bookmark_card_id ON public.card_bookmark USING btree (card_id);


--
-- Name: idx_card_bookmark_user_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_card_bookmark_user_id ON public.card_bookmark USING btree (user_id);


--
-- Name: idx_card_collection_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_card_collection_id ON public.report_card USING btree (collection_id);


--
-- Name: idx_card_creator_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_card_creator_id ON public.report_card USING btree (creator_id);


--
-- Name: idx_card_label_card_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_card_label_card_id ON public.card_label USING btree (card_id);


--
-- Name: idx_card_label_label_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_card_label_label_id ON public.card_label USING btree (label_id);


--
-- Name: idx_card_public_uuid; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_card_public_uuid ON public.report_card USING btree (public_uuid);


--
-- Name: idx_cardfavorite_card_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_cardfavorite_card_id ON public.report_cardfavorite USING btree (card_id);


--
-- Name: idx_cardfavorite_owner_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_cardfavorite_owner_id ON public.report_cardfavorite USING btree (owner_id);


--
-- Name: idx_collection_bookmark_collection_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_collection_bookmark_collection_id ON public.collection_bookmark USING btree (collection_id);


--
-- Name: idx_collection_bookmark_user_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_collection_bookmark_user_id ON public.collection_bookmark USING btree (user_id);


--
-- Name: idx_collection_location; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_collection_location ON public.collection USING btree (location);


--
-- Name: idx_collection_permission_graph_revision_user_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_collection_permission_graph_revision_user_id ON public.collection_permission_graph_revision USING btree (user_id);


--
-- Name: idx_collection_personal_owner_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_collection_personal_owner_id ON public.collection USING btree (personal_owner_id);


--
-- Name: idx_conn_impersonations_db_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_conn_impersonations_db_id ON public.connection_impersonations USING btree (db_id);


--
-- Name: idx_conn_impersonations_group_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_conn_impersonations_group_id ON public.connection_impersonations USING btree (group_id);


--
-- Name: idx_core_session_user_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_core_session_user_id ON public.core_session USING btree (user_id);


--
-- Name: idx_dashboard_bookmark_dashboard_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_dashboard_bookmark_dashboard_id ON public.dashboard_bookmark USING btree (dashboard_id);


--
-- Name: idx_dashboard_bookmark_user_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_dashboard_bookmark_user_id ON public.dashboard_bookmark USING btree (user_id);


--
-- Name: idx_dashboard_collection_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_dashboard_collection_id ON public.report_dashboard USING btree (collection_id);


--
-- Name: idx_dashboard_creator_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_dashboard_creator_id ON public.report_dashboard USING btree (creator_id);


--
-- Name: idx_dashboard_favorite_dashboard_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_dashboard_favorite_dashboard_id ON public.dashboard_favorite USING btree (dashboard_id);


--
-- Name: idx_dashboard_favorite_user_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_dashboard_favorite_user_id ON public.dashboard_favorite USING btree (user_id);


--
-- Name: idx_dashboard_public_uuid; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_dashboard_public_uuid ON public.report_dashboard USING btree (public_uuid);


--
-- Name: idx_dashboard_tab_dashboard_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_dashboard_tab_dashboard_id ON public.dashboard_tab USING btree (dashboard_id);


--
-- Name: idx_dashboardcard_card_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_dashboardcard_card_id ON public.report_dashboardcard USING btree (card_id);


--
-- Name: idx_dashboardcard_dashboard_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_dashboardcard_dashboard_id ON public.report_dashboardcard USING btree (dashboard_id);


--
-- Name: idx_dashboardcard_series_card_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_dashboardcard_series_card_id ON public.dashboardcard_series USING btree (card_id);


--
-- Name: idx_dashboardcard_series_dashboardcard_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_dashboardcard_series_dashboardcard_id ON public.dashboardcard_series USING btree (dashboardcard_id);


--
-- Name: idx_data_permissions_db_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_data_permissions_db_id ON public.data_permissions USING btree (db_id);


--
-- Name: idx_data_permissions_group_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_data_permissions_group_id ON public.data_permissions USING btree (group_id);


--
-- Name: idx_data_permissions_table_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_data_permissions_table_id ON public.data_permissions USING btree (table_id);


--
-- Name: idx_dependency_dependent_on_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_dependency_dependent_on_id ON public.dependency USING btree (dependent_on_id);


--
-- Name: idx_dependency_dependent_on_model; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_dependency_dependent_on_model ON public.dependency USING btree (dependent_on_model);


--
-- Name: idx_dependency_model; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_dependency_model ON public.dependency USING btree (model);


--
-- Name: idx_dependency_model_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_dependency_model_id ON public.dependency USING btree (model_id);


--
-- Name: idx_dimension_field_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_dimension_field_id ON public.dimension USING btree (field_id);


--
-- Name: idx_dimension_human_readable_field_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_dimension_human_readable_field_id ON public.dimension USING btree (human_readable_field_id);


--
-- Name: idx_field_entity_qualified_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_field_entity_qualified_id ON public.metabase_field USING btree ((('field_'::text || id)));


--
-- Name: idx_field_name_lower; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_field_name_lower ON public.metabase_field USING btree (lower((name)::text));


--
-- Name: idx_field_parent_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_field_parent_id ON public.metabase_field USING btree (parent_id);


--
-- Name: idx_field_table_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_field_table_id ON public.metabase_field USING btree (table_id);


--
-- Name: idx_field_usage_field_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_field_usage_field_id ON public.field_usage USING btree (field_id);


--
-- Name: idx_fieldvalues_field_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_fieldvalues_field_id ON public.metabase_fieldvalues USING btree (field_id);


--
-- Name: idx_gtap_table_id_group_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_gtap_table_id_group_id ON public.sandboxes USING btree (table_id, group_id);


--
-- Name: idx_label_slug; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_label_slug ON public.label USING btree (slug);


--
-- Name: idx_lower_email; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_lower_email ON public.core_user USING btree (lower((email)::text));


--
-- Name: idx_metabase_database_creator_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_metabase_database_creator_id ON public.metabase_database USING btree (creator_id);


--
-- Name: idx_metabase_table_db_id_schema; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_metabase_table_db_id_schema ON public.metabase_table USING btree (db_id, schema);


--
-- Name: idx_metabase_table_show_in_getting_started; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_metabase_table_show_in_getting_started ON public.metabase_table USING btree (show_in_getting_started);


--
-- Name: idx_metric_creator_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_metric_creator_id ON public.metric USING btree (creator_id);


--
-- Name: idx_metric_important_field_field_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_metric_important_field_field_id ON public.metric_important_field USING btree (field_id);


--
-- Name: idx_metric_important_field_metric_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_metric_important_field_metric_id ON public.metric_important_field USING btree (metric_id);


--
-- Name: idx_metric_show_in_getting_started; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_metric_show_in_getting_started ON public.metric USING btree (show_in_getting_started);


--
-- Name: idx_metric_table_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_metric_table_id ON public.metric USING btree (table_id);


--
-- Name: idx_model_index_creator_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_model_index_creator_id ON public.model_index USING btree (creator_id);


--
-- Name: idx_model_index_model_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_model_index_model_id ON public.model_index USING btree (model_id);


--
-- Name: idx_moderation_review_item_type_item_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_moderation_review_item_type_item_id ON public.moderation_review USING btree (moderated_item_type, moderated_item_id);


--
-- Name: idx_native_query_snippet_creator_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_native_query_snippet_creator_id ON public.native_query_snippet USING btree (creator_id);


--
-- Name: idx_parameter_card_card_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_parameter_card_card_id ON public.parameter_card USING btree (card_id);


--
-- Name: idx_parameter_card_parameterized_object_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_parameter_card_parameterized_object_id ON public.parameter_card USING btree (parameterized_object_id);


--
-- Name: idx_permissions_group_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_permissions_group_id ON public.permissions USING btree (group_id);


--
-- Name: idx_permissions_group_id_object; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_permissions_group_id_object ON public.permissions USING btree (group_id, object);


--
-- Name: idx_permissions_group_membership_group_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_permissions_group_membership_group_id ON public.permissions_group_membership USING btree (group_id);


--
-- Name: idx_permissions_group_membership_group_id_user_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_permissions_group_membership_group_id_user_id ON public.permissions_group_membership USING btree (group_id, user_id);


--
-- Name: idx_permissions_group_membership_user_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_permissions_group_membership_user_id ON public.permissions_group_membership USING btree (user_id);


--
-- Name: idx_permissions_group_name; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_permissions_group_name ON public.permissions_group USING btree (name);


--
-- Name: idx_permissions_object; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_permissions_object ON public.permissions USING btree (object);


--
-- Name: idx_permissions_revision_user_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_permissions_revision_user_id ON public.permissions_revision USING btree (user_id);


--
-- Name: idx_persisted_info_creator_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_persisted_info_creator_id ON public.persisted_info USING btree (creator_id);


--
-- Name: idx_persisted_info_database_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_persisted_info_database_id ON public.persisted_info USING btree (database_id);


--
-- Name: idx_pulse_card_card_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_pulse_card_card_id ON public.pulse_card USING btree (card_id);


--
-- Name: idx_pulse_card_dashboard_card_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_pulse_card_dashboard_card_id ON public.pulse_card USING btree (dashboard_card_id);


--
-- Name: idx_pulse_card_pulse_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_pulse_card_pulse_id ON public.pulse_card USING btree (pulse_id);


--
-- Name: idx_pulse_channel_pulse_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_pulse_channel_pulse_id ON public.pulse_channel USING btree (pulse_id);


--
-- Name: idx_pulse_channel_recipient_pulse_channel_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_pulse_channel_recipient_pulse_channel_id ON public.pulse_channel_recipient USING btree (pulse_channel_id);


--
-- Name: idx_pulse_channel_recipient_user_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_pulse_channel_recipient_user_id ON public.pulse_channel_recipient USING btree (user_id);


--
-- Name: idx_pulse_channel_schedule_type; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_pulse_channel_schedule_type ON public.pulse_channel USING btree (schedule_type);


--
-- Name: idx_pulse_collection_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_pulse_collection_id ON public.pulse USING btree (collection_id);


--
-- Name: idx_pulse_creator_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_pulse_creator_id ON public.pulse USING btree (creator_id);


--
-- Name: idx_pulse_dashboard_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_pulse_dashboard_id ON public.pulse USING btree (dashboard_id);


--
-- Name: idx_qrtz_ft_inst_job_req_rcvry; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_qrtz_ft_inst_job_req_rcvry ON public.qrtz_fired_triggers USING btree (sched_name, instance_name, requests_recovery);


--
-- Name: idx_qrtz_ft_j_g; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_qrtz_ft_j_g ON public.qrtz_fired_triggers USING btree (sched_name, job_name, job_group);


--
-- Name: idx_qrtz_ft_jg; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_qrtz_ft_jg ON public.qrtz_fired_triggers USING btree (sched_name, job_group);


--
-- Name: idx_qrtz_ft_t_g; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_qrtz_ft_t_g ON public.qrtz_fired_triggers USING btree (sched_name, trigger_name, trigger_group);


--
-- Name: idx_qrtz_ft_tg; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_qrtz_ft_tg ON public.qrtz_fired_triggers USING btree (sched_name, trigger_group);


--
-- Name: idx_qrtz_ft_trig_inst_name; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_qrtz_ft_trig_inst_name ON public.qrtz_fired_triggers USING btree (sched_name, instance_name);


--
-- Name: idx_qrtz_j_grp; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_qrtz_j_grp ON public.qrtz_job_details USING btree (sched_name, job_group);


--
-- Name: idx_qrtz_j_req_recovery; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_qrtz_j_req_recovery ON public.qrtz_job_details USING btree (sched_name, requests_recovery);


--
-- Name: idx_qrtz_t_c; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_qrtz_t_c ON public.qrtz_triggers USING btree (sched_name, calendar_name);


--
-- Name: idx_qrtz_t_g; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_qrtz_t_g ON public.qrtz_triggers USING btree (sched_name, trigger_group);


--
-- Name: idx_qrtz_t_j; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_qrtz_t_j ON public.qrtz_triggers USING btree (sched_name, job_name, job_group);


--
-- Name: idx_qrtz_t_jg; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_qrtz_t_jg ON public.qrtz_triggers USING btree (sched_name, job_group);


--
-- Name: idx_qrtz_t_n_g_state; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_qrtz_t_n_g_state ON public.qrtz_triggers USING btree (sched_name, trigger_group, trigger_state);


--
-- Name: idx_qrtz_t_n_state; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_qrtz_t_n_state ON public.qrtz_triggers USING btree (sched_name, trigger_name, trigger_group, trigger_state);


--
-- Name: idx_qrtz_t_next_fire_time; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_qrtz_t_next_fire_time ON public.qrtz_triggers USING btree (sched_name, next_fire_time);


--
-- Name: idx_qrtz_t_nft_misfire; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_qrtz_t_nft_misfire ON public.qrtz_triggers USING btree (sched_name, misfire_instr, next_fire_time);


--
-- Name: idx_qrtz_t_nft_st; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_qrtz_t_nft_st ON public.qrtz_triggers USING btree (sched_name, trigger_state, next_fire_time);


--
-- Name: idx_qrtz_t_nft_st_misfire; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_qrtz_t_nft_st_misfire ON public.qrtz_triggers USING btree (sched_name, misfire_instr, next_fire_time, trigger_state);


--
-- Name: idx_qrtz_t_nft_st_misfire_grp; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_qrtz_t_nft_st_misfire_grp ON public.qrtz_triggers USING btree (sched_name, misfire_instr, next_fire_time, trigger_group, trigger_state);


--
-- Name: idx_qrtz_t_state; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_qrtz_t_state ON public.qrtz_triggers USING btree (sched_name, trigger_state);


--
-- Name: idx_query_action_database_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_query_action_database_id ON public.query_action USING btree (database_id);


--
-- Name: idx_query_cache_updated_at; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_query_cache_updated_at ON public.query_cache USING btree (updated_at);


--
-- Name: idx_query_execution_action_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_query_execution_action_id ON public.query_execution USING btree (action_id);


--
-- Name: idx_query_execution_card_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_query_execution_card_id ON public.query_execution USING btree (card_id);


--
-- Name: idx_query_execution_card_id_started_at; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_query_execution_card_id_started_at ON public.query_execution USING btree (card_id, started_at);


--
-- Name: idx_query_execution_card_qualified_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_query_execution_card_qualified_id ON public.query_execution USING btree ((('card_'::text || card_id)));


--
-- Name: idx_query_execution_context; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_query_execution_context ON public.query_execution USING btree (context);


--
-- Name: idx_query_execution_executor_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_query_execution_executor_id ON public.query_execution USING btree (executor_id);


--
-- Name: idx_query_execution_query_hash_started_at; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_query_execution_query_hash_started_at ON public.query_execution USING btree (hash, started_at);


--
-- Name: idx_query_execution_started_at; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_query_execution_started_at ON public.query_execution USING btree (started_at);


--
-- Name: idx_query_field_card_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_query_field_card_id ON public.query_field USING btree (card_id);


--
-- Name: idx_query_field_field_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_query_field_field_id ON public.query_field USING btree (field_id);


--
-- Name: idx_recent_views_user_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_recent_views_user_id ON public.recent_views USING btree (user_id);


--
-- Name: idx_report_card_database_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_report_card_database_id ON public.report_card USING btree (database_id);


--
-- Name: idx_report_card_made_public_by_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_report_card_made_public_by_id ON public.report_card USING btree (made_public_by_id);


--
-- Name: idx_report_card_table_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_report_card_table_id ON public.report_card USING btree (table_id);


--
-- Name: idx_report_dashboard_made_public_by_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_report_dashboard_made_public_by_id ON public.report_dashboard USING btree (made_public_by_id);


--
-- Name: idx_report_dashboard_show_in_getting_started; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_report_dashboard_show_in_getting_started ON public.report_dashboard USING btree (show_in_getting_started);


--
-- Name: idx_report_dashboardcard_action_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_report_dashboardcard_action_id ON public.report_dashboardcard USING btree (action_id);


--
-- Name: idx_report_dashboardcard_dashboard_tab_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_report_dashboardcard_dashboard_tab_id ON public.report_dashboardcard USING btree (dashboard_tab_id);


--
-- Name: idx_revision_model_model_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_revision_model_model_id ON public.revision USING btree (model, model_id);


--
-- Name: idx_revision_most_recent; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_revision_most_recent ON public.revision USING btree (most_recent);


--
-- Name: idx_revision_user_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_revision_user_id ON public.revision USING btree (user_id);


--
-- Name: idx_sandboxes_card_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_sandboxes_card_id ON public.sandboxes USING btree (card_id);


--
-- Name: idx_sandboxes_permission_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_sandboxes_permission_id ON public.sandboxes USING btree (permission_id);


--
-- Name: idx_secret_creator_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_secret_creator_id ON public.secret USING btree (creator_id);


--
-- Name: idx_segment_creator_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_segment_creator_id ON public.segment USING btree (creator_id);


--
-- Name: idx_segment_show_in_getting_started; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_segment_show_in_getting_started ON public.segment USING btree (show_in_getting_started);


--
-- Name: idx_segment_table_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_segment_table_id ON public.segment USING btree (table_id);


--
-- Name: idx_session_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_session_id ON public.login_history USING btree (session_id);


--
-- Name: idx_snippet_collection_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_snippet_collection_id ON public.native_query_snippet USING btree (collection_id);


--
-- Name: idx_snippet_name; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_snippet_name ON public.native_query_snippet USING btree (name);


--
-- Name: idx_table_db_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_table_db_id ON public.metabase_table USING btree (db_id);


--
-- Name: idx_table_privileges_role; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_table_privileges_role ON public.table_privileges USING btree (role);


--
-- Name: idx_table_privileges_table_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_table_privileges_table_id ON public.table_privileges USING btree (table_id);


--
-- Name: idx_task_history_db_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_task_history_db_id ON public.task_history USING btree (db_id);


--
-- Name: idx_task_history_end_time; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_task_history_end_time ON public.task_history USING btree (ended_at);


--
-- Name: idx_task_history_started_at; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_task_history_started_at ON public.task_history USING btree (started_at);


--
-- Name: idx_timeline_collection_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_timeline_collection_id ON public.timeline USING btree (collection_id);


--
-- Name: idx_timeline_creator_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_timeline_creator_id ON public.timeline USING btree (creator_id);


--
-- Name: idx_timeline_event_creator_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_timeline_event_creator_id ON public.timeline_event USING btree (creator_id);


--
-- Name: idx_timeline_event_timeline_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_timeline_event_timeline_id ON public.timeline_event USING btree (timeline_id);


--
-- Name: idx_timeline_event_timeline_id_timestamp; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_timeline_event_timeline_id_timestamp ON public.timeline_event USING btree (timeline_id, "timestamp");


--
-- Name: idx_timestamp; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_timestamp ON public.login_history USING btree ("timestamp");


--
-- Name: idx_uniq_field_table_id_parent_id_name_2col; Type: INDEX; Schema: public; Owner: aman
--

CREATE UNIQUE INDEX idx_uniq_field_table_id_parent_id_name_2col ON public.metabase_field USING btree (table_id, name) WHERE (parent_id IS NULL);


--
-- Name: idx_uniq_table_db_id_schema_name_2col; Type: INDEX; Schema: public; Owner: aman
--

CREATE UNIQUE INDEX idx_uniq_table_db_id_schema_name_2col ON public.metabase_table USING btree (db_id, name) WHERE (schema IS NULL);


--
-- Name: idx_user_full_name; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_user_full_name ON public.core_user USING btree (((((first_name)::text || ' '::text) || (last_name)::text)));


--
-- Name: idx_user_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_user_id ON public.login_history USING btree (user_id);


--
-- Name: idx_user_id_device_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_user_id_device_id ON public.login_history USING btree (session_id, device_id);


--
-- Name: idx_user_id_timestamp; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_user_id_timestamp ON public.login_history USING btree (user_id, "timestamp");


--
-- Name: idx_user_parameter_value_user_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_user_parameter_value_user_id ON public.user_parameter_value USING btree (user_id);


--
-- Name: idx_user_qualified_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_user_qualified_id ON public.core_user USING btree ((('user_'::text || id)));


--
-- Name: idx_view_log_entity_qualified_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_view_log_entity_qualified_id ON public.view_log USING btree (((((model)::text || '_'::text) || model_id)));


--
-- Name: idx_view_log_model_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_view_log_model_id ON public.view_log USING btree (model_id);


--
-- Name: idx_view_log_timestamp; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_view_log_timestamp ON public.view_log USING btree ("timestamp");


--
-- Name: idx_view_log_user_id; Type: INDEX; Schema: public; Owner: aman
--

CREATE INDEX idx_view_log_user_id ON public.view_log USING btree (user_id);


--
-- Name: action fk_action_creator_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.action
    ADD CONSTRAINT fk_action_creator_id FOREIGN KEY (creator_id) REFERENCES public.core_user(id);


--
-- Name: action fk_action_made_public_by_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.action
    ADD CONSTRAINT fk_action_made_public_by_id FOREIGN KEY (made_public_by_id) REFERENCES public.core_user(id) ON DELETE CASCADE;


--
-- Name: action fk_action_model_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.action
    ADD CONSTRAINT fk_action_model_id FOREIGN KEY (model_id) REFERENCES public.report_card(id) ON DELETE CASCADE;


--
-- Name: api_key fk_api_key_created_by_user_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.api_key
    ADD CONSTRAINT fk_api_key_created_by_user_id FOREIGN KEY (creator_id) REFERENCES public.core_user(id);


--
-- Name: api_key fk_api_key_updated_by_id_user_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.api_key
    ADD CONSTRAINT fk_api_key_updated_by_id_user_id FOREIGN KEY (updated_by_id) REFERENCES public.core_user(id);


--
-- Name: api_key fk_api_key_user_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.api_key
    ADD CONSTRAINT fk_api_key_user_id FOREIGN KEY (user_id) REFERENCES public.core_user(id);


--
-- Name: bookmark_ordering fk_bookmark_ordering_user_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.bookmark_ordering
    ADD CONSTRAINT fk_bookmark_ordering_user_id FOREIGN KEY (user_id) REFERENCES public.core_user(id) ON DELETE CASCADE;


--
-- Name: card_bookmark fk_card_bookmark_dashboard_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.card_bookmark
    ADD CONSTRAINT fk_card_bookmark_dashboard_id FOREIGN KEY (card_id) REFERENCES public.report_card(id) ON DELETE CASCADE;


--
-- Name: card_bookmark fk_card_bookmark_user_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.card_bookmark
    ADD CONSTRAINT fk_card_bookmark_user_id FOREIGN KEY (user_id) REFERENCES public.core_user(id) ON DELETE CASCADE;


--
-- Name: report_card fk_card_collection_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.report_card
    ADD CONSTRAINT fk_card_collection_id FOREIGN KEY (collection_id) REFERENCES public.collection(id) ON DELETE SET NULL;


--
-- Name: card_label fk_card_label_ref_card_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.card_label
    ADD CONSTRAINT fk_card_label_ref_card_id FOREIGN KEY (card_id) REFERENCES public.report_card(id) ON DELETE CASCADE;


--
-- Name: card_label fk_card_label_ref_label_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.card_label
    ADD CONSTRAINT fk_card_label_ref_label_id FOREIGN KEY (label_id) REFERENCES public.label(id) ON DELETE CASCADE;


--
-- Name: report_card fk_card_made_public_by_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.report_card
    ADD CONSTRAINT fk_card_made_public_by_id FOREIGN KEY (made_public_by_id) REFERENCES public.core_user(id) ON DELETE CASCADE;


--
-- Name: report_card fk_card_ref_user_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.report_card
    ADD CONSTRAINT fk_card_ref_user_id FOREIGN KEY (creator_id) REFERENCES public.core_user(id) ON DELETE CASCADE;


--
-- Name: report_cardfavorite fk_cardfavorite_ref_card_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.report_cardfavorite
    ADD CONSTRAINT fk_cardfavorite_ref_card_id FOREIGN KEY (card_id) REFERENCES public.report_card(id) ON DELETE CASCADE;


--
-- Name: report_cardfavorite fk_cardfavorite_ref_user_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.report_cardfavorite
    ADD CONSTRAINT fk_cardfavorite_ref_user_id FOREIGN KEY (owner_id) REFERENCES public.core_user(id) ON DELETE CASCADE;


--
-- Name: collection_bookmark fk_collection_bookmark_collection_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.collection_bookmark
    ADD CONSTRAINT fk_collection_bookmark_collection_id FOREIGN KEY (collection_id) REFERENCES public.collection(id) ON DELETE CASCADE;


--
-- Name: collection_bookmark fk_collection_bookmark_user_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.collection_bookmark
    ADD CONSTRAINT fk_collection_bookmark_user_id FOREIGN KEY (user_id) REFERENCES public.core_user(id) ON DELETE CASCADE;


--
-- Name: collection fk_collection_personal_owner_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.collection
    ADD CONSTRAINT fk_collection_personal_owner_id FOREIGN KEY (personal_owner_id) REFERENCES public.core_user(id) ON DELETE CASCADE;


--
-- Name: collection_permission_graph_revision fk_collection_revision_user_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.collection_permission_graph_revision
    ADD CONSTRAINT fk_collection_revision_user_id FOREIGN KEY (user_id) REFERENCES public.core_user(id) ON DELETE CASCADE;


--
-- Name: connection_impersonations fk_conn_impersonation_db_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.connection_impersonations
    ADD CONSTRAINT fk_conn_impersonation_db_id FOREIGN KEY (db_id) REFERENCES public.metabase_database(id) ON DELETE CASCADE;


--
-- Name: connection_impersonations fk_conn_impersonation_group_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.connection_impersonations
    ADD CONSTRAINT fk_conn_impersonation_group_id FOREIGN KEY (group_id) REFERENCES public.permissions_group(id) ON DELETE CASCADE;


--
-- Name: dashboard_bookmark fk_dashboard_bookmark_dashboard_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.dashboard_bookmark
    ADD CONSTRAINT fk_dashboard_bookmark_dashboard_id FOREIGN KEY (dashboard_id) REFERENCES public.report_dashboard(id) ON DELETE CASCADE;


--
-- Name: dashboard_bookmark fk_dashboard_bookmark_user_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.dashboard_bookmark
    ADD CONSTRAINT fk_dashboard_bookmark_user_id FOREIGN KEY (user_id) REFERENCES public.core_user(id) ON DELETE CASCADE;


--
-- Name: report_dashboard fk_dashboard_collection_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.report_dashboard
    ADD CONSTRAINT fk_dashboard_collection_id FOREIGN KEY (collection_id) REFERENCES public.collection(id) ON DELETE SET NULL;


--
-- Name: dashboard_favorite fk_dashboard_favorite_dashboard_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.dashboard_favorite
    ADD CONSTRAINT fk_dashboard_favorite_dashboard_id FOREIGN KEY (dashboard_id) REFERENCES public.report_dashboard(id) ON DELETE CASCADE;


--
-- Name: dashboard_favorite fk_dashboard_favorite_user_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.dashboard_favorite
    ADD CONSTRAINT fk_dashboard_favorite_user_id FOREIGN KEY (user_id) REFERENCES public.core_user(id) ON DELETE CASCADE;


--
-- Name: report_dashboard fk_dashboard_made_public_by_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.report_dashboard
    ADD CONSTRAINT fk_dashboard_made_public_by_id FOREIGN KEY (made_public_by_id) REFERENCES public.core_user(id) ON DELETE CASCADE;


--
-- Name: report_dashboard fk_dashboard_ref_user_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.report_dashboard
    ADD CONSTRAINT fk_dashboard_ref_user_id FOREIGN KEY (creator_id) REFERENCES public.core_user(id) ON DELETE CASCADE;


--
-- Name: dashboard_tab fk_dashboard_tab_ref_dashboard_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.dashboard_tab
    ADD CONSTRAINT fk_dashboard_tab_ref_dashboard_id FOREIGN KEY (dashboard_id) REFERENCES public.report_dashboard(id) ON DELETE CASCADE;


--
-- Name: report_dashboardcard fk_dashboardcard_ref_card_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.report_dashboardcard
    ADD CONSTRAINT fk_dashboardcard_ref_card_id FOREIGN KEY (card_id) REFERENCES public.report_card(id) ON DELETE CASCADE;


--
-- Name: report_dashboardcard fk_dashboardcard_ref_dashboard_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.report_dashboardcard
    ADD CONSTRAINT fk_dashboardcard_ref_dashboard_id FOREIGN KEY (dashboard_id) REFERENCES public.report_dashboard(id) ON DELETE CASCADE;


--
-- Name: dashboardcard_series fk_dashboardcard_series_ref_card_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.dashboardcard_series
    ADD CONSTRAINT fk_dashboardcard_series_ref_card_id FOREIGN KEY (card_id) REFERENCES public.report_card(id) ON DELETE CASCADE;


--
-- Name: dashboardcard_series fk_dashboardcard_series_ref_dashboardcard_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.dashboardcard_series
    ADD CONSTRAINT fk_dashboardcard_series_ref_dashboardcard_id FOREIGN KEY (dashboardcard_id) REFERENCES public.report_dashboardcard(id) ON DELETE CASCADE;


--
-- Name: data_permissions fk_data_permissions_ref_db_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.data_permissions
    ADD CONSTRAINT fk_data_permissions_ref_db_id FOREIGN KEY (db_id) REFERENCES public.metabase_database(id) ON DELETE CASCADE;


--
-- Name: data_permissions fk_data_permissions_ref_permissions_group; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.data_permissions
    ADD CONSTRAINT fk_data_permissions_ref_permissions_group FOREIGN KEY (group_id) REFERENCES public.permissions_group(id) ON DELETE CASCADE;


--
-- Name: data_permissions fk_data_permissions_ref_table_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.data_permissions
    ADD CONSTRAINT fk_data_permissions_ref_table_id FOREIGN KEY (table_id) REFERENCES public.metabase_table(id) ON DELETE CASCADE;


--
-- Name: metabase_database fk_database_creator_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.metabase_database
    ADD CONSTRAINT fk_database_creator_id FOREIGN KEY (creator_id) REFERENCES public.core_user(id) ON DELETE SET NULL;


--
-- Name: dimension fk_dimension_displayfk_ref_field_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.dimension
    ADD CONSTRAINT fk_dimension_displayfk_ref_field_id FOREIGN KEY (human_readable_field_id) REFERENCES public.metabase_field(id) ON DELETE CASCADE;


--
-- Name: dimension fk_dimension_ref_field_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.dimension
    ADD CONSTRAINT fk_dimension_ref_field_id FOREIGN KEY (field_id) REFERENCES public.metabase_field(id) ON DELETE CASCADE;


--
-- Name: timeline_event fk_event_creator_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.timeline_event
    ADD CONSTRAINT fk_event_creator_id FOREIGN KEY (creator_id) REFERENCES public.core_user(id) ON DELETE CASCADE;


--
-- Name: timeline_event fk_events_timeline_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.timeline_event
    ADD CONSTRAINT fk_events_timeline_id FOREIGN KEY (timeline_id) REFERENCES public.timeline(id) ON DELETE CASCADE;


--
-- Name: metabase_field fk_field_parent_ref_field_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.metabase_field
    ADD CONSTRAINT fk_field_parent_ref_field_id FOREIGN KEY (parent_id) REFERENCES public.metabase_field(id) ON DELETE CASCADE;


--
-- Name: metabase_field fk_field_ref_table_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.metabase_field
    ADD CONSTRAINT fk_field_ref_table_id FOREIGN KEY (table_id) REFERENCES public.metabase_table(id) ON DELETE CASCADE;


--
-- Name: field_usage fk_field_usage_field_id_metabase_field_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.field_usage
    ADD CONSTRAINT fk_field_usage_field_id_metabase_field_id FOREIGN KEY (field_id) REFERENCES public.metabase_field(id) ON DELETE CASCADE;


--
-- Name: field_usage fk_field_usage_query_execution_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.field_usage
    ADD CONSTRAINT fk_field_usage_query_execution_id FOREIGN KEY (query_execution_id) REFERENCES public.query_execution(id) ON DELETE CASCADE;


--
-- Name: metabase_fieldvalues fk_fieldvalues_ref_field_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.metabase_fieldvalues
    ADD CONSTRAINT fk_fieldvalues_ref_field_id FOREIGN KEY (field_id) REFERENCES public.metabase_field(id) ON DELETE CASCADE;


--
-- Name: application_permissions_revision fk_general_permissions_revision_user_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.application_permissions_revision
    ADD CONSTRAINT fk_general_permissions_revision_user_id FOREIGN KEY (user_id) REFERENCES public.core_user(id);


--
-- Name: sandboxes fk_gtap_card_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.sandboxes
    ADD CONSTRAINT fk_gtap_card_id FOREIGN KEY (card_id) REFERENCES public.report_card(id) ON DELETE CASCADE;


--
-- Name: sandboxes fk_gtap_group_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.sandboxes
    ADD CONSTRAINT fk_gtap_group_id FOREIGN KEY (group_id) REFERENCES public.permissions_group(id) ON DELETE CASCADE;


--
-- Name: sandboxes fk_gtap_table_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.sandboxes
    ADD CONSTRAINT fk_gtap_table_id FOREIGN KEY (table_id) REFERENCES public.metabase_table(id) ON DELETE CASCADE;


--
-- Name: http_action fk_http_action_ref_action_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.http_action
    ADD CONSTRAINT fk_http_action_ref_action_id FOREIGN KEY (action_id) REFERENCES public.action(id) ON DELETE CASCADE;


--
-- Name: implicit_action fk_implicit_action_action_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.implicit_action
    ADD CONSTRAINT fk_implicit_action_action_id FOREIGN KEY (action_id) REFERENCES public.action(id) ON DELETE CASCADE;


--
-- Name: login_history fk_login_history_session_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.login_history
    ADD CONSTRAINT fk_login_history_session_id FOREIGN KEY (session_id) REFERENCES public.core_session(id) ON DELETE SET NULL;


--
-- Name: login_history fk_login_history_user_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.login_history
    ADD CONSTRAINT fk_login_history_user_id FOREIGN KEY (user_id) REFERENCES public.core_user(id) ON DELETE CASCADE;


--
-- Name: metric_important_field fk_metric_important_field_metabase_field_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.metric_important_field
    ADD CONSTRAINT fk_metric_important_field_metabase_field_id FOREIGN KEY (field_id) REFERENCES public.metabase_field(id) ON DELETE CASCADE;


--
-- Name: metric_important_field fk_metric_important_field_metric_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.metric_important_field
    ADD CONSTRAINT fk_metric_important_field_metric_id FOREIGN KEY (metric_id) REFERENCES public.metric(id) ON DELETE CASCADE;


--
-- Name: metric fk_metric_ref_creator_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.metric
    ADD CONSTRAINT fk_metric_ref_creator_id FOREIGN KEY (creator_id) REFERENCES public.core_user(id) ON DELETE CASCADE;


--
-- Name: metric fk_metric_ref_table_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.metric
    ADD CONSTRAINT fk_metric_ref_table_id FOREIGN KEY (table_id) REFERENCES public.metabase_table(id) ON DELETE CASCADE;


--
-- Name: model_index fk_model_index_creator_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.model_index
    ADD CONSTRAINT fk_model_index_creator_id FOREIGN KEY (creator_id) REFERENCES public.core_user(id) ON DELETE CASCADE;


--
-- Name: model_index fk_model_index_model_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.model_index
    ADD CONSTRAINT fk_model_index_model_id FOREIGN KEY (model_id) REFERENCES public.report_card(id) ON DELETE CASCADE;


--
-- Name: model_index_value fk_model_index_value_model_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.model_index_value
    ADD CONSTRAINT fk_model_index_value_model_id FOREIGN KEY (model_index_id) REFERENCES public.model_index(id) ON DELETE CASCADE;


--
-- Name: parameter_card fk_parameter_card_ref_card_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.parameter_card
    ADD CONSTRAINT fk_parameter_card_ref_card_id FOREIGN KEY (card_id) REFERENCES public.report_card(id) ON DELETE CASCADE;


--
-- Name: permissions_group_membership fk_permissions_group_group_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.permissions_group_membership
    ADD CONSTRAINT fk_permissions_group_group_id FOREIGN KEY (group_id) REFERENCES public.permissions_group(id) ON DELETE CASCADE;


--
-- Name: permissions fk_permissions_group_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.permissions
    ADD CONSTRAINT fk_permissions_group_id FOREIGN KEY (group_id) REFERENCES public.permissions_group(id) ON DELETE CASCADE;


--
-- Name: permissions_group_membership fk_permissions_group_membership_user_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.permissions_group_membership
    ADD CONSTRAINT fk_permissions_group_membership_user_id FOREIGN KEY (user_id) REFERENCES public.core_user(id) ON DELETE CASCADE;


--
-- Name: permissions_revision fk_permissions_revision_user_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.permissions_revision
    ADD CONSTRAINT fk_permissions_revision_user_id FOREIGN KEY (user_id) REFERENCES public.core_user(id) ON DELETE CASCADE;


--
-- Name: persisted_info fk_persisted_info_card_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.persisted_info
    ADD CONSTRAINT fk_persisted_info_card_id FOREIGN KEY (card_id) REFERENCES public.report_card(id) ON DELETE CASCADE;


--
-- Name: persisted_info fk_persisted_info_database_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.persisted_info
    ADD CONSTRAINT fk_persisted_info_database_id FOREIGN KEY (database_id) REFERENCES public.metabase_database(id) ON DELETE CASCADE;


--
-- Name: persisted_info fk_persisted_info_ref_creator_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.persisted_info
    ADD CONSTRAINT fk_persisted_info_ref_creator_id FOREIGN KEY (creator_id) REFERENCES public.core_user(id);


--
-- Name: pulse_card fk_pulse_card_ref_card_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.pulse_card
    ADD CONSTRAINT fk_pulse_card_ref_card_id FOREIGN KEY (card_id) REFERENCES public.report_card(id) ON DELETE CASCADE;


--
-- Name: pulse_card fk_pulse_card_ref_pulse_card_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.pulse_card
    ADD CONSTRAINT fk_pulse_card_ref_pulse_card_id FOREIGN KEY (dashboard_card_id) REFERENCES public.report_dashboardcard(id) ON DELETE CASCADE;


--
-- Name: pulse_card fk_pulse_card_ref_pulse_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.pulse_card
    ADD CONSTRAINT fk_pulse_card_ref_pulse_id FOREIGN KEY (pulse_id) REFERENCES public.pulse(id) ON DELETE CASCADE;


--
-- Name: pulse_channel_recipient fk_pulse_channel_recipient_ref_pulse_channel_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.pulse_channel_recipient
    ADD CONSTRAINT fk_pulse_channel_recipient_ref_pulse_channel_id FOREIGN KEY (pulse_channel_id) REFERENCES public.pulse_channel(id) ON DELETE CASCADE;


--
-- Name: pulse_channel_recipient fk_pulse_channel_recipient_ref_user_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.pulse_channel_recipient
    ADD CONSTRAINT fk_pulse_channel_recipient_ref_user_id FOREIGN KEY (user_id) REFERENCES public.core_user(id) ON DELETE CASCADE;


--
-- Name: pulse_channel fk_pulse_channel_ref_pulse_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.pulse_channel
    ADD CONSTRAINT fk_pulse_channel_ref_pulse_id FOREIGN KEY (pulse_id) REFERENCES public.pulse(id) ON DELETE CASCADE;


--
-- Name: pulse fk_pulse_collection_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.pulse
    ADD CONSTRAINT fk_pulse_collection_id FOREIGN KEY (collection_id) REFERENCES public.collection(id) ON DELETE SET NULL;


--
-- Name: pulse fk_pulse_ref_creator_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.pulse
    ADD CONSTRAINT fk_pulse_ref_creator_id FOREIGN KEY (creator_id) REFERENCES public.core_user(id) ON DELETE CASCADE;


--
-- Name: pulse fk_pulse_ref_dashboard_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.pulse
    ADD CONSTRAINT fk_pulse_ref_dashboard_id FOREIGN KEY (dashboard_id) REFERENCES public.report_dashboard(id) ON DELETE CASCADE;


--
-- Name: qrtz_blob_triggers fk_qrtz_blob_triggers_triggers; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.qrtz_blob_triggers
    ADD CONSTRAINT fk_qrtz_blob_triggers_triggers FOREIGN KEY (sched_name, trigger_name, trigger_group) REFERENCES public.qrtz_triggers(sched_name, trigger_name, trigger_group);


--
-- Name: qrtz_cron_triggers fk_qrtz_cron_triggers_triggers; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.qrtz_cron_triggers
    ADD CONSTRAINT fk_qrtz_cron_triggers_triggers FOREIGN KEY (sched_name, trigger_name, trigger_group) REFERENCES public.qrtz_triggers(sched_name, trigger_name, trigger_group);


--
-- Name: qrtz_simple_triggers fk_qrtz_simple_triggers_triggers; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.qrtz_simple_triggers
    ADD CONSTRAINT fk_qrtz_simple_triggers_triggers FOREIGN KEY (sched_name, trigger_name, trigger_group) REFERENCES public.qrtz_triggers(sched_name, trigger_name, trigger_group);


--
-- Name: qrtz_simprop_triggers fk_qrtz_simprop_triggers_triggers; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.qrtz_simprop_triggers
    ADD CONSTRAINT fk_qrtz_simprop_triggers_triggers FOREIGN KEY (sched_name, trigger_name, trigger_group) REFERENCES public.qrtz_triggers(sched_name, trigger_name, trigger_group);


--
-- Name: qrtz_triggers fk_qrtz_triggers_job_details; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.qrtz_triggers
    ADD CONSTRAINT fk_qrtz_triggers_job_details FOREIGN KEY (sched_name, job_name, job_group) REFERENCES public.qrtz_job_details(sched_name, job_name, job_group);


--
-- Name: query_action fk_query_action_database_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.query_action
    ADD CONSTRAINT fk_query_action_database_id FOREIGN KEY (database_id) REFERENCES public.metabase_database(id) ON DELETE CASCADE;


--
-- Name: query_action fk_query_action_ref_action_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.query_action
    ADD CONSTRAINT fk_query_action_ref_action_id FOREIGN KEY (action_id) REFERENCES public.action(id) ON DELETE CASCADE;


--
-- Name: query_field fk_query_field_card_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.query_field
    ADD CONSTRAINT fk_query_field_card_id FOREIGN KEY (card_id) REFERENCES public.report_card(id) ON DELETE CASCADE;


--
-- Name: query_field fk_query_field_field_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.query_field
    ADD CONSTRAINT fk_query_field_field_id FOREIGN KEY (field_id) REFERENCES public.metabase_field(id) ON DELETE CASCADE;


--
-- Name: recent_views fk_recent_views_ref_user_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.recent_views
    ADD CONSTRAINT fk_recent_views_ref_user_id FOREIGN KEY (user_id) REFERENCES public.core_user(id) ON DELETE CASCADE;


--
-- Name: report_card fk_report_card_ref_database_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.report_card
    ADD CONSTRAINT fk_report_card_ref_database_id FOREIGN KEY (database_id) REFERENCES public.metabase_database(id) ON DELETE CASCADE;


--
-- Name: report_card fk_report_card_ref_table_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.report_card
    ADD CONSTRAINT fk_report_card_ref_table_id FOREIGN KEY (table_id) REFERENCES public.metabase_table(id) ON DELETE CASCADE;


--
-- Name: report_dashboardcard fk_report_dashboardcard_ref_action_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.report_dashboardcard
    ADD CONSTRAINT fk_report_dashboardcard_ref_action_id FOREIGN KEY (action_id) REFERENCES public.action(id) ON DELETE CASCADE;


--
-- Name: report_dashboardcard fk_report_dashboardcard_ref_dashboard_tab_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.report_dashboardcard
    ADD CONSTRAINT fk_report_dashboardcard_ref_dashboard_tab_id FOREIGN KEY (dashboard_tab_id) REFERENCES public.dashboard_tab(id) ON DELETE CASCADE;


--
-- Name: revision fk_revision_ref_user_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.revision
    ADD CONSTRAINT fk_revision_ref_user_id FOREIGN KEY (user_id) REFERENCES public.core_user(id) ON DELETE CASCADE;


--
-- Name: secret fk_secret_ref_user_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.secret
    ADD CONSTRAINT fk_secret_ref_user_id FOREIGN KEY (creator_id) REFERENCES public.core_user(id);


--
-- Name: segment fk_segment_ref_creator_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.segment
    ADD CONSTRAINT fk_segment_ref_creator_id FOREIGN KEY (creator_id) REFERENCES public.core_user(id) ON DELETE CASCADE;


--
-- Name: segment fk_segment_ref_table_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.segment
    ADD CONSTRAINT fk_segment_ref_table_id FOREIGN KEY (table_id) REFERENCES public.metabase_table(id) ON DELETE CASCADE;


--
-- Name: core_session fk_session_ref_user_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.core_session
    ADD CONSTRAINT fk_session_ref_user_id FOREIGN KEY (user_id) REFERENCES public.core_user(id) ON DELETE CASCADE;


--
-- Name: native_query_snippet fk_snippet_collection_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.native_query_snippet
    ADD CONSTRAINT fk_snippet_collection_id FOREIGN KEY (collection_id) REFERENCES public.collection(id) ON DELETE SET NULL;


--
-- Name: native_query_snippet fk_snippet_creator_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.native_query_snippet
    ADD CONSTRAINT fk_snippet_creator_id FOREIGN KEY (creator_id) REFERENCES public.core_user(id) ON DELETE CASCADE;


--
-- Name: table_privileges fk_table_privileges_table_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.table_privileges
    ADD CONSTRAINT fk_table_privileges_table_id FOREIGN KEY (table_id) REFERENCES public.metabase_table(id) ON DELETE CASCADE;


--
-- Name: metabase_table fk_table_ref_database_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.metabase_table
    ADD CONSTRAINT fk_table_ref_database_id FOREIGN KEY (db_id) REFERENCES public.metabase_database(id) ON DELETE CASCADE;


--
-- Name: timeline fk_timeline_collection_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.timeline
    ADD CONSTRAINT fk_timeline_collection_id FOREIGN KEY (collection_id) REFERENCES public.collection(id) ON DELETE CASCADE;


--
-- Name: timeline fk_timeline_creator_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.timeline
    ADD CONSTRAINT fk_timeline_creator_id FOREIGN KEY (creator_id) REFERENCES public.core_user(id) ON DELETE CASCADE;


--
-- Name: user_parameter_value fk_user_parameter_value_user_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.user_parameter_value
    ADD CONSTRAINT fk_user_parameter_value_user_id FOREIGN KEY (user_id) REFERENCES public.core_user(id) ON DELETE CASCADE;


--
-- Name: view_log fk_view_log_ref_user_id; Type: FK CONSTRAINT; Schema: public; Owner: aman
--

ALTER TABLE ONLY public.view_log
    ADD CONSTRAINT fk_view_log_ref_user_id FOREIGN KEY (user_id) REFERENCES public.core_user(id) ON DELETE CASCADE;


--
-- Name: DATABASE metabase; Type: ACL; Schema: -; Owner: postgres
--

GRANT ALL ON DATABASE metabase TO aman WITH GRANT OPTION;


--
-- Name: SCHEMA public; Type: ACL; Schema: -; Owner: pg_database_owner
--

GRANT ALL ON SCHEMA public TO aman;


--
-- PostgreSQL database dump complete
--

--
-- Database "postgres" dump
--

\connect postgres

--
-- PostgreSQL database dump
--

-- Dumped from database version 16.3
-- Dumped by pg_dump version 16.3

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- PostgreSQL database dump complete
--

--
-- PostgreSQL database cluster dump complete
--

