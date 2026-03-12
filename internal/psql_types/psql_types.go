package psql_types

type Phone struct {
	CountryCode string
	Number      string
}

type Money struct {
	Amount        int
	CurrencyAlpha string
}
