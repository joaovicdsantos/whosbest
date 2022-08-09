create table Users (
    id serial primary key,
    username varchar(50) not null,
    password varchar(120) not null
);

create table Leaderboards (
    id serial primary key,
    title varchar(80) not null,
    description text not null,
    creator integer not null,
    constraint fk_creator foreign key (creator) references Users(id)
);

create table Competitors (
    id serial primary key,
    title varchar(80) not null,
    description text not null,
    imageUrl varchar(120) not null,
    votes integer not null,
    leaderboard integer not null,
    constraint fk_leaderboard foreign key (leaderboard) references Leaderboards(id)
);
