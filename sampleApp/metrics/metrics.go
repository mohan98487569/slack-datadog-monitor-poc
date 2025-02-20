package metrics

import (
	"fmt"
	log "sample_app/logFolder"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

const namespace = "datadog-demo"

type Metrics struct {
	service  *cloudwatch.CloudWatch
	counter1 int
}

func New(region string) (*Metrics, error) {
	log.Debug("Creating new metrics service")
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, fmt.Errorf("could not start aws session: %s", err)
	}

	log.Info("Successfully initialized AWS session")
	return &Metrics{
		service: cloudwatch.New(sess),
	}, nil
}

func (m *Metrics) IncrementCounter1() {
	m.counter1++
}

func (m *Metrics) Publish() {
	log.Info("publishing metrics")

	for {
		_, err := m.service.PutMetricData(&cloudwatch.PutMetricDataInput{
			Namespace: aws.String(namespace),
			MetricData: []*cloudwatch.MetricDatum{
				{
					MetricName: aws.String("Count"),
					Unit:       aws.String(cloudwatch.StandardUnitCount),
					Value:      aws.Float64(float64(m.counter1)),
					Dimensions: []*cloudwatch.Dimension{
						{
							Name:  aws.String("type"),
							Value: aws.String("Counter1"),
						},
					},
				},
			},
		})
		if err != nil {
			log.Errorf("Could not publish metrics: %s", err)
		} else {
			log.Infof("Successfully published metric counter1 = %d.\n", m.counter1)
		}

		// after 5 mins, send Counter1 as 0
		log.Info("Sleep for 5 mins")
		time.Sleep(5 * time.Minute)
		_, err = m.service.PutMetricData(&cloudwatch.PutMetricDataInput{
			Namespace: aws.String(namespace),
			MetricData: []*cloudwatch.MetricDatum{
				{
					MetricName: aws.String("Count"),
					Unit:       aws.String(cloudwatch.StandardUnitCount),
					Value:      aws.Float64(float64(0)),
					Dimensions: []*cloudwatch.Dimension{
						{
							Name:  aws.String("type"),
							Value: aws.String("Counter1"),
						},
					},
				},
			},
		})
		if err != nil {
			log.Errorf("Could not publish metrics: %s", err)
		} else {
			log.Infof("Successfully published metric counter1 = %d.\n", 0)
		}
	}
}
