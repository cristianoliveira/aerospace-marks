package constants

// Define here all environment variables used in the application


const (
  // Environment variables

  // EnvAeroSpaceSock is the environment variable for the AeroSpace socket path
  // default: `/tmp/bobko.aerospace-$USER.sock`
  EnvAeroSpaceSock string = "AEROSPACESOCK"

  // EnvAeroSpaceMarksDbPath is the environment variable for the AeroSpace marks database path
  // default: `$HOME/.local/state/aerospace-marks`
  EnvAeroSpaceMarksDbPath string = "AEROSPACE_MARKS_DB_PATH"

  // EnvAeroSpaceMarksLogsPath is the environment variable for the AeroSpace marks logs path
  // default: `/tmp/aerospace-marks.log`
  EnvAeroSpaceMarksLogsPath string = "AEROSPACE_MARKS_LOGS_PATH"

  // EnvAeroSpaceMarksLogsLevel is the environment variable for the AeroSpace marks logs level
  // default: `DISABLED`
  EnvAeroSpaceMarksLogsLevel string = "AEROSPACE_MARKS_LOGS_LEVEL"
)
