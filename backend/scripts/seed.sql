-- SEED DATA FOR TESTING

-- INSERT TEST USERS
INSERT INTO users (email, password_hash, full_name, role) VALUES
('admin@sustainwear.com', '$2a$10$rN7h8vZKZfZfZfZfZfZfZeJ5vY5vY5vY5vY5vY5vY5vY5vY5vY5vY', 'Admin User', 'admin'),
('donor1@example.com', '$2a$10$rN7h8vZKZfZfZfZfZfZfZeJ5vY5vY5vY5vY5vY5vY5vY5vY5vY5vY', 'John Donor', 'donor'),
('donor2@example.com', '$2a$10$rN7h8vZKZfZfZfZfZfZfZeJ5vY5vY5vY5vY5vY5vY5vY5vY5vY5vY', 'Jane Smith', 'donor'),
('charity@redcross.org', '$2a$10$rN7h8vZKZfZfZfZfZfZfZeJ5vY5vY5vY5vY5vY5vY5vY5vY5vY5vY', 'Red Cross Staff', 'charity_staff');

-- INSERT TEST ORGANISATIONS
INSERT INTO organisations (name, description, type, email, phone, address, city, county, postcode, website, status) VALUES
('British Red Cross', 'Leading humanitarian charity providing emergency response and community support', 'charity', 'info@redcross.org.uk', '07711222333', '44 Moorfields', 'London', 'Greater London', 'EC2Y 9AL', 'https://www.redcross.org.uk', 'active'),
('Oxfam UK', 'Global organisation working to end poverty and injustice', 'charity', 'info@oxfam.org.uk', '07722333444', 'John Smith Drive', 'Oxford', 'Oxfordshire', 'OX4 2JY', 'https://www.oxfam.org.uk', 'active'),
('Shelter', 'Housing and homelessness charity', 'charity', 'info@shelter.org.uk', '07733444555', '88 Old Street', 'London', 'Greater London', 'EC1V 9HU', 'https://england.shelter.org.uk', 'active');

-- INSERT TEST DONATIONS
INSERT INTO donations (donor_id, org_id, item_name, description, category, size, gender, condition, quantity, status) VALUES
(2, 1, 'Winter Coat', 'Navy blue winter coat with hood', 'Outerwear', 'L', 'Unisex', 'Good', 1, 'approved'),
(2, 1, 'Jeans', 'Blue denim jeans', 'Bottoms', '32', 'Male', 'Excellent', 2, 'approved'),
(3, 2, 'T-Shirts', 'Cotton t-shirts in various colors', 'Tops', 'M', 'Unisex', 'Good', 5, 'pending'),
(3, 1, 'Running Shoes', 'Nike running shoes, barely worn', 'Footwear', '9', 'Male', 'Excellent', 1, 'approved');

-- INSERT TEST INVENTORY
INSERT INTO inventory (donation_id, org_id, item_name, category, condition, quantity, available_qty, allocated_qty, distributed_qty, location, status) VALUES
(1, 1, 'Winter Coat', 'Outerwear', 'Good', 1, 1, 0, 0, 'Warehouse A', 'available'),
(2, 1, 'Jeans', 'Bottoms', 'Excellent', 2, 2, 0, 0, 'Warehouse A', 'available'),
(4, 1, 'Running Shoes', 'Footwear', 'Excellent', 1, 0, 1, 0, 'Warehouse B', 'allocated');