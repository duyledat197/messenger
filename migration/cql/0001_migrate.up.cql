-- create messages table
CREATE TABLE messages (
   user_id text,
   bucket int,
   message_id text,
   from_id text,
   to_id text,
   content text,
   PRIMARY KEY ((user_id, to_id, bucket), message_id)
) WITH CLUSTERING ORDER BY (message_id DESC);