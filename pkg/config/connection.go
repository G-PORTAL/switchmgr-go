package config

// Connection is a struct that holds the connection information for a switch.
// It is used to connect to the switch via SSH. The connection information is
// used as a parameter for the driver.Connect() method.
type Connection struct {
	// Host is the hostname or IP address of the switch.
	Host string
	// Port is the port number of the switch.
	Port int32

	// Username is the username used to connect to the switch.
	Username string
	// Password is the password used to connect to the switch.
	Password string
}
