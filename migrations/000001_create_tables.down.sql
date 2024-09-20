-- Drop the tables in reverse order to avoid foreign key constraint issues
DROP TABLE IF EXISTS group_days;
DROP TABLE IF EXISTS teacher_groups;
DROP TABLE IF EXISTS student_groups;
DROP TABLE IF EXISTS groups;
DROP TABLE IF EXISTS users;

-- Drop ENUM types
DROP TYPE IF EXISTS roles;
DROP TYPE IF EXISTS genders;
DROP TYPE IF EXISTS days;
