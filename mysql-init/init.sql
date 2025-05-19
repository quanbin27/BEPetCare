-- Create databases
CREATE DATABASE IF NOT EXISTS userservice;
CREATE DATABASE IF NOT EXISTS orderservice;
CREATE DATABASE IF NOT EXISTS productservice;
CREATE DATABASE IF NOT EXISTS appointmentservice;
CREATE DATABASE IF NOT EXISTS notificationservice;
CREATE DATABASE IF NOT EXISTS paymentservice;

-- Schema: appointmentservice
USE appointmentservice;

CREATE TABLE appointments (
                              id INT AUTO_INCREMENT PRIMARY KEY,
                              customer_id INT NOT NULL,
                              employee_id INT NOT NULL,
                              customer_address TEXT NOT NULL,
                              scheduled_time DATETIME(3) NOT NULL,
                              status VARCHAR(20) DEFAULT 'pending' NOT NULL,
                              created_at DATETIME(3) NULL,
                              updated_at DATETIME(3) NULL,
                              total FLOAT NOT NULL,
                              note LONGTEXT NULL,
                              branch_id INT NULL
);

CREATE INDEX idx_appointments_branch_id ON appointments (branch_id);
CREATE INDEX idx_appointments_customer_id ON appointments (customer_id);
CREATE INDEX idx_appointments_employee_id ON appointments (employee_id);

CREATE TABLE services (
                          id INT AUTO_INCREMENT PRIMARY KEY,
                          name LONGTEXT NOT NULL,
                          description TEXT NULL,
                          price FLOAT NOT NULL,
                          created_at DATETIME(3) NULL,
                          updated_at DATETIME(3) NULL,
                          img_url TEXT NULL
);

CREATE TABLE appointment_details (
                                     appointment_id INT NOT NULL,
                                     service_id INT NOT NULL,
                                     service_price FLOAT NOT NULL,
                                     quantity INT DEFAULT 0 NOT NULL,
                                     PRIMARY KEY (appointment_id, service_id),
                                     CONSTRAINT fk_appointment_details_service FOREIGN KEY (service_id) REFERENCES services (id),
                                     CONSTRAINT fk_appointments_details FOREIGN KEY (appointment_id) REFERENCES appointments (id)
);

-- Insert data for appointmentservice
INSERT INTO services (name, description, price, created_at, updated_at, img_url) VALUES
                                                                                     ('Grooming', 'Pet grooming service including bath and haircut', 50.0, NOW(), NOW(), 'https://images.unsplash.com/photo-1583512603806-6cc94e815e26'),
                                                                                     ('Vaccination', 'Basic vaccination for pets', 30.0, NOW(), NOW(), 'https://images.pexels.com/photos/6234608/pexels-photo-6234608.jpeg'),
                                                                                     ('Health Check', 'Comprehensive health check-up', 40.0, NOW(), NOW(), 'https://images.unsplash.com/photo-1583511666445-775f1f2116f5'),
                                                                                     ('Dental Cleaning', 'Professional dental cleaning for pets', 60.0, NOW(), NOW(), 'https://images.pexels.com/photos/6234610/pexels-photo-6234610.jpeg'),
                                                                                     ('Nail Trimming', 'Safe and quick nail trimming', 15.0, NOW(), NOW(), 'https://images.unsplash.com/photo-1583511655826-05700d52f4d9'),
                                                                                     ('Flea Treatment', 'Flea and tick prevention treatment', 35.0, NOW(), NOW(), 'https://images.pexels.com/photos/6234609/pexels-photo-6234609.jpeg'),
                                                                                     ('Pet Massage', 'Relaxing massage for pets', 45.0, NOW(), NOW(), 'https://images.unsplash.com/photo-1583512603868-980b467aa0fd'),
                                                                                     ('Ear Cleaning', 'Gentle ear cleaning service', 20.0, NOW(), NOW(), 'https://images.pexels.com/photos/6234611/pexels-photo-6234611.jpeg'),
                                                                                     ('Behavioral Training', 'Basic obedience training session', 70.0, NOW(), NOW(), 'https://images.unsplash.com/photo-1583511666372-33739043e7d8'),
                                                                                     ('Pet Spa Package', 'Luxury spa treatment for pets', 100.0, NOW(), NOW(), 'https://images.unsplash.com/photo-1583512603805-3cc6b41f3edb');

