{{- define "monitor.collector-config" -}}
receivers:
	otlp:
		protocols:
			grpc:
processors:
	batch:
exporters:
	jaeger:
		endpoint: jaeger:14250
		tls:
			insecure: true
service:
	pipelines:
		traces:
			receivers:
				- otlp
			processors: 
				- batch
			exporters:
				- jaeger
{{- end }}
