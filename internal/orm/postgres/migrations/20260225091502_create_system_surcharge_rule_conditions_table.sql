-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS system_surcharge_rule_conditions (
    surcharge_id UUID NOT NULL,
    condition_id UUID NOT NULL,
    PRIMARY KEY (surcharge_id, condition_id),
    CONSTRAINT fk_surcharge_rule_conditions_surcharge
        FOREIGN KEY (surcharge_id)
        REFERENCES system_surcharge_rules(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_surcharge_rule_conditions_condition
        FOREIGN KEY (condition_id)
        REFERENCES system_surcharge_conditions(id)
        ON DELETE CASCADE
);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS system_surcharge_rule_conditions;

-- +goose StatementEnd
