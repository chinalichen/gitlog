package main

import (
	"os"

	"github.com/chinalichen/gitlog/pkg/gitprocess"
	"github.com/sirupsen/logrus"
)

func main() {
	gp := gitprocess.NewGetProcessor("./")
	csv, err := gp.GitLog("./")
	if err != nil {
		logrus.Errorf("git log error: %v", err)
	}
	os.WriteFile("./gitlog.csv", []byte(csv), os.ModePerm)
}
