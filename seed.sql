-- inventory-tracker — development seed data
--
-- Loaded automatically after schema.sql on first `docker compose up`.
-- Everything here is fictional. Do not run this against a real database.
--
-- Two accounts are created, both with the password: password123
--
--   itadmin  — category_group 'IT', can reach the admin-only routes
--   staff    — category_group 'Operations', gets 403 on those routes
--
-- The hash below is a real bcrypt hash at Go's DefaultCost (10), so
-- utils.CheckPasswordHash accepts it and you can log in immediately without
-- registering anything.

INSERT INTO category_group (group_name) VALUES
    ('IT'),
    ('Operations'),
    ('Facilities');

INSERT INTO users (name, category_group, username, password) VALUES
    ('Ada Admin', 'IT',         'itadmin', '$2b$10$3lqeidi8h3ZJVfS1w1U7Wuy5qEyYK56AhfcoVcEqySJOKWYrdeRMG'),
    ('Sam Staff', 'Operations', 'staff',   '$2b$10$3lqeidi8h3ZJVfS1w1U7Wuy5qEyYK56AhfcoVcEqySJOKWYrdeRMG');

INSERT INTO category (category_name, category_group) VALUES
    ('Laptop',       'IT'),
    ('Monitor',      'IT'),
    ('Mobile Phone', 'IT'),
    ('Forklift',     'Operations'),
    ('Office Chair', 'Facilities');

-- Categories are referenced by name rather than by a hardcoded id, so this file
-- stays correct no matter what the identity sequence hands out.
INSERT INTO asset
    (asset_code, asset_name, brand, serial_number, status, location, "user", purchase_date, description, category_id)
VALUES
    ('AST-0001', 'Dell Latitude 5440',   'Dell',    'SN-4410023', 'active',      'HQ Floor 3',  'Ada Admin', '2025-03-14', 'Issued to IT support',
        (SELECT id FROM category WHERE category_name = 'Laptop')),
    ('AST-0002', 'MacBook Pro 14',       'Apple',   'SN-MBP14-77', 'active',     'HQ Floor 2',  'Sam Staff', '2025-07-01', NULL,
        (SELECT id FROM category WHERE category_name = 'Laptop')),
    ('AST-0003', 'Dell UltraSharp U2723','Dell',    'SN-U2723-01', 'active',     'HQ Floor 3',  NULL,        '2024-11-20', 'Paired with AST-0001',
        (SELECT id FROM category WHERE category_name = 'Monitor')),
    ('AST-0004', 'Samsung Galaxy A55',   'Samsung', NULL,          'in repair',  'IT Store',    NULL,        '2025-01-09', 'Screen replacement pending',
        (SELECT id FROM category WHERE category_name = 'Mobile Phone')),
    ('AST-0005', 'Toyota 8FG25 Forklift','Toyota',  'SN-FG25-004', 'active',     'Warehouse A', NULL,        '2023-06-30', 'Annual service due June',
        (SELECT id FROM category WHERE category_name = 'Forklift')),
    ('AST-0006', 'Herman Miller Aeron',  'Herman Miller', NULL,    'retired',    'Storage',     NULL,        '2019-02-15', 'Replaced under refresh programme',
        (SELECT id FROM category WHERE category_name = 'Office Chair'));
