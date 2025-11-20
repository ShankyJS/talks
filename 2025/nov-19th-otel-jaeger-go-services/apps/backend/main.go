package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var tracer trace.Tracer

var wordList = []string{
	"sunshine", "adventure", "harmony", "serenity", "wisdom",
	"courage", "freedom", "journey", "discovery", "wonder",
	"creativity", "passion", "balance", "gratitude", "resilience",
	"innovation", "excellence", "integrity", "compassion", "unity",
}

type WordsResponse struct {
	Words     []string `json:"words"`
	Timestamp string   `json:"timestamp"`
}

func initTracer() func() {
	ctx := context.Background()

	otelEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if otelEndpoint == "" {
		otelEndpoint = "localhost:4317"
	}

	conn, err := grpc.Dial(otelEndpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to create gRPC connection: %v", err)
	}

	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		log.Fatalf("Failed to create OTLP exporter: %v", err)
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName("backend-service"),
			semconv.ServiceVersion("1.0.0"),
		),
	)
	if err != nil {
		log.Fatalf("Failed to create resource: %v", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	tracer = tp.Tracer("backend-service")

	return func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}
}

func getRandomWords(ctx context.Context, count int) []string {
	_, span := tracer.Start(ctx, "getRandomWords")
	defer span.End()

	span.SetAttributes(attribute.Int("word.count", count))

	words := make([]string, count)
	for i := 0; i < count; i++ {
		words[i] = wordList[rand.Intn(len(wordList))]
	}

	span.SetAttributes(attribute.StringSlice("words.selected", words))
	return words
}

func wordsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Simulate some processing time
	time.Sleep(time.Duration(50+rand.Intn(100)) * time.Millisecond)

	words := getRandomWords(ctx, 5)

	response := WordsResponse{
		Words:     words,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	log.Printf("Returned %d words: %v", len(words), words)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	rand.Seed(time.Now().UnixNano())

	shutdown := initTracer()
	defer shutdown()

	// Wrap handlers with OpenTelemetry HTTP middleware to extract trace context
	http.Handle("/words", otelhttp.NewHandler(http.HandlerFunc(wordsHandler), "wordsHandler"))
	http.Handle("/health", otelhttp.NewHandler(http.HandlerFunc(healthHandler), "healthHandler"))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Backend service starting on port %s", port)
	log.Printf("OTEL endpoint: %s", os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"))

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
