package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"time"

	"github.com/adelowo/rivertui/config"
	"github.com/adelowo/rivertui/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var (
	// Version describes the version of the current build.
	Version = "dev"

	// Commit describes the commit of the current build.
	Commit = "none"

	// Date describes the date of the current build.
	Date = time.Now().UTC()
)

const (
	defaultConfigFilePath = "config"
	envPrefix             = "RIVERTUI_"
)

func Execute() error {

	cfg := &config.Config{}

	rootCmd := &cobra.Command{
		Use:   "malak",
		Short: `Riverqueue TUI`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {

			if cmd.Use == "version" {
				return nil
			}

			confFile, err := cmd.Flags().GetString("config")
			if err != nil {
				return err
			}

			if err := initializeConfig(cfg, confFile); err != nil {
				return err
			}

			return cfg.Validate()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			poolConfig, err := pgxpool.ParseConfig(cfg.Database.DSN)
			if err != nil {
				return fmt.Errorf("error parsing db config: %w", err)
			}

			dbPool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
			if err != nil {
				return fmt.Errorf("error connecting to db: %w", err)
			}

			if err := dbPool.Ping(context.Background()); err != nil {
				return err
			}

			client, err := river.NewClient(riverpgxv5.New(dbPool), &river.Config{})
			if err != nil {
				return err
			}

			p := tea.NewProgram(tui.New(client))

			_, err = p.Run()
			return err
		},
	}

	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Version: %s\nCommit: %s\nBuild Date: %s\n", Version, Commit, Date.Format(time.RFC3339))
		},
	}

	rootCmd.AddCommand(versionCmd)

	rootCmd.PersistentFlags().StringP("config", "c", defaultConfigFilePath, "Config file. This is in YAML")

	return rootCmd.Execute()
}

func initializeConfig(cfg *config.Config, pathToFile string) error {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	viper.AddConfigPath(filepath.Join(homePath, ".config", defaultConfigFilePath))
	viper.AddConfigPath(pathToFile)
	viper.AddConfigPath(".")

	viper.SetConfigName(defaultConfigFilePath)
	viper.SetConfigType("yml")

	viper.SetEnvPrefix(envPrefix)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	setDefaults()

	bindEnvs(viper.GetViper(), "", cfg)

	return viper.Unmarshal(cfg)
}

func bindEnvs(v *viper.Viper, prefix string, iface interface{}) {
	ifv := reflect.ValueOf(iface)
	ift := reflect.TypeOf(iface)

	if ifv.Kind() == reflect.Ptr {
		ifv = ifv.Elem()
		ift = ift.Elem()
	}

	for i := 0; i < ift.NumField(); i++ {
		fieldv := ifv.Field(i)
		t := ift.Field(i)
		name := t.Name
		tag, ok := t.Tag.Lookup("mapstructure")
		if ok {
			name = tag
		}

		path := name
		if prefix != "" {
			path = prefix + "." + name
		}

		switch fieldv.Kind() {
		case reflect.Struct:
			bindEnvs(v, path, fieldv.Addr().Interface())
		default:
			envKey := strings.ToUpper(strings.ReplaceAll(path, ".", "_"))
			if err := v.BindEnv(path, envPrefix+envKey); err != nil {
				panic(err)
			}
		}
	}
}

func setDefaults() {

}

func getLogger(cfg config.Config) (*zap.Logger, error) {
	switch cfg.Logging.Mode {
	case config.LogModeProd:
		return zap.NewProduction()
	case config.LogModeDev:
		return zap.NewDevelopment()
	default:
		return zap.NewDevelopment()
	}
}
