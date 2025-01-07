package mode

const (
	// LocalDev for local development mode.
	LocalDev = "local-dev"
	// DockerDev for dockerized development mode.
	DockerDev = "docker-dev"
)

var mode = LocalDev

func Set(newMode string) {
	mode = newMode
}

func Get() string {
	return mode
}

func IsLocalDev() bool {
	return Get() == LocalDev
}
