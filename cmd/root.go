package cmd

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/marema31/kin/cache"
	"github.com/marema31/kin/server"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	baseURL   string
	cfgFile   string
	debugMode bool
	path      string
	port      int
	quietMode bool
	rootPath  string
)

var rootCmd = &cobra.Command{
	Use:   "kin",
	Short: "Automatic home page for docker hosted web application",
	Long: `Generate home page with links to web application hosted by
	the same docker deamon from templates and docker labels.`,
	Run: runServer,
}

func runServer(cmd *cobra.Command, args []string) {
	parseArguments()

	db, err := cache.New()
	if err != nil {
		log.Fatalf("Cannot initialize cache: %v", err)
	}

	//TODO: remove this test datas
	containers := []cache.ContainerInfo{
		{Name: "Mon Site 1", URL: "http://localhost/1"},
		{Name: "Mon Site 2", URL: "http://localhost/2"},
		{Name: "Mon Site 3", URL: "http://localhost/3"},
	}

	txn := db.Db.Txn(true)

	for _, c := range containers {
		if err := txn.Insert("container", c); err != nil {
			log.Fatal(err)
		}
	}

	txn.Commit()

	//TODO: do something with the error returned
	_ = server.Run(logger, db, baseURL, rootPath, port)
}

// Execute the corresponding cobra sub-command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

//nolint: errcheck
func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&baseURL, "base", "b", "/", "base URL (default is '/')")
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.kin.yaml)")
	rootCmd.PersistentFlags().BoolVarP(&debugMode, "debug", "d", false, "log more information")
	rootCmd.PersistentFlags().StringVarP(&path, "logpath", "l", "-", "log file path (default is '-' for screen")
	rootCmd.PersistentFlags().BoolVarP(&quietMode, "quiet", "q", false, "log only errors")
	rootCmd.PersistentFlags().IntVarP(&port, "port", "p", 8080, "port to listen")
	rootCmd.PersistentFlags().StringVarP(&rootPath, "root", "r", "", "template root path (default is $HOME/.kin_root)")

	viper.BindEnv("base", "KIN_BASE")
	viper.BindEnv("log.path", "KIN_LOGPATH")
	viper.BindEnv("port", "KIN_PORT")
	viper.BindEnv("root", "KIN_ROOT")

	viper.BindPFlag("base", rootCmd.PersistentFlags().Lookup("base"))
	viper.BindPFlag("log.path", rootCmd.PersistentFlags().Lookup("logpath"))
	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("root", rootCmd.PersistentFlags().Lookup("root"))
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

	viper.SetDefault("log.json", false)
	viper.SetDefault("log.level", "info")

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
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
