-- Disable foreign key checks for easier data insertion
SET session_replication_role = 'replica';

-- Drop tables if they exist to ensure a clean start
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS address_locations CASCADE;
DROP TABLE IF EXISTS addresses CASCADE;
DROP TABLE IF EXISTS listings CASCADE;
DROP TABLE IF EXISTS listing_images CASCADE;
DROP TABLE IF EXISTS favorites CASCADE;
DROP TABLE IF EXISTS bookings CASCADE;
DROP TABLE IF EXISTS featured_ads CASCADE;
DROP TABLE IF EXISTS ad_packages CASCADE;
DROP TABLE IF EXISTS search_logs CASCADE;

-- Drop enum types if they exist
DROP TYPE IF EXISTS gender_enum;
DROP TYPE IF EXISTS user_status_enum;
DROP TYPE IF EXISTS user_type_enum;
DROP TYPE IF EXISTS identity_type_enum;
DROP TYPE IF EXISTS role_enum;
DROP TYPE IF EXISTS location_level_enum;
DROP TYPE IF EXISTS property_type_enum;
DROP TYPE IF EXISTS listing_type_enum;
DROP TYPE IF EXISTS listing_status_enum;
DROP TYPE IF EXISTS booking_status_enum;
DROP TYPE IF EXISTS ad_status_enum;
DROP TYPE IF EXISTS token_type_enum;

-- Create a trigger function to automatically update the updated_at timestamp
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create ENUM types
CREATE TYPE gender_enum AS ENUM('male', 'female', 'other');
CREATE TYPE user_status_enum AS ENUM('active', 'inactive', 'suspended', 'deleted');
CREATE TYPE user_type_enum AS ENUM('renter', 'owner');
CREATE TYPE identity_type_enum AS ENUM('personal', 'agent', 'business', 'sublease');
CREATE TYPE role_enum AS ENUM('user', 'admin', 'moderator');
CREATE TYPE location_level_enum AS ENUM('city', 'ward', 'area');
CREATE TYPE property_type_enum AS ENUM('rented_room', 'apartment', 'house', 'villa', 'dormitory', 'office', 'store', 'warehouse', 'land', 'all');
CREATE TYPE listing_type_enum AS ENUM('for_rent', 'for_sale');
CREATE TYPE listing_status_enum AS ENUM('active', 'inactive', 'pending', 'rejected', 'deleted');
CREATE TYPE booking_status_enum AS ENUM('pending', 'accepted', 'rejected', 'cancelled', 'completed');
CREATE TYPE ad_status_enum AS ENUM('pending', 'active', 'expired', 'cancelled');
CREATE TYPE token_type_enum AS ENUM('email_verification', 'password_reset');

--
-- Table structure for table `users`
--
CREATE TABLE users (
    id                     BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    first_name             VARCHAR(50),
    last_name              VARCHAR(50),
    name                   VARCHAR(100),
    email                  VARCHAR(100) NOT NULL,
    phone                  VARCHAR(20),
    date_of_birth          DATE,
    gender                 gender_enum DEFAULT 'other',
    avatar_url             VARCHAR(255),
    bio                    TEXT,
    password_hash          TEXT NOT NULL,
    salt                   VARCHAR(64),
    status                 user_status_enum DEFAULT 'inactive',
    email_verified_at      TIMESTAMP,
    last_login_at          TIMESTAMP,
    reset_password_token   VARCHAR(255),
    reset_password_expires_at TIMESTAMP,
    user_type              user_type_enum DEFAULT 'renter',
    identity_type          identity_type_enum DEFAULT 'personal',
    company_name           VARCHAR(255),
    business_license_number VARCHAR(100),
    agent_license_number   VARCHAR(100),
    verified_owner         BOOLEAN DEFAULT FALSE,
    role                   role_enum DEFAULT 'user',
    created_at             TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at             TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uni_users_email UNIQUE (email)
);
CREATE TRIGGER set_timestamp_users BEFORE UPDATE ON users FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();

