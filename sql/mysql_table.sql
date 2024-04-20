-- 创建用户表
CREATE TABLE Users (
                       user_id INT AUTO_INCREMENT PRIMARY KEY,
                       username VARCHAR(50) NOT NULL,
                       password VARCHAR(50) NOT NULL,
                       email VARCHAR(100) ,
                       phone_number VARCHAR(20)
);