package bootstrap

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/propagators/b3"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

func (c *Container) GetTracer() trace.Tracer {
	return otel.Tracer(viper.GetString("telemetry.service_name"))
}

func (c *Container) initTracer() *sdktrace.TracerProvider {
	if !viper.GetBool("telemetry.enable") {
		logrus.Debug("opentelemetry disabled")
		return nil
	}

	logrus.Debug("opentelemetry initialize")

	host := viper.GetString("telemetry.jaeger.agent_host")
	port := viper.GetString("telemetry.jaeger.agent_port")

	logrus.
		WithField("bootstrap", "jaeger").
		Debugf("trying to connect to %s:%s", host, port)

	exporter, err := jaeger.New(jaeger.WithAgentEndpoint(
		jaeger.WithAgentHost(host),
		jaeger.WithAgentPort(port),
	))

	if err != nil {
		logrus.Fatal(err)
	}

	logrus.
		WithField("bootstrap", "jaeger").
		Debugf("connected to %s:%s", host, port)

	c.trace = sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(viper.GetString("telemetry.service_name")),
			attribute.String("environment", viper.GetString("telemetry.environment")),
		)),
	)

	otel.SetTracerProvider(c.trace)
	otel.SetTextMapPropagator(
		b3.New(b3.WithInjectEncoding(b3.B3MultipleHeader | b3.B3SingleHeader)),
	)

	return c.trace
}

func (c *Container) GetTracerProvider() *sdktrace.TracerProvider {
	return c.trace
}
