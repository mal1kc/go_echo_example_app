package common

type EnvNameConst string

const (
	PqSqlEnv EnvNameConst = "APROD_APP_PQSQL_CONN_STR"
	//  connString := `host=localhost port=%s user=%s password=%s dbname=%s sslmode=disable statement_cache_mode=describe pool_max_conns=%s pool_max_conn_idle_time=%s`
)
