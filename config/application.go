package config

import (
	"errors"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"reflect"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

var ErrApplicationConfigNotFound = fmt.Errorf("application config not found")
var ErrApplicationConfigPermissionDenied = fmt.Errorf("application config permission denied")

type Application struct {
	ConfigPath string         `yaml:",omitempty" json:"configPath"` // the location where the application config was read from (either from -c or discovered while loading)
	Log        logging        `yaml:"log" json:"log" mapstructure:"log"`
	CliOptions CliOnlyOptions `yaml:"-" json:"-"`
	Quiet      bool           `yaml:"quiet" json:"quiet" mapstructure:"quiet"` // -q, indicates to not show any status output to stderr (ETUI or logging UI)
	Database   Database       `yaml:"database" json:"database" mapstructure:"database"`
}
type defaultValueLoader interface {
	loadDefaultValues(*viper.Viper)
}

type parser interface {
	parseConfigValues() error
}

func newApplicationConfig(v *viper.Viper, cliOpts CliOnlyOptions) *Application {
	config := &Application{
		CliOptions: cliOpts,
	}
	config.loadDefaultValues(v)
	return config
}

// LoadApplicationConfig populates the given viper object with application configuration discovered on disk
func LoadApplicationConfig(v *viper.Viper, cliOpts CliOnlyOptions) (*Application, error) {
	// check DEV mode first
	if cliOpts.DevMode {
		//pwd, _ := os.Getwd()
		cliOpts.ConfigPath = "../configs/test-config.yml"
	}

	// the user may not have a config, and this is OK, we can use the default config + default cobra cli values instead
	config := newApplicationConfig(v, cliOpts)
	if err := readConfig(v, cliOpts.ConfigPath); err != nil {
		return nil, err
	}

	if err := v.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("unable to parse config: %w", err)
	}
	config.ConfigPath = v.ConfigFileUsed()

	if err := config.parseConfigValues(); err != nil {
		return nil, fmt.Errorf("invalid application config: %w", err)
	}

	return config, nil
}

// init loads the default configuration values into the viper instance (before the config values are read and parsed).
func (cfg Application) loadDefaultValues(v *viper.Viper) {
	// for each field in the configuration struct, see if the field implements the defaultValueLoader interface and invoke it if it does
	value := reflect.ValueOf(cfg)
	for i := 0; i < value.NumField(); i++ {
		// note: the defaultValueLoader method receiver is NOT a pointer receiver.
		if loadable, ok := value.Field(i).Interface().(defaultValueLoader); ok {
			// the field implements defaultValueLoader, call it
			loadable.loadDefaultValues(v)
		}
	}
}

// build inflates simple config values into native objects (or other complex objects) after the config is fully read in.
func (cfg *Application) parseConfigValues() error {
	if cfg.Quiet {
		cfg.Log.LevelOpt = logrus.PanicLevel
	} else {
		if cfg.CliOptions.Verbosity > 0 {
			// set the log level implicitly
			switch v := cfg.CliOptions.Verbosity; {
			case v == 1:
				cfg.Log.LevelOpt = logrus.InfoLevel
			case v == 2:
				cfg.Log.LevelOpt = logrus.DebugLevel
			case v >= 3:
				cfg.Log.LevelOpt = logrus.TraceLevel
			default:
				cfg.Log.LevelOpt = logrus.WarnLevel
			}
			cfg.Log.Level = strconv.Itoa(int(cfg.Log.LevelOpt))
		} // set the log level explicitly
		if cfg.Log.Level != "" {
			level, err := logrus.ParseLevel(cfg.Log.Level)
			if err != nil {
				return fmt.Errorf("invalid log level: %w", err)
			}
			cfg.Log.LevelOpt = level
		}
		// set output from config

	}

	// for each field in the configuration struct, see if the field implements the parser interface
	// note: the app config is a pointer, so we need to grab the elements explicitly (to traverse the address)
	value := reflect.ValueOf(cfg).Elem()
	for i := 0; i < value.NumField(); i++ {
		// note: since the interface method of parser is a pointer receiver we need to get the value of the field as a pointer.
		if parsable, ok := value.Field(i).Addr().Interface().(parser); ok {
			// the field implements parser, call it
			if err := parsable.parseConfigValues(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (cfg Application) String() string {
	// yaml is pretty human friendly (at least when compared to json)
	appCfgStr, err := yaml.Marshal(&cfg)

	if err != nil {
		return err.Error()
	}

	return string(appCfgStr)
}

// readConfig attempts to read the given config path from disk or discover an alternate store location
// nolint:funlen
func readConfig(v *viper.Viper, configPath string) error {
	var err error
	v.AutomaticEnv()
	v.SetEnvPrefix(ApplicationName)
	// allow for nested options to be specified via environment variables
	// e.g. pod.context = APPNAME_POD_CONTEXT
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	// use explicitly the given user config
	if configPath != "" {
		v.SetConfigFile(configPath)
		if err := v.ReadInConfig(); err != nil {
			return fmt.Errorf("unable to read application config=%q : %w", configPath, err)
		}
		// don't fall through to other options if the config path was explicitly provided
		return nil
	}

	// start searching for valid configs in order...

	// 1. look for /etc/<appname>/<*>.yml (in the current directory)
	v.AddConfigPath("/etc/" + ApplicationName)
	v.SetConfigName(ApplicationName)
	if err = v.MergeInConfig(); err == nil {
		return nil
	}
	// 2. look for ~/.<appname>.yaml
	home, err := homedir.Dir()
	if err == nil {
		v.AddConfigPath(home)
		v.SetConfigName("." + ApplicationName)
		if err = v.ReadInConfig(); err == nil {
			return nil
		} else if !errors.As(err, &viper.ConfigFileNotFoundError{}) {
			return fmt.Errorf("unable to parse config=%q: %w", v.ConfigFileUsed(), err)
		}
	}

	return ErrApplicationConfigNotFound
}