INSERT INTO appointments (customer_id, employee_id, customer_address, scheduled_time, status, created_at, updated_at, total, note, branch_id) VALUES
                                                                                                                                                  (1, 2, '123 Pet Street, Hanoi', '2025-05-12 10:00:00', 'pending', NOW(), NOW(), 80.0, 'Bring pet food', 1),
                                                                                                                                                  (1, 2, '456 Dog Road, Hanoi', '2025-05-13 14:00:00', 'confirmed', NOW(), NOW(), 50.0, NULL, 1),
                                                                                                                                                  (2, 2, '789 Cat Lane, HCMC', '2025-05-14 09:00:00', 'pending', NOW(), NOW(), 115.0, 'Pet is sensitive to water', 2),
                                                                                                                                                  (3, 2, '101 Puppy Avenue, Hanoi', '2025-05-15 11:00:00', 'confirmed', NOW(), NOW(), 95.0, NULL, 1),
                                                                                                                                                  (1, 2, '222 Kitten Road, HCMC', '2025-05-16 15:00:00', 'pending', NOW(), NOW(), 150.0, 'Luxury spa package requested', 2);

INSERT INTO appointment_details (appointment_id, service_id, service_price, quantity) VALUES
                                                                                          (1, 1, 50.0, 1),
                                                                                          (1, 2, 30.0, 1),
                                                                                          (2, 1, 50.0, 1),
                                                                                          (3, 3, 40.0, 1),
                                                                                          (3, 6, 35.0, 2),
                                                                                          (4, 4, 60.0, 1),
                                                                                          (4, 5, 15.0, 2),
                                                                                          (5, 7, 45.0, 1),
                                                                                          (5, 9, 70.0, 1),
                                                                                          (5, 10, 100.0, 1);

-- Schema: notificationservice
USE notificationservice;

CREATE TABLE email_notifications (
                                     id VARCHAR(36) NOT NULL PRIMARY KEY,
                                     email VARCHAR(255) NOT NULL,
                                     subject VARCHAR(255) NOT NULL,
                                     body TEXT NOT NULL,
                                     created_at DATETIME(3) NOT NULL,
                                     sent_at DATETIME(3) NULL,
                                     status VARCHAR(50) NOT NULL,
                                     updated_at DATETIME(3) NULL
);

-- Insert data for notificationservice
INSERT INTO email_notifications (id, email, subject, body, created_at, sent_at, status, updated_at) VALUES
                                                                                                        ('uuid-1', 'customer1@example.com', 'Appointment Confirmation', 'Your appointment on 2025-05-12 is confirmed.', NOW(), NULL, 'pending', NOW()),
                                                                                                        ('uuid-2', 'customer1@example.com', 'Order Shipped', 'Your order #1 has been shipped.', NOW(), NOW(), 'sent', NOW()),
                                                                                                        ('uuid-3', 'customer2@example.com', 'Payment Received', 'Payment for order #2 received.', NOW(), NOW(), 'sent', NOW());

-- Schema: orderservice
USE orderservice;

CREATE TABLE orders (
                        id INT AUTO_INCREMENT PRIMARY KEY,
                        customer_id INT NOT NULL,
                        branch_id INT NOT NULL,
                        appointment_id INT NULL,
                        total_price FLOAT NOT NULL,
                        status LONGTEXT NOT NULL,
                        created_at DATETIME(3) NULL,
                        updated_at DATETIME(3) NULL
);

CREATE TABLE order_items (
                             order_id INT NOT NULL,
                             product_id INT NOT NULL,
                             product_type VARCHAR(191) NOT NULL,
                             quantity INT NOT NULL,
                             unit_price FLOAT NOT NULL,
                             PRIMARY KEY (order_id, product_id, product_type),
                             CONSTRAINT fk_orders_items FOREIGN KEY (order_id) REFERENCES orders (id)
);

CREATE INDEX idx_orders_branch_id ON orders (branch_id);
CREATE INDEX idx_orders_customer_id ON orders (customer_id);

