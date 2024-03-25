-- Partition by created date year and month.
CREATE TABLE student (
    `id` bigint NOT NULL AUTO_INCREMENT, `name` varchar(100) NOT NULL DEFAULT '0', `class` varchar(100) NOT NULL, `created` datetime NOT NULL DEFAULT NOW(), PRIMARY KEY (id, created)
)
PARTITION BY
    RANGE (
        YEAR(created) * 100 + MONTH(created)
    ) (
        PARTITION p202101
        VALUES
            LESS THAN (202102), PARTITION p202102
        VALUES
            LESS THAN (202103), PARTITION p202103
        VALUES
            LESS THAN (202104), PARTITION pLAST
        VALUES
            LESS THAN MAXVALUE
    );

-- Data added goes to respective partition (seprate db file).
insert into student values ( DEFAULT, 'Aman', '10th', '2021-01-01' );

insert into
    student
values (
        DEFAULT, 'Preet', '11th', '2021-02-01'
    );

insert into
    student
values (
        DEFAULT, 'Singh', '12th', '2021-03-01'
    );

select * from student where created = '2021-01-01';

-- Faster operations when partition column included.
explain select * from student where created = '2021-01-01';

-- Check Partition Row Count.
SELECT PARTITION_NAME, TABLE_ROWS
FROM INFORMATION_SCHEMA.PARTITIONS
WHERE
    TABLE_SCHEMA = Database()
    AND TABLE_NAME = 'student';

-- Can Clear Partition if to be reused.
alter table student truncate partition p202101;
-- Can Drop Partition if no longer required.
# alter table student drop partition p202103;

-- Can Add more Partitions
ALTER TABLE student REORGANIZE PARTITION pLAST INTO (
    PARTITION p202104
    VALUES
        LESS THAN (202105), PARTITION p202105
    VALUES
        LESS THAN (202106), PARTITION pLAST
    VALUES
        LESS THAN MAXVALUE
);

-- Data goes to newly added Partition
insert into
    student
values (
        DEFAULT, 'Sachin', '6th', '2021-05-01'
    );

DROP TABLE student;