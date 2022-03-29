package comptests

import (
	"log"
	"net"
	"os"
	"os/exec"
	"time"

	"github.com/ONSdigital/dp-find-insights-poc-api/pkg/database"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const DefaultDSN = "postgres://insights:insights@localhost:54322/censustest"
const DefaultPostgresPW = "mylocalsecret"

func SetupDockerDB(dsn string) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	_, _, host, port, _ := database.ParseDSN(dsn)
	user := "postgres"
	db := "postgres"

	// get password from env, or fall back to default
	var pw string
	envPW := os.Getenv("POSTGRES_PASSWORD")
	if envPW != "" {
		pw = envPW
	} else {
		pw = DefaultPostgresPW
		os.Setenv("POSTGRES_PASSWORD", DefaultPostgresPW)
	}

	dsn = database.CreatDSN(user, pw, host, port, db)

	// is docker postgres+postgis running?
	_, err := net.DialTimeout("tcp", net.JoinHostPort("localhost", port), time.Second)
	if err != nil {
		log.Println("starting postgres docker")

		go func() {
			cmd := exec.Command("docker", "run", "--rm", "--name", "postgis", "--publish", port+":5432", "-e", "POSTGRES_PASSWORD="+pw, "postgis/postgis")
			if err := cmd.Run(); err != nil {
				log.Fatalf("is docker installed and running? %v", err)
			}
		}()

		// poll for start up
		for {
			time.Sleep(time.Second)
			_, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
			if err == nil {
				log.Printf("connected to %s", dsn)
				break
			}
			log.Println("polling for started docker...")
		}

	}

	log.Println("postgres docker running")
}

func KillDockerDB() {
	cmd := exec.Command("docker", "container", "kill", "postgis")
	if err := cmd.Run(); err != nil {
		log.Print(err)
	}

	log.Fatal("exiting")
}
