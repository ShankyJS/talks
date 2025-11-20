package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
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
var backendURL string

type WordsResponse struct {
	Words     []string `json:"words"`
	Timestamp string   `json:"timestamp"`
}

const htmlTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Random Words Generator</title>
    <style>
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            display: flex;
            justify-content: center;
            align-items: center;
            padding: 20px;
        }

        .container {
            background: white;
            border-radius: 20px;
            padding: 40px;
            box-shadow: 0 20px 60px rgba(0, 0, 0, 0.3);
            max-width: 600px;
            width: 100%;
            text-align: center;
        }

        h1 {
            color: #667eea;
            margin-bottom: 30px;
            font-size: 2.5em;
        }

        button {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            border: none;
            padding: 15px 40px;
            font-size: 1.2em;
            border-radius: 50px;
            cursor: pointer;
            transition: transform 0.2s, box-shadow 0.2s;
            margin: 20px 0;
        }

        button:hover {
            transform: translateY(-2px);
            box-shadow: 0 10px 20px rgba(102, 126, 234, 0.4);
        }

        button:active {
            transform: translateY(0);
        }

        #words {
            margin-top: 30px;
            min-height: 100px;
        }

        .word {
            display: inline-block;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            padding: 10px 20px;
            margin: 10px;
            border-radius: 25px;
            font-size: 1.1em;
            animation: fadeIn 0.5s ease-in;
        }

        @keyframes fadeIn {
            from {
                opacity: 0;
                transform: translateY(20px);
            }
            to {
                opacity: 1;
                transform: translateY(0);
            }
        }

        .timestamp {
            color: #666;
            font-size: 0.9em;
            margin-top: 20px;
        }

        .loading {
            color: #667eea;
            font-size: 1.2em;
        }

        .error {
            color: #e74c3c;
            padding: 15px;
            background: #fee;
            border-radius: 10px;
            margin-top: 20px;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>🎲 Random Words</h1>
        <p>Click the button to generate random inspirational words!</p>
        <button onclick="fetchWords()">Generate Words</button>
        <div id="words"></div>
    </div>

    <script>
        async function fetchWords() {
            const wordsDiv = document.getElementById('words');
            wordsDiv.innerHTML = '<p class="loading">Loading...</p>';

            try {
                const response = await fetch('/api/words');
                if (!response.ok) {
                    throw new Error('Failed to fetch words');
                }

                const data = await response.json();

                wordsDiv.innerHTML = data.words.map(word =>
                    '<span class="word">' + word + '</span>'
                ).join('') + '<p class="timestamp">Generated at: ' + new Date(data.timestamp).toLocaleString() + '</p>';
            } catch (error) {
                wordsDiv.innerHTML = '<p class="error">Error: ' + error.message + '</p>';
            }
        }
    </script>
</body>
</html>
`

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
			semconv.ServiceName("frontend-service"),
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

	tracer = tp.Tracer("frontend-service")

	return func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}
}

func fetchWordsFromBackend(ctx context.Context) (*WordsResponse, error) {
	ctx, span := tracer.Start(ctx, "fetchWordsFromBackend")
	defer span.End()

	span.SetAttributes(attribute.String("backend.url", backendURL))

	req, err := http.NewRequestWithContext(ctx, "GET", backendURL+"/words", nil)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	// Propagate trace context
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}
	defer resp.Body.Close()

	span.SetAttributes(attribute.Int("http.status_code", resp.StatusCode))

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("backend returned status %d", resp.StatusCode)
		span.RecordError(err)
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		span.RecordError(err)
		return nil, err
	}

	var wordsResp WordsResponse
	if err := json.Unmarshal(body, &wordsResp); err != nil {
		span.RecordError(err)
		return nil, err
	}

	span.SetAttributes(attribute.Int("words.count", len(wordsResp.Words)))
	return &wordsResp, nil
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.New("home").Parse(htmlTemplate)
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	if err := tmpl.Execute(w, nil); err != nil {
		log.Printf("Template execution error: %v", err)
	}
}

func apiWordsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	words, err := fetchWordsFromBackend(ctx)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch words: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(words); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	log.Printf("Served %d words to client", len(words.Words))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	shutdown := initTracer()
	defer shutdown()

	backendURL = os.Getenv("BACKEND_URL")
	if backendURL == "" {
		backendURL = "http://localhost:8080"
	}

	// Wrap handlers with OpenTelemetry HTTP middleware to extract trace context
	http.Handle("/", otelhttp.NewHandler(http.HandlerFunc(homeHandler), "homeHandler"))
	http.Handle("/api/words", otelhttp.NewHandler(http.HandlerFunc(apiWordsHandler), "apiWordsHandler"))
	http.Handle("/health", otelhttp.NewHandler(http.HandlerFunc(healthHandler), "healthHandler"))

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Frontend service starting on port %s", port)
	log.Printf("Backend URL: %s", backendURL)
	log.Printf("OTEL endpoint: %s", os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"))

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
