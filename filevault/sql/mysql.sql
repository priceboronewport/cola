drop table if exists file_words;
drop table if exists words;
drop table if exists files;
drop table if exists hashes;

create table hashes (
  hash_id int primary key auto_increment,
  hash    varchar(64) not null unique,
  size    bigint not null
) engine=InnoDB;

create table files (
  file_id   int primary key auto_increment,
  hash_id   int not null,
  path      varchar(500) not null,
  name      varchar(500) not null,
  timestamp datetime not null,
  foreign key(hash_id) references hashes(hash_id),
  unique(hash_id, path, name)
) engine=InnoDB;

create table words (
  word_id  int primary key auto_increment,
  word     varchar(100) not null unique
) engine=InnoDB;

create table file_words (
  file_id int not null,
  word_id int not null,
  primary key(file_id, word_id)
) engine=InnoDB;

create index file_words_word_file on file_words(word_id, file_id);
