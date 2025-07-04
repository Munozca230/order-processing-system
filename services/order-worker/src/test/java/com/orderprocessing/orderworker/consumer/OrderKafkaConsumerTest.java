package com.orderprocessing.orderworker.consumer;

import com.orderprocessing.orderworker.model.OrderMessage;
import com.orderprocessing.orderworker.service.EnrichmentService;
import com.orderprocessing.orderworker.service.ValidationService;
import org.springframework.context.ApplicationEventPublisher;
import com.orderprocessing.orderworker.repository.OrderRepository;
import com.orderprocessing.orderworker.service.OrderLockService;
import reactor.core.publisher.Mono;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.mockito.Mockito;

import static org.junit.jupiter.api.Assertions.assertDoesNotThrow;

class OrderKafkaConsumerTest {

    @Test
    @DisplayName("El consumidor debería procesar el mensaje sin lanzar excepciones")
    void consume_shouldNotThrow() {
        EnrichmentService enrichment = Mockito.mock(EnrichmentService.class);
        ValidationService validation = Mockito.mock(ValidationService.class);
        ApplicationEventPublisher publisher = Mockito.mock(ApplicationEventPublisher.class);
        OrderRepository repo = Mockito.mock(OrderRepository.class);
        OrderLockService lockService = Mockito.mock(OrderLockService.class);
        Mockito.when(repo.save(Mockito.any())).thenReturn(Mono.empty());
        Mockito.when(lockService.acquire(Mockito.anyString())).thenReturn(Mono.just(true));
        Mockito.when(lockService.release(Mockito.anyString())).thenReturn(Mono.empty());
        OrderKafkaConsumer consumer = new OrderKafkaConsumer(enrichment, validation, publisher, repo, lockService);

        Mockito.when(enrichment.enrich(Mockito.any(OrderMessage.class)))
                .thenReturn(Mono.empty());

        assertDoesNotThrow(() -> consumer.consume("{\"orderId\":\"test-1\",\"customerId\":\"c1\",\"products\":[\"p1\"]}"));
    }
}
