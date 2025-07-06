package com.orderprocessing.orderworker.service;

import com.orderprocessing.orderworker.dto.OrderStatusResponse;
import com.orderprocessing.orderworker.model.FailedMessage;
import com.orderprocessing.orderworker.model.OrderDocument;
import com.orderprocessing.orderworker.model.OrderStatus;
import com.orderprocessing.orderworker.model.ProductDetails;
import com.orderprocessing.orderworker.repository.FailedMessageRepository;
import com.orderprocessing.orderworker.repository.OrderRepository;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import reactor.core.publisher.Mono;
import reactor.test.StepVerifier;

import java.time.LocalDateTime;
import java.util.List;

import static org.mockito.ArgumentMatchers.anyString;
import static org.mockito.Mockito.when;

@ExtendWith(MockitoExtension.class)
class OrderStatusServiceTest {

    @Mock
    private OrderRepository orderRepository;

    @Mock
    private FailedMessageRepository failedMessageRepository;

    @InjectMocks
    private OrderStatusService orderStatusService;

    @BeforeEach
    void setUp() {
        // Setup común si es necesario
    }

    @Test
    @DisplayName("Debería devolver estado COMPLETED para orden exitosa")
    void getOrderStatus_shouldReturnCompletedForSuccessfulOrder() {
        // Given
        String orderId = "test-order-001";
        OrderDocument orderDocument = new OrderDocument(
            "64a1b2c3d4e5f6789012345",
            orderId,
            "customer-1",
            List.of(new ProductDetails("product-1", "Test Product", 10.0)),
            OrderStatus.COMPLETED,
            LocalDateTime.now().minusMinutes(5),
            LocalDateTime.now().minusMinutes(1),
            LocalDateTime.now().minusMinutes(1),
            null,
            0
        );

        when(orderRepository.findByOrderId(orderId))
            .thenReturn(Mono.just(orderDocument));
        when(failedMessageRepository.getFailedMessage(anyString()))
            .thenReturn(Mono.empty());

        // When
        Mono<OrderStatusResponse> result = orderStatusService.getOrderStatus(orderId);

        // Then
        StepVerifier.create(result)
            .expectNextMatches(response -> 
                response.orderId().equals(orderId) &&
                response.status().equals(OrderStatus.COMPLETED.name()) &&
                response.failureReason() == null &&
                response.retryCount() == 0
            )
            .verifyComplete();
    }

    @Test
    @DisplayName("Debería devolver estado FAILED para orden fallida")
    void getOrderStatus_shouldReturnFailedForFailedOrder() {
        // Given
        String orderId = "test-order-002";
        FailedMessage failedMessage = new FailedMessage(
            orderId,
            "test content",
            "Customer not found"
        );
        failedMessage.setRetryCount(3);
        failedMessage.setNextRetryAt(LocalDateTime.now().minusMinutes(1)); // Set past time for FAILED status

        when(orderRepository.findByOrderId(orderId))
            .thenReturn(Mono.empty());
        when(failedMessageRepository.getFailedMessage(orderId))
            .thenReturn(Mono.just(failedMessage));
        when(failedMessageRepository.getFailedMessage("order_" + orderId))
            .thenReturn(Mono.empty());

        // When
        Mono<OrderStatusResponse> result = orderStatusService.getOrderStatus(orderId);

        // Then
        StepVerifier.create(result)
            .expectNextMatches(response -> 
                response.orderId().equals(orderId) &&
                response.status().equals(OrderStatus.FAILED.name()) &&
                response.failureReason().equals("Customer not found") &&
                response.retryCount() == 3
            )
            .verifyComplete();
    }

    @Test
    @DisplayName("Debería devolver estado RETRYING para orden en reintento")
    void getOrderStatus_shouldReturnRetryingForRetryingOrder() {
        // Given
        String orderId = "test-order-003";
        FailedMessage failedMessage = new FailedMessage(
            orderId,
            "test content",
            "Connection timeout"
        );
        failedMessage.setRetryCount(1);
        failedMessage.setNextRetryAt(LocalDateTime.now().plusMinutes(5));

        when(orderRepository.findByOrderId(orderId))
            .thenReturn(Mono.empty());
        when(failedMessageRepository.getFailedMessage(orderId))
            .thenReturn(Mono.just(failedMessage));
        when(failedMessageRepository.getFailedMessage("order_" + orderId))
            .thenReturn(Mono.empty());

        // When
        Mono<OrderStatusResponse> result = orderStatusService.getOrderStatus(orderId);

        // Then
        StepVerifier.create(result)
            .expectNextMatches(response -> 
                response.orderId().equals(orderId) &&
                response.status().equals(OrderStatus.RETRYING.name()) &&
                response.failureReason().equals("Connection timeout") &&
                response.retryCount() == 1 &&
                response.nextRetryAt() != null
            )
            .verifyComplete();
    }

    @Test
    @DisplayName("Debería devolver estado PROCESSING para orden desconocida")
    void getOrderStatus_shouldReturnProcessingForUnknownOrder() {
        // Given
        String orderId = "unknown-order";

        when(orderRepository.findByOrderId(orderId))
            .thenReturn(Mono.empty());
        when(failedMessageRepository.getFailedMessage(orderId))
            .thenReturn(Mono.empty());
        when(failedMessageRepository.getFailedMessage("order_" + orderId))
            .thenReturn(Mono.empty());

        // When
        Mono<OrderStatusResponse> result = orderStatusService.getOrderStatus(orderId);

        // Then
        StepVerifier.create(result)
            .expectNextMatches(response -> 
                response.orderId().equals(orderId) &&
                response.status().equals(OrderStatus.PROCESSING.name()) &&
                response.retryCount() == 0
            )
            .verifyComplete();
    }
}