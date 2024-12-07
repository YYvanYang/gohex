package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "syscall"

    "github.com/gohex/gohex/internal/infrastructure/bootstrap"
)

func main() {
    // 创建上下文，用于优雅关闭
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // 监听系统信号
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    // 创建应用实例
    app, err := bootstrap.NewApplication("configs/config.yaml")
    if err != nil {
        log.Fatalf("Failed to create application: %v", err)
    }

    // 启动应用
    go func() {
        if err := app.Start(ctx); err != nil {
            log.Printf("Failed to start application: %v", err)
            cancel()
        }
    }()

    // 等待中断信号
    <-sigChan
    log.Println("Received shutdown signal")

    // 优雅关闭
    shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), app.Config().ShutdownTimeout)
    defer shutdownCancel()

    if err := app.Stop(shutdownCtx); err != nil {
        log.Printf("Error during shutdown: %v", err)
    }
} 