-- name: AddSurchargeRuleConditions :exec
INSERT INTO system_surcharge_rule_conditions (surcharge_id, condition_id)
SELECT $1, UNNEST($2::uuid[]);

-- name: DeleteSurchargeRuleConditionsBySurchargeID :exec
DELETE FROM system_surcharge_rule_conditions
WHERE surcharge_id = $1;

-- name: GetConditionIDsBySurchargeID :many
SELECT condition_id
FROM system_surcharge_rule_conditions
WHERE surcharge_id = $1;

