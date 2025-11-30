package constants

// Define here all environment variables used in the application

const (
	// EnvAeroSpaceMarksDBPath is the environment variable for the AeroSpace marks database path
	// default: `$HOME/.local/state/aerospace-marks`
	EnvAeroSpaceMarksDBPath string = "AEROSPACE_MARKS_DB_PATH"

	// EnvAeroSpaceMarksLogsPath is the environment variable for the AeroSpace marks logs path
	// default: `/tmp/aerospace-marks.log`
	EnvAeroSpaceMarksLogsPath string = "AEROSPACE_MARKS_LOGS_PATH"

	// EnvAeroSpaceMarksLogsLevel is the environment variable for the AeroSpace marks logs level
	// default: `DISABLED`
	EnvAeroSpaceMarksLogsLevel string = "AEROSPACE_MARKS_LOGS_LEVEL"

	// EnvAeroSpaceSock is the environment variable for the AeroSpace IPC socket path.
	EnvAeroSpaceSock string = "AEROSPACESOCK"
)
