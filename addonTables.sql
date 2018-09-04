CREATE TABLE articles (
  id serial primary key,
  activity_id uuid,
  user_id uuid,
  title text,
  url text,
  cover_img uuid,
  desciption text,
  created_at timestamp without time zone DEFAULT now() NOT NULL,
  updated_at timestamp without time zone DEFAULT now() NOT NULL
);

ALTER TABLE articles ADD COLUMN content text;
ALTER TABLE articles ADD COLUMN popular boolean DEFAULT false NOT NULL;
ALTER TABLE articles alter column cover_img TYPE text;

ALTER TABLE activities ADD COLUMN popular boolean DEFAULT false NOT NULL;
ALTER TABLE users ADD COLUMN premium boolean DEFAULT false NOT NULL;

ALTER TABLE activities alter column cover_id TYPE text;
ALTER TABLE activities rename column cover_id TO cover_img;

ALTER TABLE activities ADD COLUMN tags text[];
ALTER TABLE articles ADD COLUMN tags text[];

ALTER TABLE activities ADD COLUMN city text;

CREATE TABLE privatelike (
  id serial primary key,
  activity_id uuid,
  user_id uuid,
  target_user_id uuid
);


CREATE TABLE smsverify (
  id serial primary key,
  mobile_num text,
  code text,
  isVerified boolean,
  expired_at timestamp without time zone NOT NULL
);

CREATE TABLE article_collections (
    id uuid DEFAULT uuid_generate_v4() NOT NULL,
    article_id int,
    user_id uuid,
    created_at timestamp without time zone DEFAULT now() NOT NULL
);


ALTER TABLE users ADD COLUMN cover_img text;
ALTER TABLE users ADD COLUMN avatar text;