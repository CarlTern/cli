package vulnerabilities

type IVulnerabilities interface {
	Order(args IOrderArgs) (string, error)
}

type IOrderArgs interface{}
