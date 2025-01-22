package config

import (
	"biolerplate/pkg/logger"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var cfgFile string
var outputDir string
var inputDir string
var sugar *zap.SugaredLogger
var rootCmd = &cobra.Command{
	Use:   "boilerplate",
	Short: "An example app reading configs",
	Run: func(cmd *cobra.Command, args []string) {
		sugar.Debugf("Input directory: %s\n", viper.GetString("configs.inputDir"))
		sugar.Debugf("Output directory: %s\n", viper.GetString("configs.outputDir"))
	},
}

func init() {
	logger.InitLogger(true)
	sugar = logger.GetSugar()
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/config.toml)")
	rootCmd.Flags().StringVarP(&outputDir, "outputDir", "o", "", "Output directory")
	rootCmd.Flags().StringVarP(&inputDir, "inputDir", "i", "", "Input directory")
	viper.BindPFlag("configs.outputDir", rootCmd.Flags().Lookup("outputDir"))
	viper.BindPFlag("configs.inputDir", rootCmd.Flags().Lookup("inputDir"))
	if err := rootCmd.Execute(); err != nil {
		sugar.Error(err)
	}
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("configs") // 添加配置檔案所在的路徑
		viper.SetConfigName("configs") // 設定配置檔案名稱(無需副檔名)
	}

	viper.AutomaticEnv() // 讀取匹配的環境變數

	if err := viper.ReadInConfig(); err == nil {

		sugar.Debugf("Using config file: %s\n\r", viper.ConfigFileUsed())
	} else {
		sugar.Debugf("Error reading config file: %s \n", err)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		sugar.Error(err)
	}
}
