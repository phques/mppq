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
	initDone     bool = false
	mppqProvider *mppq.Provider

	appFilesDir    string
	configFilepath string
)

//------

// Start http, mppq servers, registers androidPush with mppq
// nb: InitAppFilesDir should be called 1st
func Start() error {

	log.Println("provider.Start")

	// start http server
	if err := StartHTTP(); err != nil {
		return err
	}

	// create/start mppq provider
	mppqProvider = mppq.NewProvider()
	err := mppqProvider.Start()
	if err != nil {
		return err
	}

	// register androidPush
	registerMppqService("androidPush")

	return nil
}

// Initialize the app's files dir, copies config file there if 1st time
// called from android app
func InitAppFilesDir(appFilesDir_ string) error {
	// already done ?
	if initDone {
		return nil
	}
	initDone = true

	//## debug
	dir, _ := os.Getwd()
	log.Printf("cwd: %v\n", dir)

	appFilesDir = appFilesDir_

	// setup config file path
	configFilepath = filepath.Join(appFilesDir, configFilename)
	log.Print("config file:", configFilepath)

	// create initial (copy from assets) config.json in appFilesDir if does not exists
	// does config file exist in app files dir?
	if _, err := os.Stat(configFilepath); err != nil {
		return copyConfigFile()
	}

	return nil
}

//--- utils -----

// copy config file from assets to app filesdir
func copyConfigFile() (err error) {
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

	return nil
}

// register a service we provide with mppq
func registerMppqService(serviceName string) {

	log.Println("registerMppqService", serviceName)

	// register a service (mppqProvider must be started)
	//## PQ use 'deviceName' from config
	providerName, _ := os.Hostname() // returns 'localhost' on my Nexus 7
	mppqProvider.AddService(mppq.ServiceDef{
		ServiceName:  serviceName,
		ProviderName: providerName,
		HostPort:     httpListenPort,
		Protocol:     "jsonhttp",
	})
}
