--
-- PostgreSQL database dump
--

-- Dumped from database version 9.6.5
-- Dumped by pg_dump version 9.6.5
-- This database relations were designed by a bloody stupid freelancer
-- who successfully f**ked up everything, and I have suffered enough from this shit
-- And finally I have to change every occurance of Owner to ME so that I can bloody use it

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner:  jackey
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner:  jackey
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


--
-- Name: postgis; Type: EXTENSION; Schema: -; Owner:  jackey
--

CREATE EXTENSION IF NOT EXISTS postgis WITH SCHEMA public;


--
-- Name: EXTENSION postgis; Type: COMMENT; Schema: -; Owner:  jackey
--

COMMENT ON EXTENSION postgis IS 'PostGIS geometry, geography, and raster spatial types and functions';


--
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner:  jackey
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;


--
-- Name: EXTENSION "uuid-ossp"; Type: COMMENT; Schema: -; Owner:  jackey
--

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';


SET search_path = public, pg_catalog;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: activities; Type: TABLE; Schema: public; Owner:  jackey 
--

CREATE TABLE activities (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    title character varying(255),
    content character varying(1024),
    start_time timestamp without time zone,
    join_number integer,
    location_name character varying(255),
    status text,
    weight integer DEFAULT 0,
    cover_id uuid,
    category_id uuid,
    creator_id uuid,
    apartment_id uuid,
    latitude double precision,
    longitude double precision,
    location geography(Point,4326),
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL,
    browsing_count bigint DEFAULT 0
);


ALTER TABLE activities OWNER TO "jackey";

--
-- Name: activity_categories; Type: TABLE; Schema: public; Owner:  jackey
--

