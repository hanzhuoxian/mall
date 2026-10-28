-- Applied only when initializing an empty development MySQL volume.
CREATE DATABASE IF NOT EXISTS mall CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE DATABASE IF NOT EXISTS mall_user CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
GRANT ALL PRIVILEGES ON mall.* TO 'mall'@'%';
GRANT ALL PRIVILEGES ON mall_user.* TO 'mall'@'%';
