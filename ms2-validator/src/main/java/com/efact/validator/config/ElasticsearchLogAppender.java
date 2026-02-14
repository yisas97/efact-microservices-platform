package com.efact.validator.config;

import ch.qos.logback.classic.spi.ILoggingEvent;
import ch.qos.logback.core.AppenderBase;
import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.databind.node.ObjectNode;

import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.time.LocalDate;
import java.time.format.DateTimeFormatter;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

public class ElasticsearchLogAppender extends AppenderBase<ILoggingEvent> {

    private String elasticsearchUrl = "http://localhost:9200";
    private String indexPrefix = "ms2-logs";
    private final HttpClient httpClient;
    private final ObjectMapper objectMapper;
    private final ExecutorService executorService;

    public ElasticsearchLogAppender() {
        this.httpClient = HttpClient.newBuilder().build();
        this.objectMapper = new ObjectMapper();
        this.executorService = Executors.newFixedThreadPool(2);
    }

    public void setElasticsearchUrl(String elasticsearchUrl) {
        this.elasticsearchUrl = elasticsearchUrl;
    }

    public void setIndexPrefix(String indexPrefix) {
        this.indexPrefix = indexPrefix;
    }

    @Override
    protected void append(ILoggingEvent event) {
        CompletableFuture.runAsync(() -> sendToElasticsearch(event), executorService);
    }

    private void sendToElasticsearch(ILoggingEvent event) {
        try {
            // Crear documento JSON
            ObjectNode logDoc = objectMapper.createObjectNode();
            logDoc.put("@timestamp", event.getInstant().toString());
            logDoc.put("level", event.getLevel().toString());
            logDoc.put("message", event.getFormattedMessage());
            logDoc.put("service", "ms2-validator");
            logDoc.put("logger", event.getLoggerName());
            logDoc.put("thread", event.getThreadName());

            // Nombre del índice con fecha
            String indexName = String.format("%s-%s",
                indexPrefix,
                LocalDate.now().format(DateTimeFormatter.ofPattern("yyyy.MM.dd"))
            );

            String url = String.format("%s/%s/_doc", elasticsearchUrl, indexName);

            // Crear request
            HttpRequest request = HttpRequest.newBuilder()
                .uri(URI.create(url))
                .header("Content-Type", "application/json")
                .POST(HttpRequest.BodyPublishers.ofString(objectMapper.writeValueAsString(logDoc)))
                .build();

            // Enviar de forma asíncrona
            httpClient.sendAsync(request, HttpResponse.BodyHandlers.ofString());

        } catch (Exception e) {
            // Ignorar errores para no afectar la aplicación
        }
    }

    @Override
    public void stop() {
        executorService.shutdown();
        super.stop();
    }
}
