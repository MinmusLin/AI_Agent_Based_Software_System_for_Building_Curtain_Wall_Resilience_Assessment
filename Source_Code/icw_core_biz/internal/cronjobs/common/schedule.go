package common

import (
	"github.com/robfig/cron/v3"
)

// Schedule 按 Cron 表达式启动定时任务
func (c *CronJob) Schedule(name, expression string, fn func() error) error {
	scheduler := cron.New(
		cron.WithSeconds(),
		cron.WithChain(cron.SkipIfStillRunning(cron.DefaultLogger)),
	)

	if _, err := scheduler.AddFunc(expression, func() {
		_ = c.Run(name, fn)
	}); err != nil {
		return err
	}

	scheduler.Start()

	go func() {
		<-c.Ctx().Done()
		stopCtx := scheduler.Stop()
		<-stopCtx.Done()
	}()

	return nil
}
