-- 创建用户表
CREATE TABLE Users (
                       user_id INT AUTO_INCREMENT PRIMARY KEY,
                       username VARCHAR(50) NOT NULL,
                       password VARCHAR(50) NOT NULL,
                       email VARCHAR(100) ,
                       phone_number VARCHAR(20)
);

--创建警报表
CREATE TABLE alarm_group (
                             alarm_id  INT PRIMARY KEY NOT NULL AUTO_INCREMENT ,
                             user_id INT NOT NULL,
                             alarm_name VARCHAR(20) NOT NULL,
                             alarm_input VARCHAR(50) NOT NULL,
                             alarm_match_mode VARCHAR(10) NOT NULL,
                             alarm_description VARCHAR(500),
                             email_notification BOOLEAN NOT NULL,
                             alarm_email VARCHAR(50),
                             sms_notification BOOLEAN NOT NULL,
                             alarm_phone VARCHAR(20),
                             FOREIGN KEY (user_id) REFERENCES users(user_id)
)

--创建报警日志表
CREATE TABLE alarm_logs (
                            log_id INT PRIMARY KEY AUTO_INCREMENT,
                            alarm_id INT NOT NULL,
                            log_content VARCHAR(500) NOT NULL,
                            FOREIGN KEY (alarm_id) REFERENCES alarm_group(alarm_id)
)