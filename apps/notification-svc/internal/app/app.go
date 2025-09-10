package app

import (
	"context"
	"log"
	"sync"

	"github.com/hassiimykyta/life-rpg/apps/notification-svc/internal/consumers"
	"github.com/hassiimykyta/life-rpg/apps/notification-svc/internal/mailer"
	"github.com/hassiimykyta/life-rpg/apps/notification-svc/internal/service"
	"github.com/hassiimykyta/life-rpg/pkg/config"
	"github.com/hassiimykyta/life-rpg/pkg/helpers"
	"gopkg.in/gomail.v2"
)

type App struct {
	cfg         *config.Config
	ctx         context.Context
	cancel      context.CancelFunc
	mailSender  *mailer.MailSender
	mailBuilder mailer.MailBuilder
	userReg     *consumers.UserRegistered
	wg          sync.WaitGroup
	errCh       chan error
}

func New() (*App, error) {
	cfg, err := config.Load(config.WithSMTP())
	if err != nil {
		return nil, err
	}

	brokers := helpers.Csv(helpers.GetEnv("KAFKA_BROKERS", "kafka:9092"))
	groupID := helpers.GetEnv("KAFKA_GROUP_ID", "notification-svc")

	d := gomail.NewDialer(cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.Username, cfg.SMTP.Password)
	sender := mailer.New(d, "life-rpg@noreply.com")
	builder := mailer.NewHTMLBuilder()
	svc := service.NewNotificationService(builder, sender)

	userReg := consumers.NewUserRegistered(brokers, groupID, svc)

	ctx, cancel := context.WithCancel(context.Background())

	return &App{
		cfg:         cfg,
		ctx:         ctx,
		cancel:      cancel,
		mailSender:  sender,
		mailBuilder: builder,
		userReg:     userReg,
		errCh:       make(chan error, 1),
	}, nil
}

func (a *App) Start() error {
	log.Printf("notification-svc starting (env=%s)", a.cfg.App.Env)

	a.wg.Add(1)
	go func() {
		defer a.wg.Done()
		if err := a.userReg.Start(a.ctx); err != nil && err != context.Canceled {
			a.errCh <- err
		}
	}()

	log.Printf("notification-svc started (brokers=%v, group=%s)",
		helpers.Csv(helpers.GetEnv("KAFKA_BROKERS", "kafka:9092")),
		helpers.GetEnv("KAFKA_GROUP_ID", "notification-svc"),
	)
	return nil
}

func (a *App) Stop(_ context.Context) error {
	a.cancel()
	a.wg.Wait()

	_ = a.userReg.Close()

	close(a.errCh)
	log.Println("notification-svc stopped")
	return nil
}

func (a *App) ErrChan() <-chan error { return a.errCh }
