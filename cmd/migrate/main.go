package main

import (
    "flag"
    "log"

    "github.com/gohex/gohex/internal/infrastructure/config"
    "github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/database/mysql"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
    configPath := flag.String("config", "configs/config.yaml", "Path to config file")
    flag.Parse()

    // 加载配置
    cfg, err := config.Load(*configPath)
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // 构建 DSN
    dsn := cfg.Database.DSN()

    // 创建迁移实例
    m, err := migrate.New(
        "file://migrations",
        dsn,
    )
    if err != nil {
        log.Fatalf("Failed to create migrate instance: %v", err)
    }
    defer m.Close()

    // 执行迁移
    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        log.Fatalf("Failed to run migrations: %v", err)
    }

    log.Println("Migrations completed successfully")
} 