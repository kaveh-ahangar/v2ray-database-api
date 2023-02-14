package commands

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
	"v2ray-database-api/config"
	"v2ray-database-api/log"
)

var appConfig *config.Application

func init() {
	cobra.OnInitialize(
		InitAppConfig,
		initLogging,
		logAppConfig,
	)
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.CompletionOptions.DisableNoDescFlag = true
	rootCmd.CompletionOptions.DisableDescriptions = true
	rootCmd.DisableFlagsInUseLine = true
}
func stderrPrintLnf(message string, args ...interface{}) error {
	if !strings.HasSuffix(message, "\n") {
		message += "\n"
	}
	_, err := fmt.Fprint(os.Stderr, message)
	return err
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		_ = stderrPrintLnf(err.Error())
		os.Exit(1)
	}
}

// InitAppConfig initializes the application configuration
func InitAppConfig() {
	if ConfigPath != "" {
		persistentOpts.ConfigPath = ConfigPath
	}
	cfg, err := config.LoadApplicationConfig(viper.GetViper(), persistentOpts)
	if err != nil {
		log.Infof("failed to load application config: %+v", err)
		os.Exit(1)
	}
	appConfig = cfg
}

// logAppConfig logs the application configuration
func logAppConfig() {
	//log.Debugf("application config:\n%+v", color.Magenta.Sprint(appConfig.String()))
}

// initLogging initializes the logging
func initLogging() {

	cfg := log.LogrusConfig{
		EnableConsole: (appConfig.Log.FileLocation == "" || appConfig.CliOptions.Verbosity > 0) && !appConfig.Quiet,
		EnableFile:    appConfig.Log.FileLocation != "",
		Level:         appConfig.Log.LevelOpt,
		Structured:    appConfig.Log.Structured,
		FileLocation:  appConfig.Log.FileLocation,
	}
	//  check if enable file create it
	if cfg.EnableFile {
		//  check if file exists
		if _, err := os.Stat(cfg.FileLocation); os.IsNotExist(err) {
			//  create file
			file, err := os.Create(cfg.FileLocation)
			if err != nil {
				fmt.Println(err)
			}
			// add writeable permission
			err = os.Chmod(cfg.FileLocation, 0777)
			if err != nil {
				fmt.Println(err)
			}
			defer file.Close()
		}
	}
	logWrapper := log.NewLogrusLogger(cfg)

	SetLogger(logWrapper)

	// add a structured field to all loggers of dependencies
	//stereoscope.SetLogger(log.LogrusNestedLogger{Logger: logWrapper.Logger.WithField("from-lib", "stereoscope")},)
}
func SetLogger(logger log.Logger) {
	log.Log = logger
}