CREATE TABLE activity_categories (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    name character varying(50),
    "order" integer,
    img_key text,
    parent_id uuid,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE activity_categories OWNER TO "jackey";

--
-- Name: activity_collections; Type: TABLE; Schema: public; Owner:  jackey
--

CREATE TABLE activity_collections (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    activity_id uuid,
    user_id uuid,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE activity_collections OWNER TO "jackey";

--
-- Name: activity_comments; Type: TABLE; Schema: public; Owner:  jackey
--

CREATE TABLE activity_comments (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    content text,
    topic_id uuid,
    from_uid uuid,
    to_uid uuid,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE activity_comments OWNER TO "jackey";

--
-- Name: activity_exit_reasons; Type: TABLE; Schema: public; Owner:  jackey
--

CREATE TABLE activity_exit_reasons (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    activity_id uuid,
    reason_id uuid,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE activity_exit_reasons OWNER TO "jackey";

--
-- Name: activity_kick_histories; Type: TABLE; Schema: public; Owner:  jackey
--

CREATE TABLE activity_kick_histories (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    user_id uuid,
    activity_id uuid,
    kick_user_id uuid,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE activity_kick_histories OWNER TO "jackey";

--
-- Name: activity_members; Type: TABLE; Schema: public; Owner:  jackey
--

CREATE TABLE activity_members (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    activity_id uuid,
    user_id uuid,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE activity_members OWNER TO "jackey";

--
-- Name: apartments; Type: TABLE; Schema: public; Owner:  jackey
--

CREATE TABLE apartments (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    name text,
    address text,
    invitation_code text,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE apartments OWNER TO "jackey";

--
-- Name: auths; Type: TABLE; Schema: public; Owner:  jackey
--

CREATE TABLE auths (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    user_id uuid NOT NULL,
    phone character varying(30) NOT NULL,
    password_hash character varying(255),
    phone_verified_at timestamp without time zone,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE auths OWNER TO "jackey";

--
-- Name: exit_reasons; Type: TABLE; Schema: public; Owner:  jackey
--

CREATE TABLE exit_reasons (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    user_id uuid,
    content character varying(255),
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE exit_reasons OWNER TO "jackey";

--
-- Name: images; Type: TABLE; Schema: public; Owner:  jackey
--

CREATE TABLE images (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    type character varying(50),
    key character varying(1024),
    association_id uuid,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE images OWNER TO "jackey";

--
-- Name: messages; Type: TABLE; Schema: public; Owner:  jackey
--

CREATE TABLE messages (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    activity_id uuid,
    user_id uuid,
    action_id uuid,
    type text,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE messages OWNER TO "jackey";

--
-- Name: personal_feeds; Type: TABLE; Schema: public; Owner:  jackey
--

CREATE TABLE personal_feeds (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    activity_id uuid,
    user_id uuid,
    type text,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE personal_feeds OWNER TO "jackey";

--
-- Name: schema_info; Type: TABLE; Schema: public; Owner:  jackey
--

CREATE TABLE schema_info (
    version integer DEFAULT 0 NOT NULL
);


ALTER TABLE schema_info OWNER TO "jackey";

--
-- Name: sms_captchas; Type: TABLE; Schema: public; Owner:  jackey
--

CREATE TABLE sms_captchas (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    phone text NOT NULL,
    code text NOT NULL,
    expired_at timestamp without time zone NOT NULL,
    verified_at timestamp without time zone,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    record_id uuid
);


ALTER TABLE sms_captchas OWNER TO "jackey";

--
-- Name: sms_records; Type: TABLE; Schema: public; Owner:  jackey
--

CREATE TABLE sms_records (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    platform character varying(255),
    action character varying(255),
    response jsonb,
    content character varying(2000),
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE sms_records OWNER TO "jackey";

--
-- Name: social_accounts; Type: TABLE; Schema: public; Owner:  jackey
--

CREATE TABLE social_accounts (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    user_id uuid,
    identifier text,
    type text,
    raw jsonb,
    profile jsonb,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE social_accounts OWNER TO "jackey";

--
-- Name: tags; Type: TABLE; Schema: public; Owner:  jackey
--

CREATE TABLE tags (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    name text,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE tags OWNER TO "jackey";

--
-- Name: user_followings; Type: TABLE; Schema: public; Owner:  jackey
--

CREATE TABLE user_followings (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    user_id uuid,
    target_id uuid,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE user_followings OWNER TO "jackey";

--
-- Name: user_leanclouds; Type: TABLE; Schema: public; Owner:  jackey
--

CREATE TABLE user_leanclouds (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    user_id uuid,
    object_id text,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE user_leanclouds OWNER TO "jackey";

--
-- Name: user_ratings; Type: TABLE; Schema: public; Owner:  jackey
--

CREATE TABLE user_ratings (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    user_id uuid,
    rating_user_id uuid,
    activity_id uuid,
    type text,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE user_ratings OWNER TO "jackey";

--
-- Name: user_tags; Type: TABLE; Schema: public; Owner:  jackey
--

CREATE TABLE user_tags (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    user_id uuid,
    tag_id uuid,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE user_tags OWNER TO "jackey";

--
-- Name: users; Type: TABLE; Schema: public; Owner:  jackey
--

CREATE TABLE users (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    apartment_id uuid,
    profile jsonb,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE users OWNER TO "jackey";

--
-- Name: activities activities_pkey; Type: CONSTRAINT; Schema: public; Owner:  jackey
--

ALTER TABLE ONLY activities
    ADD CONSTRAINT activities_pkey PRIMARY KEY (id);


--
-- Name: activity_categories activity_categories_pkey; Type: CONSTRAINT; Schema: public; Owner:  jackey
--

ALTER TABLE ONLY activity_categories
    ADD CONSTRAINT activity_categories_pkey PRIMARY KEY (id);


--
-- Name: activity_collections activity_collections_pkey; Type: CONSTRAINT; Schema: public; Owner:  jackey
--

ALTER TABLE ONLY activity_collections
    ADD CONSTRAINT activity_collections_pkey PRIMARY KEY (id);


--
-- Name: activity_comments activity_comments_pkey; Type: CONSTRAINT; Schema: public; Owner:  jackey
--

ALTER TABLE ONLY activity_comments
    ADD CONSTRAINT activity_comments_pkey PRIMARY KEY (id);


--
-- Name: activity_exit_reasons activity_exit_reasons_pkey; Type: CONSTRAINT; Schema: public; Owner:  jackey
--

ALTER TABLE ONLY activity_exit_reasons
    ADD CONSTRAINT activity_exit_reasons_pkey PRIMARY KEY (id);


--
-- Name: activity_kick_histories activity_kick_histories_pkey; Type: CONSTRAINT; Schema: public; Owner:  jackey
--

ALTER TABLE ONLY activity_kick_histories
    ADD CONSTRAINT activity_kick_histories_pkey PRIMARY KEY (id);


--
-- Name: activity_members activity_members_pkey; Type: CONSTRAINT; Schema: public; Owner:  jackey
--

ALTER TABLE ONLY activity_members
    ADD CONSTRAINT activity_members_pkey PRIMARY KEY (id);


--
-- Name: apartments apartments_pkey; Type: CONSTRAINT; Schema: public; Owner:  jackey
--

ALTER TABLE ONLY apartments
    ADD CONSTRAINT apartments_pkey PRIMARY KEY (id);


--
-- Name: auths auths_pkey; Type: CONSTRAINT; Schema: public; Owner:  jackey
--

ALTER TABLE ONLY auths
    ADD CONSTRAINT auths_pkey PRIMARY KEY (id);


--
-- Name: exit_reasons exit_reasons_pkey; Type: CONSTRAINT; Schema: public; Owner:  jackey
--

ALTER TABLE ONLY exit_reasons
    ADD CONSTRAINT exit_reasons_pkey PRIMARY KEY (id);


--
-- Name: images images_pkey; Type: CONSTRAINT; Schema: public; Owner:  jackey
--

ALTER TABLE ONLY images
    ADD CONSTRAINT images_pkey PRIMARY KEY (id);


--
-- Name: messages messages_pkey; Type: CONSTRAINT; Schema: public; Owner:  jackey
--

ALTER TABLE ONLY messages
    ADD CONSTRAINT messages_pkey PRIMARY KEY (id);


--
-- Name: personal_feeds personal_feeds_pkey; Type: CONSTRAINT; Schema: public; Owner:  jackey
--

ALTER TABLE ONLY personal_feeds
    ADD CONSTRAINT personal_feeds_pkey PRIMARY KEY (id);


--
-- Name: sms_captchas sms_captchas_pkey; Type: CONSTRAINT; Schema: public; Owner:  jackey
--

ALTER TABLE ONLY sms_captchas
    ADD CONSTRAINT sms_captchas_pkey PRIMARY KEY (id);


--
-- Name: sms_records sms_records_pkey; Type: CONSTRAINT; Schema: public; Owner:  jackey
--

ALTER TABLE ONLY sms_records
    ADD CONSTRAINT sms_records_pkey PRIMARY KEY (id);


--
-- Name: social_accounts social_accounts_pkey; Type: CONSTRAINT; Schema: public; Owner:  jackey
--

ALTER TABLE ONLY social_accounts
    ADD CONSTRAINT social_accounts_pkey PRIMARY KEY (id);


--
-- Name: tags tags_pkey; Type: CONSTRAINT; Schema: public; Owner:  jackey
--

ALTER TABLE ONLY tags
    ADD CONSTRAINT tags_pkey PRIMARY KEY (id);


--
-- Name: user_followings user_followings_pkey; Type: CONSTRAINT; Schema: public; Owner:  jackey
--

ALTER TABLE ONLY user_followings
    ADD CONSTRAINT user_followings_pkey PRIMARY KEY (id);


--
-- Name: user_leanclouds user_leanclouds_pkey; Type: CONSTRAINT; Schema: public; Owner:  jackey
--

ALTER TABLE ONLY user_leanclouds
    ADD CONSTRAINT user_leanclouds_pkey PRIMARY KEY (id);


--
-- Name: user_ratings user_ratings_pkey; Type: CONSTRAINT; Schema: public; Owner:  jackey
--

ALTER TABLE ONLY user_ratings
    ADD CONSTRAINT user_ratings_pkey PRIMARY KEY (id);


--
-- Name: user_tags user_tags_pkey; Type: CONSTRAINT; Schema: public; Owner:  jackey
--

ALTER TABLE ONLY user_tags
    ADD CONSTRAINT user_tags_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner:  jackey
--

ALTER TABLE ONLY users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: activities_gix; Type: INDEX; Schema: public; Owner:  jackey
--

CREATE INDEX activities_gix ON activities USING gist (location);


--
-- Name: activities_status_category_id_start_time_weight_created_at_inde; Type: INDEX; Schema: public; Owner:  jackey
--

CREATE INDEX activities_status_category_id_start_time_weight_created_at_inde ON activities USING btree (status, category_id, start_time, weight, created_at);


--
-- Name: activity_collections_activity_id_user_id_index; Type: INDEX; Schema: public; Owner:  jackey
--

CREATE INDEX activity_collections_activity_id_user_id_index ON activity_collections USING btree (activity_id, user_id);


--
-- Name: activity_collections_user_id_activity_id_index; Type: INDEX; Schema: public; Owner:  jackey
--

CREATE INDEX activity_collections_user_id_activity_id_index ON activity_collections USING btree (user_id, activity_id);


--
-- Name: activity_comments_topic_id_index; Type: INDEX; Schema: public; Owner:  jackey
--

CREATE INDEX activity_comments_topic_id_index ON activity_comments USING btree (topic_id);


--
-- Name: activity_exit_reasons_activity_id_reason_id_index; Type: INDEX; Schema: public; Owner:  jackey
--

CREATE INDEX activity_exit_reasons_activity_id_reason_id_index ON activity_exit_reasons USING btree (activity_id, reason_id);


--
-- Name: activity_kick_histories_activity_id_kick_user_id_index; Type: INDEX; Schema: public; Owner:  jackey
--

CREATE UNIQUE INDEX activity_kick_histories_activity_id_kick_user_id_index ON activity_kick_histories USING btree (activity_id, kick_user_id);


--
-- Name: activity_members_activity_id_user_id_index; Type: INDEX; Schema: public; Owner:  jackey
--

CREATE INDEX activity_members_activity_id_user_id_index ON activity_members USING btree (activity_id, user_id);


--
-- Name: activity_members_user_id_activity_id_index; Type: INDEX; Schema: public; Owner:  jackey
--

CREATE INDEX activity_members_user_id_activity_id_index ON activity_members USING btree (user_id, activity_id);


--
-- Name: apartments_invitation_code_index; Type: INDEX; Schema: public; Owner:  jackey
--

CREATE INDEX apartments_invitation_code_index ON apartments USING btree (invitation_code);


--
-- Name: auths_phone_index; Type: INDEX; Schema: public; Owner:  jackey
--

CREATE UNIQUE INDEX auths_phone_index ON auths USING btree (phone);


--
-- Name: auths_user_id_index; Type: INDEX; Schema: public; Owner:  jackey
--

CREATE UNIQUE INDEX auths_user_id_index ON auths USING btree (user_id);


--
-- Name: exit_reasons_content_index; Type: INDEX; Schema: public; Owner:  jackey
--

CREATE INDEX exit_reasons_content_index ON exit_reasons USING btree (content);


--
-- Name: images_type_association_id_index; Type: INDEX; Schema: public; Owner:  jackey
--

CREATE INDEX images_type_association_id_index ON images USING btree (type, association_id);


--
-- Name: images_type_key_index; Type: INDEX; Schema: public; Owner:  jackey
--

CREATE INDEX images_type_key_index ON images USING btree (type, key);


--
-- Name: messages_user_id_index; Type: INDEX; Schema: public; Owner:  jackey
--

CREATE INDEX messages_user_id_index ON messages USING btree (user_id);


--
-- Name: personal_feeds_user_id_index; Type: INDEX; Schema: public; Owner:  jackey
--

CREATE INDEX personal_feeds_user_id_index ON personal_feeds USING btree (user_id);


--
-- Name: sms_captchas_phone_index; Type: INDEX; Schema: public; Owner:  jackey
--

CREATE INDEX sms_captchas_phone_index ON sms_captchas USING btree (phone);


--
-- Name: sms_records_created_at_index; Type: INDEX; Schema: public; Owner:  jackey
--

CREATE INDEX sms_records_created_at_index ON sms_records USING btree (created_at);


--
-- Name: social_accounts_identifier_index; Type: INDEX; Schema: public; Owner:  jackey
--

CREATE INDEX social_accounts_identifier_index ON social_accounts USING btree (identifier);


--
-- Name: social_accounts_user_id_index; Type: INDEX; Schema: public; Owner:  jackey
--

CREATE INDEX social_accounts_user_id_index ON social_accounts USING btree (user_id);


--
-- Name: tags_name_index; Type: INDEX; Schema: public; Owner:  jackey
--

CREATE INDEX tags_name_index ON tags USING btree (name);


--
-- Name: user_followings_target_id_index; Type: INDEX; Schema: public; Owner:  jackey
--

CREATE INDEX user_followings_target_id_index ON user_followings USING btree (target_id);


--
-- Name: user_followings_user_id_index; Type: INDEX; Schema: public; Owner:  jackey
--

CREATE INDEX user_followings_user_id_index ON user_followings USING btree (user_id);


--
-- Name: user_leanclouds_user_id_index; Type: INDEX; Schema: public; Owner:  jackey
--

CREATE UNIQUE INDEX user_leanclouds_user_id_index ON user_leanclouds USING btree (user_id);


--
-- Name: user_ratings_user_id_type_index; Type: INDEX; Schema: public; Owner:  jackey
--

CREATE INDEX user_ratings_user_id_type_index ON user_ratings USING btree (user_id, type);


--
-- Name: user_tags_user_id_tag_id_index; Type: INDEX; Schema: public; Owner:  jackey
--

CREATE INDEX user_tags_user_id_tag_id_index ON user_tags USING btree (user_id, tag_id);


--
-- PostgreSQL database dump complete
--

