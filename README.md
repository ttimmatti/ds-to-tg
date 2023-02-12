# ds-to-tg
parses messages from discord channels and sends updates to telegram channels
  
  
  
Init:  
docker run -d --name=ds_to_tg -e POSTGRES_HOST_AUTH_METHOD=trust -e POSTGRES_USER=<user> -e POSTGRES_PASSWORD=<pass> -e POSTGRES_DB=ds_to_tg -p <port>:5432 postgres:12.5-alpine

dbs:
1. msgs  
create table msgs(channel_id text,msg_id text,timestamp text);  

2. channels  
create table channels(channel_id text not null primary key,name text,tg_channel_id text);


  
Tg:  
1. Add the bot to channel  
2. Type /help  
3. Type /add <channel_id> <name>  

The bot will now send updates from that Discord channel to current tg group/channel
