package web

type Config struct {
	Port              string `envconfig:"PORT" default:"8080"`
	Secret            string `envconfig:"API_SECRET"`
	TokenHourLifespan int    `envconfig:"TOKEN_HOUR_LIFESPAN"`
}
