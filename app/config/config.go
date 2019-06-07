package config

var conf *Config

type Config struct {
	Dev        bool
	DBpath     string
	StatDBpath string

	// per user api tokenkey or something doogle calls it
	FirebaseKeys map[string][]string

	//app/project key to push notification
	AppFirebaseKeyFile string
}

func (c *Config) Update() {

}

func init() {
	// todo get that from arg/Env/config file
	conf = &Config{
		Dev:        true,
		DBpath:     "mybolt.db",
		StatDBpath: "sysstat.db",
		FirebaseKeys: map[string][]string{
			"batman47": {"fegnEF0AXtY:APA91bG4f6R6S0I1vtAkf7ngd0z6Vo3aaUiMnCMpy7pmgDZF0aplQ41tt4F4ww0FRhK1BEkZFnEk1nEa79D0hFeGk5ydYldwjSX67P17a71sbCT9iwiJ5JLmXizEOz9xVGzA9i8Ux3M9"},
		},
		AppFirebaseKeyFile: "./tcloud-42ebf-firebase-adminsdk-ma9t8-d5a2581857.json",
	}
}

func GetConfig() *Config {
	return conf
}
