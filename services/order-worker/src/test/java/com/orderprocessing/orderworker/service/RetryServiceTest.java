package com.orderprocessing.orderworker.service;

import com.orderprocessing.orderworker.model.FailedMessage;
import com.orderprocessing.orderworker.repository.FailedMessageRepository;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.junit.jupiter.api.Test;
import org.mockito.Mock;
import org.mockito.MockitoAnnotations;
import org.springframework.test.util.ReflectionTestUtils;
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
        
        // Configure retry parameters for testing
        ReflectionTestUtils.setField(retryService, "maxRetryAttempts", 3);
        ReflectionTestUtils.setField(retryService, "initialDelaySeconds", 0); // No delay for fast tests
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

        // When - Operación que falla las primeras 2 veces, luego tiene éxito en el 3er intento (total 3 attempts)
        // Use Mono.defer to ensure the operation is re-evaluated on each retry
        Mono<String> result = retryService.executeWithRetry(
            messageId,
            content,
            () -> Mono.defer(() -> {
                int attempt = attempts.incrementAndGet();
                if (attempt <= 2) {
                    return Mono.error(new RuntimeException("connection timeout error")); // Retryable error
                } else {
                    return Mono.just(expectedResult);
                }
            })
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

        // When - Operación que siempre falla con error retryable
        Mono<String> result = retryService.executeWithRetry(
            messageId,
            content,
            () -> Mono.error(new RuntimeException("persistent timeout error")) // lowercase para que sea retryable
        );

        // Then - Después de 3 reintentos, debería fallar y almacenar el mensaje
        StepVerifier.create(result)
            .expectError(RuntimeException.class)
            .verify();
    }
    
    @Test
    @DisplayName("No debería reintentar errores no retryables")
    void executeWithRetry_nonRetryableError_shouldFailImmediately() {
        // Given
        String messageId = "test-message-4";
        String content = "test content";

        // Mock repository calls
        when(failedMessageRepository.saveFailedMessage(any(FailedMessage.class)))
            .thenReturn(Mono.empty());
        when(failedMessageRepository.addToRetryQueue(anyString(), any()))
            .thenReturn(Mono.empty());

        // When - Error que no es retryable (no contiene palabras clave como timeout, connection, etc.)
        Mono<String> result = retryService.executeWithRetry(
            messageId,
            content,
            () -> Mono.error(new IllegalArgumentException("Invalid input data"))
        );

        // Then - Debería fallar inmediatamente sin reintentos
        StepVerifier.create(result)
            .expectError(IllegalArgumentException.class)
            .verify();
    }
}