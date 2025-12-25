CREATE TABLE accounts (
  id BIGSERIAL PRIMARY KEY,
  owner TEXT NOT NULL,
  balance BIGINT NOT NULL DEFAULT 0,
  currency TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

ALTER TABLE accounts ADD CONSTRAINT check_balance_positive CHECK (balance >= 0);

CREATE TABLE transfers (
  id BIGSERIAL PRIMARY KEY,
  from_account_id BIGINT NOT NULL REFERENCES accounts(id),
  to_account_id BIGINT NOT NULL REFERENCES accounts(id),
  amount BIGINT NOT NULL, -- Số tiền chuyển
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);