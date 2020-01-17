-- Create user database
create database books with encoding 'UTF8' template template0;;
grant all privileges on database books to postgres;

-- Connect to user database
\c books

-- Create table in user database
create table if not exists BookInfo(
	id serial primary key not null,
	title varchar(50),
	author varchar(50)
);

-- Insert test data in user table
insert into BookInfo (author, title)
values('Best Author 1', 'Boook 1');

insert into BookInfo (author, title)
values('Best Author 2', 'Boook 1111');

insert into BookInfo (author, title)
values('Best Author 3', 'Boook 111');
