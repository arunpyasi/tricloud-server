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
			"batman47": {"dTqHkMo922o:APA91bEh-FO0FzXmR05xp_txgLUja7KeFoHjFd26nzJwA9UVPpU5NhUDZS36ozSF-h3v1ft_pYqwFWBp6bm7mCcR0r3ZjcR_MN4e7MygCeEDG4R6TT8WW4rx5HBNT7KnheosqYH5JVdp"},
		},
		AppFirebaseKeyFile: "./.meta/gcm/tcloud-42ebf-firebase-adminsdk-ma9t8-d5a2581857.json",
	}
}

func GetConfig() *Config {
	return conf
}
