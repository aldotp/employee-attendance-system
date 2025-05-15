package bootstrap

func (b *Bootstrap) BuildRestBootstrap() *Bootstrap {
	// set dependencies
	b.setConfig()
	b.setPostgresDB()
	b.setRestApiRepository()
	b.setLogger()
	b.setJWTToken()
	b.setCache()
	b.SetMinio()
	// b.setGCS()
	// b.setRabbitMQ()

	return b
}
