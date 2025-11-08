package main

import (
	"fmt"
	"os/exec"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	cmd := exec.Command(
		"tern",
		"migrate",
		"--migrations",
		"./internal/store/pgstore/migrations",
		"--config",
		"./internal/store/pgstore/migrations/tern.conf",
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("command execution failed", err)
		fmt.Println("output", string(out))
		panic(err)
	}

	fmt.Println("command executed successfully", out)

}
