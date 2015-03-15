// Package provider
package provider

import (
	//	"fmt"
	"github.com/phques/mppq"
	"golang.org/x/mobile/app"
	"io"
	"log"
	"os"
	"path/filepath"
)

const (
	configFilename string = "config.json"
)

var (
	provider       *mppq.Provider
	appFilesDir    string
	configFilepath string
)

//------

// called in main, through app.Run(Start callback)
func Start() {

	log.Println("provider.Start")

	// create/start mppq provider
	provider = mppq.NewProvider()
	provider.Start()

}

// called in main, through app.Run(Start callback)
func Register(serviceName string) {

	log.Println("provider.Register", serviceName)

	// register a service (provider main loop must be running)
	hostname, _ := os.Hostname() // returns 'localhost' on my Nexus 7
	provider.AddService(mppq.ServiceDef{
		ServiceName:  serviceName,
		ProviderName: hostname,
		HostPort:     1234,
		Protocol:     "jsonrpc1",
	})
}

//------- methods for Android App -----

// Initialize the app's files dir, copies config file there if 1st time
// called from android app
func InitAppFilesDir(appFilesDir_ string) {
	appFilesDir = appFilesDir_

	// setup config file path
	configFilepath = filepath.Join(appFilesDir, configFilename)
	log.Print("config file:", configFilepath)

	// create initial (copy from assets) config.json in appFilesDir if does not exists
	// does config file exist in app files dir?
	if _, err := os.Stat(configFilepath); err != nil {
		copyConfigFile()
	}
}

//--- utils -----

// copy config file from assets to app filesdir
func copyConfigFile() {
	// open src config file from assets
	srcFile, err := app.Open(configFilename)
	if err != nil {
		log.Printf("copyConfigFile, error opening source : %v\n", err)
		return
	}
	defer srcFile.Close()

	// create/open dest config file
	destFile, err := os.Create(configFilepath)
	if err != nil {
		log.Printf("copyConfigFile, error opening dest : %v\n", err)
		return
	}
	defer destFile.Close()

	// copy
	nbCopied, err := io.Copy(destFile, srcFile)
	if err == nil {
		log.Printf("copyConfigFile, copied %v bytes\n", nbCopied)
	} else {
		log.Printf("copyConfigFile, error copying : %v\n", err)
	}
}
