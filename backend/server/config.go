package server

type config struct {
	Stage      string `env:"STAGE" envDefault:"local"`
	BucketURL  string `env:"BUCKET_URL" envDefault:"http://localhost:8080/images"`
	BucketDir  string `env:"BUCKET_DIR" envDefault:"../../public/images"`
	PublicDir  string `env:"PUBLIC_DIR" envDefault:"../../public"`
	DBHost     string `env:"DATABASE_HOST" envDefault:"localhost:3306"`
	DBUser     string `env:"DATABASE_USER" envDefault:"root"`
	DBPassword string `env:"DATABASE_PASSWORD" envDefault:"root"`
	DBName     string `env:"DATABASE_NAME" envSeparator:":" envDefault:"myapps"`
}
