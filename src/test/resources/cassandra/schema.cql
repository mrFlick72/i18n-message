CREATE KEYSPACE IF NOT EXISTS i18n_messages WITH replication = {'class':'SimpleStrategy','replication_factor':'1'};

CREATE TABLE IF NOT EXISTS i18n_messages.APPLICATION (id varchar primary key, name varchar,defaultLanguage varchar)