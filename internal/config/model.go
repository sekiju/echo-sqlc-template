package config

type Application struct {
	Address     string
	Development bool
}

type Database struct {
	Uri        string
	Migrations string
}

type Smtp struct {
	From     string
	Port     int
	Host     string
	Username string
	Password string
}

type Storage struct {
	Host   string
	Daemon string
}

type Config struct {
	Application *Application
	Database    *Database
	Smtp        *Smtp
	Storage     *Storage
}
