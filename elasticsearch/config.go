package elasticsearch

// Configuration ...
type Configuration struct {
	Name      string   `yaml:"name"`
	Addresses []string `yaml:"addresses"`
	Username  string   `yaml:"username"`
	Password  string   `yaml:"password"`
}
