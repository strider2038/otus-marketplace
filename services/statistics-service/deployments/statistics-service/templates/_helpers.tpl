{{/*
Expand the name of the chart.
*/}}
{{- define "statistics-service.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "statistics-service.fullname" -}}
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
{{- define "statistics-service.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "statistics-service.labels" -}}
helm.sh/chart: {{ include "statistics-service.chart" . }}
{{ include "statistics-service.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "statistics-service.selectorLabels" -}}
app.kubernetes.io/name: {{ include "statistics-service.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "statistics-service.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "statistics-service.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Application secrets
*/}}
{{- define "databaseUrl" }}
{{- with .Values.secrets.postgres }}
{{- printf "postgres://%s:%s@%s:%s/%s?%s" .user .password .host .port .dbname .sslmode | b64enc }}
{{- end }}
{{- end }}

{{- define "kafkaConsumerUrl" }}
{{- with .Values.secrets.kafka.consumer }}
{{- printf "%s" .url | b64enc }}
{{- end }}
{{- end }}
