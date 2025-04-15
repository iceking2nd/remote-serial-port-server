/*
Copyright © 2025 Daniel Wu<wxc@wxccs.org>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"context"
	"fmt"
	"github.com/iceking2nd/remote-serial-port-server/app/routers"
	"github.com/iceking2nd/remote-serial-port-server/global"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/szuecs/gin-glog"
	"github.com/toorop/gin-logrus"
)

var (
	gListenAddress string
	gListenPort    uint
	gLogFile       string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "remote-serial-port-server",
	Short:   "A web terminal for local serial port",
	Version: global.Version,
	/*Long: `A longer description that spans multiple lines and likely contains
	examples and usage of using your application. For example:

	Cobra is a CLI library for Go that empowers applications.
	This application is a tool to generate the needed files
	to quickly create a Cobra application.`,*/
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		if global.LogLevel < 5 {
			gin.SetMode(gin.ReleaseMode)
		}
		corsConfig := cors.Config{
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "HEAD"},
			AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-API-Key"},
			AllowCredentials: true,
			MaxAge:           12 * time.Hour,
			AllowAllOrigins:  true,
		}

		global.APIKey = uuid.New().String()

		r := gin.Default()
		r.Use(ginglog.Logger(3 * time.Second))
		r.Use(cors.New(corsConfig))
		r.Use(ginlogrus.Logger(global.Log), gin.Recovery())

		root := r.Group("/")
		routers.SetupRouter(root)

		s := http.Server{Handler: r}

		go func() {
			ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", gListenAddress, gListenPort))
			if err != nil {
				log.Fatalf("create listener error: %s\n", err.Error())
			}
			fmt.Printf("listening on %s ...\n", ln.Addr().String())
			err = s.Serve(ln)
			if err != nil {
				log.Fatalf("http serve error: %s\n", err.Error())
			}
		}()

		signalChan := make(chan os.Signal)
		signal.Notify(signalChan, os.Interrupt, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT)
		ticker := time.NewTicker(time.Millisecond)
		for {
			select {
			case sig := <-signalChan:
				log.Println("Get Signal:", sig)
				log.Println("Shutdown Server ...")
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				if err := s.Shutdown(ctx); err != nil {
					log.Fatal("Closing web service error: ", err)
				}
				log.Println("Server exiting")
				os.Exit(0)
			case <-ticker.C:
				//do sth
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.remote-serial-port-server.yaml)")
	rootCmd.PersistentFlags().StringVarP(&gListenAddress, "listen-address", "l", "127.0.0.1", "listen address")
	rootCmd.PersistentFlags().UintVarP(&gListenPort, "listen-port", "p", 0, "listen port")
	rootCmd.PersistentFlags().StringVar(&gLogFile, "log-file", "", "logging file")
	rootCmd.PersistentFlags().Uint32Var(&global.LogLevel, "log-level", 3, "log level (0 - 6, 3 = warn , 5 = debug)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.SetVersionTemplate(fmt.Sprintf(`{{with .Name}}{{printf "%%s version information: " .}}{{end}}
   {{printf "Version:    %%s" .Version}}
   Build Time:		%s
   Git Revision:	%s
   Go version:		%s
   OS/Arch:			%s/%s
`, global.BuildTime, global.GitCommit, runtime.Version(), runtime.GOOS, runtime.GOARCH))
}

func initConfig() {
	global.Log = logrus.New()
	var logWriter io.Writer
	if gLogFile == "" {
		logWriter = os.Stdout
	} else {
		logFileHandle, err := os.OpenFile(gLogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			panic(err.Error())
		}
		logWriter = io.MultiWriter(os.Stdout, logFileHandle)
	}
	global.Log.SetOutput(logWriter)
	global.Log.SetLevel(logrus.Level(global.LogLevel))
}
