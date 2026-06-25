CREATE TABLE draws (
    draw_id BIGSERIAL PRIMARY KEY,
    draw_date TIMESTAMP NOT NULL
);

CREATE TABLE prizes (
    prize_id BIGSERIAL PRIMARY KEY,
    draw_id BIGINT NOT NULL REFERENCES draws(draw_id),
    prize_type VARCHAR(20) NOT NULL,
    winning_number VARCHAR(6) NOT NULL
);

CREATE INDEX idx_winning_number ON prizes(winning_number);