-- Insert data for orderservice
INSERT INTO orders (customer_id, branch_id, appointment_id, total_price, status, created_at, updated_at) VALUES
                                                                                                             (1, 1, 1, 150.0, 'pending', NOW(), NOW()),
                                                                                                             (1, 1, NULL, 80.0, 'shipped', NOW(), NOW());

INSERT INTO order_items (order_id, product_id, product_type, quantity, unit_price) VALUES
                                                                                       (1, 1, 'food', 2, 50.0),
                                                                                       (1, 2, 'accessory', 1, 50.0),
                                                                                       (2, 3, 'medicine', 1, 80.0);

-- Schema: paymentservice
USE paymentservice;

CREATE TABLE payments (
                          id INT AUTO_INCREMENT PRIMARY KEY,
                          order_id INT NULL,
                          appointment_id INT NULL,
                          amount FLOAT NOT NULL,
                          description TEXT NULL,
                          status VARCHAR(20) NOT NULL,
                          method VARCHAR(20) NOT NULL,
                          created_at DATETIME(3) NULL,
                          updated_at DATETIME(3) NULL
);

-- Insert data for paymentservice
INSERT INTO payments (order_id, appointment_id, amount, description, status, method, created_at, updated_at) VALUES
                                                                                                                 (1, NULL, 150.0, 'Payment for order #1', 'completed', 'credit_card', NOW(), NOW()),
                                                                                                                 (NULL, 1, 80.0, 'Payment for appointment #1', 'completed', 'cash', NOW(), NOW()),
                                                                                                                 (2, NULL, 80.0, 'Payment for order #2', 'pending', 'bank_transfer', NOW(), NOW());

-- Schema: productservice
USE productservice;

CREATE TABLE accessories (
                             id INT AUTO_INCREMENT PRIMARY KEY,
                             name VARCHAR(255) NOT NULL,
                             description VARCHAR(500) NULL,
                             price FLOAT NOT NULL,
                             created_at DATETIME(3) NULL,
                             updated_at DATETIME(3) NULL,
                             is_attachable TINYINT(1) DEFAULT 0 NULL,
                             img_url LONGTEXT NULL
);

CREATE TABLE branch_products (
                                 branch_id INT NOT NULL,
                                 product_id INT NOT NULL,
                                 product_type VARCHAR(50) NOT NULL,
                                 stock_quantity INT NOT NULL,
                                 PRIMARY KEY (branch_id, product_id)
);

CREATE TABLE branches (
                          id INT AUTO_INCREMENT PRIMARY KEY,
                          name VARCHAR(255) NOT NULL,
                          location VARCHAR(500) NULL
);

CREATE TABLE foods (
                       id INT AUTO_INCREMENT PRIMARY KEY,
                       name VARCHAR(255) NOT NULL,
                       description VARCHAR(500) NULL,
                       price FLOAT NOT NULL,
                       created_at DATETIME(3) NULL,
                       updated_at DATETIME(3) NULL,
                       is_attachable TINYINT(1) DEFAULT 0 NULL,
                       img_url LONGTEXT NULL
);

CREATE TABLE medicines (
                           id INT AUTO_INCREMENT PRIMARY KEY,
                           name VARCHAR(255) NOT NULL,
                           description VARCHAR(500) NULL,
                           price FLOAT NOT NULL,
                           created_at DATETIME(3) NULL,
                           updated_at DATETIME(3) NULL,
                           is_attachable TINYINT(1) DEFAULT 0 NULL,
                           img_url LONGTEXT NULL
);

-- Insert data for productservice
INSERT INTO branches (name, location) VALUES
                                          ('Hanoi Branch', '123 Main Street, Hanoi'),
                                          ('HCMC Branch', '456 Central Road, HCMC'),
                                          ('Da Nang Branch', '789 Beach Road, Da Nang'),
                                          ('Hai Phong Branch', '101 Harbor Street, Hai Phong');

