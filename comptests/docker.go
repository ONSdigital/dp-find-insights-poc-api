package comptests

import (
	"log"
	"net"
	"os"
	"os/exec"
	"time"

	"github.com/ONSdigital/dp-find-insights-poc-api/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const DefaultDSN = "postgres://insights:insights@localhost:54322/censustest"

func SetupDockerDB(dsn string) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	_, _, host, port, _ := model.ParseDSN(dsn)
	user := "postgres"
	db := "postgres"
	pw := os.Getenv("POSTGRES_PASSWORD")
	dsn = model.CreatDSN(user, pw, host, port, db)

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
				log.Println("connected")
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
