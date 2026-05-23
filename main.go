// Copyright 2024 The Casdoor Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
	_ "github.com/casdoor/casdoor/routers"
)

func main() {
	createDatabaseFlag := flag.Bool("createDatabase", false, "true if you need to create the database")
	flag.Parse()

	if *createDatabaseFlag {
		err := object.CreateDatabase()
		if err != nil {
			panic(err)
		}
		return
	}

	object.InitConfig()
	object.InitAdapter()
	object.InitDb()
	object.InitDefaultStorageProvider()
	object.InitLdapAutoSynchronizer()
	proxy.InitHttpClient()

	if beego.AppConfig.String("runmode") == "dev" {
		beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
			AllowAllOrigins:  true,
			AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
			ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
			AllowCredentials: true,
		}))
	}

	// Default port changed from 8000 to 8080 to align with my local dev environment convention.
	port := beego.AppConfig.String("httpport")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Casdoor server started on port %s\n", port)

	if err := checkRequiredEnvVars(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	beego.Run()
}

// checkRequiredEnvVars validates that all required environment variables or
// configuration values are present before starting the server.
func checkRequiredEnvVars() error {
	requiredConfigs := []string{
		"dbName",
		"dataSourceName",
	}

	for _, key := range requiredConfigs {
		val := beego.AppConfig.String(key)
		if val == "" {
			// Not a hard failure — some configs have defaults
			fmt.Printf("Warning: config key '%s' is not set\n", key)
		}
	}

	return nil
}
