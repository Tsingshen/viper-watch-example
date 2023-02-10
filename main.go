package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Person struct {
	Name    string `yaml:"name,omitempty"`
	Version string `yaml:"version,omitempty"`
}

var stopCh chan struct{} = make(chan struct{}, 1)
var ch chan struct{} = make(chan struct{}, 1)

var (
	ConfigPath string = "./config"
	ConfigName string = "config.yaml"
	ConfigType string = "yaml"
)

func main() {

	var configFile string
	v := viper.NewWithOptions(viper.KeyDelimiter("::"))
	v.SetConfigType(ConfigType)

	configFile = "./config/config.yaml"
	v.SetConfigFile(configFile)

	lc := &Person{}

	err1 := v.ReadInConfig()
	if err1 != nil {
		log.Println("read in config err: ", err1)
		os.Exit(1)
	}

	if err := v.Unmarshal(lc); err != nil {
		panic(err)
	}

	v.OnConfigChange(func(in fsnotify.Event) {
                // when watch k8s configMap, it should be set to fsnotify.Create,
                // or you test it by fsnotify.Write or others
		if in.Has(fsnotify.Create) {

			if err := v.Unmarshal(lc); err != nil {
				panic(err)
			}

			stopCh <- struct{}{}
			go func() {
				log.Printf("config update, name=%s,op=%v, content=%v", in.Name, in.Op, *lc)
				lc.someFunc(stopCh)
			}()
		}

	})
	v.WatchConfig()

	lc.someFunc(stopCh)
	<-ch
}

func (p *Person) someFunc(ch <-chan struct{}) {
	log.Printf("get ch in, p=%v, and wait <-ch", p)
	<-ch
}
