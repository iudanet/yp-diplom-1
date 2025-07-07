package main

import "github.com/iudanet/yp-diplom-1/internal/config"

func main() {
	cfg := config.New()
	cfg.FlagParse()
	cfg.EnvParse()
}
