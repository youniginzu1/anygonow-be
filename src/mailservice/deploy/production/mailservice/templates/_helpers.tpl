{{/*
Expand the name of the chart.
*/}}
{{- define "mailservice.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "mailservice.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "mailservice.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "mailservice.labels" -}}
helm.sh/chart: {{ include "mailservice.chart" . }}
{{ include "mailservice.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "mailservice.selectorLabels" -}}
app.kubernetes.io/name: {{ include "mailservice.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "mailservice.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "mailservice.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
List env variables
*/}}
{{- define "helpers.list-env-variables"}}
{{- range $key, $val := .Values.env.secret }}
- name: {{ $key }}
  valueFrom:
    configMapRef:
      name: env
      key: {{ $key }}
{{- end}}
{{- end }}

{{/*
Convert env to configmap
*/}}
{{- define "helpers.convert-env" -}}
  {{- if ne .line "" }}
    {{- if (not (hasPrefix "#" .line )) }}
      {{- $arr := (regexSplit "=" .line 2) -}}
        {{- index $arr 0 -}}: {{ index $arr 1 | quote -}}
    {{- end }}
  {{- end }}
{{- end }}
