package draken

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/rs/zerolog/log"
	"golang.org/x/time/rate"
)

type R2 struct {
	AccountId       string
	AccessKeyId     string
	AccessKeySecret string
	Limiter         *rate.Limiter
	Client          *s3.Client
	Context         context.Context
	Cancel          context.CancelFunc
}

func NewR2(accountId, accessKeyId, accessKeySecret string) *R2 {
	config.WithRequestChecksumCalculation(0)
	config.WithResponseChecksumValidation(0)
	endpoint := fmt.Sprintf("https://%s.eu.r2.cloudflarestorage.com", accountId)

	ctx, cancel := context.WithCancel(context.Background())
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyId, accessKeySecret, "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		log.Fatal().Err(err).Msgf("Failed to load R2 configuration for endpoint %s", endpoint)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
	})

	log.Debug().Msgf("R2 configuration loaded for endpoint %s", endpoint)
	return &R2{
		AccountId:       accountId,
		AccessKeyId:     accessKeyId,
		AccessKeySecret: accessKeySecret,
		// 1 request every 2 seconds
		Limiter: rate.NewLimiter(rate.Every(2*time.Second), 1),
		Client:  client,
		Context: ctx,
		Cancel:  cancel,
	}
}

func (d *Draken) initR2() {
	if !d.Config.R2.Enabled {
		log.Debug().Msgf("R2 is disabled in the config, skipping...")
		return
	}
	log.Debug().Msgf("Initializing R2...")
	d.R2 = NewR2(d.Config.R2.AccountId, d.Config.R2.AccessKeyId, d.Config.R2.AccessKeySecret)
	log.Info().Msgf("R2 initialized.")
}
