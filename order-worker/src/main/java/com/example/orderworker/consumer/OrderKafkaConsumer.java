package com.example.orderworker.consumer;

import com.example.orderworker.model.OrderMessage;
import com.example.orderworker.service.EnrichmentService;
import com.example.orderworker.service.ValidationService;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.stereotype.Component;
import reactor.core.publisher.Mono;

import java.io.IOException;

@Component
public class OrderKafkaConsumer {

    private static final Logger logger = LoggerFactory.getLogger(OrderKafkaConsumer.class);
    private final ObjectMapper mapper = new ObjectMapper();
    private final EnrichmentService enrichmentService;
    private final ValidationService validationService;

    public OrderKafkaConsumer(EnrichmentService enrichmentService, ValidationService validationService) {
        this.enrichmentService = enrichmentService;
        this.validationService = validationService;
    }

    @KafkaListener(topics = "orders", groupId = "order-worker-group")
    public void consume(String message) {
        logger.info("Received order message: {}", message);
        try {
            OrderMessage order = mapper.readValue(message, OrderMessage.class);
            enrichmentService.enrich(order)
                    .flatMap(validationService::validate)
                    .doOnNext(enriched -> logger.info("Enriched and validated order: {}", enriched))
                    .doOnError(err -> logger.error("Order processing failed", err))
                    .subscribe();
        } catch (IOException e) {
            logger.error("Failed to deserialize order", e);
        }
    }
}