INSERT INTO foods (name, description, price, created_at, updated_at, is_attachable, img_url) VALUES
                                                                                                 ('Premium Dog Food', 'High-quality dog food', 50.0, NOW(), NOW(), 1, 'https://images.pexels.com/photos/11265036/pexels-photo-11265036.jpeg'),
                                                                                                 ('Cat Food', 'Nutritious cat food', 40.0, NOW(), NOW(), 1, 'https://images.unsplash.com/photo-1592194993150-7d7f1f9e7c24'),
                                                                                                 ('Puppy Formula', 'Special formula for puppies', 45.0, NOW(), NOW(), 1, 'https://images.pexels.com/photos/11265037/pexels-photo-11265037.jpeg'),
                                                                                                 ('Senior Cat Food', 'Food for senior cats', 42.0, NOW(), NOW(), 1, 'https://images.unsplash.com/photo-1592194993150-7d7f1f9e7c24'),
                                                                                                 ('Grain-Free Dog Food', 'Grain-free diet for dogs', 55.0, NOW(), NOW(), 1, 'https://images.pexels.com/photos/11265038/pexels-photo-11265038.jpeg'),
                                                                                                 ('Kitten Food', 'Nutrient-rich kitten food', 38.0, NOW(), NOW(), 1, 'https://images.unsplash.com/photo-1592194993150-7d7f1f9e7c24'),
                                                                                                 ('Organic Pet Food', 'Organic food for pets', 60.0, NOW(), NOW(), 1, 'https://images.pexels.com/photos/11265039/pexels-photo-11265039.jpeg'),
                                                                                                 ('Low-Fat Dog Food', 'Low-fat diet for overweight dogs', 48.0, NOW(), NOW(), 1, 'https://images.pexels.com/photos/11265040/pexels-photo-11265040.jpeg');

INSERT INTO accessories (name, description, price, created_at, updated_at, is_attachable, img_url) VALUES
                                                                                                       ('Dog Collar', 'Stylish dog collar', 20.0, NOW(), NOW(), 1, 'https://images.unsplash.com/photo-1583511655857-d19b37a5a71a'),
                                                                                                       ('Pet Bed', 'Comfortable pet bed', 50.0, NOW(), NOW(), 0, 'https://images.pexels.com/photos/7210294/pexels-photo-7210294.jpeg'),
                                                                                                       ('Cat Scratching Post', 'Durable scratching post for cats', 35.0, NOW(), NOW(), 0, 'https://images.unsplash.com/photo-1583512603867-66b6b1f1f718'),
                                                                                                       ('Dog Leash', 'Strong leash for dogs', 15.0, NOW(), NOW(), 1, 'https://images.pexels.com/photos/7210295/pexels-photo-7210295.jpeg'),
                                                                                                       ('Pet Carrier', 'Portable pet carrier', 60.0, NOW(), NOW(), 0, 'https://images.unsplash.com/photo-1583511666407-9f07c0b9054e'),
                                                                                                       ('Chew Toy', 'Safe chew toy for dogs', 10.0, NOW(), NOW(), 1, 'https://images.pexels.com/photos/7210296/pexels-photo-7210296.jpeg'),
                                                                                                       ('Cat Tunnel', 'Fun tunnel for cats', 25.0, NOW(), NOW(), 0, 'https://images.unsplash.com/photo-1583512603867-66b6b1f1f718'),
                                                                                                       ('Grooming Brush', 'Brush for pet grooming', 18.0, NOW(), NOW(), 1, 'https://images.pexels.com/photos/7210297/pexels-photo-7210297.jpeg');

