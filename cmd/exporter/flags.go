package updater

type Flags struct {
	SQLLitePath *string
	ETHRPCURL   *string
	FromBlcok   *string
	ToBlock     *string

	LogLevel *string
}
