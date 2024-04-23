package stoke

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// Gets the current tracer
func getTracer() trace.Tracer {
	return otel.GetTracerProvider().
		Tracer("stoke-client", 
			trace.WithInstrumentationVersion("0.1"),
		)
}

// Adds a token and all of it's field information to a given span
// Always returns false to allow incorperating it a chain of ors (|)
func AddTokenToSpan(jwtToken *jwt.Token, span trace.Span) bool {
	var traceAttrs []attribute.KeyValue
	issuer, err := jwtToken.Claims.GetIssuer()
	if err == nil && issuer != "" {
		traceAttrs = append(traceAttrs, attribute.String("jwt-issuer", issuer))
	}

	subject, err := jwtToken.Claims.GetSubject()
	if err == nil && subject != "" {
		traceAttrs = append(traceAttrs, attribute.String("jwt-subject", subject))
	}

	audience, err := jwtToken.Claims.GetAudience()
	if err == nil && len(audience) != 0 {
		traceAttrs = append(traceAttrs, attribute.String("jwt-audience", fmt.Sprint(audience)))
	}

	expires, err := jwtToken.Claims.GetExpirationTime()
	if err == nil && expires != nil {
		traceAttrs = append(traceAttrs, attribute.String("jwt-expires", expires.String()))
	}

	issued, err := jwtToken.Claims.GetIssuedAt()
	if err == nil && issued != nil {
		traceAttrs = append(traceAttrs, attribute.String("jwt-issued", issued.String()))
	}

	notbefore, err := jwtToken.Claims.GetNotBefore()
	if err == nil  && issued != nil {
		traceAttrs = append(traceAttrs, attribute.String("jwt-notbefore", notbefore.String()))
	}

	claims, ok := jwtToken.Claims.(*Claims)
	if ok {
		for key, value := range claims.StokeClaims {
			traceAttrs = append(traceAttrs, attribute.String("stoke-jwt-" + key, value))
		}
	}
	traceAttrs = append(traceAttrs, attribute.Bool("parsable", ok))

	span.AddEvent("Parsed Token", trace.WithAttributes(traceAttrs...))

	return false
}