INSERT INTO medicines (name, description, price, created_at, updated_at, is_attachable, img_url) VALUES
                                                                                                     ('Flea Treatment', 'Flea and tick treatment', 80.0, NOW(), NOW(), 1, 'https://images.pexels.com/photos/6234609/pexels-photo-6234609.jpeg'),
                                                                                                     ('Pain Relief', 'Pain relief for pets', 60.0, NOW(), NOW(), 0, 'https://images.pexels.com/photos/6234612/pexels-photo-6234612.jpeg'),
                                                                                                     ('Antibiotics', 'Broad-spectrum antibiotics', 70.0, NOW(), NOW(), 0, 'https://images.pexels.com/photos/6234613/pexels-photo-6234613.jpeg'),
                                                                                                     ('Dewormer', 'Deworming treatment for pets', 25.0, NOW(), NOW(), 1, 'https://images.pexels.com/photos/6234614/pexels-photo-6234614.jpeg'),
                                                                                                     ('Skin Ointment', 'Ointment for skin issues', 30.0, NOW(), NOW(), 1, 'https://images.pexels.com/photos/6234615/pexels-photo-6234615.jpeg'),
                                                                                                     ('Eye Drops', 'Eye drops for pets', 20.0, NOW(), NOW(), 1, 'https://images.pexels.com/photos/6234616/pexels-photo-6234616.jpeg'),
                                                                                                     ('Heartworm Prevention', 'Monthly heartworm prevention', 50.0, NOW(), NOW(), 1, 'https://images.pexels.com/photos/6234617/pexels-photo-6234617.jpeg'),
                                                                                                     ('Anti-Itch Spray', 'Spray for itching relief', 28.0, NOW(), NOW(), 1, 'https://images.pexels.com/photos/6234618/pexels-photo-6234618.jpeg');

-- Insert data for branch_products with unique product_id
INSERT INTO branch_products (branch_id, product_id, product_type, stock_quantity) VALUES
                                                                                      (1, 1, 'food', 100),        -- Premium Dog Food
                                                                                      (1, 2, 'accessory', 50),    -- Dog Collar
                                                                                      (1, 3, 'medicine', 20),     -- Flea Treatment
                                                                                      (2, 4, 'food', 80),         -- Cat Food
                                                                                      (2, 5, 'accessory', 30),    -- Pet Bed
                                                                                      (2, 6, 'medicine', 15),     -- Pain Relief
                                                                                      (3, 7, 'food', 60),         -- Puppy Formula
                                                                                      (3, 8, 'accessory', 25),    -- Cat Scratching Post
                                                                                      (3, 9, 'medicine', 10),     -- Antibiotics
                                                                                      (4, 10, 'food', 50),        -- Senior Cat Food
                                                                                      (4, 11, 'accessory', 20),   -- Dog Leash
                                                                                      (4, 12, 'medicine', 8);     -- Dewormer

-- Schema: userservice
USE userservice;

CREATE TABLE employee_branches (
                                   user_id INT NOT NULL,
                                   branch_id INT NOT NULL,
                                   PRIMARY KEY (user_id, branch_id)
);

CREATE TABLE roles (
                       id INT AUTO_INCREMENT PRIMARY KEY,
                       name VARCHAR(191) NOT NULL,
                       CONSTRAINT uni_roles_name UNIQUE (name)
);

CREATE TABLE user_roles (
                            user_id INT NOT NULL,
                            role_id INT NOT NULL,
                            PRIMARY KEY (user_id, role_id)
);

CREATE TABLE users (
                       id INT AUTO_INCREMENT PRIMARY KEY,
                       email VARCHAR(191) NOT NULL,
                       password LONGTEXT NOT NULL,
                       name LONGTEXT NOT NULL,
                       branch_id INT NULL,
                       created_at DATETIME(3) NULL,
                       phone_number VARCHAR(11) NULL,
                       address LONGTEXT NULL,
                       CONSTRAINT uni_users_email UNIQUE (email)
);

CREATE INDEX idx_users_branch_id ON users (branch_id);

-- Insert data for userservice
INSERT INTO roles (name) VALUES
                             ('customer'),
                             ('employee'),
                             ('admin');

INSERT INTO users (email, password, name, branch_id, created_at, phone_number, address) VALUES
                                                                                            ('customer1@example.com', 'hashed_password_1', 'John Doe', NULL, NOW(), '0123456789', '123 Pet Street, Hanoi'),
                                                                                            ('employee1@example.com', 'hashed_password_2', 'Jane Smith', 1, NOW(), '0987654321', '456 Work Road, Hanoi'),
                                                                                            ('admin1@example.com', 'hashed_password_3', 'Admin User', NULL, NOW(), '0112233445', NULL);

INSERT INTO user_roles (user_id, role_id) VALUES
                                              (1, 1), -- John Doe as customer
                                              (2, 2), -- Jane Smith as employee
                                              (3, 3); -- Admin User as admin

INSERT INTO employee_branches (user_id, branch_id) VALUES
    (2, 1); -- Jane Smith works at Hanoi Branch