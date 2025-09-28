-- 1. Таблица пользователей (соответствует вашей модели User)
CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    fio VARCHAR(100) NOT NULL,
    login VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(10) NOT NULL,
    contacts VARCHAR(100),
    cargo_weight DECIMAL(10,2),
    containers_20ft_count INTEGER DEFAULT 0,
    containers_40ft_count INTEGER DEFAULT 0,
    is_moderator BOOLEAN DEFAULT FALSE
);

-- 2. Таблица кораблей (соответствует вашей модели Ship)
CREATE TABLE ships (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    capacity DECIMAL(10,2),               
    length DECIMAL(10,2),
    width DECIMAL(10,2),
    draft DECIMAL(10,2),
    cranes INTEGER,
    containers INTEGER,
    features TEXT,
    photo_url VARCHAR(500),
    is_active BOOLEAN DEFAULT TRUE
);

-- 3. Таблица заявок (соответствует вашей модели Request)
CREATE TABLE requests (
    id SERIAL PRIMARY KEY,
    status VARCHAR(20) DEFAULT 'черновик',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    user_id INTEGER NOT NULL REFERENCES users(user_id),
    
    formed_at TIMESTAMP NULL,
    completed_at TIMESTAMP NULL,
    moderator_id INTEGER NULL REFERENCES users(user_id),
    
    containers_20ft_count INTEGER DEFAULT 0,
    containers_40ft_count INTEGER DEFAULT 0,
    comment TEXT,
    loading_speed VARCHAR(100),
    loading_time DECIMAL(10,2)
);

-- 4. М-М таблица (соответствует вашей модели ShipInRequest)
CREATE TABLE request_ships (
    request_id INTEGER NOT NULL REFERENCES requests(id),
    ship_id INTEGER NOT NULL REFERENCES ships(id),
    count INTEGER DEFAULT 1 NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (request_id, ship_id)
);

-- 5. Ограничение одной черновой заявки
CREATE UNIQUE INDEX one_draft_request_per_user 
ON requests (user_id) 
WHERE status = 'черновик';

-- 6. Тестовый пользователь
INSERT INTO users (fio, login, password, is_moderator) VALUES 
('Агапова Анна Денисовна', 'login001', 'password', true);

-- 7. Все корабли как в вашем коде
INSERT INTO ships (name, capacity, length, width, draft, cranes, containers, features, photo_url, is_active) VALUES 
('Ever Ace', 23992, 400, 61.53, 17.0, 6, 11996, 'самый большой в мире, двигатель Wartsila 70950 кВт', 'ever-ace.png', true),
('FESCO Diomid', 3108, 195, 32.20, 11.0, 3, 536, 'построен в 2010 г., судно класса Ice1 (для Арктики), дизельный двигатель, используется на Северном морском пути', 'fesco-diomid.png', true),
('HMM Algeciras', 23964, 399.9, 61.0, 16.5, 7, 11982, 'двигатель MAN B&W 11G95ME-C9.5 мощностью 64 000 кВт, двойные двигатели, система рекуперации энергии, класс DNV GL', 'hmm-algeciras.png', true),
('MSC Gulsun', 23756, 399.9, 61.4, 16.0, 7, 11878, 'первый в мире контейнеровоз, вмещающий более 23 000 TEU, двигатель MAN B&W 11G95ME-C9.5, класс DNV GL', 'msc-gulsun.png', true);

-- 8. Демо-заявка
INSERT INTO requests (status, user_id, comment) VALUES 
('черновик', 1, 'Демо-заявка для тестирования');