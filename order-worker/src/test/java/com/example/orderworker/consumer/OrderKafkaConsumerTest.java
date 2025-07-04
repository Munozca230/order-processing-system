package com.example.orderworker.consumer;

import com.example.orderworker.model.OrderMessage;
import com.example.orderworker.service.EnrichmentService;
import com.example.orderworker.service.ValidationService;
import org.springframework.context.ApplicationEventPublisher;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.mockito.Mockito;
import reactor.core.publisher.Mono;

import static org.junit.jupiter.api.Assertions.assertDoesNotThrow;

class OrderKafkaConsumerTest {

    @Test
    @DisplayName("El consumidor debería procesar el mensaje sin lanzar excepciones")
    void consume_shouldNotThrow() {
        EnrichmentService enrichment = Mockito.mock(EnrichmentService.class);
        ValidationService validation = Mockito.mock(ValidationService.class);
        ApplicationEventPublisher publisher = Mockito.mock(ApplicationEventPublisher.class);
        OrderKafkaConsumer consumer = new OrderKafkaConsumer(enrichment, validation, publisher);

        Mockito.when(enrichment.enrich(Mockito.any(OrderMessage.class)))
                .thenReturn(Mono.empty());

        assertDoesNotThrow(() -> consumer.consume("{\"orderId\":\"test-1\",\"customerId\":\"c1\",\"products\":[\"p1\"]}"));
    }
}
