package config

var (
	cfgUsername = "root"
	cfgPassword = "itsshawn@007@"
	cfgHost     = "127.0.0.1:3306"
	cfgDB       = "parking"
	redisHost   = "127.0.0.1"
	redisPort   = "6379"
	jwtSecret   = "987654"
)

// Functions to access these variables if needed
func Username() string {
	return cfgUsername
}

func Password() string {
	return cfgPassword
}

func Host() string {
	return cfgHost
}

func DB() string {
	return cfgDB
}
func RedisHost() string {
	return redisHost
}
func RedisPort() string {
	return redisPort
}
func JwtSecret() string {
	return jwtSecret
}

// func TokenPassword() string {
// 	return tokenPassword
// }
