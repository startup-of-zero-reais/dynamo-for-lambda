package expressions

import "github.com/startup-of-zero-reais/dynamo-for-lambda/domain"

const (
	GreaterThan        = domain.Condition("Gt")
	GreaterThanOrEqual = domain.Condition("Ge")
	LessThan           = domain.Condition("Lt")
	LessThanOrEqual    = domain.Condition("Le")
	Equal              = domain.Condition("Eq")
	Between            = domain.Condition("Bt")
	StartsWith         = domain.Condition("Sw")
)
