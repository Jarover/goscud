package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/Jarover/goscud/database"
	"github.com/Jarover/goscud/models"
	"github.com/Jarover/goscud/readconfig"
	"github.com/Jarover/goscud/utils"
	_ "github.com/nakagami/firebirdsql"
)

// Вычисляем продолжительность со старта события
func duration(msg string, start time.Time) {
	log.Printf("%v: %v\n", msg, time.Since(start))
}

// Читаем флаги и окружение
func readFlag(configFlag *readconfig.Flag) {
	flag.StringVar(&configFlag.ConfigFile, "f", readconfig.GetEnv("CONFIGFILE", readconfig.GetDefaultConfigFile()), "config file")
	flag.StringVar(&configFlag.Host, "h", readconfig.GetEnv("HOST", ""), "host")
	flag.UintVar(&configFlag.Port, "p", uint(readconfig.GetEnvInt("PORT", 0)), "port")
	flag.Parse()
}

func readEVN(lastId, limit uint) uint {

	rows, err := database.GetDB().Query("SELECT FIRST ? ID,DT,DVS,USR FROM FB_EVN WHERE USR>0 AND ID > ?", limit, lastId)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		ns := new(models.FbEnv)
		err = rows.Scan(&ns.ID, &ns.DT, &ns.DVS, &ns.USR)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(ns)
		lastId = ns.ID
	}
	return lastId
}

func main() {
	var configFlag readconfig.Flag
	start := time.Now()

	err := readconfig.Version.ReadVersionFile(utils.GetDir() + "/version.json")
	if err != nil {
		fmt.Println(err)
	}

	logFile, err := os.OpenFile(readconfig.GetBaseFile()+".log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	log.Printf(
		"Starting quicklog program...\ncommit: %s, build time: %s, release: %s",
		readconfig.Version.Commit,
		readconfig.Version.BuildTime,
		readconfig.Version.VersionStr())

	readFlag(&configFlag)

	// log.Println(configFlag)

	Config, err := readconfig.ReadConfig(configFlag.ConfigFile)
	if err != nil {
		log.Fatal(err)
	}

	database.InitDB(Config.Host, Config.Db, Config.User, Config.Pass)
	defer database.GetDB().Close()

	lastId, err := utils.ReadState(readconfig.GetBaseFile() + "_config.json")
	if err != nil {
		log.Println("Not Config file")
	}
	log.Printf("Start ID is %d", lastId)
	readEVN(lastId, Config.Limit)

	utils.SaveState(readconfig.GetBaseFile()+"_config.json", lastId)
	duration("Время работы программы", start)
}
