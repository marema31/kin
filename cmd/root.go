package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/marema31/kin/cache"
	"github.com/marema31/kin/collector"
	"github.com/marema31/kin/server"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	baseURL   string
	cfgFile   string
	ctx       context.Context
	debugMode bool
	logjson   bool
	path      string
	port      int
	quietMode bool
	rootPath  string
	swarmMode bool
)

var rootCmd = &cobra.Command{
	Use:   "kin",
	Short: "Automatic home page for docker hosted web application",
	Long: `Generate home page with links to web application hosted by
	the same docker deamon from templates and docker labels.`,
	RunE: runServer,
}

func runServer(cmd *cobra.Command, args []string) error {
	parseArguments()

	db, err := cache.New()
	if err != nil {
		return fmt.Errorf("cannot initialize cache: %w", err)
	}

	end := make(chan bool, 1)
	ctx, cancel := context.WithCancel(ctx)
	g, ctx := errgroup.WithContext(ctx)

	// Manage graceful shutdown of http server on context cancellation
	g.Go(func() error {
		select {
		case <-ctx.Done():
			return server.Shutdown(ctx)
		case <-end:
			return nil
		}
	})

	// Start the http server
	g.Go(func() error {
		err := server.Run(ctx, log, db, baseURL, rootPath, port)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			end <- true
			return err
		}
		return nil
	})

	// Start the collector server
	g.Go(func() error {
		err := collector.Run(ctx, log, db, swarmMode)
		if err != nil {
			cancel()
			return err
		}
		return nil
	})

	return g.Wait()
}

// Execute the corresponding cobra sub-command.
func Execute(c context.Context) error {
	ctx = c
	return rootCmd.Execute()
}

//nolint: errcheck
func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&baseURL, "base", "b", "/", "base URL ")
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default $HOME/.kin.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&debugMode, "debug", "d", false, "log more information")
	rootCmd.PersistentFlags().BoolVarP(&logjson, "json", "j", false, "JSON-formatted logs")
	rootCmd.PersistentFlags().StringVarP(&path, "logpath", "l", "-", "log file path ")
	rootCmd.PersistentFlags().BoolVarP(&quietMode, "quiet", "q", false, "log only errors")
	rootCmd.PersistentFlags().IntVarP(&port, "port", "p", 8080, "port to listen")
	rootCmd.PersistentFlags().StringVarP(&rootPath, "root", "r", "", "template root path (default $HOME/.kin_root)")
	rootCmd.PersistentFlags().BoolVarP(&swarmMode, "swarm", "s", false, "Docker swarm")

	viper.BindEnv("base", "KIN_BASE")
	viper.BindEnv("log.path", "KIN_LOGPATH")
	viper.BindEnv("log.level", "KIN_LOGLEVEL")
	viper.BindEnv("log.json", "KIN_LOGJSON")
	viper.BindEnv("port", "KIN_PORT")
	viper.BindEnv("root", "KIN_ROOT")
	viper.BindEnv("swarm", "KIN_SWARM")

	viper.BindPFlag("base", rootCmd.PersistentFlags().Lookup("base"))
	viper.BindPFlag("log.path", rootCmd.PersistentFlags().Lookup("logpath"))
	viper.BindPFlag("log.json", rootCmd.PersistentFlags().Lookup("json"))
	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("root", rootCmd.PersistentFlags().Lookup("root"))
	viper.BindPFlag("swarm", rootCmd.PersistentFlags().Lookup("swarm"))
}

func initConfig() {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defaultRootPath := home + "/.kin_root"
	if fi, err := os.Stat(defaultRootPath); err == nil && fi.Mode().IsDir() {
		viper.SetDefault("rootpath", defaultRootPath)
	}

	viper.SetDefault("log.level", "info")

	switch {
	case cfgFile != "":
		viper.SetConfigFile(cfgFile)
	case os.Getenv("KIN_CONFIG") != "":
		viper.SetConfigFile(os.Getenv("KIN_CONFIG"))
	default:
		viper.AddConfigPath(home)
		viper.SetConfigName(".kin")
	}

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func parseArguments() {
	if debugMode {
		viper.Set("log.level", "debug")
	}

	if quietMode {
		viper.Set("log.level", "error")
	}

	configureLogging()

	baseURL = viper.GetString("base")

	_, err := url.ParseRequestURI(baseURL)
	if err != nil || strings.Contains(baseURL, ":") {
		log.Fatal("Base URL non valid, it can only be a relative URL")
	}

	rootPath = viper.GetString("root")
	port = viper.GetInt("port")

	if 1 > port || port > 65536 {
		log.Fatal("Port non valid")
	}

	if fi, err := os.Stat(rootPath); err == nil && fi.Mode().IsDir() {
		if fi, err := os.Stat(filepath.Join(rootPath, "index.html")); err == nil && !fi.Mode().IsDir() {
			log.Infof("Will use %s as templates root", rootPath)
		} else if fi, err := os.Stat(filepath.Join(rootPath, "index.html.tpl")); err == nil && !fi.Mode().IsDir() {
			log.Infof("Will use %s as templates root", rootPath)
		} else {
			log.Fatalf("%s exists but does not contain index.html file", rootPath)
		}
	} else {
		log.Infof("Template root %s does not exists, will use my own templates", rootPath)
		rootPath = ""
	}
}
