CREATE TABLE IF NOT EXISTS users
(
    user_id     INT AUTO_INCREMENT
        PRIMARY KEY,
    username    VARCHAR(50)                                                         NULL,
    telephone   VARCHAR(50)                                                         NULL,
    email       VARCHAR(50)                                                         NULL,
    password    VARCHAR(255)                                                        NULL,
    level       ENUM ('ADMIN', 'SUPPORT')                 DEFAULT 'SUPPORT'         NULL,
    role_action ENUM ('ALL', 'INSERT', 'UPDATE', 'QUERY') DEFAULT 'QUERY'           NULL,
    status      ENUM ('ACTIVE', 'DISABLE', 'SUSPEND')     DEFAULT 'ACTIVE'          NULL,
    op_id       INT                                                                 NULL,
    created_at  DATETIME                                  DEFAULT CURRENT_TIMESTAMP NOT NULL,
    -- updated_at  DATETIME                                                            NULL ON UPDATE CURRENT_TIMESTAMP,
    updated_at  DATETIME                                                            NULL,
    deleted_at  DATETIME                                                            NULL
);
CREATE INDEX IF NOT EXISTS users_index
    ON users (level, role_action, status, op_id, deleted_at);

-- 
CREATE TABLE IF NOT EXISTS accounts
(
    account_id INT AUTO_INCREMENT
        PRIMARY KEY,
    balance    DECIMAL(10, 2) DEFAULT 0.00  NOT NULL,
    created_at DATETIME       DEFAULT NOW() NOT NULL,
    deleted_at DATETIME                     NULL
)
AUTO_INCREMENT = 1001;

-- 
CREATE TABLE IF NOT EXISTS pockets
(
    pocket_id   INT AUTO_INCREMENT
        PRIMARY KEY,
    pocket_name VARCHAR(100)                 NULL,
    amount      DECIMAL(10, 2) DEFAULT 0.00  NULL,
    currency    VARCHAR(3)     DEFAULT 'USD' NULL,
    account_id  INT                          NULL,
    created_at  DATETIME       DEFAULT NOW() NULL,
    updated_at  DATETIME                     NULL,
    deleted_at  DATETIME                     NULL
)
AUTO_INCREMENT = 1001;

CREATE INDEX pockets_index
    ON pockets (pocket_name, currency);

-- 
CREATE TABLE IF NOT EXISTS transactions
(
    txn_id         INT AUTO_INCREMENT
        PRIMARY KEY,
    from_pocket_id INT                    NOT NULL,
    to_pocket_id   INT                    NOT NULL,
    amount         DECIMAL(10, 2)         NOT NULL,
    txn_date       DATETIME DEFAULT NOW() NOT NULL
)
    AUTO_INCREMENT = 1001;

CREATE INDEX transactions_index
    ON transactions (from_pocket_id, to_pocket_id);
