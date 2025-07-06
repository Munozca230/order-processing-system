package com.orderprocessing.orderworker.controller;

import com.orderprocessing.orderworker.dto.OrderStatusResponse;
import com.orderprocessing.orderworker.model.OrderStatus;
import com.orderprocessing.orderworker.service.OrderStatusService;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.springframework.test.web.reactive.server.WebTestClient;
import reactor.core.publisher.Mono;

import java.time.LocalDateTime;

import static org.mockito.ArgumentMatchers.anyString;
import static org.mockito.Mockito.when;

@ExtendWith(MockitoExtension.class)
class OrderStatusControllerTest {

    @Mock
    private OrderStatusService orderStatusService;

    @InjectMocks
    private OrderStatusController orderStatusController;

    private WebTestClient webTestClient;

    @BeforeEach
    void setUp() {
        webTestClient = WebTestClient.bindToController(orderStatusController).build();
    }

    @Test
    @DisplayName("Debería devolver el estado de una orden completada")
    void getOrderStatus_shouldReturnCompletedOrder() {
        // Given
        String orderId = "test-order-001";
        OrderStatusResponse expectedResponse = new OrderStatusResponse(
            orderId,
            OrderStatus.COMPLETED.name(),
            LocalDateTime.now().minusMinutes(5),
            LocalDateTime.now().minusMinutes(1),
            LocalDateTime.now().minusMinutes(1),
            null,
            0,
            null
        );

        when(orderStatusService.getOrderStatus(orderId))
            .thenReturn(Mono.just(expectedResponse));

        // When & Then
        webTestClient.get()
            .uri("/api/orders/{orderId}/status", orderId)
            .exchange()
            .expectStatus().isOk()
            .expectBody(OrderStatusResponse.class)
            .value(response -> {
                assert response.orderId().equals(orderId);
                assert response.status().equals(OrderStatus.COMPLETED.name());
                assert response.failureReason() == null;
                assert response.retryCount() == 0;
            });
    }

    @Test
    @DisplayName("Debería devolver el estado de una orden fallida")
    void getOrderStatus_shouldReturnFailedOrder() {
        // Given
        String orderId = "test-order-002";
        OrderStatusResponse expectedResponse = new OrderStatusResponse(
            orderId,
            OrderStatus.FAILED.name(),
            LocalDateTime.now().minusMinutes(5),
            null,
            LocalDateTime.now().minusMinutes(1),
            "Customer not found",
            3,
            null
        );

        when(orderStatusService.getOrderStatus(orderId))
            .thenReturn(Mono.just(expectedResponse));

        // When & Then
        webTestClient.get()
            .uri("/api/orders/{orderId}/status", orderId)
            .exchange()
            .expectStatus().isOk()
            .expectBody(OrderStatusResponse.class)
            .value(response -> {
                assert response.orderId().equals(orderId);
                assert response.status().equals(OrderStatus.FAILED.name());
                assert response.failureReason().equals("Customer not found");
                assert response.retryCount() == 3;
            });
    }

    @Test
    @DisplayName("Debería devolver el estado de una orden en procesamiento")
    void getOrderStatus_shouldReturnProcessingOrder() {
        // Given
        String orderId = "test-order-003";
        OrderStatusResponse expectedResponse = OrderStatusResponse.processing(orderId);

        when(orderStatusService.getOrderStatus(orderId))
            .thenReturn(Mono.just(expectedResponse));

        // When & Then
        webTestClient.get()
            .uri("/api/orders/{orderId}/status", orderId)
            .exchange()
            .expectStatus().isOk()
            .expectBody(OrderStatusResponse.class)
            .value(response -> {
                assert response.orderId().equals(orderId);
                assert response.status().equals(OrderStatus.PROCESSING.name());
                assert response.retryCount() == 0;
            });
    }

    @Test
    @DisplayName("Debería devolver 404 cuando no se encuentra la orden")
    void getOrderStatus_shouldReturn404ForNotFoundOrder() {
        // Given
        String orderId = "non-existent-order";

        when(orderStatusService.getOrderStatus(orderId))
            .thenReturn(Mono.empty());

        // When & Then
        webTestClient.get()
            .uri("/api/orders/{orderId}/status", orderId)
            .exchange()
            .expectStatus().isNotFound();
    }
}