--
-- Table structure for table `tokens`
--
CREATE TABLE tokens (
    id                      BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id                 BIGINT NOT NULL,
    token_type              token_type_enum NOT NULL,
    token                   VARCHAR(255) NOT NULL,
    expires_at              TIMESTAMP NOT NULL,
    created_at              TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at              TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
CREATE TRIGGER set_timestamp_tokens BEFORE UPDATE ON tokens FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();

--
-- Cấu trúc bảng `address_locations`
--
CREATE TABLE address_locations (
    id            BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name          VARCHAR(255) NOT NULL,
    level         location_level_enum NOT NULL,
    parent_id     BIGINT,
    FOREIGN KEY (parent_id) REFERENCES address_locations(id)
);

--
-- Table structure for table `addresses`
--
CREATE TABLE addresses (
    id              BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    city_id         BIGINT NOT NULL,
    ward_id         BIGINT NOT NULL,
    area_id         BIGINT,
    street          VARCHAR(100),
    house_number    VARCHAR(50),
    building_name   VARCHAR(100),
    floor_number    INT,
    room_number     VARCHAR(20),
    latitude        DECIMAL(10,8),
    longitude       DECIMAL(10,8),
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (city_id) REFERENCES address_locations(id),
    FOREIGN KEY (ward_id) REFERENCES address_locations(id),
    FOREIGN KEY (area_id) REFERENCES address_locations(id)
);
CREATE TRIGGER set_timestamp_addresses BEFORE UPDATE ON addresses FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();

--
-- Table structure for table `listings`
--
CREATE TABLE listings (
    id                BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    owner_id          BIGINT REFERENCES users(id),
    title             VARCHAR(255) NOT NULL,
    description       TEXT,
    property_type     property_type_enum DEFAULT 'rented_room',
	is_shared		  BOOLEAN DEFAULT FALSE,
    price             DECIMAL(15,2) NOT NULL,
    area_m2           DECIMAL(10,2),
    address_id        BIGINT REFERENCES addresses(id) NOT NULL,
    contact_phone     VARCHAR(20),
    contact_email     VARCHAR(255),
    contact_name      VARCHAR(100),
    num_bedrooms      INT,
    num_bathrooms     INT,
    num_floors        INT,
    has_balcony       BOOLEAN DEFAULT FALSE,
    has_parking       BOOLEAN DEFAULT FALSE,
    amenities         JSON,
    pet_allowed       BOOLEAN DEFAULT FALSE,
    allowed_pet_types JSON,
    latitude          DECIMAL(10,8), -- REMOVE IGNORE
    longitude         DECIMAL(10,8), -- REMOVE IGNORE
    listing_type      listing_type_enum DEFAULT 'for_rent',
    deposit_amount    DECIMAL(15,2),
    status            listing_status_enum DEFAULT 'pending',
    is_featured       BOOLEAN DEFAULT FALSE,
    view_count        INT DEFAULT 0,
    created_at        TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at        TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    published_at      TIMESTAMP,
    expires_at        TIMESTAMP
);
CREATE TRIGGER set_timestamp_listings BEFORE UPDATE ON listings FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();

--
-- Table structure for table `listing_images`
--
CREATE TABLE listing_images (
    id            BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    listing_id    BIGINT NOT NULL,
    image_url     VARCHAR(255) NOT NULL,
    is_main       BOOLEAN DEFAULT FALSE,
    sort_order    INT DEFAULT 0,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (listing_id) REFERENCES listings(id) ON DELETE CASCADE
);
CREATE TRIGGER set_timestamp_listing_images BEFORE UPDATE ON listing_images FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();

--
-- Table structure for table `favorites`
--
CREATE TABLE favorites (
    user_id    BIGINT NOT NULL REFERENCES users(id),
    listing_id BIGINT NOT NULL REFERENCES listings(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, listing_id)
);

--
-- Table structure for table `bookings`
--
CREATE TABLE bookings (
    id            BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    listing_id    BIGINT NOT NULL,
    renter_id     BIGINT NOT NULL,
    scheduled_time TIMESTAMP NOT NULL,
    status        booking_status_enum DEFAULT 'pending',
    message_from_renter TEXT,
    response_message_from_owner TEXT,
    responded_at TIMESTAMP,
    responded_by BIGINT,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (listing_id) REFERENCES listings(id),
    FOREIGN KEY (renter_id) REFERENCES users(id),
    FOREIGN KEY (responded_by) REFERENCES users(id)
);
CREATE TRIGGER set_timestamp_bookings BEFORE UPDATE ON bookings FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();

--
-- Table structure for table `ad_packages`
--
CREATE TABLE ad_packages (
    id                  BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name                VARCHAR(100) NOT NULL,
    duration_days       INT NOT NULL,
    price               DECIMAL(15,2) NOT NULL,
    priority_level      INT DEFAULT 1,
    display_locations   JSON,
    promoted_areas      JSON,
    description         TEXT,
    is_active           BOOLEAN DEFAULT TRUE,
    created_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
CREATE TRIGGER set_timestamp_ad_packages BEFORE UPDATE ON ad_packages FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();

--
-- Table structure for table `featured_ads`
--
CREATE TABLE featured_ads (
    id            BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    listing_id    BIGINT NOT NULL,
    ad_package_id BIGINT,
    start_date    TIMESTAMP NOT NULL,
    end_date      TIMESTAMP NOT NULL,
    status        ad_status_enum DEFAULT 'pending',
    payment_txn_id VARCHAR(255),
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (listing_id) REFERENCES listings(id),
	FOREIGN KEY (ad_package_id) REFERENCES ad_packages(id)
);
CREATE TRIGGER set_timestamp_featured_ads BEFORE UPDATE ON featured_ads FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();

--
-- Table structure for table `search_logs`
--
CREATE TABLE search_logs (
    id                BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id           BIGINT,
    keyword           VARCHAR(255),
    location_id       BIGINT,
    min_price         DECIMAL(15,2),
    max_price         DECIMAL(15,2),
    min_area_m2       DECIMAL(10,2),
    max_area_m2       DECIMAL(10,2),
    num_bedrooms      INT,
    num_bathrooms     INT,
    num_floors        INT,
    property_type     property_type_enum DEFAULT 'all',
    pet_allowed       BOOLEAN,
    has_balcony       BOOLEAN,
    has_parking       BOOLEAN,
    additional_filters JSON,
    search_count      INT DEFAULT 1,
    created_at        TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at        TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (location_id) REFERENCES addresses(id)
);
CREATE TRIGGER set_timestamp_search_logs BEFORE UPDATE ON search_logs FOR EACH ROW EXECUTE PROCEDURE trigger_set_timestamp();

--
-- Chèn dữ liệu cho TP.HCM (cấp thành phố)
--
INSERT INTO address_locations (name, level, parent_id) VALUES ('TP.HCM', 'city', NULL);

INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Sài Gòn', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Tân Định', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Bến Thành', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Cầu Ông Lãnh', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Bàn Cờ', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Xuân Hòa', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Nhiêu Lộc', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Xóm Chiếu', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Khánh Hội', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Vĩnh Hội', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Chợ Quán', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường An Đông', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Chợ Lớn', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Bình Tây', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Bình Tiên', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Bình Phú', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Phú Lâm', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Tân Thuận', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Phú Thuận', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Tân Mỹ', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Tân Hưng', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Chánh Hưng', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Phú Định', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Bình Đông', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Diên Hồng', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Vườn Lài', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Hòa Hưng', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Minh Phụng', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Bình Thới', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Hòa Bình', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Phú Thọ', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Đông Hưng Thuận', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Trung Mỹ Tây', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Tân Thới Hiệp', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Thới An', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường An Phú Đông', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường An Lạc', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Bình Tân', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Tân Tạo', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Bình Trị Đông', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Bình Hưng Hòa', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Gia Định', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Bình Thạnh', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Bình Lợi Trung', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Thạnh Mỹ Tây', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Bình Quới', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Hạnh Thông', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường An Nhơn', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Gò Vấp', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường An Hội Đông', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Thông Tây Hội', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường An Hội Tây', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Đức Nhuận', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Cầu Kiệu', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Phú Nhuận', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Tân Sơn Hòa', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Tân Sơn Nhất', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Tân Hòa', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Bảy Hiền', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Tân Bình', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Tân Sơn', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Tây Thạnh', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Tân Sơn Nhì', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Phú Thọ Hòa', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Tân Phú', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Phú Thạnh', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Hiệp Bình', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Thủ Đức', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Tam Bình', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Linh Xuân', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Tăng Nhơn Phú', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Long Bình', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Long Phước', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Long Trường', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Cát Lái', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Bình Trưng', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Phước Long', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường An Khánh', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Dĩ An', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường An Phú', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Bình Hòa', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Lái Thiêu', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Thuận An', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Thuận Giao', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Thủ Dầu Một', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Phú Lợi', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Tam Thắng', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Rạch Dừa', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Phước Thắng', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Long Hương', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Tam Long', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Phường Phú Mỹ', 'ward', 1);

INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Vĩnh Lộc', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Tân Vĩnh Lộc', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Bình Lợi', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Tân Nhựt', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Bình Chánh', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Hưng Long', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Bình Hưng', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Bình Khánh', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã An Thới Đông', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Cần Giờ', 'ward', 1); 
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Củ Chi', 'ward', 1); 
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Tân An Hội', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Thái Mỹ', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã An Nhơn Tây', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Nhuận Đức', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Phú Hòa Đông', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Bình Mỹ', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Đông Thạnh', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Hóc Môn', 'ward', 1); 
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Xuân Thới Sơn', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Bà Điểm', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Nhà Bè', 'ward', 1); 
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Hiệp Phước', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Thường Tân', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Bắc Tân Uyên', 'ward', 1); 
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Phú Giáo', 'ward', 1); 
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Phước Hòa', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Phước Thành', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã An Long', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Trừ Văn Thố', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Bàu Bàng', 'ward', 1); 
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Long Hòa', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Thanh An', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Dầu Tiếng', 'ward', 1); 
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Minh Thạnh', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Châu Pha', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Long Hải', 'ward', 1); 
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Long Điền', 'ward', 1); 
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Phước Hải', 'ward', 1); 
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Đất Đỏ', 'ward', 1); 
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Nghĩa Thành', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Ngãi Giao', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Kim Long', 'ward', 1); 
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Châu Đức', 'ward', 1); 
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Bình Giã', 'ward', 1); 
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Xuân Sơn', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Hồ Tràm', 'ward', 1); 
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Xuyên Mộc', 'ward', 1); 
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Hòa Hội', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Bàu Lâm', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Đông Hòa', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Tân Đông Hiệp', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Chánh Hiệp', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Bình Dương', 'ward', 1); 
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Hòa Lợi', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Phú An', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Tây Nam', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Long Nguyên', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Bến Cát', 'ward', 1); 
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Chánh Phú Hòa', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Vĩnh Tân', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Bình Cơ', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Tân Uyên', 'ward', 1); 
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Tân Hiệp', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Tân Khánh', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Vũng Tàu', 'ward', 1); 
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Bà Rịa', 'ward', 1); 
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Tân Hải', 'ward', 1);
INSERT INTO address_locations (name, level, parent_id) VALUES ('Xã Tân Phước', 'ward', 1);

INSERT INTO address_locations (name, level, parent_id) VALUES ('Đặc khu Côn Đảo', 'ward', 1);

-- Enable foreign key checks again
SET session_replication_role = 'origin';


