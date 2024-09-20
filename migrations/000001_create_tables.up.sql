CREATE TYPE roles AS ENUM ('admin','student','teacher','support');
CREATE TYPE genders AS ENUM ('male','female');



CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hh_id VARCHAR(100) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    role roles NOT NULL DEFAULT 'student',
    password_hash VARCHAR(255) NOT NULL,
    profile_image VARCHAR(100),
    phone_number VARCHAR(20) UNIQUE NOT NULL,
    parents_phone_number VARCHAR(20) UNIQUE NOT NULL,
    gender genders NOT NULL,
    date_of_birth DATE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS groups (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL,
    subject_id UUID NOT NULL,
    room VARCHAR(100),
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
)