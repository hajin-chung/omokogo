CREATE TABLE user (
	id varchar(10) NOT NULL PRIMARY KEY,
	name TEXT NOT NULL,
	password TEXT NOT NULL,
	email TEXT NOT NULL,
	score INTEGER DEFAULT 0 NOT NULL,
	status INTEGER DEFAULT 0 NOT NULL,
	gameId varchar(10)
);

CREATE TABLE game (
	id varchar(10) NOT NULL PRIMARY KEY,
	userId1 varchar(10) NOT NULL,
	userId2 varchar(10) NOT NULL,
	status INTEGER DEFAULT 0 NOT NULL
);

CREATE TABLE stone (
	gameId varchar(10) NOT NULL,
	x INTEGER NOT NULL,
	y INTEGER NOT NULL,
	placedAt INTEGER DEFAULT (cast(strftime('%s', 'now') as int)) NOT NULL
);
