-- create messages table
CREATE TABLE messages (
   user_id text,
   bucket int,
   message_id text,
   channel_id text,
   content text,
   PRIMARY KEY ((channel_id, bucket), message_id)
) WITH CLUSTERING ORDER BY (message_id DESC);