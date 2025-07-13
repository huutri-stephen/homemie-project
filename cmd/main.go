package main

import (
    "mihome/config"

    "mihome/internal"
    "mihome/internal/infra"
)

func main() {
    cfg := config.LoadConfig()
    db := infra.InitDB(cfg)

	infra.SeedData(db)

    r := internal.NewRouter(db, cfg)
    r.Run(":" + cfg.Server.Port)
}
