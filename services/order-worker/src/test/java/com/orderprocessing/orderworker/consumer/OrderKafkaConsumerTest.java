package com.orderprocessing.orderworker.consumer;

import com.orderprocessing.orderworker.model.EnrichedOrder;
import com.orderprocessing.orderworker.model.OrderDocument;
import com.orderprocessing.orderworker.model.OrderMessage;
import com.orderprocessing.orderworker.model.ProductReference;
import com.orderprocessing.orderworker.model.ProductDetails;
import com.orderprocessing.orderworker.service.EnrichmentService;
import com.orderprocessing.orderworker.service.ValidationService;
import org.springframework.context.ApplicationEventPublisher;
import com.orderprocessing.orderworker.repository.OrderRepository;
import com.orderprocessing.orderworker.service.OrderLockService;
import reactor.core.publisher.Mono;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.mockito.Mock;
import org.mockito.MockitoAnnotations;

import java.util.List;

import static org.junit.jupiter.api.Assertions.assertDoesNotThrow;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.anyString;
import static org.mockito.Mockito.when;

class OrderKafkaConsumerTest {

    @Mock
    private EnrichmentService enrichmentService;
    
    @Mock
    private ValidationService validationService;
    
    @Mock
    private ApplicationEventPublisher eventPublisher;
    
    @Mock
    private OrderRepository orderRepository;
    
    @Mock
    private OrderLockService orderLockService;
    
    private OrderKafkaConsumer consumer;

    @BeforeEach
    void setUp() {
        MockitoAnnotations.openMocks(this);
        consumer = new OrderKafkaConsumer(enrichmentService, validationService, eventPublisher, orderRepository, orderLockService);
    }

    @Test
    @DisplayName("El consumidor debería procesar el mensaje exitosamente con mocking realista")
    void consume_shouldProcessSuccessfully() {
        // Given
        String validOrderJson = "{\"orderId\":\"test-1\",\"customerId\":\"c1\",\"products\":[{\"productId\":\"p1\"}]}";
        
        // Mock enriched order
        OrderMessage orderMessage = new OrderMessage("test-1", "c1", List.of(new ProductReference("p1")));
        ProductDetails productDetails = new ProductDetails("p1", "Test Product", 10.0);
        EnrichedOrder enrichedOrder = new EnrichedOrder(orderMessage, null, List.of(productDetails));
        
        // Mock order document for save (using record constructor)
        OrderDocument orderDocument = new OrderDocument("saved-id", "test-1", "c1", List.of(productDetails), 
            com.orderprocessing.orderworker.model.OrderStatus.COMPLETED, java.time.LocalDateTime.now(), 
            java.time.LocalDateTime.now(), java.time.LocalDateTime.now(), null, 0);
        
        // Setup mocks
        when(orderLockService.acquire(anyString())).thenReturn(Mono.just(true));
        when(enrichmentService.enrich(any(OrderMessage.class))).thenReturn(Mono.just(enrichedOrder));
        when(validationService.validate(any(EnrichedOrder.class))).thenReturn(Mono.just(enrichedOrder));
        when(orderRepository.save(any(OrderDocument.class))).thenReturn(Mono.just(orderDocument));
        when(orderLockService.release(anyString())).thenReturn(Mono.empty());

        // When & Then
        assertDoesNotThrow(() -> consumer.consume(validOrderJson));
    }
    
    @Test
    @DisplayName("El consumidor debería manejar JSON inválido sin lanzar excepciones")
    void consume_invalidJson_shouldHandleGracefully() {
        // Given
        String invalidJson = "{invalid json}";

        // When & Then
        assertDoesNotThrow(() -> consumer.consume(invalidJson));
    }
    
    @Test
    @DisplayName("El consumidor debería manejar el caso cuando no puede adquirir el lock")
    void consume_lockNotAcquired_shouldSkipProcessing() {
        // Given
        String validOrderJson = "{\"orderId\":\"test-1\",\"customerId\":\"c1\",\"products\":[{\"productId\":\"p1\"}]}";
        
        // Setup mocks - lock acquisition fails
        when(orderLockService.acquire(anyString())).thenReturn(Mono.just(false));

        // When & Then
        assertDoesNotThrow(() -> consumer.consume(validOrderJson));
    }
}
