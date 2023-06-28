-- Partition by created date year and month.
CREATE TABLE student
(
    `id`            bigint       NOT NULL AUTO_INCREMENT,
    `name`          varchar(100) NOT NULL DEFAULT '0',
    `class`         varchar(100) NOT NULL,
    `created`       datetime     NOT NULL DEFAULT NOW(),
    `partition_key` bigint       NOT NULL,
    PRIMARY KEY (id, partition_key)
) PARTITION BY LIST (partition_key ) (
    PARTITION p202101 VALUES IN (202101),
    PARTITION p202102 VALUES IN (202102)
    );


-- Data added goes to respective partition (seprate db file).
insert into student
values (DEFAULT, 'Aman', '10th', DEFAULT, 202101);
insert into student
values (DEFAULT, 'Preet', '11th', DEFAULT, 202101);
insert into student
values (DEFAULT, 'Singh', '12th', DEFAULT, 202102);

select *
from student
where partition_key = 202102;

-- Faster operations when partition column included.
explain
select *
from student
where partition_key = 202102;

-- Check Partition Row Count.
SELECT PARTITION_NAME, TABLE_ROWS
FROM INFORMATION_SCHEMA.PARTITIONS
WHERE TABLE_SCHEMA = Database()
  AND TABLE_NAME = 'student';

-- Can Clear Partition if to be reused.
alter table student truncate partition p202101;
-- Can Drop Partition if no longer required.
# alter table student drop partition p202103;

-- Can Add more Partitions
alter table student
    add partition (
        PARTITION p202103 VALUES IN (202103),
        PARTITION p202104 VALUES IN (202104)
        );

-- Data goes to newly added Partition
insert into student
values (DEFAULT, 'Sachin', '12th', DEFAULT, 202103);


DROP TABLE student;