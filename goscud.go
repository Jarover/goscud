package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

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

	DbaUrl := Config.User + ":" + Config.Pass + "@" + Config.Host + "/" + Config.Db

	conn, _ := sql.Open("firebirdsql", DbaUrl)
	log.Println(DbaUrl)
	//conn, _ := sql.Open("firebirdsql", "SYSDBA:masterkey@skud-serv/skud")
	defer conn.Close()

	lastId, err := utils.ReadState(readconfig.GetBaseFile() + "_config.json")
	if err != nil {
		log.Println("Not Config file")
	}
	log.Printf("Start ID is %d", lastId)

	rows, err := conn.Query("SELECT FIRST ? ID,DT FROM FB_EVN WHERE ID > ?", Config.Limit, lastId)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		ns := new(models.FbEnv)
		err = rows.Scan(&ns.ID, &ns.DT)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(ns)
		lastId = ns.ID
	}
	utils.SaveState(readconfig.GetBaseFile()+"_config.json", lastId)
	duration("Время работы программы", start)
}
