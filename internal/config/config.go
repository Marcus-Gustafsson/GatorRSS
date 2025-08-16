package config

import (
	"fmt"
	"os"
)

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() Config {
	// Export a Read function that reads the JSON file found at ~/.gatorconfig.json and returns a Config struct. CHECK!
	// It should read the file from the HOME directory, then decode the JSON string into a new Config struct.
	// I used os.UserHomeDir to get the location of HOME.
	homePath, err := os.UserHomeDir()

	if err != nil {
		fmt.Println(err)
		return Config{}
	}

	jsonPath := homePath + "/.gatorconfig.json"

	fmt.Printf("DBG: jsonPath = %v \n", jsonPath)
	// Open our jsonFile
	jsonFile, err := os.Open(jsonPath)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("DBG: Successfully Opened .gatorconfig.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()



	return Config{}
}
