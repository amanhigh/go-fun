-- Create Sample Student Table

CREATE TABLE
    `student` (
        `id` int(11) NOT NULL AUTO_INCREMENT,
        `name` varchar(255) NOT NULL,
        `age` int(11) NOT NULL,
        `gender` varchar(10) NOT NULL,
        PRIMARY KEY (`id`)
    )