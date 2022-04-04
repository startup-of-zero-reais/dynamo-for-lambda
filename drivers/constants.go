package drivers

import "github.com/startup-of-zero-reais/dynamo-for-lambda/domain"

const (
	// STANDARD é uma classe de provisionamento do Storage
	STANDARD = domain.TableClass("STANDARD")
	// INFREQUENT_ACCESS é uma classe de provisionamento do Storage
	INFREQUENT_ACCESS = domain.TableClass("INFREQUENT_ACCESS")

	GET    = domain.Action("GET")
	PUT    = domain.Action("PUT")
	QUERY  = domain.Action("QUERY")
	UPDATE = domain.Action("UPDATE")
	DELETE = domain.Action("DELETE")

	prod = domain.Environment("production")
	stg  = domain.Environment("staging")
	dev  = domain.Environment("development")
)
