package com.orderprocessing.orderworker.service;

import com.orderprocessing.orderworker.model.FailedMessage;
import com.orderprocessing.orderworker.repository.FailedMessageRepository;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.scheduling.annotation.Scheduled;
import org.springframework.stereotype.Service;
import reactor.core.publisher.Mono;
import reactor.util.retry.Retry;

import java.time.Duration;
import java.util.function.Supplier;

@Service
public class RetryService {

    private static final Logger logger = LoggerFactory.getLogger(RetryService.class);
    
    private final FailedMessageRepository failedMessageRepository;
    
    @Value("${app.retry.max-attempts:5}")
    private int maxRetryAttempts;
    
    @Value("${app.retry.initial-delay:1}")
    private int initialDelaySeconds;

    public RetryService(FailedMessageRepository failedMessageRepository) {
        this.failedMessageRepository = failedMessageRepository;
    }

    public <T> Mono<T> executeWithRetry(String messageId, String messageContent, Supplier<Mono<T>> operation) {
        return operation.get()
            .retryWhen(Retry.backoff(maxRetryAttempts, Duration.ofSeconds(initialDelaySeconds))
                .filter(throwable -> isRetryableError(throwable))
                .doBeforeRetry(retrySignal -> {
                    logger.warn("Retrying operation for message {}, attempt {}: {}", 
                        messageId, retrySignal.totalRetries() + 1, retrySignal.failure().getMessage());
                }))
            .onErrorResume(throwable -> {
                logger.error("Max retries exceeded for message {}, storing as failed", messageId, throwable);
                return storeFailedMessage(messageId, messageContent, throwable.getMessage())
                    .then(Mono.error(throwable));
            });
    }

    public Mono<Void> storeFailedMessage(String messageId, String content, String errorMessage) {
        FailedMessage failedMessage = new FailedMessage(messageId, content, errorMessage);
        
        return failedMessageRepository.saveFailedMessage(failedMessage)
            .then(failedMessageRepository.addToRetryQueue(messageId, failedMessage.getNextRetryAt()))
            .doOnSuccess(unused -> logger.info("Stored failed message {} for retry", messageId))
            .doOnError(error -> logger.error("Failed to store failed message {}", messageId, error));
    }

    @Scheduled(fixedDelay = 30000) // Check every 30 seconds
    public void processRetryQueue() {
        logger.debug("Processing retry queue");
        
        failedMessageRepository.getMessagesReadyForRetry()
            .flatMap(this::processRetryMessage)
            .subscribe(
                unused -> {},
                error -> logger.error("Error processing retry queue", error)
            );
    }

    private Mono<Void> processRetryMessage(String messageId) {
        return failedMessageRepository.getFailedMessage(messageId)
            .filter(failedMessage -> failedMessage.getRetryCount() < maxRetryAttempts)
            .flatMap(failedMessage -> {
                logger.info("Retrying failed message {}, attempt {}", messageId, failedMessage.getRetryCount() + 1);
                
                failedMessage.incrementRetry();
                
                return failedMessageRepository.saveFailedMessage(failedMessage)
                    .then(failedMessageRepository.addToRetryQueue(messageId, failedMessage.getNextRetryAt()))
                    .then(failedMessageRepository.removeFromRetryQueue(messageId));
            })
            .switchIfEmpty(
                // Max retries exceeded, move to dead letter queue
                failedMessageRepository.removeFromRetryQueue(messageId)
                    .doOnSuccess(unused -> logger.warn("Message {} moved to dead letter queue after max retries", messageId))
            );
    }

    private boolean isRetryableError(Throwable throwable) {
        // Retry on network errors, timeouts, 5xx HTTP errors
        String message = throwable.getMessage().toLowerCase();
        return message.contains("timeout") || 
               message.contains("connection") || 
               message.contains("5") || // 5xx errors
               throwable instanceof java.net.ConnectException ||
               throwable instanceof java.util.concurrent.TimeoutException;
    }
}