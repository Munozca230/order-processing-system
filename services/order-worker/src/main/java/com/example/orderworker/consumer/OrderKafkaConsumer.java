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
        logger.info("🔄 RECEIVED order message: {}", message);
        try {
            OrderMessage order = mapper.readValue(message, OrderMessage.class);
            logger.info("📦 PARSED OrderMessage: orderId={}, customerId={}, products={}", 
                order.orderId(), order.customerId(), order.products());
            
            lockService.acquire(order.orderId())
                    .flatMap(acquired -> {
                        if (!acquired) {
                            logger.warn("🔒 Order {} is already being processed, skipping", order.orderId());
                            return Mono.empty();
                        }
                        logger.info("🔓 ACQUIRED lock for order: {}", order.orderId());
                        
                        return enrichmentService.enrich(order)
                                .doOnSuccess(enriched -> logger.info("✅ ENRICHMENT SUCCESS for order: {}, products enriched: {}", 
                                    enriched.order().orderId(), enriched.products().size()))
                                .doOnError(error -> logger.error("❌ ENRICHMENT FAILED for order: {}", order.orderId(), error))
                                .flatMap(enriched -> {
                                    logger.info("🔍 STARTING validation for order: {}", enriched.order().orderId());
                                    return validationService.validate(enriched)
                                            .doOnSuccess(valid -> logger.info("✅ VALIDATION SUCCESS for order: {}", valid.order().orderId()))
                                            .doOnError(error -> logger.error("❌ VALIDATION FAILED for order: {}", enriched.order().orderId(), error));
                                })
                                .flatMap(valid -> {
                                    logger.info("💾 STARTING MongoDB save for order: {}", valid.order().orderId());
                                    OrderDocument doc = OrderDocument.from(valid);
                                    logger.info("📄 Created OrderDocument: id={}, orderId={}, customerId={}, products={}", 
                                        doc.id(), doc.orderId(), doc.customerId(), doc.products().size());
                                    
                                    return orderRepository.save(doc)
                                            .retryWhen(reactor.util.retry.Retry.backoff(5, java.time.Duration.ofSeconds(1)))
                                            .doOnSuccess(saved -> logger.info("✅ MONGODB SAVE SUCCESS: id={}, orderId={}", saved.id(), saved.orderId()))
                                            .doOnError(error -> logger.error("❌ MONGODB SAVE FAILED for order: {}", valid.order().orderId(), error))
                                            .thenReturn(valid);
                                })
                                .doOnNext(enriched -> {
                                    logger.info("🎉 ORDER PROCESSING COMPLETED: {}", enriched.order().orderId());
                                    publisher.publishEvent(new ProcessedOrderEvent(this, enriched));
                                })
                                .doOnError(processingError -> {
                                    logger.error("💥 ORDER PROCESSING PIPELINE FAILED for order: {}", order.orderId(), processingError);
                                })
                                .doFinally(sig -> {
                                    logger.info("🔓 RELEASING lock for order: {}", order.orderId());
                                    lockService.release(order.orderId()).subscribe();
                                });
                    })
                    .doOnError(err -> logger.error("💥 CONSUMER ERROR for message: {}", message, err))
                    .subscribe();
        } catch (IOException e) {
            logger.error("💥 JSON DESERIALIZATION FAILED for message: {}", message, e);
        }
    }
}
