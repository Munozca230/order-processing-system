package com.example.orderworker.consumer;

import com.example.orderworker.model.OrderMessage;
import com.example.orderworker.service.EnrichmentService;
import com.example.orderworker.service.ValidationService;
import com.example.orderworker.repository.OrderRepository;
import com.example.orderworker.model.OrderDocument;
import com.example.orderworker.event.ProcessedOrderEvent;
import com.example.orderworker.service.OrderLockService;
import reactor.core.publisher.Mono;
import org.springframework.context.ApplicationEventPublisher;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.stereotype.Component;

import java.io.IOException;

@Component
public class OrderKafkaConsumer {

    private static final Logger logger = LoggerFactory.getLogger(OrderKafkaConsumer.class);
    private final ObjectMapper mapper = new ObjectMapper();
    private final EnrichmentService enrichmentService;
    private final ValidationService validationService;
    private final ApplicationEventPublisher publisher;
    private final OrderRepository orderRepository;
    private final OrderLockService lockService;

    public OrderKafkaConsumer(EnrichmentService enrichmentService,
                              ValidationService validationService,
                              ApplicationEventPublisher publisher,
                              OrderRepository orderRepository,
                              OrderLockService lockService) {
        this.enrichmentService = enrichmentService;
        this.validationService = validationService;
        this.publisher = publisher;
        this.orderRepository = orderRepository;
        this.lockService = lockService;
    }

    @KafkaListener(topics = "orders", groupId = "order-worker-group")
    public void consume(String message) {
        logger.info("Received order message: {}", message);
        try {
            OrderMessage order = mapper.readValue(message, OrderMessage.class);
            lockService.acquire(order.orderId())
                    .flatMap(acquired -> {
                        if (!acquired) {
                            logger.info("Order {} is already being processed, skipping", order.orderId());
                            return Mono.empty();
                        }
                        return enrichmentService.enrich(order)
                                .flatMap(validationService::validate)
                                .flatMap(valid -> orderRepository.save(OrderDocument.from(valid))
                                        .retryWhen(reactor.util.retry.Retry.backoff(5, java.time.Duration.ofSeconds(1)))
                                        .thenReturn(valid))
                                .doOnNext(enriched -> {
                                    logger.info("Enriched and validated order: {}", enriched);
                                    publisher.publishEvent(new ProcessedOrderEvent(this, enriched));
                                })
                                .doFinally(sig -> lockService.release(order.orderId()).subscribe());
                    })
                    .doOnError(err -> logger.error("Order processing failed", err))
                    .subscribe();
        } catch (IOException e) {
            logger.error("Failed to deserialize order", e);
        }
    }
}
