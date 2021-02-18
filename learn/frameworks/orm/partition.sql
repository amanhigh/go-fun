-- Partition by created date into 12 partitions
CREATE TABLE student
(
    `id`      bigint       NOT NULL AUTO_INCREMENT,
    `name`    varchar(100) NOT NULL DEFAULT '0',
    `class`   varchar(100) NOT NULL,
    `created` datetime     NOT NULL DEFAULT NOW(),
    PRIMARY KEY (`id`, `created`)
) PARTITION BY HASH ( MONTH(created) ) PARTITIONS 12;


insert into student
values (DEFAULT, 'Aman', '10th', DEFAULT);
insert into student
values (DEFAULT, 'Preet', '11th', DEFAULT);
insert into student
values (DEFAULT, 'Singh', '12th', DEFAULT);

select *
from student
where created between DATE_SUB(NOW(), INTERVAL 2 HOUR) and DATE_ADD(NOW(), INTERVAL 2 HOUR);

-- Faster operations when partition column included.
explain
select *
from student
where created between DATE_SUB(NOW(), INTERVAL 2 HOUR) and DATE_ADD(NOW(), INTERVAL 2 HOUR);

DROP TABLE student;