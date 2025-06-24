-- 万知文站数据库初始化脚本

-- 创建数据库
CREATE DATABASE IF NOT EXISTS `wz_wenzhan` 
DEFAULT CHARACTER SET utf8mb4 
COLLATE utf8mb4_unicode_ci;

USE `wz_wenzhan`;

-- 执行建表脚本
SOURCE ./001_create_tables.sql;

-- 创建上传目录索引
INSERT INTO `folders` (`name`, `user_id`, `parent_id`) VALUES
('我的文档', 1, NULL),
('工作文档', 1, 1),
('个人笔记', 1, 1),
('图片资源', 1, NULL);

-- 插入示例文档
INSERT INTO `documents` (`title`, `content`, `type`, `status`, `user_id`, `folder_id`, `tags`) VALUES
('欢迎使用万知文站', '这是您的第一个文档，开始您的知识管理之旅吧！', 'note', 2, 1, 1, '欢迎,指南'),
('项目计划模板', '# 项目计划\n\n## 目标\n\n## 时间安排\n\n## 资源分配', 'note', 1, 1, 2, '模板,计划'),
('会议纪要模板', '# 会议纪要\n\n**时间**: \n**参与人员**: \n**议题**: \n**决议**: ', 'note', 1, 1, 2, '模板,会议');

COMMIT;
