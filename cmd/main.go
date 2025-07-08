package main

import (
    "mihome/api/v1"
    "mihome/config"
    "mihome/infra"
)

func main() {
    cfg := config.LoadConfig()
    db := infra.InitDB(cfg)

	infra.SeedData(db)

    r := v1.NewRouter(db, cfg)
    r.Run(":" + cfg.Server.Port)
}
