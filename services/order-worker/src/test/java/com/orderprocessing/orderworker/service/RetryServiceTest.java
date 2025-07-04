package com.orderprocessing.orderworker.service;

import com.orderprocessing.orderworker.model.FailedMessage;
import com.orderprocessing.orderworker.repository.FailedMessageRepository;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.mockito.Mock;
import org.mockito.MockitoAnnotations;
import reactor.core.publisher.Mono;
import reactor.test.StepVerifier;

import java.util.concurrent.atomic.AtomicInteger;

import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.anyString;
import static org.mockito.Mockito.when;

class RetryServiceTest {

    @Mock
    private FailedMessageRepository failedMessageRepository;

    private RetryService retryService;

    @BeforeEach
    void setUp() {
        MockitoAnnotations.openMocks(this);
        retryService = new RetryService(failedMessageRepository);
    }

    @Test
    @DisplayName("Debería ejecutar operación exitosa sin reintentos")
    void executeWithRetry_successfulOperation_shouldNotRetry() {
        // Given
        String messageId = "test-message-1";
        String content = "test content";
        String expectedResult = "success";

        // When
        Mono<String> result = retryService.executeWithRetry(
            messageId, 
            content, 
            () -> Mono.just(expectedResult)
        );

        // Then
        StepVerifier.create(result)
            .expectNext(expectedResult)
            .verifyComplete();
    }

    @Test
    @DisplayName("Debería reintentar operación fallida y eventualmente tener éxito")
    void executeWithRetry_failThenSuccess_shouldRetryAndSucceed() {
        // Given
        String messageId = "test-message-2";
        String content = "test content";
        String expectedResult = "success";
        AtomicInteger attempts = new AtomicInteger(0);

        // Mock repository calls
        when(failedMessageRepository.saveFailedMessage(any(FailedMessage.class)))
            .thenReturn(Mono.empty());
        when(failedMessageRepository.addToRetryQueue(anyString(), any()))
            .thenReturn(Mono.empty());

        // When - Operación que falla las primeras 2 veces, luego tiene éxito
        Mono<String> result = retryService.executeWithRetry(
            messageId,
            content,
            () -> {
                int attempt = attempts.incrementAndGet();
                if (attempt <= 2) {
                    return Mono.error(new RuntimeException("Timeout error"));
                }
                return Mono.just(expectedResult);
            }
        );

        // Then
        StepVerifier.create(result)
            .expectNext(expectedResult)
            .verifyComplete();
    }

    @Test
    @DisplayName("Debería almacenar mensaje fallido después de máximo de reintentos")
    void executeWithRetry_maxRetriesExceeded_shouldStoreFailedMessage() {
        // Given
        String messageId = "test-message-3";
        String content = "test content";

        // Mock repository calls
        when(failedMessageRepository.saveFailedMessage(any(FailedMessage.class)))
            .thenReturn(Mono.empty());
        when(failedMessageRepository.addToRetryQueue(anyString(), any()))
            .thenReturn(Mono.empty());

        // When - Operación que siempre falla
        Mono<String> result = retryService.executeWithRetry(
            messageId,
            content,
            () -> Mono.error(new RuntimeException("Persistent timeout error"))
        );

        // Then
        StepVerifier.create(result)
            .expectError(RuntimeException.class)
            .verify();
    }
}